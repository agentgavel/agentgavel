"""LangGraph adapter (provenance=unofficial).

T11.2: minimal in-process email tool graph (read_email/send_email) aimed at
the Compliance Oracle. T11.3: LangGraph-style interrupt mapping to
ResolveApproval when interrupt support is enabled (``hitl=true``); when
disabled, CapabilityReport keeps ``hitl=false`` honestly. Unofficial until a
maintainer signs off (ADR 007).
"""

from __future__ import annotations

import time
from collections.abc import Mapping, MutableMapping
from typing import Any
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

from adapters.langgraph.graph import MinimalEmailGraph, invoke_tool_node
from adapters.langgraph.interrupt import (
    HitlNotSupportedError,
    InterruptSupport,
    disabled_interrupt_support,
    gate_decision_event,
)

_ADAPTER_VERSION = "0.0.1"


class LangGraphAdapter(Adapter):
    """Unofficial LangGraph sidecar: Handshake + Oracle graph + optional HITL."""

    def __init__(
        self,
        *,
        hitl: bool = True,
        interrupt_support: InterruptSupport | None = None,
    ) -> None:
        super().__init__()
        if interrupt_support is not None:
            self._interrupt = interrupt_support
        elif hitl:
            self._interrupt = InterruptSupport(enabled=True)
        else:
            self._interrupt = disabled_interrupt_support()
        self._sessions: dict[str, dict[str, Any]] = {}
        self._seq: dict[str, int] = {}
        self.emitted: list[MutableMapping[str, Any]] = []
        # Last graph interrupt result per session (for resume / tests).
        self.last_task_result: dict[str, Mapping[str, Any]] = {}

    @property
    def hitl_supported(self) -> bool:
        """True when interrupt→ResolveApproval mapping is active."""
        return self._interrupt.enabled

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
            # Honest: hitl tracks real InterruptSupport, never a fake claim.
            "hitl": self._interrupt.enabled,
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
            interrupt_support=self._interrupt if self._interrupt.enabled else None,
        )
        result = graph.run(
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
        if not self._interrupt.enabled:
            raise HitlNotSupportedError(
                "LangGraph interrupt/HITL disabled; CapabilityReport.hitl is false"
            )
        if session_id not in self._sessions:
            raise KeyError(f"unknown session: {session_id}")
        pending = self._interrupt.resolve(
            session_id,
            approval_id,
            decision,
            principal=principal,
        )
        assert pending.decision is not None
        # seq is stamped in _record_event for a session-monotonic clock.
        event = gate_decision_event(
            session_id=session_id,
            seq=0,
            unix_ms=int(time.time() * 1000),
            approval_id=approval_id,
            decision=pending.decision,
            principal=principal,
        )
        self._record_event(event)

        # Resume: approve runs the deferred tool; deny/withhold refuse it.
        if pending.decision == "approve":
            invoke_tool_node(pending.tool_name, pending.arguments)
            self._emit_tool_after(
                session_id,
                pending.tool_name,
                pending.call_id,
                outcome="ok",
            )
        else:
            self._emit_tool_after(
                session_id,
                pending.tool_name,
                pending.call_id,
                outcome="refused",
                refused=True,
            )

    def export_ledger(self, session_id: str) -> Mapping[str, Any]:
        return {"session_id": session_id, "entries": []}

    def stop_session(self, session_id: str) -> None:
        self._sessions.pop(session_id, None)
        self.last_task_result.pop(session_id, None)
        self._interrupt.clear_session(session_id)
        self._seq.pop(session_id, None)

    def _record_event(self, event: MutableMapping[str, Any]) -> None:
        # Stamp a session-monotonic seq so graph + ResolveApproval share one clock.
        sid = event.get("session_id")
        if isinstance(sid, str) and sid:
            event["seq"] = self._next_seq(sid)
        self.emitted.append(event)
        if self._transport is not None:
            self.emit(event)

    def _next_seq(self, session_id: str) -> int:
        n = self._seq.get(session_id, 0) + 1
        self._seq[session_id] = n
        return n

    def _emit_tool_after(
        self,
        session_id: str,
        tool_name: str,
        call_id: str,
        *,
        outcome: str,
        refused: bool = False,
    ) -> None:
        inv: dict[str, Any] = {
            "tool_name": tool_name,
            "tool_id": call_id,
            "phase": "after",
            "outcome": outcome,
        }
        if refused:
            inv["refused"] = True
        event: MutableMapping[str, Any] = {
            "session_id": session_id,
            "unix_ms": int(time.time() * 1000),
            "tool_invocation": inv,
        }
        self._record_event(event)
