package world

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
)

type Resolver interface {
	Resolve(context.Context, Resource) (any, error)
}

type Graph struct {
	definition *Definition
	resolver   Resolver
}

func NewGraph(definition *Definition, resolver Resolver) *Graph {
	return &Graph{definition: definition, resolver: resolver}
}

type Session struct {
	id         string
	definition *Definition
	resolver   Resolver
	mu         sync.RWMutex
	reachable  map[string]*reachableHandle
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

var ErrAbsent = errors.New("handle is not reachable")

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
	}
	roots := make([]RootHandle, 0, len(graph.definition.Roots))
	for _, root := range graph.definition.Roots {
		ref, err := opaqueID("h_")
		if err != nil {
			return nil, nil, fmt.Errorf("create root handle: %w", err)
		}
		handle := Handle{Ref: ref, Type: root.Type, Label: root.Label}
		session.reachable[ref] = &reachableHandle{Handle: handle, resource: root.Resource}
		roots = append(roots, RootHandle{Name: root.Name, Handle: handle})
	}
	return session, roots, nil
}

func (session *Session) ID() string {
	return session.id
}

func (session *Session) ReachableCount() int {
	session.mu.RLock()
	defer session.mu.RUnlock()
	return len(session.reachable)
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
		return Access{Status: AccessStale, Problem: &Problem{Code: "stale_state", Message: "live resource state changed"}}
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
