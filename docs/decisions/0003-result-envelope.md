# 0003: Result envelope and world definition

- **Status:** Accepted for Phase 0 review
- **Date:** 2026-08-17
- **Schemas:** `spec/result.schema.json`, `spec/world.schema.json`

## Decision

Every verb invocation returns one versioned result envelope. Calls and alternatives are structured objects; execution, reachability, and attribution never parse rendered text. Unknown fields are rejected throughout protocol-owned objects.

The envelope distinguishes four outcomes:

- `ok` contains a typed result and no error or refusal;
- `fail` contains a typed error for an attempted verb;
- `refused` contains a teaching refusal and a reachable structured alternative or an explicit reason that no alternative exists;
- `absent` contains the generic `absent` error, no frontier, and no refusal that could reveal whether a reference was guessed, revoked, cross-session, or nonexistent.

Every result includes the complete world-build digest, session and act identifiers, handle and verb, compact rendered text, explicit handle grants/revocations, and a frontier capped at three entries. Each frontier call contains a reachable handle, verb, fully bound arguments, and an optional state digest that is revalidated at invocation.

Opaque handle references are session-scoped. World files name root bindings for authors, but the daemon issues opaque references when a session starts. A result changes reachability only through `handles.grant` and `handles.revoke`.

## World definitions

A world definition declares:

- named root bindings and their live resource kind;
- handle types and the verbs attached to each type;
- argument and result schemas per verb;
- handle types that each verb may grant or revoke;
- state-sensitive refusal rules;
- authored transition rules with bound call templates.

Bindings may use a literal, a JSON Pointer into the typed result, or a JSON Pointer into live state. A target handle comes from the current handle, a named root, or an explicitly granted handle in the result. The loader must reject missing handle types, verbs not attached to the selected type, invalid JSON Pointers, unbound arguments, and result grants not declared by the verb. Those are graph-semantic checks beyond what JSON Schema alone can express.

A refusal `when` schema evaluates a structured feature object. `args` is always present; `state` and `state_digest` are present after live resolution; `failure` is present only when evaluating a typed verb failure. Request/state refusals run before the handler, and failure-shaped refusals run before an ordinary `fail` envelope is returned. Human-readable result text is never a refusal input.

## Canonicalization and digests

Protocol JSON is canonicalized with RFC 8785 before hashing. Digests use lowercase SHA-256 encoded as `sha256:` plus 64 hexadecimal characters.

A world-build digest covers an ordered manifest containing:

1. daemon executable digest;
2. registered verb implementation digests;
3. result and world schema digests;
4. world-definition digest;
5. authored transition and refusal digests;
6. active frontier-weight digest.

Episode records store the world-build digest and the canonical digest of each request and result. Rendered `text` remains part of the result envelope and its digest, but no machine decision depends on parsing it.

## Compatibility

Version 1 rejects unknown fields rather than silently accepting a newer meaning. A future incompatible change increments `v` and ships a new schema. Readers may support multiple versions explicitly; they may not reinterpret a version 1 field.

## Verification

`go run ./cmd/validate-spec --repo-root ../..` compiles both Draft 2020-12 schemas and validates the checked-in success, refusal, absence, and world examples. Phase 1 adds semantic loader tests for cross-reference and reachability rules.
