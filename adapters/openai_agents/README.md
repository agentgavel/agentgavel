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
That SDK pulls a heavier tool/LLM stack than CI and fixtures need.
Instead, `adapters.openai_agents.agent.MinimalEmailAgent` is an in-process
stub with `read_email` / `send_email` tools that:

1. Points the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Executes the matching tool.
3. Records `tool_invocation` before/after events and `context_attestation`
   of the prompt (ADR 005) via callback / `Adapter.emit` when a transport
   is attached.

CapabilityReport (honest): `observability=true`, `context_mode=attestation`,
`hitl=false` (no `needs_approval` → ResolveApproval yet), `ledger=false`,
`tenancy=false`.

Swap in real OpenAI Agents SDK later if needed; the observation contract
stays the same.

## Capability honesty

Handshake reports only what this sidecar actually implements. Ledger
export stays `false` until a real session ledger exists. Observability and
`context_mode=attestation` reflect the tool/attestation event sink.
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
