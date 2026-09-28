> **Planning document.** Describes designs and roadmaps, not necessarily what is implemented. See [docs/README.md](../README.md) for current docs.

# FortMox GPU Enhancement Plan
## Features, Implementation, Architecture & Roadmap

---

## Part 1: Your GPU Switching Requirements (Recap)

### Features You Asked For

#### 1. **Interactive GPU Method Selection**
**What you need:**
- Choose between GPU passthrough modes:
  - Full GPU Passthrough (100% performance, single VM)
  - SR-IOV GPU (90% performance, shared across multiple VMs)
  - virtIO GPU (50% performance, unlimited VMs)
  - vGPU (NVIDIA only, 95% performance, multiple VMs)

**User experience:**
- Simple menu or GUI dialog
- Show available options based on hardware
- Explain trade-offs for each method

#### 2. **Multi-GPU Setup with Mixed Modes**
**What you need:**
- If you have 2+ GPUs, assign different modes to each
- Example scenario:
  ```
  GPU 1 (RTX 4090)      → Full Passthrough    → Gaming VM
  GPU 2 (RX 7900)       → SR-IOV (4 VFs)      → Clean VM + Research VM
  GPU 3 (Intel Arc)     → virtIO              → Additional VMs
  ```
- Each GPU independent configuration

**User experience:**
- List all detected GPUs
- Configure each GPU separately
- Preview configuration before applying
- Save/load GPU profiles

#### 3. **Single GPU with Mode Switching**
**What you need:**
- If you only have 1 GPU, switch between modes
- Scenario: Sometimes want full passthrough (gaming), sometimes want shared (work)
- No permanent setting lock

**User experience:**
- Switch modes with minimal downtime
- Possibly without full reboot (where technically possible)
- Save multiple profiles per GPU
- Switch between saved profiles

#### 4. **Per-VM GPU Assignment**
**What you need:**
- Specify which GPU goes to which VM
- Not automatic, but user-controlled
- Example:
  ```
  VM 100 (Clean)    ← GPU 2 (SR-IOV VF0)
  VM 101 (Gaming)   ← GPU 1 (Full Passthrough)
  VM 102 (Research) ← GPU 2 (SR-IOV VF1)
  ```

**User experience:**
- Configuration file or GUI form
- Dropdown for each VM showing available GPU options
- Validation (prevent assigning same GPU twice if full passthrough)

#### 5. **Simple Intuitive Interface**
**What you need:**
- NOT command-line only
- Easy to use for non-expert users
- Clear visual feedback
- Guided setup process

**User experience:**
- Menu-driven or GUI-based
- Step-by-step wizards
- Color-coded status indicators
- Help text for each option
- Real-time validation and error messages

---

## Part 2: Implementation Strategy

### Architecture Overview

```
FortMox GPU Management System
├─ Legacy Bash Scripts (Current - use tonight)
│  └─ scripts/gpu-setup.sh (existing)
│
├─ NEW: Python-Based System (Future - gradual migration)
│  ├─ fortmox_gpu/
│  │  ├─ __init__.py
│  │  ├─ core/
│  │  │  ├─ gpu_detector.py      (GPU detection)
│  │  │  ├─ gpu_manager.py       (GPU configuration)
│  │  │  ├─ vm_manager.py        (VM GPU assignment)
│  │  │  └─ config_parser.py     (YAML config handling)
│  │  ├─ ui/
│  │  │  ├─ cli_menu.py          (Command-line menu)
│  │  │  ├─ gui_gtk.py           (GTK GUI - Linux)
│  │  │  ├─ gui_qt.py            (Qt GUI - Cross-platform)
│  │  │  └─ web_ui.py            (Web dashboard)
│  │  ├─ profiles/
│  │  │  ├─ profile_manager.py   (Save/load profiles)
│  │  │  └─ presets.yaml         (Built-in presets)
│  │  └─ utils/
│  │     ├─ logger.py
│  │     ├─ validators.py
│  │     └─ helpers.py
│  │
│  ├─ CLI Entry Points:
│  │  ├─ fortmox-gpu-setup      (Interactive menu)
│  │  ├─ fortmox-gpu-gui        (GUI launcher)
│  │  ├─ fortmox-gpu-config     (Configuration tool)
│  │  ├─ fortmox-gpu-switch     (Runtime switching)
│  │  └─ fortmox-gpu-profile    (Profile management)
│  │
│  └─ Configuration Files:
│     ├─ /etc/fortmox/gpu-config.yaml
│     ├─ /etc/fortmox/gpu-profiles/
│     └─ ~/.fortmox/local-gpu-config.yaml
│
└─ Backward Compatibility Wrapper
   ├─ scripts/gpu-setup.sh → calls Python version
   ├─ Detects mode preference (bash or python)
   └─ Fallback to bash if python not available
```

