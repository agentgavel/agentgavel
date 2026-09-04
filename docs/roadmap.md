# AgentGavel Roadmap

Living progress story for AgentGavel (RFC 0001). Detailed tasks live in
`docs/plan.md` and `docs/plans/`.

## Now (frontier)

- **Shipped (2026-09-04):** T12.8 cut `v0.1.0` tag + release. 5 binary
  assets uploaded (darwin/linux, amd64/arm64) + checksums, verified via
  `gh release view`. **v0.1.0 released:**
  https://github.com/agentgavel/agentgavel/releases/tag/v0.1.0.
  **E12 complete.** Wave 13 (T13.0, next-horizon planning) is now
  unblocked.

- **Shipped (2026-09-03):** T12.6 RFC Status → Implemented (v0.1 scope)
  (this PR). Brand mark in `brand/` (#81). Module rename (#80).

- **Shipped (2026-09-03):** T12.7 full make test/lint (#76).

- **Shipped (2026-09-03):** T12.4 v0.1 smoke doc (#75).

- **Shipped (2026-09-03):** T11.7 LangGraph ruff (#74). **E11 complete.**

- **Shipped (2026-09-03):** T11.5 LangGraph SEC-001/007 E2E (#73).

- **Shipped (2026-09-03):** T10.7 Sire ruff (#67); T11.2 LangGraph tools (#66); T11.3 HITL (#69); T11.4 events (#70); T11.6 LangGraph README (#71). **E10 complete.**

- **Shipped (2026-09-03):** T10.6 Sire unofficial README/ratification (#64).

- **Shipped (2026-09-03):** T10.5 Sire oracle E2E SEC-002 (#62).

- **Shipped (2026-09-03):** T12.1 CI matrix (#60); T12.3 VERSION/ldflags (#61).

- **Shipped (2026-09-03):** T12.2 GoReleaser static binaries (#56).

- **Shipped (2026-09-03):** T9.6 cmd gofmt (#53); T10.4 ExportLedger honesty (#54). **E9 complete.**

- **Shipped (2026-09-03):** T9.4 run --fingerprint reload (#51).

- **Shipped (2026-09-03):** T9.5 CLI binary tests (#48).

- **Shipped (2026-09-03):** T10.3 Sire ResolveApproval + gate_decision (#46).

- **Shipped (2026-09-03):** T9.2 AgentGavel run oracle security (#44); T11.1 LangGraph scaffold (#43).

- **Shipped (2026-09-03):** T10.2 Sire mockable lifecycle client (#42).

- **Shipped (2026-09-03):** T8.11 suites/security gofmt already clean (#39).

- **Shipped (2026-09-03):** T10.1 unofficial Sire adapter scaffold (#36).

- **Shipped (2026-09-03):** T7.6 ruff on sdk/python (#33); T8.10 SuiteOracleFake E2E (#34).

- **Shipped (2026-09-03):** Wave close — T7.5 (#30), T8.9 (#27), T12.5 (#26).

- **Shipped (2026-09-03):** T8.9 SEC-007 composite fuzz (#27).

- **Shipped (2026-09-03):** T12.5 v0.1 plan-execution stub in docs/devlog.md (#26).

- **Shipped (2026-09-03):** Wave 7 batch — T7.4 FakeAdapter example, T8.5 SEC-003, T3.7 gofmt (already clean).

- **Next:** After this wave: T8.11 lint; T9.2 CLI run; T7.6 ruff if still open.

- **Shipped (2026-09-03):** Wave 6 complete — T6.8 StartFuzzMode; E6 closed.

- **Shipped (2026-09-03):** Wave 6 batch 2 — T3.6 IntegrationNoop, T7.3 callbacks, T8.4/6/7/8.

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
- **Next:** T7.4 FakeAdapter example; T8.5/T8.9 (need T6.8); Wave 7 SEC scenarios.

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
