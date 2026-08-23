// Package surface exposes Phoenix through one constant Model Context Protocol
// tool and admits every invocation through the same bounded execution path.
package surface

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nesste/phoenix/internal/activate"
	"github.com/nesste/phoenix/internal/episode"
	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/teach"
	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
)

const (
	ToolName        = "act"
	ToolDescription = "Invoke one verb on a live Phoenix handle, or provide intent to receive up to three ready-to-run reachable calls."
	defaultMaxArgs  = 64 * 1024
)

var (
	handlePattern = regexp.MustCompile(`^h_[A-Za-z0-9_-]{16,}$`)
	verbPattern   = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
	digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

var actInputSchema = json.RawMessage(`{
	"type":"object",
	"additionalProperties":false,
	"required":["handle"],
	"properties":{
		"handle":{"type":"string","pattern":"^h_[A-Za-z0-9_-]{16,}$"},
		"verb":{"type":"string","pattern":"^[a-z][a-z0-9_]{0,63}$"},
		"args":{"type":"object"},
		"state":{"type":"string","pattern":"^sha256:[0-9a-f]{64}$"},
		"intent":{"type":"string","minLength":1,"maxLength":800}
	}
}`)

type Input struct {
	Handle string         `json:"handle"`
	Verb   string         `json:"verb"`
	Args   map[string]any `json:"args"`
	State  string         `json:"state,omitempty"`
	Intent string         `json:"intent,omitempty"`
}

type Status string

const (
	StatusOK      Status = "ok"
	StatusFail    Status = "fail"
	StatusRefused Status = "refused"
	StatusAbsent  Status = "absent"
)

type Error struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type HandleGrant struct {
	Ref   string `json:"ref"`
	Type  string `json:"type"`
	Label string `json:"label"`
	State string `json:"state"`
}

type HandleDelta struct {
	Grant  []HandleGrant `json:"grant"`
	Revoke []string      `json:"revoke"`
}

type Call = frontier.Call

type FrontierEntry = frontier.Entry

type Refusal = teach.Refusal

type Envelope struct {
	V          int             `json:"v"`
	WorldBuild string          `json:"world_build"`
	SessionID  string          `json:"session_id"`
	ActID      string          `json:"act_id"`
	Handle     string          `json:"handle"`
	Verb       string          `json:"verb"`
	Status     Status          `json:"status"`
	Result     any             `json:"result"`
	Error      *Error          `json:"error"`
	Text       string          `json:"text"`
	Handles    HandleDelta     `json:"handles"`
	Frontier   []FrontierEntry `json:"frontier"`
	Refusal    *Refusal        `json:"refusal"`
}

type Config struct {
	WorldBuild       string
	Graph            *world.Graph
	Executor         *verb.Executor
	Frontier         *frontier.Engine
	Activation       *activate.Engine
	Teacher          *teach.Engine
	Episodes         EpisodeLog
	Warning          io.Writer
	MaxArgsBytes     int
	AfterAct         func(context.Context, int, Input, Envelope) error
	SuppressFrontier bool
	SuppressTeaching bool
}

type EpisodeLog interface {
	StartEpisode(context.Context, string, string) (string, error)
	BeginAct(context.Context, episode.StartedAct) error
	CompleteAct(context.Context, episode.CompletedAct) error
}

type Admission struct {
	worldBuild       string
	graph            *world.Graph
	executor         *verb.Executor
	frontier         *frontier.Engine
	activation       *activate.Engine
	teacher          *teach.Engine
	episodes         EpisodeLog
	warning          io.Writer
	maxArgsBytes     int
	afterAct         func(context.Context, int, Input, Envelope) error
	suppressFrontier bool
	suppressTeaching bool
}

type Session struct {
	graph     *world.Session
	episodeID string
	mu        sync.Mutex
	suggested map[string]episode.SuggestionLink
	stateful  map[string]statefulSuggestion
	pending   []FrontierEntry
	executed  int
}

type statefulSuggestion struct {
	state     string
	link      episode.SuggestionLink
	ambiguous bool
}

func New(config Config) (*Admission, error) {
	if !digestPattern.MatchString(config.WorldBuild) {
		return nil, fmt.Errorf("world build must be a sha256 digest")
	}
	if config.Graph == nil {
		return nil, fmt.Errorf("world graph is required")
	}
	if config.Executor == nil {
		return nil, fmt.Errorf("verb executor is required")
	}
	if config.Frontier == nil {
		return nil, fmt.Errorf("frontier engine is required")
	}
	if config.Activation == nil {
		return nil, fmt.Errorf("activation engine is required")
	}
	if config.Teacher == nil {
		return nil, fmt.Errorf("teaching engine is required")
	}
	if config.Episodes != nil && config.Warning == nil {
		return nil, fmt.Errorf("warning writer is required with episode logging")
	}
	if config.MaxArgsBytes <= 0 {
		config.MaxArgsBytes = defaultMaxArgs
	}
	return &Admission{
		worldBuild: config.WorldBuild, graph: config.Graph,
		executor: config.Executor, frontier: config.Frontier, activation: config.Activation,
		teacher: config.Teacher, episodes: config.Episodes, warning: config.Warning,
		maxArgsBytes: config.MaxArgsBytes, afterAct: config.AfterAct,
		suppressFrontier: config.SuppressFrontier, suppressTeaching: config.SuppressTeaching,
	}, nil
}

func (admission *Admission) StartSession() (*Session, []world.RootHandle, error) {
	graphSession, roots, err := admission.graph.StartSession()
	if err != nil {
		return nil, nil, err
	}
	session := &Session{
		graph: graphSession, suggested: make(map[string]episode.SuggestionLink),
		stateful: make(map[string]statefulSuggestion),
	}
	if admission.episodes != nil {
		episodeID, startErr := admission.episodes.StartEpisode(context.Background(), graphSession.ID(), admission.worldBuild)
		if startErr != nil {
			admission.warn("episode logging disabled for session: %v", startErr)
		} else {
			session.episodeID = episodeID
		}
	}
	return session, roots, nil
}

func (session *Session) ID() string {
	if session == nil || session.graph == nil {
		return ""
	}
	return session.graph.ID()
}

func (session *Session) EpisodeID() string {
	if session == nil {
		return ""
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	return session.episodeID
}

func (admission *Admission) Act(ctx context.Context, session *Session, input Input) Envelope {
	base := admission.baseEnvelope(session, input)
	if input.Intent != "" && (input.Verb != "" || input.Args != nil || input.State != "") {
		return failure(base, "invalid_act", "intent and executable fields are mutually exclusive", nil)
	}
	if session != nil && input.Intent == "" && input.State == "" {
		input = session.completeSuggestedState(input)
	}
	if session != nil && input.Intent == "" {
		session.clearPending()
	}
	base = admission.baseEnvelope(session, input)
	if session == nil || session.graph == nil {
		return failure(base, "invalid_session", "session is unavailable", nil)
	}
	if input.Intent != "" {
		return admission.orient(ctx, session, input, base)
	}
	if failed := admission.validateExecutionInput(base, input); failed != nil {
		return *failed
	}
	return admission.runAct(ctx, session, input, base)
}

func (admission *Admission) validateExecutionInput(base Envelope, input Input) *Envelope {
	encodedArgs, err := json.Marshal(input.Args)
	if err != nil || input.Args == nil {
		failed := failure(base, "invalid_arguments", "act arguments must be a JSON object", nil)
		return &failed
	}
	if len(encodedArgs) > admission.maxArgsBytes {
		failed := failure(base, "arguments_too_large", "act arguments exceeded the size limit", map[string]any{
			"limit_bytes": admission.maxArgsBytes,
		})
		return &failed
	}
	return nil
}

func (admission *Admission) runAct(ctx context.Context, session *Session, input Input, base Envelope) (output Envelope) {
	logContext := context.WithoutCancel(ctx)
	if admission.beginEpisodeAct(logContext, session, base, input) {
		defer func() {
			admission.completeEpisodeAct(logContext, session, base.ActID, output)
			session.rememberSuggestions(output)
		}()
	} else {
		defer func() { session.rememberSuggestions(output) }()
	}
	defer func() {
		index := session.nextExecutedIndex()
		if admission.afterAct != nil {
			if err := admission.afterAct(context.WithoutCancel(ctx), index, input, output); err != nil {
				output = failure(base, "state_event_failed", "configured state event could not be applied", nil)
			}
		}
	}()

	return admission.executePrepared(ctx, session, input, base)
}

func (admission *Admission) executePrepared(ctx context.Context, session *Session, input Input, base Envelope) Envelope {
	access := session.graph.Prepare(ctx, input.Handle, input.Verb, input.State)
	switch access.Status {
	case world.AccessAbsent:
		return absent(base)
	case world.AccessStale:
		observation := teach.Observation{
			HandleType: access.Target.Type, Handle: input.Handle, Verb: input.Verb,
			Args: input.Args, State: access.State,
			Failure: map[string]any{"code": access.Problem.Code, "message": access.Problem.Message},
		}
		if shaped := admission.evaluateRefusal(base, session, observation); shaped != nil {
			return *shaped
		}
		return failure(base, access.Problem.Code, access.Problem.Message, nil)
	case world.AccessFailed:
		return failure(base, access.Problem.Code, access.Problem.Message, nil)
	case world.AccessReady:
		return admission.executeReady(ctx, session, input, base, access)
	default:
		return failure(base, "execution_failed", "verb execution failed", nil)
	}
}

func (admission *Admission) executeReady(ctx context.Context, session *Session, input Input, base Envelope, access world.Access) Envelope {
	observation := teach.Observation{
		HandleType: access.Target.Type, Handle: input.Handle, Verb: input.Verb,
		Args: input.Args, State: access.State,
	}
	if shaped := admission.evaluateRefusal(base, session, observation); shaped != nil {
		return *shaped
	}
	result := admission.executor.Execute(ctx, verb.Request{
		SessionID: session.ID(), HandleType: access.Target.Type, Verb: input.Verb,
		Resource: access.Target.Resource, State: access.State, Args: input.Args,
		Reachable: session.graph,
	})
	if result.Status == verb.StatusFail {
		return admission.failedExecution(base, session, input, access, observation, result)
	}
	base.Status = StatusOK
	base.Result = result.Value
	base.Text = renderSuccess(access.Target.Type, input.Verb, result.Value)
	if !admission.suppressFrontier {
		base.Frontier = admission.frontier.Compute(session.graph, frontier.Observation{
			HandleType: access.Target.Type, Handle: input.Handle, Verb: input.Verb,
			Status: string(StatusOK), Result: result.Value, State: access.State,
		})
	}
	return base
}

func (admission *Admission) failedExecution(base Envelope, session *Session, input Input, access world.Access, observation teach.Observation, result verb.Result) Envelope {
	observation.Failure = failureFeatures(result.Error)
	if shaped := admission.evaluateRefusal(base, session, observation); shaped != nil {
		return *shaped
	}
	failed := failure(base, result.Error.Code, result.Error.Message, result.Error.Details)
	if !admission.suppressFrontier {
		failed.Frontier = admission.frontier.Compute(session.graph, frontier.Observation{
			HandleType: access.Target.Type, Handle: input.Handle, Verb: input.Verb,
			Status: string(StatusFail), Result: failureFeatures(result.Error), State: access.State,
		})
	}
	return failed
}

func (admission *Admission) orient(ctx context.Context, session *Session, input Input, base Envelope) (output Envelope) {
	if len([]byte(input.Intent)) > admission.maxArgsBytes {
		return failure(base, "intent_too_large", "intent exceeded the size limit", map[string]any{"limit_bytes": admission.maxArgsBytes})
	}
	logContext := context.WithoutCancel(ctx)
	if admission.beginEpisodeAct(logContext, session, base, input) {
		defer func() {
			admission.completeEpisodeAct(logContext, session, base.ActID, output)
			session.rememberSuggestions(output)
		}()
	} else {
		defer func() { session.rememberSuggestions(output) }()
	}
	base.Status = StatusOK
	base.Result = map[string]any{"matched": false}
	base.Text = "ok orient: requested capability is unavailable in the reachable world; do not repeat this intent; finish without an executable act if no other ready call exists"
	base.Frontier = session.pendingSuggestions()
	if len(base.Frontier) == 0 && !session.hasExecutedAct() {
		base.Frontier = admission.activation.Compute(session.graph, input.Handle, input.Intent)
	}
	if len(base.Frontier) > 0 {
		base.Result = map[string]any{"matched": true}
		base.Text = fmt.Sprintf("ok orient: %d ready call(s)", len(base.Frontier))
	}
	return base
}

func (session *Session) nextExecutedIndex() int {
	session.mu.Lock()
	defer session.mu.Unlock()
	index := session.executed
	session.executed++
	return index
}

func (session *Session) hasExecutedAct() bool {
	session.mu.Lock()
	defer session.mu.Unlock()
	return session.executed > 0
}

func (admission *Admission) beginEpisodeAct(ctx context.Context, session *Session, base Envelope, input Input) bool {
	episodeID := session.EpisodeID()
	if admission.episodes == nil || episodeID == "" {
		return false
	}
	verbName, arguments := input.Verb, input.Args
	if input.Intent != "" {
		verbName, arguments = "orient", map[string]any{}
	}
	started := episode.StartedAct{
		EpisodeID: episodeID, ActID: base.ActID, WorldBuild: admission.worldBuild,
		Request: episode.ActRequest{
			Handle: input.Handle, Verb: verbName, Args: arguments, State: input.State, Intent: input.Intent,
		},
		SuggestionTaken: session.takeSuggestion(input),
	}
	if err := admission.episodes.BeginAct(ctx, started); err != nil {
		admission.disableEpisode(session, "begin act", err)
		return false
	}
	return true
}

func (admission *Admission) completeEpisodeAct(ctx context.Context, session *Session, actID string, result Envelope) {
	episodeID := session.EpisodeID()
	if episodeID == "" {
		return
	}
	if err := admission.episodes.CompleteAct(ctx, episode.CompletedAct{
		EpisodeID: episodeID, ActID: actID, Result: result,
	}); err != nil {
		admission.disableEpisode(session, "complete act", err)
	}
}

func (admission *Admission) disableEpisode(session *Session, operation string, err error) {
	session.mu.Lock()
	session.episodeID = ""
	session.mu.Unlock()
	admission.warn("episode logging disabled after %s: %v", operation, err)
}

func (admission *Admission) warn(format string, arguments ...any) {
	if admission.warning != nil {
		fmt.Fprintf(admission.warning, "phoenix warning: "+format+"\n", arguments...)
	}
}

func (session *Session) takeSuggestion(input Input) *episode.SuggestionLink {
	key, ok := inputSuggestionKey(input)
	if !ok {
		return nil
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	link, exists := session.suggested[key]
	if !exists {
		return nil
	}
	delete(session.suggested, key)
	return &link
}

func (session *Session) rememberSuggestions(envelope Envelope) {
	if session == nil {
		return
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	pending := append([]FrontierEntry(nil), envelope.Frontier...)
	if len(pending) == 0 && envelope.Refusal != nil && envelope.Refusal.Instead != nil {
		pending = []FrontierEntry{{
			Call: *envelope.Refusal.Instead, Why: envelope.Refusal.Why,
			Provenance: "authored", Score: 1,
		}}
	}
	session.pending = pending
	// Omitted-state restoration is scoped to the current pending frontier.
	// A call remembered from an earlier response must never regain a retired
	// precondition merely because its handle, verb, and arguments match.
	session.stateful = make(map[string]statefulSuggestion)
	for index, entry := range pending {
		key, err := suggestionKey(entry.Call)
		if err == nil {
			link := episode.SuggestionLink{ActID: envelope.ActID, Index: index}
			if index < len(envelope.Frontier) {
				session.suggested[key] = link
			}
			if entry.Call.State != nil {
				unstated := entry.Call
				unstated.State = nil
				if unstatedKey, keyErr := suggestionKey(unstated); keyErr == nil {
					value := statefulSuggestion{state: *entry.Call.State, link: link}
					if existing, exists := session.stateful[unstatedKey]; exists {
						value.ambiguous = existing.ambiguous || existing.state != value.state
					}
					session.stateful[unstatedKey] = value
				}
			}
		}
	}
}

func (session *Session) pendingSuggestions() []FrontierEntry {
	session.mu.Lock()
	defer session.mu.Unlock()
	return append([]FrontierEntry(nil), session.pending...)
}

func (session *Session) clearPending() {
	session.mu.Lock()
	defer session.mu.Unlock()
	session.pending = nil
	session.stateful = make(map[string]statefulSuggestion)
}

func (session *Session) completeSuggestedState(input Input) Input {
	call := frontier.Call{Handle: input.Handle, Verb: input.Verb, Args: input.Args}
	key, err := suggestionKey(call)
	if err != nil {
		return input
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	suggestion, exists := session.stateful[key]
	if !exists || suggestion.ambiguous {
		return input
	}
	delete(session.stateful, key)
	input.State = suggestion.state
	return input
}

func inputSuggestionKey(input Input) (string, bool) {
	call := frontier.Call{Handle: input.Handle, Verb: input.Verb, Args: input.Args}
	if input.State != "" {
		call.State = &input.State
	}
	key, err := suggestionKey(call)
	return key, err == nil
}

func suggestionKey(call frontier.Call) (string, error) {
	encoded, err := json.Marshal(call)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (admission *Admission) evaluateRefusal(base Envelope, session *Session, observation teach.Observation) *Envelope {
	if admission.suppressTeaching {
		return nil
	}
	refusal, matched, err := admission.teacher.Evaluate(session.graph, observation)
	if err != nil {
		failed := failure(base, "refusal_invalid", "authored refusal could not bind a reachable alternative", nil)
		return &failed
	}
	if !matched {
		return nil
	}
	shaped := refused(base, refusal)
	return &shaped
}

func failureFeatures(failure *verb.Failure) map[string]any {
	features := map[string]any{"code": failure.Code, "message": failure.Message}
	if failure.Details != nil {
		features["details"] = failure.Details
	}
	return features
}

func (admission *Admission) baseEnvelope(session *Session, input Input) Envelope {
	verb := input.Verb
	if input.Intent != "" {
		verb = "orient"
	}
	return Envelope{
		V: 1, WorldBuild: admission.worldBuild, SessionID: session.ID(), ActID: "a_" + rand.Text(),
		Handle: input.Handle, Verb: verb, Result: nil, Error: nil,
		Handles:  HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}},
		Frontier: []FrontierEntry{}, Refusal: nil,
	}
}

func absent(envelope Envelope) Envelope {
	envelope.Status = StatusAbsent
	envelope.Error = &Error{Code: "absent", Message: "act is not reachable"}
	envelope.Text = "absent: act is not reachable"
	return envelope
}

func failure(envelope Envelope, code, message string, details map[string]any) Envelope {
	envelope.Status = StatusFail
	envelope.Error = &Error{Code: code, Message: message, Details: details}
	envelope.Text = "failed " + envelope.Verb + ": " + message
	return envelope
}

func refused(envelope Envelope, refusal teach.Refusal) Envelope {
	envelope.Status = StatusRefused
	envelope.Refusal = &refusal
	envelope.Text = "refused " + refusal.What + ": " + refusal.Why
	if refusal.Instead != nil {
		envelope.Text += "\ninstead: " + refusal.Instead.Handle + "." + refusal.Instead.Verb
	}
	return envelope
}

func renderSuccess(handleType, name string, value any) string {
	text := "ok " + handleType + "." + name
	members, ok := value.(map[string]any)
	if !ok || len(members) == 0 {
		return text
	}
	if name == "build" && members["exit_code"] == float64(0) {
		text += ": build succeeded"
	}
	keys := make([]string, 0, len(members))
	for key := range members {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, 3)
	for _, key := range keys {
		if len(parts) == 3 {
			break
		}
		if rendered, useful := renderScalar(members[key]); useful {
			parts = append(parts, key+"="+rendered)
		}
	}
	if len(parts) > 0 {
		separator := ": "
		if strings.Contains(text, ":") {
			separator = ", "
		}
		text += separator + strings.Join(parts, ", ")
	}
	if len(text) > 240 {
		return text[:237] + "..."
	}
	return text
}

func renderScalar(value any) (string, bool) {
	switch typed := value.(type) {
	case string:
		return strings.ReplaceAll(typed, "\n", " "), true
	case bool:
		return fmt.Sprint(typed), true
	case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprint(typed), true
	case []any:
		return fmt.Sprintf("%d items", len(typed)), true
	default:
		return "", false
	}
}

type MCP struct {
	server    *mcp.Server
	admission *Admission
	sessions  sync.Map
}

func NewMCP(version string, admission *Admission) *MCP {
	if admission == nil {
		panic("surface admission is required")
	}
	adapter := &MCP{admission: admission}
	adapter.server = mcp.NewServer(
		&mcp.Implementation{Name: "phoenix", Title: "Phoenix", Version: version},
		&mcp.ServerOptions{Capabilities: &mcp.ServerCapabilities{Tools: &mcp.ToolCapabilities{}}},
	)
	adapter.server.AddTool(&mcp.Tool{
		Name: ToolName, Description: ToolDescription, InputSchema: actInputSchema,
	}, adapter.handle)
	return adapter
}

func (adapter *MCP) Server() *mcp.Server {
	return adapter.server
}

func (adapter *MCP) handle(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if request == nil || request.Params == nil || request.Session == nil {
		return toolError("invalid act arguments"), nil
	}
	input, err := decodeInput(request.Params.Arguments)
	if err != nil {
		return toolError("invalid act arguments"), nil
	}
	session, err := adapter.session(request.Session)
	if err != nil {
		return toolError("Phoenix session is unavailable"), nil
	}
	envelope := adapter.admission.Act(ctx, session, input)
	return &mcp.CallToolResult{
		Content:           []mcp.Content{&mcp.TextContent{Text: envelope.Text}},
		StructuredContent: envelope,
		IsError:           false,
	}, nil
}

func (adapter *MCP) session(serverSession *mcp.ServerSession) (*Session, error) {
	if existing, ok := adapter.sessions.Load(serverSession); ok {
		return existing.(*Session), nil
	}
	created, _, err := adapter.admission.StartSession()
	if err != nil {
		return nil, err
	}
	actual, loaded := adapter.sessions.LoadOrStore(serverSession, created)
	if !loaded {
		go func() {
			_ = serverSession.Wait()
			adapter.sessions.CompareAndDelete(serverSession, created)
		}()
	}
	return actual.(*Session), nil
}

func decodeInput(raw json.RawMessage) (Input, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var input Input
	if err := decoder.Decode(&input); err != nil {
		return Input{}, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return Input{}, fmt.Errorf("trailing JSON data")
	}
	execution := verbPattern.MatchString(input.Verb) && input.Args != nil && input.Intent == ""
	orientation := input.Verb == "" && input.Args == nil && strings.TrimSpace(input.Intent) != "" && len(input.Intent) <= 800 && input.State == ""
	if !handlePattern.MatchString(input.Handle) || (!execution && !orientation) {
		return Input{}, fmt.Errorf("invalid act identity")
	}
	if input.State != "" && !digestPattern.MatchString(input.State) {
		return Input{}, fmt.Errorf("invalid state digest")
	}
	return input, nil
}

func toolError(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: message}},
		IsError: true,
	}
}
