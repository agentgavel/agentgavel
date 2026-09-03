# OpenAI Agents SDK adapter (unofficial)

Unofficial AgentGavel sidecar targeting the [OpenAI Agents SDK](https://github.com/openai/openai-agents-python).

## Unofficial status

This adapter ships with **`provenance=unofficial`** on every scorecard.
It is community/author-built, not an OpenAI-maintainer contribution.

A low score behind this adapter is a claim about the adapter as much as
about the OpenAI Agents SDK. Do not treat unofficial results as ratified
framework rankings.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default until OpenAI Agents SDK maintainers review/contribute, or an independent external reviewer signs off under ADR 007. |
| **provisional** | Possible after independent external review (contract honesty, no oracle special-casing, event completeness) when maintainers decline. Expires after 180 days unless renewed or upgraded. |
| **ratified** | Preferred path: OpenAI Agents SDK maintainers review or contribute the adapter. |

Until that sign-off lands, Handshake and scorecards keep
`provenance: unofficial`.

## Dependency choice

This package does **not** depend on the `openai-agents` PyPI package.
That SDK pulls a heavier tool/LLM stack than CI and fixtures need for a
Handshake scaffold (T13.11). Lifecycle methods are in-process stubs so
`python -m adapters.openai_agents` starts stdio without installing the
real framework.

A later wave (T13.17) adds a minimal Oracle tool path and may optionally
pull `openai-agents` only when wiring real agents. Document any dependency
change in the same PR that flips CapabilityReport flags.

## Capability honesty

Scaffold CapabilityReport (honest empty support):

- `hitl: false` — `needs_approval` / interrupt → ResolveApproval not wired yet
- `tenancy: false`
- `ledger: false` — `ExportLedger` returns empty `entries`
- `observability: false` — no `tool_invocation` event sink yet
- `context_mode: none`

Do not treat stub flags as OpenAI Agents SDK product limitations.

## Run

```bash
# From adapters/openai_agents after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.openai_agents --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.openai_agents   # stdio serve
```

## Test

```bash
cd adapters/openai_agents
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
