# E9 -- CLI run, report, fingerprint

Acceptance: `AgentGavel run` and `AgentGavel report` produce a human scorecard and machine JSON including GSI, flags, provenance, and full fingerprint fields from RFC 4.11.
fidelity: executable

## Wave 1

- [x] T9.1 Implement fingerprint struct and hasher (scenario-version, framework-version, config-hash, adapter-version, model, seed-set)  Owner: pool  Est: 45m  kind: agent  verifies: [UC-015]  acc: [go test ./internal/engine -run Fingerprint stable hash for identical inputs]  completed: 2026-09-03
  - deps: [T3.4]

- [x] T9.2 CLI run: flags for adapter command, suite, scenarios, modes, seeds, out dir  Owner: pool  Est: 90m  kind: agent  verifies: [UC-001, UC-002]  acc: [AgentGavel run --adapter <fake> --suite security --seeds 25 --mode oracle writes results and exit 0 on all-pass fixture]  completed: 2026-09-03
  - deps: [T3.6, T8.10, T4.4]

- [x] T9.3 CLI report: render scorecard text + JSON from results dir  Owner: pool  Est: 60m  kind: agent  verifies: [UC-014]  acc: [AgentGavel report <run-id> prints GSI grade pillars and Catastrophic flags]  completed: 2026-09-03
  - deps: [T5.6, T9.1]

- [x] T9.4 Reproduce path: run --fingerprint file reloads pins  Owner: pool  Est: 45m  kind: agent  verifies: [UC-015]  acc: [AgentGavel run --fingerprint results/x/fingerprint.json uses same seed-set]  completed: 2026-09-03
  - deps: [T9.1, T9.2]

- [x] T9.5 CLI tests (exec of binary) for help, version, run dry-run validation errors  Owner: pool  Est: 45m  kind: agent  verifies: [UC-001]  acc: [go test ./cmd/AgentGavel -run CLI covers missing adapter error]  completed: 2026-09-03
  - deps: [T9.2]

- [ ] T9.6 Lint/format cmd  Owner: pool  Est: 20m  kind: agent  verifies: [infrastructure]  acc: [gofmt -l cmd is empty]
  - deps: [T9.5]
