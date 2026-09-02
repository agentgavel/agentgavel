# E3 -- Engine core

Acceptance: Engine can launch a fake adapter, run a noop scenario across N seeds with isolation and timeouts, and write a raw run artifact.
fidelity: executable

## Wave 1

- [ ] T3.1 Scaffold internal/engine orchestrator interfaces: Scenario, Runner, SeedScheduler  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001]  acc: [packages compile; Scenario interface requires Setup Probe Observe Predicate Score methods or equivalent named hooks]
  - deps: [T2.4]

- [ ] T3.2 Implement adapter process launcher (exec sidecar, wire stdio, cancel on context)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-003, UC-004]  acc: [go test ./internal/engine -run LaunchFakeAdapter Handshake succeeds and StopSession cleans up process]
  - deps: [T2.4, T3.1]

- [ ] T3.3 Implement seed scheduler: parallel workers, per-seed timeout, deterministic seed set  Owner: pool  Est: 75m  kind: agent  verifies: [UC-002]  acc: [go test ./internal/engine -run SeedScheduler runs 25 noop seeds with max 4 workers and records all seed IDs]
  - deps: [T3.1]

- [ ] T3.4 Implement run artifact writer: JSON results dir with per-scenario and fingerprint stub  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001, UC-015]  acc: [engine write produces results/<run-id>/summary.json readable by encoding/json]
  - deps: [T3.1]

- [ ] T3.5 Wire oracle vs model mode session config (base_url injection)  Owner: pool  Est: 45m  kind: agent  verifies: [UC-002]  acc: [SessionConfig in launcher sets ModelBaseURL from mode; unit test asserts oracle URL when mode=oracle]
  - deps: [T3.2]
  - Decision rationale: docs/adr/003-hard-soft-oracle.md

- [ ] T3.6 Engine integration test with fake adapter completing SubmitTask and emitting events  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001, UC-004]  acc: [go test ./internal/engine -run IntegrationNoopScenario -count=1 passes]
  - deps: [T3.2, T3.3, T3.4]

- [ ] T3.7 Lint/format engine package  Owner: pool  Est: 20m  kind: agent  verifies: [infrastructure]  acc: [gofmt -l internal/engine is empty]
  - deps: [T3.6]
