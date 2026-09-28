# Features and Status

FortMox is alpha. Configuration and planning are usable for development and review, but the project is not production-ready and has not been validated on real Proxmox hardware.

| Area | Current behavior | Boundary |
|------|------------------|----------|
| Config | Typed YAML schema, strict unknown-key checks, validation and dependency ordering | Does not certify hardware compatibility |
| Detect | CPU model/count, virtualization flags, IOMMU groups, GPU listing and laptop/desktop estimate | Depends on Linux sysfs/procfs and `lspci` |
| Verify | All configured templates, enabled VM IDs, Proxmox VM listing and firewall service checks | VM/firewall runtime checks require Proxmox tools and configured services |
| Deploy | Config-driven host change plan | Go and Python deploy commands are dry-run/plan-only; no host setup is applied |
| VM | Builds `qm create` commands from enabled config entries and templates | Does not apply every template field or install guests; actual create is opt-in |
| GPU | Plans virtIO, SR-IOV, vGPU or full passthrough strategy | Does not configure host drivers or per-VM assignment yet |
| Power | Plans CPU governors; Go CLI supports explicit root-only apply | Runtime support depends on CPU driver and available governors |
| Python scripts | Config-driven detect, verify, VM planning/creation and power helpers | PyYAML required; deployment and GPU utilities are informational only |

## Configuration Files

The shipped configuration consists of `config/system.yaml` and five files in `config/vm-templates/`: `opnsense-microvm.yaml`, `clean-vm.yaml`, `gaming-vm.yaml`, `research-vm.yaml`, and `tools-vm.yaml`. VM enablement and template paths are read from the system config; there is no separate hardcoded VM list.

## Commands

```bash
fortmox config validate
fortmox config plan
fortmox detect
fortmox verify
fortmox deploy --dry-run
fortmox vm create --dry-run
fortmox gpu --method auto
fortmox power balanced
```

See the [quick reference](quick-reference.md) for build commands and opt-in operations. Historical feature promises are preserved in the [archive](../archive/features-legacy.md), not as current implementation claims.
