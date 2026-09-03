"""OpenAI Agents SDK adapter (provenance=unofficial).

T13.11 scaffold: Handshake + lifecycle stubs. No ``openai-agents`` PyPI
dependency yet — a later wave (T13.17) wires the minimal Oracle tool path.
Unofficial until a maintainer or independent reviewer signs off (ADR 007).
"""

from __future__ import annotations

from collections.abc import Mapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

_ADAPTER_VERSION = "0.0.1"


class HitlNotSupportedError(RuntimeError):
    """Raised when ResolveApproval is called while CapabilityReport.hitl is false."""


class OpenAIAgentsAdapter(Adapter):
    """Unofficial OpenAI Agents SDK sidecar: Handshake + honest stub capabilities."""

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
            "adapter_name": "openai_agents",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: unofficial until maintainer / external ratification.
            "provenance": "unofficial",
            # Scaffold: no needs_approval / interrupt mapping yet (T13.17+).
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": "openai-agents",
            "framework_version": "stub-0.0.1",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = f"openai-agents-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = dict(config)
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        # Scaffold: no Oracle tool path yet (T13.17).
        del task

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        del session_id, approval_id, decision, principal
        raise HitlNotSupportedError(
            "OpenAI Agents SDK HITL not wired; CapabilityReport.hitl is false"
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)
