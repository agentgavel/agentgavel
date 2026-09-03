# LangGraph adapter (unofficial)

Unofficial AgentGavel sidecar targeting LangGraph-shaped agent graphs.

**Provenance:** `unofficial` — community/author-built until a maintainer
signs off (see `docs/adr/007-adapter-ratification.md`). Full status
documentation lands in T11.6.

## Dependency choice (T11.2)

This package does **not** depend on the `langgraph` PyPI package. That
stack pulls LangChain and is heavy for CI/fixtures. Instead,
`adapters.langgraph.graph.MinimalEmailGraph` is an in-process stub with
`read_email` / `send_email` tool nodes that:

1. Points the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Executes the matching tool node.
3. Records `tool_invocation` before/after events (via callback /
   `Adapter.emit` when a transport is attached).

Swap in real LangGraph later if needed; the observation contract stays
the same.

## Run

```bash
# From adapters/langgraph after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.langgraph --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.langgraph   # stdio serve
```

## Test

```bash
cd adapters/langgraph
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
