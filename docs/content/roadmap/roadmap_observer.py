#!/usr/bin/env python3
"""
Poll docs/content/roadmap for file changes and log them to observer.md.
Uses simple sleep-based polling to avoid external dependencies.
"""

import argparse
import time
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Monitor docs/content/roadmap and log file changes."
    )
    parser.add_argument(
        "--interval",
        type=float,
        default=2.0,
        help="Polling interval in seconds (default: 2.0)",
    )
    return parser.parse_args()


def ensure_log_file(log_path: Path) -> None:
    if not log_path.exists():
        log_path.write_text("# Roadmap Observer\n\n", encoding="utf-8")


def log_event(log_path: Path, message: str) -> None:
    timestamp = time.strftime("%Y-%m-%d %H:%M:%S")
    with log_path.open("a", encoding="utf-8") as log_file:
        log_file.write(f"- [{timestamp}] {message}\n")


def snapshot(target_dir: Path, ignored: set[str]) -> dict[str, int]:
    state: dict[str, int] = {}
    for entry in target_dir.iterdir():
        if not entry.is_file() or entry.name in ignored:
            continue
        try:
            state[entry.name] = entry.stat().st_mtime_ns
        except FileNotFoundError:
            continue
    return state


def main() -> None:
    args = parse_args()
    target_dir = Path(__file__).resolve().parent
    log_path = target_dir / "observer.md"
    ignored = {log_path.name, Path(__file__).name}

    ensure_log_file(log_path)

    previous_state = snapshot(target_dir, ignored)
    log_event(
        log_path,
        f"Observer started; tracking {len(previous_state)} file(s); interval {args.interval:.2f}s.",
    )

    while True:
        time.sleep(args.interval)
        current_state = snapshot(target_dir, ignored)

        for name in sorted(current_state):
            if name not in previous_state:
                log_event(log_path, f"Detected new file {name}")
            elif current_state[name] != previous_state[name]:
                log_event(log_path, f"Detected modification in {name}")

        for name in sorted(previous_state):
            if name not in current_state:
                log_event(log_path, f"Detected removal of {name}")

        previous_state = current_state


if __name__ == "__main__":
    main()
