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
- Frontier is now v0.3 (E14 executable, 7/24 after Wave 22): REL-v0 shipped
  as predicates + FakeAdapter E2E (ADR 010); remaining work is `rubber-stamp`
  (ADR 011), REL CLI wiring, GSI pillar mapping, and the GitHub Pages
  leaderboard (ADR 006 + addendum).
- Keep v1.0 as an outline epic expanded after the v0.3 tag (rolling wave).

Non goals (v0.3):
- Public submission bounty / harness red-team (v1.0 / E15).
- Per-framework exploit code (forbidden by RFC section 0).
- Hosted submission API beyond static Pages + `report --publish` JSON.

Constraints and assumptions:
- Module/repo: `github.com/agentgavel/agentgavel`.
- Neutrality rules in RFC section 0 are binding.
- Open questions Q3-Q7 resolved in ADRs 002, 004, 005, 006, 007; SEC-008
  semantic judge in ADR 009; REL-v0 scenarios in ADR 010; `rubber-stamp`
  verdict contract (no GSI, both-N/A exits 1) in ADR 011.
- No hosted production service beyond GitHub Pages dashboard; release "done"
  means tagged GitHub release + CI green + documented local smoke.

Success metrics:
- `AgentGavel run` / `run --ci` / `rubber-stamp` produce fingerprint + scorecard
  for FakeAdapter with documented exit codes.
- SEC-001..010 and REL-001..003 predicates covered by automated tests.
- Soft rates use >=25 seeds with Wilson intervals.
- Dashboard Opt-in vs Unratified tabs labeled per ADR 006; sample provenance
  honest; the site is live on GitHub Pages before `v0.3.0` (T14.20).
- Scenario governance comment window applies before REL/SEC/GOV changes publish.

## 2. Discovery Summary

Work type: Engineering (greenfield).

Graph scan: no `.code-review-graph/graph.db` yet; skipped. Manual scan: only
`docs/RFC-0001.md` present. Use cases derived from the RFC.

Use cases: 27 total (see manifest; UC-031 added for REL-v0). Manifest:
`.claude/scratch/usecases-manifest.json`.

Gaps to close in v0.3 (re-scanned 2026-09-04 after Wave 22):
- `rubber-stamp` CLI (T14.7, ADR 011) and its manual (T14.8).
- `run --suite reliability` is rejected by `cmd/AgentGavel/run.go` although
  `reliability.RunOracleFakeREL` exists; `protocol.ScenarioNA` has no REL
  entries (T14.17).
- `internal/metrics.ScenarioPillar` drops REL-* IDs from GSI; design.md places
  them in the resilience pillar (T14.18).
- Static dashboard + entry schema + samples + `report --publish` + index
  (T14.9–T14.11), Pages workflow (T14.19), dashboard CI check (T14.21).
- No Pages site exists yet (`gh api .../pages` → 404): enabling it is a
  founder action and `v0.3.0` waits on the live URL (T14.20).
- Docs: Pages publish path, v0.3 smoke, README, scenario governance for REL
  (T14.12, T14.13, T14.22, T14.23).

Research notes (from RFC + ADRs + v0.2 apply):
- Tech: Go engine; JSON-RPC stdio; Python SDK; Oracle HTTP; SEC-v2; GOV-v0; REL-v0.
- Risk: REL definitions were RFC-outline only until ADR 010; soft-run cost; Pages spam (ADR 006 signatures deferred to v1.0).
- Arch: sidecar contract; CapabilityReport N/A; unofficial provenance (ADR 007);
  `--ci` exit-code map reused by rubber-stamp (ADR 011); static site needs an
  `index.json` written by the same code path as the entries.

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
| D14 | `rubber-stamp` CLI | SEC-002 + SEC-006; exit 0/1/2; both-N/A→1 (ADR 011) |
| D15 | Leaderboard dashboard | Opt-in + Unratified tabs (ADR 006); live on Pages (T14.20) |
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

### E14 -- v0.3 reliability, rubber-stamp, leaderboard  -> docs/plans/E14-v03-reliability-leaderboard.md  (11/24)

### E15 -- v1.0 public submission and harness red-team  -> docs/plans/E15-v10-public-process.md  (0/1)

## 5. Parallel Work

Tracks (v0.3):
- Track Q: REL-v0 catalog + fixtures (T14.1–T14.2) — done
- Track R: REL predicates + E2E (T14.3–T14.6) — done
- Track S: rubber-stamp + REL CLI/GSI wiring + governance doc (T14.7, T14.8, T14.17, T14.18, T14.23)
- Track T: dashboard / publish / Pages (T14.9–T14.12, T14.19, T14.21, T14.22)
- Track U: smoke + quality + live + tag (T14.13–T14.16, T14.20)

Sync points: T14.7 + T14.17 (CLI surfaces), T14.11 (publish schema), T14.15
(quality gate), T14.20 (live leaderboard) before T14.16.

### Wave 13–19: E13 (done)
- T13.0–T13.25

### Wave 20: E14 planning (done)
- T14.0

