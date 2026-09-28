# Go CLI implementation status

FortMox uses one host-side Cobra CLI binary. An API service is not part of the current scope.

## Layout

```
cmd/fortmox/            main + cobra commands
internal/config/        typed system.yaml schema + validation
internal/hardware/      CPU, GPU, IOMMU, laptop vs desktop
internal/verify/        checks ported from verify.sh
internal/deploy/        host deployment planning (apply path not implemented)
internal/vm/            templates to qm calls (VM list read from config)
internal/gpu/           passthrough, SR-IOV, vGPU, virtIO
internal/power/         power profiles
```

Commands: `fortmox detect`, `verify`, `deploy`, `vm create`, `gpu`.

## Progress

- [x] Config schema, strict parsing, validation and dependency-ordered VM planning.
- [x] `detect` reports CPU, virtualization, IOMMU, GPU and form factor.
- [x] `verify` checks configuration/templates, enabled VM presence and Proxmox firewall service state.
- [x] `vm create` reads enabled VMs and resource values from config/templates; dry-run is the default.
- [x] GPU strategy planning and CPU governor profile planning; applying a power profile requires root.
- [ ] Host `deploy` apply path for kernel, storage and networking. `fortmox deploy` is plan-only.
- [ ] Full GPU device binding and per-VM multi-GPU assignment.
- [ ] VM creation does not yet consume every template field or install guest operating systems.

## Decisions

- CLI first, run directly on the host; no API service.
- Bash scripts have been removed; current script utilities in `scripts/` use Python 3.
- Cobra provides the CLI command tree.
- OPNsense has two switches in system.yaml (`vms.opnsense.enabled`, `microvm.opnsense_firewall.enabled`). It deploys if either is on, and validation warns when they disagree. Pick one place and remove the other when the schema is next revised.
