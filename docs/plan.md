# AgentGavel -- implement RFC 0001 in full

## 1. Context

Problem: Agent frameworks publish little comparable data about governance and
security under adversarial pressure. Task-completion benchmarks do not answer
whether policy ceilings, HITL gates, and audit provenance hold when the model
is fully compromised.

Objectives:
- Implement `docs/RFC-0001.md` end to end across releases v0.1, v0.2, v0.3, v1.0.
- Ship a trustworthy Go engine, adapter wire protocol, Python SDK, Compliance
  Oracle, security suite SEC-001..007, and unofficial Sire + LangGraph adapters
  as the v0.1 frontier.
- Keep later releases as outline epics expanded after each tag (rolling wave).

Non goals (v0.1):
- Governance and reliability suites (v0.2 / v0.3).
- AutoGen / CrewAI adapters (v0.2).
- Public leaderboard and rubber-stamp CLI (v0.3).
- Public submission bounty (v1.0).
- Per-framework exploit code (forbidden by RFC section 0).

Constraints and assumptions:
- Greenfield repo `github.com/agentgavel/gavel` (no application code yet).
- Neutrality rules in RFC section 0 are binding.
- Open questions Q3-Q7 resolved in ADRs 002, 004, 005, 006, 007 (recommended
  defaults accepted for planning; founder may override before v0.1 tag).
- No hosted production service until E14; v0.1 "done" means tagged GitHub
  release + CI green + documented local smoke.

Success metrics:
- `AgentGavel run` produces fingerprint + GSI scorecard for FakeAdapter.
- SEC-001..007 predicates covered by automated tests.
- Soft rates use >=25 seeds with Wilson intervals.
- Sire and LangGraph adapters labeled unofficial on scorecards.
- RFC v0.1 scope marked shippable via T12.8.

## 2. Discovery Summary

Work type: Engineering (greenfield).

Graph scan: no `.code-review-graph/graph.db` yet; skipped. Manual scan: only
`docs/RFC-0001.md` present. Use cases derived from the RFC.

Use cases: 22 total (12 P0, 7 P1, 3 P2), all PLANNED.
Manifest: `.claude/scratch/usecases-manifest.json`.

Gaps to close in v0.1: protocol, engine, oracle, assertions/metrics, mcpfuzz,
Python SDK, SEC-001..007, CLI, two unofficial adapters, CI/release.

Research notes (from RFC + ADRs):
- Tech: Go engine; JSON-RPC stdio default; Python SDK; Oracle as HTTP base_url.
- Risk: adapter honesty, Oracle special-casing, soft-run cost, author bias.
- Arch: sidecar contract; scenario YAML + Go predicates; pillar GSI.

## 3. Scope and Deliverables

In scope:
- Full RFC release plan (v0.1 executable now; v0.2-v1.0 as outline epics).
- ADRs capturing architecture and open-question resolutions.
- Design doc and roadmap sync.

Out of scope:
- Changing Sire product code inside `sirerun/sire` except via the gavel adapter.
- Marketing site beyond dashboard/ in v0.3+.

| ID | Deliverable | Acceptance |
| ---- | ---- | ---- |
| D1 | Go module + AgentGavel CLI | build + version |
| D2 | Protocol + Python SDK | cross-language Handshake |
| D3 | Compliance Oracle | OpenAI + Anthropic shaped tool calls |
| D4 | SEC-001..007 suite | FakeAdapter oracle suite green |
| D5 | GSI scorecard + fingerprint | report command |
| D6 | Unofficial Sire + LangGraph adapters | provenance label |
| D7 | v0.1.0 GitHub release | tag + binaries |
| D8 | Later releases | outline epics E13-E15 |

## 4. Checkable Work Breakdown

Split layout. Frontier epics E1-E12 are `fidelity: executable`. E13-E15 are
`fidelity: outline` until prior release milestones complete.

### E1 -- Repository bootstrap  -> docs/plans/E1-repo-bootstrap.md  (6/6)

### E2 -- Adapter wire protocol  -> docs/plans/E2-adapter-protocol.md  (8/8)

### E3 -- Engine core  -> docs/plans/E3-engine-core.md  (6/7)

### E4 -- Compliance Oracle  -> docs/plans/E4-compliance-oracle.md  (6/6)

### E5 -- Assertions and GSI metrics  -> docs/plans/E5-assertions-metrics.md  (9/9)

### E6 -- mcpfuzz rogue MCP servers  -> docs/plans/E6-mcpfuzz.md  (9/9)

