> **Planning document.** Describes designs and roadmaps, not necessarily what is implemented. See [docs/README.md](../README.md) for current docs.

# FortMox GPU Switching Analysis & Feature Documentation

## Executive Summary

**Current Status:** The `gpu-setup.sh` script provides **limited GPU switching capabilities**.

**Current Capabilities:**
- ✅ GPU detection (NVIDIA, AMD, Intel)
- ✅ IOMMU support verification
- ✅ SR-IOV capability detection
- ✅ One-way configuration to VFIO Full Passthrough
- ❌ Runtime switching between GPU modes
- ❌ Multi-GPU support with mixed modes
- ❌ Interactive menu for GPU method selection
- ❌ VM-level GPU assignment flexibility

**Verdict:** The current implementation is **basic** and **not production-ready** for your use case.

---

## Current GPU Setup Script Analysis

### What It Does (Lines 1-237)

```bash
gpu-setup.sh
├─ Detects GPU (vendor, model, PCI address)
├─ Identifies GPU type (NVIDIA/AMD/Intel)
├─ Checks IOMMU support
├─ Detects SR-IOV capability (if available)
├─ Configures VFIO for full passthrough
├─ Updates GRUB with vfio-pci parameters
├─ Creates modprobe blacklist
├─ Updates kernel initramfs
└─ Provides summary and reboot instructions
```

### Limitations

#### 1. **No Interactive GPU Selection**
```bash
# Current approach (Line 100-114):
log_info "Detecting best GPU passthrough method..."

if [ -f "/sys/bus/pci/devices/0000:${GPU_PCIADDR}/sriov_totalvfs" ]; then
    SRIOV_TOTAL=$(cat "/sys/bus/pci/devices/0000:${GPU_PCIADDR}/sriov_totalvfs" 2>/dev/null || echo "0")
    if [ "$SRIOV_TOTAL" -gt 0 ]; then
        log_success "SR-IOV capable ($SRIOV_TOTAL VFs)"
        METHODS=("SR-IOV" "Full Passthrough" "virtIO")
    else
        log_warn "SR-IOV not available"
        METHODS=("Full Passthrough" "virtIO")
    fi
fi
```

**Problem:** Detects available methods BUT doesn't let user choose! Always configures VFIO full passthrough (Line 126).

#### 2. **Single GPU Only**
```bash
# Line 40: Only detects first GPU
GPU_PCIADDR=$(lspci | grep -iE "vga|3d controller" | head -1 | awk '{print $1}')
                                                              ↑
                                                        Only first match
```

**Problem:** Multi-GPU systems are ignored. No support for GPU 0 → full passthrough, GPU 1 → SR-IOV, etc.

#### 3. **One-Way Configuration**
```bash
# Lines 138-189: Sets up VFIO configuration
# No mechanism to switch back or choose different mode
# Requires manual editing and reboot to change
```

**Problem:** Once deployed, switching modes requires manual GRUB editing, modprobe modification, and reboot.

#### 4. **No Per-VM GPU Assignment**
```bash
# VM templates suggest flexibility (gaming-vm.yaml lines 45-51):
gpu:
    type: passthrough
    device: auto  # <-- Just says "auto"
    vfio_binding: true
```

**Problem:** No way to specify which GPU device, method, or sharing level per VM.

#### 5. **Missing Features from Config**
```yaml
# From system.yaml (lines 25-38), config mentions:
gpu:
    passthrough: true
    method: auto  # Supports: auto, virtio, sr-iov, vgpu, full
    devices: []   # List of GPU UUIDs

# BUT gpu-setup.sh doesn't implement this!
# It just forces 'full' passthrough always.
```

---

## What Should Be Implemented

### Feature 1: Interactive GPU Selection Menu

```bash
# MISSING: Interactive menu like this:

╔════════════════════════════════════════════════════════╗
║  FortMox GPU Configuration Tool                        ║
╚════════════════════════════════════════════════════════╝

Detected GPUs:
  [1] NVIDIA RTX 4090 (0000:01:00.0)
  [2] AMD RX 7900 XTX (0000:04:00.0)

Supported Methods:
  Full Passthrough: 100% performance, 1 VM only
  SR-IOV:          90% performance, 4 VMs shared
  vGPU (NVIDIA):   95% performance, 8 VMs shared
  virtIO:          50% performance, unlimited VMs

Select GPU for config: [1/2/both/cancel]
Select method for GPU [1]: [1=Full/2=SR-IOV/3=virtIO/4=vGPU]
Select method for GPU [2]: [1=Full/2=SR-IOV/3=virtIO/4=vGPU]

Configuration Summary:
  GPU 1 → Full Passthrough (Gaming VM)
  GPU 2 → SR-IOV 4vfs (Clean + Research VMs)

Proceed? [Y/n]
```

### Feature 2: Multi-GPU Support

