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

## Learnings from Wave 21–22 (bind into the remaining tasks)

- `suites/reliability.RunOracleFakeREL` exists but the CLI still rejects
  `--suite reliability` (`cmd/AgentGavel/run.go` accepts only `security`),
  so UC-031's CLI interface is unwired. T14.17 closes that gap before smoke.
- `internal/metrics.ScenarioPillar` has no REL-* entries: `ComputeGSI`
  silently drops unmapped IDs, so a REL run reports a GSI that ignores REL.
  `docs/design.md` already places REL-001..003 in the resilience pillar for
  v0.3; T14.18 wires it. Keep the RFC §6 weights unchanged.
- `protocol.ScenarioNA` drives N/A for SEC/GOV only. REL-001 needs
  `hitl=false`, REL-002/REL-003 need `ledger=false` (same driver as
  SEC-009/010; a framework without receipts has no config-hash binding).
- A gate that cannot observe an approval store must not print green:
  `rubber-stamp` with both scenarios N/A exits 1 with a `not_applicable`
  reason (ADR 011). Do not extend the exit map with a new code.
- A static site cannot list a directory: `dashboard/data/index.json` is the
  entry list, and `report --publish` is the only writer of both the entry
  file and the index (one code path, no hand-edited index).
- No GitHub Pages site exists yet (`gh api repos/agentgavel/agentgavel/pages`
  → 404). Enabling it is a founder-visible, outward-facing action (T14.20),
  and `v0.3.0` (T14.16) waits on the live URL, not only on green CI.
- Until v1.0 signatures, `report --publish` writes only the Unratified tab;
  Opt-in entries are hand-labeled `sample: true` (ADR 006 addendum).
- Test names in `cmd/AgentGavel` follow `Test<Command><Behavior>`
  (`TestCIModeExitMapper`, `TestReportJSON`); new tests keep that shape so
  `-run RubberStamp` / `-run ReportPublish` select them.

## File-touch map (remaining tasks)

| Task | Creates / modifies |
| ---- | ---- |
| T14.7 | `cmd/AgentGavel/rubber_stamp.go`, `rubber_stamp_test.go`, `main.go` (dispatch + help), `suites/security/oracle_fake.go` (`OracleFakeResult.NA`) |
| T14.17 | `cmd/AgentGavel/run.go`, `run_reliability_test.go`, `internal/protocol/capability.go` (+test), `suites/reliability/oracle_fake.go` (`Capabilities`) |
| T14.18 | `internal/metrics/gsi.go` (+test), `internal/report/report_test.go` |
| T14.9 | `dashboard/index.html`, `dashboard/style.css`, `dashboard/app.js`, `dashboard/data/schema.json`, `dashboard/data/index.json`, `dashboard/README.md` |
| T14.10 | `dashboard/data/sample-*.json`, `dashboard/data/index.json` |
| T14.11 | `internal/publish/{entry.go,entry_test.go}`, `cmd/AgentGavel/report.go`, `report_publish_test.go` |
| T14.19 | `.github/workflows/pages.yml` |
| T14.21 | `scripts/check-dashboard.sh`, `.github/workflows/ci.yml` (job `dashboard`) |
| T14.8 / T14.12 / T14.13 / T14.22 / T14.23 | `docs/manual/rubber-stamp.md`, `docs/manual/leaderboard-pages.md`, `docs/manual/v0.3-smoke.md`, `README.md`, `docs/manual/scenario-governance.md` |

## Wave 1 -- catalog + fixtures

- [x] T14.0 PLAN: expand E14 to executable fidelity (informed by v0.2)  Owner: pool  Est: 1h  kind: plan  delivers: [plans/E14-v03-reliability-leaderboard.md at fidelity: executable]  deps: [T13.17]  acc: [parse_plan sees E14 with >= 5 tasks, fidelity executable]  completed: 2026-09-04

- [x] T14.1 Scaffold suites/reliability with REL-v0 loader listing REL-001..003  Owner: pool  Est: 45m  kind: agent  verifies: [UC-031]  acc: [go test ./suites/reliability -run Load lists REL-001 through REL-003 and reports version REL-v0]  deps: [T14.0]  completed: 2026-09-04

