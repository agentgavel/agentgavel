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
- v0.2 shipped (`v0.2.0`): SEC-008..010, GOV-v0 scaffold, six RFC §8.1 unofficial
  adapters, `run --ci`, scenario governance.
- Frontier is now v0.3 (E14 executable): REL-v0 (ADR 010), `rubber-stamp`,
  GitHub Pages leaderboard (ADR 006).
- Keep v1.0 as an outline epic expanded after the v0.3 tag (rolling wave).

Non goals (v0.3):
- Public submission bounty / harness red-team (v1.0 / E15).
- Per-framework exploit code (forbidden by RFC section 0).
- Hosted submission API beyond static Pages + `report --publish` JSON.

Constraints and assumptions:
- Module/repo: `github.com/agentgavel/agentgavel`.
- Neutrality rules in RFC section 0 are binding.
- Open questions Q3-Q7 resolved in ADRs 002, 004, 005, 006, 007; SEC-008
  semantic judge in ADR 009; REL-v0 scenarios in ADR 010.
- No hosted production service beyond GitHub Pages dashboard; release "done"
  means tagged GitHub release + CI green + documented local smoke.

Success metrics:
- `AgentGavel run` / `run --ci` / `rubber-stamp` produce fingerprint + scorecard
  for FakeAdapter with documented exit codes.
- SEC-001..010 and REL-001..003 predicates covered by automated tests.
- Soft rates use >=25 seeds with Wilson intervals.
- Dashboard Opt-in vs Unratified tabs labeled per ADR 006; sample provenance honest.
- Scenario governance comment window applies before REL/SEC/GOV changes publish.

## 2. Discovery Summary

Work type: Engineering (greenfield).

Graph scan: no `.code-review-graph/graph.db` yet; skipped. Manual scan: only
`docs/RFC-0001.md` present. Use cases derived from the RFC.

Use cases: 27 total (see manifest; UC-031 added for REL-v0). Manifest:
`.claude/scratch/usecases-manifest.json`.

Gaps to close in v0.3: REL-001..003 (ADR 010), `rubber-stamp` CLI, static
dashboard with Opt-in/Unratified tabs, `report --publish`, Pages publish docs.

Research notes (from RFC + ADRs + v0.2 apply):
- Tech: Go engine; JSON-RPC stdio; Python SDK; Oracle HTTP; SEC-v2; GOV-v0; REL-v0.
- Risk: REL definitions were RFC-outline only until ADR 010; soft-run cost; Pages spam (ADR 006 signatures deferred to v1.0).
- Arch: sidecar contract; CapabilityReport N/A; unofficial provenance (ADR 007);
  `--ci` exit-code map reused by rubber-stamp.

## 3. Scope and Deliverables

In scope:
- Full RFC release plan (v0.1–v0.2 complete; v0.3 executable; v1.0 outline).
- ADRs including 009 (SEC-008 judge) and 010 (REL-v0).
- Design doc and roadmap sync.

Out of scope:
- Changing Sire product code inside `sirerun/sire` except via the adapter.
- Marketing site beyond `dashboard/` static Pages.
- v1.0 signed-submission API and harness bounty (E15).

| ID | Deliverable | Acceptance |
| ---- | ---- | ---- |
| D1–D7 | v0.1 harness + release | tag `v0.1.0` shipped |
| D8 | SEC-008..010 (SEC-v2) | FakeAdapter oracle green |
| D9 | Governance suite scaffold | GOV-v0 + GOV-001 stub |
| D10 | Unofficial §8.1 adapters (×6) | provenance label |
| D11 | `run --ci` + scenario governance docs | exit codes + comment window |
| D12 | v0.2.0 GitHub release | tag + binaries |
| D13 | REL-001..003 (REL-v0) | FakeAdapter oracle green |
| D14 | `rubber-stamp` CLI | SEC-002 + SEC-006; exit 0/1/2 |
| D15 | Leaderboard dashboard | Opt-in + Unratified tabs (ADR 006) |
| D16 | v0.3.0 GitHub release | tag + binaries |
| D17 | v1.0 | outline epic E15 |

