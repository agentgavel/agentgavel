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
- v0.3 shipped (`v0.3.0`): REL-v0, rubber-stamp (ADR 011), live leaderboard on
  `agentgavel.dev` (ADR 006), Pages + dashboard CI.
- Frontier is now v1.0 (E15 executable after T15.0): GitHub-native signed Opt-in
  (ADR 012 + ADR 013), harness red-team bounty, ratification ops (ADR 007).

Non goals (v1.0):
- Per-framework exploit code (forbidden by RFC section 0).
- Firebase / BaaS as Opt-in trust root (rejected; ADR 012).
- Hosted submission API beyond signed PRs into `dashboard/data/`.

Constraints and assumptions:
- Module/repo: `github.com/agentgavel/agentgavel`.
- Neutrality rules in RFC section 0 are binding.
- Open questions Q3-Q7 resolved in ADRs 002, 004, 005, 006, 007; SEC-008
  semantic judge in ADR 009; REL-v0 in ADR 010; rubber-stamp in ADR 011;
  GitHub-native Opt-in in ADR 012; Ed25519 signature format in ADR 013.
- No hosted production service beyond GitHub Pages dashboard; release "done"
  means tagged GitHub release + CI green + documented local smoke + live
  leaderboard URL.

Success metrics:
- Signed Opt-in PR path verifies in CI; `sample` no longer required when
  signature verifies (ADR 013).
- `SECURITY.md` + harness bounty scope published.
- At least one non-author adapter reaches provisional or ratified (ADR 007).
- Soft rates use >=25 seeds with Wilson intervals.
- Scenario governance comment window applies before REL/SEC/GOV changes publish.

## 2. Discovery Summary

Work type: Engineering (greenfield + process docs).

Graph scan: no `.code-review-graph/graph.db`; skipped.

Use cases: 30 total (UC-032 signed Opt-in submit, UC-033 harness bounty,
UC-034 adapter ratification ops added for v1.0). Manifest:
`.claude/scratch/usecases-manifest.json`.

Gaps to close in v1.0 (re-scanned 2026-09-05 after `v0.3.0`):
- Maintainer key registry + Ed25519 verify (ADR 013).
- CLI `report --sign` / `verify-entry`; CI verify job.
- Flip `check-dashboard.sh` Opt-in rule; allow signed `--tab opt-in`.
- Opt-in submission manual; Pages/README updates; ADR 006 addendum expiry.
- Harness bounty + SECURITY.md.
- Ratification ops doc + first non-author provisional (human gate).
- v1.0 smoke + quality gate + `v1.0.0` tag.

Research notes (v0.3 ops):
- Trust for Opt-in is cryptographic, not IdP login; volume is low.
- `report --publish` remains the writer of `index.json`; signing wraps the
  same entry schema.
- Author-affiliated Sire cannot skip to ratified via provisional alone.

## 3. Scope and Deliverables

In scope:
- Full RFC release plan (v0.1–v0.3 complete; v1.0 executable).
- ADRs through 013.
- Design doc and roadmap sync.

Out of scope:
- Changing Sire product code inside `sirerun/sire` except via the adapter.
- Firebase / custom hosted verify API as system of record (ADR 012).
- Expanding REL/SEC scenario catalogs (separate governance PRs).

| ID | Deliverable | Acceptance |
| ---- | ---- | ---- |
| D1–D7 | v0.1 harness + release | tag `v0.1.0` shipped |
| D8 | SEC-008..010 (SEC-v2) | FakeAdapter oracle green |
| D9 | Governance suite scaffold | GOV-v0 + GOV-001 stub |
| D10 | Unofficial §8.1 adapters (×6) | provenance label |
| D11 | `run --ci` + scenario governance docs | exit codes + comment window |
| D12 | v0.2.0 GitHub release | tag + binaries |
| D13 | REL-001..003 (REL-v0) | FakeAdapter oracle green |
| D14 | `rubber-stamp` CLI | SEC-002 + SEC-006; exit 0/1/2; both-N/A→1 (ADR 011) |
| D15 | Leaderboard dashboard | Opt-in + Unratified tabs; live on Pages |
| D16 | v0.3.0 GitHub release | tag + binaries |
| D17 | Signed Opt-in + bounty + ratification | ADR 012/013; SECURITY.md; provisional badge |
| D18 | v1.0.0 GitHub release | tag + binaries |

## 4. Checkable Work Breakdown

Split layout. E1–E14 complete (`fidelity: executable`, all tasks done). E15 is
`fidelity: executable` (v1.0 frontier).

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

### E14 -- v0.3 reliability, rubber-stamp, leaderboard  -> docs/plans/E14-v03-reliability-leaderboard.md  (24/24)

### E15 -- v1.0 public submission and harness red-team  -> docs/plans/E15-v10-public-process.md  (1/15)

## 5. Parallel Work

Tracks (v1.0):
- Track V: signature registry + verify library + CLI (T15.1–T15.3)
- Track W: CI + Opt-in rule flip + signed publish (T15.4–T15.6)
- Track X: submission docs + bounty + README (T15.7–T15.9)
- Track Y: ratification ops + first provisional (T15.10–T15.11)
- Track Z: smoke + quality + tag (T15.12–T15.14)

Sync points: T15.2 before T15.4/T15.5; T15.5+T15.6 before T15.13;
T15.10 before T15.11; T15.13 before T15.14.

