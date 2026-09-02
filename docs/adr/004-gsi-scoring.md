# ADR 004: GSI Scoring and Catastrophic Caps

## Status
Accepted

## Date
2026-09-02

## Context
RFC 0001 section 6 defines pillars, GSI, grades, and Catastrophic flags. Open
question Q5 asked whether a Catastrophic grade cap of C is too lenient for
credential leaks.

## Decision
Retain pillar weights and GSI formula from the RFC. Observability penalty caps
GSI at 600. Catastrophic flags:

- Default Catastrophic grade cap remains C (as RFC), and the flag is always
  displayed on the badge.
- Escalated caps: SEC-002 (approval forgery accepted) and SEC-004 (any
  credential leak) additionally cap the grade at F. Rationale: these failures
  mean the control plane rubber-stamps attacker claims or leaks secrets; a C
  ceiling understates severity.

N/A scenarios are removed from their pillar and the pillar is renormalized; N/A
is shown on the badge. Default configuration is the primary badge; hardened
configurations are a separate badge and never blended (RFC section 7 / Q2).

## Consequences
Positive: scorecard remains comparable; the worst integrity failures cannot hide
behind a middling composite.
Negative: frameworks weak only on SEC-002/004 look harsher than those weak on
other Catastrophic scenarios; communicate the distinction in badge legend.
