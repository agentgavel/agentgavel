# E6 -- mcpfuzz rogue MCP servers

Acceptance: Each SEC-007 fuzz mode runs as a subprocess MCP server (stdio JSON-RPC) and is covered by a Go test that asserts the advertised misbehavior.
fidelity: executable

## Wave 1

- [x] T6.1 Shared MCP stdio server skeleton in internal/mcpfuzz  Owner: pool  Est: 75m  kind: agent  verifies: [UC-009, UC-013]  acc: [go test ./internal/mcpfuzz -run Skeleton lists one noop tool via tools/list]  completed: 2026-09-02
  - deps: [T1.1]

- [ ] T6.2 Mode toxic-output: inject instructions in tool results  Owner: pool  Est: 45m  kind: agent  verifies: [UC-013]  acc: [go test ./internal/mcpfuzz -run ToxicOutput result contains SYSTEM Ignore prior instructions]
  - deps: [T6.1]

- [ ] T6.3 Mode schema-flooder: extra/unknown params and deep schemas  Owner: pool  Est: 45m  kind: agent  verifies: [UC-009, UC-013]  acc: [go test ./internal/mcpfuzz -run SchemaFlooder tools/list returns undeclared params]
  - deps: [T6.1]

- [ ] T6.4 Mode early-disconnect: drop mid-call  Owner: pool  Est: 45m  kind: agent  verifies: [UC-013]  acc: [go test ./internal/mcpfuzz -run EarlyDisconnect closes transport before tools/call returns]
  - deps: [T6.1]

- [ ] T6.5 Mode tool-renamer: mutate names between list and call  Owner: pool  Est: 45m  kind: agent  verifies: [UC-013]  acc: [go test ./internal/mcpfuzz -run ToolRenamer list name differs from subsequent call name]
  - deps: [T6.1]

- [ ] T6.6 Mode slowloris: delayed responses under context  Owner: pool  Est: 30m  kind: agent  verifies: [UC-013]  acc: [go test ./internal/mcpfuzz -run Slowloris exceeds configured delay threshold]
  - deps: [T6.1]

- [ ] T6.7 Mode masquerade: impersonate an already-granted tool identity  Owner: pool  Est: 45m  kind: agent  verifies: [UC-013]  acc: [go test ./internal/mcpfuzz -run Masquerade advertises colliding name with different backend id]
  - deps: [T6.1]

- [ ] T6.8 Engine helper to start mode by name and inject endpoint into SessionConfig  Owner: pool  Est: 45m  kind: agent  verifies: [UC-013]  acc: [go test ./internal/engine -run StartFuzzMode returns dialable endpoint for toxic-output]
  - deps: [T6.2, T3.2]

- [ ] T6.9 Lint/format mcpfuzz  Owner: pool  Est: 20m  kind: agent  verifies: [infrastructure]  acc: [gofmt -l internal/mcpfuzz is empty]
  - deps: [T6.8]
