#!/usr/bin/env python3
"""Report configured GPU passthrough; VFIO changes require an explicit plan review."""

import argparse

from fortmox_scripts import hardware_report, load_config


def main() -> int:
    """Report configured and detected GPUs without changing drivers or VMs."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("-c", "--config", default="config/system.yaml")
    args = parser.parse_args()
    config, _ = load_config(args.config)
    settings = config.get("hardware", {}).get("gpu", {})
    print(f"Passthrough enabled: {settings.get('passthrough', False)}")
    print(f"Method: {settings.get('method', 'auto')}")
    print(f"Configured devices: {', '.join(settings.get('devices', [])) or 'auto-detect'}")
    print("Detected devices:")
    for device in hardware_report()["gpus"]:
        print(f"  {device}")
    print("No driver or host configuration was changed.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
