# Sire adapter (unofficial)

Scaffold for the unofficial AgentGavel sidecar targeting Sire.

**Provenance:** `unofficial` — author-affiliated adapters cannot self-ratify
(see `docs/adr/007-adapter-ratification.md`). Full status documentation lands
in T10.6.

## Run

```bash
# From adapters/sire after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.sire --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.sire   # stdio serve
```

Lifecycle mapping to real Sire APIs is T10.2+.
