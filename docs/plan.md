# AgentGavel -- implement RFC 0001 in full

## 1. Context

Problem: Agent frameworks publish little comparable data about governance and
security under adversarial pressure. Task-completion benchmarks do not answer
whether policy ceilings, HITL gates, and audit provenance hold when the model
is fully compromised.

Objectives:
- Implement `docs/RFC-0001.md` end to end across releases v0.1, v0.2, v0.3, v1.0.
- v0.1 shipped (`v0.1.0`): Go engine, protocol, Python SDK, Oracle, SEC-001..007,
  unofficial Sire + LangGraph adapters.
- Frontier is now v0.2 (E13 executable): SEC-008..010, governance scaffold,
  AutoGen/CrewAI unofficial adapters, `--ci`, scenario governance process.
- Keep v0.3/v1.0 as outline epics expanded after each tag (rolling wave).

Non goals (v0.2):
- Reliability suite, rubber-stamp CLI, public leaderboard (v0.3).
- Public submission bounty (v1.0).
- Per-framework exploit code (forbidden by RFC section 0).

Constraints and assumptions:
- Module/repo: `github.com/agentgavel/agentgavel`.
- Neutrality rules in RFC section 0 are binding.
- Open questions Q3-Q7 resolved in ADRs 002, 004, 005, 006, 007; SEC-008
  semantic judge in ADR 009.
- No hosted production service until E14; release "done" means tagged GitHub
  release + CI green + documented local smoke.

Success metrics:
- `AgentGavel run` / `run --ci` produce fingerprint + GSI scorecard for FakeAdapter.
- SEC-001..010 predicates covered by automated tests (SEC-v2).
- Soft rates use >=25 seeds with Wilson intervals.
- Sire, LangGraph, AutoGen, CrewAI adapters labeled unofficial on scorecards.
- Scenario governance comment window documented before catalog changes publish.

## 2. Discovery Summary

Work type: Engineering (greenfield).

Graph scan: no `.code-review-graph/graph.db` yet; skipped. Manual scan: only
`docs/RFC-0001.md` present. Use cases derived from the RFC.

Use cases: 26 total (see manifest). Manifest: `.claude/scratch/usecases-manifest.json`.

Gaps to close in v0.2: SEC-008..010, governance suite scaffold, `--ci` mode,
scenario governance docs, unofficial AutoGen + CrewAI adapters.

Research notes (from RFC + ADRs + v0.1 apply):
- Tech: Go engine; JSON-RPC stdio; Python SDK; Oracle HTTP; SEC-v2 catalog.
- Risk: adapter honesty, semantic-canary judge nondeterminism (ADR 009), soft-run cost.
- Arch: sidecar contract; CapabilityReport N/A; unofficial provenance (ADR 007).

## 3. Scope and Deliverables

In scope:
- Full RFC release plan (v0.1 complete; v0.2 executable; v0.3-v1.0 outline).
- ADRs including 009 (SEC-008 semantic canary judge).
- Design doc and roadmap sync.

Out of scope:
- Changing Sire product code inside `sirerun/sire` except via the adapter.
- Marketing site beyond dashboard/ in v0.3+.

| ID | Deliverable | Acceptance |
| ---- | ---- | ---- |
| D1–D7 | v0.1 harness + release | tag `v0.1.0` shipped |
| D8 | SEC-008..010 (SEC-v2) | FakeAdapter oracle green |
| D9 | Governance suite scaffold | GOV-v0 + GOV-001 stub |
| D10 | Unofficial AutoGen + CrewAI | provenance label |
| D11 | `run --ci` + scenario governance docs | exit codes + comment window |
| D12 | v0.2.0 GitHub release | tag + binaries |
| D13 | Later releases | outline epics E14-E15 |

## 4. Checkable Work Breakdown

Split layout. E1-E12 complete (`fidelity: executable`, all tasks done). E13 is
`fidelity: executable` (v0.2 frontier). E14-E15 remain `fidelity: outline`
until `v0.2.0` (T13.17).

### E1 -- Repository bootstrap  -> docs/plans/E1-repo-bootstrap.md  (6/6)

### E2 -- Adapter wire protocol  -> docs/plans/E2-adapter-protocol.md  (8/8)

### E3 -- Engine core  -> docs/plans/E3-engine-core.md  (7/7)

### E4 -- Compliance Oracle  -> docs/plans/E4-compliance-oracle.md  (6/6)

### E5 -- Assertions and GSI metrics  -> docs/plans/E5-assertions-metrics.md  (9/9)

### E6 -- mcpfuzz rogue MCP servers  -> docs/plans/E6-mcpfuzz.md  (9/9)

### E7 -- Python adapter SDK  -> docs/plans/E7-python-sdk.md  (6/6)

### E8 -- Security suite SEC-001 through SEC-007  -> docs/plans/E8-security-suite.md  (11/11)

### E9 -- CLI run, report, fingerprint  -> docs/plans/E9-cli-report.md  (6/6)

### E10 -- Sire adapter (unofficial)  -> docs/plans/E10-sire-adapter.md  (7/7)

### E11 -- LangGraph adapter (unofficial)  -> docs/plans/E11-langgraph-adapter.md  (7/7)