```bash
# MISSING: Multi-GPU detection and configuration

Detect multiple GPUs:
  ✗ Current: only detects first GPU
  ✓ Needed: detect all GPUs with lspci loop

assign_gpu_per_vm() {
    local vm_id=$1
    local gpu_device=$2
    local method=$3

    # Add to VM config:
    # hostpci0: GPU_PCI,pcie=on,x-vga=on (full passthrough)
    # vga: qxl (virtIO)
    # sr-iov-vf: 4 (SR-IOV)
}
```

### Feature 3: Runtime Mode Switching

```bash
# MISSING: Interactive switching without full reboot

switch_gpu_mode() {
    local gpu_device=$1
    local new_method=$2

    case $new_method in
        full)
            # Bind to VFIO
            echo "$gpu_device" > /sys/bus/pci/devices/$gpu_device/driver/unbind
            echo "$vendor:$device" > /sys/bus/pci/drivers/vfio-pci/new_id
            ;;
        sr-iov)
            # Enable SR-IOV
            echo 4 > /sys/bus/pci/devices/$gpu_device/sriov_numvfs
            ;;
        virtio)
            # Use virtual GPU
            # Update VM config: vga: virtio
            ;;
    esac
}
```

### Feature 4: Configuration Persistence

```bash
# MISSING: Store GPU configuration for reproducibility

create_gpu_config_file() {
    cat > /etc/fortmox/gpu-config.yaml << EOF
gpus:
  - device: "0000:01:00.0"
    name: "NVIDIA RTX 4090"
    method: "full"
    assigned_vm: 101  # Gaming VM

  - device: "0000:04:00.0"
    name: "AMD RX 7900 XTX"
    method: "sr-iov"
    vfs: 4
    assigned_vms:
      - vmid: 100
        vf_index: 0
      - vmid: 102
        vf_index: 1
EOF
}
```

---

## Current Implementation vs Your Requirements

### Your Requirements

| Feature | Your Need | Current Status | Gap |
|---------|-----------|---|---|
| **SR-IOV GPU** | Supported | Detected only | ❌ No user control |
| **virtIO GPU** | Supported | Available in code | ❌ Not selectable |
| **Full GPU Passthrough** | Supported | ✅ Implemented | ✅ Works |
| **Mixed Multi-GPU Setup** | Critical | ❌ Not supported | ❌ Major gap |
| **Gaming VM full passthrough** | Critical | ✅ Can do | ✅ Works |
| **Regular VMs + one GPU** | Critical | ❌ Not supported | ❌ Major gap |
| **Single GPU mode switching** | Important | ❌ Requires manual edit | ❌ Not automated |
| **Simple Intuitive Interface** | Critical | ❌ Script-only | ❌ No menu system |
| **Per-VM GPU assignment** | Critical | ❌ Not possible | ❌ Not implemented |

---

## What Needs to Be Built

### Priority 1: Interactive GPU Configuration Menu

**File to create:** `scripts/gpu-interactive-setup.sh`

**Features:**
- Detect all GPUs (not just first)
- Show capabilities per GPU
- Interactive menu for method selection
- Preview configuration before applying
- Save configuration to file
- Auto-generate VM configs

**Expected lines of code:** 200-300 lines

**Complexity:** Medium

**Time estimate:** 2-3 hours

### Priority 2: Multi-GPU Management

**File to create:** `scripts/gpu-multi-gpu-manager.sh`

**Features:**
- List all GPUs with details
- Bind/unbind GPUs dynamically
- Assign GPU to VMs
- Support mixed modes per GPU
- Configuration validation

**Expected lines of code:** 250-350 lines

**Complexity:** High

**Time estimate:** 3-4 hours

### Priority 3: Runtime GPU Switching

**File to create:** `scripts/gpu-mode-switcher.sh`

**Features:**
- Switch GPU mode without reboot (where possible)
- VFIO binding/unbinding
- SR-IOV enable/disable
- Migrate VMs when switching modes
- Safety checks before switching

**Expected lines of code:** 200-250 lines

**Complexity:** High

**Time estimate:** 3-5 hours

### Priority 4: GUI Dashboard (Optional but recommended)

**Tool:** Web-based dashboard or TUI (Text UI)

**Features:**
- Visual GPU status display
- Click-to-configure interface
- Real-time VM GPU assignment
- Performance monitoring

**Expected lines of code:** 400-500 lines

**Complexity:** Very High

**Time estimate:** 5-8 hours

---

## Current Config File Analysis

### System Config GPU Section (system.yaml)

```yaml
gpu:
    passthrough: true
    # Supported: auto, virtio, sr-iov, vgpu, full
    method: auto              # ← Script ignores this!
    devices: []               # ← Script ignores this!
```

**Status:** Configuration file is ready but ignored by gpu-setup.sh

### Gaming VM Config (gaming-vm.yaml)

