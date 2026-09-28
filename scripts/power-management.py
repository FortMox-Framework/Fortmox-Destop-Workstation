#!/usr/bin/env python3
"""Inspect or set CPU frequency governors; dry-run is the default."""

import argparse
import glob
import os
from pathlib import Path


def main() -> int:
    """Plan governor changes, applying them only when explicitly requested."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("profile", nargs="?", choices=("performance", "balanced", "powersave"), default="balanced")
    parser.add_argument("--dry-run", action="store_true", default=True, help="show governor changes only (default)")
    parser.add_argument("--apply", action="store_false", dest="dry_run", help="apply governor changes; requires root")
    args = parser.parse_args()
    paths = sorted(glob.glob("/sys/devices/system/cpu/cpu[0-9]*/cpufreq/scaling_governor"))
    if not paths:
        print("No CPU frequency governor controls found.")
        return 0 if args.dry_run else 1
    if not args.dry_run and os.geteuid() != 0:
        parser.error("--apply requires root; re-run with sudo")
    for raw_path in paths:
        path = Path(raw_path)
        available_path = path.with_name("scaling_available_governors")
        available = available_path.read_text().split() if available_path.exists() else []
        governor = args.profile
        if args.profile == "balanced":
            governor = next((candidate for candidate in ("schedutil", "ondemand", "powersave") if candidate in available), "powersave")
        if available and governor not in available:
            print(f"[WARN] {governor} is unavailable for {path.parent.parent.name}; supported: {' '.join(available)}")
            continue
        print(f"{path}: {governor}")
        if not args.dry_run:
            path.write_text(governor + "\n", encoding="ascii")
    if args.dry_run:
        print("Dry run only; pass --apply to write governors.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
