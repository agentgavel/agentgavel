"""LangGraph adapter (provenance=unofficial).

T11.2 adds a minimal in-process email tool graph (read_email/send_email)
that points ``model_base_url`` at the Compliance Oracle. Capability flags
stay honest: HITL/ledger/observability remain false until later tasks map
real interrupt/ledger/hooks (T11.3+). Unofficial until a maintainer signs
off (ADR 007).
"""

from __future__ import annotations

from collections.abc import Mapping, MutableMapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

from adapters.langgraph.graph import MinimalEmailGraph

_ADAPTER_VERSION = "0.0.1"


class LangGraphAdapter(Adapter):
    """Unofficial LangGraph sidecar: Handshake + Oracle-backed email graph."""

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
            "adapter_name": "langgraph",
            "adapter_version": _ADAPTER_VERSION,
            # ADR 007: unofficial until maintainer ratification.
            "provenance": "unofficial",
            # Honest stubs — HITL/ledger/observability land in later tasks.
            "hitl": False,
            "tenancy": False,
            "ledger": False,
            "observability": False,
            "context_mode": "none",
            "framework_name": "langgraph",
            "framework_version": "stub-0.0.1",
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = f"langgraph-sess-{uuid4().hex[:12]}"
        self._sessions[session_id] = dict(config)
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
        graph = MinimalEmailGraph(
            model_base_url=base_url,
            session_id=session_id,
            on_event=self._record_event,
        )
        graph.run(
            str(task.get("prompt") or ""),
            probe_directive=directive,
        )

    def resolve_approval(
        self,
        session_id: str,
        approval_id: str,
        decision: str | int,
        *,
        principal: str | None = None,
    ) -> None:
        del session_id, approval_id, decision, principal
        return None

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)

    def _record_event(self, event: MutableMapping[str, Any]) -> None:
        self.emitted.append(event)
        if self._transport is not None:
            self.emit(event)
