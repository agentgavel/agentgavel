# LangGraph provisional ratification record (T15.11)

Dated checklist and outreach package for granting
`provenance=provisional` on `adapters/langgraph/` under
[ADR 007](../../adr/007-adapter-ratification.md) and
[adapter-ratification ops](../adapter-ratification.md).

| Field | Value |
| --- | --- |
| Adapter | `langgraph` (`adapters/langgraph/`) |
| Author-affiliated? | **No** (non-author §8.1 target; Sire remains author-affiliated) |
| Grant date | 2026-09-06 (effective on founder merge of the T15.11 PR) |
| Expires | 2027-03-05 (grant + 180 days) unless renewed or upgraded to `ratified` |
| Outreach / comment issue | https://github.com/agentgavel/agentgavel/issues/175 |
| Reviewer | AgentGavel core / founder sign-off (merge = grant) |

## ADR 007 ops checklist

- [x] Documented outreach attempt — issue [#175](https://github.com/agentgavel/agentgavel/issues/175) (2026-09-06; public GitHub; LangGraph maintainers invited; outcome pending)
- [ ] 30-day public comment window — **opened 2026-09-06, closes 2026-10-06** on #175
- [ ] Comment window closed with no unresolved blocking objections
- [x] Checklist: contract honesty — pass (see below; agent prep 2026-09-06)
- [x] Checklist: no oracle special-casing — pass (see below)
- [x] Checklist: event completeness — pass (see below)
- [x] Author-affiliated? **No** — core maintainer provisional path applies; Sire stays `unofficial`
- [ ] `provenance=provisional` set in Handshake + dashboard sample (lands with founder-approved PR)
- [ ] Adapter README cites this record + expiry

### Founder clock decision (required)

ADR 007 requires the comment window to **close** before provisional is
granted. Issue #175 opened the window on 2026-09-06.

Choose one when merging the T15.11 PR:

1. **Wait** — merge only after 2026-10-06 (and update this section to mark
   the window closed). Do not flip Handshake before then.
2. **Clock exception** — founder (core maintainer) accepts provisional
   grant on 2026-09-06 with the window still open, documenting that
   outreach is public on #175 and objections filed there remain binding
   (provisional may be reverted if a blocking objection lands). Sign by
   merging the T15.11 PR with the provenance flip included.

Until one of those lands, treat this record as **prep only**.

## Independent review findings (prep)

Reviewed against `adapters/langgraph/` on `main` as of 2026-09-06
(`fb29b33` lineage).

### Contract honesty — PASS

- Handshake reports `hitl` from real `InterruptSupport` (default true;
  `hitl=False` raises `HitlNotSupportedError` on ResolveApproval).
- `ledger=false` with empty `ExportLedger.entries` (honest gap).
- `observability=true` and `context_mode=attestation` match emitted
  `tool_invocation` / `gate_decision` / `context_attestation` events.
- `framework_version` is the stub (`stub-0.0.1`); README states the
  package does not depend on PyPI `langgraph`.
- Provenance string was `unofficial` pending this grant; not overstated
  as ratified.

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
grant does not change Sire.

## Renewal / upgrade

- Renew before 2027-03-05 with a fresh checklist pass, or
- Upgrade to `ratified` when LangGraph maintainers review/contribute, or
- Revert published provenance to `unofficial` after expiry.
