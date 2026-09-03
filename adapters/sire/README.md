# Sire adapter (unofficial)

Unofficial AgentGavel sidecar targeting Sire.

**Provenance:** `unofficial` — author-affiliated adapters cannot self-ratify
(see `docs/adr/007-adapter-ratification.md`). Full status documentation lands
in T10.6.

## Lifecycle mapping (best-effort)

Sire's public API (`https://api.sire.run/api/v1`) exposes Workers and Runs, not
AgentGavel sessions. The adapter drives a mockable `SireClient`:

| AgentGavel | Default stub | Documented live HTTP (`HttpSireClient`) |
| --- | --- | --- |
| `start_session` | in-memory bind | `GET /workers/{workerId}` |
| `submit_task` | record prompt | `POST /workers/{workerId}/run` (seed = prompt) |
| `stop_session` | mark stopped | `POST /runs/{runId}/cancel` |

A live client still needs a bearer token and `extra["sire_worker_id"]` (or
`HttpSireClient(worker_id=...)`). Point the worker's model `base_url` at the
Compliance Oracle. `ResolveApproval` is not wired yet (T10.3). HITL / ledger /
observability stay `false` until those surfaces exist.

Default construction uses `StubSireClient` so the sidecar starts without Sire.

## Run

```bash
# From adapters/sire after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.sire --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.sire   # stdio serve

cd adapters/sire && PYTHONPATH=src:../../sdk/python/src pytest
```
