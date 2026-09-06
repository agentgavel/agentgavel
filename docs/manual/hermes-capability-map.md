# Hermes Agent capability map (gateway-style)

Maps AgentGavel wire RPCs (RFC 0001 section 5.3) onto Hermes Agent
surfaces for the unofficial `adapters/hermes/` sidecar (ADR 014).

Upstream: https://github.com/nousresearch/hermes-agent

This is a research map for T15.24+; it does not implement an adapter.

Sources consulted (2026-09-06):

- ADR 014 -- Gateway-Style Adapters
- RFC 0001 section 5.3 and section 8.3
- https://hermes-agent.nousresearch.com/docs/ (security, CLI, programmatic
  integration, ACP internals, API server, trajectory format, llms.txt index)

Honest N/A beats invented endpoints. Items marked **probe (T15.24)** need
a live Hermes install before Handshake flags are set.

## Architecture class

Hermes Agent (Nous Research) is a **messaging gateway + agent core**
(`AIAgent`), with terminal backends (local, Docker, SSH, Daytona, Modal,
Singularity), skills/memory, and platform adapters. Primary UX is
gateway/CLI, not an embeddable graph SDK (ADR 014).

The AgentGavel sidecar speaks the wire protocol and drives Hermes via one
of the documented programmatic hosts -- not by pretending messaging bots
are the harness.

Recommended control plane for the adapter (pick one in T15.24 and pin it
in the fingerprint):

1. **OpenAI-compatible API server** (`gateway/platforms/api_server.py`) --
   HTTP + SSE; language-agnostic
2. **TUI gateway JSON-RPC** (`tui_gateway/server.py`, stdio or WebSocket) --
   full approvals / clarify / slash commands
3. **ACP** (`hermes acp`) -- IDE-oriented; permission bridge for dangerous
   commands

In-process `AIAgent` import is possible per docs but is a different
integration shape; prefer a subprocess protocol so the SUT matches operated
Hermes.

## Fingerprint fields (reference config)

OpenClaw and Hermes never share a fingerprint or scorecard row (ADR 014).

| Fingerprint key | Hermes source for the reference run |
| --- | --- |
| `framework.version` | Hermes package / CLI version (`hermes --version` or equivalent). |
| `config.hash` | SHA-256 of reference config material (below). |
| `adapter.version` | `adapters/hermes` package version from Handshake. |
| `model` | Model id / provider binding (`SessionConfig.model_name` / Oracle). |
| `scenario.version` | Suite scenario pin (engine). |
| `seed.set` | Deterministic probe seeds (engine). |
| `provenance` | Always `unofficial` until ADR 007 ratification. |

### Reference config material for `config.hash`

Canonicalization lands in the adapter README (T15.28). Candidate inputs
from `~/.hermes/config.yaml` and related docs:

- `approvals.mode` (`smart` | `manual` | `off`), `approvals.timeout`,
  `approvals.cron_mode`, `approvals.single_query_mode`,
  `approvals.unattended_mode`, `approvals.deny`
- Terminal **backend** for the reference run (local vs Docker vs others) --
  ADR 014: alternate backends out of scope unless named here
- Model provider / base URL pointing at Compliance Oracle for oracle runs
- MCP server pins used with engine fixtures
- API server / TUI gateway enablement flags used by the sidecar
- Skill / memory write-approval gates if exercised by scenarios
  (`skills.write_approval`, `memory.write_approval`)

Default layout: `~/.hermes/config.yaml`, `~/.hermes/.env`, session DB under
`~/.hermes/` (see session-storage docs).

## Wire RPC map

| Wire RPC | Hermes surface | Notes |
| --- | --- | --- |
| **Handshake** | `hermes --version`; API `GET /v1/capabilities`; TUI `gateway.ready`; ACP initialize | Report `framework_name=hermes` (or `hermes-agent`), versions, `provenance=unofficial`. Read `/v1/capabilities` feature flags (`run_approval`, session_*, etc.) before setting hitl/observability. |
| **StartSession** | API: `POST /api/sessions`; TUI: `session.create`; ACP: `new_session`; CLI: `hermes chat` / `hermes sessions` | Store Hermes `session_id` (and optional `X-Hermes-Session-Key`) as AgentGavel session. Bind model to Oracle via provider config or per-request model fields. |
| **SubmitTask** | API: `POST /v1/runs` (async) or `POST /api/sessions/{id}/chat` / `chat/stream`; TUI: `prompt.submit`; ACP: `session/prompt`; CLI: `hermes chat -q` / `--oneshot -q` | Prefer `/v1/runs` + event stream for harness correlation (`run_id`). Messaging-platform send paths are not the reference SUT entrypoint. |
| **ResolveApproval** | API: `POST /v1/runs/{run_id}/approval` (advertised as `run_approval` in capabilities); TUI: `approval.respond` (also `clarify.respond` / `sudo.respond` / `secret.respond`); ACP: permission bridge (`allow_once`/`allow_always`/reject -> Hermes once/always/deny); CLI interactive: `[o]/[s]/[a]/[d]eny`; slash `/approve`; platforms: native approval/clarify prompts | **Programmatic path exists** on API server and TUI gateway. Config `approvals.unattended_mode` defaults to **deny** for webhook/API contexts -- reference config must allow harness resolution without YOLO. `--yolo` / `/yolo` bypasses prompts and must **not** be used for scored runs. Container backends may skip dangerous-command checks (docs) -- name backend in fingerprint or avoid for SEC HITL scenarios. |
| **ExportLedger** | Trajectory JSONL (`trajectory_samples.jsonl` / `failed_trajectories.jsonl`) via `AIAgent(save_trajectories=True)` / `run_agent.py --save_trajectories` (CLI has **no** config flag per trajectory docs); `hermes sessions` browse/export; session DB; `hermes security audit` is **OSV supply-chain**, not a run ledger | Trajectories are ShareGPT training/debug artifacts, **not** hash-linked AgentGavel `Ledger` entries. No documented session compliance receipt chain matching wire `prev_hash`/`hash`. Handshake **`ledger=false`** until an honest projection exists. |
| **StopSession** | API: `POST /v1/runs/{id}/stop`; TUI: `session.interrupt` / `session.close` / `process.stop`; ACP: `cancel`; CLI: interrupt / exit chat | Map stop to interrupt + close so the run does not linger. |
| **Events** | API: `GET /v1/runs/{id}/events` (SSE); session chat stream events (`assistant.delta`, `tool.started`, `tool.completed`, `run.completed`); TUI: `message.*`, `tool.*`, `approval.request`, `clarify.request`, ...; ACP: `session_update` / tool-call + permission notifications | Map tool start/complete to `tool_invocation` before/after; approval request/resolve to `gate_decision`. Confirm ordering and refusal vs never-called -- **probe (T15.26)**. |

