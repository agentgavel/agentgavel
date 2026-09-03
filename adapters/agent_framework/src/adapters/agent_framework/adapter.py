"""Microsoft Agent Framework adapter (provenance=unofficial).

T13.13 scaffold: Handshake + session lifecycle stubs. Targets Microsoft
Agent Framework (AutoGen successor). Capability flags stay honestly false
until T13.19 wires a real tool path / HITL / events. Unofficial until a
maintainer signs off (ADR 007).
"""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

_ADAPTER_VERSION = "0.0.1"
_FRAMEWORK_NAME = "microsoft-agent-framework"
_FRAMEWORK_VERSION = "stub-0.0.1"


class HitlNotSupportedError(RuntimeError):
    """Raised when ResolveApproval is called but Handshake hitl is false."""


class AgentFrameworkAdapter(Adapter):
    """Unofficial Microsoft Agent Framework sidecar scaffold."""

    def __init__(self) -> None:
        super().__init__()
        self._sessions: dict[str, dict[str, Any]] = {}

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        del engine_version  # reserved for future negotiation
        return {
            "adapter_protocol_version": engine_protocol_version or "1.0",
            "adapter_name": "agent_framework",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: unofficial until maintainer ratification.
            "provenance": "unofficial",
            # Honest scaffold: no HITL / ledger / events yet (T13.19).
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": _FRAMEWORK_NAME,
            "framework_version": _FRAMEWORK_VERSION,
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = f"maf-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = dict(config)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Scaffold: no Oracle / tool path yet (T13.19).
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
            "Microsoft Agent Framework HITL disabled; CapabilityReport.hitl is false"
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)
