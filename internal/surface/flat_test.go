package surface

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nesste/phoenix/internal/world"
)

func TestFlatMCPExposesCompleteSchemasAndUsesSharedAdmission(t *testing.T) {
	definition := flatTestDefinition()
	admission := testAdmission(t, 1024)
	adapter, err := NewFlatMCP("test", admission, definition)
	if err != nil {
		t.Fatal(err)
	}
	clientSession, cleanup := connectServer(t, adapter.Server())
	defer cleanup()

	listed, err := clientSession.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed.Tools) != 1 || listed.Tools[0].Name != "repo_inspect" {
		t.Fatalf("flat tools = %#v, want repo_inspect", listed.Tools)
	}
	assertFlatToolSchema(t, listed.Tools[0])
	assertFlatToolResult(t, clientSession)
}

func assertFlatToolSchema(t *testing.T, tool *mcp.Tool) {
	t.Helper()
	encodedSchema, err := json.Marshal(tool.InputSchema)
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(encodedSchema, &schema); err != nil {
		t.Fatal(err)
	}
	properties, _ := schema["properties"].(map[string]any)
	if schema["additionalProperties"] != false || properties["detail"] == nil {
		t.Fatalf("input schema = %s, want complete inspect schema", encodedSchema)
	}
}

func assertFlatToolResult(t *testing.T, clientSession *mcp.ClientSession) {
	t.Helper()
	result, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "repo_inspect", Arguments: json.RawMessage(`{"detail":"brief"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError || result.StructuredContent == nil {
		t.Fatalf("flat result = %#v, want typed success", result)
	}
	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["summary"] != "repository brief" || payload["world_build"] != nil || payload["frontier"] != nil {
		t.Fatalf("flat payload leaked Phoenix envelope or lost result: %#v", payload)
	}
}

func TestFlatMCPPlainFailureDoesNotExposeTeachingOrFrontier(t *testing.T) {
	definition := flatTestDefinition()
	admission := testAdmissionWithPolicy(t, true, true)
	adapter, err := NewFlatMCP("test", admission, definition)
	if err != nil {
		t.Fatal(err)
	}
	clientSession, cleanup := connectServer(t, adapter.Server())
	defer cleanup()

	result, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "repo_inspect", Arguments: json.RawMessage(`{"detail":"blocked"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(result.StructuredContent)
	var payload map[string]any
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["status"] != "fail" || payload["frontier"] != nil || payload["refusal"] != nil {
		t.Fatalf("flat failure = %#v, want plain typed failure", payload)
	}
}

func flatTestDefinition() *world.Definition {
	return &world.Definition{
		V: 1, ID: "flat_test",
		Roots: []world.Root{{
			Name: "repo", Type: "repo", Label: "current repository",
			Resource: world.Resource{Kind: "logical", Value: "repo"},
		}},
		HandleTypes: map[string]world.HandleType{"repo": {Verbs: map[string]world.Verb{
			"inspect": {ArgsSchema: inspectArgs, ResultSchema: inspectResult},
		}}},
	}
}

func connectServer(t *testing.T, server *mcp.Server) (*mcp.ClientSession, func()) {
	t.Helper()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "flat-surface-test", Version: "test"}, nil)
	clientSession, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		serverSession.Close()
		t.Fatal(err)
	}
	return clientSession, func() {
		clientSession.Close()
		serverSession.Close()
	}
}