### Implementation Phases

#### Phase 0: Backward Compatibility (Tonight)
- Keep current `gpu-setup.sh` working as-is
- Add compatibility wrapper/shim
- Document both paths

#### Phase 1: Python Core (This Week - 6-8 hours)
- Build `fortmox_gpu` Python package
- Implement GPU detector and manager
- Command-line menu interface (TUI)
- Basic profile support

#### Phase 2: GUI (Next Week - 8-10 hours)
- GTK/Qt GUI implementation
- Interactive wizard for GPU setup
- Real-time status monitoring
- Visual GPU assignment

#### Phase 3: Advanced Features (Later - ongoing)
- Web dashboard
- Runtime GPU switching
- Automated profiles
- Integration with Proxmox API

---

## Part 3: Detailed Feature Implementation

### Feature 1: Interactive GPU Method Selection

#### How It Works

**Bash Version (Tonight):**
```bash
# Simple menu
echo "Select GPU passthrough method:"
echo "1) Full Passthrough (100% perf, 1 VM)"
echo "2) SR-IOV (90% perf, 4+ VMs)"
echo "3) virtIO (50% perf, unlimited VMs)"
read -p "Choice: " method
```

**Python Version (This Week):**
```python
class GPUMethodSelector:
    """Interactive GPU method selection with validation"""

    METHODS = {
        'full': {
            'name': 'Full GPU Passthrough',
            'performance': 100,
            'max_vms': 1,
            'description': 'Direct GPU access, best for gaming',
            'requirements': ['IOMMU', 'VFIO support']
        },
        'sr-iov': {
            'name': 'SR-IOV Virtual Functions',
            'performance': 90,
            'max_vms': 4,  # Can be higher
            'description': 'Share GPU across multiple VMs',
            'requirements': ['IOMMU', 'SR-IOV capable GPU']
        },
        'vgpu': {
            'name': 'NVIDIA vGPU',
            'performance': 95,
            'max_vms': 8,
            'description': 'NVIDIA-only, professional use',
            'requirements': ['NVIDIA GRID GPU', 'License']
        },
        'virtio': {
            'name': 'Virtual GPU',
            'performance': 50,
            'max_vms': -1,  # Unlimited
            'description': 'Universal, works everywhere',
            'requirements': ['QEMU/KVM']
        }
    }

    def display_options(self, gpu_device):
        """Show user-friendly GPU method options"""
        # Check which methods supported
        # Display comparison table
        # Get user choice with validation

    def validate_choice(self, gpu, method):
        """Verify method is available for this GPU"""
        # Check hardware capabilities
        # Verify kernel modules loaded
        # Confirm BIOS settings
```

