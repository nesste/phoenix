package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nesste/phoenix/internal/surface"
)

func TestStateEventsApplyAtConfiguredExecutableAct(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	target := filepath.Join(root, "state.txt")
	if err := os.WriteFile(target, []byte("before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	planPath := filepath.Join(t.TempDir(), "events.json")
	if err := os.WriteFile(planPath, []byte(`{"v":1,"events":[{"after_act":0,"path":"state.txt","content":"after\n"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	hook, err := loadStateEvents(planPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := hook(context.Background(), 0, surface.Input{}, surface.Envelope{}); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "after\n" {
		t.Fatalf("state contents = %q, want applied event", contents)
	}
}

func TestStateEventsRejectPathsOutsideSandbox(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	planPath := filepath.Join(t.TempDir(), "events.json")
	if err := os.WriteFile(planPath, []byte(`{"v":1,"events":[{"after_act":0,"path":"../outside.txt","content":"x"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadStateEvents(planPath); err == nil {
		t.Fatal("loadStateEvents accepted an escaping path")
	}
}

func TestFlatArmAppliesStateEventAfterComputingResult(t *testing.T) {
	repositoryRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	sandbox := t.TempDir()
	t.Chdir(sandbox)
	if err := os.WriteFile("state.txt", []byte("before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	planPath := filepath.Join(t.TempDir(), "events.json")
	if err := os.WriteFile(planPath, []byte(`{"v":1,"events":[{"after_act":0,"path":"state.txt","content":"after\n"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	assembled, err := assembleSurface(serveOptions{
		arm: "A", serverVersion: "test",
		worldPath:       filepath.Join(repositoryRoot, "worlds", "dev-repo", "world.json"),
		schemaPath:      filepath.Join(repositoryRoot, "spec", "world.schema.json"),
		stateEventsPath: planPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer assembled.Close()

	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := assembled.server.Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "flat-event-test", Version: "test"}, nil)
	clientSession, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()
	result, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "repo_read", Arguments: json.RawMessage(`{"path":"state.txt"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(encoded, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["content"] != "before\n" {
		t.Fatalf("flat result content = %q, want pre-event bytes", payload["content"])
	}
	contents, err := os.ReadFile("state.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "after\n" {
		t.Fatalf("state event contents = %q, want post-result replacement", contents)
	}
}
