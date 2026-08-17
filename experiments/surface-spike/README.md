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

The checked-in MCP configs resolve the executable from `${CLAUDE_PROJECT_DIR}`. Build the executable before invoking Claude Code from the repository root.

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

The local MCP, schema, and O(1) checks pass. The model-facing run is pending because the first attempt returned an account session-limit HTTP 429 before inference. Do not mark Task 0.2 accepted until `results.json` contains successful pinned `act` and `eval` runs, real token deltas, malformed-call measurements, and a Linux rerun.
