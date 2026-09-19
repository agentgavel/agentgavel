"""CLI entry: ``python -m adapters.openai_agents`` serves stdio; ``--help`` prints usage."""

from __future__ import annotations

import argparse

from adapters.openai_agents.runtime_env import adapter_from_env


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="python -m adapters.openai_agents",
        description=(
            "Unofficial AgentGavel OpenAI Agents adapter (stdio JSON-RPC). "
            "Provenance is always unofficial (ADR 007). "
            "Set AGENTGAVEL_OPENAI_AGENTS_RUNTIME=live after pip install '.[live]'."
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