- [x] T14.2 Fixtures for REL-001..003: demotion probe, drifted replay plan, ledger-gap corpus  Owner: pool  Est: 60m  kind: agent  verifies: [UC-031]  acc: [fixtures/ files exist for demotion, replay-drift, and ledger-gap; no framework-specific exploit code]  deps: [T14.0]  completed: 2026-09-04

## Wave 2 -- REL predicates + E2E

- [x] T14.3 Implement REL-001 demotion latch predicate (N/A without hitl)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-031]  acc: [go test ./suites/reliability -run REL001 covers re-escalate=Fail, clean demotion=100, hitl=false→N/A; ADR 010]  deps: [T14.1, T14.2, T5.3]  completed: 2026-09-04

- [x] T14.4 Implement REL-002 replay/config-drift predicate  Owner: pool  Est: 90m  kind: agent  verifies: [UC-031]  acc: [go test ./suites/reliability -run REL002 detects drifted replay accept as Fail; clean reject/flag=100; no binding→N/A]  deps: [T14.1, T14.2, T5.3]  completed: 2026-09-04

- [x] T14.5 Implement REL-003 ledger completeness predicate  Owner: pool  Est: 90m  kind: agent  verifies: [UC-031]  acc: [go test ./suites/reliability -run REL003 scores gaps Fail; complete ledger=100; ledger=false→N/A]  deps: [T14.1, T14.2, T5.2]  completed: 2026-09-04

- [x] T14.6 End-to-end FakeAdapter oracle run for REL-001..003  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001, UC-031]  acc: [go test ./suites/reliability -run SuiteOracleFakeREL -count=1 exits 0 with REL-001..003 rows]  deps: [T14.3, T14.4, T14.5, T8.10]  completed: 2026-09-04

## Wave 3 -- rubber-stamp CLI + REL CLI wiring (Wave 23, 4 agents)

- [x] T14.7 Add AgentGavel rubber-stamp: SEC-002 and SEC-006 only, CI exit-code map, all-N/A fails closed (ADR 011)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-021]  acc: [go test ./cmd/AgentGavel -run RubberStamp -count=1 passes: FakeAdapter run writes summary.json whose scenario keys are exactly SEC-002 and SEC-006; exit map pass→0 fail→1 catastrophic→2 and both-N/A→1 with stderr containing not_applicable; `AgentGavel help` lists rubber-stamp]  deps: [T13.8, T8.4, T8.8]  completed: 2026-09-04
  - Flags: `--adapter` (required), `--out`/`--root`, `--run-id`, `--seeds` (default 25), `--fingerprint`. No `--suite`, no `--scenarios`.
  - Calls `security.RunOracleFake` with `Scenarios: []string{"SEC-002","SEC-006"}` and `Capabilities` from `pingAdapter`; add `NA []string` to `security.OracleFakeResult` so the CLI can see N/A rows without re-reading summary.json.
  - stdout: absolute `summary.json` path only (same shape as `run --ci`). Never print GSI or a grade — two scenarios do not make a scorecard.
  - Exit: `ciExitCode(!AllPass, Catastrophic)`; if every selected scenario is N/A, exit 1 and write `rubber-stamp: not_applicable (hitl=false): no approval store to check` to stderr.
  - Tests: `TestRubberStampFakeAdapterWritesOnlySEC002SEC006`, `TestRubberStampExitMap` (table over a pure helper, includes the all-N/A case), extend `TestCLI_Help` expectations.

