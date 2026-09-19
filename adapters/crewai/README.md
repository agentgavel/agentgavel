# CrewAI adapter (unofficial)

Unofficial AgentGavel sidecar targeting CrewAI-shaped multi-agent crews.

## Unofficial status

This adapter ships with **`provenance=unofficial`** on every scorecard.
It is community/author-built, not a CrewAI-maintainer contribution.

A low score behind this adapter is a claim about the adapter as much as
about CrewAI. Do not treat unofficial results as ratified framework
rankings.

### Ratification path (ADR 007)

Full policy: [`docs/adr/007-adapter-ratification.md`](../../docs/adr/007-adapter-ratification.md).

| Label | Meaning for this adapter |
| --- | --- |
| **unofficial** (current) | Default until CrewAI maintainers review/contribute, or an independent external reviewer signs off under ADR 007. |
| **provisional** | Possible after independent external review (contract honesty, no oracle special-casing, event completeness) when maintainers decline. Expires after 180 days unless renewed or upgraded. |
| **ratified** | Preferred path: CrewAI maintainers review or contribute the adapter. |

Until that sign-off lands, Handshake and scorecards keep
`provenance: unofficial`.

## Dependency choice (stub vs live)

**Default (`runtime=stub`):** this package does **not** require the `crewai`
PyPI package. CI uses `MinimalEmailCrew`. Handshake reports
`runtime=stub` / `framework_version=stub-0.0.1`. Stub scores are **not**
product rankings ([ADR 015](../../docs/adr/015-stub-vs-live-runtime.md)).

**Optional live (`runtime=live`):** install the extra and set the env var so
Handshake reports the real package version (Oracle tool path still uses
`MinimalEmailCrew` until a full CrewAI graph lands):

```bash
cd adapters/crewai
pip install -e '.[live]'
export AGENTGAVEL_CREWAI_RUNTIME=live
PYTHONPATH=src:../../sdk/python/src python -m adapters.crewai
```

Without the package, `AGENTGAVEL_CREWAI_RUNTIME=live` fails closed.

Both modes use Oracle-compatible `read_email` / `send_email` tool nodes that:

1. Point the model client at a Compliance Oracle `base_url`
   (`POST …/v1/chat/completions` with `X-AgentGavel-Probe-Directive`).
2. Execute the matching tool node.
3. Record `tool_invocation` before/after events and `context_attestation`
   of the prompt (ADR 005) via callback / `Adapter.emit` when a transport
   is attached.

CapabilityReport (honest): `observability=true`, `context_mode=attestation`,
`hitl=false`, `ledger=false`, `tenancy=false`.

Swap in real CrewAI later if needed; the observation contract stays the same.

## Capability honesty

Handshake reports only what this sidecar actually implements today:

- `hitl: false` — no ResolveApproval / `human_input` mapping yet
- `tenancy: false`
- `ledger: false` — ExportLedger returns an empty Ledger shape
- `observability: true` — `tool_invocation` + attestation event sink (T13.21)
- `context_mode: attestation`

Missing capabilities score N/A (never silent Fail). Do not treat stub
flags as CrewAI product limitations.

## Run

```bash
# From adapters/crewai after editable install, or with PYTHONPATH:
PYTHONPATH=src:../../sdk/python/src python -m adapters.crewai --help
PYTHONPATH=src:../../sdk/python/src python -m adapters.crewai   # stdio serve
```

## Test

```bash
cd adapters/crewai
PYTHONPATH=src:../../sdk/python/src python -m pytest tests/ -q
```
