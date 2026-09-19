"""T17.10: Hermes live env bootstrap — stub default, fail-closed probe."""

from __future__ import annotations

import json
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Any

import pytest

from adapters.hermes.adapter import HermesAdapter
from adapters.hermes.client import (
    HermesClientError,
    HttpHermesClient,
    StubHermesClient,
    client_from_env,
)


def test_stub_default_runtime() -> None:
    report = HermesAdapter().handshake("1.0")
    assert report["runtime"] == "stub"
    assert report["framework_version"] == "stub"


def test_client_from_env_defaults_stub(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.delenv("AGENTGAVEL_HERMES_API_BASE", raising=False)
    monkeypatch.delenv("HERMES_API_BASE", raising=False)
    client = client_from_env()
    assert isinstance(client, StubHermesClient)
    assert HermesAdapter(client=client).handshake("1.0")["runtime"] == "stub"


class _CapsHandler(BaseHTTPRequestHandler):
    def log_message(self, format: str, *args: Any) -> None:  # noqa: A003
        del format, args

    def do_GET(self) -> None:  # noqa: N802
        if self.path != "/v1/capabilities":
            self.send_response(404)
            self.end_headers()
            return
        body = json.dumps({"version": "hermes-test-9"}).encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


@pytest.fixture()
def hermes_base() -> Any:
    server = HTTPServer(("127.0.0.1", 0), _CapsHandler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    host, port = server.server_address[:2]
    try:
        yield f"http://{host}:{port}"
    finally:
        server.shutdown()
        thread.join(timeout=5)


def test_client_from_env_live_after_probe(
    monkeypatch: pytest.MonkeyPatch,
    hermes_base: str,
) -> None:
    monkeypatch.setenv("AGENTGAVEL_HERMES_API_BASE", hermes_base)
    client = client_from_env()
    assert isinstance(client, HttpHermesClient)
    report = HermesAdapter(client=client).handshake("1.0")
    assert report["runtime"] == "live"
    assert report["framework_version"] == "hermes-test-9"


def test_client_from_env_fail_closed(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("AGENTGAVEL_HERMES_API_BASE", "http://127.0.0.1:1")
    with pytest.raises(HermesClientError):
        client_from_env()
