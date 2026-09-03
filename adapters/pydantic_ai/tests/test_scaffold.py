"""Smoke tests for the Pydantic AI adapter scaffold (T13.12 / UC-027)."""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path

import pytest

from adapters.pydantic_ai.adapter import HitlNotSupportedError, PydanticAIAdapter

_SRC = Path(__file__).resolve().parents[1] / "src"
_SDK_SRC = Path(__file__).resolve().parents[3] / "sdk" / "python" / "src"


def test_handshake_provenance_unofficial() -> None:
    report = PydanticAIAdapter().handshake("1.0", engine_version="0.0.0-dev")
    assert report["adapter_name"] == "pydantic_ai"
    assert report["framework_name"] == "pydantic-ai"
    assert report["provenance"] == "unofficial"
    # Honest CapabilityReport: no false support claims (N/A over Fail).
    assert report["hitl"] is False
    assert report["tenancy"] is False
    assert report["ledger"] is False
    assert report["observability"] is False
    assert report["context_mode"] == "none"


def test_lifecycle_stubs_do_not_crash() -> None:
    adapter = PydanticAIAdapter()
    session = adapter.start_session({})
    sid = session["id"]
    assert sid
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    ledger = adapter.export_ledger(sid)
    assert ledger["session_id"] == sid
    assert ledger["entries"] == []
    adapter.stop_session(sid)


def test_resolve_approval_refuses_when_hitl_false() -> None:
    adapter = PydanticAIAdapter()
    sid = adapter.start_session({})["id"]
    with pytest.raises(HitlNotSupportedError):
        adapter.resolve_approval(sid, "appr-1", "approve")


def test_module_help_exits_zero() -> None:
    env = os.environ.copy()
    env["PYTHONPATH"] = f"{_SRC}{os.pathsep}{_SDK_SRC}"
    proc = subprocess.run(
        [sys.executable, "-m", "adapters.pydantic_ai", "--help"],
        capture_output=True,
        text=True,
        env=env,
        check=False,
    )
    assert proc.returncode == 0
    assert "usage" in proc.stdout.lower() or "Pydantic AI adapter" in proc.stdout
