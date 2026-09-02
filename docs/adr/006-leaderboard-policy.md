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