## 4. Checkable Work Breakdown

Split layout. E1-E13 complete (`fidelity: executable`, all tasks done). E14 is
`fidelity: executable` (v0.3 frontier). E15 remains `fidelity: outline` until
`v0.3.0` (T14.16).

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

### E13 -- v0.2 expansion  -> docs/plans/E13-v02-expansion.md  (26/26)

### E14 -- v0.3 reliability, rubber-stamp, leaderboard  -> docs/plans/E14-v03-reliability-leaderboard.md  (7/17)

### E15 -- v1.0 public submission and harness red-team  -> docs/plans/E15-v10-public-process.md  (0/1)

## 5. Parallel Work

Tracks (v0.3):
- Track Q: REL-v0 catalog + fixtures (T14.1–T14.2)
- Track R: REL predicates + E2E (T14.3–T14.6)
- Track S: rubber-stamp CLI + docs (T14.7–T14.8)
- Track T: dashboard / Pages (T14.9–T14.12)
- Track U: Smoke + quality + tag (T14.13–T14.16)

Sync points: T14.6 (REL E2E), T14.7 (`rubber-stamp`), T14.15 (quality gate).

### Wave 13–19: E13 (done)
- T13.0–T13.25

### Wave 20: E14 planning (done)
- T14.0

### Wave 21: REL catalog + fixtures (2 agents)
- T14.1, T14.2

### Wave 22: REL predicates + E2E (4 agents)
- T14.3, T14.4, T14.5, T14.6

### Wave 23: rubber-stamp (2 agents)
- T14.7, T14.8

### Wave 24: dashboard / leaderboard (4 agents)
- T14.9, T14.10, T14.11, T14.12

### Wave 25: v0.3 smoke + quality + ship (4 agents / humans)
- T14.13, T14.14, T14.15, T14.16

## Roadmap
- **Now:** Wave 22 done (T14.3–T14.6); Wave 23 next (T14.7–T14.8 rubber-stamp)
- **Next:** rubber-stamp + dashboard → T14.16 `v0.3.0`

## 6. Timeline and Milestones

| ID | Milestone | Depends | Exit criteria |
| ---- | ---- | ---- | ---- |
| M1 | Protocol + Oracle usable | Waves 1-4 | T2.4, T4.4 green |
| M2 | Engine + SDK cross-talk | Waves 5-6 | T7.5 green |
| M3 | Security suite SEC-001..007 | Waves 7-8 | T8.10 green |
| M4 | Unofficial adapters | Waves 9-10 | T10.5, T11.5 green |
| M5 | v0.1.0 release | Waves 11-12 | T12.8 done |
| M6 | v0.2 planned | T13.0 | E13 executable |
| M7 | v0.2.0 release | Waves 14-19 | T13.25 done |
| M8 | v0.3 planned | T14.0 | E14 executable |
| M9 | v0.3.0 release | Waves 21-25 | T14.16 done |
| M10 | v1.0 planned | T15.0 | E15 executable |

## 7. Risk Register

