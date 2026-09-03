"""CrewAI adapter (provenance=unofficial).

T13.15 scaffold + T13.21 minimal Oracle tool path (read_email/send_email)
mirroring LangGraph T11.2. Capability flags stay honest: observability and
context attestation are real; hitl / tenancy / ledger remain false until
wired. Unofficial until a maintainer signs off (ADR 007).
"""

from __future__ import annotations

from collections.abc import Mapping, MutableMapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

from adapters.crewai.crew import MinimalEmailCrew

_ADAPTER_VERSION = "0.0.1"


class HitlNotSupportedError(RuntimeError):
    """ResolveApproval called while CapabilityReport.hitl is false."""


class CrewAIAdapter(Adapter):
    """Unofficial CrewAI sidecar: Handshake + Oracle crew tools."""

    def __init__(self) -> None:
        super().__init__()
        self._sessions: dict[str, dict[str, Any]] = {}
        self._seq: dict[str, int] = {}
        self.emitted: list[MutableMapping[str, Any]] = []
        self.last_task_result: dict[str, Mapping[str, Any]] = {}

    def handshake(
        self,
        engine_protocol_version: str,
        *,
        engine_version: str | None = None,
    ) -> Mapping[str, Any]:
        del engine_version  # reserved for future negotiation
        return {
            "adapter_protocol_version": engine_protocol_version or "1.0",
            "adapter_name": "crewai",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: unofficial until maintainer ratification.
            "provenance": "unofficial",
            # Honest: tool_invocation + attestation wired; HITL/ledger not.
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": True,
            "context_mode": "attestation",
            "framework_name": "crewai",
            "framework_version": "stub-0.0.1",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = f"crewai-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = dict(config)
        self._seq[session_id] = 0
        return {"id": session_id}

    def submit_task(self, session_id: str, task: Mapping[str, Any]) -> None:
        config = self._sessions.get(session_id)
        if config is None:
            raise KeyError(f"unknown session: {session_id}")
        base_url = str(config.get("model_base_url") or "").strip()
        if not base_url:
            # No Oracle binding yet — lifecycle no-op (scaffold behavior).
            return
        meta = task.get("metadata") if isinstance(task.get("metadata"), Mapping) else {}
        directive = meta.get("probe_directive") if isinstance(meta, Mapping) else None
        if directive is not None and not isinstance(directive, Mapping):
            raise TypeError("task.metadata.probe_directive must be a mapping")
        crew = MinimalEmailCrew(
            model_base_url=base_url,
            session_id=session_id,
            on_event=self._record_event,
        )
        result = crew.run(
            str(task.get("prompt") or ""),
            probe_directive=directive,
        )
        self.last_task_result[session_id] = result

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
            "CrewAI HITL/human_input mapping not wired; CapabilityReport.hitl is false"
        )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)
        self.last_task_result.pop(session_id, None)
        self._seq.pop(session_id, None)

    def _next_seq(self, session_id: str) -> int:
        n = self._seq.get(session_id, 0) + 1
        self._seq[session_id] = n
        return n

    def _record_event(self, event: MutableMapping[str, Any]) -> None:
        # Stamp a session-monotonic seq so crew + adapter share one clock.
        sid = event.get("session_id")
        if isinstance(sid, str) and sid:
            event["seq"] = self._next_seq(sid)
        self.emitted.append(event)
        if self._transport is not None:
            self.emit(event)
