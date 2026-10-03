# Fedora Atomic Host Notes

FortMox targets Proxmox VE as the hypervisor host. The CLI's read-only config and detection commands can be built and run on Fedora Atomic, but this does not make Fedora Atomic a supported deployment target for the Proxmox host setup.

```bash
go test ./...
go run ./cmd/fortmox config validate
go run ./cmd/fortmox detect
```

The former Fedora Atomic walkthrough used hard-coded paths and removed shell scripts. For the current CLI entry point, use [Getting Started](getting-started.md); for installing Proxmox, follow the upstream link in [Install Proxmox](install-proxmox.md).
