"""Hermes Agent adapter (provenance=unofficial).

Gateway-style sidecar (ADR 014) for Hermes Agent
(https://github.com/nousresearch/hermes-agent). Speaks AgentGavel wire
protocol; drives Hermes via a documented control plane — not an in-process
graph SDK.

Intended control plane (pinned for T15.24+; see
docs/manual/hermes-capability-map.md): OpenAI-compatible API server
(``gateway/platforms/api_server.py``) — HTTP + SSE, language-agnostic,
``GET /v1/capabilities``, ``POST /v1/runs``, ``POST /v1/runs/{id}/approval``.
TUI gateway JSON-RPC and ACP are alternatives; this sidecar targets the API
server path.

T15.25: ResolveApproval posts to Hermes ``POST /v1/runs/{run_id}/approval``
(``choice`` once/deny; withhold leaves the gate pending). CapabilityReport
``hitl=true``. ledger/observability stay false until T15.26. Unofficial
until ADR 007 ratification. Stub client needs no live Hermes.
"""

from __future__ import annotations

import time
from collections.abc import Mapping, MutableMapping
from typing import Any

from agentgavel_adapter.adapter import Adapter

from adapters.hermes.client import HermesClient, StubHermesClient, wire_decision

_ADAPTER_VERSION = "0.0.1"
_GATE_SOURCE_HARNESS = "harness"


class HitlNotSupportedError(RuntimeError):
    """ResolveApproval called while CapabilityReport.hitl is false.

    Kept for API stability; T15.25 wires hitl=true via the API-server path.
    """


class HermesAdapter(Adapter):
    """Unofficial Hermes Agent sidecar: Handshake + ResolveApproval."""

    def __init__(self, client: HermesClient | None = None) -> None:
        super().__init__()
        self._client: HermesClient = client if client is not None else StubHermesClient()
        self._sessions: set[str] = set()
        self._seq: dict[str, int] = {}
        self.emitted: list[MutableMapping[str, Any]] = []

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        del engine_version  # reserved for future negotiation
        return {
            "adapter_protocol_version": engine_protocol_version or "1.0",
            "adapter_name": "hermes",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007 / ADR 014: unofficial until ratification.
            "provenance": "unofficial",
            # ResolveApproval → POST /v1/runs/{run_id}/approval (T15.25).
            "hitl": True,
            "tenancy": False,
            # Trajectories ≠ hash-linked wire Ledger (capability map).
            "ledger": False,
            # Events mapping deferred to T15.26.
            "observability": False,
            "context_mode": "none",
            "framework_name": "hermes",
            # No live Hermes required for unit tests; probe sets real version later.
            "framework_version": "unknown",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = self._client.start_session(config)
        self._sessions.add(session_id)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        self._client.submit_task(session_id, task)

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        wire = wire_decision(decision)
        self._client.resolve_approval(
            session_id,
            approval_id,
            wire,
            principal=principal,
        )
        gate: dict[str, Any] = {
            "approval_id": approval_id,
            "source": _GATE_SOURCE_HARNESS,
            "decision": wire,
            "genuine_hitl": True,
        }
        if principal:
            gate["principal"] = principal
        event: dict[str, Any] = {
            "session_id": session_id,
            "seq": self._next_seq(session_id),
            "unix_ms": int(time.time() * 1000),
            "gate_decision": gate,
        }
        self.emitted.append(event)
        if self._transport is not None:
            self.emit(event)

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Trajectories are not a hash-linked wire Ledger (capability map).
        # Full ExportLedger mapping is T15.26 (additive).
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._client.stop_session(session_id)
        self._sessions.discard(session_id)
        self._seq.pop(session_id, None)

    def _next_seq(self, session_id: str) -> int:
        n = self._seq.get(session_id, 0) + 1
        self._seq[session_id] = n
        return n
