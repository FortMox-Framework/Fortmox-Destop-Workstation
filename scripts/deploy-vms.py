#!/usr/bin/env python3
"""Create enabled configured VMs; dry-run is the default."""

import argparse
import shutil
import subprocess

from fortmox_scripts import enabled_vms, load_config, load_template, qm_create_args


def main() -> int:
    """Create enabled VMs, or only the requested enabled VM, when explicitly applied."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("-c", "--config", default="config/system.yaml")
    parser.add_argument("--dry-run", action="store_true", default=True, help="print qm commands only (default)")
    parser.add_argument("--apply", action="store_false", dest="dry_run", help="create the VMs on this Proxmox host")
    parser.add_argument("vm", nargs="?", help="create only this enabled VM")
    args = parser.parse_args()
    config, root = load_config(args.config)
    selected = enabled_vms(config)
    if args.vm:
        selected = [(name, entry) for name, entry in selected if name == args.vm]
        if not selected:
            parser.error(f"VM {args.vm!r} is not enabled in {args.config}")
    if not args.dry_run and not shutil.which("qm"):
        parser.error("qm is required for --apply; run on a Proxmox host")

    for name, entry in selected:
        path = entry.get("template_file") or config.get("microvm", {}).get("opnsense_firewall", {}).get("template_file")
        command = ["qm", *qm_create_args(name, load_template(root, path))]
        print("$ " + " ".join(command))
        if not args.dry_run:
            if subprocess.run(command, check=False).returncode:
                return 1
    if args.dry_run:
        print("Dry run only; no VMs created. Pass --apply to execute on Proxmox.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
