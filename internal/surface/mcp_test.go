package surface

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

const testWorldBuild = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestActReturnsCompactTextAndStructuredEnvelope(t *testing.T) {
	admission := testAdmission(t, 1024)
	session, roots, err := admission.StartSession()
	if err != nil {
		t.Fatal(err)
	}

	envelope := admission.Act(context.Background(), session, Input{
		Handle: roots[0].Ref,
		Verb:   "inspect",
		Args:   map[string]any{"detail": "brief"},
	})
	if envelope.Status != StatusOK {
		t.Fatalf("status = %q, want ok: %#v", envelope.Status, envelope.Error)
	}
	if envelope.SessionID != session.ID() || envelope.WorldBuild != testWorldBuild {
		t.Fatalf("envelope identity = %#v, want session and world build", envelope)
	}
	if envelope.Text == "" || strings.Contains(envelope.Text, "{") {
		t.Fatalf("text = %q, want compact rendering", envelope.Text)
	}
	if len(envelope.Frontier) != 1 || envelope.Handles.Grant == nil || envelope.Handles.Revoke == nil {
		t.Fatalf("frontier/delta not encoded explicitly: %#v", envelope)
	}
	if got := envelope.Frontier[0]; got.Call.Handle != roots[0].Ref || got.Call.Verb != "inspect" || got.Call.Args["detail"] != "repository brief" {
		t.Fatalf("frontier call = %#v, want reachable call with bound result", got)
	}
	validateEnvelope(t, envelope)
}

func TestUnknownHandleIsGenericAbsent(t *testing.T) {
	admission := testAdmission(t, 1024)
	session, _, err := admission.StartSession()
	if err != nil {
		t.Fatal(err)
	}

	envelope := admission.Act(context.Background(), session, Input{
		Handle: "h_0123456789abcdef",
		Verb:   "inspect",
		Args:   map[string]any{},
	})
	if envelope.Status != StatusAbsent || envelope.Error == nil {
		t.Fatalf("result = %#v, want absent error", envelope)
	}
	if envelope.Error.Code != "absent" || envelope.Error.Message != "handle is not reachable" {
		t.Fatalf("error = %#v, want generic absence", envelope.Error)
	}
	if envelope.Result != nil || len(envelope.Frontier) != 0 || envelope.Refusal != nil {
		t.Fatalf("absent result leaked topology: %#v", envelope)
	}
	validateEnvelope(t, envelope)
}

func TestAdmissionRejectsOversizedArguments(t *testing.T) {
	admission := testAdmission(t, 32)
	session, roots, err := admission.StartSession()
	if err != nil {
		t.Fatal(err)
	}

	envelope := admission.Act(context.Background(), session, Input{
		Handle: roots[0].Ref,
		Verb:   "inspect",
		Args:   map[string]any{"detail": strings.Repeat("x", 64)},
	})
	if envelope.Status != StatusFail || envelope.Error == nil || envelope.Error.Code != "arguments_too_large" {
		t.Fatalf("result = %#v, want arguments_too_large failure", envelope)
	}
	validateEnvelope(t, envelope)
}

func TestMCPRejectsMalformedCallAtAdmission(t *testing.T) {
	malformed := []json.RawMessage{
		json.RawMessage(`{"handle":7,"verb":"inspect","args":{}}`),
		json.RawMessage(`{"handle":"h_0123456789abcdef","verb":"inspect"}`),
		json.RawMessage(`{"handle":"h_0123456789abcdef","verb":"inspect","args":{},"extra":true}`),
	}
	for index, arguments := range malformed {
		adapter := NewMCP("test", testAdmission(t, 1024))
		clientSession, _, cleanup := connect(t, adapter)
		result, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
			Name: ToolName, Arguments: arguments,
		})
		if err != nil {
			cleanup()
			t.Fatalf("case %d call malformed act: %v", index, err)
		}
		if !result.IsError || len(result.Content) != 1 || result.StructuredContent != nil {
			cleanup()
			t.Fatalf("case %d malformed result = %#v, want unstructured tool error", index, result)
		}
		text, ok := result.Content[0].(*mcp.TextContent)
		if !ok || text.Text != "invalid act arguments" {
			cleanup()
			t.Fatalf("case %d error content = %#v, want generic admission error", index, result.Content)
		}
		cleanup()
	}
}

func TestDecodeInputRejectsTrailingJSON(t *testing.T) {
	_, err := decodeInput(json.RawMessage(`{"handle":"h_0123456789abcdef","verb":"inspect","args":{}} true`))
	if err == nil {
		t.Fatal("decodeInput accepted trailing JSON")
	}
}