**GUI Version:**
```
┌─────────────────────────────────────────────┐
│  GPU Method Selection                       │
│                                             │
│  GPU: NVIDIA RTX 4090                       │
│                                             │
│  ┌─────────────────────────────────────┐  │
│  │ ✓ Full Passthrough      100% | ○●  │  │
│  │ ✓ SR-IOV                 90% | ●◐  │  │
│  │ ✗ vGPU              (not avail)    │  │
│  │ ✓ virtIO                 50% | ●  │  │
│  └─────────────────────────────────────┘  │
│                                             │
│  Selected: Full Passthrough                 │
│  VMs: 1 max | Performance: Excellent       │
│  Suitable for: Gaming, ML, HPC              │
│                                             │
│  [← Previous] [Next →] [Cancel]            │
└─────────────────────────────────────────────┘
```

#### Implementation Details

**Python GPU Detector:**
```python
class GPUDetector:
    """Detect GPU capabilities and supported methods"""

    def detect_gpus(self):
        """Find all GPUs in system"""
        gpus = []
        for line in self.run_lspci():
            gpu = self.parse_gpu_info(line)
            gpus.append(gpu)
        return gpus

    def get_supported_methods(self, gpu):
        """Determine which passthrough methods available"""
        methods = []

        # All GPUs support virtIO
        methods.append('virtio')

        # Check for SR-IOV
        if self.has_sriov(gpu):
            methods.append('sr-iov')
            gpu['sriov_vfs'] = self.get_sriov_count(gpu)

        # Check for vGPU (NVIDIA only)
        if gpu['vendor'] == 'NVIDIA' and self.has_vgpu_license():
            methods.append('vgpu')

        # All can do full passthrough
        methods.append('full')

        return methods

    def validate_iommu(self):
        """Check IOMMU is enabled"""
        # Check kernel cmdline
        # Verify driver loaded
        # Test IOMMU groups

    def get_iommu_groups(self, gpu):
        """Get IOMMU group info for GPU"""
        # Read /sys/kernel/iommu_groups
        # Return devices in same group
        # Warn if other devices grouped
```

---

### Feature 2: Multi-GPU Setup with Mixed Modes

#### How It Works

**Configuration Flow:**
```
1. Detect all GPUs
   ↓
2. For each GPU:
   - Show capabilities
   - Let user choose method
   - Validate choice
   ↓
3. Save configuration to YAML
   ↓
4. Assign VMs to GPUs
   ↓
5. Apply kernel/GRUB changes
   ↓
6. Reboot with new config
```

**Example YAML Configuration:**
```yaml
# /etc/fortmox/gpu-config.yaml
version: 2

gpus:
  - index: 0
    pci_address: "0000:01:00.0"
    vendor: nvidia
    model: "RTX 4090"
    method: "full"
    assigned_vm:
      vmid: 101
      vm_name: gaming-vm

  - index: 1
    pci_address: "0000:04:00.0"
    vendor: amd
    model: "RX 7900 XTX"
    method: "sr-iov"
    sriov_vfs: 4
    assigned_vms:
      - vmid: 100
        vm_name: clean-vm
        vf_index: 0
      - vmid: 102
        vm_name: research-vm
        vf_index: 1
      - vmid: 103
        vm_name: tools-vm
        vf_index: 2

  - index: 2
    pci_address: "0000:08:00.0"
    vendor: intel
    model: "Arc A770"
    method: "virtio"
    assigned_vm:
      vmid: 100
      vm_name: clean-vm  # Can share virtIO with other VMs
```

