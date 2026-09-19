# E16 -- First live benchmarks

Acceptance: Operators can run oracle-mode security + reliability suites against
FakeAdapter (harness) and a live Sire worker; published Unratified rows carry
`runtime` + `provenance` per ADR 015; a benchmark ops manual exists.
fidelity: executable

## Wave 40 -- Runtime disclosure (ADR 015)

- [x] T16.0 Document ADR 015 stub/live/harness runtime in design + protocol notes  Owner: pool  Est: 30m  kind: agent  verifies: [UC-037]  lane: agent  delivers: [docs/adr/015-stub-vs-live-runtime.md referenced from design.md]  deps: []  acc: [docs/adr/015-stub-vs-live-runtime.md exists and design.md invariant lists runtime stub|live|harness]  completed: 2026-09-18

- [x] T16.1 Add CapabilityReport.runtime to proto + Go/Python wire types  Owner: pool  Est: 75m  kind: agent  verifies: [UC-037]  deps: [T16.0]  acc: [go test ./internal/protocol and pytest sdk/python cover Handshake decoding runtime=stub|live|harness; missing runtime rejected or defaulted per ADR 015 migration rule]  completed: 2026-09-18

- [x] T16.2 Copy runtime into fingerprint + report/publish scorecard JSON  Owner: pool  Est: 60m  kind: agent  verifies: [UC-037, UC-038]  deps: [T16.1]  acc: [AgentGavel report --json includes runtime; dashboard schema accepts runtime; check-dashboard.sh exits 0 on committed tree]  completed: 2026-09-18

- [x] T16.3 Set runtime on FakeAdapter (harness) and all in-tree adapters (stub or documented live)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-037]  deps: [T16.1]  acc: [FakeAdapter Handshake runtime=harness; sire/langgraph/openclaw/hermes and §8.1 stubs report runtime=stub until live wiring; go test FakeAdapter Handshake asserts harness]  completed: 2026-09-18

- [x] T16.4 Lint/format protocol + adapter Handshake changes  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  deps: [T16.2, T16.3]  acc: [gofmt and ruff clean on touched trees]  completed: 2026-09-18

## Wave 41 -- Harness baseline runs

- [x] T16.5 Add docs/manual/benchmark-ops.md for FakeAdapter full SEC+REL oracle baseline  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001, UC-038]  lane: agent  deps: [T16.2]  acc: [manual has copy-paste make build, FakeAdapter build, run --suite security and reliability --mode oracle, report, and states runtime=harness is not a product ranking]  completed: 2026-09-18

- [x] T16.6 Record FakeAdapter oracle baseline summary under results/ or gitignored path + pointer in benchmark-ops  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001]  deps: [T16.3, T16.5]  acc: [documented commands exit 0; summary.json lists SEC and REL rows; path cited in benchmark-ops.md]  completed: 2026-09-18

- [x] T16.7 Optional: publish FakeAdapter Unratified sample labeled harness (not a framework row)  Owner: pool  Est: 45m  kind: agent  verifies: [UC-038]  deps: [T16.2, T16.6]  acc: [dashboard sample or docs state FakeAdapter harness-only; check-dashboard.sh exits 0]  lane: agent  completed: 2026-09-18

## Wave 42 -- Sire live wiring

- [x] T16.8 Env/bootstrap HttpSireClient (token + worker id + base URL) without changing stub default  Owner: pool  Est: 90m  kind: agent  verifies: [UC-039]  deps: [T16.3]  acc: [documented env vars switch python -m adapters.sire to HttpSireClient; without env, stub remains; pytest covers stub default and live construction with fake requester]  completed: 2026-09-18

- [x] T16.9 When live, set runtime=live and probe framework_version; keep ledger/observability honesty  Owner: pool  Est: 60m  kind: agent  verifies: [UC-037, UC-039]  deps: [T16.8]  acc: [Handshake runtime=live only on HttpSireClient path; ledger=false and observability=false unless code proves otherwise; unit test asserts]  completed: 2026-09-18

- [ ] T16.10 Human: provision Sire worker + API token for Oracle-pointed dogfood (non-prod)  Owner: founder  Est: 60m  kind: human  verifies: [UC-039]  deps: [T16.8]  delivers: [scratch credential pointer + worker id for local runs; never commit secrets]  blocked: needs founder Sire account access

- [x] T16.11 Lint/format Sire adapter live bootstrap  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  deps: [T16.8, T16.9]  acc: [ruff clean on adapters/sire]  completed: 2026-09-18

## Wave 43 -- First live Sire benchmarks + publish

- [ ] T16.12 Run Sire live oracle security suite (SEC catalog) with documented seeds; write summary  Owner: pool  Est: 90m  kind: agent  verifies: [UC-039, UC-001]  deps: [T16.9, T16.10]  acc: [AgentGavel run --adapter sire live --suite security --mode oracle exits documented code; summary has runtime=live provenance=unofficial; logged in benchmark-ops or devlog]  blocked-by: [T16.10]

- [ ] T16.13 Run Sire live oracle reliability suite REL-001..003  Owner: pool  Est: 60m  kind: agent  verifies: [UC-039]  deps: [T16.12]  acc: [REL rows present or honest N/A per CapabilityReport; no silent Fail]

- [ ] T16.14 rubber-stamp against live Sire (SEC-002/006 or fail-closed N/A)  Owner: pool  Est: 45m  kind: agent  verifies: [UC-021, UC-039]  deps: [T16.12]  acc: [rubber-stamp exit 0/1/2 matches ADR 011; documented in benchmark-ops]

- [ ] T16.15 Publish first Unratified Sire live scorecard (unofficial + live) via report --publish PR  Owner: pool  Est: 60m  kind: agent  verifies: [UC-038, UC-039]  deps: [T16.12, T16.2]  acc: [dashboard/data entry has runtime=live provenance=unofficial; Pages index lists it after merge; check-dashboard green]

- [ ] T16.16 Append Tier-3 run notes to docs/devlog.md (scores, N/A reasons, fingerprint hashes)  Owner: pool  Est: 30m  kind: agent  verifies: [UC-039]  lane: agent  deps: [T16.12, T16.13]  delivers: [devlog entry for first live Sire oracle runs]

## Wave 44 -- Close E16 quality gate

- [x] T16.17 make test + make lint + check-dashboard green on E16 branch  Owner: pool  Est: 45m  kind: agent  verifies: [infrastructure]  deps: [T16.4, T16.11, T16.15]  acc: [make test && make lint && bash scripts/check-dashboard.sh all exit 0]  completed: 2026-09-18

- [x] T16.18 Update README Quick Start with benchmark-ops link and stub vs live warning  Owner: pool  Est: 30m  kind: agent  verifies: [UC-037]  lane: agent  deps: [T16.5, T16.15]  acc: [README links docs/manual/benchmark-ops.md and ADR 015]  completed: 2026-09-18
