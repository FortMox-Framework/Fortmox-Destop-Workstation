# Fedora Atomic Host Notes

FortMox targets Proxmox VE as the hypervisor host. The CLI's read-only config and detection commands can be built and run on Fedora Atomic, but this does not make Fedora Atomic a supported deployment target for the Proxmox host setup.

```bash
go test ./...
go run ./cmd/fortmox config validate
go run ./cmd/fortmox detect
```

The old walkthrough used hard-coded workstation paths and shell scripts that have since been removed. It is retained as [historical material](../archive/fedora-atomic-host-legacy.md). For a current install path, use [Getting Started](getting-started.md).