**Python Implementation:**
```python
class MultiGPUConfigurator:
    """Configure multiple GPUs with different methods"""

    def detect_all_gpus(self):
        """Find all GPUs in system"""
        detector = GPUDetector()
        return detector.detect_gpus()

    def interactive_setup(self):
        """Step through each GPU interactively"""
        gpus = self.detect_all_gpus()

        config = {
            'version': 2,
            'gpus': [],
            'vms': {}
        }

        for i, gpu in enumerate(gpus):
            print(f"\n=== GPU {i}: {gpu['name']} ===")
            print(f"PCI: {gpu['pci_address']}")

            # Get available methods
            methods = self.get_supported_methods(gpu)
            selected = self.show_method_menu(methods)

            # Get assignment
            assignment = self.assign_to_vms(gpu, selected)

            gpu['method'] = selected
            gpu['assigned_vms'] = assignment
            config['gpus'].append(gpu)

        return config

    def validate_assignment(self, config):
        """Check for conflicts"""
        # Can't assign same full-passthrough GPU to 2 VMs
        # Can't assign same SR-IOV VF to 2 VMs
        # Can assign same virtIO GPU to multiple VMs
        # Can assign same VM multiple GPUs

        errors = []
        for gpu in config['gpus']:
            if gpu['method'] == 'full' and len(gpu['assigned_vms']) > 1:
                errors.append(f"GPU {gpu['pci']} (full) assigned to {len(gpu['assigned_vms'])} VMs")

        return errors

    def generate_kernel_config(self, config):
        """Create GRUB/kernel parameters from config"""
        params = []

        # IOMMU
        if self.is_intel():
            params.append("intel_iommu=on")
        else:
            params.append("amd_iommu=on")
        params.append("iommu=pt")

        # VFIO for full passthrough GPUs
        for gpu in config['gpus']:
            if gpu['method'] == 'full':
                params.append(f"vfio-pci.ids={gpu['vendor_id']}:{gpu['device_id']}")

        return params
```

#### GUI Mockup for Multi-GPU Setup

```
┌─────────────────────────────────────────────────────────────┐
│  Multi-GPU Configuration Wizard                             │
│                                                              │
│  Step 2 of 4: Configure GPU Methods                         │
│                                                              │
│  GPU List:                                                  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ GPU 1: NVIDIA RTX 4090 [0000:01:00.0]               │  │
│  │   Current method: [Full Passthrough ▼]              │  │
│  │   Status: ✓ Supported                                │  │
│  │   Assign to: VM 101 (gaming-vm) [Change]             │  │
│  │                                                       │  │
│  │ GPU 2: AMD RX 7900 [0000:04:00.0]                   │  │
│  │   Current method: [SR-IOV (4 VFs) ▼]               │  │
│  │   Status: ✓ Supported (4 virtual functions)          │  │
│  │   Assigned VMs:                                       │  │
│  │     □ VM 100 (clean-vm) - VF 0   [Remove]          │  │
│  │     □ VM 102 (research-vm) - VF 1 [Remove]         │  │
│  │     [+ Add VM to this GPU]                           │  │
│  │                                                       │  │
│  │ GPU 3: Intel Arc A770 [0000:08:00.0]                │  │
│  │   Current method: [Virtual GPU (virtIO) ▼]         │  │
│  │   Status: ✓ Supported                                │  │
│  │   Assigned VMs: [+ Add VM]                           │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ⓘ Note: Full Passthrough can only be assigned to 1 VM     │
│          SR-IOV and virtIO can be shared                    │
│                                                              │
│  [◄ Previous] [Next ►] [Cancel] [Save Config]             │
└─────────────────────────────────────────────────────────────┘
```

---

### Feature 3: Single GPU Mode Switching

#### How It Works

**Use Case:**
```
Scenario: You have 1 GPU, want to switch between:
- Gaming mode (full passthrough) during gaming sessions
- Work mode (SR-IOV) for multiple VMs during work

Current problem: Requires manual GRUB edit + reboot
New solution: Click menu, select mode, apply changes
```

