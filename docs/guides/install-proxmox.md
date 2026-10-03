# Install Proxmox VE

FortMox does not install Proxmox VE. Use the current [official Proxmox VE download and installation documentation](https://www.proxmox.com/en/downloads) for version-specific requirements and installer steps; this repository does not pin an installer release or maintain verified screenshots and commands.

## Before installing

- Install only on hardware dedicated to the hypervisor. The installer can erase the selected target disk; confirm the disk carefully and back up any data you need.
- Check current Proxmox hardware and networking requirements in the official documentation. This guide's former device names, network defaults, and installer walkthrough are not reliable instructions for your hardware.
- Treat host installation and FortMox configuration as separate tasks. FortMox does not set up Proxmox networking, firewall rules, or VM isolation.

## After installation

From the FortMox repository root, build the CLI and inspect its read-only checks and plans:

```bash
go test ./...
go build -o fortmox ./cmd/fortmox
./fortmox config validate
./fortmox verify
./fortmox deploy --dry-run
./fortmox vm create --dry-run
```

`deploy` prints a plan and does not apply host changes. VM creation is dry-run by default; `./fortmox vm create --dry-run=false` invokes `qm create` for configured VMs, but does not install guest operating systems or configure networking and isolation. Only consider that opt-in command on the intended Proxmox host after reviewing the generated commands. See [Deploy FortMox](deploy-fortmox.md) for current CLI behavior and limitations.