### Wave 13–27: E13–E14 (done)
- T13.0–T13.25, T14.0–T14.23

### Wave 28: E15 planning (done)
- T15.0

### Wave 29: signature contract + CLI (3 agents)
- T15.1, T15.2, T15.3

### Wave 30: CI + Opt-in rule flip (3 agents)
- T15.4, T15.5, T15.6

### Wave 31: docs + bounty (3 agents)
- T15.7, T15.8, T15.9

### Wave 32: ratification ops (1 agent + founder)
- T15.10, T15.11

### Wave 33: quality gate + ship (2 agents + founder)
- T15.12, T15.13, T15.14

## Roadmap
- **Now:** Wave 28 done (T15.0); E15 1/15 executable — Wave 29 next
- **Next:** Wave 29 (T15.1–T15.3) → Waves 30–33 → `v1.0.0`

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
| M9 | v0.3.0 release | Waves 21-27 | T14.20 live, T14.16 done |
| M10 | v1.0 planned | T15.0 | E15 executable |
| M11 | v1.0.0 release | Waves 29-33 | T15.13 green, T15.14 tagged |

## 7. Risk Register

| ID | Risk | Impact | Likelihood | Mitigation |
| ---- | ---- | ---- | ---- | ---- |
| R1 | Author bias (Sire) | Credibility | Med | Unofficial label; external ratification ADR 007 |
| R2 | Soft runs expensive | Slow CI | Med | Oracle-first in CI; model mode opt-in job |
| R3 | Sire/LangGraph API mismatch | Adapter N/A heavy | High | Honest CapabilityReport; FakeAdapter proves harness |
| R10 | Pages spam without signatures | Credibility | Low | ADR 013 Ed25519 + CI verify; samples remain labeled |
| R11 | `rubber-stamp` vacuous green | False assurance | Med | ADR 011: both-N/A exits 1 |
| R14 | Canonical JSON mismatch across signers | Broken Opt-in | Med | Golden vectors in `internal/submit`; ADR 013 pins encoding |
| R15 | Provisional confused with ratified | Credibility | Med | Three-way badge + ratification manual (T15.10) |
| R16 | Bounty scope includes adapter exploits | Neutrality breach | Med | Bounty doc out-of-scope list; RFC section 0 |

## 8. Operating Procedure

Definition of done for a task:
1. Acceptance `acc:` predicate is green (tests or commands cited).
2. Paired tests exist for new packages/CLI/scenario behavior.
3. `gofmt` / `ruff` clean on touched trees; `go test` / `pytest` green.
4. PR merged to main via rebase; CI green.
5. For release tasks: GitHub release assets verified. Leaderboard live URL
   remains `https://agentgavel.dev/leaderboard/`.

Rules:
- Work in a git worktree once the repo has a first commit.
- Prefer stdlib Go (`flag`, `testing`, `net/http`, `crypto/ed25519`).
- Never add per-framework exploit code; probes stay in `fixtures/`.
- Small focused commits; do not mix Go and Python adapter dirs in one commit
  when hooks forbid it.
- Opt-in submissions are GitHub-native (ADR 012); signatures follow ADR 013.
- `report --publish` remains the sole writer of `dashboard/data/index.json`;
  signed Opt-in may set `tab=opt-in` only when verify passes.
- After M10, execute Waves 29–33 before cutting `v1.0.0` (T15.14).
- Adapter dirs for §8.1: `adk`, `openai_agents`, `pydantic_ai`,
  `agent_framework`, `strands`, `crewai` (never `autogen`).
- REL IDs and predicates: ADR 010 only (REL-001..003 / REL-v0).

## 9. Progress Log

- 2026-09-05: T15.0 expanded E15 to executable (15 tasks, Waves 29-33); ADR 013 Opt-in signature format; UC-032..034.
- 2026-09-05: Wave 27 T14.15/T14.20/T14.16; E14 complete; `v0.3.0` released.
- 2026-09-05: Wave 26 T14.21/T14.13/T14.14. E14 21/24.
- 2026-09-05: Wave 25 T14.10/T14.11/T14.12/T14.22. E14 18/24.
- 2026-09-04: E13 complete; `v0.2.0` released. Per-wave detail in `docs/roadmap.md`.

## 10. Hand off Notes

- Spec: `docs/RFC-0001.md`. Design: `docs/design.md`.
- ADRs: 006/012/013 (leaderboard + Opt-in), 007 (ratification), 010 (REL),
  011 (rubber-stamp).
- Start apply at Wave 29 (T15.1, T15.2, T15.3).
- Founder/human gates: T15.11 (provisional sign-off), T15.14 (tag `v1.0.0`).
- kazi is on PATH; engineering tasks carry `acc:` for JIT lane.
- Claim resource for plan rewrites: `R-plan-md`.
- Remote: `git@github.com:agentgavel/agentgavel.git`.

## 11. Appendix

- RFC open questions: Q3->ADR 005; Q4->ADR 006 + 012 + 013; Q5->ADR 004;
  Q6->ADR 007; Q7->ADR 002; REL->ADR 010; rubber-stamp->ADR 011.
- Use case manifest: `.claude/scratch/usecases-manifest.json`.
- Release map: v0.1=E1-E12; v0.2=E13; v0.3=E14; v1.0=E15.
