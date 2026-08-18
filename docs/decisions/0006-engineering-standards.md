# 0006: Phase 1 engineering standards

- **Status:** Accepted for Phase 1 implementation
- **Date:** 2026-08-18
- **Scope:** Root Phoenix production module

## Decision

Phoenix Phase 1 uses a Linux-amd64, CGO-disabled Go build and one reproducible `make quality` entry point. A change is not complete when only its focused test passes: the root test suite, formatting check, static analysis, dependency checks, schema validation, complexity and duplication limits, and production build must all pass.

The Phase 0 experiment modules remain independently reproducible with their accepted commands. Root quality invokes the accepted schema validator but does not rewrite or open either sealed evaluation tranche.

## Pinned runtime and dependencies

| Component | Pin | Rule |
| --- | --- | --- |
| Go language | `1.26.0` | Declared in root `go.mod`. |
| Go toolchain | `go1.26.6` | Declared in root `go.mod`; CI must not float to a newer patch. Phase 0 evidence remains pinned to `go1.26.5`; production advanced to the security-fixed patch after the Task 1.1 vulnerability gate found reachable standard-library advisories in `1.26.5`. |
| MCP Go SDK | `github.com/modelcontextprotocol/go-sdk v1.6.1` | Preserves the independently accepted Phase 0 surface. Upgrade only through a reviewed decision and repeated surface checks. |
| JSON canonicalization | `github.com/gowebpki/jcs v1.0.1` | RFC 8785 canonicalization for world and live-state digests. |
| JSON Schema | `github.com/santhosh-tekuri/jsonschema/v6 v6.0.3` | Draft 2020-12 validation. |
| SQLite | `modernc.org/sqlite v1.56.0` | Pure-Go driver reserved for the Phase 1 episode store; no `mattn/go-sqlite3` or other CGO driver. |
| Staticcheck | `honnef.co/go/tools/cmd/staticcheck v0.7.0` | Invoked by exact version from `make lint`. |
| Govulncheck | `golang.org/x/vuln/cmd/govulncheck v1.7.0` | Invoked by exact version from `make dependency-check`. |
| Gocyclo | `github.com/fzipp/gocyclo/cmd/gocyclo v0.6.0` | Fails above the complexity limit. |
| Dupl | `github.com/mibk/dupl v1.1.0` | Fails at duplicate blocks of 100 tokens or more. |

All production builds set `CGO_ENABLED=0`, `GOOS=linux`, and `GOARCH=amd64` unless a later decision adds a target. Build flags remove host paths and VCS metadata. The generated manifest records the target and the fact that CGO is disabled.

Direct runtime dependencies must use exact module versions in `go.mod`; transitive versions are fixed by `go.sum`. Dependency updates are isolated changes with release-note review, `go mod verify`, `govulncheck`, the full quality gate, and any protocol-compatibility reruns affected by the update. A vulnerability finding is not suppressed silently.

## Formatting and analysis

- `gofmt` is authoritative. Committed Go files must produce no output from the formatting check.
- `go vet ./...` and pinned Staticcheck must pass with CGO disabled.
- No function may have cyclomatic complexity greater than 15.
- No duplicate block of 100 tokens or more may remain without an approved exception.
- New exported identifiers require useful Go documentation. Raw stack traces, secrets, and unbounded outputs must not cross a protocol boundary.

## Tests

Production behavior is developed test-first: add the smallest focused failing test, confirm its intended failure, implement the behavior, rerun the focused test, then run `make quality`. Tests must be deterministic and must not depend on real project paths, ambient credentials, network access, sealed labels, or validation and held-out outcomes.

Temporary filesystem and process tests use per-test directories and bounded contexts. Concurrent code receives race-detector coverage in a dedicated supported environment; the race build is not the production artifact and must not weaken the CGO-disabled production policy. Failures at process or persistence boundaries must be represented as typed terminal states once those layers exist.

## Maintainability review

Reviewers check that each change:

1. belongs to the current phase and does not scaffold deferred learning, amendments, shared state, a second runtime, or UI;
2. preserves the single `act` surface and does not add standing instructional prose;
3. keeps protocol decisions in structured data rather than rendered text;
4. has one clear package owner, no import cycle, and no avoidable global mutable state;
5. keeps functions normally below 80 logical lines and files below 500 logical lines;
6. handles errors at the layer that can add useful context and never ignores an error without a documented reason;
7. updates tests, schemas, decision records, and the build-manifest inputs when their owned behavior changes.

Limits in this section are review triggers, not permission to split code mechanically or hide duplication. Reviewers should prefer a cohesive exception over a misleading refactor.

## Exceptions

An exception requires a checked-in decision or review record with all fields below. The author of the exception cannot be its only approver.

```yaml
rule:
scope:
rationale:
risk:
compensating_control:
owner:
reviewer:
approved_on:
expires_on_or_condition:
```

Expired exceptions fail review. Permanent exceptions must state why an expiry condition cannot be made meaningful.

## Make targets and exit status

Every target returns `0` only when all named checks pass and a nonzero status on tool, validation, test, or build failure.

| Target | Contract |
| --- | --- |
| `make test` | Runs the root Go tests once with CGO disabled. |
| `make lint` | Checks formatting, `go vet`, and pinned Staticcheck. |
| `make build` | Produces a stripped, path-trimmed Linux-amd64 daemon at `bin/phoenix`. |
| `make validate-spec` | Runs the accepted Draft 2020-12 validator over checked-in result and world examples. |
| `make quality` | CI entry point: tests, lint, dependency verification and vulnerability scan, complexity and duplication checks, schema validation, build, and manifest generation. |

`make manifest` writes `build/manifest.json`. From Task 1.3 onward it hashes the daemon, accepted schemas, the typed registry/executor boundary, and every source file in the registered dev-repo verb set. It still records an empty authored-rule set plus absent world-definition and active-weight artifacts until those components are assembled. Gate 1A freezes the complete manifest before validation opens.