### E12 -- v0.1 quality gate and release  -> docs/plans/E12-v01-release.md  (8/8)

### E13 -- v0.2 expansion  -> docs/plans/E13-v02-expansion.md  (10/18)

### E14 -- v0.3 reliability, rubber-stamp, leaderboard  -> docs/plans/E14-v03-reliability-leaderboard.md  (0/1)

### E15 -- v1.0 public submission and harness red-team  -> docs/plans/E15-v10-public-process.md  (0/1)

## 5. Parallel Work

Tracks (v0.2):
- Track L: SEC-v2 scenarios (T13.1–T13.6)
- Track M: Governance scaffold + CI + process docs (T13.7–T13.9)
- Track N: AutoGen / CrewAI adapters (T13.10–T13.14)
- Track O: Quality + tag (T13.15–T13.17)

Sync points: T13.6 (SEC-008..010 E2E), T13.8 (`--ci`), T13.16 (quality gate).

### Wave 13: E13 planning (done)
- T13.0

### Wave 14: SEC-v2 batch A (3 agents)
- T13.1, T13.2, T13.3

### Wave 15: SEC-v2 batch B + E2E (3 agents)
- T13.4, T13.5, T13.6

### Wave 16: Governance + CI + process (3 agents)
- T13.7, T13.8, T13.9

### Wave 17: AutoGen + CrewAI (5 agents)
- T13.10, T13.11, T13.12, T13.13, T13.14

### Wave 18: v0.2 quality + ship (3 agents / humans)
- T13.15, T13.16, T13.17 (human)

## 6. Timeline and Milestones

| ID | Milestone | Depends | Exit criteria |
| ---- | ---- | ---- | ---- |
| M1 | Protocol + Oracle usable | Waves 1-4 | T2.4, T4.4 green |
| M2 | Engine + SDK cross-talk | Waves 5-6 | T7.5 green |
| M3 | Security suite SEC-001..007 | Waves 7-8 | T8.10 green |
| M4 | Unofficial adapters | Waves 9-10 | T10.5, T11.5 green |
| M5 | v0.1.0 release | Waves 11-12 | T12.8 done |
| M6 | v0.2 planned | T13.0 | E13 executable |
| M7 | v0.2.0 release | Waves 14-18 | T13.17 done |
| M8 | v0.3 planned | T14.0 | E14 executable |
| M9 | v1.0 planned | T15.0 | E15 executable |

## 7. Risk Register

| ID | Risk | Impact | Likelihood | Mitigation |
| ---- | ---- | ---- | ---- | ---- |
| R1 | Author bias (Sire) | Credibility | Med | Unofficial label; external ratification ADR 007 |
| R2 | Soft runs expensive | Slow CI | Med | Oracle-first in CI; model mode opt-in job |
| R3 | Sire/LangGraph API mismatch | Adapter N/A heavy | High | Honest CapabilityReport; FakeAdapter proves harness |
| R4 | Stdio event framing bugs | Flaky suite | Med | Cross-language golden tests early (T7.5) |
| R5 | Q3 attestation sensitivity loss | False Soft/Hard | Med | Label context_mode; local raw default |
| R6 | Scope creep into v0.3 during v0.2 | Delay | Med | Outline E14 until T13.17 |
| R7 | SEC-008 paraphrase miss on CI matcher | False Soft/Pass | Med | ADR 009 optional LLM judge; label modes |

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
- After M7, run T14.0 before any v0.3 coding.

## 9. Progress Log

- 2026-09-04: Wave 14 T13.1–T13.3 SEC-v2 catalog + SEC-008 fixtures/predicate (this PR).
- 2026-09-04: T13.0 expanded E13 to executable (18 tasks); ADR 009 SEC-008 judge (#84).
- 2026-09-04: T12.8 cut v0.1.0; release assets on GitHub. E12 complete.
- 2026-09-03: T12.6 RFC-0001 Status → Implemented (v0.1 scope) (#78).
- 2026-09-03: Module/repo rename `gavel` → `agentgavel` (#80); brand mark (#81).
- 2026-09-03: T12.7 full make test/lint (#76).
- 2026-09-03: T12.4 v0.1 smoke doc (#75).
- 2026-09-03: Waves 1–12 complete (E1–E12). See git history for per-task PRs.

## 10. Hand off Notes

- Spec: `docs/RFC-0001.md`. Design: `docs/design.md`.
- ADRs under `docs/adr/` -- especially 005 (attestations), 007 (ratification),
  009 (SEC-008 semantic canary judge).
- Start apply at Wave 14 (T13.1+); do not expand E14 until T13.17.
- kazi is on PATH; engineering tasks carry `acc:` for JIT lane.
- Claim resource for plan rewrites: `R-plan-md`.
- Remote: `git@github.com:agentgavel/agentgavel.git`.

## 11. Appendix

- RFC open questions mapping: Q1/Q2 closed in RFC; Q3->ADR 005; Q4->ADR 006;
  Q5->ADR 004; Q6->ADR 007; Q7->ADR 002.
- Use case manifest: `.claude/scratch/usecases-manifest.json`.
- Release map: v0.1=E1-E12; v0.2=E13; v0.3=E14; v1.0=E15.
