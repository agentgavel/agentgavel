# E4 -- Compliance Oracle

Acceptance: Oracle HTTP server accepts OpenAI-shaped and Anthropic-shaped chat completions and always returns tool/function calls matching the probe directive supplied via request metadata or a side channel documented in package docs.
fidelity: executable

## Wave 1

- [x] T4.1 Design Oracle probe directive binding (header or system channel) and document it  Owner: pool  Est: 45m  kind: agent  verifies: [UC-002]  lane: agent  delivers: [internal/oracle README or doc comment]  acc: [documented binding lets a client force a specific tool name+args in the next completion]  completed: 2026-09-02
  - deps: [T1.1]
  - Decision rationale: docs/adr/003-hard-soft-oracle.md

- [x] T4.2 Implement OpenAI-compatible /v1/chat/completions handler that emits the directed tool call  Owner: pool  Est: 90m  kind: agent  verifies: [UC-002]  acc: [go test ./internal/oracle -run OpenAIToolCall returns choices with tool_calls matching directive]  completed: 2026-09-03
  - deps: [T4.1]

- [x] T4.3 Implement Anthropic-shaped messages handler with equivalent tool-use content blocks  Owner: pool  Est: 75m  kind: agent  verifies: [UC-002]  acc: [go test ./internal/oracle -run AnthropicToolUse returns content tool_use matching directive]  completed: 2026-09-03
  - deps: [T4.1]

- [x] T4.4 Add oracle subcommand or engine-managed lifecycle: AgentGavel oracle --listen  Owner: pool  Est: 45m  kind: agent  verifies: [UC-002]  acc: [AgentGavel oracle --listen 127.0.0.1:0 prints addr; curl health returns 200]  completed: 2026-09-03
  - deps: [T4.2, T1.1]

- [x] T4.5 Negative tests: without directive, Oracle refuses to invent actions (explicit error)  Owner: pool  Est: 30m  kind: agent  verifies: [UC-002]  acc: [go test ./internal/oracle -run MissingDirective returns 4xx]  completed: 2026-09-03
  - deps: [T4.2]

- [x] T4.6 Lint/format oracle package  Owner: pool  Est: 20m  kind: agent  verifies: [infrastructure]  acc: [gofmt -l internal/oracle is empty]  completed: 2026-09-03
  - deps: [T4.5]
