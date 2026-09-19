# AgentGavel -- live product runs (post E16 disclosure)

## 1. Context

Problem: E16 shipped ADR 015 runtime disclosure, FakeAdapter harness baselines,
and Sire live env bootstrap. Sire Wave 43 live oracle publish remains blocked
on founder credentials (T16.10). Product-comparable rankings still need a
non-Sire `runtime=live` path — preferably LangGraph after real outreach
(#8992).

Objectives:
- Expand and execute E17: live LangGraph (real PyPI package + interrupt).
- Optional OpenClaw/Hermes live probes when a gateway is operated.
- Keep stub default for CI; never market stub as product ranking.
- Leave Soft model-mode to E18; leave T15.11 until 2026-10-18.

Non goals:
- Completing Sire Wave 43 without credentials (still blocked on T16.10).
- Soft campaigns (>=25 seeds) — E18.
- Ratifying unofficial adapters early.
- Treating gateway stubs as live while hitl/withhold gaps remain.

Constraints:
- Worktrees; claim `R-plan-md` for this file.
- Secrets only in gitignored scratch / 1Password.
- kazi available; engineering tasks carry `acc:`.

Success metrics:
- LangGraph Handshake can report `runtime=live` with real package version.
- At least one Unratified LangGraph live scorecard path documented or published.
- Stub CI path remains green without `langgraph` installed.
- T15.11 still gated to 2026-10-18.

## 2. Discovery Summary

Work type: Engineering + ops runbooks.

Use cases (prior UC-001..039). This plan adds:
- UC-040 Run live LangGraph oracle benchmarks (real package)
- UC-041 Probe live OpenClaw/Hermes gateways with honest CapabilityReport

Wiring gaps:
- LangGraph adapter is stub-only (`framework_version=stub-0.0.1`).
- OpenClaw/Hermes hitl=false until withhold map / operated gateway.
- E16 T16.10–T16.16 still open (founder Sire creds).

## 3. Scope and Deliverables

In scope:
- E17 executable waves 45–48 (LangGraph live + optional gateways).
- benchmark-ops LangGraph live section; README runtime update.
- E16 remainder stays open until T16.10 unblocks Wave 43.

Out of scope:
- Soft model-mode (E18).
- Changing LangGraph upstream product code.

| ID | Deliverable | Acceptance |
| ---- | ---- | ---- |
| D26 | Live LangGraph adapter path | runtime=live + real interrupt/resume |
| D27 | LangGraph live Unratified evidence | summary + publish path |
| D28 | Gateway live probe scaffolding | env-gated; honest hitl |

## 4. Checkable Work Breakdown

### E15 -- v1.0 Opt-in, red-team, OpenClaw, Hermes  -> docs/plans/E15-v10-public-process.md  (30/31)

### E16 -- First live benchmarks  -> docs/plans/E16-first-live-benchmarks.md  (13/18)

### E17 -- Live LangGraph and gateway product runs  -> docs/plans/E17-live-langgraph-gateways.md  (13/14; T17.11 human gateway open)

### E18 -- Soft model-mode and broader live adapter set  -> docs/plans/E18-soft-and-broader-live.md  (8/10; Soft runs blocked on T18.3)

Completed earlier epics (audit): E1–E14 under `docs/plans/`.

## 5. Parallel Work

Tracks (E17):
- Track L: LangGraph live T17.1–T17.8 (after T17.0)
- Track G: Gateway probes T17.9–T17.11 (parallel after T17.0; human T17.11)
- Track Q: quality gate T17.12–T17.14

E16 Track P (Wave 43) remains blocked on T16.10 — resume when creds appear.

### Wave 45: LangGraph live bootstrap
- T17.0, T17.1, T17.2, T17.3, T17.4

### Wave 46: LangGraph oracle evidence
- T17.5, T17.6, T17.7, T17.8

### Wave 47: Gateway probes
- T17.9, T17.10, T17.11 (human)

### Wave 48: E17 quality gate
- T17.12, T17.13, T17.14

## 6. Timeline and Milestones

| ID | Milestone | Depends | Exit criteria |
| ---- | ---- | ---- | ---- |
| M14 | First live Unratified Sire row | Wave 43 | T16.15 (blocked T16.10) |
| M16 | LangGraph provisional | T15.11 | after 2026-10-18 or #8992 |
| M17 | E17 planned | T17.0 | done 2026-09-19 |
| M18 | Live LangGraph bootstrap | Wave 45 | T17.4 green |
| M19 | LangGraph live Unratified evidence | Wave 46 | T17.7 |
| M20 | E17 closed | Wave 48 | T17.12 green |

## 7. Risk Register

| ID | Risk | Impact | Likelihood | Mitigation |
| ---- | ---- | ---- | ---- | ---- |
| R1 | Author bias (Sire) | Credibility | Med | unofficial + runtime=live |
| R20 | Stub scores misread as product | Credibility | High | ADR 015 + README |
| R23 | T16.10 blocked long | Schedule | High | E17 LangGraph live unblocks non-Sire path |
| R24 | langgraph PyPI heavy for CI | CI time | Med | optional extra; stub default |
| R25 | Gateway withhold unmapped | False hitl | High | keep hitl=false until map exists |
| R15 | Provisional confused with ratified | Credibility | Med | wait #175 window |

## 8. Operating Procedure

Same as E16: worktrees, no per-framework exploits, honest runtime/provenance,
gofmt/ruff, rebase PRs.

## 9. Progress Log

- 2026-09-19: T18.9 — frontier updated. Soft deferred (T18.5); all §8.1 Python adapters have optional live env bootstrap (langgraph/crewai/adk/pydantic_ai/openai_agents/strands/agent_framework + openclaw/hermes probes). Remaining agent work blocked on founder: T16.10 Sire, T17.11 gateway, T18.3 Soft model API.
- 2026-09-19: Soft Unratified deferred (T18.5) — `--mode model` not implemented; no model API creds; Hard/oracle LangGraph live sample already covers non-Sire Unratified without Soft spend. Pydantic AI live env bootstrap added.
- 2026-09-19: E18 expanded (Soft runbook + model-mode fail-closed test + CI ruff matrix for langgraph/hermes). Soft implementation still blocked on T18.3 model creds; T18.6 CrewAI/ADK live pending.
- 2026-09-19: Wave 47 scaffolding — OpenClaw Gateway health probe + Hermes capabilities probe (env-gated, fail-closed); hitl stays false on OpenClaw. T17.11 human gateway + T16.10 Sire creds still open.
- 2026-09-19: Waves 45–46 — live LangGraph bootstrap + SEC-001/007 oracle evidence + Unratified sample (`runtime=live`). E17 9/14; Wave 47 gateway probes next; E16 Wave 43 still blocked on T16.10.
- 2026-09-19: T17.0 — E17 expanded to executable (14 tasks, waves 45–48); frontier shifts to live LangGraph while E16 Wave 43 waits on T16.10.
- 2026-09-18: Waves 41-42 + quality gate landed (benchmark-ops, FakeAdapter baselines, Sire env live bootstrap). T16.10/T16.12-16 blocked on founder Sire credentials.
- 2026-09-18: Wave 40 shipped -- T16.1-T16.4 runtime disclosure (ADR 015 wire + adapters + dashboard).
- 2026-09-18: /plan first live benchmarks -- trim focus to E15 remainder + E16 executable + E17/E18 outline; ADR 015; UC-037..039.
- 2026-09-18: T15.11 real outreach langchain-ai/langgraph#8992; window closes 2026-10-18.
- 2026-09-07: v1.0.0 tagged (T15.14); E15 30/31.

## 10. Hand off Notes

- Next /apply: E16 Wave 43 when T16.10 Sire creds appear; Soft T18.4 when T18.3 model API exists; T17.11 gateway when operated.
- Founder: T16.10 Sire creds; T17.11 gateway URL; T18.3 Soft model endpoint; T15.11 after 2026-10-18.
- PR #182 holds E17/E18 live bootstrap work (CI green). Merge when ready.

## 11. Appendix

- Spec: docs/RFC-0001.md
- Design: docs/design.md
- ADR 015: docs/adr/015-stub-vs-live-runtime.md
- Roadmap: docs/roadmap.md
- Use cases: .claude/scratch/usecases-manifest.json
- Release map: v0.1=E1-E12; v0.2=E13; v0.3=E14; v1.0=E15; live=E16+
