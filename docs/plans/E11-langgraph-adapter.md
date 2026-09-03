# E11 -- LangGraph adapter (unofficial)

Acceptance: adapters/langgraph sidecar runs a minimal LangGraph graph under AgentGavel, labeled unofficial, and produces scores for scenarios it can observe.
fidelity: executable

## Wave 1

- [x] T11.1 Scaffold adapters/langgraph with pyproject and Adapter subclass  Owner: pool  Est: 60m  kind: agent  verifies: [UC-017, UC-006]  acc: [module entry starts stdio server]  completed: 2026-09-03
  - deps: [T7.3]

- [ ] T11.2 Minimal graph with tool nodes for read_email/send_email fixtures  Owner: pool  Est: 90m  kind: agent  verifies: [UC-017, UC-007]  lane: agent  acc: [pytest runs graph once with Oracle base_url and records a tool call event]
  - deps: [T11.1, T4.2]

- [ ] T11.3 HITL interrupt mapping to ResolveApproval when supported; else capability hitl=false  Owner: pool  Est: 75m  kind: agent  verifies: [UC-005, UC-017]  acc: [CapabilityReport.hitl reflects real interrupt support; test covers both paths]
  - deps: [T11.2]

- [ ] T11.4 Event hooks: tool_invocation before/after, gate_decision, context attestation helper  Owner: pool  Est: 60m  kind: agent  verifies: [UC-004, UC-010]  acc: [integration emits before and after tool_invocation ordered correctly]
  - deps: [T11.2]

- [ ] T11.5 End-to-end AgentGavel run against LangGraph for SEC-001 and SEC-007 toxic-output  Owner: pool  Est: 60m  kind: agent  verifies: [UC-017, UC-007, UC-013]  acc: [scorecard provenance=unofficial; SEC-001 and SEC-007 rows present]
  - deps: [T11.3, T11.4, T9.2]

- [ ] T11.6 adapters/langgraph/README unofficial labeling  Owner: pool  Est: 30m  kind: agent  delivers: [adapters/langgraph/README]  verifies: [UC-017]  acc: [README states unofficial]
  - deps: [T11.1]

- [ ] T11.7 Ruff lint adapters/langgraph  Owner: pool  Est: 20m  kind: agent  verifies: [infrastructure]  acc: [ruff check adapters/langgraph exits 0]
  - deps: [T11.5]
