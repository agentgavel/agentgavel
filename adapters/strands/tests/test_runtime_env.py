"""Strands stub default + live package probe (ADR 015)."""

from __future__ import annotations

import pytest

from adapters.strands.adapter import StrandsAdapter
from adapters.strands.runtime_env import adapter_from_env, runtime_from_env


def test_default_runtime_stub() -> None:
    report = StrandsAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["framework_version"] == "stub-0.0.1"


def test_adapter_from_env_defaults_stub(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_STRANDS_RUNTIME", raising=False)
    monkeypatch.delenv("STRANDS_RUNTIME", raising=False)
    assert runtime_from_env() == "stub"
    assert adapter_from_env().handshake("1.0")["runtime"] == "stub"


def test_live_fail_closed(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("AGENTGAVEL_STRANDS_RUNTIME", "live")

    def _boom() -> str:
        raise RuntimeError(
            "strands-agents package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-strands[live]'"
        )

    monkeypatch.setattr("adapters.strands.runtime_env.require_strands", _boom)
    with pytest.raises(RuntimeError, match="strands-agents package required"):
        adapter_from_env()
