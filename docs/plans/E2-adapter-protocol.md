# E2 -- Adapter wire protocol

Acceptance: Handshaking JSON-RPC stdio client/server round-trip in Go tests; proto schema documents all RFC 5.3 methods; version negotiation rejects incompatible majors.
fidelity: executable

## Wave 1

- [x] T2.1 Author proto/adapter.proto with Handshake, StartSession, SubmitTask, ResolveApproval, Events, ExportLedger, StopSession and message types  Owner: pool  Est: 90m  kind: agent  verifies: [UC-003, UC-004, UC-005]  lane: agent  acc: [proto/adapter.proto defines all seven RPCs and Event oneof kinds from RFC 5.3]  completed: 2026-09-03
  - deps: [T1.1]
  - Decision rationale: docs/adr/002-adapter-transport.md

- [x] T2.2 Implement internal/protocol types + JSON codec (stdlib encoding/json) mirroring proto  Owner: pool  Est: 90m  kind: agent  verifies: [UC-003, UC-004]  acc: [go test ./internal/protocol -run Codec round-trips HandshakeRequest and Event tool_invocation]  completed: 2026-09-03
  - deps: [T2.1]

- [x] T2.3 Implement stdio JSON-RPC 2.0 framing (newline-delimited request/response + event notifications)  Owner: pool  Est: 90m  kind: agent  verifies: [UC-003, UC-004]  acc: [go test ./internal/protocol -run Stdio launches two ends and completes Handshake]  completed: 2026-09-03
  - deps: [T2.2]

- [ ] T2.4 Implement session lifecycle helpers: start, submit, resolve, export, stop with timeouts/context  Owner: pool  Est: 60m  kind: agent  verifies: [UC-004, UC-005]  acc: [go test ./internal/protocol -run SessionLifecycle passes with a fake adapter server]
  - deps: [T2.3]

- [ ] T2.5 Version negotiation: Handshake carries engine and adapter protocol versions; incompatible major returns error  Owner: pool  Est: 45m  kind: agent  verifies: [UC-003]  acc: [go test ./internal/protocol -run VersionReject fails Handshake when major differs]
  - deps: [T2.3]

- [x] T2.6 Document CapabilityReport fields (hitl, tenancy, ledger, observability, context_mode) and N/A mapping  Owner: pool  Est: 45m  kind: agent  verifies: [UC-003]  delivers: [docs snippet or package comment for CapabilityReport]  acc: [CapabilityReport struct fields documented and covered by a table-driven test of N/A mapping helpers]  completed: 2026-09-03
  - deps: [T2.2]

- [ ] T2.7 Protocol package tests for Events before/after tool_invocation ordering invariant  Owner: pool  Est: 45m  kind: agent  verifies: [UC-004]  acc: [go test ./internal/protocol -run ToolInvocationOrder fails if after precedes before]
  - deps: [T2.4]

- [ ] T2.8 Lint/format protocol package  Owner: pool  Est: 20m  kind: agent  verifies: [infrastructure]  acc: [gofmt -l internal/protocol proto is empty]
  - deps: [T2.7]
