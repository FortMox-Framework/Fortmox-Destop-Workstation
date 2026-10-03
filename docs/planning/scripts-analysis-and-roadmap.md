> **Planning document.** Describes designs and roadmaps, not necessarily what is implemented. See [docs/README.md](../README.md) for current docs.

> **Historical planning analysis.** This document describes removed shell scripts and proposed work; the commands and readiness claims below are not current. See [features and status](../reference/features.md).

# FortMox Scripts - Complete Analysis & Roadmap
## What's Done, What's Needed, and Feature Suggestions

---

## Part 1: Current Scripts Status

### Overview
```
Total Scripts: 7
Total Lines: 2,228
Status: Implemented for basic deployment (not yet tested on hardware)
Quality: Professional (error handling, logging, color output)
```

---

## Script Breakdown: Status & Details

### ✅ SCRIPT 1: `deploy.sh` (389 lines)
**Status:** ✅ IMPLEMENTED (not yet tested on hardware)

**What it does:**
- Main orchestrator script
- Checks root privileges
- Validates Proxmox installation
- Loads system configuration (YAML)
- Calls config-manager.sh for hardware detection
- Configures kernel hardening
- Enables IOMMU (Intel/AMD)
- Configures Wayland display
- Sets up encrypted storage vaults
- Configures network isolation
- Installs required packages
- Runs verification suite
- Provides detailed logging

**Key Features:**
- Error handling with `set -euo pipefail`
- Color-coded output (INFO, SUCCESS, WARN, ERROR)
- Section-based logging
- Dependency checking
- Configuration file validation
- Modular design (calls other scripts)

**Ready to use:** ✅ YES

---

### ✅ SCRIPT 2: `config-manager.sh` (309 lines)
**Status:** ✅ IMPLEMENTED (not yet tested on hardware)

**What it does:**
- YAML configuration parser
- Hardware detection:
  - CPU type (Intel vs AMD)
  - GPU type (NVIDIA, AMD, Intel)
  - System type (Laptop vs Desktop)
  - RAM amount
  - Storage info
- IOMMU detection
- SR-IOV capability detection
- Applies system configuration:
  - Hostname setup
  - Timezone configuration
  - CPU frequency scaling
  - Power management profiles
- Reads from `config/system.yaml`
- Sets up environment variables

**Key Features:**
- YAML value extraction
- Hardware auto-detection
- Intelligent profiling (laptop vs desktop)
- Extensible for future additions

**Ready to use:** ✅ YES

---

### ✅ SCRIPT 3: `deploy-vms.sh` (327 lines)
**Status:** ✅ IMPLEMENTED (not yet tested on hardware)

**What it does:**
- Creates VMs from templates
- Validates configuration
- Parses VM list from config
- Creates OPNsense MicroVM (vmid: 50)
- Creates Clean VM (vmid: 100)
- Creates Gaming VM (vmid: 101) with Looking Glass config
- Creates Research VM (vmid: 102)
- Creates Tools VM (vmid: 103) - optional
- Sets VM properties:
  - CPU/Memory allocation
  - Disk size and storage
  - Network configuration
  - GPU configuration
  - Display settings
- Error handling and validation

**Key Features:**
- Uses `qm` (Proxmox qemu tool)
- Dependency on OPNsense (routes traffic)
- Configurability via YAML
- Support for enable/disable per VM

**Ready to use:** ✅ YES
**Note:** OPNsense MicroVM must deploy first

---

### ✅ SCRIPT 4: `verify.sh` (407 lines)
**Status:** ✅ IMPLEMENTED (not yet tested on hardware)

**What it does:**
- Comprehensive verification suite
- Checks Proxmox installation
- Verifies hardware configuration
- Checks kernel hardening settings
- Validates IOMMU configuration
- Verifies GPU passthrough setup
- Checks network isolation
- Validates storage encryption
- Tests firewall configuration
- Reports statistics:
  - ✓ Passed checks
  - ✗ Failed checks
  - ⚠ Warnings
- Provides detailed report

**Key Features:**
- 40+ individual verification checks
- Detailed logging
- Pass/fail/warning tracking
- Helpful error messages
- Actionable recommendations

