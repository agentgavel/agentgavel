"""Protocol Event builders for OpenClaw Gateway frame mapping (T15.19).

Keeps ``tool_invocation`` / ``gate_decision`` shapes aligned with
``proto/adapter.proto`` without touching ResolveApproval (T15.18 owns HITL).

OpenClaw Gateway frames (capability map): ``session.tool``,
``session.approval``, ``exec.approval.requested`` / ``exec.approval.resolved``.
Audit RPC is metadata-only — never invent ``ledger_append`` / hash-linked
entries here. CapabilityReport keeps ``ledger=false`` and, until a live
Gateway subscription proves before/after completeness, ``observability=false``.
"""

from __future__ import annotations

import json
import time
from collections.abc import Mapping, MutableMapping
from typing import Any

GATE_SOURCE_HARNESS = "harness"
GATE_SOURCE_STORE = "store"
GATE_SOURCE_TOOL_OUTPUT = "tool_output"
GATE_SOURCE_LLM = "llm"

_DECISION_BY_INT = {
    1: "approve",
    2: "deny",
    3: "withhold",
}
_DECISION_NAMES = frozenset({"approve", "deny", "withhold"})

# Gateway event type → wire kind (honest N/A for unknown / audit-only frames).
_GATEWAY_TOOL_TYPES = frozenset(
    {
        "session.tool",
        "tool",
        "tool_invocation",
    }
)
_GATEWAY_GATE_TYPES = frozenset(
    {
        "session.approval",
        "exec.approval.resolved",
        "exec.approval.requested",
        "plugin.approval.resolved",
        "gate_decision",
    }
)


def normalize_decision(decision: str | int) -> str:
    """Normalize ResolveApproval / Gateway decision to AgentGavel wire names."""
    if isinstance(decision, bool):
        raise ValueError(f"unknown decision {decision!r}")
    if isinstance(decision, int):
        try:
            return _DECISION_BY_INT[decision]
        except KeyError as exc:
            raise ValueError(f"unknown decision enum {decision!r}") from exc
    name = str(decision).strip()
    if name.startswith("DECISION_"):
        name = name.removeprefix("DECISION_").lower()
    else:
        name = name.lower()
    if name not in _DECISION_NAMES:
        raise ValueError(f"unknown decision {decision!r}")
    return name


def build_tool_invocation(
    tool_name: str,
    tool_id: str,
    phase: str,
    *,
    arguments: Mapping[str, Any] | None = None,
    outcome: str | None = None,
    error: str | None = None,
    refused: bool = False,
) -> dict[str, Any]:
    """Build a ``tool_invocation`` payload (before | after)."""
    if phase not in ("before", "after"):
        raise ValueError(f"phase must be before|after, got {phase!r}")
    inv: dict[str, Any] = {
        "tool_name": tool_name,
        "tool_id": tool_id,
        "phase": phase,
    }
    if arguments is not None:
        inv["arguments_json"] = json.dumps(dict(arguments), separators=(",", ":"))
    if outcome is not None:
        inv["outcome"] = outcome
    if error is not None:
        inv["error"] = error
    if refused:
        inv["refused"] = True
    return inv


def build_gate_decision(
    approval_id: str,
    decision: str | int,
    *,
    source: str = GATE_SOURCE_HARNESS,
    principal: str | None = None,
    genuine_hitl: bool = False,
) -> dict[str, Any]:
    """Build a ``gate_decision`` payload for HITL / Gateway approval frames."""
    gate: dict[str, Any] = {
        "approval_id": approval_id,
        "source": source,
        "decision": normalize_decision(decision),
        "genuine_hitl": genuine_hitl,
    }
    if principal:
        gate["principal"] = principal
    return gate


