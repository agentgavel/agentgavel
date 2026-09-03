# E7 -- Python adapter SDK

Acceptance: `sdk/python` installs editable; a sample adapter completes Handshake and emits one tool_invocation event over stdio against the Go protocol test client.
fidelity: executable

## Wave 1

- [x] T7.1 Scaffold sdk/python package (pyproject.toml, AgentGavel adapter base class)  Owner: pool  Est: 60m  kind: agent  verifies: [UC-006]  acc: [pip install -e sdk/python succeeds in a fresh venv]  completed: 2026-09-03
  - deps: [T2.1]
  - Decision rationale: docs/adr/001-go-engine-sidecar-adapters.md

- [x] T7.2 Implement JSON-RPC stdio transport loop and emit() buffering  Owner: pool  Est: 90m  kind: agent  verifies: [UC-003, UC-004, UC-006]  acc: [pytest sdk/python/tests/test_transport.py::test_handshake passes against a Go or Python peer]  completed: 2026-09-03
  - deps: [T7.1, T2.3]

- [x] T7.3 Implement callback dispatch for start_session, submit_task, resolve_approval, export_ledger, stop_session  Owner: pool  Est: 75m  kind: agent  verifies: [UC-004, UC-005, UC-006]  acc: [pytest covers each callback invoked from a corresponding JSON-RPC method]  completed: 2026-09-03
  - deps: [T7.2]

- [x] T7.4 Ship example FakeAdapter used by Go integration tests  Owner: pool  Est: 45m  kind: agent  verifies: [UC-006]  acc: [python -m agentgavel_adapter.examples.fake serves Handshake to Go test]  completed: 2026-09-03
  - deps: [T7.3]

- [x] T7.5 Cross-language integration test: Go engine + Python FakeAdapter  Owner: pool  Est: 60m  kind: agent  verifies: [UC-001, UC-006]  acc: [go test ./internal/engine -run PythonFakeAdapter passes]  completed: 2026-09-03
  - deps: [T7.4, T3.6]

- [ ] T7.6 Python lint: ruff check + format on sdk/python  Owner: pool  Est: 30m  kind: agent  verifies: [infrastructure]  acc: [ruff check sdk/python exits 0]
  - deps: [T7.5]
