# Architecture & Design Philosophy

> **Target design, not current implementation.** The CLI does not yet apply the hardening, encryption, Wayland, network-isolation, or full GPU features shown below. See [features and status](features.md) for implemented behavior.

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│          Bare Metal Hardware (Laptop or Desktop)             │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ Proxmox Host (Debloated, Hardened, Lightweight)          ││
│  │                                                          ││
│  │  ┌────────────────────────────────────────────────────┐ ││
│  │  │ Wayland Compositor (Hikari/Woodland/dwl)           │ ││
│  │  │ • Stacking support for easy window management      │ ││
│  │  │ • Lightweight (minimal CPU/battery impact)         │ ││
│  │  │ • Direct input device access                       │ ││
│  │  └────────────────────────────────────────────────────┘ ││
│  │                                                          ││
│  │  ┌────────────────────────────────────────────────────┐ ││
│  │  │ Hypervisor Kernel (Hardened)                        │ ││
│  │  │ • IOMMU/VT-d enabled for device isolation          │ ││
│  │  │ • AppArmor/SELinux for process isolation           │ ││
│  │  │ • Hardware-accelerated encryption (AES-NI)         │ ││
│  │  │ • Power management (Intel P-State / AMD P-State)   │ ││
│  │  │ • Debloated services (only essentials)             │ ││
│  │  └────────────────────────────────────────────────────┘ ││
│  │                                                          ││
│  │  ┌────────────────────────────────────────────────────┐ ││
│  │  │ Storage Layer (Encrypted & Performant)              │ ││
│  │  │ • Hardware-accelerated encryption (no perf loss)   │ ││
│  │  │ • TRIM/Discard for SSD efficiency                  │ ││
│  │  │ • Separate vaults for different security zones     │ ││
│  │  └────────────────────────────────────────────────────┘ ││
│  │                                                          ││
│  │  ┌────────────────────────────────────────────────────┐ ││
│  │  │ Network Stack (Isolated & Monitored)                │ ││
│  │  │ • Virtual bridges per compartment                  │ ││
│  │  │ • Hardware-level packet filtering (no CPU impact)  │ ││
│  │  │ • Optional air-gapping for Research VM             │ ││
│  │  └────────────────────────────────────────────────────┘ ││
│  └─────────────────────────────────────────────────────────┘│
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ Compartmentalized VMs (Separate Security Zones)          ││
│  │                                                          ││
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ││
│  │  │  Clean VM    │  │  Gaming VM   │  │ Research VM  │  ││
│  │  │              │  │              │  │              │  ││
│  │  │ • General    │  │ • GPU Full   │  │ • Isolated   │  ││
│  │  │   computing  │  │   Passthrough│  │ • Air-gapped │  ││
│  │  │ • Browsing   │  │ • Full perf  │  │ • Analysis   │  ││
│  │  │ • Work tasks │  │ • Gaming     │  │ • Testing    │  ││
│  │  │ • Normal net │  │ • Normal net │  │ • No network │  ││
│  │  └──────────────┘  └──────────────┘  └──────────────┘  ││
│  │         │                  │                 │          ││
│  │         └──────────────────┴─────────────────┘          ││
│  │                    (Shared Storage)                      ││
│  │                                                          ││
│  │  ┌──────────────────────────────────────────────────┐  ││
│  │  │ Tools VM (Optional)                              │  ││
│  │  │ • Network analysis and testing tools             │  ││
│  │  │ • Isolated from other VMs                        │  ││
│  │  └──────────────────────────────────────────────────┘  ││
│  │                                                          ││
│  └─────────────────────────────────────────────────────────┘│
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Universal Laptop/Desktop Design

### Automatic Detection
The current read-only `fortmox detect` command reports:
- **Form Factor**: Laptop (battery detected) vs Desktop (AC only)
- **CPU**: Intel vs AMD (affects IOMMU naming)
- **GPU**: Available GPUs and passthrough capabilities
- **Storage**: SSD vs HDD (affects encryption overhead)

### Laptop-Specific Optimizations
```yaml
laptop_optimizations:
  cpu:
    - intel_pstate or amd_pstate: enable dynamic frequency
    - turbo_boost: enable for performance when needed
    - idle_states: use deeper C-states for battery life

  gpu:
    - dynamic_power_management: scale down when idle
    - clock_gating: disable unused GPU units

  storage:
    - write_cache: intelligent buffering
    - TRIM_schedule: hourly for SSD health

  display:
    - backlight: adaptive brightness
    - refresh_rate: dynamic based on content
```

### Desktop Full Performance
```yaml
desktop_optimizations:
  cpu:
    - performance_mode: all cores at max frequency
    - turbo_boost: always enabled
    - idle_states: minimal C-states for latency

  gpu:
    - maximum_clocks: sustained max performance
    - no_power_gating: disable for consistency

  storage:
    - high_performance_mode: maximize throughput

  display:
    - constant_refresh_rate: optimized for monitor
```

## Compartmentalization Strategy

### Security Zones

**Zone 1: Clean VM**
- Purpose: Daily work, browsing, general computing
- Network: Normal internet access
- GPU: virtIO GPU (low perf, isolated)
- Threat Level: Low
- Compromise Risk: Contained to this VM

