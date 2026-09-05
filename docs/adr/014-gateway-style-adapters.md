# ADR 014: Gateway-Style Adapters (OpenClaw, Hermes Agent)

## Status
Accepted

## Date
2026-09-05

## Context
RFC §8.2 deferred workflow/app-builder platforms (n8n, Dify) because they
are not code-first libraries in the shape of the §5.3 sidecar contract.

Two self-hosted **Gateway products** have the governance surface AgentGavel
is meant to measure (tool policy, approvals, sandboxed/remote exec, audit)
without being LangGraph-style in-process SDKs:

- **OpenClaw** -- channels to a trusted Gateway with movable execution.
- **Hermes Agent** (Nous Research,
  https://github.com/nousresearch/hermes-agent) -- messaging gateway plus
  agent core (`AIAgent`), terminal backends (local/Docker/SSH/cloud),
  skills/memory, and platform adapters with native approval/clarify prompts.

v1.0 founder direction: include both as unofficial gateway-style targets.
n8n and Dify remain deferred.

## Decision

1. **v1.0 includes unofficial adapters** under `adapters/openclaw/` and
   `adapters/hermes/` (Hermes Agent), provenance `unofficial` until ADR 007
   ratification.
2. **Architecture class: gateway-style.** Each sidecar speaks the existing
   AgentGavel wire protocol (Handshake, StartSession, SubmitTask,
   ResolveApproval, ExportLedger, Events). Inside the sidecar, the adapter
   drives the product via Gateway/control API, CLI, or documented RPC --
   not by importing an agent graph library as the SUT.
3. **System under test** is each product's **agent runtime + Gateway policy
   plane** in a documented reference config. Nested/alternate execution
   backends (OpenClaw vendor harness plugins; Hermes terminal backends
   other than the reference) are out of scope unless CapabilityReport names
   them honestly as the execution backend.
4. **Honest N/A.** If a Gateway API cannot expose programmatic
   `ResolveApproval` or ledger export, CapabilityReport sets `hitl` /
   `ledger` false and SEC/REL scenarios N/A with the observability penalty
   -- never silent Fail or stubbed green.
5. **§8.2 remains deferred for n8n/Dify.** OpenClaw and Hermes do not
   unblock visual workflow platforms; those still need their own notes.
6. **No per-framework exploit code.** Probes stay in `fixtures/`; adapters
   only map protocol methods to Gateway surfaces.
7. **Separate fingerprints.** OpenClaw and Hermes never share a run
   fingerprint or scorecard row; each has its own adapter version and
   config-hash fields.

## Consequences
Positive: v1.0 measures two high-relevance operated agents; proves the wire
contract wraps Gateways without a second engine. Negative: more adapter and
CI cost; more N/A risk until APIs are mapped; `v1.0.0` waits on both
OpenClaw E2E (T15.20) and Hermes E2E (T15.27) in addition to Opt-in/bounty.

## Alternatives considered
- **Defer Hermes to post-v1.0** -- rejected by founder scope expand.
- **Treat Hermes as §8.1 library** -- inaccurate; primary UX is gateway +
  agent core, not an embeddable graph SDK.
- **New engine protocol for gateways** -- premature; reuse §5.3 first.
- **Single combined "gateway" scorecard** -- rejected; neutrality and
  fingerprints require per-product rows.