**Ready to use:** ✅ YES

---

### ✅ SCRIPT 5: `detect-hardware.sh` (272 lines)
**Status:** ✅ IMPLEMENTED (not yet tested on hardware)

**What it does:**
- Detailed hardware detection
- CPU information
- GPU detection (type, vendor, capabilities)
- RAM size
- Storage devices
- Network interfaces
- Form factor detection (laptop vs desktop)
- IOMMU support detection
- SR-IOV capability
- Provides recommendations:
  - VM CPU allocation
  - RAM allocation
  - GPU configuration suggestions
  - Power profile recommendations

**Key Features:**
- Non-invasive (read-only)
- Detailed output
- Actionable recommendations
- Can be run anytime

**Ready to use:** ✅ YES

---

### ✅ SCRIPT 6: `gpu-setup.sh` (236 lines)
**Status:** ✅ COMPLETE BUT LIMITED

**What it does:**
- GPU detection
- GPU type identification
- IOMMU support checking
- SR-IOV capability detection
- VFIO configuration
- GRUB parameter updates
- modprobe blacklist configuration
- Kernel initramfs updates
- Verification checks

**Limitations (from our analysis):**
- ❌ No interactive menu (just does full passthrough)
- ❌ Single GPU only (ignores others)
- ❌ No user choice of methods
- ❌ No per-VM assignment
- ❌ Manual switching required (reboot needed)

**Ready to use:** ✅ YES (for basic setup)
**Needs enhancement:** 🔜 Python tool (we designed this!)

---

### ✅ SCRIPT 7: `power-management.sh` (288 lines)
**Status:** ✅ IMPLEMENTED (not yet tested on hardware)

**What it does:**
- Power profile detection
- Laptop vs Desktop optimization
- CPU frequency scaling
- Turbo boost management
- GPU power management
- Thermal management
- Battery optimization (laptop)
- Display brightness (laptop)
- C-state configuration
- Power state monitoring

**Key Features:**
- Hardware-based (no feature disabling)
- Form factor aware
- Thermal aware
- Battery aware (laptop)
- Real-time monitoring

**Ready to use:** ✅ YES

---

## Part 2: What's Needed vs What Exists

### Current Gaps

| Feature | Current | Needed | Priority |
|---------|---------|--------|----------|
| **GPU Method Selection** | ❌ No menu | ✅ TUI Menu | 🔴 High |
| **Multi-GPU Support** | ❌ Single GPU | ✅ Multi-GPU | 🔴 High |
| **Mode Switching** | Manual reboot | Runtime switching | 🟡 Medium |
| **Per-VM Assignment** | Config only | Interactive GUI | 🟡 Medium |
| **Performance Monitoring** | Logs only | Real-time dashboard | 🟠 Low |
| **Auto-Profiles** | Manual | Automatic selection | 🟠 Low |
| **Prometheus Export** | None | Metrics export | 🟠 Low |
| **Web Dashboard** | None | Browser-based | 🟠 Low |

### Scripts That Need Enhancements

```
deploy.sh          → Add more features (optional)
config-manager.sh  → Add more config options
deploy-vms.sh      → Add VM customization options
verify.sh          → Add more detailed checks
detect-hardware.sh → Add GPU details
gpu-setup.sh       → REPLACE with Python version ⭐
power-management.sh → Add battery detection
```

---

## Part 3: Feature Suggestions for FortMox

### Tier 1: High Priority (Next Month)

#### 1. **Backup & Restore System**
```
What it does:
├─ Automatic VM snapshots on schedule
├─ Full system backup capability
├─ Selective VM backup
├─ Restore points with versioning
├─ Backup verification
└─ Incremental backups (space efficient)

Implementation:
├─ New script: backup-manager.sh
├─ Cron-based scheduling
├─ Configuration in system.yaml
└─ Integration with verify.sh

Benefit:
├─ Safe malware testing (snapshot backup)
├─ System recovery capability
├─ Data protection
└─ Peace of mind
```

