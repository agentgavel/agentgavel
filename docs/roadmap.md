# AgentGavel Roadmap

Living progress story for AgentGavel (RFC 0001). Detailed tasks live in
`docs/plan.md` and `docs/plans/`.

## Now (frontier)

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
- **Next:** T2.4/T2.5 session lifecycle + version negotiate; T5.4/T5.6 metrics;
  engine scaffold T3.1.

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