| ID | Risk | Impact | Likelihood | Mitigation |
| ---- | ---- | ---- | ---- | ---- |
| R1 | Author bias (Sire) | Credibility | Med | Unofficial label; external ratification ADR 007 |
| R2 | Soft runs expensive | Slow CI | Med | Oracle-first in CI; model mode opt-in job |
| R3 | Sire/LangGraph API mismatch | Adapter N/A heavy | High | Honest CapabilityReport; FakeAdapter proves harness |
| R4 | Stdio event framing bugs | Flaky suite | Med | Cross-language golden tests early (T7.5) |
| R5 | Q3 attestation sensitivity loss | False Soft/Hard | Med | Label context_mode; local raw default |
| R6 | Scope creep into v0.3 during v0.2 | Delay | Med | Outline E14 until T13.25 (closed) |
| R8 | Six §8.1 adapters stretch v0.2 | Delay | Med | Parallel waves 17–18; LangGraph-style stubs keep CI light |
| R7 | SEC-008 paraphrase miss on CI matcher | False Soft/Pass | Med | ADR 009 optional LLM judge; label modes |
| R9 | REL-v0 too thin vs Fault Recovery pillar | Weak GSI | Med | ADR 010 pins three scenarios; comment window before REL-v1 |
| R10 | Pages spam without v1.0 signatures | Credibility | Low | ADR 006 Unratified tab; samples only until E15 |

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
- After M8, execute Waves 21–25 before cutting `v0.3.0` (T14.16).
- After M9, run T15.0 before any v1.0 coding.
- Adapter dirs for §8.1: `adk`, `openai_agents`, `pydantic_ai`,
  `agent_framework`, `strands`, `crewai` (never `autogen`).
- REL IDs and predicates: ADR 010 only (REL-001..003 / REL-v0).

## 9. Progress Log

- 2026-09-04: Wave 21 T14.1 REL-v0 loader + T14.2 REL fixtures shipped.
- 2026-09-04: T14.0 expanded E14 to executable; ADR 010 REL-v0.
- 2026-09-04: T13.25 cut v0.2.0; GitHub release assets uploaded. E13 complete (26/26).
- 2026-09-04: T13.24 make test && make lint green on clean tree; E13 25/26.
- 2026-09-04: Wave 19 T13.22 smoke + T13.23 gofmt/ruff clean shipped; E13 24/26.
- 2026-09-04: Wave 18 T13.16–T13.21 §8.1 adapter tool paths shipped (adk, openai_agents, pydantic_ai, agent_framework, strands, crewai).
- 2026-09-04: Wave 17 T13.10–T13.15 §8.1 adapter scaffolds shipped (adk, openai_agents, pydantic_ai, agent_framework, strands, crewai).
- 2026-09-04: Retargeted E13 Wave 4–6 to RFC §8.1 (six adapters; T13.10–T13.25).
- 2026-09-04: Wave 16 T13.7–T13.9 GOV + `--ci` + scenario-governance docs.
- 2026-09-04: Wave 14 T13.1–T13.3 SEC-v2 catalog + SEC-008 fixtures/predicate.
- 2026-09-04: T13.0 expanded E13 to executable; ADR 009 SEC-008 judge (#84).
- 2026-09-04: T12.8 cut v0.1.0; release assets on GitHub. E12 complete.
- 2026-09-03: T12.6 RFC-0001 Status → Implemented (v0.1 scope) (#78).
- 2026-09-03: Module/repo rename `gavel` → `agentgavel` (#80); brand mark (#81).
- 2026-09-03: T12.7 full make test/lint (#76).
- 2026-09-03: T12.4 v0.1 smoke doc (#75).
- 2026-09-03: Waves 1–12 complete (E1–E12). See git history for per-task PRs.

## 10. Hand off Notes

- Spec: `docs/RFC-0001.md`. Design: `docs/design.md`.
- ADRs under `docs/adr/` -- especially 005 (attestations), 006 (leaderboard),
  007 (ratification), 009 (SEC-008 judge), 010 (REL-v0).
- Start apply at Wave 21 (T14.1–T14.2); E14 is executable.
- kazi is on PATH; engineering tasks carry `acc:` for JIT lane.
- Claim resource for plan rewrites: `R-plan-md`.
- Remote: `git@github.com:agentgavel/agentgavel.git`.

## 11. Appendix

- RFC open questions mapping: Q1/Q2 closed in RFC; Q3->ADR 005; Q4->ADR 006;
  Q5->ADR 004; Q6->ADR 007; Q7->ADR 002; REL suite->ADR 010.
- Use case manifest: `.claude/scratch/usecases-manifest.json`.
- Release map: v0.1=E1-E12; v0.2=E13; v0.3=E14; v1.0=E15.