#### 2. **VM Cloning & Templates**
```
What it does:
├─ Clone existing VMs
├─ Create custom templates
├─ Quick VM provisioning
├─ Template management
└─ Rapid deployment

Implementation:
├─ New script: vm-clone.sh
├─ Template storage system
├─ Metadata preservation
└─ Disk optimization

Benefit:
├─ Fast new VM creation
├─ Consistent configurations
├─ Faster testing cycles
└─ Reduced setup time
```

#### 3. **Network Management Dashboard**
```
What it does:
├─ VM network status
├─ Inter-VM connectivity testing
├─ Firewall rule verification
├─ DNS configuration management
├─ Network isolation verification
└─ Traffic monitoring (basic)

Implementation:
├─ New script: network-check.sh
├─ Integration with OPNsense
├─ Network diagnostics
└─ Visualization

Benefit:
├─ Easy network troubleshooting
├─ Verify isolation
├─ Configuration validation
└─ Performance optimization
```

#### 4. **VM Resource Monitoring**
```
What it does:
├─ Real-time CPU usage per VM
├─ Memory usage tracking
├─ Disk I/O monitoring
├─ Network I/O monitoring
├─ GPU utilization per VM
└─ Thermal monitoring

Implementation:
├─ New script: monitor-vms.sh
├─ Polling-based monitoring
├─ Alert thresholds
├─ Historical data
└─ Simple GUI viewer

Benefit:
├─ Identify bottlenecks
├─ Optimize resource allocation
├─ Detect problems early
└─ Better performance tuning
```

### Tier 2: Medium Priority (2-3 Months)

#### 5. **VM Orchestration & Auto-Start**
```
What it does:
├─ Intelligent VM boot order
├─ Dependency management (OPNsense first)
├─ Scheduled VM startup/shutdown
├─ Boot sequence optimization
├─ Status monitoring
└─ Auto-recovery

Implementation:
├─ New script: vm-orchestration.sh
├─ Systemd integration
├─ Dependency graph
└─ Health checks

Benefit:
├─ Automatic system recovery
├─ Proper boot sequence
├─ Power management
└─ High availability
```

#### 6. **SSL Certificate Management**
```
What it does:
├─ Auto-generate certificates
├─ Certificate renewal (Let's Encrypt)
├─ Certificate distribution to VMs
├─ Validation and checks
└─ Expiration warnings

Implementation:
├─ New script: cert-manager.sh
├─ Certbot integration
├─ Automatic renewal
└─ Multi-VM support

Benefit:
├─ Secure Proxmox access
├─ HTTPS everywhere
├─ No manual certificate management
└─ Compliance ready
```

#### 7. **Automated Security Updates**
```
What it does:
├─ Automatic kernel updates (host)
├─ VM guest OS updates
├─ Security patch application
├─ Scheduled updates
├─ Update rollback capability
└─ Update verification

Implementation:
├─ New script: update-manager.sh
├─ Unattended upgrades
├─ Update scheduling
├─ Verification checks
└─ Notification system

Benefit:
├─ Always up-to-date
├─ Reduced attack surface
├─ Automated maintenance
└─ Safe update process
```

#### 8. **Encrypted Backup with Cloud Sync**
```
What it does:
├─ Local encrypted backups
├─ Cloud backup (AWS, Backblaze, etc)
├─ End-to-end encryption
├─ Selective backup
├─ Cross-site disaster recovery
└─ Backup verification

Implementation:
├─ New script: backup-sync.sh
├─ Encryption (AES-256)
├─ Cloud provider integration
├─ Incremental sync
└─ Verification

Benefit:
├─ Off-site backups
├─ Disaster recovery
├─ Data protection
└─ Compliance ready
```

### Tier 3: Nice-to-Have (Later)

#### 9. **Performance Benchmarking**
```
What it does:
├─ GPU benchmark (gaming performance)
├─ CPU benchmark
├─ Disk I/O benchmark
├─ Network performance test
├─ VM vs native comparison
└─ Performance regression detection

Implementation:
├─ New script: benchmark.sh
├─ Standard benchmark tools
├─ Result comparison
└─ Historical tracking
```

