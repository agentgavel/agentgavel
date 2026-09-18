# ADR 015: Stub vs live adapter disclosure

## Status
Accepted

## Date
2026-09-18

## Context

v1.0 ships many unofficial adapters that are in-process stubs (no real
framework package). Operators can already run suites and publish Unratified
scorecards. Without an explicit runtime class, a stub scorecard can be
misread as a measurement of the named product (LangGraph, CrewAI, ADK, etc.).

FakeAdapter is an intentional lab fixture. Sire and gateway adapters can
talk to live systems. The harness needs one honest label for all three.

## Decision

1. CapabilityReport gains a required string field `runtime` with values:
   - `stub` -- in-process fixture shaped like the framework; does not load
     the real product package or live control plane.
   - `live` -- adapter drives a real framework package, hosted API, or
     gateway named in Handshake.
   - `harness` -- FakeAdapter / engine self-test only; never a product ranking.

2. Fingerprints and published scorecards MUST copy `runtime` (and keep
   `provenance`) so Pages rows show both.

3. Publishing rules:
   - `runtime=harness` may publish to Unratified only with framework name
     FakeAdapter (or equivalent lab label).
   - `runtime=stub` may publish Unratified for CI/smoke evidence but MUST
     keep `provenance=unofficial` and MUST NOT be marketed as a product
     ranking.
   - Product-comparable Unratified or Opt-in claims require `runtime=live`
     plus the usual provenance rules (ADR 007).

4. Existing stub adapters set `runtime=stub` without changing observation
   contracts. Promoting an adapter to `live` is a separate change that
   updates Handshake, README, and tests.

## Consequences

Positive: buyers and maintainers can tell lab stubs from product runs.
Negative: every adapter Handshake and dashboard schema must learn one
field; old published rows without `runtime` are treated as `stub` for
display until republished.
