#!/usr/bin/env python3
"""Show the host deployment plan. This command never mutates the host."""

import argparse

from fortmox_scripts import enabled_vms, load_config


def main() -> int:
    """Print the host deployment plan without applying system changes."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("-c", "--config", default="config/system.yaml")
    parser.add_argument("--dry-run", action="store_true", default=True, help="show planned changes (only supported mode)")
    args = parser.parse_args()
    config, _ = load_config(args.config)
    steps = [
        "validate host requirements and detect hardware",
        "configure kernel IOMMU and hardening settings",
        "prepare configured storage encryption vaults",
        "configure isolated virtual networking and firewall",
    ]
    if config.get("hardware", {}).get("gpu", {}).get("passthrough"):
        steps.append("prepare GPU passthrough")
    steps.extend(f"create configured VM {name}" for name, _ in enabled_vms(config))
    for number, step in enumerate(steps, start=1):
        print(f"{number}. {step}")
    print("Dry run only; host deployment operations are not yet implemented by this Python command.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
