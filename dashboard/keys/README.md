# Maintainer key registry (ADR 013)

Opt-in leaderboard rows with `sample=false` must carry an Ed25519 signature
that verifies against an **active** entry in `registry.json`.

## Register a key

1. Generate an Ed25519 keypair (32-byte seed / 64-byte private key).
2. Open a PR that appends an object to `registry.json`:

```json
{
  "key_id": "your-org-2026-09",
  "framework": "Exact Framework Display Name",
  "alg": "ed25519",
  "public_key_b64": "<base64 of 32-byte public key>",
  "added_at": "2026-09-05T00:00:00Z",
  "status": "active"
}
```

3. After merge, sign scorecards with `AgentGavel report --sign` and open a
   data PR per `docs/manual/opt-in-submission.md` (Wave 31).

## Status

- `active` -- may verify Opt-in submissions
- `revoked` -- rejected by verify; keep the row for audit

The `framework` field on the key must match the entry's `framework` field
exactly (no sample suffix for real Opt-in rows).

## Fixture key

`example-framework-test-1` is a **test** key for CI and unit tests. Do not use
it for production Opt-in rows. The matching private key lives only under
`internal/submit/testdata/` (tests).
