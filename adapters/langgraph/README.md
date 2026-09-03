# LangGraph adapter (unofficial)

Scaffold for the unofficial AgentGavel sidecar targeting LangGraph.

**Provenance:** `unofficial` — community/author-built until a maintainer
signs off (see `docs/adr/007-adapter-ratification.md`). Full status
documentation lands in T11.6.

## Run

```bash
# From adapters/langgraph after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.langgraph --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.langgraph   # stdio serve
```

Graph tools and lifecycle mapping are T11.2+.
