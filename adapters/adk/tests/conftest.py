"""Shared fixtures for adapters/adk tests."""

from __future__ import annotations

import json
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Any

import pytest

from adapters.adk.graph import HEADER_PROBE_DIRECTIVE


class _OracleHandler(BaseHTTPRequestHandler):
    """Minimal OpenAI-shaped Compliance Oracle for unit tests."""

    def log_message(self, format: str, *args: Any) -> None:  # noqa: A003
        del format, args

    def do_GET(self) -> None:  # noqa: N802
        if self.path == "/healthz":
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"ok")
            return
        self.send_response(404)
        self.end_headers()

    def do_POST(self) -> None:  # noqa: N802
        if self.path != "/v1/chat/completions":
            self.send_response(404)
            self.end_headers()
            return
        length = int(self.headers.get("Content-Length", "0"))
        _ = self.rfile.read(length)
        raw_dir = self.headers.get(HEADER_PROBE_DIRECTIVE)
        if not raw_dir:
            self._json(
                400,
                {
                    "error": {
                        "message": "missing probe directive",
                        "type": "invalid_request_error",
                    }
                },
            )
            return
        directive = json.loads(raw_dir)
        tool_name = directive["tool_name"]
        arguments = directive.get("arguments") or {}
        self._json(
            200,
            {
                "id": "chatcmpl-test-oracle",
                "object": "chat.completion",
                "model": "oracle",
                "choices": [
                    {
                        "index": 0,
                        "message": {
                            "role": "assistant",
                            "content": None,
                            "tool_calls": [
                                {
                                    "id": "call_oracle_1",
                                    "type": "function",
                                    "function": {
                                        "name": tool_name,
                                        "arguments": json.dumps(arguments),
                                    },
                                }
                            ],
                        },
                        "finish_reason": "tool_calls",
                    }
                ],
            },
        )

    def _json(self, status: int, payload: dict[str, Any]) -> None:
        body = json.dumps(payload).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


@pytest.fixture()
def oracle_base_url() -> Any:
    server = HTTPServer(("127.0.0.1", 0), _OracleHandler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    host, port = server.server_address[:2]
    try:
        yield f"http://{host}:{port}"
    finally:
        server.shutdown()
        thread.join(timeout=5)
