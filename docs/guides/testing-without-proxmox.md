# Testing Without Proxmox

The Go CLI's config validation and hardware detection can run on a regular Linux machine. Proxmox VM and firewall checks require `qm` and `pve-firewall` and will be reported as skipped when unavailable.

```bash
go test ./...
go run ./cmd/fortmox config validate
go run ./cmd/fortmox config plan
go run ./cmd/fortmox detect
go run ./cmd/fortmox verify
```

All commands above are read-only. Deployment and VM dry runs can also be inspected without applying changes:

```bash
go run ./cmd/fortmox deploy --dry-run
go run ./cmd/fortmox vm create --dry-run
```

Do not use `fortmox vm create --dry-run=false` except on the target Proxmox host after reviewing the generated commands. The project has not been certified on physical hardware; use a disposable test host before relying on any security or isolation behavior.
