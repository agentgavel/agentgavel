"""Smoke tests for the Microsoft Agent Framework adapter scaffold (T13.13)."""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path

import pytest

from adapters.agent_framework.adapter import AgentFrameworkAdapter, HitlNotSupportedError

_SRC = Path(__file__).resolve().parents[1] / "src"
_SDK_SRC = Path(__file__).resolve().parents[3] / "sdk" / "python" / "src"


def test_handshake_provenance_unofficial() -> None:
    report = AgentFrameworkAdapter().handshake("1.0", engine_version="0.0.0-dev")
    assert report["adapter_name"] == "agent_framework"
    assert report["framework_name"] == "microsoft-agent-framework"
    assert report["provenance"] == "unofficial"
    # Honest capabilities after T13.19 tool path.
    assert report["hitl"] is False
    assert report["ledger"] is False
    assert report["observability"] is True
    assert report["tenancy"] is False
    assert report["context_mode"] == "attestation"


def test_lifecycle_stubs_do_not_crash() -> None:
    adapter = AgentFrameworkAdapter()
    session = adapter.start_session({})
    sid = session["id"]
    assert sid
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    with pytest.raises(HitlNotSupportedError):
        adapter.resolve_approval(sid, "appr-1", "approve")
    ledger = adapter.export_ledger(sid)
    assert ledger["session_id"] == sid
    assert ledger["entries"] == []
    adapter.stop_session(sid)


def test_module_help_exits_zero() -> None:
    env = os.environ.copy()
    env["PYTHONPATH"] = f"{_SRC}{os.pathsep}{_SDK_SRC}"
    proc = subprocess.run(
        [sys.executable, "-m", "adapters.agent_framework", "--help"],
        capture_output=True,
        text=True,
        env=env,
        check=False,
    )
    assert proc.returncode == 0
    assert "usage" in proc.stdout.lower() or "Agent Framework" in proc.stdout
