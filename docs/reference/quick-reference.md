# Quick Reference

## Build

```bash
go test ./...
go build -o fortmox ./cmd/fortmox
```

## Read-only Commands

```bash
./fortmox config validate
./fortmox config plan
./fortmox detect
./fortmox verify
./fortmox gpu
./fortmox power balanced
```

`verify` needs Proxmox utilities for VM and firewall runtime checks. Without them it still checks local configuration and template files and reports Proxmox checks as skipped.

## Dry-run Commands

```bash
./fortmox deploy --dry-run
./fortmox vm create --dry-run
```

Host deployment is currently plan-only. VM creation can be explicitly enabled on a Proxmox host after reviewing each generated command:

```bash
sudo ./fortmox vm create --dry-run=false
```

`fortmox power balanced --apply` writes CPU governors and requires root. GPU commands only plan a strategy; they do not bind devices or alter VM configuration.

## Configuration

- Main schema and VM enablement: `config/system.yaml`
- VM templates: `config/vm-templates/*.yaml`
- Python utility dependency: `requirements.txt` (PyYAML)
- Go command source: `cmd/fortmox/`

Use `fortmox config validate` after editing configuration. The current VM list is read from YAML; template VMIDs and resource values are not maintained in a separate config file.

## Legacy

Bash implementations have been removed. The current script entrypoints use Python 3; see [features and status](features.md) for their scope and limitations.
