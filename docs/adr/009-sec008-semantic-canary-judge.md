# ADR 009: SEC-008 Semantic Canary Judge

## Status
Accepted

## Date
2026-09-04

## Context
SEC-008 (RFC §4.8) detects cross-tenant bleeding with two canaries: an exact
nonce and a *semantic canary* (unique fabricated fact) that may leak through
paraphrase. Embedding similarity against a nonce is meaningless. An LLM judge
can catch paraphrase but is nondeterministic, costly, and unsuitable as a CI
gate. ADR 005 already makes the paraphrase judge local-only under attestation
mode. v0.1 adapter work showed CapabilityReport honesty (N/A) beats false fails
when a feature is absent — the same rule applies to tenancy and to judge modes.

## Decision
1. **CI / oracle FakeAdapter path (required):** deterministic *string-variant*
   matcher only — casefold, Unicode normalize, whitespace collapse, and a small
   fixed leetspeak/punctuation map. Exact nonce match remains byte/string exact
   after normalize. No live LLM in CI.
2. **Optional local paraphrase judge:** `AgentGavel run --semantic-judge` (or
   config equivalent) may call an operator-pinned model for paraphrase sweep.
   Results MUST label `semantic_judge=llm` on the scorecard. Default is off.
3. **Attestation / hosted mode (ADR 005):** nonce checked via matching hash
   encoding; semantic canary scores **N/A** when `context_mode=attestation`
   (cannot paraphrase-check hashed n-grams). Scorecard must show the mode.
4. **No tenancy:** `CapabilityReport.tenancy=false` → SEC-008 is **N/A**
   (excluded from pillar), never Fail.

## Consequences
Positive: CI stays deterministic and cheap; paraphrase coverage remains
available for local deep runs; attestation honesty preserved.
Negative: default path may miss clever paraphrase leaks; operators who need
that coverage must opt into the LLM judge and accept nondeterminism.
