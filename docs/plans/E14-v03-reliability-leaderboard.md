# E14 -- v0.3 reliability, rubber-stamp, leaderboard

Acceptance: Reliability suite ships (REL-001..003 against FakeAdapter);
`AgentGavel rubber-stamp` runs SEC-002 and SEC-006 only with CI-compatible
exit codes; static `dashboard/` on GitHub Pages shows Opt-in and Unratified
tabs per ADR 006 with sample scorecards and honest provenance labels.
fidelity: executable

## Learnings from v0.2 (bind into tasks)

- CapabilityReport honesty: missing hitl / ledger → N/A, never silent Fail
  (Sire/LangGraph/§8.1 adapters). REL-001/003 inherit the same rule (ADR 010).
- `--ci` exit-code map is the contract for non-interactive gates: Fail→1,
  Catastrophic→2. `rubber-stamp` must reuse that map, not invent a third.
- Scenario governance comment window (`docs/manual/scenario-governance.md`)
  applies to REL catalog changes before published leaderboard scores move.
- FakeAdapter oracle E2E first; do not block REL predicates on adapter ports.
- Soft rates stay ≥25 seeds with Wilson intervals when a REL path is Soft.
- Provenance on dashboard samples stays `unofficial` until ADR 007
  ratification; Unratified tab never promotes into Opt-in (ADR 006).
- Prefer static HTML/JSON under `dashboard/` — no app server in v0.3.
- REL definitions are pinned in ADR 010 (demotion latch, replay/drift,
  ledger completeness). Do not invent alternate REL IDs mid-wave.

## Wave 1 -- catalog + fixtures

- [x] T14.0 PLAN: expand E14 to executable fidelity (informed by v0.2)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E14-v03-reliability-leaderboard.md at fidelity: executable]  deps: [T13.17]  acc: [parse_plan sees E14 with >= 5 tasks, fidelity executable]  completed: 2026-09-04

- [x] T14.1 Scaffold suites/reliability with REL-v0 loader listing REL-001..003  Owner: pool  Est: 45m  kind: agent  verifies: [UC-031]  acc: [go test ./suites/reliability -run Load lists REL-001 through REL-003 and reports version REL-v0]  deps: [T14.0]  completed: 2026-09-04

- [x] T14.2 Fixtures for REL-001..003: demotion probe, drifted replay plan, ledger-gap corpus  Owner: pool  Est: 60m  kind: agent  verifies: [UC-031]  acc: [fixtures/ files exist for demotion, replay-drift, and ledger-gap; no framework-specific exploit code]  deps: [T14.0]  completed: 2026-09-04

## Wave 2 -- REL predicates + E2E

- [x] T14.3 Implement REL-001 demotion latch predicate (N/A without hitl)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-031]  acc: [go test ./suites/reliability -run REL001 covers re-escalate=Fail, clean demotion=100, hitl=false→N/A; ADR 010]  deps: [T14.1, T14.2, T5.3]  completed: 2026-09-04

- [x] T14.4 Implement REL-002 replay/config-drift predicate  Owner: pool  Est: 90m  kind: agent  verifies: [UC-031]  acc: [go test ./suites/reliability -run REL002 detects drifted replay accept as Fail; clean reject/flag=100; no binding→N/A]  deps: [T14.1, T14.2, T5.3]  completed: 2026-09-04

- [x] T14.5 Implement REL-003 ledger completeness predicate  Owner: pool  Est: 90m  kind: agent  verifies: [UC-031]  acc: [go test ./suites/reliability -run REL003 scores gaps Fail; complete ledger=100; ledger=false→N/A]  deps: [T14.1, T14.2, T5.2]  completed: 2026-09-04

- [ ] T14.6 End-to-end FakeAdapter oracle run for REL-001..003  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001, UC-031]  acc: [go test ./suites/reliability -run SuiteOracleFakeREL -count=1 exits 0 with REL-001..003 rows]  deps: [T14.3, T14.4, T14.5, T8.10]

## Wave 3 -- rubber-stamp CLI

- [ ] T14.7 Add AgentGavel rubber-stamp: SEC-002 and SEC-006 only, CI exit-code map  Owner: pool  Est: 90m  kind: agent  verifies: [UC-021]  acc: [go test ./cmd/AgentGavel -run RubberStamp asserts only SEC-002/SEC-006 rows and Fail→exit 1 Catastrophic→exit 2]  deps: [T13.8, T8.4, T8.8]

- [ ] T14.8 Document rubber-stamp usage and exit codes in docs/manual/rubber-stamp.md  Owner: pool  Est: 30m  kind: agent  delivers: [docs/manual/rubber-stamp.md]  verifies: [UC-021]  lane: agent  acc: [doc names SEC-002 and SEC-006 only and documents exit 0/1/2]  deps: [T14.7]

## Wave 4 -- dashboard / leaderboard

- [ ] T14.9 Scaffold dashboard/ static site with Opt-in and Unratified tabs (ADR 006)  Owner: pool  Est: 75m  kind: agent  verifies: [UC-022]  lane: agent  acc: [dashboard/index.html (or equivalent) exposes Opt-in and Unratified sections distinguishable without JS]  deps: [T14.0]

- [ ] T14.10 Sample scorecards in dashboard/data with unofficial provenance labels  Owner: pool  Est: 45m  kind: agent  verifies: [UC-022]  lane: agent  acc: [at least one opt-in sample and one unratified sample JSON; provenance field present and unofficial unless marked sample-ratified]  deps: [T14.9]

- [ ] T14.11 AgentGavel report --publish writes scorecard JSON consumable by dashboard/data  Owner: pool  Est: 75m  kind: agent  verifies: [UC-022]  acc: [go test ./cmd/AgentGavel -run ReportPublish writes a JSON file matching dashboard sample schema]  deps: [T9.3, T14.9]

- [ ] T14.12 Document GitHub Pages publish path for dashboard/  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/leaderboard-pages.md]  verifies: [UC-022]  lane: agent  acc: [doc states Pages source path, Opt-in vs Unratified rules citing ADR 006, and that unpublished drafts do not appear]  deps: [T14.10, T13.9]

## Wave 5 -- smoke + quality + ship

- [ ] T14.13 Add docs/manual/v0.3-smoke.md for REL suite, rubber-stamp, and dashboard  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/v0.3-smoke.md]  verifies: [UC-021, UC-022, UC-031]  lane: agent  acc: [smoke doc lists commands for REL FakeAdapter, rubber-stamp, and local dashboard open]  deps: [T14.6, T14.7, T14.10]

- [ ] T14.14 gofmt clean on suites/reliability and cmd/AgentGavel rubber-stamp/publish paths  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  acc: [gofmt -l suites/reliability cmd/AgentGavel is empty]  deps: [T14.6, T14.7, T14.11]

- [ ] T14.15 Full make test && make lint on clean tree after E14 landings  Owner: pool  Est: 45m  kind: agent  verifies: [infrastructure]  acc: [make test and make lint exit 0]  deps: [T14.14, T14.8, T14.12, T14.13]

- [ ] T14.16 Cut v0.3.0 tag after main green (founder)  Owner: pool  Est: 30m  kind: human  verifies: [infrastructure]  acc: [git tag v0.3.0 exists on origin and release assets uploaded]  deps: [T14.15]
  - Note: human confirms public v0.3 messaging and tag push.
