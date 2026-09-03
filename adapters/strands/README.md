# AWS Strands Agents adapter (unofficial)

Unofficial AgentGavel sidecar targeting AWS Strands Agents.

## Unofficial status

This adapter ships with **`provenance=unofficial`** on every scorecard.
It is community/author-built, not an AWS Strands Agents-maintainer contribution.

A low score behind this adapter is a claim about the adapter as much as
about AWS Strands Agents. Do not treat unofficial results as ratified framework
rankings.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default until AWS Strands Agents maintainers review/contribute, or an independent external reviewer signs off under ADR 007. |
| **provisional** | Possible after independent external review (contract honesty, no oracle special-casing, event completeness) when maintainers decline. Expires after 180 days unless renewed or upgraded. |
| **ratified** | Preferred path: AWS Strands Agents maintainers review or contribute the adapter. |

Until that sign-off lands, Handshake and scorecards keep
`provenance: unofficial`.

## Dependency choice

This package does **not** depend on the `strands-agents` PyPI package. That
stack is heavy for CI/fixtures. Instead,
`adapters.strands.graph.MinimalEmailGraph` is an in-process stub with
`read_email` / `send_email` tool nodes that:

1. Points the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Executes the matching tool node.
3. Records `tool_invocation` before/after events and `context_attestation`
   of the prompt (ADR 005) via callback / `Adapter.emit` when a transport
   is attached.

CapabilityReport (honest): `observability=true`, `context_mode=attestation`,
`hitl=false` (no interrupt yet), `ledger=false`, `tenancy=false`.

Swap in real AWS Strands Agents later if needed; the observation contract
stays the same.

## Capability honesty

Handshake reports only what this sidecar actually implements today:

- `hitl: false` — no interrupt / ResolveApproval wiring yet
- `tenancy: false`
- `ledger: false` — ExportLedger returns an empty Ledger shape
- `observability: true` — `tool_invocation` + `context_attestation` events
- `context_mode: attestation`

Do not treat stub flags as AWS Strands Agents product limitations.

## Run

```bash
# From adapters/strands after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.strands --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.strands   # stdio serve
```

## Test

```bash
cd adapters/strands
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
