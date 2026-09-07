# LangGraph adapter (provisional)

AgentGavel sidecar targeting LangGraph-shaped agent graphs.

## Provisional status (T15.11)

This adapter ships with **`provenance=provisional`** on Handshake and
scorecards after ADR 007 independent review.

| Field | Value |
| --- | --- |
| Grant date | 2026-09-06 |
| Expires | 2027-03-05 (180 days) |
| Checklist / outreach record | [`docs/manual/ratification/langgraph-provisional-2026-09-06.md`](../../docs/manual/ratification/langgraph-provisional-2026-09-06.md) |
| Public comment / outreach | https://github.com/agentgavel/agentgavel/issues/175 |

A low score behind this adapter is still a claim about the adapter as much
as about LangGraph. **Provisional is not ratified** — prefer maintainer
ratification when LangGraph engage on #175.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).
Ops checklist: [`docs/manual/adapter-ratification.md`](../../docs/manual/adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** | Pre-T15.11 default |
| **provisional** (current) | Core grant after checklist; expires 2027-03-05 unless renewed or upgraded |
| **ratified** | Preferred: LangGraph maintainers review or contribute the adapter |

Author-affiliated **Sire** remains `unofficial` (cannot skip to ratified via
this path).

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
`hitl` tracks interrupt support, `ledger=false`.

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
