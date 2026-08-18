package world

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Resolver interface {
	Resolve(context.Context, Resource) (any, error)
}

type Graph struct {
	definition *Definition
	resolver   Resolver
	rootRefs   map[string]string
}

func NewGraph(definition *Definition, resolver Resolver) *Graph {
	return &Graph{definition: definition, resolver: resolver}
}

// NewGraphWithRootReferences creates an isolated graph whose first-level
// handles are supplied by the caller. Evaluation runners use it to hand the
// same opaque references to a fresh runtime and its fresh Phoenix process.
func NewGraphWithRootReferences(definition *Definition, resolver Resolver, references map[string]string) (*Graph, error) {
	if definition == nil {
		return nil, fmt.Errorf("world definition is required")
	}
	if len(references) != len(definition.Roots) {
		return nil, fmt.Errorf("one opaque reference is required for every root")
	}
	copyRefs := make(map[string]string, len(references))
	seen := make(map[string]struct{}, len(references))
	for _, root := range definition.Roots {
		ref, exists := references[root.Name]
		if !exists || !strings.HasPrefix(ref, "h_") || len(ref) < 18 {
			return nil, fmt.Errorf("root %q has no valid opaque reference", root.Name)
		}
		if _, duplicate := seen[ref]; duplicate {
			return nil, fmt.Errorf("root reference %q is duplicated", ref)
		}
		seen[ref] = struct{}{}
		copyRefs[root.Name] = ref
	}
	return &Graph{definition: definition, resolver: resolver, rootRefs: copyRefs}, nil
}

type Session struct {
	id         string
	definition *Definition
	resolver   Resolver
	mu         sync.RWMutex
	reachable  map[string]*reachableHandle
	roots      map[string]string
}

type reachableHandle struct {
	Handle
	resource Resource
}