def make_event(
    session_id: str,
    seq: int,
    *,
    unix_ms: int | None = None,
    tool_invocation_payload: Mapping[str, Any] | None = None,
    gate_decision_payload: Mapping[str, Any] | None = None,
) -> MutableMapping[str, Any]:
    """Assemble a protocol Event with exactly one kind payload."""
    kinds = [
        ("tool_invocation", tool_invocation_payload),
        ("gate_decision", gate_decision_payload),
    ]
    present = [(k, v) for k, v in kinds if v is not None]
    if len(present) != 1:
        raise ValueError("make_event requires exactly one event kind payload")
    kind, payload = present[0]
    event: MutableMapping[str, Any] = {
        "session_id": session_id,
        "seq": seq,
        "unix_ms": int(time.time() * 1000) if unix_ms is None else unix_ms,
        kind: dict(payload),
    }
    return event


def empty_ledger(session_id: str) -> dict[str, Any]:
    """Honest empty Ledger shape (capability map: audit is not hash-linked)."""
    return {"session_id": session_id, "entries": []}


def map_gateway_frame(
    session_id: str,
    seq: int,
    frame: Mapping[str, Any],
    *,
    unix_ms: int | None = None,
) -> MutableMapping[str, Any] | None:
    """Map an OpenClaw Gateway event frame to a wire Event, or None if N/A.

    Unknown / audit-only / incomplete frames return None (do not invent
    ledger_append or stub tool phases). Callers keep observability=false until
    a live subscription proves ordered before/after coverage.
    """
    event_type = str(frame.get("type") or frame.get("event") or frame.get("kind") or "").strip()
    payload = frame.get("payload")
    if not isinstance(payload, Mapping):
        payload = frame

    if event_type in _GATEWAY_TOOL_TYPES or ("tool_name" in payload and "phase" in payload):
        tool_name = str(payload.get("tool_name") or payload.get("name") or "")
        tool_id = str(payload.get("tool_id") or payload.get("id") or tool_name)
        phase = str(payload.get("phase") or "").lower()
        if phase not in ("before", "after") or not tool_name:
            return None
        args = payload.get("arguments")
        if args is not None and not isinstance(args, Mapping):
            args = None
        inv = build_tool_invocation(
            tool_name,
            tool_id,
            phase,
            arguments=args,
            outcome=str(payload["outcome"]) if payload.get("outcome") is not None else None,
            error=str(payload["error"]) if payload.get("error") is not None else None,
            refused=bool(payload.get("refused", False)),
        )
        return make_event(
            session_id,
            seq,
            unix_ms=unix_ms,
            tool_invocation_payload=inv,
        )

    if event_type in _GATEWAY_GATE_TYPES or ("approval_id" in payload and "decision" in payload):
        approval_id = str(payload.get("approval_id") or payload.get("id") or "")
        if not approval_id or payload.get("decision") is None:
            # Requested-only frames without a decision are N/A for gate_decision.
            return None
        try:
            gate = build_gate_decision(
                approval_id,
                payload["decision"],  # type: ignore[arg-type]
                source=str(payload.get("source") or GATE_SOURCE_HARNESS),
                principal=(
                    str(payload["principal"]) if payload.get("principal") is not None else None
                ),
                genuine_hitl=bool(payload.get("genuine_hitl", False)),
            )
        except ValueError:
            return None
        return make_event(
            session_id,
            seq,
            unix_ms=unix_ms,
            gate_decision_payload=gate,
        )

    return None


def tool_invocation_phases(
    events: list[Mapping[str, Any]],
) -> list[tuple[str, str, str]]:
    """Return (tool_id_or_name, phase, tool_name) for tool_invocation events."""
    out: list[tuple[str, str, str]] = []
    for ev in events:
        inv = ev.get("tool_invocation")
        if not isinstance(inv, Mapping):
            continue
        name = str(inv.get("tool_name") or "")
        tid = str(inv.get("tool_id") or "") or name
        phase = str(inv.get("phase") or "")
        out.append((tid, phase, name))
    return out


def assert_tool_invocation_order(events: list[Mapping[str, Any]]) -> None:
    """Raise AssertionError if an after-phase precedes its before (UC-004)."""
    seen_before: set[str] = set()
    for tid, phase, name in tool_invocation_phases(events):
        key = tid or name
        if phase == "before":
            seen_before.add(key)
        elif phase == "after":
            if key not in seen_before:
                raise AssertionError(f"tool_invocation after precedes before for {key!r}")
