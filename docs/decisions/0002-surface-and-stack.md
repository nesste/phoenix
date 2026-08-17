# 0002: Phase 0 surface and stack

- **Status:** Provisional; Claude model round trip and Linux rerun pending
- **Date:** 2026-08-17
- **Experiment:** `experiments/surface-spike`

## Pins

- Go language version: 1.26
- Go toolchain: 1.26.5
- MCP Go SDK: `github.com/modelcontextprotocol/go-sdk` v1.6.1
- MCP protocol negotiated by that stable SDK: up to 2025-11-25
- Transport: stdio
- Claude Code: 2.1.229
- Model: `claude-sonnet-5`
- Claude effort: low
- Session persistence: disabled
- Maximum turns per probe: 1 for the baseline, 2 for tool calls
- Maximum spend per probe: USD 0.20

The official MCP Go SDK lists v1.6.1 as its current stable release. Support for the 2026-07-28 MCP protocol is available only in the v1.7 prerelease line as of this decision, so Phoenix does not pin that prerelease. The SDK and protocol compatibility table is maintained in the [official MCP Go SDK repository](https://github.com/modelcontextprotocol/go-sdk), and the current Go release history is maintained by the [Go project](https://go.dev/doc/devel/release).

Claude Code accepts stdio MCP configuration through `--mcp-config`, supports a strict configuration mode, and can pin `claude-sonnet-5` by its canonical model ID. The invocation follows the [Claude Code CLI](https://code.claude.com/docs/en/cli-usage) and [MCP configuration](https://code.claude.com/docs/en/mcp) references.

## Candidates

The spike implements two one-tool surfaces over the same six simulated verbs and two handles:

- `act(handle, verb, args)` uses separate structured fields;
- `eval(expression)` accepts one restricted object-graph expression such as `repo.status()`.

Both return the same compact text and structured frontier. The `eval` spike is deliberately restricted and does not execute arbitrary code.

## Evidence so far

Both candidates compile and complete an MCP initialize/list/call round trip through an actual stdio subprocess. The returned text and structured content both contain a legible, fully bound next call.

The serialized MCP tool list is flat across synthetic worlds:

| Surface | 10 verbs | 100 verbs | 1,000 verbs |
| --- | ---: | ---: | ---: |
| `act` | 1,621 bytes | 1,621 bytes | 1,621 bytes |
| `eval` | 1,490 bytes | 1,490 bytes | 1,490 bytes |
| flat tools | 3,861 bytes | 38,601 bytes | 386,001 bytes |

This proves the structural O(1) property for the one-tool candidates. It does not substitute for the runtime token measurement required by the Phase 0 gate.

The first pinned Claude Code probe reached the service but returned HTTP 429 because the account session limit had been reached. It consumed no model tokens and made no tool call. Token cost, model-visible rendering, malformed-call rate, and take-up of the returned frontier therefore remain unmeasured.

## Provisional decision

Keep Go, the stable official MCP Go SDK, and stdio. Keep structured `act` as the provisional surface: its schema costs 131 more serialized bytes than `eval`, but it preserves typed validation and avoids parsing an expression language. Do not freeze this choice until the same pinned Claude Code run:

1. completes an `act` call and follows its returned frontier;
2. completes the equivalent `eval` call;
3. records actual standing input-token deltas against the empty-MCP baseline;
4. records malformed-call outcomes over the pre-committed repetitions;
5. reruns the stdio probe on Linux.

## Known failure modes

- A daemon crash closes the stdio connection. The Phase 0 spike has no persistence or reconciliation, so the client receives no trustworthy result and must reconnect.
- The spike does not cap result size. Production admission and result rendering must reject or truncate oversized data before persistence.
- Claude Code starts one stdio subprocess per configured session. Phase 1 sessions remain isolated; shared mutable worlds are deferred.
- An MCP client may cache the tool list. Measurements must use fresh sessions and record whether cache creation or cache reads contributed to reported input tokens.
- The current Windows result establishes portability of the spike only. Linux remains the target acceptance platform.

## Gate effect

Task 0.2 is implemented but not accepted. Phase 0 remains closed until the pending model and Linux measurements are recorded in `experiments/surface-spike/results.json` and this decision changes from provisional to accepted or rejected.
