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

## Dependency choice

This package does **not** depend on the `pydantic-ai` PyPI package. That
stack pulls LLM providers and tooling extras unsuitable for CI/scaffold
smoke. The sidecar is an in-process stub that speaks the AgentGavel
stdio JSON-RPC contract.

T13.12 is Handshake + lifecycle stubs only. A minimal Oracle tool path
(and any real Pydantic AI wiring) lands in a later task (T13.18). Swap in
real Pydantic AI later if needed; the observation contract stays the same.

## Capability honesty

Handshake reports only what this sidecar actually implements:

- `hitl: false` — no ResolveApproval / deferred-tools mapping yet
- `tenancy: false`
- `ledger: false` — ExportLedger returns an empty Ledger shape
- `observability: false` — no `tool_invocation` / event sink yet
- `context_mode: none`

Missing capabilities score N/A (never silent Fail). Do not treat stub
flags as Pydantic AI product limitations.

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
