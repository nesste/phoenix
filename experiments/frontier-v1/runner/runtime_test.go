package main

import (
	"errors"
	"testing"
)

func TestInspectRuntimeOutputCapturesCompleteResultMetrics(t *testing.T) {
	output := []byte("{\"type\":\"system\",\"subtype\":\"init\",\"mcp_servers\":[{\"name\":\"phoenix\",\"status\":\"connected\"}]}\n" +
		"{\"type\":\"assistant\",\"message\":{}}\n" +
		"{\"type\":\"result\",\"subtype\":\"success\",\"result\":\"done\",\"num_turns\":3,\"duration_api_ms\":1250,\"total_cost_usd\":0.0125,\"usage\":{\"input_tokens\":10,\"cache_creation_input_tokens\":20,\"cache_read_input_tokens\":30,\"output_tokens\":40}}\n")
	observation, err := inspectRuntimeOutput(output)
	if err != nil {
		t.Fatal(err)
	}
	metrics := observation.metrics
	if !observation.modelOutputSeen || observation.finalMessage != "done" || metrics.Turns != 3 ||
		metrics.InputTokens != 10 || metrics.CacheCreationInputTokens != 20 ||
		metrics.CacheReadInputTokens != 30 || metrics.OutputTokens != 40 || metrics.TotalTokens != 100 ||
		metrics.CostUSD != 0.0125 || metrics.APITimeMS != 1250 {
		t.Fatalf("runtime observation = %#v", observation)
	}
}

func TestInspectRuntimeOutputAggregatesModelUsageFallback(t *testing.T) {
	output := []byte("{\"type\":\"assistant\"}\n" +
		"{\"type\":\"result\",\"result\":\"done\",\"modelUsage\":{\"model-a\":{\"inputTokens\":2,\"cacheCreationInputTokens\":3,\"cacheReadInputTokens\":5,\"outputTokens\":7,\"costUSD\":0.01},\"model-b\":{\"inputTokens\":11,\"outputTokens\":13,\"costUSD\":0.02}}}\n")
	observation, err := inspectRuntimeOutput(output)
	if err != nil {
		t.Fatal(err)
	}
	metrics := observation.metrics
	if metrics.InputTokens != 13 || metrics.CacheCreationInputTokens != 3 ||
		metrics.CacheReadInputTokens != 5 || metrics.OutputTokens != 20 ||
		metrics.TotalTokens != 41 || roundUSD(metrics.CostUSD) != 0.03 {
		t.Fatalf("model usage metrics = %#v", metrics)
	}
}

func TestRuntimeFailureClassificationRetriesOnlyEligiblePreTokenInfrastructure(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		message     string
		observation runtimeObservation
		code        string
		retry       bool
	}{
		{name: "provider", err: errors.New("exit"), message: "provider status 503", observation: runtimeObservation{resultError: true}, code: "provider_transient", retry: true},
		{name: "mcp", err: errors.New("exit"), message: "failed to connect to MCP server", code: "mcp_connection", retry: true},
		{name: "post token provider", err: errors.New("exit"), message: "provider status 503", observation: runtimeObservation{modelOutputSeen: true}, code: "runtime_error"},
		{name: "agent", err: errors.New("exit"), message: "agent failed", code: "runtime_error"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			failure := classifyRuntimeFailure(test.err, test.message, test.observation)
			if failure.Code != test.code || failure.RetryEligible != test.retry {
				t.Fatalf("failure = %#v", failure)
			}
		})
	}
}

func TestResultCapsAreTerminalITTFailures(t *testing.T) {
	for subtype, code := range map[string]string{
		"error_max_turns":      "turn_limit",
		"error_max_budget_usd": "cost_cap",
	} {
		failure := resultFailure(runtimeObservation{resultSubtype: subtype, modelOutputSeen: true})
		if failure == nil || failure.Code != code || !failure.CapHit || failure.RetryEligible {
			t.Fatalf("%s failure = %#v", subtype, failure)
		}
	}
}