func TestConcurrentSessionsRemainIsolated(t *testing.T) {
	admission := testAdmission(t, 1024)
	first, firstRoots, err := admission.StartSession()
	if err != nil {
		t.Fatal(err)
	}
	second, secondRoots, err := admission.StartSession()
	if err != nil {
		t.Fatal(err)
	}

	type invocation struct {
		session *Session
		handle  string
		want    Status
	}
	calls := []invocation{
		{first, firstRoots[0].Ref, StatusOK},
		{second, secondRoots[0].Ref, StatusOK},
		{first, secondRoots[0].Ref, StatusAbsent},
		{second, firstRoots[0].Ref, StatusAbsent},
	}
	var wait sync.WaitGroup
	for index, call := range calls {
		wait.Add(1)
		go func(index int, call invocation) {
			defer wait.Done()
			got := admission.Act(context.Background(), call.session, Input{
				Handle: call.handle, Verb: "inspect", Args: map[string]any{},
			})
			if got.Status != call.want {
				t.Errorf("call %d status = %q, want %q", index, got.Status, call.want)
			}
		}(index, call)
	}
	wait.Wait()
}

func TestMCPServesOnlyActWithStructuredOutput(t *testing.T) {
	adapter := NewMCP("test", testAdmission(t, 1024))
	clientSession, serverSession, cleanup := connect(t, adapter)
	defer cleanup()
	session, roots, err := adapter.admission.StartSession()
	if err != nil {
		t.Fatal(err)
	}
	adapter.sessions.Store(serverSession, session)

	tools, err := clientSession.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Name != ToolName {
		t.Fatalf("tools = %#v, want only act", tools.Tools)
	}
	result, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
		Name: ToolName,
		Arguments: json.RawMessage(fmt.Sprintf(
			`{"handle":%q,"verb":"inspect","args":{"detail":"brief"}}`, roots[0].Ref,
		)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError || result.StructuredContent == nil || len(result.Content) != 1 {
		t.Fatalf("tool result = %#v, want text plus structured envelope", result)
	}
	var envelope Envelope
	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(encoded, &envelope); err != nil {
		t.Fatal(err)
	}
	validateEnvelope(t, envelope)
}

func testAdmission(t *testing.T, maxArgs int) *Admission {
	t.Helper()
	definition := &world.Definition{
		V:  1,
		ID: "test",
		Roots: []world.Root{{
			Name: "repo", Type: "repo", Label: "repository",
			Resource: world.Resource{Kind: "logical", Value: "repo"},
		}},
		HandleTypes: map[string]world.HandleType{
			"repo": {Verbs: map[string]world.Verb{
				"inspect": {ArgsSchema: inspectArgs, ResultSchema: inspectResult},
			}},
		},
		Transitions: []world.Transition{{
			ID: "inspect_again",
			Match: world.Match{
				HandleType: "repo", Verb: "inspect", Status: "ok",
				ResultWhen: json.RawMessage(`{"required":["summary"]}`),
			},
			Suggestions: []world.Suggestion{{
				Call: world.CallTemplate{
					Handle: world.HandleSelector{Source: "self"}, Verb: "inspect",
					Args: map[string]world.Binding{"detail": resultBinding("/summary")},
				},
				Why: "inspect the current summary", Score: 1,
			}},
		}},
	}
	registry := verb.NewRegistry()
	err := registry.Register(verb.Definition{
		HandleType: "repo", Name: "inspect", ArgsSchema: inspectArgs, ResultSchema: inspectResult,
		Handler: verb.HandlerFunc(func(_ context.Context, request verb.Request) (any, error) {
			detail, _ := request.Args["detail"].(string)
			return map[string]any{"summary": "repository " + detail}, nil
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	frontierEngine, err := frontier.New(definition)
	if err != nil {
		t.Fatal(err)
	}
	admission, err := New(Config{
		WorldBuild:   testWorldBuild,
		Graph:        world.NewGraph(definition, staticResolver{}),
		Executor:     verb.NewExecutor(registry, verb.Options{}),
		Frontier:     frontierEngine,
		MaxArgsBytes: maxArgs,
	})
	if err != nil {
		t.Fatal(err)
	}
	return admission
}

func resultBinding(pointer string) world.Binding {
	return world.Binding{ResultPointer: &pointer}
}

var inspectArgs = json.RawMessage(`{
	"type":"object",
	"additionalProperties":false,
	"properties":{"detail":{"type":"string"}}
}`)

var inspectResult = json.RawMessage(`{
	"type":"object",
	"additionalProperties":false,
	"required":["summary"],
	"properties":{"summary":{"type":"string"}}
}`)

type staticResolver struct{}

func (staticResolver) Resolve(context.Context, world.Resource) (any, error) {
	return map[string]any{"revision": 1}, nil
}

func validateEnvelope(t *testing.T, envelope Envelope) {
	t.Helper()
	encoded, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	instance, err := jsonschema.UnmarshalJSON(strings.NewReader(string(encoded)))
	if err != nil {
		t.Fatal(err)
	}
	schema, err := jsonschema.NewCompiler().Compile("../../spec/result.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := schema.Validate(instance); err != nil {
		t.Fatalf("envelope does not match frozen schema: %v\n%s", err, encoded)
	}
}

func connect(t *testing.T, adapter *MCP) (*mcp.ClientSession, *mcp.ServerSession, func()) {
	t.Helper()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := adapter.Server().Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "surface-test", Version: "test"}, nil)
	clientSession, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		serverSession.Close()
		t.Fatal(err)
	}
	return clientSession, serverSession, func() {
		clientSession.Close()
		serverSession.Close()
	}
}
