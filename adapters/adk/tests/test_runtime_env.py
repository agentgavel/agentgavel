"""ADK stub default + live package probe (ADR 015)."""

from __future__ import annotations

import pytest

from adapters.adk.adapter import AdkAdapter
from adapters.adk.runtime_env import adapter_from_env, runtime_from_env


def test_default_runtime_stub() -> None:
    report = AdkAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["framework_version"] == "stub-0.0.1"


def test_adapter_from_env_defaults_stub(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_ADK_RUNTIME", raising=False)
    monkeypatch.delenv("ADK_RUNTIME", raising=False)
    assert runtime_from_env() == "stub"
    assert adapter_from_env().handshake("1.0")["runtime"] == "stub"


def test_live_fail_closed(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("AGENTGAVEL_ADK_RUNTIME", "live")

    def _boom() -> str:
        raise RuntimeError(
            "google-adk package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-adk[live]'"
        )

    monkeypatch.setattr("adapters.adk.runtime_env.require_adk", _boom)
    with pytest.raises(RuntimeError, match="google-adk package required"):
        adapter_from_env()
