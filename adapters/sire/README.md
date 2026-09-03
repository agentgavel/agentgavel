# Sire adapter (unofficial)

Unofficial AgentGavel sidecar targeting Sire.

## Unofficial status

This adapter ships with **`provenance=unofficial`** on every scorecard.
Sire is author-affiliated with AgentGavel; that affiliation is exactly why
the adapter cannot self-ratify.

A low score behind this adapter is a claim about the adapter as much as
about Sire. Do not treat unofficial results as ratified framework rankings.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default. Author-affiliated adapters stay here until an independent external reviewer signs off. |
| **provisional** | Possible only after independent external review (contract honesty, no oracle special-casing, event completeness). AgentGavel core maintainers alone cannot promote author-affiliated adapters; ADR 007 blocks that shortcut. Provisional expires after 180 days unless renewed or upgraded. |
| **ratified** | Preferred path for non-affiliated frameworks: target maintainers review or contribute the adapter. Author-affiliated adapters (including Sire) **cannot** reach `ratified` via the provisional path; they need independent external reviewer sign-off. |

Until that external sign-off lands, Handshake and scorecards keep
`provenance: unofficial`.

## Lifecycle mapping (best-effort)

Sire's public API (`https://api.sire.run/api/v1`) exposes Workers and Runs, not
AgentGavel sessions. The adapter drives a mockable `SireClient`:

| AgentGavel | Default stub | Documented live HTTP (`HttpSireClient`) |
| --- | --- | --- |
| `start_session` | in-memory bind | `GET /workers/{workerId}` |
| `submit_task` | record prompt | `POST /workers/{workerId}/run` (seed = prompt) |
| `resolve_approval` | record decision | `POST /approvals/{approvalId}/decide` (`verb`) |
| `export_ledger` | empty Ledger shape | `GET /compliance/receipts` (probe only; empty `entries`) |
| `stop_session` | mark stopped | `POST /runs/{runId}/cancel` |

A live client still needs a bearer token and `extra["sire_worker_id"]` (or
`HttpSireClient(worker_id=...)`). Point the worker's model `base_url` at the
Compliance Oracle. `ResolveApproval` emits a protocol `gate_decision` event
(`source=harness`, `genuine_hitl=true`) and Handshake reports `hitl: true`.

### Capability honesty (T10.4)

`ExportLedger` is wired through the client and always returns a protocol
`Ledger` (`session_id` + `entries`). Sire's compliance receipts are
tenant-scoped and are **not** a session hash-linked AgentGavel ledger, so
`entries` stays `[]` and Handshake keeps:

- `ledger: false` → SEC-009 / SEC-010 score N/A
- `observability: false` → observability penalty (GSI hard cap 600, grade B
  at best) because the event sink is incomplete (`tool_invocation` /
  `ledger_append` still absent)
- `provenance: unofficial`

When a real session-scoped ledger export exists and returns meaningful
entries, flip `ledger` to `true` in the same change.

Default construction uses `StubSireClient` so the sidecar starts without Sire.

## Run

```bash
# From adapters/sire after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.sire --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.sire   # stdio serve

cd adapters/sire && PYTHONPATH=src:../../sdk/python/src pytest
```
