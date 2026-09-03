# Microsoft Agent Framework adapter (unofficial)

Unofficial AgentGavel sidecar targeting **Microsoft Agent Framework**
(the AutoGen successor). Directory name is `agent_framework`; Handshake
reports `framework_name=microsoft-agent-framework`.

## Unofficial status

This adapter ships with **`provenance=unofficial`** on every scorecard.
It is community/author-built, not a Microsoft Agent Framework maintainer
contribution. Do not treat unofficial results as ratified framework rankings.

Microsoft Agent Framework replaces AutoGen as the RFC §8.1 target
(never `autogen` as an adapter dir). This scaffold is still unofficial
until ADR 007 ratification.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default until Microsoft Agent Framework maintainers review/contribute, or an independent external reviewer signs off under ADR 007. |
| **provisional** | Possible after independent external review (contract honesty, no oracle special-casing, event completeness) when maintainers decline. Expires after 180 days unless renewed or upgraded. |
| **ratified** | Preferred path: framework maintainers review or contribute the adapter. |

Until that sign-off lands, Handshake and scorecards keep
`provenance: unofficial`.

## Dependency choice

This package does **not** depend on the Microsoft Agent Framework / AutoGen
PyPI packages. Those stacks are heavy for CI/fixtures. Instead,
`adapters.agent_framework.graph.MinimalEmailGraph` is an in-process stub
with `read_email` / `send_email` tool nodes that:

1. Points the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Executes the matching tool node.
3. Records `tool_invocation` before/after events and `context_attestation`
   of the prompt (ADR 005) via callback / `Adapter.emit` when a transport
   is attached.

CapabilityReport (honest): `observability=true`, `context_mode=attestation`,
`hitl=false`, `ledger=false`, `tenancy=false`.

Swap in the real framework later if needed; the observation contract stays
the same.

## Capability honesty

Handshake reports only what this sidecar actually implements today:

- `hitl: false` — no ResolveApproval / interrupt mapping yet
- `ledger: false` — no session hash-linked ledger
- `observability: true` — `tool_invocation` + attestation event sink (T13.19)
- `tenancy: false`
- `context_mode: attestation`

Do not treat stub flags as Microsoft Agent Framework product limitations.
Flip flags only when real support lands in the same change.

## Run

```bash
# From adapters/agent_framework after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.agent_framework --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.agent_framework   # stdio serve
```

## Test

```bash
cd adapters/agent_framework
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
