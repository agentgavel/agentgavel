# LangGraph adapter (unofficial)

Unofficial AgentGavel sidecar targeting LangGraph-shaped agent graphs.

## Unofficial status

This adapter ships with **`provenance=unofficial`** on every scorecard.
It is community/author-built, not a LangGraph-maintainer contribution.

A low score behind this adapter is a claim about the adapter as much as
about LangGraph. Do not treat unofficial results as ratified framework
rankings.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default until LangGraph maintainers review/contribute, or an independent external reviewer signs off under ADR 007. |
| **provisional** | Possible after independent external review (contract honesty, no oracle special-casing, event completeness) when maintainers decline. Expires after 180 days unless renewed or upgraded. |
| **ratified** | Preferred path: LangGraph maintainers review or contribute the adapter. |

Until that sign-off lands, Handshake and scorecards keep
`provenance: unofficial`.

## Dependency choice

This package does **not** depend on the `langgraph` PyPI package. That
stack pulls LangChain and is heavy for CI/fixtures. Instead,
`adapters.langgraph.graph.MinimalEmailGraph` is an in-process stub with
`read_email` / `send_email` tool nodes that:

1. Points the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Executes the matching tool node.
3. Records `tool_invocation` before/after events, `context_attestation`
   of the prompt (ADR 005), and `gate_decision` on ResolveApproval
   (via callback / `Adapter.emit` when a transport is attached).

CapabilityReport (honest): `observability=true`, `context_mode=attestation`,
`hitl=false` until T11.3 maps interrupts, `ledger=false`.

Swap in real LangGraph later if needed; the observation contract stays
the same.

## HITL / interrupts

When interrupt support is enabled (default), gated tools (`send_email`)
pause before side effects — LangGraph-style `interrupt()` — and wait for
harness `ResolveApproval`. Handshake reports `hitl: true` and a
`gate_decision` event is emitted on resolve (`source=harness`,
`genuine_hitl=true`).

Construct with `LangGraphAdapter(hitl=False)` for the honest unsupported
path: `hitl: false` and `ResolveApproval` raises `HitlNotSupportedError`.

## Capability honesty

Handshake reports only what this sidecar actually implements. Ledger
export stays `false` until a real session ledger exists. Observability and
`context_mode=attestation` reflect the tool/gate/attestation event sink.
Do not treat stub flags as LangGraph product limitations.

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
