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

T15.24: Handshake + lifecycle stubs only. hitl/ledger/observability stay
false until live probe (T15.25 ResolveApproval, T15.26 Events). Unofficial
until ADR 007 ratification. Scaffold tests do not require a live Hermes.
"""

from __future__ import annotations

from collections.abc import Mapping, MutableMapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

_ADAPTER_VERSION = "0.0.1"


class HitlNotSupportedError(RuntimeError):
    """ResolveApproval called while CapabilityReport.hitl is false."""


class HermesAdapter(Adapter):
    """Unofficial Hermes Agent sidecar: Handshake + deferred Gateway mapping."""

    def __init__(self) -> None:
        super().__init__()
        self._sessions: dict[str, dict[str, Any]] = {}
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
            # Conservative until API-server approval / events are probed.
            # Missing capabilities score N/A (never silent Fail / stub green).
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": "hermes",
            # No live Hermes required for scaffold; probe sets real version later.
            "framework_version": "unknown",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = f"hermes-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = dict(config)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # T15.24: no live Hermes / API-server submit yet.
        del task

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        del approval_id, decision, principal
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        raise HitlNotSupportedError(
            "Hermes API-server ResolveApproval not wired (T15.25); "
            "CapabilityReport.hitl is false"
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Trajectories are not a hash-linked wire Ledger (capability map).
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)
