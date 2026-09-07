# Fallback Discipline Contract

> Extracted development contract. Canonical route: [`../../DEVELOPMENT.md`](../../DEVELOPMENT.md).


## The Boundary Rule

Fail-loud discipline applies strictly at **untrusted external boundaries**, where data crosses between external domains (uncontrolled) and internal runtime logic (controlled):

- **Model Provider Wire:** Serialization and deserialization of tool calls, thinking blocks, and message payloads.
- **Untyped Tool Arguments:** Dynamic model-generated JSON / `map[string]any` inputs.
- **Auth & Credential Exchanges:** OAuth tokens, refresh responses, and account identifiers.
- **Platform Isolation Gates:** Linux namespaces, FUSE mounts, and sandboxing primitives.
- **Cross-Process / Persistence Surfaces:** Authoritative state files (`current.jsonl`, `session.json`, exit markers).

**Rule:** At an external boundary, never absorb malformed, missing, or unrecognized data with a synthetic default that masks the failure. If the payload is invalid or unsupported, reject or surface the error immediately.

## Internal Trusted Logic: Idiomatic & Resilient

Inside internal domain code (fully typed, repository-controlled):
- Standard Go zero-values (`""`, `0`, `nil`, empty slices/maps) are natural and expected.
- Resilient fallbacks, optional config defaults, and graceful degradation are encouraged.
- Do not introduce panics or dogmatic error propagation where an idiomatic default value is safe and well-defined.

## Key Boundary Anti-Patterns & Antidotes

1. **Untyped Tool-Arg Coercion:** Never use unchecked type assertions like `v := args[k].(T)`. Use present-aware helpers (e.g. `IntArg`, `StringArg`) to distinguish *absent* (apply valid default) from *present-but-wrong-type* (reject call).
2. **Provider Wire Serialization:** Never discard JSON marshal/unmarshal errors on API wire paths. Ensure zero-value maps/slices serialize correctly according to provider schema (e.g., `{}` rather than `"null"` for tool arguments).
3. **Auth Refresh Overwrite:** An external refresh response omitting optional fields must never overwrite existing valid credentials with empty values.
4. **Platform Isolation Fallbacks:** If a requested security or isolation primitive is unsupported on the host platform, fail closed with a clear configuration error rather than silently running unisolated.
5. **Authoritative Persistence:** Never discard a write error for a canonical
   cross-process surface. Propagate the error, and publish concurrently read
   state atomically so readers cannot observe a partial record.
6. **Unknown Boundary Cases:** Do not use a silent `default` branch to absorb an
   unrecognized external wire or enum case; reject it or surface the drift.
