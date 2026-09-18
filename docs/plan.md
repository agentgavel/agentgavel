# AgentGavel -- first live benchmarks (post v1.0)

## 1. Context

Problem: AgentGavel v1.0.0 is released and can run suites, but almost all
framework adapters are in-process stubs. Running them today measures the
harness shape, not product control planes. Publishing those rows as
framework rankings would damage credibility.

Objectives:
- Disclose stub vs live vs harness on every Handshake and scorecard (ADR 015).
- Establish a FakeAdapter oracle baseline (harness proof).
- Wire Sire to a live worker + Oracle and publish the first Unratified
  `runtime=live` / `provenance=unofficial` scorecard.
- Keep T15.11 LangGraph provisional gated until 2026-10-18 after real
  outreach (langchain-ai/langgraph#8992 / agentgavel#175).
- Defer live LangGraph/gateway product runs and Soft model-mode to outline
  epics E17/E18.

Non goals (this frontier):
- Treating stub adapter scores as product rankings.
- Soft model-mode (>=25 seeds) campaigns (E18).
- Live LangGraph PyPI or live OpenClaw/Hermes gateways (E17).
- Expanding SEC/REL/GOV scenario catalogs.
- n8n / Dify adapters (still RFC section 8.2 deferred).
- Ratifying author-affiliated Sire without independent review (ADR 007).

Constraints and assumptions:
- Module/repo: `github.com/agentgavel/agentgavel`.
- Neutrality rules in RFC section 0 remain binding.
- Secrets for Sire live runs stay in gitignored scratch or 1Password; never
  commit tokens.
- Work in a task worktree; claim `R-plan-md` when rewriting this file.
- kazi is available; engineering tasks carry `acc:` predicates.

Success metrics:
- CapabilityReport + published JSON include `runtime` in {stub, live, harness}.
- FakeAdapter full oracle SEC+REL baseline documented and reproducible.
- At least one Sire live Unratified row on agentgavel.dev with
  runtime=live and provenance=unofficial.
- README warns operators that stub != product ranking.
- T15.11 remains blocked until 2026-10-18 (or earlier only if #8992 ratifies).

## 2. Discovery Summary

Work type: Engineering + operations (benchmark operator runbooks).

Graph: `.code-review-graph` optional; manual CLI/adapter scan used.

Use cases: prior manifest UC-001..036 (v0.1-v1.0). This plan adds:
- UC-037 Disclose adapter runtime class on Handshake and scorecards
- UC-038 Publish Unratified rows that carry runtime + provenance
- UC-039 Run live Sire oracle benchmarks (dogfood)

Wiring gaps:
- No `runtime` field on CapabilityReport yet (ADR 015 new).
- Sire HttpSireClient exists but has no env bootstrap; default is stub.
- LangGraph/CrewAI/ADK/etc. explicitly do not depend on real packages.
- OpenClaw/Hermes live paths need operated gateways (N/A-heavy stubs today).

Research findings:
- FakeAdapter path is the proven harness proof (README / v0.1-v0.3 smoke).
- Fingerprints already hash framework.version; stubs often report
  stub-0.0.1 -- insufficient alone for Pages disclosure (ADR 015).
- Unratified publish path: `report --publish` + dashboard PR
  (`docs/manual/leaderboard-pages.md`).

## 3. Scope and Deliverables

In scope:
- ADR 015 + wire `runtime` through protocol, fingerprints, dashboard.
- FakeAdapter baseline + benchmark-ops manual.
- Sire live env bootstrap + first live oracle runs + Unratified publish.
- Plan trim: E1-E14 complete at v1.0; E15 remains for T15.11 only;
  E16 frontier; E17/E18 outline.

Out of scope:
- Product marketing as ranked leaderboards from stub runs.
- Changing Sire product code outside the adapter.
- Firebase / hosted verify APIs.

| ID | Deliverable | Acceptance |
| ---- | ---- | ---- |
| D21 | ADR 015 + runtime on wire/scorecards | Handshake + report JSON + dashboard schema |
| D22 | FakeAdapter oracle baseline manual | docs/manual/benchmark-ops.md copy-paste green |
| D23 | Sire live bootstrap | env-gated HttpSireClient; stub default preserved |
| D24 | First Sire live Unratified publish | runtime=live provenance=unofficial on Pages |
| D25 | E17/E18 outline epics | planning tasks only until E16 gate |

## 4. Checkable Work Breakdown

Split layout. E1-E14 complete (v0.1-v0.3 / early v1.0 work; retained under
`docs/plans/` for audit). Frontier is E16. E15 still open on human T15.11.

### E15 -- v1.0 Opt-in, red-team, OpenClaw, Hermes  -> docs/plans/E15-v10-public-process.md  (30/31)

### E16 -- First live benchmarks  -> docs/plans/E16-first-live-benchmarks.md  (5/18)

### E17 -- Live LangGraph and gateway product runs  -> docs/plans/E17-live-langgraph-gateways.md  (0/1)

### E18 -- Soft model-mode and broader live adapter set  -> docs/plans/E18-soft-and-broader-live.md  (0/1)

Completed earlier epics (audit only; not frontier):
E1-E14 under `docs/plans/E1-*` .. `E14-*` (all tasks done through v0.3 / v1.0 tag).

## 5. Parallel Work

Tracks (E16):
- Track R: runtime disclosure T16.0-T16.4
- Track H: harness baseline T16.5-T16.7 (after T16.2)
- Track S: Sire live T16.8-T16.11 (after T16.3); human T16.10
- Track P: publish + close T16.12-T16.18 (after T16.10)

Sync: T16.1 before T16.2/T16.3; T16.10 before T16.12; T16.15 before T16.17.

### Wave 40: runtime disclosure
- T16.0, T16.1, T16.2, T16.3, T16.4

### Wave 41: FakeAdapter baseline
- T16.5, T16.6, T16.7

### Wave 42: Sire live wiring
- T16.8, T16.9, T16.10 (human), T16.11

### Wave 43: live runs + publish
- T16.12, T16.13, T16.14, T16.15, T16.16

### Wave 44: E16 quality gate
- T16.17, T16.18

## 6. Timeline and Milestones

| ID | Milestone | Depends | Exit criteria |
| ---- | ---- | ---- | ---- |
| M11 | v1.0.0 release | E15 tag | done 2026-09-07 |
| M12 | Runtime disclosure | Wave 40 | T16.4 green |
| M13 | Harness baseline documented | Wave 41 | T16.6 green |
| M14 | First live Unratified Sire row | Wave 43 | T16.15 merged |
| M15 | E16 closed | Wave 44 | T16.17 green |
| M16 | LangGraph provisional | T15.11 | after 2026-10-18 or #8992 ratify |
| M17 | E17 planned | T17.0 | executable fidelity |

## 7. Risk Register

| ID | Risk | Impact | Likelihood | Mitigation |
| ---- | ---- | ---- | ---- | ---- |
| R1 | Author bias (Sire) | Credibility | Med | unofficial + runtime=live; ADR 007 |
| R20 | Stub scores misread as product | Credibility | High | ADR 015 runtime field + README warning |
| R21 | Sire live secrets in git | Security | Med | scratch/1Password only; T16.10 human |
| R22 | Sire observability/ledger gaps | Low GSI | High | honest N/A; publish anyway with labels |
| R23 | T16.10 blocked long | Schedule | Med | FakeAdapter baseline still ships |
| R15 | Provisional confused with ratified | Credibility | Med | three-way badge; wait for #175 window |

## 8. Operating Procedure

Definition of done for a task:
1. `acc:` predicate green.
2. Paired tests for protocol/adapter/CLI changes.
3. gofmt/ruff clean; go test / pytest green on touched packages.
4. PR via rebase; CI green.
5. Live publishes: dashboard check + Pages URL verified.

Rules:
- Worktrees for code changes.
- Never per-framework exploit code.
- Do not publish stub rows as product rankings.
- Soft rates still require >=25 seeds when E18 starts.
- Scenario governance comment windows for SEC/REL/GOV ID changes.

## 9. Progress Log

- 2026-09-18: Wave 40 shipped -- T16.1-T16.4 runtime disclosure (ADR 015 wire + adapters + dashboard).
- 2026-09-18: /plan first live benchmarks -- trim focus to E15 remainder + E16 executable + E17/E18 outline; ADR 015; UC-037..039. Prior E1-E14 left as completed audit files.
- 2026-09-18: T15.11 real outreach langchain-ai/langgraph#8992; window closes 2026-10-18.
- 2026-09-07: v1.0.0 tagged (T15.14); E15 30/31.

## 10. Hand off Notes

- Next /apply frontier: Wave 40 (T16.0 already has ADR file drafted -- confirm and wire).
- Founder: T16.10 Sire dogfood credentials; T15.11 after 2026-10-18.
- Do not expand E17 until T16.17 green.

## 11. Appendix

- Spec: docs/RFC-0001.md
- Design: docs/design.md
- ADR 015: docs/adr/015-stub-vs-live-runtime.md
- Roadmap story: docs/roadmap.md
- Use cases: .claude/scratch/usecases-manifest.json
- Release map: v0.1=E1-E12; v0.2=E13; v0.3=E14; v1.0=E15; live benchmarks=E16+
