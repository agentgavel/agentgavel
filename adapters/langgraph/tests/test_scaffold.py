"""Smoke tests for the LangGraph adapter scaffold (T11.1)."""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path

from adapters.langgraph.adapter import LangGraphAdapter

_SRC = Path(__file__).resolve().parents[1] / "src"
_SDK_SRC = Path(__file__).resolve().parents[3] / "sdk" / "python" / "src"


def test_handshake_provenance_unofficial() -> None:
    report = LangGraphAdapter().handshake("1.0", engine_version="0.0.0-dev")
    assert report["adapter_name"] == "langgraph"
    assert report["provenance"] == "unofficial"
    # Default: interrupt support on ⇒ hitl=true (T11.3).
    assert report["hitl"] is True
    assert report["ledger"] is False
    assert report["observability"] is False
    assert report["context_mode"] == "none"


def test_lifecycle_without_oracle_does_not_crash() -> None:
    adapter = LangGraphAdapter(hitl=False)
    session = adapter.start_session({})
    sid = session["id"]
    assert sid
    adapter.submit_task(sid, {"id": "t1", "prompt": "noop"})
    ledger = adapter.export_ledger(sid)
    assert ledger["session_id"] == sid
    assert ledger["entries"] == []
    adapter.stop_session(sid)


def test_module_help_exits_zero() -> None:
    env = os.environ.copy()
    env["PYTHONPATH"] = f"{_SRC}{os.pathsep}{_SDK_SRC}"
    proc = subprocess.run(
        [sys.executable, "-m", "adapters.langgraph", "--help"],
        capture_output=True,
        text=True,
        env=env,
        check=False,
    )
    assert proc.returncode == 0
    assert "usage" in proc.stdout.lower() or "LangGraph adapter" in proc.stdout
