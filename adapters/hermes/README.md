# Hermes Agent adapter (unofficial)

Unofficial AgentGavel **gateway-style** sidecar (ADR 014) for
[Hermes Agent](https://github.com/nousresearch/hermes-agent) (Nous Research).

Scaffold only (T15.24). Full operator README lands in T15.28. Capability map:
[`docs/manual/hermes-capability-map.md`](../../docs/manual/hermes-capability-map.md).

## Unofficial status

Handshake and scorecards always report **`provenance=unofficial`** until
ADR 007 ratification. Do not treat unofficial results as ratified rankings.

## Control plane (pinned)

This sidecar targets Hermes' **OpenAI-compatible API server**
(`gateway/platforms/api_server.py`): HTTP + SSE, `GET /v1/capabilities`,
`POST /v1/runs`, `POST /v1/runs/{id}/approval`. TUI gateway and ACP are
documented alternatives; they are not the reference path for this adapter.

## Run as AgentGavel sidecar

```bash
cd adapters/hermes
PYTHONPATH=src:../../sdk/python/src python3 -m adapters.hermes
# or:
AgentGavel run --adapter "python3 -m adapters.hermes" ...
```

Stdio JSON-RPC via `agentgavel_adapter`. No live Hermes process is required
for Handshake scaffold tests.

## Capability honesty (T15.24)

| Flag | Value | Why |
| --- | --- | --- |
| `hitl` | `false` | ResolveApproval deferred to T15.25 |
| `ledger` | `false` | Trajectories ≠ hash-linked wire Ledger |
| `observability` | `false` | Events mapping deferred to T15.26 |
| `tenancy` | `false` | Unproven for reference config |
| `context_mode` | `none` | No prompt attestation yet |
| `framework_name` | `hermes` | |
| `framework_version` | `unknown` | Set after live `hermes --version` / capabilities probe |
