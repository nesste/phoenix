package spike

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestExecuteReturnsBoundFrontier(t *testing.T) {
	result := Execute(ActInput{Handle: RepoHandle, Verb: "status"})
	if result.Status != "ok" {
		t.Fatalf("status = %q, want ok", result.Status)
	}
	if len(result.Frontier) != 1 {
		t.Fatalf("frontier length = %d, want 1", len(result.Frontier))
	}
	call := result.Frontier[0].Call
	if call.Handle != TestsHandle || call.Verb != "run" || call.Args == nil {
		t.Fatalf("frontier call = %#v, want bound tests.run call", call)
	}
}

func TestExecuteDoesNotRevealUnknownHandles(t *testing.T) {
	result := Execute(ActInput{Handle: "h_guessed", Verb: "status"})
	if result.Status != "absent" {
		t.Fatalf("status = %q, want absent", result.Status)
	}
	if result.Refusal != nil || len(result.Frontier) != 0 {
		t.Fatalf("absent result leaked guidance: %#v", result)
	}
}

func TestSurfaceSchemaGrowth(t *testing.T) {
	act10 := schemaBytes(t, "act", 10)
	act1000 := schemaBytes(t, "act", 1000)
	if act10 != act1000 {
		t.Fatalf("act schema bytes changed with world size: 10=%d 1000=%d", act10, act1000)
	}

	flat10 := schemaBytes(t, "flat", 10)
	flat100 := schemaBytes(t, "flat", 100)
	flat1000 := schemaBytes(t, "flat", 1000)
	if !(flat10 < flat100 && flat100 < flat1000) {
		t.Fatalf("flat schemas did not grow: 10=%d 100=%d 1000=%d", flat10, flat100, flat1000)
	}
}

func TestFrontierWhyLinesFitBudget(t *testing.T) {
	for _, call := range []ActInput{
		{Handle: RepoHandle, Verb: "status"},
		{Handle: RepoHandle, Verb: "build"},
		{Handle: TestsHandle, Verb: "run"},
		{Handle: TestsHandle, Verb: "list"},
	} {
		for _, entry := range Execute(call).Frontier {
			if len(entry.Why) > 80 {
				t.Errorf("why line has %d bytes: %q", len(entry.Why), entry.Why)
			}
		}
	}
}

func schemaBytes(t *testing.T, surface string, count int) int {
	t.Helper()
	ctx := context.Background()
	server, err := NewServer(surface, count)
	if err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "schema-probe", Version: "0.1.0"}, nil)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()
	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(tools.Tools)
	if err != nil {
		t.Fatal(err)
	}
	return len(encoded)
}
