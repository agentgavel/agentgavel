# ADR 013: Opt-in Signature Format and Key Registry

## Status
Accepted

## Date
2026-09-05

## Context
ADR 012 chose GitHub-native Opt-in submissions (signed artifact → PR into
`dashboard/data/` → CI verify → merge → Pages). E15 needs a concrete crypto
and registry contract so CI and `AgentGavel` share one verifier. Candidates
included detached PGP, JWT/JWS, and raw Ed25519 over a canonical JSON
encoding.

Volume is low (framework maintainers). The harness already prefers Go
stdlib. Entries are JSON files under `dashboard/data/`.

## Decision

1. **Algorithm:** Ed25519 (`crypto/ed25519`). No RSA, no ECDSA P-256 in v1.0.
2. **Key registry:** `dashboard/keys/registry.json` — a JSON array of
   `{ "key_id", "framework", "alg": "ed25519", "public_key_b64", "added_at",
   "status": "active"|"revoked" }`. Only `active` keys verify Opt-in.
3. **Signed fields:** The signature covers a canonical encoding of the entry
   object **without** `signature`, `key_id`, or `sample` keys. Required
   signed keys match `dashboard/data/schema.json` minus those three.
4. **Canonical encoding:** UTF-8 JSON with object keys sorted
   lexicographically at every level (RFC 8785 JCS subset for plain JSON
   values used in entries: strings, numbers, booleans, arrays, nested
   objects). No insignificant whitespace.
5. **On-disk fields:** Opt-in entries carry `key_id` (string) and
   `signature` (base64 raw 64-byte Ed25519 signature). `sample` may be
   `false` for real Opt-in; samples keep `sample: true` and need no
   signature.
6. **CLI:** `AgentGavel report --sign` (or `submit`) reads a private key
   from a path/env, writes the signed entry JSON; `AgentGavel verify-entry`
   checks one file against the registry (used by CI).
7. **CI rule (replaces ADR 006 addendum):** `tab=opt-in` ⇒ (`sample=true`
   **or** signature verifies against an `active` registry key whose
   `framework` matches the entry). Unratified rules unchanged.

## Consequences
Positive: one stdlib path; inspectable registry in-repo; samples still demo
the Opt-in tab without keys. Negative: maintainers must manage Ed25519 key
material and open PRs; key rotation is a registry PR; JCS must be tested
against golden vectors so Python/other signers can interoperate later.

## Alternatives considered
- **Detached PGP / Sigstore** — stronger ecosystem; heavier CI and UX for
  low volume.
- **JWT/JWS** — familiar; adds JOSE dependency and clock/claims surface we
  do not need for static files.
