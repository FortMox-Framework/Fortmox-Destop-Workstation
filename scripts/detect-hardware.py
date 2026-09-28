#!/usr/bin/env python3
"""Read-only CPU, GPU, IOMMU and form-factor report."""

from fortmox_scripts import hardware_report


def main() -> int:
    """Print a best-effort report of the current Linux host."""
    result = hardware_report()
    print(f"Form factor: {result['form_factor']}")
    print(f"CPU: {result['cpu']} ({result['logical_cpus']} logical CPUs)")
    print(f"Virtualization: {result['virtualization']}")
    print(f"IOMMU groups: {result['iommu']}")
    print("GPUs:")
    for device in result["gpus"]:
        print(f"  {device}")
    if not result["gpus"]:
        print("  none detected (lspci may be unavailable)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
