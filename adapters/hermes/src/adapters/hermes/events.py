"""Protocol Event builders + Hermes API SSE frame mapping (T15.26).

Maps Hermes run/session stream frames onto AgentGavel wire Events
(``tool_invocation`` before/after, ``gate_decision``) per
``docs/manual/hermes-capability-map.md``. Trajectories / session export are
not hash-linked ledgers — see ``empty_ledger``.
"""

from __future__ import annotations

import json
import time
from collections.abc import Mapping, MutableMapping, Sequence
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

# Hermes chat/run SSE type strings (capability map + API server docs).
_TOOL_STARTED = frozenset(
    {
        "tool.started",
        "tool_started",
        "tool.start",
        "hermes.tool.progress",
    }
)
_TOOL_COMPLETED = frozenset(
    {
        "tool.completed",
        "tool_completed",
        "tool.complete",
        "tool.end",
    }
)
_GATE_TYPES = frozenset(
    {
        "approval.resolved",
        "approval.resolve",
        "approval.respond",
        "approval.decision",
        "gate_decision",
        "run.approval.resolved",
    }
)


def empty_ledger(session_id: str) -> dict[str, Any]:
    """Return a protocol ``Ledger`` with no hash-linked entries."""
    return {"session_id": session_id, "entries": []}


def normalize_decision(decision: str | int) -> str:
    """Normalize a decision to AgentGavel wire names (approve|deny|withhold)."""
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
    # Hermes / TUI synonyms.
    if name in {"allow", "allow_once", "allow_always", "o", "s", "a"}:
        return "approve"
    if name in {"reject", "d", "denied"}:
        return "deny"
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
    """Build a ``gate_decision`` payload for approval resolve / SSE frames."""
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


def tool_invocation_phases(
    events: Sequence[Mapping[str, Any]],
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


def assert_tool_invocation_order(events: Sequence[Mapping[str, Any]]) -> None:
    """Raise AssertionError if an after-phase precedes its before (UC-004)."""
    seen_before: set[str] = set()
    for tid, phase, name in tool_invocation_phases(events):
        key = tid or name
        if phase == "before":
            seen_before.add(key)
        elif phase == "after":
            if key not in seen_before:
                raise AssertionError(f"tool_invocation after precedes before for {key!r}")


def _frame_type(frame: Mapping[str, Any]) -> str:
    for key in ("type", "event", "kind", "name"):
        raw = frame.get(key)
        if isinstance(raw, str) and raw.strip():
            return raw.strip()
    data = frame.get("data")
    if isinstance(data, Mapping):
        for key in ("type", "event", "kind"):
            raw = data.get(key)
            if isinstance(raw, str) and raw.strip():
                return raw.strip()
    return ""


def _payload(frame: Mapping[str, Any]) -> Mapping[str, Any]:
    data = frame.get("data")
    if isinstance(data, Mapping):
        return data
    return frame


def _tool_fields(payload: Mapping[str, Any]) -> tuple[str, str]:
    name = str(
        payload.get("tool_name") or payload.get("name") or payload.get("tool") or "unknown_tool"
    )
    tid = str(payload.get("tool_id") or payload.get("call_id") or payload.get("id") or name)
    return name, tid


def map_hermes_frame(frame: Mapping[str, Any]) -> list[dict[str, Any]]:
    """Map one Hermes SSE/API frame to zero or more wire kind payloads.

    Returns a list of single-key dicts: ``{"tool_invocation": ...}`` or
    ``{"gate_decision": ...}``. Unknown / non-safety frames yield ``[]``.
    """
    if not isinstance(frame, Mapping):
        return []
    ftype = _frame_type(frame)
    payload = _payload(frame)
    if ftype in _TOOL_STARTED:
        name, tid = _tool_fields(payload)
        args = payload.get("arguments")
        if args is not None and not isinstance(args, Mapping):
            args = {"value": args}
        return [
            {
                "tool_invocation": build_tool_invocation(
                    name,
                    tid,
                    "before",
                    arguments=args if isinstance(args, Mapping) else None,
                )
            }
        ]
    if ftype in _TOOL_COMPLETED:
        name, tid = _tool_fields(payload)
        outcome = payload.get("outcome")
        if outcome is None and payload.get("output") is not None:
            outcome = str(payload.get("output"))
        error = payload.get("error")
        refused = bool(payload.get("refused") or payload.get("denied"))
        return [
            {
                "tool_invocation": build_tool_invocation(
                    name,
                    tid,
                    "after",
                    outcome=None if outcome is None else str(outcome),
                    error=None if error is None else str(error),
                    refused=refused,
                )
            }
        ]
    if ftype in _GATE_TYPES or ("approval_id" in payload and payload.get("decision") is not None):
        approval_id = str(
            payload.get("approval_id") or payload.get("id") or frame.get("approval_id") or "unknown"
        )
        decision = payload.get("decision")
        if decision is None:
            decision = payload.get("verb") or payload.get("action")
        if decision is None:
            return []
        source = str(payload.get("source") or GATE_SOURCE_STORE)
        principal = payload.get("principal")
        genuine = bool(payload.get("genuine_hitl", False))
        return [
            {
                "gate_decision": build_gate_decision(
                    approval_id,
                    decision,  # type: ignore[arg-type]
                    source=source,
                    principal=None if principal is None else str(principal),
                    genuine_hitl=genuine,
                )
            }
        ]
    return []
