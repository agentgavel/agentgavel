"""Smoke tests for the OpenClaw adapter scaffold (T15.17 / T15.18 N/A)."""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path

import pytest

from adapters.openclaw.adapter import HitlNotSupportedError, OpenClawAdapter

_SRC = Path(__file__).resolve().parents[1] / "src"
_SDK_SRC = Path(__file__).resolve().parents[3] / "sdk" / "python" / "src"


def test_handshake_provenance_and_versions() -> None:
    report = OpenClawAdapter().handshake("1.0", engine_version="0.0.0-dev")
    assert report["adapter_name"] == "openclaw"
    assert report["framework_name"] == "openclaw"
    assert report["provenance"] == "unofficial"
    assert report["adapter_version"]
    assert report["framework_version"]
    assert report["adapter_protocol_version"] == "1.0"
    # Conservative until live Gateway probe (capability map / zero-stub).
    assert report["hitl"] is False
    assert report["ledger"] is False
    assert report["observability"] is False
    assert report["tenancy"] is False
    assert report["context_mode"] == "none"


def test_lifecycle_scaffold_does_not_crash() -> None:
    adapter = OpenClawAdapter()
    session = adapter.start_session({})
    sid = session["id"]
    assert sid.startswith("openclaw-sess-")
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    ledger = adapter.export_ledger(sid)
    assert ledger["session_id"] == sid
    assert ledger["entries"] == []
    with pytest.raises(HitlNotSupportedError):
        adapter.resolve_approval(sid, "appr-1", "approve")
    adapter.stop_session(sid)


def test_module_help_exits_zero() -> None:
    env = os.environ.copy()
    env["PYTHONPATH"] = f"{_SRC}{os.pathsep}{_SDK_SRC}"
    proc = subprocess.run(
        [sys.executable, "-m", "adapters.openclaw", "--help"],
        capture_output=True,
        text=True,
        env=env,
        check=False,
    )
    assert proc.returncode == 0
    assert "usage" in proc.stdout.lower() or "OpenClaw adapter" in proc.stdout
