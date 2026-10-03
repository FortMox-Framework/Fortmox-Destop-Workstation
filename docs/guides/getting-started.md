# Getting Started

FortMox is an early-stage Go CLI and configuration for a Proxmox workstation. It has not been certified on production hosts. Start with read-only checks and review generated commands before allowing system changes.

## Build and validate

From the repository root, using Go 1.22 or later:

```bash
go test ./...
go build -o fortmox ./cmd/fortmox
./fortmox config validate
./fortmox config plan
```

Edit `config/system.yaml` and the matching files under `config/vm-templates/` for the hardware and VMs you actually intend to use. The typed schema validates names, dependencies and template paths; it does not guarantee that a host can safely apply every setting.

## Inspect the host

```bash
./fortmox detect
./fortmox verify
./fortmox gpu
./fortmox power balanced
```

These commands do not mutate host settings. Proxmox VM and firewall checks are skipped or warned when the commands are run on a non-Proxmox Linux system.

## Review changes

```bash
./fortmox deploy --dry-run
./fortmox vm create --dry-run
```

`deploy` currently prints a plan only. VM creation is dry-run by default; `./fortmox vm create --dry-run=false` invokes `qm create` and must only be used on the intended Proxmox host after reviewing the output. The generated VM commands are a starting implementation and do not apply every setting in each template.

For Proxmox installation guidance, see [install-proxmox](install-proxmox.md). For current command coverage and limitations, see [features](../reference/features.md) and [quick reference](../reference/quick-reference.md).
