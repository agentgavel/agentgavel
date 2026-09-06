# OpenClaw adapter (unofficial)

Unofficial AgentGavel sidecar targeting OpenClaw's Gateway policy plane
(ADR 014). Provenance is always `unofficial` until ADR 007 ratification.

<<<<<<< HEAD
=======
Handshake scaffold (T15.17) plus Event helpers / ExportLedger honesty
(T15.19). ResolveApproval / live Gateway HITL lands in T15.18. Capability
flags stay conservative until proven:

| Flag | Value | Honesty |
| --- | --- | --- |
| `hitl` | `false` | ResolveApproval not wired (T15.18) |
| `ledger` | `false` | `audit.activity.list` is metadata-only, not hash-linked AgentGavel Ledger |
| `observability` | `false` | No live Gateway subscription; Event helpers exist but stream incomplete |

`ExportLedger` returns `{session_id, entries: []}`. Helpers
`emit_tool_invocation` / `emit_gate_decision` / `ingest_gateway_frame` map
Gateway frames when available; scaffold lifecycle does not invent Events.

>>>>>>> bb7e26d (feat(adapters): OpenClaw Events helpers and ExportLedger honesty (T15.19))
Capability map: [`docs/manual/openclaw-capability-map.md`](../../docs/manual/openclaw-capability-map.md).

## Capability honesty (T15.18)

| Flag | Value | Why |
| --- | --- | --- |
| `hitl` | `false` | **Documented N/A.** No live Gateway in the reference path; AgentGavel `withhold` has no `exec.approval.resolve` enum (`allow-once` / `allow-always` / `deny` only). |
| `ledger` | `false` | Audit RPC ≠ hash-linked wire Ledger |
| `observability` | `false` | Events mapping deferred to T15.19 |
| `tenancy` | `false` | Unproven for reference config |
| `context_mode` | `none` | No prompt attestation yet |

`ResolveApproval` raises `HitlNotSupportedError` (never a no-op). With
`hitl=false`, SEC-002 / SEC-005 / SEC-006 score **N/A** via engine
`ScenarioNA` — not silent Fail (ADR 011 / ADR 014). Flip `hitl` to true
only after a live Gateway probe can drive approve/deny/withhold as the
harness human without auto-approve or silent timeout approve.

### Intended Gateway map (when hitl becomes true)

| AgentGavel Decision | OpenClaw `exec.approval.resolve` / `approval.resolve` |
| --- | --- |
| `approve` | `allow-once` (harness never uses `allow-always`) |
| `deny` | `deny` |
| `withhold` | **unmapped** — gap that keeps `hitl=false` until a documented surface exists |

Requires operator scope `operator.approvals`. Pure mapping helpers:
`adapters.openclaw.approvals` (no Gateway I/O).

## Run

```bash
PYTHONPATH=src:../../sdk/python/src python -m adapters.openclaw --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.openclaw   # stdio serve
```

Harness: `AgentGavel --adapter "python3 -m adapters.openclaw"`.

## Test

```bash
cd adapters/openclaw
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