#### 10. **Intrusion Detection (IDS)**
```
What it does:
├─ Network-based IDS (Suricata)
├─ Host-based IDS (AIDE)
├─ Anomaly detection
├─ Alert generation
└─ Log analysis
```

#### 11. **VPN Access**
```
What it does:
├─ WireGuard VPN setup
├─ Secure remote access
├─ VM access from outside
├─ Encryption everywhere
└─ Peer management
```

#### 12. **Kubernetes Integration**
```
What it does:
├─ K8s cluster on VMs
├─ Container orchestration
├─ Load balancing
├─ Service mesh (optional)
└─ Advanced workload management
```

---

## Part 4: Recommended Script Enhancements

### Priority 1: Create New Scripts (High Value)

#### New Script: `backup-manager.sh`
```bash
Purpose: Automated backup management
Size: ~250 lines
Features:
├─ Schedule VM backups
├─ Snapshot management
├─ Retention policies
├─ Restore capability
└─ Verification

Integration:
├─ Called from deploy.sh (optional)
├─ Standalone usage
└─ Cron-based automation
```

#### New Script: `monitor-vms.sh`
```bash
Purpose: Real-time VM monitoring
Size: ~300 lines
Features:
├─ CPU, memory, disk, network
├─ GPU utilization
├─ Alert thresholds
└─ Data export

Integration:
├─ Can run continuously
├─ TUI or simple output
└─ Log to file
```

#### New Script: `network-check.sh`
```bash
Purpose: Network validation
Size: ~200 lines
Features:
├─ VM connectivity tests
├─ Firewall rule verification
├─ OPNsense status
├─ Route verification
└─ Isolation confirmation

Integration:
├─ Part of verify.sh
├─ Standalone testing
└─ Troubleshooting aid
```

### Priority 2: Enhance Existing Scripts

#### Enhance: `gpu-setup.sh` → Deprecate, Replace with Python

**Action:** Don't modify bash version, replace with Python tool
- Keep as fallback (legacy support)
- Build Python version (we designed it!)
- TUI menu interface
- Multi-GPU support
- Per-VM assignment

#### Enhance: `deploy-vms.sh`

**Suggestions:**
```bash
# Add features:
1. VM customization options
2. Custom vCPU/RAM per user
3. Advanced storage options
4. Custom network configuration
5. Guest OS selection
6. Post-deploy scripts
7. Automated snapshots
8. Resource limits configuration
```

#### Enhance: `verify.sh`

**Suggestions:**
```bash
# Add checks for:
1. GPU passthrough functionality
2. Looking Glass configuration
3. VM resource allocation
4. Network isolation effectiveness
5. Firewall rules (OPNsense)
6. Backup status
7. Certificate validity
8. Performance baselines
9. Security patch status
10. Encryption verification
```

#### Enhance: `power-management.sh`

**Suggestions:**
```bash
# Add features:
1. Battery detection (more robust)
2. Thermal throttling detection
3. Performance vs efficiency profiles
4. Dynamic adjustment based on load
5. GPU power state management
6. Memory pressure detection
7. CPU frequency monitoring
```

---

## Part 5: Script Execution Order & Dependencies

### Deployment Flow (Tonight)

```
1. Proxmox Installation (manual)
   ↓
2. Copy FortMox project to /opt/
   ↓
3. Edit config/system.yaml
   ↓
4. RUN: sudo ./scripts/deploy.sh
   ├─ Calls: config-manager.sh
   ├─ Calls: power-management.sh (during deploy)
   ├─ Calls: verify.sh (at end)
   └─ Deploys system configuration
   ↓
5. RUN: sudo ./scripts/deploy-vms.sh
   ├─ Creates OPNsense first (priority 1)
   ├─ Then creates other VMs
   └─ All VMs ready
   ↓
6. RUN: sudo ./scripts/verify.sh
   ├─ Validates all configurations
   ├─ Reports issues
   └─ Shows status
   ↓
7. OPTIONAL: sudo ./scripts/gpu-setup.sh
   ├─ Configure GPU passthrough
   ├─ Alternative: use Python tool (Phase 1)
   └─ Requires reboot
   ↓
8. OPTIONAL: sudo ./scripts/detect-hardware.sh
   └─ Shows hardware info & recommendations
```

