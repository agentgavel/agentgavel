"""T15.18: ResolveApproval N/A path — hitl=false, loud HitlNotSupportedError."""

from __future__ import annotations

import pytest

from adapters.openclaw.adapter import HitlNotSupportedError, OpenClawAdapter
from adapters.openclaw.approvals import (
    GATEWAY_RESOLVE_METHODS,
    OPENCLAW_ALLOW_ONCE,
    OPENCLAW_DENY,
    ApprovalMapError,
    openclaw_resolve_decision,
    wire_decision,
)


def test_handshake_hitl_false_documented_na() -> None:
    report = OpenClawAdapter().handshake("1.0")
    assert report["hitl"] is False
    assert report["provenance"] == "unofficial"


@pytest.mark.parametrize("decision", ["approve", "deny", "withhold", 1, 2, 3])
def test_resolve_approval_raises_hitl_not_supported(decision: str | int) -> None:
    adapter = OpenClawAdapter()
    sid = adapter.start_session({})["id"]
    with pytest.raises(HitlNotSupportedError, match=r"hitl is false"):
        adapter.resolve_approval(sid, "appr-1", decision, principal="harness")


def test_resolve_approval_unknown_session() -> None:
    adapter = OpenClawAdapter()
    with pytest.raises(KeyError, match="unknown session"):
        adapter.resolve_approval("missing", "appr-1", "approve")


def test_resolve_approval_error_mentions_sec002_na() -> None:
    """SEC-002 must not silently Fail: loud N/A path is explicit in the error."""
    adapter = OpenClawAdapter()
    sid = adapter.start_session({})["id"]
    with pytest.raises(HitlNotSupportedError, match=r"SEC-002"):
        adapter.resolve_approval(sid, "appr-1", "approve")


def test_wire_decision_normalize() -> None:
    assert wire_decision("approve") == "approve"
    assert wire_decision(2) == "deny"
    assert wire_decision("DECISION_WITHHOLD") == "withhold"


def test_openclaw_map_approve_deny() -> None:
    assert openclaw_resolve_decision("approve") == OPENCLAW_ALLOW_ONCE
    assert openclaw_resolve_decision("deny") == OPENCLAW_DENY
    assert openclaw_resolve_decision(1) == OPENCLAW_ALLOW_ONCE


def test_openclaw_map_withhold_unmapped() -> None:
    with pytest.raises(ApprovalMapError, match=r"withhold"):
        openclaw_resolve_decision("withhold")


def test_gateway_resolve_methods_documented() -> None:
    assert "exec.approval.resolve" in GATEWAY_RESOLVE_METHODS
    assert "approval.resolve" in GATEWAY_RESOLVE_METHODS
