# Pydantic AI adapter (unofficial)

Unofficial AgentGavel sidecar targeting Pydantic AI agents.

## Unofficial status

This adapter ships with **`provenance=unofficial`** on every scorecard.
It is community/author-built, not a Pydantic AI-maintainer contribution.

A low score behind this adapter is a claim about the adapter as much as
about Pydantic AI. Do not treat unofficial results as ratified framework
rankings.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default until Pydantic AI maintainers review/contribute, or an independent external reviewer signs off under ADR 007. |
| **provisional** | Possible after independent external review (contract honesty, no oracle special-casing, event completeness) when maintainers decline. Expires after 180 days unless renewed or upgraded. |
| **ratified** | Preferred path: Pydantic AI maintainers review or contribute the adapter. |

Until that sign-off lands, Handshake and scorecards keep
`provenance: unofficial`.

## Dependency choice (T13.18)

This package does **not** depend on the `pydantic-ai` PyPI package. That
stack pulls LLM providers and tooling extras unsuitable for CI/fixtures.
Instead, `adapters.pydantic_ai.agent.MinimalEmailAgent` is an in-process
stub with `read_email` / `send_email` tools that:

1. Points the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Executes the matching tool.
3. Records `tool_invocation` before/after events (via callback /
   `Adapter.emit` when a transport is attached).

CapabilityReport (honest): `hitl=false`, `tenancy=false`, `ledger=false`,
`observability=false`, `context_mode=none` until later tasks wire real
deferred-tools / ledger / event hooks. Missing capabilities score N/A
(never silent Fail).

Swap in real Pydantic AI later if needed; the observation contract stays
the same.

## Capability honesty

Handshake reports only what this sidecar actually implements:

- `hitl: false` — no ResolveApproval / deferred-tools mapping yet
- `tenancy: false`
- `ledger: false` — ExportLedger returns an empty Ledger shape
- `observability: false` — local capture for tests; harness event-hook
  claim lands later
- `context_mode: none`

Do not treat stub flags as Pydantic AI product limitations.

## Run

```bash
# From adapters/pydantic_ai after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.pydantic_ai --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.pydantic_ai   # stdio serve
```

## Test

```bash
cd adapters/pydantic_ai
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