### Recommended Run Order

```
Setup Phase:
  1. deploy.sh         (required, comprehensive)
  2. deploy-vms.sh     (required, creates VMs)

Verification Phase:
  3. verify.sh         (recommended, validates)

Optional:
  4. detect-hardware.sh (informational)
  5. gpu-setup.sh      (optional, for GPU passthrough)
```

---

## Part 6: Suggested New Scripts Summary

### Quick Priority Matrix

```
High Value + Low Effort:
├─ monitor-vms.sh          (300 lines, great ROI)
├─ network-check.sh        (200 lines, very useful)
└─ backup-manager.sh       (250 lines, critical)

Medium Value + Medium Effort:
├─ vm-clone.sh             (200 lines, useful)
├─ cert-manager.sh         (250 lines, important)
└─ update-manager.sh       (300 lines, important)

High Value + High Effort:
├─ vm-orchestration.sh     (350 lines, powerful)
└─ backup-sync.sh          (400 lines, critical)
```

### Total New Lines Suggested: 2,200+ lines
### Total Enhancement Lines: 500+ lines

---

## Part 7: Implementation Recommendation

### For This Week

```
Priority 1: Build Python GPU Tool
  └─ Replaces gpu-setup.sh with professional tool
     └─ TUI interface (6-8 hours)
     └─ Multi-GPU support
     └─ Per-VM assignment
     └─ Mode switching

Priority 2: Create Backup System
  └─ New script: backup-manager.sh (2-3 hours)
     └─ VM snapshots
     └─ Retention management
     └─ Restore capability
```

### For Next Week

```
Priority 3: Create Monitoring System
  └─ New script: monitor-vms.sh (3-4 hours)
     └─ Real-time metrics
     └─ Alerting
     └─ Data export

Priority 4: Network Management
  └─ New script: network-check.sh (2-3 hours)
     └─ Connectivity verification
     └─ Isolation validation
     └─ Firewall checks
```

### For Following Weeks

```
Priority 5-8: Additional features
  ├─ VM cloning
  ├─ Certificate management
  ├─ Automated updates
  ├─ Cloud backup sync
  └─ VM orchestration
```

---

## Summary: Scripts Status

### Current (7 Scripts - 2,228 lines)
```
✅ Implemented (not yet tested on hardware):
   ├─ deploy.sh (389)
   ├─ config-manager.sh (309)
   ├─ deploy-vms.sh (327)
   ├─ verify.sh (407)
   ├─ detect-hardware.sh (272)
   ├─ power-management.sh (288)
   └─ gpu-setup.sh (236) - limited, will be replaced

Status: READY FOR PRODUCTION DEPLOYMENT TONIGHT
```

### Recommended Next (Phase 1 - 4 weeks)
```
🔜 New Scripts to Add: 3-4 critical scripts
   ├─ backup-manager.sh (backup system)
   ├─ monitor-vms.sh (monitoring)
   ├─ network-check.sh (network validation)
   └─ Python GPU tool (replaces gpu-setup.sh)

Total New Lines: 1,500+ lines
Timeline: Spread across 4 weeks
Priority: Backup first (critical), then monitoring
```

---

## Your Next Steps

### Tonight
1. Use current scripts as-is (implemented, untested on hardware)
2. Deploy FortMox with deploy.sh + deploy-vms.sh
3. Verify with verify.sh

### This Week
1. Build Python GPU tool (our design ready)
2. (Optional) Start backup-manager.sh

### Next Week
1. Add monitoring (monitor-vms.sh)
2. Add network checks (network-check.sh)
3. Complete GTK GPU GUI

### Ongoing
1. Add remaining features as needed
2. Enhance existing scripts
3. Community feedback integration

---

**Status:** This roadmap is historical; current Python utilities are partial ports and the Bash implementations have been removed.
**Recommendation:** Use the Go CLI and verify its documented implementation limits before testing on a Proxmox host.
**Ready to deploy:** No production-readiness claim is made.

