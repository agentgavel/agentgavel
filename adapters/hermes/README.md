# Hermes Agent adapter (unofficial)

Unofficial AgentGavel **gateway-style** sidecar for
[Hermes Agent](https://github.com/nousresearch/hermes-agent) (Nous Research).

Architecture class and scoring rules:
[ADR 014](../../docs/adr/014-gateway-style-adapters.md).
Capability map (wire RPC → Hermes surfaces):
[`docs/manual/hermes-capability-map.md`](../../docs/manual/hermes-capability-map.md).

Handshake and scorecards always report **`provenance=unofficial`** until
[ADR 007](../../docs/adr/007-adapter-ratification.md) ratification. Do not
treat unofficial results as ratified rankings. OpenClaw and Hermes never
share a fingerprint or scorecard row (ADR 014).

## Control plane (v1.0 reference)

| Path | Status |
| --- | --- |
| **OpenAI-compatible API server** (`gateway/platforms/api_server.py`) | **In** — reference control plane for this adapter |
| TUI gateway JSON-RPC (`tui_gateway/server.py`) | **Out** of v1.0 reference (documented alternative) |
| ACP (`hermes acp`) | **Out** of v1.0 reference (documented alternative) |
| Messaging platforms (Telegram, Discord, …) as harness entry | **Out** |

The sidecar speaks AgentGavel stdio JSON-RPC and drives Hermes over HTTP +
SSE: `GET /v1/capabilities`, `POST /v1/runs`, `GET /v1/runs/{id}/events`,
`POST /v1/runs/{id}/approval`, `POST /v1/runs/{id}/stop`.

## Terminal backends (v1.0 reference)

Hermes supports several terminal backends (local, Docker, SSH, Daytona,
Modal, Singularity). Only the **local** backend is in the v1.0 reference
config.

| Backend | Status |
| --- | --- |
| **local** | **In** — reference execution backend |
| Docker | **Out** unless named in `config.hash` / CapabilityReport |
| SSH | **Out** unless named |
| Daytona / Modal / Singularity | **Out** unless named |

Container backends may skip dangerous-command checks (Hermes docs). Do not
use them for scored SEC HITL runs unless the fingerprint explicitly names
the backend and an alternate approval gate still covers scenario tools
(capability map / ADR 014).

## Run against a local Hermes API-server gateway

1. Install Hermes Agent upstream and configure a provider (see
   [hermes-agent](https://github.com/nousresearch/hermes-agent) and
   `hermes setup --portal` if using Nous Portal).

2. Enable the API server in `~/.hermes/.env` (or `gateway.api_server` in
   `~/.hermes/config.yaml`):

   ```bash
   API_SERVER_ENABLED=true
   API_SERVER_KEY=change-me-local-dev   # Bearer token; required
   # API_SERVER_PORT=8642               # default
   # API_SERVER_HOST=127.0.0.1          # default
   ```

3. Pin reference governance (do **not** use `--yolo` / `HERMES_YOLO_MODE`
   for scored runs). Approvals must remain genuine HITL — e.g.
   `approvals.mode` of `smart` or `manual`, not `off`. Prefer
   `approvals.unattended_mode` semantics that leave unanswered gates to
   deny-after-timeout so harness `withhold` is observable.

4. Start the gateway:

   ```bash
   hermes gateway
   # expect: [API Server] API server listening on http://127.0.0.1:8642
   ```

5. Smoke the control plane:

   ```bash
   curl -sS -H "Authorization: Bearer $API_SERVER_KEY" \
     http://127.0.0.1:8642/v1/capabilities
   ```

6. Run the AgentGavel sidecar (stdio JSON-RPC). Default client is
   in-memory `StubHermesClient` (no live Hermes). For a live API server,
   inject `HttpHermesClient` with a requester targeting
   `http://127.0.0.1:8642` and `API_SERVER_KEY` auth
   (`DEFAULT_HERMES_API_BASE` in `adapters.hermes.client`).

   ```bash
   cd adapters/hermes
   PYTHONPATH=src:../../sdk/python/src python3 -m adapters.hermes
   # harness:
   AgentGavel run --adapter "python3 -m adapters.hermes" ...
   ```

Stub path does not require a live Hermes process (unit tests / Handshake).

## Fingerprint fields (reference config)

Run fingerprints use RFC §4.11 / engine keys. Canonical material for
`config.hash` is the SHA-256 of the ordered reference inputs below
(normalize YAML/JSON with stable key order; exclude secrets such as
`API_SERVER_KEY` and `~/.hermes/.env` credentials — hash enablement flags
and non-secret pins only).

| Fingerprint key | Hermes source for the reference run |
| --- | --- |
| `framework.version` | Hermes package / CLI (`hermes --version`) or `GET /v1/capabilities` |
| `config.hash` | SHA-256 of reference config material (table below) |
| `adapter.version` | `adapters/hermes` Handshake `adapter_version` (currently `0.0.1`) |
| `model` | Model id / provider binding (`SessionConfig.model_name` / Oracle) |
| `scenario.version` | Suite scenario pin (engine) |
| `seed.set` | Deterministic probe seeds (engine) |
| `provenance` | Always `unofficial` until ADR 007 ratification |

### `config.hash` material

| Input | Reference expectation |
| --- | --- |
| Control plane | API server enabled (`API_SERVER_ENABLED` / `gateway.api_server.enabled`) |
| Terminal backend | **local** |
| `approvals.mode` | `smart` or `manual` (not `off`) |
| `approvals.timeout` | Pin the value used for the run |
| `approvals.cron_mode` / `approvals.single_query_mode` / `approvals.unattended_mode` / `approvals.deny` | Include when set |
| Model provider / base URL | Oracle binding for oracle runs |
| MCP server pins | Those used with engine fixtures |
| Skill / memory write-approval gates | `skills.write_approval`, `memory.write_approval` if exercised |
| YOLO | Must be **off** for scored runs |

Default layout: `~/.hermes/config.yaml`, `~/.hermes/.env`, session DB under
`~/.hermes/`.

## ResolveApproval

| AgentGavel Decision | Hermes `POST /v1/runs/{run_id}/approval` |
| --- | --- |
| `approve` | `{"choice": "once", "request_id": <approval_id>}` |
| `deny` | `{"choice": "deny", "request_id": <approval_id>}` |
| `withhold` | **No POST** — leave pending (Hermes unanswered → deny after timeout) |

`approval_id` maps to Hermes `request_id`. Bind the active Hermes `run_id`
on the session (`StubHermesClient` creates one on `submit_task`;
`HttpHermesClient` accepts `metadata.hermes_run_id` until full
`POST /v1/runs` is wired).

## Capability honesty

| Flag | Value | Why |
| --- | --- | --- |
| `hitl` | `true` | ResolveApproval → API-server `/v1/runs/{id}/approval` |
| `ledger` | `false` | Trajectories ≠ hash-linked wire Ledger (`ExportLedger` returns empty `entries`) |
| `observability` | `true` | API/SSE frames map to `tool_invocation` + `gate_decision` via `ingest_hermes_event` |
| `tenancy` | `false` | Unproven for reference config |
| `context_mode` | `none` | No prompt attestation yet |
| `framework_name` | `hermes` | |
| `framework_version` | `unknown` | Set after live `hermes --version` / capabilities probe |
| `provenance` | `unofficial` | ADR 007 / ADR 014 |

Missing capabilities score **N/A** (never silent Fail). SEC-002 is scored
when `hitl=true`. Live `GET /v1/runs/{id}/events` subscribe may still be
partial; unit tests feed Hermes-shaped frames into `ingest_hermes_event`.

## Test

```bash
cd adapters/hermes
PYTHONPATH=src:../../sdk/python/src python3 -m pytest tests/ -q
```
