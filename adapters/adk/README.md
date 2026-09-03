# Google ADK adapter (unofficial)

Unofficial AgentGavel sidecar targeting Google's Agent Development Kit (ADK).

## Unofficial status

This adapter ships with **`provenance=unofficial`** on every scorecard.
It is community/author-built, not a Google ADK-maintainer contribution.

A low score behind this adapter is a claim about the adapter as much as
about Google ADK. Do not treat unofficial results as ratified framework
rankings.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default until Google ADK maintainers review/contribute, or an independent external reviewer signs off under ADR 007. |
| **provisional** | Possible after independent external review (contract honesty, no oracle special-casing, event completeness) when maintainers decline. Expires after 180 days unless renewed or upgraded. |
| **ratified** | Preferred path: Google ADK maintainers review or contribute the adapter. |

Until that sign-off lands, Handshake and scorecards keep
`provenance: unofficial`.

## Dependency choice

This package does **not** depend on the `google-adk` PyPI package. That
stack is heavy for CI/fixtures. Instead,
`adapters.adk.graph.MinimalEmailGraph` is an in-process stub with
`read_email` / `send_email` tool nodes that:

1. Points the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Executes the matching tool node.
3. Records `tool_invocation` before/after events and `context_attestation`
   of the prompt (ADR 005) via callback / `Adapter.emit` when a transport
   is attached.

CapabilityReport (honest): `observability=true`, `context_mode=attestation`,
`hitl=false`, `tenancy=false`, `ledger=false`.

Swap in real Google ADK later if needed; the observation contract stays
the same.

## Capability honesty

Handshake reports only what this sidecar actually implements today:

- `hitl: false` — no Tool Confirmation / ResolveApproval wiring yet
- `tenancy: false`
- `ledger: false` — ExportLedger returns an empty Ledger shape
- `observability: true` — `tool_invocation` before/after + attestation emit
- `context_mode: attestation`

Do not treat stub flags as Google ADK product limitations.

## Run

```bash
# From adapters/adk after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.adk --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.adk   # stdio serve
```

## Test

```bash
cd adapters/adk
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
