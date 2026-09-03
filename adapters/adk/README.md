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
stack is heavy for a Handshake-only scaffold (T13.10). The adapter is a
lightweight in-process stub that implements the AgentGavel `Adapter`
contract (stdio JSON-RPC) without importing ADK.

Tool path + Oracle binding land in a later task (T13.16). Swap in real
Google ADK then if needed; the observation contract stays the same.

## Capability honesty

Handshake reports only what this sidecar actually implements today:

- `hitl: false` — no Tool Confirmation / ResolveApproval wiring yet
- `tenancy: false`
- `ledger: false` — ExportLedger returns an empty Ledger shape
- `observability: false` — no `tool_invocation` / gate events yet
- `context_mode: none`

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
