# OpenClaw adapter (unofficial)

Unofficial AgentGavel sidecar targeting OpenClaw's Gateway policy plane
(ADR 014). Provenance is always `unofficial` until ADR 007 ratification.

This is a Handshake scaffold (T15.17). ResolveApproval, Events, and live
Gateway probes land in later tasks. Capability flags stay conservative:
`hitl=false`, `ledger=false`, `observability=false` until proven.

Capability map: [`docs/manual/openclaw-capability-map.md`](../../docs/manual/openclaw-capability-map.md).

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
