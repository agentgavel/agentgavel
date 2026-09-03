"""Smoke tests for the CrewAI adapter scaffold (T13.15 / T13.21 / UC-030)."""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path

import pytest

from adapters.crewai.adapter import CrewAIAdapter, HitlNotSupportedError

_SRC = Path(__file__).resolve().parents[1] / "src"
_SDK_SRC = Path(__file__).resolve().parents[3] / "sdk" / "python" / "src"


def test_handshake_provenance_unofficial() -> None:
    report = CrewAIAdapter().handshake("1.0", engine_version="0.0.0-dev")
    assert report["adapter_name"] == "crewai"
    assert report["framework_name"] == "crewai"
    assert report["provenance"] == "unofficial"
    # Honest capabilities after T13.21 tool path.
    assert report["hitl"] is False
    assert report["tenancy"] is False
    assert report["ledger"] is False
    assert report["observability"] is True
    assert report["context_mode"] == "attestation"


def test_lifecycle_stubs_do_not_crash() -> None:
    adapter = CrewAIAdapter()
    session = adapter.start_session({})
    sid = session["id"]
    assert sid
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    ledger = adapter.export_ledger(sid)
    assert ledger["session_id"] == sid
    assert ledger["entries"] == []
    adapter.stop_session(sid)


def test_resolve_approval_refuses_when_hitl_false() -> None:
    adapter = CrewAIAdapter()
    sid = adapter.start_session({})["id"]
    with pytest.raises(HitlNotSupportedError):
        adapter.resolve_approval(sid, "appr-1", "approve")


def test_module_help_exits_zero() -> None:
    env = os.environ.copy()
    env["PYTHONPATH"] = f"{_SRC}{os.pathsep}{_SDK_SRC}"
    proc = subprocess.run(
        [sys.executable, "-m", "adapters.crewai", "--help"],
        capture_output=True,
        text=True,
        env=env,
        check=False,
    )
    assert proc.returncode == 0
    assert "usage" in proc.stdout.lower() or "CrewAI adapter" in proc.stdout
