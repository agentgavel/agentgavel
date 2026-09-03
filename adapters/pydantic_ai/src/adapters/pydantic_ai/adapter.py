"""Pydantic AI adapter (provenance=unofficial).

T13.12: lightweight Handshake + lifecycle stubs. Does not depend on the
``pydantic-ai`` PyPI package (heavy transitive tree). CapabilityReport keeps
hitl/tenancy/ledger/observability false honestly until a later task wires
real support (T13.18 tool path). Unofficial until a maintainer signs off
(ADR 007).
"""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

_ADAPTER_VERSION = "0.0.1"


class HitlNotSupportedError(RuntimeError):
    """ResolveApproval called while CapabilityReport.hitl is false."""


class PydanticAIAdapter(Adapter):
    """Unofficial Pydantic AI sidecar: Handshake + honest false capabilities."""

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
            "adapter_name": "pydantic_ai",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: unofficial until maintainer ratification.
            "provenance": "unofficial",
            # Honest scaffold: no HITL / ledger / event sink yet (T13.18).
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": "pydantic-ai",
            "framework_version": "stub-0.0.1",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = f"pydantic-ai-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = dict(config)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Scaffold: no Oracle / agent wiring yet (T13.18).
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
            "Pydantic AI deferred-tools/HITL mapping not wired; "
            "CapabilityReport.hitl is false"
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)
