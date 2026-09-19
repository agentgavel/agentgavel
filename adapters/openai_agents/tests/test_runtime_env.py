"""OpenAI Agents stub default + live package probe (ADR 015)."""

from __future__ import annotations

import pytest

from adapters.openai_agents.adapter import OpenAIAgentsAdapter
from adapters.openai_agents.runtime_env import adapter_from_env, runtime_from_env


def test_default_runtime_stub() -> None:
    report = OpenAIAgentsAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["framework_version"] == "stub-0.0.1"


def test_adapter_from_env_defaults_stub(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_OPENAI_AGENTS_RUNTIME", raising=False)
    monkeypatch.delenv("OPENAI_AGENTS_RUNTIME", raising=False)
    assert runtime_from_env() == "stub"
    assert adapter_from_env().handshake("1.0")["runtime"] == "stub"


def test_live_fail_closed(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("AGENTGAVEL_OPENAI_AGENTS_RUNTIME", "live")

    def _boom() -> str:
        raise RuntimeError(
            "openai-agents package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-openai-agents[live]'"
        )

    monkeypatch.setattr(
        "adapters.openai_agents.runtime_env.require_openai_agents",
        _boom,
    )
    with pytest.raises(RuntimeError, match="openai-agents package required"):
        adapter_from_env()
