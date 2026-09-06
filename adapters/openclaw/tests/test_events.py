"""T15.19: tool_invocation / gate_decision emission + Gateway N/A paths."""

from __future__ import annotations

import pytest

from adapters.openclaw.adapter import OpenClawAdapter
from adapters.openclaw.events import (
    assert_tool_invocation_order,
    build_gate_decision,
    build_tool_invocation,
    make_event,
    map_gateway_frame,
)


def test_lifecycle_without_gateway_emits_no_events() -> None:
    """N/A path: scaffold has no Gateway subscription → empty Events buffer."""
    adapter = OpenClawAdapter()
    report = adapter.handshake("1.0")
    assert report["observability"] is False
    assert report["ledger"] is False

    sid = adapter.start_session({})["id"]
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    assert adapter.emitted == []
    adapter.stop_session(sid)


def test_emit_tool_invocation_before_after_ordered() -> None:
    """When capabilities allow (helpers wired), stream includes tool_invocation."""
    adapter = OpenClawAdapter()
    sid = adapter.start_session({})["id"]

    adapter.emit_tool_invocation(
        sid,
        "exec",
        "tool-1",
        "before",
        arguments={"cmd": "ls"},
    )
    adapter.emit_tool_invocation(
        sid,
        "exec",
        "tool-1",
        "after",
        outcome="ok",
    )

    assert_tool_invocation_order(list(adapter.emitted))
    tool_events = [e for e in adapter.emitted if "tool_invocation" in e]
    assert len(tool_events) == 2
    assert tool_events[0]["tool_invocation"]["phase"] == "before"
    assert tool_events[0]["tool_invocation"]["tool_name"] == "exec"
    assert "arguments_json" in tool_events[0]["tool_invocation"]
    assert tool_events[1]["tool_invocation"]["phase"] == "after"
    assert tool_events[1]["tool_invocation"]["outcome"] == "ok"
    assert tool_events[0]["seq"] == 1
    assert tool_events[1]["seq"] == 2
    assert tool_events[0]["session_id"] == sid


def test_emit_gate_decision_shape() -> None:
    """gate_decision Event shape for ResolveApproval / Gateway approval frames."""
    adapter = OpenClawAdapter()
    sid = adapter.start_session({})["id"]

    event = adapter.emit_gate_decision(
        sid,
        "appr-1",
        "deny",
        principal="operator",
        genuine_hitl=True,
    )

    assert "gate_decision" in event
    gate = event["gate_decision"]
    assert gate["approval_id"] == "appr-1"
    assert gate["decision"] == "deny"
    assert gate["source"] == "harness"
    assert gate["genuine_hitl"] is True
    assert gate["principal"] == "operator"
    assert adapter.emitted[0] is event


def test_ingest_gateway_session_tool_frame() -> None:
    adapter = OpenClawAdapter()
    sid = adapter.start_session({})["id"]

    event = adapter.ingest_gateway_frame(
        sid,
        {
            "type": "session.tool",
            "payload": {
                "tool_name": "bash",
                "tool_id": "t-9",
                "phase": "before",
                "arguments": {"script": "echo hi"},
            },
        },
    )
    assert event is not None
    assert event["tool_invocation"]["tool_name"] == "bash"
    assert event["tool_invocation"]["phase"] == "before"


def test_ingest_gateway_approval_resolved_frame() -> None:
    adapter = OpenClawAdapter()
    sid = adapter.start_session({})["id"]

    event = adapter.ingest_gateway_frame(
        sid,
        {
            "type": "exec.approval.resolved",
            "payload": {
                "approval_id": "exec-42",
                "decision": "approve",
                "genuine_hitl": True,
            },
        },
    )
    assert event is not None
    assert event["gate_decision"]["approval_id"] == "exec-42"
    assert event["gate_decision"]["decision"] == "approve"


def test_ingest_gateway_unknown_or_audit_frame_is_na() -> None:
    """Honest N/A: audit / unknown frames do not invent Events or burn seq."""
    adapter = OpenClawAdapter()
    sid = adapter.start_session({})["id"]

    assert (
        adapter.ingest_gateway_frame(
            sid,
            {"type": "audit.activity", "payload": {"id": "row-1"}},
        )
        is None
    )
    assert (
        adapter.ingest_gateway_frame(
            sid,
            {"type": "exec.approval.requested", "payload": {"approval_id": "a1"}},
        )
        is None
    )
    assert adapter.emitted == []

    # After N/A, a real tool frame still starts at seq 1.
    event = adapter.ingest_gateway_frame(
        sid,
        {
            "type": "session.tool",
            "payload": {"tool_name": "x", "tool_id": "1", "phase": "before"},
        },
    )
    assert event is not None
    assert event["seq"] == 1


def test_map_gateway_frame_rejects_incomplete_tool() -> None:
    assert map_gateway_frame("s", 1, {"type": "session.tool", "payload": {}}) is None
    assert (
        map_gateway_frame(
            "s",
            1,
            {"type": "session.tool", "payload": {"tool_name": "x", "phase": "during"}},
        )
        is None
    )


def test_builders_and_make_event_validation() -> None:
    inv = build_tool_invocation("t", "id", "before")
    gate = build_gate_decision("a", 2)
    assert gate["decision"] == "deny"
    with pytest.raises(ValueError):
        build_tool_invocation("t", "id", "middle")
    with pytest.raises(ValueError):
        make_event("s", 1)
    with pytest.raises(ValueError):
        make_event(
            "s",
            1,
            tool_invocation_payload=inv,
            gate_decision_payload=gate,
        )


def test_emit_unknown_session_raises() -> None:
    adapter = OpenClawAdapter()
    with pytest.raises(KeyError):
        adapter.emit_tool_invocation("missing", "t", "1", "before")
    with pytest.raises(KeyError):
        adapter.emit_gate_decision("missing", "a", "approve")
    with pytest.raises(KeyError):
        adapter.ingest_gateway_frame("missing", {"type": "session.tool"})
