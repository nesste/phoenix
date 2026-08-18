package surface

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nesste/phoenix/internal/activate"
	"github.com/nesste/phoenix/internal/frontier"
	"github.com/nesste/phoenix/internal/teach"
	"github.com/nesste/phoenix/internal/verb"
	"github.com/nesste/phoenix/internal/world"
)

const standingTokenBudget = 600

func TestStandingSurfaceIsFlatAndUnderBudget(t *testing.T) {
	sizes := []int{10, 100, 1000}
	var baseline string
	for _, size := range sizes {
		serialized := serializedSurface(t, size)
		if baseline == "" {
			baseline = serialized
		} else if serialized != baseline {
			t.Fatalf("standing surface grew with world: 10 verbs=%d bytes, %d verbs=%d bytes", len(baseline), size, len(serialized))
		}
		if tokens := asciiByteTokenCeiling(serialized); tokens > standingTokenBudget {
			t.Fatalf("%d-verb standing surface = %d tokens, budget = %d", size, tokens, standingTokenBudget)
		}
	}
}

func TestStandingSurfaceContainsNoInstructionsOrRuleText(t *testing.T) {
	adapter := syntheticMCP(t, 1000)
	clientSession, _, cleanup := connect(t, adapter)
	defer cleanup()
	if instructions := clientSession.InitializeResult().Instructions; instructions != "" {
		t.Fatalf("standing instructions = %q, want empty", instructions)
	}
	tools, err := clientSession.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Description != ToolDescription {
		t.Fatalf("served description changed: %#v", tools.Tools)
	}
}

func serializedSurface(t *testing.T, verbCount int) string {
	t.Helper()
	adapter := syntheticMCP(t, verbCount)
	clientSession, _, cleanup := connect(t, adapter)
	defer cleanup()
	tools, err := clientSession.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	payload := struct {
		Instructions string      `json:"instructions"`
		Tools        []*mcp.Tool `json:"tools"`
	}{Instructions: clientSession.InitializeResult().Instructions, Tools: tools.Tools}
	encoded, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

func syntheticMCP(t *testing.T, verbCount int) *MCP {
	t.Helper()
	verbs := make(map[string]world.Verb, verbCount)
	for index := 0; index < verbCount; index++ {
		verbs[fmt.Sprintf("verb_%04d", index)] = world.Verb{ArgsSchema: inspectArgs, ResultSchema: inspectResult}
	}
	definition := &world.Definition{
		V: 1, ID: "synthetic",
		Roots:       []world.Root{{Name: "root", Type: "synthetic", Label: "synthetic", Resource: world.Resource{Kind: "logical", Value: "root"}}},
		HandleTypes: map[string]world.HandleType{"synthetic": {Verbs: verbs}},
		Transitions: []world.Transition{},
	}
	frontierEngine, err := frontier.New(definition)
	if err != nil {
		t.Fatal(err)
	}
	activationEngine, err := activate.New(definition)
	if err != nil {
		t.Fatal(err)
	}
	teachingEngine, err := teach.New(definition)
	if err != nil {
		t.Fatal(err)
	}
	admission, err := New(Config{
		WorldBuild:   testWorldBuild,
		Graph:        world.NewGraph(definition, staticResolver{}),
		Executor:     verb.NewExecutor(verb.NewRegistry(), verb.Options{}),
		Frontier:     frontierEngine,
		Activation:   activationEngine,
		Teacher:      teachingEngine,
		MaxArgsBytes: 1024,
	})
	if err != nil {
		t.Fatal(err)
	}
	return NewMCP("test", admission)
}

// asciiByteTokenCeiling treats every serialized byte as a token. The served
// surface is ASCII, so this is stricter than the accepted runtime measurement.
func asciiByteTokenCeiling(value string) int {
	return len(value)
}