### Wave 21: REL catalog + fixtures (2 agents)
- T14.1, T14.2

### Wave 22: REL predicates + E2E (4 agents)
- T14.3, T14.4, T14.5, T14.6

### Wave 23: rubber-stamp + REL CLI/GSI wiring (4 agents)
- T14.7, T14.17, T14.18, T14.23

### Wave 24: rubber-stamp docs + dashboard scaffold + Pages workflow (3 agents)
- T14.8, T14.9, T14.19

### Wave 25: dashboard samples + publish + docs (4 agents)
- T14.10, T14.11, T14.12, T14.22

### Wave 26: dashboard CI check + v0.3 smoke + gofmt (3 agents)
- T14.21, T14.13, T14.14

### Wave 27: quality gate + live leaderboard + ship (2 agents + founder)
- T14.15, T14.20, T14.16

## Roadmap
- **Now:** Wave 23 done (T14.7, T14.17, T14.18, T14.23); Wave 24 next (T14.8, T14.9, T14.19)
- **Next:** Waves 24–26 dashboard/publish/Pages + smoke → Wave 27 live
  leaderboard (T14.20) → T14.16 `v0.3.0` → T15.0 expands E15

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
| R10 | Pages spam without v1.0 signatures | Credibility | Low | ADR 006 addendum: publish writes Unratified only; opt-in ⇒ sample (CI check T14.21) |
| R11 | `rubber-stamp` vacuous green on hitl=false adapters | False assurance | Med | ADR 011: both-N/A exits 1 with not_applicable |
| R12 | Pages enablement is a founder action; v0.3.0 stalls | Delay | Med | T14.20 scheduled with T14.15; workflow + docs land first so the action is one command |
| R13 | REL rows silently excluded from GSI | Misleading scorecard | High | T14.18 pillar mapping + report test before smoke |

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
- After M8, execute Waves 21–27 before cutting `v0.3.0` (T14.16); the tag
  waits on the live leaderboard (T14.20), not only on green CI.
- `rubber-stamp` follows ADR 011: SEC-002 + SEC-006 only, `ciExitCode`
  unchanged, both-N/A exits 1, never prints GSI.
- `report --publish` writes Unratified entries only until E15 (ADR 006
  addendum); it is the sole writer of `dashboard/data/index.json`.
- After M9, run T15.0 before any v1.0 coding.
- Adapter dirs for §8.1: `adk`, `openai_agents`, `pydantic_ai`,
  `agent_framework`, `strands`, `crewai` (never `autogen`).
- REL IDs and predicates: ADR 010 only (REL-001..003 / REL-v0); they score in
  the resilience pillar with RFC §6 weights unchanged (design.md).

## 9. Progress Log

- 2026-09-04: /plan refine after Wave 22 — E14 7/24: refined T14.7–T14.16, added T14.17–T14.23 (REL CLI wiring, REL pillar mapping, Pages workflow, publish index, dashboard CI check, README, governance-for-REL); waves re-cut to 23–27; ADR 011 (rubber-stamp verdict) + ADR 006 addendum.
- 2026-09-04: Wave 22 T14.3–T14.6 REL-001..003 predicates + FakeAdapter oracle E2E shipped (#112, #113).
- 2026-09-04: Wave 21 T14.1 REL-v0 loader + T14.2 REL fixtures shipped.
- 2026-09-04: T14.0 expanded E14 to executable; ADR 010 REL-v0.
- 2026-09-04: E13 complete (26/26); `v0.2.0` released (T13.25). Per-wave detail in `docs/roadmap.md`.
- 2026-09-03: E1–E12 complete; `v0.1.0` released (T12.8). Per-task PRs in git history.

## 10. Hand off Notes

- Spec: `docs/RFC-0001.md`. Design: `docs/design.md`.
- ADRs under `docs/adr/` -- especially 005 (attestations), 006 (leaderboard),
  006 addendum (v0.3 publish tab rule), 007 (ratification), 009 (SEC-008 judge), 010 (REL-v0),
  011 (rubber-stamp verdict).
- Start apply at Wave 23 (T14.7, T14.17, T14.18, T14.23); E14 is executable.
- Founder actions in E14: T14.20 (enable Pages, verify live URL) and T14.16
  (tag `v0.3.0`). Everything else is agent/kazi work.
- kazi is on PATH; engineering tasks carry `acc:` for JIT lane.
- Claim resource for plan rewrites: `R-plan-md`.
- Remote: `git@github.com:agentgavel/agentgavel.git`.

## 11. Appendix

- RFC open questions mapping: Q1/Q2 closed in RFC; Q3->ADR 005; Q4->ADR 006;
  Q5->ADR 004; Q6->ADR 007; Q7->ADR 002; REL suite->ADR 010;
  rubber-stamp verdict->ADR 011; v0.3 publish tab rule->ADR 006 addendum.
- Use case manifest: `.claude/scratch/usecases-manifest.json`.
- Release map: v0.1=E1-E12; v0.2=E13; v0.3=E14; v1.0=E15.
