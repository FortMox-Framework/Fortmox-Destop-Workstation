#!/usr/bin/env python3
"""Validate and summarize the typed system configuration using PyYAML."""

import argparse

from fortmox_scripts import hardware_report, load_config


def main() -> int:
    """Load the selected config and print a concise hardware-aware summary."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("-c", "--config", default="config/system.yaml")
    args = parser.parse_args()
    config, _ = load_config(args.config)
    system = config.get("system", {})
    hardware = config.get("hardware", {})
    print(f"Config: {args.config}")
    print(f"Hostname: {system.get('hostname', 'unset')}")
    print(f"Configured type: {system.get('type', 'unset')}; detected type: {hardware_report()['form_factor']}")
    print(f"GPU passthrough: {hardware.get('gpu', {}).get('passthrough', False)} ({hardware.get('gpu', {}).get('method', 'unset')})")
    print(f"Configured VMs: {', '.join(name for name, _ in config.get('vms', {}).items())}")
    print("YAML loaded successfully; run `fortmox config validate` for schema validation.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