**Implementation:**
```python
class GPUModeSwitcher:
    """Switch GPU mode between passthrough methods"""

    def get_current_mode(self, gpu):
        """Detect current GPU passthrough mode"""
        # Check /sys/bus/pci/devices/{gpu}/driver
        # Returns: 'full', 'sr-iov', 'virtio', 'native'

    def switch_mode(self, gpu, new_mode):
        """Change GPU passthrough mode"""
        current = self.get_current_mode(gpu)

        if current == new_mode:
            print(f"Already in {new_mode} mode")
            return

        # Check for running VMs using GPU
        running_vms = self.get_vms_using_gpu(gpu)
        if running_vms:
            print(f"VMs using GPU: {running_vms}")
            print("Stop VMs before switching modes")
            return False

        # Perform switch
        self.unbind_gpu(gpu)
        self.bind_gpu_to_method(gpu, new_mode)

        # Optionally update GRUB for persistence
        if self.should_persist():
            self.update_grub_config(gpu, new_mode)
            print("Reboot recommended for full persistence")

        return True

    def unbind_gpu(self, gpu):
        """Unbind GPU from current driver"""
        pci_addr = gpu['pci_address']
        driver_path = f"/sys/bus/pci/devices/{pci_addr}/driver"
        if os.path.exists(driver_path):
            self.write_sysfs(f"{driver_path}/unbind", pci_addr)

    def bind_gpu_to_method(self, gpu, method):
        """Bind GPU to new driver/method"""
        if method == 'full':
            self.bind_vfio(gpu)
        elif method == 'sr-iov':
            self.enable_sriov(gpu)
        elif method == 'virtio':
            self.use_native_driver(gpu)
```

**CLI Interface:**
```bash
$ sudo fortmox-gpu-switch

Current GPU Status:
  GPU 0: NVIDIA RTX 4090 (0000:01:00.0)
  Current mode: Full Passthrough
  Status: Ready (no VMs attached)

Available modes:
  1. Full Passthrough (current)
  2. SR-IOV (4 VFs)
  3. Virtual GPU (virtIO)

Select new mode [1-3]: 2

Switching to SR-IOV mode...
  ✓ Unbinding GPU from VFIO
  ✓ Enabling SR-IOV (4 VFs)
  ✓ Configuring virtual functions
  ✓ Mode switched successfully

Note: Reboot recommended for permanent change
Apply GRUB changes? [Y/n]: y

After reboot, SR-IOV will be active.
```

---

### Feature 4: Per-VM GPU Assignment

#### Configuration Format

**Option A: YAML Configuration**
```yaml
# Part of system.yaml or separate vm-gpu-assignments.yaml
vm_gpu_assignments:
  vm-100:  # clean-vm
    gpus:
      - device: "0000:04:00.0"
        method: "sr-iov"
        vf_index: 0
        mode: "headless"

  vm-101:  # gaming-vm
    gpus:
      - device: "0000:01:00.0"
        method: "full"
        mode: "3d"
        options:
          x-vga: on
          pcie: on
          rombar: 1

  vm-102:  # research-vm
    gpus:
      - device: "0000:04:00.0"
        method: "sr-iov"
        vf_index: 1
        mode: "headless"
```

**Option B: Interactive Assignment**
```python
class VMGPUAssigner:
    """Assign specific GPUs to VMs"""

    def show_assignment_matrix(self, gpus, vms):
        """Display GPU-to-VM assignment table"""
        # Show in GUI or TUI
        # Allow drag-drop or click assignment

    def assign_gpu_to_vm(self, vm_id, gpu_device, method):
        """Add GPU assignment to VM config"""
        # Read VM config file
        # Add hostpci device or vga parameter
        # Validate configuration
        # Write updated config

    def validate_assignments(self):
        """Check for conflicts"""
        # Full passthrough GPU can't be assigned twice
        # SR-IOV VF can't be assigned twice
        # virtIO GPU can be assigned multiple times
```

#### GUI Assignment Interface

