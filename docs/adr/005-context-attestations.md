# ADR 005: Hosted Context Privacy Attestations

## Status
Accepted

## Date
2026-09-02

## Context
SEC-004 and SEC-008 need to inspect LLM context for leaks. Open question Q3
raised privacy concerns for hosted runs if raw context leaves the customer
environment.

## Decision
Default execution is local-only and may use raw `context_snapshot` events when
the operator opts in. For hosted or privacy-sensitive runs, adapters emit hashed
n-gram attestations of context instead of raw text. Leak checks for SEC-004 and
SEC-008 run against those attestations (exact credential encodings hashed the
same way; semantic canary paraphrase judge remains local-only or N/A when raw
context is unavailable). Sensitivity loss is accepted for hosted mode and
documented on the scorecard as `context_mode=attestation`.

## Consequences
Positive: hosted participation without shipping prompts off-box.
Negative: attestation mode may miss some encodings or paraphrase leaks;
scorecards must label the mode so results are not silently compared.
