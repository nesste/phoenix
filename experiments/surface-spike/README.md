# Phoenix surface spike

This Phase 0 experiment compares a structured `act` tool with a restricted `eval` surface over the same two handles and six simulated verbs. It contains no learning, persistence, filesystem mutation, or production daemon code.

## Pinned environment

- Go 1.26.5 (`go 1.26.0` language version)
- MCP Go SDK v1.6.1
- Claude Code 2.1.229
- `claude-sonnet-5`, low effort
- stdio transport

## Run locally

From this directory:

```powershell
$env:GOTOOLCHAIN = 'go1.26.5'
go test ./...
go build -o ./bin/phoenix-spike.exe ./cmd/phoenix-spike
go run ./cmd/probe-stdio --exe ./bin/phoenix-spike.exe --surface act
go run ./cmd/probe-stdio --exe ./bin/phoenix-spike.exe --surface eval
go run ./cmd/measure-schema
go run ./cmd/validate-spec --repo-root ../..
```

The checked-in MCP configs use a repository-relative executable path. Build the executable and invoke Claude Code from the repository root.

```powershell
claude -p "Call phoenix act on h_repo_demo with verb status and empty args. Then call the first frontier suggestion. Report both result texts exactly." `
  --model claude-sonnet-5 `
  --effort low `
  --max-turns 2 `
  --max-budget-usd 0.20 `
  --no-session-persistence `
  --output-format json `
  --strict-mcp-config `
  --mcp-config ./experiments/surface-spike/mcp-act.json `
  --tools "" `
  --allowedTools mcp__phoenix__act `
  --system-prompt "Use only the configured Phoenix tool. Follow the user request exactly."
```

Run the corresponding `eval` probe with `mcp-eval.json`, `mcp__phoenix__eval`, and `repo.status()`.

## Acceptance

The local MCP, schema, and O(1) checks pass. On 2026-08-18, the pinned runtime completed 20 fresh one-call sessions per surface with zero malformed calls. Each observed rate is 0%, with a two-sided Wilson 95% interval of `[0, 0.161125]`. A `CGO_ENABLED=0` rerun on Ubuntu 24.04.4 under WSL2 passed the Go suite and stdio round trips. Across 100 Linux `act` calls, daemon-only latency had median 636.5 microseconds and p95 837 microseconds.

The evidence is complete but not independently accepted. Do not mark Task 0.2 accepted or open Phase 0 until an independent reviewer verifies `results.json`, including the final session IDs, Wilson calculation, Linux sample, and gate arithmetic.