```
┌──────────────────────────────────────────────────────┐
│  GPU Assignment Matrix                               │
│                                                       │
│              GPU 1        GPU 2        GPU 3         │
│              (Full)       (SR-IOV)     (virtIO)      │
│  ┌──────────┬──────────┬──────────┬──────────┐     │
│  │ VM 100   │          │ ✓ VF0    │          │     │
│  │ (clean)  │          │          │          │     │
│  ├──────────┼──────────┼──────────┼──────────┤     │
│  │ VM 101   │ ✓        │          │          │     │
│  │ (gaming) │          │          │          │     │
│  ├──────────┼──────────┼──────────┼──────────┤     │
│  │ VM 102   │          │ ✓ VF1    │          │     │
│  │(research)│          │          │          │     │
│  ├──────────┼──────────┼──────────┼──────────┤     │
│  │ VM 103   │          │          │ ✓ shared │     │
│  │ (tools)  │          │          │          │     │
│  └──────────┴──────────┴──────────┴──────────┘     │
│                                                       │
│  Click on any cell to assign/unassign GPU            │
│  ✓ = assigned  (remove by clicking)                  │
│  empty = not assigned  (click to assign)             │
│                                                       │
│  [← Back] [Next] [Cancel]                           │
└──────────────────────────────────────────────────────┘
```

---

### Feature 5: Simple Intuitive Interface

#### Multi-Interface Strategy

**1. TUI (Terminal User Interface) - For SSH/Console**
```
Simple menu-based navigation:
- Colored text
- ASCII art diagrams
- Step-by-step wizards
- Progress indicators
```

**2. Web GUI - For Local/Remote**
```
Browser-based dashboard:
- Responsive design
- Real-time status
- Drag-drop GPU assignment
- Configuration preview
```

**3. Desktop GUI - For Desktop Users**
```
GTK or Qt application:
- Native look and feel
- System integration
- Keyboard shortcuts
- Accessibility features
```

#### Main Menu Example

```
╔════════════════════════════════════════════════════════╗
║                                                        ║
║     FortMox GPU Management Tool                        ║
║     Version 2.0 (Python)                              ║
║                                                        ║
╠════════════════════════════════════════════════════════╣
║                                                        ║
║  1. GPU Setup Wizard       → Configure GPUs for VMs   ║
║  2. Mode Switcher          → Switch GPU modes         ║
║  3. Assignment Manager     → Manage VM-GPU mapping    ║
║  4. Profile Manager        → Save/load configurations ║
║  5. Status Monitor         → View GPU status          ║
║  6. Settings               → Configure preferences    ║
║  7. Documentation          → Help & guides            ║
║  8. Exit                                              ║
║                                                        ║
║  Selection [1-8]: _                                   ║
║                                                        ║
╚════════════════════════════════════════════════════════╝
```

---

## Part 4: Future Enhancement Roadmap

### What You Want to Add Later

Since you mentioned "other changes I would like to make in the future version", here's the roadmap structure:

#### V2.0 (Current Plan - This Week)
```
✓ GPU method selection
✓ Multi-GPU support
✓ Mode switching
✓ Per-VM assignment
✓ Basic GUI
✓ Profile saving
```

#### V2.1 (Coming Soon - Next Month)
```
□ Web dashboard
□ Proxmox API integration
□ Real-time monitoring
□ Temperature/power monitoring
□ Advanced profile builder
```

#### V3.0 (Future Enhancement)
```
□ Machine learning-based GPU allocation
□ Automatic failover
□ Load balancing across GPUs
□ Performance benchmarking
□ GPU resource limits
```

### Your Future Change Placeholder

**Add your planned features here:**

We'll create a dedicated section in the code for your future enhancements:

```python
# fortmox_gpu/future_features.py
"""
Future enhancements planned:

1. [YOUR FEATURE 1 - describe here]
2. [YOUR FEATURE 2 - describe here]
3. [YOUR FEATURE 3 - describe here]

These will be implemented with:
- Modular design to easily add features
- Plugin system for extensibility
- Clear documentation for integration
"""
```

---

## Part 5: Implementation Plan for Tonight & Beyond

### Timeline

