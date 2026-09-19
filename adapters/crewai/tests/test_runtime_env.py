"""T18.6: CrewAI stub default + live package probe (ADR 015)."""

from __future__ import annotations

import importlib.util

import pytest

from adapters.crewai.adapter import CrewAIAdapter
from adapters.crewai.runtime_env import adapter_from_env, runtime_from_env


def test_default_runtime_stub() -> None:
    report = CrewAIAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["framework_version"] == "stub-0.0.1"


def test_adapter_from_env_defaults_stub(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_CREWAI_RUNTIME", raising=False)
    monkeypatch.delenv("CREWAI_RUNTIME", raising=False)
    assert runtime_from_env() == "stub"
    assert adapter_from_env().handshake("1.0")["runtime"] == "stub"


def test_live_fail_closed_without_package(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("AGENTGAVEL_CREWAI_RUNTIME", "live")

    def _boom() -> str:
        raise RuntimeError(
            "crewai package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-crewai[live]'"
        )

    monkeypatch.setattr("adapters.crewai.runtime_env.require_crewai", _boom)
    with pytest.raises(RuntimeError, match="crewai package required"):
        adapter_from_env()


@pytest.mark.skipif(
    importlib.util.find_spec("crewai") is None,
    reason="crewai optional extra not installed",
)
def test_live_sets_runtime_and_version(monkeypatch: pytest.MonkeyPatch) -> None:
    from importlib.metadata import version

    monkeypatch.setenv("AGENTGAVEL_CREWAI_RUNTIME", "live")
    report = adapter_from_env().handshake("1.0")
    assert report["runtime"] == "live"
    assert report["framework_version"] == version("crewai")
