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
	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
)

const (
	ToolName        = "act"
	ToolDescription = "Invoke one verb on a live Phoenix handle."
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
	"required":["handle","verb","args"],
	"properties":{
		"handle":{"type":"string","pattern":"^h_[A-Za-z0-9_-]{16,}$"},
		"verb":{"type":"string","pattern":"^[a-z][a-z0-9_]{0,63}$"},
		"args":{"type":"object"},
		"state":{"type":"string","pattern":"^sha256:[0-9a-f]{64}$"}
	}
}`)

type Input struct {
	Handle string         `json:"handle"`
	Verb   string         `json:"verb"`
	Args   map[string]any `json:"args"`
	State  string         `json:"state,omitempty"`
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

type Call struct {
	Handle string         `json:"handle"`
	Verb   string         `json:"verb"`
	Args   map[string]any `json:"args"`
	State  *string        `json:"state,omitempty"`
}

type FrontierEntry struct {
	Call       Call    `json:"call"`
	Why        string  `json:"why"`
	Provenance string  `json:"provenance"`
	Score      float64 `json:"score"`
}

type Refusal struct {
	What          string `json:"what"`
	Why           string `json:"why"`
	Instead       *Call  `json:"instead"`
	NoAlternative string `json:"no_alternative,omitempty"`
}

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
	WorldBuild   string
	Graph        *world.Graph
	Executor     *verb.Executor
	MaxArgsBytes int
}

type Admission struct {
	worldBuild   string
	graph        *world.Graph
	executor     *verb.Executor
	maxArgsBytes int
}

type Session struct {
	graph *world.Session
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
	if config.MaxArgsBytes <= 0 {
		config.MaxArgsBytes = defaultMaxArgs
	}
	return &Admission{
		worldBuild: config.WorldBuild, graph: config.Graph,
		executor: config.Executor, maxArgsBytes: config.MaxArgsBytes,
	}, nil
}

func (admission *Admission) StartSession() (*Session, []world.RootHandle, error) {
	graphSession, roots, err := admission.graph.StartSession()
	if err != nil {
		return nil, nil, err
	}
	return &Session{graph: graphSession}, roots, nil
}

func (session *Session) ID() string {
	if session == nil || session.graph == nil {
		return ""
	}
	return session.graph.ID()
}

func (admission *Admission) Act(ctx context.Context, session *Session, input Input) Envelope {
	base := admission.baseEnvelope(session, input)
	if session == nil || session.graph == nil {
		return failure(base, "invalid_session", "session is unavailable", nil)
	}
	encodedArgs, err := json.Marshal(input.Args)
	if err != nil || input.Args == nil {
		return failure(base, "invalid_arguments", "act arguments must be a JSON object", nil)
	}
	if len(encodedArgs) > admission.maxArgsBytes {
		return failure(base, "arguments_too_large", "act arguments exceeded the size limit", map[string]any{
			"limit_bytes": admission.maxArgsBytes,
		})
	}

	access := session.graph.Prepare(ctx, input.Handle, input.Verb, input.State)
	switch access.Status {
	case world.AccessAbsent:
		return absent(base)
	case world.AccessStale, world.AccessFailed:
		return failure(base, access.Problem.Code, access.Problem.Message, nil)
	case world.AccessReady:
		result := admission.executor.Execute(ctx, verb.Request{
			SessionID: session.ID(), HandleType: access.Target.Type, Verb: input.Verb,
			Resource: access.Target.Resource, State: access.State, Args: input.Args,
			Reachable: session.graph,
		})
		if result.Status == verb.StatusFail {
			return failure(base, result.Error.Code, result.Error.Message, result.Error.Details)
		}
		base.Status = StatusOK
		base.Result = result.Value
		base.Text = renderSuccess(access.Target.Type, input.Verb, result.Value)
		return base
	default:
		return failure(base, "execution_failed", "verb execution failed", nil)
	}
}

func (admission *Admission) baseEnvelope(session *Session, input Input) Envelope {
	return Envelope{
		V: 1, WorldBuild: admission.worldBuild, SessionID: session.ID(), ActID: "a_" + rand.Text(),
		Handle: input.Handle, Verb: input.Verb, Result: nil, Error: nil,
		Handles:  HandleDelta{Grant: []HandleGrant{}, Revoke: []string{}},
		Frontier: []FrontierEntry{}, Refusal: nil,
	}
}

func absent(envelope Envelope) Envelope {
	envelope.Status = StatusAbsent
	envelope.Error = &Error{Code: "absent", Message: "handle is not reachable"}
	envelope.Text = "absent: handle is not reachable"
	return envelope
}

func failure(envelope Envelope, code, message string, details map[string]any) Envelope {
	envelope.Status = StatusFail
	envelope.Error = &Error{Code: code, Message: message, Details: details}
	envelope.Text = "failed " + envelope.Verb + ": " + message
	return envelope
}

func renderSuccess(handleType, name string, value any) string {
	text := "ok " + handleType + "." + name
	members, ok := value.(map[string]any)
	if !ok || len(members) == 0 {
		return text
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
		text += ": " + strings.Join(parts, ", ")
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
		IsError:           envelope.Status != StatusOK,
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
	if !handlePattern.MatchString(input.Handle) || !verbPattern.MatchString(input.Verb) || input.Args == nil {
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
