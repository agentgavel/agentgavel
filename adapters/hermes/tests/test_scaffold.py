"""Smoke + ResolveApproval tests for the Hermes Agent adapter (T15.24/T15.25)."""

from __future__ import annotations

import os
import subprocess
import sys
from collections.abc import Mapping
from pathlib import Path
from typing import Any

import pytest

from adapters.hermes.adapter import HermesAdapter
from adapters.hermes.client import (
    HermesClientError,
    HttpHermesClient,
    PATH_RUN_APPROVAL,
    StubHermesClient,
    hermes_approval_choice,
    wire_decision,
)

_SRC = Path(__file__).resolve().parents[1] / "src"
_SDK_SRC = Path(__file__).resolve().parents[3] / "sdk" / "python" / "src"


def test_handshake_hitl_true_after_resolve_wiring() -> None:
    report = HermesAdapter().handshake("1.0", engine_version="0.0.0-dev")
    assert report["adapter_name"] == "hermes"
    assert report["adapter_version"] == "0.0.1"
    assert report["framework_name"] == "hermes"
    assert report["framework_version"] == "unknown"
    assert report["adapter_protocol_version"] == "1.0"
    assert report["provenance"] == "unofficial"
    # T15.25: API-server ResolveApproval wired (capability map).
    assert report["hitl"] is True
    assert report["tenancy"] is False
    assert report["ledger"] is False
    assert report["observability"] is False
    assert report["context_mode"] == "none"


def test_lifecycle_stubs_do_not_crash() -> None:
    adapter = HermesAdapter()
    session = adapter.start_session({})
    sid = session["id"]
    assert sid.startswith("hermes-sess-")
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    ledger = adapter.export_ledger(sid)
    assert ledger["session_id"] == sid
    assert ledger["entries"] == []
    adapter.stop_session(sid)


def test_wire_decision_and_hermes_choice_mapping() -> None:
    assert wire_decision("approve") == "approve"
    assert wire_decision(2) == "deny"
    assert wire_decision("DECISION_WITHHOLD") == "withhold"
    assert hermes_approval_choice("approve") == "once"
    assert hermes_approval_choice("deny") == "deny"
    assert hermes_approval_choice("withhold") is None
    with pytest.raises(HermesClientError):
        wire_decision("yolo")


def test_resolve_approval_calls_client_and_buffers_gate_decision() -> None:
    client = StubHermesClient()
    adapter = HermesAdapter(client=client)
    sid = adapter.start_session({"run_mode": "oracle"})["id"]
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    adapter.resolve_approval(sid, "req_1", "approve", principal="harness")

    assert ("resolve_approval", (sid, "req_1", "approve", "harness", "once")) in client.calls
    assert len(adapter.emitted) == 1
    event = adapter.emitted[0]
    assert event["session_id"] == sid
    assert event["seq"] == 1
    assert event["gate_decision"] == {
        "approval_id": "req_1",
        "source": "harness",
        "decision": "approve",
        "genuine_hitl": True,
        "principal": "harness",
    }


def test_resolve_approval_withhold_records_without_hermes_choice() -> None:
    client = StubHermesClient()
    adapter = HermesAdapter(client=client)
    sid = adapter.start_session({})["id"]
    adapter.resolve_approval(sid, "req_w", "withhold")
    assert ("resolve_approval", (sid, "req_w", "withhold", None, None)) in client.calls
    assert adapter.emitted[0]["gate_decision"]["decision"] == "withhold"


def test_resolve_approval_unknown_session() -> None:
    adapter = HermesAdapter()
    with pytest.raises(KeyError, match="unknown session"):
        adapter.resolve_approval("missing", "appr-1", "deny")


class _FakeRequester:
    def __init__(self) -> None:
        self.calls: list[tuple[str, str, Mapping[str, Any] | None]] = []

    def request(
        self,
        method: str,
        path: str,
        *,
        json: Mapping[str, Any] | None = None,
    ) -> Mapping[str, Any] | None:
        self.calls.append((method, path, json))
        return {"object": "hermes.run.approval_response", "resolved": 1}


def test_http_client_maps_resolve_approval_to_runs_approval_path() -> None:
    transport = _FakeRequester()
    client = HttpHermesClient(transport)
    sid = client.start_session({})
    client.submit_task(sid, {"id": "t1", "metadata": {"hermes_run_id": "run_abc"}})
    client.resolve_approval(sid, "req_42", "deny", principal="harness")

    method, path, body = transport.calls[-1]
    assert method == "POST"
    assert path == PATH_RUN_APPROVAL.format(run_id="run_abc")
    assert body == {"choice": "deny", "request_id": "req_42"}


def test_http_client_withhold_skips_post() -> None:
    transport = _FakeRequester()
    client = HttpHermesClient(transport)
    sid = client.start_session({})
    client.submit_task(sid, {"metadata": {"hermes_run_id": "run_abc"}})
    client.resolve_approval(sid, "req_w", "withhold")
    assert transport.calls == []


def test_http_client_approve_maps_to_once() -> None:
    transport = _FakeRequester()
    client = HttpHermesClient(transport)
    sid = client.start_session({})
    client.submit_task(sid, {"metadata": {"hermes_run_id": "run_xyz"}})
    client.resolve_approval(sid, "req_a", "approve")
    assert transport.calls[-1][2] == {"choice": "once", "request_id": "req_a"}


def test_http_client_without_requester_fails_loudly() -> None:
    client = HttpHermesClient()
    with pytest.raises(HermesClientError, match="no requester"):
        client.start_session({})


def test_module_help_exits_zero() -> None:
    env = os.environ.copy()
    env["PYTHONPATH"] = f"{_SRC}{os.pathsep}{_SDK_SRC}"
    proc = subprocess.run(
        [sys.executable, "-m", "adapters.hermes", "--help"],
        capture_output=True,
        text=True,
        env=env,
        check=False,
    )
    assert proc.returncode == 0
    assert "usage" in proc.stdout.lower() or "Hermes" in proc.stdout
