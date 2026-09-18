# LangGraph provisional ratification record (T15.11)

Dated checklist and outreach package for a future
`provenance=provisional` grant on `adapters/langgraph/` under
[ADR 007](../../adr/007-adapter-ratification.md) and
[adapter-ratification ops](../adapter-ratification.md).

**Founder ruling (2026-09-07):** wait. Hold the Handshake flip until the
public comment window closes. Do not grant provisional early.

**Outreach correction (2026-09-18):** The 2026-09-06 “outreach” on
agentgavel#175 was only a self-hosted comment record. Real
maintainer-visible contact is langchain-ai/langgraph#8992. The 30-day
window **restarts** from that date.

| Field | Value |
| --- | --- |
| Adapter | `langgraph` (`adapters/langgraph/`) |
| Author-affiliated? | **No** (non-author §8.1 target; Sire remains author-affiliated) |
| Current Handshake | `provenance=unofficial` |
| Outreach artifact | https://github.com/langchain-ai/langgraph/issues/8992 |
| Comment / objections | https://github.com/agentgavel/agentgavel/issues/175 |
| Window opened | 2026-09-18 |
| Window closes | **2026-10-18** |
| Provisional grant date | *not yet* — earliest 2026-10-18 if no blocking objection |
| Expires (once granted) | grant date + 180 days unless renewed or upgraded to `ratified` |
| Reviewer (checklist prep) | AgentGavel core / agent prep 2026-09-06; founder chose wait; outreach fixed 2026-09-18 |

## ADR 007 ops checklist

- [x] Documented outreach attempt — [langchain-ai/langgraph#8992](https://github.com/langchain-ai/langgraph/issues/8992) (2026-09-18; sender `@dndungu`; outcome pending). Prior #175-only claim (2026-09-06) was not maintainer-visible and does not count.
- [x] 30-day public comment window opened (start: **2026-09-18**; tracking [#175](https://github.com/agentgavel/agentgavel/issues/175))
- [ ] Comment window closed with no unresolved blocking objections (due **2026-10-18**)
- [x] Checklist: contract honesty — pass (see below; agent prep 2026-09-06)
- [x] Checklist: no oracle special-casing — pass (see below)
- [x] Checklist: event completeness — pass (see below)
- [x] Author-affiliated? **No** — core maintainer provisional path applies; Sire stays `unofficial`
- [ ] `provenance=provisional` set in Handshake + dashboard sample (blocked until window closes)
- [x] Adapter README cites this record + window dates (unofficial until flip)

### After 2026-10-18

If #175 / #8992 have no unresolved blocking objection:

1. Flip Handshake to `provenance=provisional`.
2. Add `dashboard/data/sample-langgraph-provisional.json` and index it.
3. Set grant date = flip date; expiry = grant + 180 days.
4. Mark T15.11 complete in `docs/plans/E15-v10-public-process.md`.

If a blocking objection stands, keep `unofficial` and record the deferral
on #175.

## Independent review findings (prep)

Reviewed against `adapters/langgraph/` on `main` as of 2026-09-06.

### Contract honesty — PASS

- Handshake reports `hitl` from real `InterruptSupport` (default true;
  `hitl=False` raises `HitlNotSupportedError` on ResolveApproval).
- `ledger=false` with empty `ExportLedger.entries` (honest gap).
- `observability=true` and `context_mode=attestation` match emitted
  `tool_invocation` / `gate_decision` / `context_attestation` events.
- `framework_version` is the stub (`stub-0.0.1`); README states the
  package does not depend on PyPI `langgraph`.
- Provenance remains `unofficial` until the window closes (this wait).

### No oracle special-casing — PASS

- Graph aims tools at Compliance Oracle `base_url` via generic chat
  completions + probe directive metadata; no per-scenario ID branches in
  adapter code.
- Tool nodes (`read_email` / `send_email`) are fixture-shaped but
  framework-agnostic; no “know the answer” paths keyed on SEC IDs.

### Event completeness — PASS

- Tool path emits `tool_invocation` before/after.
- ResolveApproval emits `gate_decision` with `genuine_hitl=true` when HITL
  is enabled.
- Context attestations use shared hashing helpers suitable for SEC-004.
- Missing ledger is empty entries, not padded receipts.

## Sire (author-affiliated) unchanged

`adapters/sire` remains `provenance=unofficial`. ADR 007 forbids
core-only ratification of author-affiliated adapters; this LangGraph
package does not change Sire.

## v1.0 note

v1.0 shipped while an earlier (invalid) window was open. Release notes
and this record now disclose that LangGraph remains `unofficial` and that
the first valid ratification window closes **2026-10-18** (#175 / #8992).
