# ADR 006: Leaderboard Submission Policy

## Status
Accepted

## Date
2026-09-02

## Context
Open question Q4: should leaderboard runs be opt-in only, and should AgentGavel
also publish unsolicited runs against public releases?

## Decision
Primary leaderboard tab is opt-in: maintainer-signed submissions with fingerprint
and adapter provenance. A clearly separate "Unratified / unsolicited" tab may
publish AgentGavel-operated runs against public releases. Unsolicited runs never
appear on the primary tab and always show provenance=unofficial unless the
adapter was previously ratified. Malicious auto-submissions are rejected by
requiring signatures bound to registered maintainer keys (v1.0 submission
process).

## Consequences
Positive: transparency without letting spam or hostile submissions define the
headline ranking.
Negative: two-tab UX must stay obvious; unsolicited runs may anger maintainers
who prefer silence -- mitigate with clear labeling and opt-out contact path.

## Addendum (2026-09-04, v0.3 static dashboard)
Until the v1.0 signed submission process exists, there is no way to verify
that a submission came from a maintainer. Therefore in v0.3:
- `AgentGavel report --publish` writes entries with `tab: "unratified"` only
  and rejects `--tab opt-in`.
- Entries on the Opt-in tab are hand-authored samples carrying
  `sample: true` and a framework name that says "(sample)". A CI check
  (`scripts/check-dashboard.sh`) enforces `tab=opt-in ⇒ sample=true` so a
  real run cannot be promoted by editing JSON.
- Every entry carries the ADR 007 three-way provenance
  (`ratified` / `provisional` / `unofficial`); unratified entries produced by
  `report --publish` copy provenance from the run's Handshake.
This addendum expires when E15 lands signatures; the CI rule then flips to
"opt-in requires a verified signature".