- [x] T14.17 Wire `run --suite reliability` (REL-v0) with capability-driven N/A and `--ci` exit map  Owner: pool  Est: 75m  kind: agent  verifies: [UC-031, UC-020]  acc: [AgentGavel run --suite reliability --adapter <fakeadapter> --mode oracle --ci writes summary.json containing REL-001, REL-002, REL-003 rows and exits 0; go test ./internal/protocol -run ScenarioNA covers hitl=false→REL-001 and ledger=false→REL-002,REL-003]  deps: [T14.6, T13.8]  completed: 2026-09-04
  - `cmd/AgentGavel/run.go`: accept `security|reliability`; the reliability branch calls `reliability.RunOracleFakeREL` with the Handshake `CapabilityReport` and honors `--scenarios`, `--seeds`, `--fingerprint`, `--ci`.
  - `suites/reliability.OracleFakeOptions` gains `Capabilities *protocol.CapabilityReport`; `RunOracleFakeREL` applies `protocol.ScenarioNA` like the security runner.
  - `internal/protocol/capability.go`: `ScenarioNA` adds `REL-001` (hitl=false) and `REL-002`, `REL-003` (ledger=false); `ObservabilityPenalty` unchanged.
  - Tests: `TestRunReliabilityFakeAdapterWritesRELRows` (CLI E2E), `TestRunUnsupportedSuiteStillExits2`, `TestScenarioNAReliability`.

- [x] T14.18 Map REL-001..003 into the Fault Recovery (resilience) pillar for GSI  Owner: pool  Est: 45m  kind: agent  verifies: [UC-031]  acc: [go test ./internal/metrics -run REL -count=1 passes: ComputeGSI with SEC-007=100 and REL-001=0 yields resilience pillar 50 and REL IDs are no longer dropped; go test ./internal/report -run RELRows shows REL N/A rows in Document.NA]  deps: [T14.6, T5.3]  completed: 2026-09-04
  - `internal/metrics/gsi.go`: add `"REL-001"`, `"REL-002"`, `"REL-003"` → `"resilience"` in `ScenarioPillar`. Pillar weights stay per RFC §6 (`docs/design.md` "Fault Recovery 15% (SEC-007 plus REL-001..003 in v0.3)").
  - Tests: `TestComputeGSIRELResiliencePillar`, `TestComputeGSIRELNAExcluded`; `internal/report` test `TestReportRELRows` loading a summary with REL rows.

- [x] T14.23 Extend docs/manual/scenario-governance.md to REL scenarios (ADR 010)  Owner: pool  Est: 30m  kind: agent  delivers: [docs/manual/scenario-governance.md covering REL]  verifies: [UC-031]  lane: agent  acc: [scenario-governance.md names REL scenarios, suites/reliability/, an owning domain for REL, and ADR 010; existing SEC/GOV lifecycle text is unchanged]  deps: [T13.9, T14.1]  completed: 2026-09-04
  - Add REL to the lifecycle intro and the Ownership list; state that REL-v0 (REL-001..003) is published under ADR 010 and that REL-v1 changes need the comment window before leaderboard scores move.

## Wave 4 -- rubber-stamp docs + dashboard scaffold + Pages workflow (Wave 24, 3 agents)

- [x] T14.8 Document rubber-stamp usage and exit codes in docs/manual/rubber-stamp.md  Owner: pool  Est: 30m  kind: agent  delivers: [docs/manual/rubber-stamp.md]  verifies: [UC-021]  lane: agent  acc: [doc names SEC-002 and SEC-006 only, documents exit 0/1/2 and the both-N/A exit 1 rule, cites ADR 011, and includes a copy-paste FakeAdapter command plus a GitHub Actions step]  deps: [T14.7]  completed: 2026-09-05
  - Sections: What it checks (SEC-002 approval forgery, SEC-006 timeout auto-approve), Usage, Exit codes table, CI snippet, What it does not do (no GSI, no grade, not a substitute for `run`).

- [x] T14.9 Scaffold dashboard/ static site with Opt-in and Unratified tabs (ADR 006) and the entry schema  Owner: pool  Est: 75m  kind: agent  verifies: [UC-022]  lane: agent  acc: [dashboard/index.html contains section elements with id="opt-in" and id="unratified" whose visible headings cite ADR 006; dashboard/data/schema.json parses and lists required keys run_id, framework, adapter, adapter_version, provenance, tab, sample, gsi, grade, pillars, catastrophic, na, fingerprint, generated_at; dashboard/data/index.json parses as a JSON array]  deps: [T14.0]  completed: 2026-09-05
  - Vanilla HTML/CSS/JS, no build step, no framework. `app.js` fetches `data/index.json`, then each entry, and renders one table per tab. `<noscript>` text explains the two tabs and links `data/index.json`.
  - Provenance badge is three-way (`ratified` / `provisional` / `unofficial`, ADR 007) and a `sample` entry is visibly tagged "sample".
  - `dashboard/data/schema.json`: JSON Schema draft-07; enums `provenance ∈ {ratified, provisional, unofficial}`, `tab ∈ {opt-in, unratified}`. `index.json` starts as `[]`.
  - `dashboard/README.md`: how to serve locally (`python3 -m http.server -d dashboard 8000`) and that `report --publish` is the only writer of `data/`.

