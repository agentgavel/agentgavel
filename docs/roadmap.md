# AgentGavel Roadmap

Living progress story for AgentGavel (RFC 0001). Detailed tasks live in
`docs/plan.md` and `docs/plans/`.

## Now (frontier)

- **In flight (2026-09-03):** Wave 6 batch 2 — T3.6 T6.8 engine; T8.4 T8.6 T8.7 T8.8 scenarios.

- **Shipped (2026-09-03):** Wave 6 batch 1 — T3.5 oracle SessionConfig, T7.2 Python
  stdio transport, T8.3 SEC-001, T9.3 CLI report.

- **Shipped (2026-09-03):** Wave 5 — T3.2 adapter process launcher, T7.1 Python SDK
  scaffold, T8.1 SEC-v1 suite loader, T9.1 run fingerprint hasher.

- **Shipped (2026-09-03):** Wave 4 — session lifecycle + version negotiate, E2 complete;
  attestation leaks; GSI/hard-soft/Wilson metrics (E5 complete); engine Scenario/SeedScheduler/artifacts.

- **v0.1 harness**: Go engine, JSON-RPC stdio adapter protocol, Python SDK,
  Compliance Oracle, assertions + GSI, mcpfuzz, SEC-001..007, CLI run/report,
  unofficial Sire and LangGraph adapters, CI + tagged release.
  Plan epics: E1-E12.
- **Shipped (2026-09-03):** Wave 3 — protocol types/codec/stdio, CapabilityReport
  N/A mapping, `AgentGavel oracle --listen`, assertions (tools/gate/leak/partial).
- **Shipped (2026-09-03):** Wave 2 — Makefile/CI, adapter.proto, OpenAI+Anthropic
  Oracle handlers, mcpfuzz modes toxic/schema/early-disconnect/renamer/slowloris/masquerade.
- **Shipped (2026-09-02):** Wave 1 complete (T1.1, T1.2, T1.5, T8.2, T4.1, T6.1).
- **Next:** After batch 1: T3.6 T6.8 engine; T8.4/6/7/8 scenarios.

## Planned

- **v0.2**: SEC-008..010, governance suite, AutoGen/CrewAI adapters, `--ci`
  mode, scenario governance process. Epic E13 (outline until v0.1 ships).
- **v0.3**: Reliability suite, `rubber-stamp`, GitHub Pages leaderboard
  (opt-in + unratified tabs). Epic E14.
- **v1.0**: Public signed submissions, harness red-team bounty, ratification
  operations. Epic E15.

## Done

- (none yet)

## References

- Spec: `docs/RFC-0001.md`
- Design: `docs/design.md`
- Plan: `docs/plan.md`
- ADRs: `docs/adr/`
