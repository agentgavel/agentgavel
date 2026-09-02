# ADR 007: Adapter Ratification When Maintainers Decline

## Status
Accepted

## Date
2026-09-02

## Context
Open question Q6: who ratifies an adapter when a framework's maintainers
decline to engage? Section 0 requires provenance labels for credibility.

## Decision
Preferred path: framework maintainers review or contribute the adapter
(ratified). If they decline or are unreachable after a documented outreach
attempt and a 30-day public comment window on the adapter PR:

1. AgentGavel core maintainers may grant provisional ratification after an
   independent review checklist (contract honesty, no oracle special-casing,
   event completeness).
2. The badge shows `provenance=provisional` (distinct from `ratified` and
   `unofficial`).
3. Provisional status expires after 180 days unless renewed or upgraded to
   ratified by maintainers.

Author-affiliated adapters (including Sire) cannot skip to ratified via this
path; they remain unofficial or provisional only via independent external
reviewer sign-off.

## Consequences
Positive: benchmarks are not blocked forever by non-engagement; labels stay
honest.
Negative: provisional may be confused with ratified -- badge UI must make the
three-way distinction obvious.
