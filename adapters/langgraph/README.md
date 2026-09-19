# LangGraph adapter (unofficial — provisional pending)

AgentGavel sidecar targeting LangGraph-shaped agent graphs.

## Status (T15.11 wait)

Handshake and scorecards currently report **`provenance=unofficial`**.

The ADR 007 provisional package is prepared and waiting on the public
comment window (restarted after real maintainer outreach):

| Field | Value |
| --- | --- |
| Outreach (LangGraph) | https://github.com/langchain-ai/langgraph/issues/8992 |
| Comment / objections | https://github.com/agentgavel/agentgavel/issues/175 |
| Window opened | 2026-09-18 |
| Window closes | **2026-10-18** |
| Checklist / ops record | [`docs/manual/ratification/langgraph-provisional-2026-09-06.md`](../../docs/manual/ratification/langgraph-provisional-2026-09-06.md) |

After 2026-10-18, if no blocking objection stands on #175 / #8992, core
maintainers may flip Handshake (and a dashboard sample) to
`provenance=provisional` with grant date / 180-day expiry recorded in
that file. **Provisional is not ratified.**

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).
Ops checklist: [`docs/manual/adapter-ratification.md`](../../docs/manual/adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default; window open through 2026-10-18 |
| **provisional** | Eligible after window closes + checklist already passed in the record |
| **ratified** | Preferred: LangGraph maintainers review or contribute |

Author-affiliated **Sire** remains `unofficial` (cannot skip to ratified via
the provisional path).

## Dependency choice (stub vs live)

**Default (`runtime=stub`):** this package does **not** require the
`langgraph` PyPI package. CI uses
`adapters.langgraph.graph.MinimalEmailGraph`, an in-process stub with
`read_email` / `send_email` tool nodes. Handshake reports
`runtime=stub` and `framework_version=stub-0.0.1`. **Stub scores are not
product rankings** ([ADR 015](../../docs/adr/015-stub-vs-live-runtime.md)).

**Optional live (`runtime=live`):** install the extra and set the env var
so the sidecar drives a real LangGraph `StateGraph` with
`interrupt()` / `Command(resume=...)`:

```bash
cd adapters/langgraph
pip install -e '.[live]'
export AGENTGAVEL_LANGGRAPH_RUNTIME=live
# Handshake: runtime=live, framework_version=<installed langgraph>
PYTHONPATH=src:../../sdk/python/src python -m adapters.langgraph
```

Without the package, `AGENTGAVEL_LANGGRAPH_RUNTIME=live` fails closed
(never silently reports live while running the stub).

Both modes:

1. Point the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Execute the matching tool node (after HITL when gated).
3. Record `tool_invocation` before/after, `context_attestation`
   of the prompt (ADR 005), and `gate_decision` on ResolveApproval.

CapabilityReport (honest): `observability=true`, `context_mode=attestation`,
`hitl` tracks interrupt support, `ledger=false`, `provenance=unofficial`.

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
