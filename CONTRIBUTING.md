# Contributing to AgentGavel

Thank you for helping measure agent-framework governance under adversarial
pressure.

## Neutrality

Scenarios are code, not judgment. Predicates live in deterministic validators;
fixtures are framework-agnostic strings and protocol mocks. Do not add
per-framework exploit code.

## Adapters

Implement adapters via the Python SDK callbacks (see `sdk/python/` once
shipped). Adapters are labeled ratified, provisional, or unofficial. See
[docs/adr/007-adapter-ratification.md](docs/adr/007-adapter-ratification.md).

## Development (v0.1)

```bash
export GOWORK=off   # parent go.work may list unrelated modules
go build -o AgentGavel ./cmd/AgentGavel
./AgentGavel version
```

Further suite and adapter instructions will land with the v0.1 release.

## Scenario changes

After v0.1, new or changed scenarios go through the RFC process with a public
comment window before they affect published scores (RFC section 0).
