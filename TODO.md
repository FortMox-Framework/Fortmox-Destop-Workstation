# TODO

This list tracks work needed before FortMox can be treated as a dependable Proxmox host manager. The current project is alpha; CI must remain read-only and dry-run only.

## Implemented in Current Source

These checks describe code that exists; they do not imply complete host setup or validation on physical Proxmox hardware.

- [x] `system.yaml` loading uses a typed Go schema and rejects unknown config fields.
- [x] Config validation checks core enums, hostname/timezone, DNS IPs, vault paths, VM dependencies, enabled-template presence, and duplicate enabled VMIDs.
- [x] VM planning produces deterministic dependency order and detects dependency cycles.
- [x] `fortmox detect` reports CPU model/count, virtualization flags, IOMMU-group presence, display devices, and an estimated laptop/desktop form factor.
- [x] `fortmox verify` checks local config/templates and performs best-effort host/Proxmox tool, storage, bridge, enabled-VM, and firewall service checks.
- [x] `fortmox vm create` builds a basic `qm create` command from VMID/name/cores/max-memory/disk plus a fixed default network; dry-run is the default.
- [x] `fortmox gpu` prints a strategy plan; `auto` currently selects virtIO and no GPU/driver changes are made.
- [x] `fortmox power` plans CPU governors and offers root-gated apply; other power controls are not implemented.
- [x] Python tests cover enabled-VM order, dependency cycles, all five template files, basic `qm` arguments, and input immutability.
- [x] GitHub Actions runs Go/Python checks and builds Linux amd64/arm64 artifacts for pull requests.

## Source and Schema Coverage

- [ ] Maintain a field-coverage table for every `Config` field and every VM-template key: validated, consumed by a command, deliberately manual, or unsupported. Do not count parsing alone as implementation.
- [ ] Make Python commands use the same strict schema/validation as Go, or remove duplicate config-sensitive Python entrypoints; current Python scripts use permissive `yaml.safe_load` and can act on configs Go rejects.
- [ ] Define and strictly validate the full VM-template schema; `LoadTemplate` currently reads only the VMID/name/CPU/memory/disk header and accepts unknown template keys.
- [ ] Expand validation tests for template relationships and fields not currently checked, including CPU sockets, memory bounds, disk settings, and each supported configuration section.

## Before Host Deployment

- [ ] Implement a guarded Go host-deploy path for kernel/IOMMU, storage, and network changes; keep dry-run as the default.
- [ ] Audit every typed config field for a real consumer; implement or remove inert settings such as hostname/timezone/debug, CPU core/IOMMU/scaling and storage tuning, security policies, DNS, additional-firewall instances/priority, advanced options, logging, backups, and development dry-run settings. Several are parsed or partially validated but never applied.
- [ ] Add host-specific configuration for bridges, storage targets, and device selection instead of relying on fixed `vmbr0` assumptions.
- [ ] Configure the OPNsense WAN/LAN topology, VM NIC attachment, forwarding, NAT, and isolation rules from config; current VM creation sends every VM to `vmbr0`.
- [ ] Add end-to-end network tests proving the Research VM is air-gapped and other VM traffic follows the configured firewall; service status alone does not prove traffic is filtered.
- [ ] Add tests around every host mutation using fake commands/filesystems; never run these tests against the CI runner's host.
- [ ] Add injectable host-probe/command interfaces and tests for hardware detection and verification; current hardware detection has no unit tests and verification invokes host tools directly.
- [ ] Add explicit confirmation, privilege checks, backups, and rollback guidance before applying host changes.
- [ ] Implement the documented host hardening controls: sysctl settings, module restrictions, AppArmor/SELinux, SSH policy, host firewall, audit logging, and integrity monitoring.
- [ ] Implement encrypted storage/vault creation, key handling, recovery, backup/restore, and safe failure behavior; current YAML settings do not create encrypted storage.
- [ ] Implement configured backup scheduling, retention, and restore verification; current backup YAML is not consumed.
- [ ] Implement configured logging/audit output and monitoring, or remove those fields and claims until support exists.
- [ ] Implement or remove claims about Proxmox service debloating; no debloat operation is currently provided.

## VM Management

- [x] Map the currently supported VM-template subset to `qm`: VMID, name, CPU cores, maximum memory, and disk storage/size.
- [ ] Map remaining template fields to `qm` arguments, including configured networks, firmware, GPU, and boot options; the current command uses a fixed default network.
- [ ] Apply or explicitly report unsupported guest settings: Looking Glass/ivshmem, SPICE, input devices, snapshots, autostart, shared storage, and security features.
- [ ] Cover all template fields not consumed by `CreateArgs`, including CPU sockets/type, memory min/ballooning, disk format/cache/type, display, template NIC/DNS settings, storage volumes/shares, boot/firmware, limits, snapshots, guest packages/users/firewall, and VM security features.
- [ ] Validate the Looking Glass host/guest workflow, ivshmem sizing, SPICE input/audio, and guest compatibility before describing it as supported.
- [ ] Define and implement guest OS installation/configuration, package setup, user creation, and guest hardening, or keep those as clearly manual procedures.
- [ ] Check for existing VMIDs and offer a clear skip/fail policy before invoking `qm create`.
- [ ] Require suitable privilege and an explicit confirmation for VM creation; preflight storage, bridges, VMID availability, and required Proxmox tools before invoking `qm`.
- [ ] Add injectable command-runner and CLI tests for dry-run output, `qm` failures, duplicate VMIDs, invalid/disabled VM names, and Proxmox-unavailable behavior.
- [ ] Implement or remove the configured OPNsense `priority` field; current deploy ordering follows dependencies and alphabetical ties instead.
- [ ] Verify firewall routing and isolation behavior, not just VM existence and service status.

## GPU and Power

- [ ] Implement hardware-aware GPU selection and per-VM assignment for full passthrough, SR-IOV, vGPU, and virtIO.
- [ ] Validate IOMMU groups, supported devices, and required host drivers before offering an apply operation.
- [x] Unit-test governor candidate selection for available profiles.
- [ ] Test sysfs profile planning, unsupported governors, write failures, and rollback against mocked files.
- [ ] Reject a requested governor when the kernel does not report supported governors, verify writes, and roll back earlier CPU-policy writes if a later write fails.
- [ ] Implement or remove claims for turbo control, GPU dynamic power, adaptive brightness, battery-aware profiles, and CPU idle-state management; only CPU frequency governors are currently applied.
- [ ] Implement the documented host Wayland compositor setup, or keep it explicitly outside the CLI's supported scope.

## Python Utilities

- [ ] Add Python tests for malformed YAML, bad template shapes/resources, missing templates, and `qm`/sysfs failures; the current five tests cover ordering, cycle detection, template presence, command values, and input immutability.
- [ ] Decide which Python utilities remain useful after their Go CLI equivalents are complete; remove duplicate entrypoints only after parity is verified.

## Release Readiness

- [ ] Audit all guides and architecture pages for unverified claims, remove “complete/fully supported” language, and distinguish intended design, manual instructions, and tested CLI behavior.
- [ ] Recheck performance, security, isolation, and compatibility numbers in the guides; retain only measurements backed by reproducible tests.
- [ ] Review every setup guide against the current CLI and mark unsupported procedures clearly.
- [ ] Document supported Proxmox, Linux, Go, and Python versions and test on a disposable Proxmox host.
- [ ] Require the Go checks, Python checks, and both build-matrix jobs in GitHub branch protection before merging pull requests.
- [ ] Define a license before publishing the repository publicly; no license has been selected yet.
- [ ] Add versioning and release packaging after the CLI behavior and config format stabilize.