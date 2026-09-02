# ADR 002: Adapter Transport Default

## Status
Accepted

## Date
2026-09-02

## Context
RFC 0001 section 5.3 named gRPC as the default transport and JSON-RPC 2.0 over
stdio as an alternative. Open question Q7 noted that gRPC adds a dependency to
every sidecar and may be unwelcome for simple Python adapters.

## Decision
JSON-RPC 2.0 over stdio is the default transport for v0.1. The `proto/` schema
remains the single source of truth for message shapes; codegen may emit both
stdio JSON codecs and optional gRPC stubs. gRPC is an optional transport for
long-running hosted adapters in later releases. Handshake negotiates protocol
version and transport capabilities.

## Consequences
Positive: Python SDK can ship with stdlib or minimal deps; easier local CI;
matches "one boring binary talks to a child process" mental model.
Negative: streaming Events over stdio needs framing (length-prefixed or NDJSON
with careful buffering); gRPC users wait for an optional path.
