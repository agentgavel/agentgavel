# E12 -- v0.1 quality gate and release

Acceptance: CI green on Go+Python; tagged v0.1.0 release artifacts include AgentGavel binary checksums; RFC status note updated to Implemented for v0.1 scope; smoke run documented in docs/manual/v0.1-smoke.md.
fidelity: executable

## Wave 1

- [x] T12.1 Expand CI matrix: Go test/race, Python pytest SDK, ruff, fixture path checks  Owner: pool  Est: 75m  kind: agent  verifies: [infrastructure]  acc: [ci.yml jobs all required and pass on a PR branch]  completed: 2026-09-03
  - deps: [T1.4, T7.5, T8.10]

- [x] T12.2 Add goreleaser or GitHub release workflow for static AgentGavel binaries  Owner: pool  Est: 60m  kind: agent  verifies: [infrastructure]  acc: [tag workflow builds linux/darwin amd64 arm64 artifacts]  completed: 2026-09-03
  - deps: [T9.2]

- [x] T12.3 VERSION file / ldflags version injection verified in binary  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  acc: [AgentGavel version prints release version when ldflags set]  completed: 2026-09-03
  - deps: [T1.1, T12.2]

- [x] T12.4 Write docs/manual/v0.1-smoke.md with exact commands for Fake, Sire, LangGraph  Owner: pool  Est: 45m  kind: agent  delivers: [docs/manual/v0.1-smoke.md]  verifies: [UC-001, UC-016, UC-017]  acc: [smoke doc lists commands and expected provenance labels]  completed: 2026-09-03
  - deps: [T10.5, T11.5, T9.3]

- [x] T12.5 Append docs/devlog.md entry for v0.1 plan execution start/finish template  Owner: pool  Est: 20m  kind: agent  delivers: [devlog seed entry]  verifies: [infrastructure]  acc: [devlog has dated stub section for v0.1]  completed: 2026-09-03
  - deps: [T1.5]

- [ ] T12.6 Update RFC-0001 status line to note v0.1 implementation in progress/complete when gated  Owner: pool  Est: 20m  kind: human  delivers: [RFC status bump]  verifies: [infrastructure]  acc: [RFC Status field reflects v0.1 shipping state]
  - deps: [T12.4]
  - Note: human confirms public messaging.

- [x] T12.7 Full make test && make lint on clean tree; fix stragglers  Owner: pool  Est: 60m  kind: agent  verifies: [infrastructure]  acc: [make test and make lint exit 0]  completed: 2026-09-03
  - deps: [T12.1, T8.11, T9.6, T10.7, T11.7]

- [ ] T12.8 Cut v0.1.0 tag after main green (founder)  Owner: pool  Est: 30m  kind: human  verifies: [infrastructure]  acc: [git tag v0.1.0 exists on origin and release assets uploaded]
  - deps: [T12.7, T12.2]
