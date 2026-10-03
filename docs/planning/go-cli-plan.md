> **Planning/status snapshot.** Current behavior is summarized here, but [features and status](../reference/features.md) is the source of truth for implementation support.

# Go CLI implementation status

FortMox uses one host-side Cobra CLI binary. An API service is not part of the current scope.

## Layout

```
cmd/fortmox/            main + cobra commands
internal/config/        typed system.yaml schema + validation
internal/hardware/      CPU, GPU, IOMMU, laptop vs desktop
internal/verify/        read-only local and best-effort Proxmox checks
internal/deploy/        host deployment planning (apply path not implemented)
internal/vm/            basic qm create arguments (dry-run by default)
internal/gpu/           strategy planning only (auto selects virtIO)
internal/power/         CPU governor planning and explicit apply
```

Commands: `fortmox config validate`, `config plan`, `detect`, `verify`, `deploy`, `vm create`, `gpu`, and `power`.

## Progress

- [x] Config schema, strict parsing, validation and dependency-ordered VM planning.
- [x] `detect` reports CPU, virtualization flags, IOMMU-group presence, PCI display devices, and estimated form factor.
- [x] `verify` checks configuration/templates and performs best-effort host/Proxmox checks; it does not test network isolation.
- [x] `vm create` reads enabled VMs and resource values from config/templates; dry-run is the default.
- [x] GPU strategy planning only; CPU governor profile application requires root.
- [ ] Host `deploy` apply path for kernel, storage and networking. `fortmox deploy` is plan-only.
- [ ] Full GPU device binding and per-VM multi-GPU assignment.
- [ ] VM creation does not yet consume every template field or install guest operating systems.

## Decisions

- CLI first, run directly on the host; no API service.
- Bash scripts have been removed; current script utilities in `scripts/` use Python 3.
- Cobra provides the CLI command tree.
- OPNsense has two switches in system.yaml (`vms.opnsense.enabled`, `microvm.opnsense_firewall.enabled`). It deploys if either is on, and validation warns when they disagree. Pick one place and remove the other when the schema is next revised.