- [x] T14.19 Add .github/workflows/pages.yml deploying dashboard/ to GitHub Pages  Owner: pool  Est: 45m  kind: agent  verifies: [UC-022]  acc: [.github/workflows/pages.yml exists with on.push.branches main + paths dashboard/** + workflow_dispatch, permissions pages:write and id-token:write, actions/upload-pages-artifact path dashboard, and actions/deploy-pages; `gh workflow view pages.yml` resolves after merge]  deps: [T14.9]  completed: 2026-09-05
  - Two jobs (`build` uploads the artifact, `deploy` runs `actions/deploy-pages@v4` with environment `github-pages`), `concurrency: group: pages, cancel-in-progress: false`.
  - The workflow cannot succeed until Pages is enabled (T14.20); do not mark T14.19 done by a green run — the file and its parse are the acceptance.

## Wave 5 -- dashboard content + publish + docs (Wave 25, 4 agents)

- [x] T14.10 Sample scorecards in dashboard/data with honest provenance and sample labels  Owner: pool  Est: 45m  kind: agent  verifies: [UC-022]  lane: agent  acc: [dashboard/data contains at least one entry with tab=opt-in and one with tab=unratified, every entry has sample=true and validates against schema.json required keys and enums, every unratified entry has provenance=unofficial, and index.json lists exactly the entry files]  deps: [T14.9]  completed: 2026-09-05
  - Unratified sample: derived from a real FakeAdapter `run --suite security` summary (provenance `unofficial`, framework `FakeAdapter (sample)`).
  - Opt-in sample: framework `Example Framework (sample)`, provenance `ratified`, `sample: true` — it demonstrates the tab, not a real ratification (ADR 006 addendum).

- [x] T14.11 AgentGavel report --publish writes a dashboard entry and updates index.json (Unratified only in v0.3)  Owner: pool  Est: 75m  kind: agent  verifies: [UC-022]  acc: [go test ./cmd/AgentGavel -run ReportPublish -count=1 passes: report --publish --dashboard <dir> <run-id> writes <dir>/data/<run-id>.json with tab=unratified, sample=false, provenance copied from summary.json, gsi/grade/pillars/na/catastrophic/fingerprint from the report Document, and appends the filename to <dir>/data/index.json exactly once; --tab opt-in exits 2 citing ADR 006]  deps: [T9.3, T14.9]  completed: 2026-09-05
  - New package `internal/publish`: `Entry` struct (mirrors `schema.json`), `FromDocument(report.Document, framework, adapter string) Entry`, `Validate(Entry) error`, `Write(dir string, e Entry) (path string, err error)` which also rewrites `index.json` (sorted, deduplicated).
  - `cmd/AgentGavel/report.go`: flags `--publish`, `--dashboard <dir>` (default `dashboard`), `--framework`, `--adapter-name`; `--publish` implies JSON output of the written path.
  - Tests: `TestReportPublishWritesEntryAndIndex`, `TestReportPublishIdempotentIndex`, `TestReportPublishRejectsOptIn`, `internal/publish` `TestEntryValidateEnums`.

- [x] T14.12 Document GitHub Pages publish path for dashboard/  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/leaderboard-pages.md]  verifies: [UC-022]  lane: agent  acc: [doc states the Pages source is pages.yml deploying dashboard/ from main, the public URL pattern, Opt-in vs Unratified rules citing ADR 006 and the three-way provenance of ADR 007, that report --publish writes only Unratified until v1.0, that samples are labeled, and that unmerged drafts do not appear]  deps: [T14.10, T14.19, T13.9]  completed: 2026-09-05

- [x] T14.22 README: document v0.3 commands and the leaderboard  Owner: pool  Est: 30m  kind: agent  delivers: [README.md v0.3 section]  verifies: [UC-021, UC-022, UC-031]  lane: agent  acc: [README.md mentions rubber-stamp, run --suite reliability, and report --publish, links docs/manual/leaderboard-pages.md and docs/manual/rubber-stamp.md, and repeats the unofficial-provenance caveat]  deps: [T14.7, T14.17, T14.11]  completed: 2026-09-05

## Wave 6 -- dashboard check + smoke + fmt (Wave 26, 3 agents)

- [x] T14.21 Add scripts/check-dashboard.sh and a CI job that validates dashboard/data  Owner: pool  Est: 45m  kind: agent  verifies: [UC-022]  lane: agent  acc: [bash scripts/check-dashboard.sh exits 0 on the committed tree and exits 1 when given a temp entry with tab=opt-in and sample=false or an index.json that omits an entry file; .github/workflows/ci.yml has a dashboard job that runs the script]  deps: [T14.10, T14.11]  completed: 2026-09-05
  - Checks: every `dashboard/data/*.json` except `schema.json`/`index.json` has the required keys and enum values; `tab=opt-in ⇒ sample=true` (v0.3 rule); `index.json` equals the sorted entry file list; `leaderboard/index.html` (or legacy `index.html`) contains both section ids. Bash + `python3 -c` JSON only (no new dependencies), mirroring `scripts/check-fixtures.sh`.

- [x] T14.13 Add docs/manual/v0.3-smoke.md for REL suite, rubber-stamp, publish, and dashboard  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/v0.3-smoke.md]  verifies: [UC-021, UC-022, UC-031]  lane: agent  acc: [smoke doc lists copy-paste commands for run --suite reliability --ci, rubber-stamp, report --publish into a temp dashboard dir, and a local dashboard serve, each with the expected stdout and exit code, plus the automated go test mirrors]  deps: [T14.17, T14.7, T14.11, T14.10]  completed: 2026-09-05
  - Follow the `v0.2-smoke.md` shape: prerequisites, one section per surface, "Automated mirrors" listing `-run` selectors.

- [x] T14.14 gofmt clean on suites/reliability, cmd/AgentGavel, internal/publish, internal/metrics, internal/protocol  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  acc: [gofmt -l suites/reliability cmd/AgentGavel internal/publish internal/metrics internal/protocol prints nothing]  deps: [T14.7, T14.11, T14.17, T14.18]  completed: 2026-09-05

## Wave 7 -- quality gate + live leaderboard + ship (Wave 27, 2 agents + founder)

- [ ] T14.15 Full make test && make lint on clean tree after E14 landings  Owner: pool  Est: 45m  kind: agent  verifies: [infrastructure]  acc: [make test and make lint exit 0 and bash scripts/check-dashboard.sh exits 0 on main]  deps: [T14.14, T14.8, T14.12, T14.13, T14.21, T14.22, T14.23]

- [ ] T14.20 Enable GitHub Pages (workflow source) and verify the leaderboard live (founder)  Owner: pool  Est: 30m  kind: human  verifies: [UC-022]  acc: [curl -sf https://agentgavel.dev/leaderboard/ returns HTML containing id="opt-in" and id="unratified", curl -sf https://agentgavel.dev/ returns the marketing home, and curl -sf https://agentgavel.dev/data/index.json returns a JSON array with the sample entries]  deps: [T14.19, T14.10, T14.21]
  - Founder action (outward-facing): `gh api -X POST repos/agentgavel/agentgavel/pages -f build_type=workflow`, then run `pages.yml` via `workflow_dispatch`, then verify with `curl` and an `agent-browser` screenshot. Record the URL in `docs/manual/leaderboard-pages.md` if it differs from the pattern.

- [ ] T14.16 Cut v0.3.0 tag after main green (founder)  Owner: pool  Est: 30m  kind: human  verifies: [infrastructure]  acc: [git tag v0.3.0 exists on origin and release assets uploaded]  deps: [T14.15, T14.20]
  - Note: human confirms public v0.3 messaging and tag push. Release notes name REL-v0 (ADR 010), `rubber-stamp` (ADR 011), and the live leaderboard URL.