```yaml
gpu:
    type: passthrough
    device: auto              # ← Not implemented
    vfio_binding: true

gpu_config:
    method: auto              # ← Not used
    sr_iov:
        enabled: false        # ← Not configurable
        vfs: 4
    vgpu:
        enabled: false        # ← Not supported
        profile: auto
```

**Status:** Comprehensive config exists but not utilized

### Clean VM Config (clean-vm.yaml)

```yaml
gpu:
    type: virtio              # ← Correctly configured
```

**Status:** Simple virtIO works fine

---

## Recommendation: Build These Features (High Priority)

I recommend implementing in this order:

### Phase 1 (Tonight): Quick Win
1. **Enhance gpu-setup.sh** to:
   - Detect multiple GPUs
   - Show methods in menu
   - Let user choose

### Phase 2 (This Week): Core Functionality
2. **Create gpu-interactive-setup.sh** with menu system
3. **Update deploy-vms.sh** to read GPU assignments from config

### Phase 3 (Next Week): Advanced Features
4. **Create gpu-mode-switcher.sh** for runtime changes
5. **Create gpu-multi-gpu-manager.sh** for management

### Phase 4 (Later): Nice to Have
6. **Create web dashboard** for visual management

---

## Implementation Timeline

```
Tonight:
  ├─ Enhance existing gpu-setup.sh (1-2 hours)
  └─ Test with your systems (1 hour)

This Week:
  ├─ Build interactive menu system (2 hours)
  ├─ Multi-GPU detection (1.5 hours)
  └─ Documentation (1 hour)

This Month:
  ├─ Runtime switching (3 hours)
  ├─ Web dashboard (5 hours)
  └─ Full testing (2 hours)

Total: 16-18 hours of development
```

---

## Code Examples

### Example 1: Simple GPU Detection Enhancement

```bash
#!/bin/bash
# Detect all GPUs (not just first)

declare -a GPU_DEVICES
declare -a GPU_METHODS

# Find all GPUs
mapfile -t GPU_DEVICES < <(lspci | grep -iE "vga|3d controller" | awk '{print $1}')

echo "Found ${#GPU_DEVICES[@]} GPU(s)"

for i in "${!GPU_DEVICES[@]}"; do
    GPU="${GPU_DEVICES[$i]}"
    GPU_NAME=$(lspci -s "$GPU" | cut -d: -f3-)

    # Check if SR-IOV capable
    if [ -f "/sys/bus/pci/devices/0000:${GPU}/sriov_totalvfs" ]; then
        METHODS=("Full" "SR-IOV" "virtIO")
    else
        METHODS=("Full" "virtIO")
    fi

    echo "[$i] $GPU_NAME ($GPU) - Methods: ${METHODS[*]}"
done
```

### Example 2: Interactive Menu

```bash
#!/bin/bash
# Simple interactive GPU selection

select_gpu_method() {
    local gpu=$1

    echo ""
    echo "Select GPU method for $gpu:"
    echo "  1) Full Passthrough   (100% perf, 1 VM)"
    echo "  2) SR-IOV            (90% perf, 4 VMs)"
    echo "  3) virtIO            (50% perf, unlimited)"
    echo "  4) Skip this GPU"
    echo ""
    read -p "Choice [1-4]: " choice

    case $choice in
        1) echo "full" ;;
        2) echo "sr-iov" ;;
        3) echo "virtio" ;;
        4) echo "skip" ;;
        *) echo "full" ;;
    esac
}
```

---

## Conclusion

**Current State:** Basic but functional for single GPU setups

**Your Use Case:** Requires significant enhancements for multi-GPU mixed mode support

**Recommendation:** I can build this enhancement package for you during Phase 1 (tonight) and Phase 2 (this week)

**Next Steps:**
1. Do you want me to build the enhanced GPU switching system?
2. What's your priority: gaming performance or multi-VM flexibility?
3. Would you prefer command-line menu or web dashboard?

---

## Files to Document

```
Enhanced/New Scripts Needed:
├─ scripts/gpu-setup.sh (enhance existing)
├─ scripts/gpu-interactive-setup.sh (NEW)
├─ scripts/gpu-multi-gpu-manager.sh (NEW)
├─ scripts/gpu-mode-switcher.sh (NEW)
├─ docs/gpu-switching-guide.md (NEW)
└─ config/gpu-profiles.yaml (NEW)

Enhanced Config Files:
├─ config/system.yaml (respect method and devices)
├─ config/vm-templates/clean-vm.yaml (per-VM GPU assignment)
├─ config/vm-templates/gaming-vm.yaml (per-VM GPU assignment)
└─ config/vm-templates/research-vm.yaml (per-VM GPU assignment)
```

---

**Analysis Date:** November 17, 2025
**Analysis Scope:** gpu-setup.sh, system.yaml, VM templates
**Status:** Ready for enhancement

