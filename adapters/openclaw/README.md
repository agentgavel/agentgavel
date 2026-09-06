# OpenClaw adapter (unofficial)

Unofficial AgentGavel **gateway-style** sidecar (ADR 014) targeting
OpenClaw's Gateway policy plane. Handshake and scorecards always report
**`provenance=unofficial`** until ADR 007 ratification. Do not treat
unofficial results as ratified rankings.

Handshake scaffold (T15.17), ResolveApproval honest N/A (T15.18), and
Event helpers / ExportLedger honesty (T15.19). Capability map:
[`docs/manual/openclaw-capability-map.md`](../../docs/manual/openclaw-capability-map.md).

## Unofficial status

Provenance is always `unofficial` (ADR 007 / ADR 014). Nested Codex /
Claude / Copilot (or other vendor harness) plugins inside OpenClaw are
**out of scope for the v1.0 reference config** unless
`CapabilityReport` / the run fingerprint explicitly names them as the
execution backend (ADR 014). The SUT is the Gateway policy plane in the
documented reference config — not channel plugins or nested vendor
harnesses.

## Capability honesty (T15.18 / T15.19)

| Flag | Value | Why |
| --- | --- | --- |
| `hitl` | `false` | **Documented N/A.** No live Gateway in the reference path; AgentGavel `withhold` has no `exec.approval.resolve` enum (`allow-once` / `allow-always` / `deny` only). |
| `ledger` | `false` | `audit.activity.list` is metadata-only, not a hash-linked AgentGavel Ledger |
| `observability` | `false` | No live Gateway subscription; Event helpers exist but the stream is incomplete |
| `tenancy` | `false` | Unproven for reference config |
| `context_mode` | `none` | No prompt attestation yet |
| `framework_name` | `openclaw` | |
| `framework_version` | `unprobed` | Set after live `openclaw --version` / Gateway health probe |

`ExportLedger` returns `{session_id, entries: []}`. Helpers
`emit_tool_invocation` / `emit_gate_decision` / `ingest_gateway_frame` map
Gateway frames when available; scaffold lifecycle does not invent Events.

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

## Fingerprint fields (reference config)

OpenClaw must never share a fingerprint or scorecard row with Hermes
(ADR 014). Run fingerprints use RFC §4.11 / `internal/engine` keys:

| Fingerprint key | OpenClaw source for the reference run |
| --- | --- |
| `framework.version` | OpenClaw package / Gateway version (`openclaw --version` or Gateway health/status RPC) |
| `config.hash` | SHA-256 of the **reference** config material below (not the whole home directory) |
| `adapter.version` | `adapters/openclaw` package version from Handshake |
| `model` | Model id bound for the run (`SessionConfig.model_name` / Oracle `base_url`) |
| `scenario.version` | Suite scenario pin (engine) |
| `seed.set` | Deterministic probe seeds (engine) |
| `provenance` | Always `unofficial` until ADR 007 ratification |

### Reference config material for `config.hash`

Hash only fields that change governance behavior for the scored run.
Canonical inputs (names from OpenClaw docs / `~/.openclaw/openclaw.json`):

- `agents.defaults` model / workspace bindings used by the reference agent
- `tools.exec.*` (host, security, ask / askFallback) — exec approval posture
- Host exec-approvals policy snapshot (`exec.approvals.get` /
  `openclaw approvals get` / `openclaw exec-policy show`)
- `logging.audit.enabled`, `logging.audit.messages`,
  `logging.audit.executionIdentity` — audit ledger posture
- Sandbox / node binding for the reference session (`tools.exec.host`,
  session `execNode` if used)
- Gateway auth / operator-scope assumptions for the harness client
  (scope names only — never secrets)

Default state path: `~/.openclaw/openclaw.json` (override via
`OPENCLAW_CONFIG_PATH` / `--profile`).

## Run against a local Gateway

Start OpenClaw Gateway locally (Control UI defaults to
`http://127.0.0.1:18789/`), then point the sidecar at stdio JSON-RPC:

```bash
# Operator / debug against a live Gateway
openclaw --version
openclaw gateway status   # or: openclaw health
openclaw approvals get
openclaw exec-policy show
openclaw audit

# Sidecar (stdio serve; scaffold does not require Gateway for Handshake)
cd adapters/openclaw
PYTHONPATH=src:../../sdk/python/src python -m adapters.openclaw --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.openclaw

# Harness
AgentGavel --adapter "python3 -m adapters.openclaw"
```

Until a live probe flips capability flags, Handshake stays conservative
(`hitl` / `ledger` / `observability` false). Missing capabilities score
**N/A** (never silent Fail).

## Test

```bash
cd adapters/openclaw
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
