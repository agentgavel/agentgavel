"""T17.1 / T17.3: stub default + live env bootstrap (ADR 015)."""

from __future__ import annotations

import importlib.util

import pytest

from adapters.langgraph.adapter import LangGraphAdapter
from adapters.langgraph.runtime_env import adapter_from_env, runtime_from_env


def test_default_handshake_runtime_stub() -> None:
    report = LangGraphAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["framework_version"] == "stub-0.0.1"
    assert report["provenance"] == "unofficial"


def test_runtime_from_env_defaults_stub(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_LANGGRAPH_RUNTIME", raising=False)
    monkeypatch.delenv("LANGGRAPH_RUNTIME", raising=False)
    assert runtime_from_env() == "stub"
    adapter = adapter_from_env()
    assert adapter.handshake("1.0")["runtime"] == "stub"


def test_adapter_from_env_live_fail_closed_without_package(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    monkeypatch.setenv("AGENTGAVEL_LANGGRAPH_RUNTIME", "live")

    def _boom() -> None:
        raise RuntimeError(
            "langgraph package required for runtime=live; "
            "install with: pip install 'agentgavel-adapter-langgraph[live]'"
        )

    monkeypatch.setattr(
        "adapters.langgraph.runtime_env.require_langgraph",
        _boom,
    )
    with pytest.raises(RuntimeError, match="langgraph package required"):
        adapter_from_env()


@pytest.mark.skipif(
    importlib.util.find_spec("langgraph") is None,
    reason="langgraph optional extra not installed",
)
def test_adapter_from_env_live_sets_runtime_and_version(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    from importlib.metadata import version

    monkeypatch.setenv("AGENTGAVEL_LANGGRAPH_RUNTIME", "live")
    adapter = adapter_from_env()
    report = adapter.handshake("1.0")
    assert report["runtime"] == "live"
    assert report["framework_version"] == version("langgraph")
    assert report["framework_version"] != "stub-0.0.1"