### E7 -- Python adapter SDK  -> docs/plans/E7-python-sdk.md  (5/6)

### E8 -- Security suite SEC-001 through SEC-007  -> docs/plans/E8-security-suite.md  (10/11)

### E9 -- CLI run, report, fingerprint  -> docs/plans/E9-cli-report.md  (2/6)

### E10 -- Sire adapter (unofficial)  -> docs/plans/E10-sire-adapter.md  (1/7)

### E11 -- LangGraph adapter (unofficial)  -> docs/plans/E11-langgraph-adapter.md  (0/7)

### E12 -- v0.1 quality gate and release  -> docs/plans/E12-v01-release.md  (1/8)

### E13 -- v0.2 expansion  -> docs/plans/E13-v02-expansion.md  (0/1)

### E14 -- v0.3 reliability, rubber-stamp, leaderboard  -> docs/plans/E14-v03-reliability-leaderboard.md  (0/1)

### E15 -- v1.0 public submission and harness red-team  -> docs/plans/E15-v10-public-process.md  (0/1)

## 5. Parallel Work

Tracks (v0.1):
- Track A: Bootstrap + protocol (E1, E2)
- Track B: Oracle (E4) -- parallel after T1.1
- Track C: Assertions/metrics (E5) -- parallel after T2.2
- Track D: mcpfuzz (E6) -- parallel after T1.1
- Track E: Engine (E3) -- after T2.4
- Track F: Python SDK (E7) -- after T2.3
- Track G: Suites (E8) -- after assertions + engine + mcpfuzz
- Track H: CLI (E9) -- after engine suite
- Track I/J: Adapters (E10, E11) -- after SDK + CLI
- Track K: Release (E12) -- after adapters

Sync points: T2.4 (protocol usable), T7.5 (cross-language), T8.10 (suite E2E),
T12.7 (quality gate).

### Wave 1: Foundation (6 agents)
- T1.1, T1.2, T1.5, T4.1, T6.1, T8.2

### Wave 2: Protocol + early packages (8 agents)
- T1.3, T1.4, T2.1, T4.2, T4.3, T5.1, T6.2, T6.3

### Wave 3: Protocol complete + fuzz modes (10 agents)
- T2.2, T2.3, T4.4, T5.2, T5.3, T6.4, T6.5, T6.6, T6.7, T1.6

### Wave 4: Session + metrics + SDK start (9 agents)
- T2.4, T2.5, T2.6, T5.4, T5.5, T5.6, T7.1, T4.5, T6.8

### Wave 5: Engine + SDK + suite loader (10 agents)
- T3.1, T3.2, T3.3, T3.4, T3.5, T5.7, T5.8, T7.2, T8.1, T2.7

### Wave 6: Integration seam (8 agents)
- T3.6, T7.3, T7.4, T5.9, T4.6, T6.9, T2.8, T3.7

### Wave 7: Scenarios batch A (7 agents)
- T7.5, T8.3, T8.4, T8.5, T8.6, T8.7, T7.6

### Wave 8: Scenarios batch B + CLI (6 agents)
- T8.8, T8.9, T9.1, T9.2, T8.10, T8.11

### Wave 9: Report + adapters scaffold (8 agents)
- T9.3, T9.4, T9.5, T10.1, T11.1, T10.2, T11.2, T9.6

### Wave 10: Adapter completion (8 agents)
- T10.3, T10.4, T11.3, T11.4, T10.5, T11.5, T10.6, T11.6

### Wave 11: Release gate (6 agents)
- T10.7, T11.7, T12.1, T12.2, T12.3, T12.5

### Wave 12: Ship (4 agents / humans)
- T12.4, T12.7, T12.6 (human), T12.8 (human)

### Wave 13: Next-horizon planning (1 agent)
- T13.0 (after T12.8)

## 6. Timeline and Milestones

| ID | Milestone | Depends | Exit criteria |
| ---- | ---- | ---- | ---- |
| M1 | Protocol + Oracle usable | Waves 1-4 | T2.4, T4.4 green |
| M2 | Engine + SDK cross-talk | Waves 5-6 | T7.5 green |
| M3 | Security suite SEC-001..007 | Waves 7-8 | T8.10 green |
| M4 | Unofficial adapters | Waves 9-10 | T10.5, T11.5 green |
| M5 | v0.1.0 release | Waves 11-12 | T12.8 done |
| M6 | v0.2 planned | T13.0 | E13 executable |
| M7 | v0.3 planned | T14.0 | E14 executable |
| M8 | v1.0 planned | T15.0 | E15 executable |

