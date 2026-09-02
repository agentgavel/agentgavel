# ADR 001: Go Engine with Sidecar Adapters

## Status
Accepted

## Date
2026-09-02

## Context
AgentGavel must produce trustworthy, reproducible scores. Target frameworks
are mostly Python (LangGraph, AutoGen, CrewAI) while the author also builds
Sire. An in-process Python harness would couple the engine to every target's
dependency tree and weaken auditability of the scoring core.

## Decision
Implement the harness core in Go as a static binary. Communicate with targets
only through sidecar processes that speak a versioned wire protocol. Ship a
Python SDK so Python targets implement callbacks, not transport. The engine
never imports a target framework.

## Consequences
Positive: auditable scoring core; one CI binary; language-neutral adapter
ratification; concurrent fuzz and seed scheduling with goroutines.
Negative: protocol and SDK work precedes the first scenario result; adapters
need process lifecycle management.
