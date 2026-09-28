# Deploy FortMox

This guide describes the current CLI boundary. FortMox is alpha and has not been certified on real Proxmox hosts. Back up host configuration and review every command before applying changes.

## Validate and inspect

From the repository root:

```bash
go test ./...
go build -o fortmox ./cmd/fortmox
./fortmox config validate
./fortmox config plan
./fortmox detect
./fortmox verify
```

`verify` is read-only. It checks configured VM templates, enabled VM IDs when `qm` is available, and Proxmox firewall state when configured. On non-Proxmox Linux, Proxmox-only checks are reported as skipped.

## Review deployment and VM commands

```bash
./fortmox deploy --dry-run
./fortmox vm create --dry-run
```

The `deploy` command is plan-only: kernel, storage, networking and host firewall changes are not applied by the Go CLI. VM creation reads enabled VM names and resource values from `config/system.yaml` and each referenced template. Dry-run is the default.

Only on the intended Proxmox host, after reviewing the output, run:

```bash
sudo ./fortmox vm create --dry-run=false
```

This invokes `qm create`; it does not configure every option in the VM templates, install guest operating systems, or set up OPNsense inside the guest. The generated VM configuration is an initial implementation, not a production deployment recipe.

## Python utilities

The Python 3 scripts under `scripts/` are helper utilities. Install their dependency with `python3 -m pip install -r requirements.txt`. `scripts/deploy.py` only prints a plan; `scripts/deploy-vms.py --apply` creates configured VMs; `scripts/power-management.py --apply` writes CPU governors. The other scripts are read-only. Bash script implementations have been removed; use the Python 3 entrypoints or Go CLI.

For host installation steps, see [Install Proxmox](install-proxmox.md). For known implementation boundaries, see [Features](../reference/features.md).
