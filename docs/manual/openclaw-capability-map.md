# OpenClaw capability map (gateway-style)

Maps AgentGavel wire RPCs (RFC 0001 section 5.3) onto OpenClaw Gateway /
CLI / config surfaces for the unofficial `adapters/openclaw/` sidecar
(ADR 014). This is a research map for T15.17+; it does not implement an
adapter.

Sources consulted (2026-09-06):

- ADR 014 -- Gateway-Style Adapters
- RFC 0001 section 5.3 (adapter contract) and section 8.3
- https://docs.openclaw.ai/ (Gateway protocol, exec approvals, CLI,
  configuration, audit ledger RPC)

Honest N/A beats invented endpoints. Items marked **probe (T15.17)** need
a live Gateway before Handshake flags are set.

## Architecture class

OpenClaw is a self-hosted **Gateway**: channels and operator clients talk
to one Gateway process; the Gateway owns sessions, tool policy, exec
approvals, and audit. The AgentGavel sidecar speaks the wire protocol and
drives OpenClaw via Gateway WebSocket RPC (primary), optional Admin HTTP
RPC plugin, and CLI wrappers -- not by importing an in-process agent
graph library (ADR 014).

Nested Codex / Claude / Copilot (or other) vendor harness plugins are
**out of scope** for the v1.0 reference config unless
`CapabilityReport` / fingerprint explicitly names them as the execution
backend.

## Fingerprint fields (reference config)

Run fingerprints use RFC section 4.11 / `internal/engine` keys. OpenClaw
must never share a fingerprint or scorecard row with Hermes (ADR 014).

| Fingerprint key | OpenClaw source for the reference run |
| --- | --- |
| `framework.version` | OpenClaw package / Gateway version (`openclaw --version` or Gateway health/status RPC). |
| `config.hash` | SHA-256 of the **reference** config material (below), not the whole home directory. |
| `adapter.version` | `adapters/openclaw` package version from Handshake. |
| `model` | Model id bound for the run (`SessionConfig.model_name` / Oracle `base_url`). |
| `scenario.version` | Suite scenario pin (engine). |
| `seed.set` | Deterministic probe seeds (engine). |
| `provenance` | Always `unofficial` until ADR 007 ratification (Handshake `CapabilityReport`). |

### Reference config material for `config.hash`

Hash only the fields that change governance behavior for the scored run.
Document the exact canonicalization in the adapter README (T15.21).
Candidate inputs (names from OpenClaw docs / `~/.openclaw/openclaw.json`):

- `agents.defaults` model / workspace bindings used by the reference agent
- `tools.exec.*` (host, security, ask / askFallback) -- exec approval posture
- Host exec-approvals policy snapshot as returned by `exec.approvals.get`
  or `openclaw approvals get` / `openclaw exec-policy show`
- `logging.audit.enabled`, `logging.audit.messages`,
  `logging.audit.executionIdentity` -- audit ledger posture
- Sandbox / node binding for the reference session (`tools.exec.host`,
  session `execNode` if used)
- Gateway auth / operator-scope assumptions for the harness client (not
  secrets; scope names only)

Default state path: `~/.openclaw/openclaw.json` (override via
`OPENCLAW_CONFIG_PATH` / `--profile`).

## Wire RPC map