**Zone 2: Gaming VM**
- Purpose: Gaming and GPU-intensive workloads
- Network: Normal internet access
- GPU: Full passthrough (maximum performance)
- Threat Level: Medium (untrusted games/software)
- Compromise Risk: Isolated to this VM, cannot affect other VMs

**Zone 3: Research VM**
- Purpose: Malware analysis, security research
- Network: Air-gapped (no network)
- GPU: virtIO GPU (no direct access)
- Threat Level: Critical (assumed malicious)
- Compromise Risk: Maximum isolation, no network escape possible

**Zone 4: Tools VM (Optional)**
- Purpose: Network analysis, penetration testing
- Network: Isolated network segment
- GPU: virtIO GPU
- Threat Level: Medium (testing tools)
- Compromise Risk: Isolated network, no cross-VM access

### Data Flow

```
Internet
  │
  ├──→ Clean VM (normal access)
  │     └──→ Can access shared storage (read-only for sensitive)
  │
  ├──→ Gaming VM (normal access)
  │     └──→ Isolated from Research VM
  │
  └──→ Tools VM (isolated network)
        └──→ No direct access to other VMs

Research VM (air-gapped)
  └──→ No network access
  └──→ No access to other VMs
  └──→ Isolated storage
```

## GPU Passthrough Strategy

All GPU methods supported with intelligent fallback:

### Method 1: virtIO GPU (Most Compatible)
- Virtual GPU rendered by host
- Lower performance, but universal compatibility
- No host GPU needed
- Use Case: Clean VM, Research VM (safe)

### Method 2: SR-IOV Passthrough
- GPU shares physical interface with host
- Medium performance, good isolation
- Requires GPU hardware support
- Use Case: Gaming VM (when available)

### Method 3: vGPU Passthrough (NVIDIA)
- Hardware-level GPU virtualization
- Multiple VMs share single GPU
- Good performance and isolation
- Use Case: Gaming VM (NVIDIA cards)

### Method 4: Full GPU Passthrough
- Exclusive GPU device to single VM
- Maximum performance (native gaming)
- One GPU per VM limitation
- Use Case: Gaming VM (preferred when possible)

### Automatic Method Selection
```
if gpu_count >= 2:
  method = full_passthrough  # Each VM gets own GPU
elif gpu_supports_vgpu:
  method = vgpu              # NVIDIA cards
elif gpu_supports_sr_iov:
  method = sr_iov            # AMD/Intel SR-IOV capable
else:
  method = virtio            # Universal fallback
```

## Security Hardening Layers

### Layer 1: Host Hypervisor
- Minimal kernel modules (debloat)
- AppArmor confinement
- IOMMU/VT-d for device isolation
- Secure boot support
- SELinux if available

### Layer 2: VM Isolation
- vCPU pinning to prevent side-channels
- Memory isolation (no ballooning)
- vIOMMU for nested isolation
- No shared virtual devices

### Layer 3: Network Isolation
- Virtual bridges per compartment
- Hardware-level packet filtering
- IP route isolation
- Optional air-gapping for Research VM

### Layer 4: Storage Encryption
- Hardware-accelerated (AES-NI)
- LUKS2 for encrypted volumes
- Separate keys per compartment
- Key derivation from system entropy

## Battery vs Performance Balance

### No Compromise Design
- Battery optimizations use **hardware features** (P-States, ACPI)
- Performance uses **kernel scheduling** (not disabled features)
- GPU power management is **automatic** (clocks scale as needed)
- Storage I/O is **queued intelligently** (not throttled)

**Result**: Same security and core functionality on both laptop and desktop, with intelligent power scaling on laptops that doesn't impact performance when needed.

## Declarative Configuration

All system settings are in `config/system.yaml`:

```yaml
system:
  type: auto              # auto-detect or specify: laptop, desktop
  hostname: proxmox-ws
  timezone: UTC

hardware:
  cpu:
    cores: auto           # auto-detect
    iommu: true
  gpu:
    passthrough: true
    method: auto          # auto, virtio, sr-iov, vgpu, full
  storage:
    encryption: true
    compression: false

wayland:
  compositor: auto        # auto, hikari, woodland, dwl
  stacking: true

power:
  profile: auto           # auto, balanced, performance

security:
  hardening_level: paranoid
  selinux: true
  apparmor: true

vms:
  templates: [clean, gaming, research, tools]
```

This YAML is processed by the config management layer and translated into concrete system configurations.

## Performance Characteristics

### Expected Performance Impact

**Hypervisor Overhead**: <3% CPU
- Proxmox minimal (debloated)
- Kernel hardening: negligible with hardware support
- Encryption: 0% with AES-NI (hardware accelerated)

**GPU Performance**: Near-native
- Full passthrough: 98-100% native performance
- SR-IOV: 95-98% native performance
- vGPU: 90-95% native performance
- virtIO: 60-80% (for low-perf needs)

**Battery Life**: 90%+ of bare metal
- Power management: hardware-managed
- Idle efficiency: equivalent to native
- Active performance: same as native when needed

---

See `docs/` directory for detailed configuration guides.
