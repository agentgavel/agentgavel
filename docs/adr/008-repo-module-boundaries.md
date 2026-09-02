# ADR 008: Repository Module and Release Boundaries

## Status
Accepted

## Date
2026-09-02

## Context
Greenfield repo `github.com/agentgavel/gavel` must match RFC layout while
staying CI-friendly. Release plan spans v0.1 through v1.0.

## Decision
- Module path: `github.com/agentgavel/gavel` (Go module at repo root).
- Binary name: `AgentGavel` (CLI under `cmd/AgentGavel`).
- Adapters are separate packages under `adapters/<name>/` with their own
  dependency manifests (e.g. `pyproject.toml`) so Python deps never enter the
  Go module.
- Suites version as `SEC-vN` artifacts under `suites/security/`.
- v0.1 ships SEC-001..007, protocol, Python SDK, Oracle, unofficial Sire and
  LangGraph adapters. Later suites and adapters land in subsequent tagged
  releases per RFC section 8.
- Go standard library preferred for CLI (`flag`), testing, and HTTP where
  practical; protobuf tooling only as needed for schema codegen.

## Consequences
Positive: clear ownership boundaries; adapters version independently in the
fingerprint.
Negative: multi-language CI matrix from day one.
