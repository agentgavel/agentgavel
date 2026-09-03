"""AWS Strands Agents adapter (provenance=unofficial).

T13.14: Handshake + session lifecycle scaffold. No strands-agents / tool path
yet (T13.20). CapabilityReport stays honest: hitl/tenancy/ledger/observability
false, context_mode=none. Unofficial until a maintainer signs off (ADR 007).
"""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

_ADAPTER_VERSION = "0.0.1"
_FRAMEWORK_VERSION = "stub-0.0.1"


class HitlNotSupportedError(RuntimeError):
    """Raised when ResolveApproval is called but Handshake reports hitl=false."""


class StrandsAdapter(Adapter):
    """Unofficial AWS Strands sidecar: Handshake + lifecycle (no tool path yet)."""

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
            "adapter_name": "strands",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: unofficial until maintainer ratification.
            "provenance": "unofficial",
            # Honest scaffold: no interrupt flow / events / ledger yet.
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": "aws-strands",
            "framework_version": _FRAMEWORK_VERSION,
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = f"strands-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = dict(config)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Scaffold: record task only; Oracle tool path lands in T13.20.
        self._sessions[session_id]["last_task"] = dict(task)

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
            "AWS Strands Agents HITL/interrupt not wired; CapabilityReport.hitl is false"
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)