Primary control plane: Gateway WebSocket RPC
(https://docs.openclaw.ai/gateway/protocol). CLI equivalents wrap the same
methods where listed. Optional `POST /api/v1/admin/rpc` (admin-http-rpc
plugin, **disabled by default**) allowlists a subset of methods -- prefer
WebSocket for the sidecar unless probe shows otherwise.

| Wire RPC | OpenClaw surface | Notes |
| --- | --- | --- |
| **Handshake** | `openclaw --version` / `openclaw gateway status` / `openclaw health`; Gateway `connect` + advertised method list; optional `agents.list` | CapabilityReport: `framework_name=openclaw`, version fields, `provenance=unofficial`. Set `hitl` / `ledger` / `observability` only after probe confirms the mappings below. |
| **StartSession** | Gateway `sessions.create` (optional nested initial message); CLI: `openclaw sessions` family / Control UI session create | Map returned session key / id into AgentGavel `SessionId`. Point model binding at Compliance Oracle via OpenClaw model config for the reference agent (`SessionConfig.model_base_url`). MCP fixtures: OpenClaw `mcp` CLI / Gateway MCP config -- **probe (T15.17)** for mounting engine-started fixture endpoints. |
| **SubmitTask** | Gateway `sessions.send` or `chat.send` with task prompt; CLI: `openclaw message send` / agent chat paths | Prefer session-scoped send so Events and approvals correlate to one `sessionKey`. `agent.wait` can block for terminal snapshot after submit. |
| **ResolveApproval** | Gateway `approval.resolve` (kind-agnostic) or `exec.approval.resolve`; events `exec.approval.requested` / `session.approval` via `sessions.messages.subscribe` with `includeApprovals: true`; CLI: channel `/approve`, Control UI; policy: `openclaw approvals get|set`, `openclaw exec-policy show|set|preset` | Requires operator scope `operator.approvals` (write scope does **not** subsume it). First-answer-wins. Intended map: approve→`allow-once`, deny→`deny`; **withhold unmapped** (no Gateway enum). **T15.18 probe:** no live Gateway in reference path + withhold gap → Handshake **`hitl=false`** (documented N/A); `ResolveApproval` raises `HitlNotSupportedError` (SEC-002/005/006 N/A, never silent Fail). |
| **ExportLedger** | Gateway `audit.activity.list` (preferred) or legacy `audit.list`; CLI: `openclaw audit`; related: `audit.run.inspect`, task ledger `tasks.list` | Activity ledger is **metadata-only** (no prompts, tool args, or bodies), 30-day / 100k-record cap, best-effort (may drop). Not a hash-linked AgentGavel `Ledger` (`prev_hash` / `hash`). Until an honest projection exists, Handshake **`ledger=false`** (SEC-009/010 N/A). Do not invent hash linkage. |
| **StopSession** | Gateway `sessions.abort` (cancel active work); archive via `sessions.patch` with `archived: true` + `expectedSessionId`; optional `sessions.delete` / `sessions.reset` | Prefer abort then archive for clean teardown. Confirm which call leaves no live run -- **probe (T15.17)**. |
| **Events** | Gateway event frames: `session.tool`, `session.message`, `session.operation`, `session.approval`, `exec.approval.requested` / `exec.approval.resolved`, `plugin.approval.*`; subscribe: `sessions.subscribe`, `sessions.messages.subscribe` | Map to wire `tool_invocation` (before/after), `gate_decision`, `session_error`. `ledger_append` only if audit rows are projected. Best-effort delivery (slow subscribers may drop frames) -- affects `observability` honesty. |

## CapabilityReport expectations (pre-scaffold)

| Flag | Expected until proven otherwise | Reason |
| --- | --- | --- |
| `hitl` | **false** (T15.18 documented N/A) | Public docs expose `approval.resolve` / `exec.approval.resolve`, but reference sidecar has no live Gateway and `withhold` has no resolve enum. Flip to **true** only after a probe drives approve/deny/withhold as the harness human (ADR 014). |
| `ledger` | **false** | Audit RPC is metadata-only and not AgentGavel hash-linked ledger shape. |
| `observability` | **probe (T15.19)** | Tool and approval events exist; completeness vs before/after `tool_invocation` and drop semantics unknown until wired. |
| `tenancy` | **false** unless multi-agent isolation is claimed | Multi-agent routing exists; cross-tenant SEC-008 mapping is unproven -- do not set true without evidence. |
| `context_mode` | `none` or `attestation` only if adapter hashes prompts | Prefer attestation (ADR 005) if emitting context events; do not claim raw context export without a surface. |
| `provenance` | `unofficial` | ADR 007 / ADR 014. |

## Gaps that force `hitl=false`

**T15.18 (2026-09-06):** Handshake keeps **`hitl=false`** (documented N/A):

1. No live OpenClaw Gateway / `openclaw` CLI in the reference sidecar path
   (cannot obtain `operator.approvals` or call resolve RPCs).
2. Gateway resolve enums are only `allow-once` / `allow-always` / `deny` —
   AgentGavel **`withhold` has no mapping** (approve→`allow-once`,
   deny→`deny` documented in `adapters/openclaw` approvals helpers).
3. Approvals that only work via human chat `/approve` without RPC would
   also force N/A (not exercised here — no Gateway).

Never stub green SEC-002: `ResolveApproval` must raise loudly while
`hitl=false` so scenarios score N/A via `ScenarioNA`, not silent Fail.

## Gaps that force `ledger=false`

- No public hash-chained session compliance ledger matching wire `Ledger`.
- `audit.activity.list` is explicitly metadata-only and best-effort.

## CLI cheat sheet (operator / debug)

```text
openclaw --version
openclaw gateway status | health | call
openclaw sessions ...
openclaw approvals get | set
openclaw exec-policy show | set | preset
openclaw audit
openclaw config get | set | validate
openclaw dashboard          # Control UI, default http://127.0.0.1:18789/
```

## Out of scope (v1.0 reference)

- Channel plugins as the SUT (WhatsApp, Telegram, etc.) -- Gateway policy
  plane is the SUT
- Cloud-worker / node desktop observe paths unless named in fingerprint
- Admin HTTP RPC as the only transport (optional convenience)
- Per-framework exploit code (fixtures stay in `fixtures/` per ADR 014)

## Follow-ups

- T15.17 -- scaffold sidecar + Handshake against a local Gateway
- T15.18 -- ResolveApproval mapping or honest `hitl=false` (**done:** N/A)
- T15.19 -- Events + CapabilityReport honesty
- T15.20 -- Oracle E2E SEC-002 (or rubber-stamp N/A path per ADR 011)
- T15.21 -- `adapters/openclaw/README.md` with this map summarized
