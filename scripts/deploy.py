#!/usr/bin/env python3
"""Show the host deployment plan. This command never mutates the host."""

import argparse

from fortmox_scripts import enabled_vms, load_config


def main() -> int:
    """Print the host deployment plan without applying system changes."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("-c", "--config", default="config/system.yaml")
    args = parser.parse_args()
    config, _ = load_config(args.config)
    steps = [
        "host requirement checks and hardware detection are not implemented by deploy",
        "kernel IOMMU and host hardening changes are not implemented",
        "encrypted storage setup is not implemented",
        "network isolation and firewall configuration are not implemented",
    ]
    if config.get("hardware", {}).get("gpu", {}).get("passthrough"):
        steps.append("GPU strategy planning only; no device changes are made")
    steps.extend(f"VM {name} is handled separately by deploy-vms.py" for name, _ in enabled_vms(config))
    for number, step in enumerate(steps, start=1):
        print(f"{number}. {step}")
    print("Dry run only; host deployment operations are not yet implemented by this Python command.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
