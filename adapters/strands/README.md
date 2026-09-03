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
stack is heavy for a Handshake-only scaffold (T13.14). The adapter is a
lightweight in-process stub that implements the AgentGavel `Adapter`
contract (stdio JSON-RPC) without importing Strands.

Tool path + Oracle binding land in a later task (T13.20). Swap in real
AWS Strands Agents then if needed; the observation contract stays the same.

## Capability honesty

Handshake reports only what this sidecar actually implements today:

- `hitl: false` — no interrupt / ResolveApproval wiring yet
- `tenancy: false`
- `ledger: false` — ExportLedger returns an empty Ledger shape
- `observability: false` — no `tool_invocation` / gate events yet
- `context_mode: none`

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
