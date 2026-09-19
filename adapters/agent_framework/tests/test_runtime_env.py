"""Agent Framework stub default + live package probe (ADR 015)."""

from __future__ import annotations

import pytest

from adapters.agent_framework.adapter import AgentFrameworkAdapter
from adapters.agent_framework.runtime_env import adapter_from_env, runtime_from_env


def test_default_runtime_stub() -> None:
    report = AgentFrameworkAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["framework_version"] == "stub-0.0.1"


def test_adapter_from_env_defaults_stub(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_AGENT_FRAMEWORK_RUNTIME", raising=False)
    monkeypatch.delenv("AGENT_FRAMEWORK_RUNTIME", raising=False)
    assert runtime_from_env() == "stub"
    assert adapter_from_env().handshake("1.0")["runtime"] == "stub"


def test_live_fail_closed(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("AGENTGAVEL_AGENT_FRAMEWORK_RUNTIME", "live")

    def _boom() -> str:
        raise RuntimeError(
            "agent-framework package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-agent-framework[live]'"
        )

    monkeypatch.setattr(
        "adapters.agent_framework.runtime_env.require_agent_framework",
        _boom,
    )
    with pytest.raises(RuntimeError, match="agent-framework package required"):
        adapter_from_env()
