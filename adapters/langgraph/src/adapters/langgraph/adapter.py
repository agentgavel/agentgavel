"""LangGraph adapter (provenance=unofficial).

T11.2: minimal in-process email tool graph (read_email/send_email) aimed at
the Compliance Oracle. T11.3: LangGraph-style interrupt mapping to
ResolveApproval when interrupt support is enabled (``hitl=true``); when
disabled, CapabilityReport keeps ``hitl=false`` honestly. T11.4: event hooks
(``tool_invocation`` before/after, ``gate_decision``, hashed context
attestations per ADR 005). T17.x: optional ``runtime=live`` path uses the real
``langgraph`` package (ADR 015). Unofficial until ADR 007 window closes
2026-10-18 (T15.11); see
docs/manual/ratification/langgraph-provisional-2026-09-06.md.
"""

from __future__ import annotations

import time
from collections.abc import Mapping, MutableMapping
from typing import Any, Literal
from uuid import uuid4

from agentgavel_adapter.adapter import Adapter

from adapters.langgraph.events import build_tool_invocation, make_event
from adapters.langgraph.graph import MinimalEmailGraph, invoke_tool_node
from adapters.langgraph.interrupt import (
    HitlNotSupportedError,
    InterruptSupport,
    disabled_interrupt_support,
    gate_decision_event,
    wire_decision,
)

_ADAPTER_VERSION = "0.0.1"

RuntimeClass = Literal["stub", "live"]


class LangGraphAdapter(Adapter):
    """Unofficial LangGraph sidecar (provisional pending 2026-10-18).

    Handshake + Oracle graph + optional HITL. Default ``runtime=stub`` keeps
    CI free of the ``langgraph`` PyPI dependency; set ``runtime=live`` (via
    constructor or :func:`adapter_from_env`) to drive a real StateGraph.
    """

    def __init__(
        self,
        *,
        hitl: bool = True,
        interrupt_support: InterruptSupport | None = None,
        runtime: RuntimeClass = "stub",
        framework_version: str | None = None,
    ) -> None:
        super().__init__()
        if runtime not in ("stub", "live"):
            raise ValueError(f"runtime must be stub|live, got {runtime!r}")
        if interrupt_support is not None:
            self._interrupt = interrupt_support
        elif hitl:
            self._interrupt = InterruptSupport(enabled=True)
        else:
            self._interrupt = disabled_interrupt_support()
        self._runtime: RuntimeClass = runtime
        if runtime == "live":
            self._framework_version = framework_version or "unknown"
        else:
            self._framework_version = framework_version or "stub-0.0.1"
        self._sessions: dict[str, dict[str, Any]] = {}
        self._seq: dict[str, int] = {}
        self._live_graphs: dict[str, Any] = {}
        self.emitted: list[MutableMapping[str, Any]] = []
        # Last graph interrupt result per session (for resume / tests).
        self.last_task_result: dict[str, Mapping[str, Any]] = {}

    @property
    def hitl_supported(self) -> bool:
        """True when interrupt→ResolveApproval mapping is active."""
        return self._interrupt.enabled

    @property
    def runtime(self) -> RuntimeClass:
        return self._runtime

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
            # ADR 007: unofficial until #175 closes 2026-10-18 (T15.11 / #8992).
            "provenance": "unofficial",
            # ADR 015: stub unless live StateGraph path is selected.
            "runtime": self._runtime,
            # Honest: hitl tracks real InterruptSupport, never a fake claim.
            "hitl": self._interrupt.enabled,
            "tenancy": False,
            "ledger": False,
            # tool_invocation before/after + gate_decision + context_attestation.
            "observability": True,
            "context_mode": "attestation",
            "framework_name": "langgraph",
            "framework_version": self._framework_version,
        }

    def start_session(self, config: Mapping[str, Any]) -> Mapping[str, Any]:
        session_id = f"langgraph-sess-{uuid4().hex[:12]}"
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

        if self._runtime == "live":
            from adapters.langgraph.live_graph import LiveEmailGraph

            graph = LiveEmailGraph(
                model_base_url=base_url,
                session_id=session_id,
                on_event=self._record_event,
                gated_tools=(
                    self._interrupt.gated_tools if self._interrupt.enabled else frozenset()
                ),
            )
            result = graph.run(
                str(task.get("prompt") or ""),
                probe_directive=directive,
            )
            if result.get("status") == "interrupted":
                # Mirror stub: register pending so ResolveApproval can resume.
                aid = str(result.get("approval_id") or "")
                if aid and self._interrupt.enabled:
                    self._interrupt.request(
                        session_id,
                        str(result.get("tool_name") or ""),
                        result.get("arguments") or {},
                        str(result.get("call_id") or ""),
                        approval_id=aid,
                    )
                self._live_graphs[session_id] = graph
            self.last_task_result[session_id] = result
            return

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

        if self._runtime == "live":
            live = self._live_graphs.get(session_id)
            if live is None:
                raise KeyError(f"no live interrupt for session: {session_id}")
            if live.pending_approval_id != approval_id:
                raise KeyError(
                    f"approval mismatch: got {approval_id!r}, want {live.pending_approval_id!r}"
                )
            wire = wire_decision(decision)
            # Mark InterruptSupport resolved for consistency with stub path.
            self._interrupt.resolve(
                session_id,
                approval_id,
                wire,
                principal=principal,
            )
            event = gate_decision_event(
                session_id=session_id,
                seq=0,
                unix_ms=int(time.time() * 1000),
                approval_id=approval_id,
                decision=wire,
                principal=principal,
            )
            self._record_event(event)
            result = live.resume(wire)
            self.last_task_result[session_id] = result
            self._live_graphs.pop(session_id, None)
            return

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
        self._live_graphs.pop(session_id, None)
        self._interrupt.clear_session(session_id)
        self._seq.pop(session_id, None)

    def _next_seq(self, session_id: str) -> int:
        n = self._seq.get(session_id, 0) + 1
        self._seq[session_id] = n
        return n

    def _record_event(self, event: MutableMapping[str, Any]) -> None:
        # Stamp a session-monotonic seq so graph + ResolveApproval share one clock.
        sid = event.get("session_id")
        if isinstance(sid, str) and sid:
            event["seq"] = self._next_seq(sid)
        self.emitted.append(event)
        if self._transport is not None:
            self.emit(event)

    def _emit_tool_after(
        self,
        session_id: str,
        tool_name: str,
        call_id: str,
        *,
        outcome: str,
        refused: bool = False,
    ) -> None:
        inv = build_tool_invocation(
            tool_name,
            call_id,
            "after",
            outcome=outcome,
            refused=refused,
        )
        event = make_event(
            session_id,
            0,
            tool_invocation_payload=inv,
        )
        self._record_event(event)
