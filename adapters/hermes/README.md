# Hermes Agent adapter (unofficial)

Unofficial AgentGavel **gateway-style** sidecar (ADR 014) for
[Hermes Agent](https://github.com/nousresearch/hermes-agent) (Nous Research).

ResolveApproval is wired (T15.25). Events + ExportLedger honesty land in
T15.26. Full operator README lands in T15.28.
Capability map:
[`docs/manual/hermes-capability-map.md`](../../docs/manual/hermes-capability-map.md).

## Unofficial status

Handshake and scorecards always report **`provenance=unofficial`** until
ADR 007 ratification. Do not treat unofficial results as ratified rankings.

## Control plane (pinned)

This sidecar targets Hermes' **OpenAI-compatible API server**
(`gateway/platforms/api_server.py`): HTTP + SSE, `GET /v1/capabilities`,
`POST /v1/runs`, `POST /v1/runs/{id}/approval`. TUI gateway and ACP are
documented alternatives; they are not the reference path for this adapter.

## ResolveApproval (T15.25)

| AgentGavel Decision | Hermes `POST /v1/runs/{run_id}/approval` |
| --- | --- |
| `approve` | `{"choice": "once", "request_id": <approval_id>}` |
| `deny` | `{"choice": "deny", "request_id": <approval_id>}` |
| `withhold` | **No POST** — leave pending (Hermes unanswered → deny after timeout) |

`approval_id` maps to Hermes `request_id`. The active Hermes `run_id` must be
bound on the session (stub creates one on `submit_task`; HTTP client accepts
`metadata.hermes_run_id` until full `POST /v1/runs` lands).

Default client is in-memory `StubHermesClient` (no live Hermes). Inject
`HttpHermesClient` with a requester for a real API server.

## Run as AgentGavel sidecar

```bash
cd adapters/hermes
PYTHONPATH=src:../../sdk/python/src python3 -m adapters.hermes
# or:
AgentGavel run --adapter "python3 -m adapters.hermes" ...
```

Stdio JSON-RPC via `agentgavel_adapter`. Stub path does not require a live
Hermes process.

## Capability honesty (T15.25 / T15.26)

| Flag | Value | Why |
| --- | --- | --- |
| `hitl` | `true` | ResolveApproval → API-server `/v1/runs/{id}/approval` |
| `ledger` | `false` | Trajectories ≠ hash-linked wire Ledger (`ExportLedger` returns empty `entries`) |
| `observability` | `true` | API/SSE frames map to `tool_invocation` + `gate_decision` via `ingest_hermes_event` |
| `tenancy` | `false` | Unproven for reference config |
| `context_mode` | `none` | No prompt attestation yet |
| `framework_name` | `hermes` | |
| `framework_version` | `unknown` | Set after live `hermes --version` / capabilities probe |

Missing capabilities score **N/A** (never silent Fail). SEC-002 is scored when
`hitl=true`. Live `GET /v1/runs/{id}/events` subscribe is still deferred; unit
tests feed Hermes-shaped frames into `ingest_hermes_event`. Full operator
README: T15.28.