```
TONIGHT (Session 1 - ~1 hour):
├─ Document this plan (✓ DONE)
├─ Keep legacy bash script working
├─ Modify for tonight's installation
└─ Use current gpu-setup.sh as-is

THIS WEEK (Phase 1 - 6-8 hours):
├─ Create fortmox_gpu Python package
├─ Implement GPU detector core
├─ Build TUI menu system
├─ Basic profile support
├─ Backward compatibility wrapper
└─ Full testing

NEXT WEEK (Phase 2 - 8-10 hours):
├─ Build GUI (GTK/Qt)
├─ Multi-GPU setup wizard
├─ Mode switching implementation
├─ Per-VM assignment interface
└─ Integration testing

LATER (Phase 3 - ongoing):
├─ Web dashboard
├─ Proxmox API integration
├─ Advanced monitoring
└─ Your custom features
```

### File Structure to Create

```
FortMoxDesktop Workstation/
├─ scripts/
│  ├─ gpu-setup.sh (existing - keep for tonight)
│  ├─ gpu-setup-wrapper.sh (NEW - compatibility layer)
│  └─ gpu-setup-new (NEW - calls Python)
│
├─ fortmox_gpu/ (NEW - Python package)
│  ├─ __init__.py
│  ├─ cli.py (entry point)
│  ├─ core/
│  │  ├─ __init__.py
│  │  ├─ gpu_detector.py
│  │  ├─ gpu_manager.py
│  │  ├─ vm_manager.py
│  │  └─ config_parser.py
│  ├─ ui/
│  │  ├─ __init__.py
│  │  ├─ tui_menu.py (Terminal UI)
│  │  ├─ gui_main.py (Desktop GUI)
│  │  └─ web_ui.py (Web dashboard)
│  ├─ profiles/
│  │  ├─ __init__.py
│  │  ├─ manager.py
│  │  └─ presets.yaml
│  └─ utils/
│     ├─ __init__.py
│     ├─ logger.py
│     └─ helpers.py
│
├─ docs/
│  ├─ docs/planning/gpu-switching-analysis.md (existing)
│  ├─ docs/planning/gpu-enhancement-plan.md (this file)
│  ├─ GPU-PYTHON-GUIDE.md (NEW - Python version guide)
│  ├─ GPU-LEGACY-GUIDE.md (NEW - Bash version guide)
│  └─ GPU-MIGRATION-GUIDE.md (NEW - How to migrate)
│
├─ tests/
│  ├─ test_gpu_detector.py
│  ├─ test_gpu_manager.py
│  └─ test_ui_menus.py
│
├─ setup.py (NEW - Python package installation)
└─ requirements.txt (NEW - Python dependencies)
```

---

## Part 6: Your Input Needed

Before we start building, please clarify:

### Question 1: GUI Preference
- **Option A:** TUI (Terminal menu) - lightweight, SSH-friendly
- **Option B:** Desktop GUI (GTK/Qt) - native feel
- **Option C:** Web UI (Browser) - remote access
- **Option D:** All three (more work but most flexible)

### Question 2: Timeline
- **Option A:** Build everything now before installation (4-5 hours)
- **Option B:** Use bash tonight, add Python later this week
- **Option C:** Use bash tonight, Python is optional someday

### Question 3: Future Features
What other changes do you want in the new Python version?
- Performance monitoring?
- Automated profile selection?
- Integration with other tools?
- Something else?

### Question 4: Profile Support
Should we support:
- **Pre-built profiles** (Gaming, Workstation, Research)?
- **Custom profiles** (user-defined configurations)?
- **Auto-switching profiles** (based on running VMs)?
- **Scheduled profiles** (time-based switching)?

---

## Conclusion

This plan covers:
✓ Your exact GPU switching requirements
✓ Implementation strategy (both bash and Python)
✓ Architecture design
✓ Code examples for each feature
✓ Timeline and file structure
✓ Roadmap for future enhancements

**Next step:** Answer the 4 questions above, then we'll proceed!

---

**Document Created:** November 17, 2025
**Status:** Ready for Review & Planning
**Next Action:** Get your input on implementation approach