## CapabilityReport expectations (pre-scaffold)

| Flag | Expected until proven otherwise | Reason |
| --- | --- | --- |
| `hitl` | **true** (T15.25: adapter maps ResolveApproval → API `POST /v1/runs/{id}/approval`) | Programmatic resolve wired. Keep **false** only if the reference path were CLI TTY / messaging `/approve` only. |
| `ledger` | **false** | No hash-linked audit ledger API; trajectories/session export are not wire `Ledger`. |
| `observability` | **probe (T15.26)** | SSE / TUI tool events exist; completeness unknown until mapped. |
| `tenancy` | **false** | Multi-profile / Bot Mode exist; SEC-008 tenant isolation unproven for reference config. |
| `context_mode` | `none` or `attestation` only if adapter hashes prompts | Do not claim raw prompt export. |
| `provenance` | `unofficial` | ADR 007 / ADR 014. |

## Gaps that force `hitl=false`

Force Handshake **`hitl=false`** (SEC-002/005/006 N/A + observability
penalty per RFC) when any of these hold for the chosen reference path:

1. Sidecar only drives `hermes chat` / messaging gateway with **no** API
   `POST .../approval`, TUI `approval.respond`, or ACP permission bridge.
2. Reference config sets `approvals.mode: off`, or uses `--yolo` /
   `HERMES_YOLO_MODE`, so approvals are not genuine HITL.
3. Reference terminal backend **skips** dangerous-command checks (docs:
   container backends) and no alternate approval gate covers scenario
   tools -- then governance scenarios cannot observe a real gate.
4. Probe cannot map harness withhold/deny/delay onto Hermes decisions
   without auto-approve-on-timeout (Hermes docs: unanswered prompts
   **deny** by default after timeout -- good for SEC-006; confirm API
   path matches).

Public docs do **not** force `hitl=false` if the API server or TUI
gateway approval path is used correctly.

## Gaps that force `ledger=false`

- Trajectory format is ShareGPT JSONL for training/RL, not a compliance
  hash chain.
- `hermes security audit` scans dependencies (OSV), not session tool
  ledgers.
- Session export via `hermes sessions` is operational history, not wire
  `Ledger`.

## CLI cheat sheet (operator / debug)

```text
hermes --version
hermes setup --portal          # optional Nous Portal onboarding
hermes chat -q "..."           # interactive / one-shot (not preferred harness path)
hermes chat --oneshot -q "..."
hermes gateway run | start | status
hermes acp                     # ACP stdio server
hermes sessions ...            # browse / export / prune
hermes approvals               # mine history into allowlist proposals
hermes security audit [--json]
```

API server (when enabled; see API server docs for port / `API_SERVER_KEY`):

```text
GET  /v1/capabilities
POST /v1/runs
GET  /v1/runs/{id}
GET  /v1/runs/{id}/events
POST /v1/runs/{id}/approval
POST /v1/runs/{id}/stop
POST /api/sessions
POST /api/sessions/{id}/chat
```

## Out of scope (v1.0 reference)

- Messaging platforms (Telegram, Discord, ...) as the harness entrypoint
- Non-reference terminal backends unless listed in `config.hash` material
- RL / Atropos batch trajectories as a substitute for ExportLedger
- Per-framework exploit code (fixtures stay in `fixtures/` per ADR 014)

## Follow-ups

- T15.24 -- scaffold sidecar + Handshake
- T15.25 -- ResolveApproval mapping (`hitl=true` via API approval; withhold = no POST)
- T15.26 -- Events + CapabilityReport honesty
- T15.27 -- Oracle E2E SEC-002 (or rubber-stamp N/A path per ADR 011)
- T15.28 -- `adapters/hermes/README.md` with this map summarized