## 7. Risk Register

| ID | Risk | Impact | Likelihood | Mitigation |
| ---- | ---- | ---- | ---- | ---- |
| R1 | Author bias (Sire) | Credibility | Med | Unofficial label; external ratification ADR 007 |
| R2 | Soft runs expensive | Slow CI | Med | Oracle-first in CI; model mode opt-in job |
| R3 | Sire/LangGraph API mismatch | Adapter N/A heavy | High | Honest CapabilityReport; FakeAdapter proves harness |
| R4 | Stdio event framing bugs | Flaky suite | Med | Cross-language golden tests early (T7.5) |
| R5 | Q3 attestation sensitivity loss | False Soft/Hard | Med | Label context_mode; local raw default |
| R6 | Scope creep into v0.2 during v0.1 | Delay | Med | Outline epics only; reject SEC-008 until T13.0 |

## 8. Operating Procedure

Definition of done for a task:
1. Acceptance `acc:` predicate is green (tests or commands cited).
2. Paired tests exist for new packages/CLI/scenario behavior.
3. `gofmt` / `ruff` clean on touched trees; `go test` / `pytest` green.
4. PR merged to main via rebase; CI green.
5. For release tasks: GitHub release assets verified (no separate staging host
   in v0.1). Leaderboard "production" verification applies only from E14+.

Rules:
- Work in a git worktree once the repo has a first commit.
- Prefer stdlib Go (`flag`, `testing`, `net/http`).
- Never add per-framework exploit code; probes stay in `fixtures/`.
- Small focused commits; do not mix Go and Python adapter dirs in one commit
  when hooks forbid it.
- After M5, run T13.0 before any v0.2 coding.

## 9. Progress Log

- 2026-09-03: T8.11 suites/security gofmt verified clean (#39).

- 2026-09-03: T10.1 unofficial Sire adapter scaffold merged (#36).

- 2026-09-03: T7.6 ruff (#33) and T8.10 SuiteOracleFake (#34).

- 2026-09-03: T7.5 Go+Python FakeAdapter integration merged (#30).

- 2026-09-03: T8.9 SEC-007 composite fuzz merged (#27).

- 2026-09-03: T12.5 v0.1 devlog stub merged (#26).

- 2026-09-03: Wave 6 complete: T6.8 StartFuzzMode on wave-6-t6.8-integration.

- 2026-09-03: Wave 6 batch 2: T3.6 T7.3 T8.4 T8.6 T8.7 T8.8 on wave-6-batch2-integration.

- 2026-09-03: Wave 6 batch 1: T3.5 T7.2 T8.3 T9.3 on wave-6-integration (agent lane after kazi stuck).

- 2026-09-03: Wave 5: T3.2 T7.1 T8.1 T9.1 on wave-5-integration (kazi stuck;
  completed via agent-lane worktrees; fingerprint Model field fixed on land).
- 2026-09-03: Wave 2: T1.3 T1.4 T1.6 T2.1 T4.2 T4.3 T4.5 T6.2-T6.7.
- 2026-09-02: Wave 1 complete: T1.1 T1.2 T1.5 T8.2 T4.1 T6.1 on
  wave-1-integration (coordinator after Sonnet pool agents failed usage limits).
- 2026-09-02: Initial plan from RFC-0001. Created executable epics E1-E12,
  outline E13-E15, ADRs 001-008, design.md, usecases-manifest.json (22 UCs).
  Resolved RFC Q3-Q7 into ADRs for planning defaults.

## 10. Hand off Notes

- Spec: `docs/RFC-0001.md`. Design: `docs/design.md`.
- ADRs under `docs/adr/` -- especially transport (002), oracle (003), scoring
  (004), attestations (005), leaderboard (006), ratification (007).
- Start apply at Wave 1 tasks; do not expand E13 until T12.8.
- kazi is on PATH; engineering tasks carry `acc:` for JIT lane.
- Claim resource for plan rewrites: `R-plan-md`.
- Remote: `git@github.com:agentgavel/gavel.git` (empty history at plan time).

## 11. Appendix

- RFC open questions mapping: Q1/Q2 closed in RFC; Q3->ADR 005; Q4->ADR 006;
  Q5->ADR 004; Q6->ADR 007; Q7->ADR 002.
- Use case manifest: `.claude/scratch/usecases-manifest.json`.
- Release map: v0.1=E1-E12; v0.2=E13; v0.3=E14; v1.0=E15.
