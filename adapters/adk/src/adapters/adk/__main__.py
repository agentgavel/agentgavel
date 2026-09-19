"""CLI entry: ``python -m adapters.adk`` serves stdio; ``--help`` prints usage."""

from __future__ import annotations

import argparse

from adapters.adk.runtime_env import adapter_from_env


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="python -m adapters.adk",
        description=(
            "Unofficial AgentGavel Google ADK adapter (stdio JSON-RPC). "
            "Provenance is always unofficial (ADR 007). "
            "Set AGENTGAVEL_ADK_RUNTIME=live after pip install '.[live]'."
        ),
    )
    parser.add_argument(
        "--serve",
        action="store_true",
        help="run the stdio JSON-RPC serve loop (default with no flags)",
    )
    return parser


def main(argv: list[str] | None = None) -> None:
    """Parse argv; ``--help`` prints usage, otherwise start stdio serve()."""
    _build_parser().parse_args(argv)
    adapter_from_env().serve()


if __name__ == "__main__":
    main()
