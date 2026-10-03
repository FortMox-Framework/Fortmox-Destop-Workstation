#!/usr/bin/env python3
"""Read-only Proxmox host, template, VM and firewall checks."""

import argparse
import shutil
import subprocess

from fortmox_scripts import all_templates, enabled_vms, load_config, load_template


def main() -> int:
    """Run local and Proxmox checks without modifying the host."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("-c", "--config", default="config/system.yaml")
    args = parser.parse_args()
    config, root = load_config(args.config)
    failures = 0

    def report(ok: bool | None, message: str) -> None:
        nonlocal failures
        # None means this check is unavailable or advisory, not a failure.
        label = "PASS" if ok else "WARN" if ok is None else "FAIL"
        failures += int(ok is False)
        print(f"[{label}] {message}")

    for name, path in all_templates(config):
        try:
            load_template(root, path)
            report(True, f"template exists and parses: {name} ({path})")
        except (OSError, ValueError) as error:
            report(False, f"template invalid for {name}: {error}")
    if not shutil.which("qm"):
        report(None, "qm unavailable; VM existence and status checks skipped")
    else:
        listing = subprocess.run(["qm", "list"], text=True, capture_output=True, check=False)
        if listing.returncode:
            report(False, f"qm list failed: {listing.stderr.strip()}")
        else:
            existing = {line.split()[0] for line in listing.stdout.splitlines()[1:] if line.split() and line.split()[0].isdigit()}
            for name, entry in enabled_vms(config):
                template = load_template(root, entry["template_file"])
                vmid = str(template["vm"].get("vmid", ""))
                report(vmid in existing, f"enabled VM {name} (vmid {vmid}) {'exists' if vmid in existing else 'is missing'}")

    firewall = config.get("microvm", {}).get("opnsense_firewall", {})
    if config.get("microvm", {}).get("enabled") and firewall.get("enabled"):
        vmid = str(firewall.get("vmid", 0))
        if shutil.which("qm"):
            status = subprocess.run(["qm", "status", vmid], text=True, capture_output=True, check=False)
            report(status.returncode == 0 and "running" in status.stdout, f"firewall VM {vmid} {'is running' if status.returncode == 0 and 'running' in status.stdout else 'is not running'}")
        else:
            report(None, "firewall VM status skipped because qm is unavailable")
        if shutil.which("pve-firewall"):
            status = subprocess.run(["pve-firewall", "status"], text=True, capture_output=True, check=False)
            report(status.returncode == 0 and "running" in status.stdout.lower(), f"Proxmox firewall service: {status.stdout.strip() or status.stderr.strip()}")
        else:
            report(None, "pve-firewall unavailable; service check skipped")
    return 1 if failures else 0


if __name__ == "__main__":
    raise SystemExit(main())