type Handle struct {
	Ref   string `json:"ref"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type RootHandle struct {
	Name string `json:"name"`
	Handle
}

type AccessStatus string

const (
	AccessReady  AccessStatus = "ready"
	AccessAbsent AccessStatus = "absent"
	AccessStale  AccessStatus = "stale"
	AccessFailed AccessStatus = "fail"
)

type Problem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type LiveState struct {
	Digest string `json:"digest"`
	Value  any    `json:"value"`
}

type Target struct {
	Type     string
	Label    string
	Resource Resource
	Verb     Verb
}

type Access struct {
	Status  AccessStatus
	Target  *Target
	State   *LiveState
	Problem *Problem
}

type GrantRequest struct {
	Type     string
	Label    string
	Resource Resource
}

type Mutation struct {
	Grants  []GrantRequest
	Revokes []string
}

type Delta struct {
	Grant  []Handle `json:"grant"`
	Revoke []string `json:"revoke"`
}

type ReachableMatch struct {
	Handle string `json:"handle"`
	Type   string `json:"type"`
	Label  string `json:"label"`
	Verb   string `json:"verb"`
}

var ErrAbsent = errors.New("act is not reachable")

func (graph *Graph) StartSession() (*Session, []RootHandle, error) {
	if graph == nil || graph.definition == nil {
		return nil, nil, fmt.Errorf("world definition is required")
	}
	if graph.resolver == nil {
		return nil, nil, fmt.Errorf("live resource resolver is required")
	}
	id, err := opaqueID("s_")
	if err != nil {
		return nil, nil, fmt.Errorf("create session id: %w", err)
	}
	session := &Session{
		id:         id,
		definition: graph.definition,
		resolver:   graph.resolver,
		reachable:  make(map[string]*reachableHandle, len(graph.definition.Roots)),
		roots:      make(map[string]string, len(graph.definition.Roots)),
	}
	roots := make([]RootHandle, 0, len(graph.definition.Roots))
	for _, root := range graph.definition.Roots {
		ref := graph.rootRefs[root.Name]
		if ref == "" {
			ref, err = opaqueID("h_")
			if err != nil {
				return nil, nil, fmt.Errorf("create root handle: %w", err)
			}
		}
		handle := Handle{Ref: ref, Type: root.Type, Label: root.Label}
		session.reachable[ref] = &reachableHandle{Handle: handle, resource: root.Resource}
		session.roots[root.Name] = ref
		roots = append(roots, RootHandle{Name: root.Name, Handle: handle})
	}
	return session, roots, nil
}

// RootHandle resolves an authored root name to this session's opaque handle.
func (session *Session) RootHandle(name string) (Handle, bool) {
	session.mu.RLock()
	defer session.mu.RUnlock()
	ref, exists := session.roots[name]
	if !exists {
		return Handle{}, false
	}
	handle, reachable := session.reachable[ref]
	if !reachable {
		return Handle{}, false
	}
	return handle.Handle, true
}

// ReachableHandle returns public identity only for a currently live handle.
func (session *Session) ReachableHandle(ref string) (Handle, bool) {
	session.mu.RLock()
	defer session.mu.RUnlock()
	handle, exists := session.reachable[ref]
	if !exists {
		return Handle{}, false
	}
	return handle.Handle, true
}

func (session *Session) ID() string {
	return session.id
}

func (session *Session) ReachableCount() int {
	session.mu.RLock()
	defer session.mu.RUnlock()
	return len(session.reachable)
}

// FindReachable searches only handles and attached verbs already reachable in
// this session. It is intentionally not a global registry search.
func (session *Session) FindReachable(query string) []ReachableMatch {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return []ReachableMatch{}
	}
	session.mu.RLock()
	defer session.mu.RUnlock()
	var matches []ReachableMatch
	for ref, handle := range session.reachable {
		verbs := session.definition.HandleTypes[handle.Type].Verbs
		for verb := range verbs {
			text := strings.ToLower(handle.Type + " " + handle.Label + " " + verb)
			if strings.Contains(text, query) {
				matches = append(matches, ReachableMatch{Handle: ref, Type: handle.Type, Label: handle.Label, Verb: verb})
			}
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Handle == matches[j].Handle {
			return matches[i].Verb < matches[j].Verb
		}
		return matches[i].Handle < matches[j].Handle
	})
	return matches
}

func (session *Session) Prepare(ctx context.Context, ref, verbName, expectedState string) Access {
	session.mu.RLock()
	handle, verb, exists := session.lookup(ref, verbName)
	session.mu.RUnlock()
	if !exists {
		return absentAccess()
	}

	value, err := session.resolver.Resolve(ctx, handle.resource)
	if err != nil {
		return Access{Status: AccessFailed, Problem: &Problem{Code: "resource_unavailable", Message: "live resource could not be resolved"}}
	}
	digest, err := digestValue(value)
	if err != nil {
		return Access{Status: AccessFailed, Problem: &Problem{Code: "state_invalid", Message: "live resource state is not canonical JSON"}}
	}

	session.mu.RLock()
	current := session.reachable[ref]
	stillReachable := current == handle
	session.mu.RUnlock()
	if !stillReachable {
		return absentAccess()
	}
	if expectedState != "" && expectedState != digest {
		return Access{
			Status:  AccessStale,
			Target:  &Target{Type: handle.Type, Label: handle.Label, Resource: handle.resource, Verb: verb},
			State:   &LiveState{Digest: digest, Value: value},
			Problem: &Problem{Code: "stale_state", Message: "live resource state changed"},
		}
	}
	return Access{
		Status: AccessReady,
		Target: &Target{Type: handle.Type, Label: handle.Label, Resource: handle.resource, Verb: verb},
		State:  &LiveState{Digest: digest, Value: value},
	}
}

func (session *Session) ApplyDelta(sourceRef, verbName string, mutation Mutation) (Delta, error) {
	session.mu.Lock()
	defer session.mu.Unlock()
	_, verb, exists := session.lookup(sourceRef, verbName)
	if !exists {
		return Delta{}, ErrAbsent
	}
	if err := session.validateMutation(verbName, verb, mutation); err != nil {
		return Delta{}, err
	}
	return session.applyMutation(mutation)
}

func (session *Session) validateMutation(verbName string, verb Verb, mutation Mutation) error {
	allowedGrants := make(map[string]struct{}, len(verb.Grants))
	for _, handleType := range verb.Grants {
		allowedGrants[handleType] = struct{}{}
	}
	for _, request := range mutation.Grants {
		if _, allowed := allowedGrants[request.Type]; !allowed {
			return fmt.Errorf("%s does not declare grants of handle type %q", verbName, request.Type)
		}
		if request.Label == "" || !validResource(request.Resource) {
			return fmt.Errorf("grant of handle type %q has invalid label or resource", request.Type)
		}
	}
	if len(mutation.Revokes) > 0 && !verb.MayRevoke {
		return fmt.Errorf("%s does not declare handle revocation", verbName)
	}
	seenRevokes := make(map[string]struct{}, len(mutation.Revokes))
	for _, ref := range mutation.Revokes {
		if _, duplicate := seenRevokes[ref]; duplicate {
			return fmt.Errorf("handle %q is revoked more than once", ref)
		}
		seenRevokes[ref] = struct{}{}
		if _, reachable := session.reachable[ref]; !reachable {
			return ErrAbsent
		}
	}
	return nil
}

func (session *Session) applyMutation(mutation Mutation) (Delta, error) {
	granted := make([]*reachableHandle, 0, len(mutation.Grants))
	for _, request := range mutation.Grants {
		ref, err := opaqueID("h_")
		if err != nil {
			return Delta{}, fmt.Errorf("create granted handle: %w", err)
		}
		handle := Handle{Ref: ref, Type: request.Type, Label: request.Label}
		granted = append(granted, &reachableHandle{Handle: handle, resource: request.Resource})
	}

	delta := Delta{Grant: make([]Handle, 0, len(granted)), Revoke: append([]string(nil), mutation.Revokes...)}
	for _, handle := range granted {
		session.reachable[handle.Ref] = handle
		delta.Grant = append(delta.Grant, handle.Handle)
	}
	for _, ref := range mutation.Revokes {
		delete(session.reachable, ref)
	}
	return delta, nil
}

func (session *Session) lookup(ref, verbName string) (*reachableHandle, Verb, bool) {
	handle, exists := session.reachable[ref]
	if !exists {
		return nil, Verb{}, false
	}
	descriptor, exists := session.definition.HandleTypes[handle.Type]
	if !exists {
		return nil, Verb{}, false
	}
	verb, exists := descriptor.Verbs[verbName]
	if !exists {
		return nil, Verb{}, false
	}
	return handle, verb, true
}

func absentAccess() Access {
	return Access{Status: AccessAbsent, Problem: &Problem{Code: "absent", Message: ErrAbsent.Error()}}
}

func validResource(resource Resource) bool {
	if resource.Value == "" {
		return false
	}
	return resource.Kind == "path" || resource.Kind == "process" || resource.Kind == "logical"
}

func opaqueID(prefix string) (string, error) {
	random := make([]byte, 18)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	return prefix + base64.RawURLEncoding.EncodeToString(random), nil
}
