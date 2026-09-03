# E10 -- Sire adapter (unofficial)

Acceptance: adapters/sire sidecar starts via Python SDK, reports CapabilityReport, completes Handshake with engine, and runs SEC-001..007 with provenance=unofficial on the scorecard. Missing capabilities are honest N/A.
fidelity: executable

## Wave 1

- [x] T10.1 Scaffold adapters/sire with pyproject and Adapter subclass  Owner: pool  Est: 60m  kind: agent  verifies: [UC-016, UC-006]  acc: [python -m adapters.sire --help or module entry starts stdio server]  completed: 2026-09-03
  - deps: [T7.3]
  - Decision rationale: docs/adr/007-adapter-ratification.md

- [x] T10.2 Map Sire session lifecycle to AgentGavel callbacks (best-effort against public Sire APIs or documented stubs)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-016, UC-004]  lane: agent  acc: [unit tests with mocked Sire client cover start/submit/stop]  completed: 2026-09-03
  - deps: [T10.1]
  - Risk: Sire API surface may need stubs if local Sire is unavailable; record gaps as N/A capabilities.

- [x] T10.3 Wire ResolveApproval and event emission from Sire hooks  Owner: pool  Est: 75m  kind: agent  verifies: [UC-005, UC-016]  acc: [integration test emits gate_decision on ResolveApproval]  completed: 2026-09-03
  - deps: [T10.2]

- [ ] T10.4 ExportLedger mapping or capability hitl/ledger false with observability penalty path  Owner: pool  Est: 45m  kind: agent  verifies: [UC-016, UC-014]  acc: [CapabilityReport fields match what ExportLedger can actually provide]
  - deps: [T10.2]

- [ ] T10.5 End-to-end AgentGavel run --adapter sire --mode oracle on SEC-002 at minimum  Owner: pool  Est: 60m  kind: agent  verifies: [UC-016, UC-008]  acc: [scorecard JSON shows provenance=unofficial and a numeric SEC-002 score or N/A]
  - deps: [T10.3, T10.4, T9.2]

- [ ] T10.6 Document unofficial status and ratification path in adapters/sire/README  Owner: pool  Est: 30m  kind: agent  delivers: [adapters/sire/README]  verifies: [UC-016]  acc: [README states unofficial and links docs/adr/007-adapter-ratification.md]
  - deps: [T10.1]

- [ ] T10.7 Ruff lint adapters/sire  Owner: pool  Est: 20m  kind: agent  verifies: [infrastructure]  acc: [ruff check adapters/sire exits 0]
  - deps: [T10.5]
