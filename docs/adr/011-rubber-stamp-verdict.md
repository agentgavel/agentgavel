# ADR 011: `rubber-stamp` Verdict Contract

## Status
Accepted

## Date
2026-09-04

## Context
RFC §8 item 3 ships `AgentGavel rubber-stamp` in v0.3 as "a fast standalone
check running only SEC-002 and SEC-006". v0.2 established the `run --ci`
exit-code map (pass→0, fail→1, catastrophic→2, catastrophic wins) as the
contract for non-interactive gates. Two questions were open before T14.7
could be executed without guessing:

1. What does `rubber-stamp` print, given that two scenarios from one pillar
   cannot produce a meaningful GSI or grade?
2. What happens when the adapter's `CapabilityReport` reports `hitl=false`,
   which makes both SEC-002 and SEC-006 honestly N/A? Under `run --ci` an
   all-N/A run exits 0. For a gate whose only purpose is to certify "an
   LLM-context claim cannot satisfy your approval gate", exit 0 with no
   scenario scored is a vacuous green.

## Decision
1. `rubber-stamp` is a thin wrapper over the security oracle runner with the
   scenario filter fixed to `SEC-002` and `SEC-006`. It accepts `--adapter`,
   `--out`/`--root`, `--run-id`, `--seeds`, `--fingerprint`; it does not
   accept `--suite` or `--scenarios`.
2. Output matches `run --ci`: stdout is the absolute path to `summary.json`
   and nothing else. `rubber-stamp` never prints a GSI or a grade. A
   scorecard needs the full suite; `report` remains the way to render one.
3. Exit codes reuse `ciExitCode` unchanged: pass→0, fail→1, catastrophic→2.
   No fourth code is added.
4. When every selected scenario is N/A (in v0.3 this means `hitl=false`),
   `rubber-stamp` exits 1 and writes a single stderr line beginning
   `rubber-stamp: not_applicable` naming the missing capability. The
   `summary.json` still records the N/A rows so the reason is inspectable.
   `run --ci` keeps its existing all-N/A→0 behavior: `run` reports what it
   observed across a suite; `rubber-stamp` certifies one property and must
   fail closed when it cannot observe it.

Alternatives considered: exit 0 with a warning (rejected: a CI gate that
passes with zero evidence is the failure mode the honest-N/A rule exists to
prevent); a new exit code 3 for not-applicable (rejected: breaks the
"one exit map" learning from v0.2 and every CI snippet that treats non-zero
as fail would already do the right thing with 1).

## Consequences
Positive: `rubber-stamp` is executable by a low-reasoning session from the
task line alone; CI users get a fail-closed gate that cannot be satisfied by
declining to expose HITL; the exit map stays a single contract.
Negative: frameworks without HITL cannot pass `rubber-stamp` at all, which
is the intended reading of RFC §3.2 ("fails unconditionally"), but the
manual must say so plainly so it is not mistaken for a harness bug.
