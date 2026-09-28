> **Planning document.** Describes designs and roadmaps, not necessarily what is implemented. See [docs/README.md](../README.md) for current docs.

# FortMox GPU Tool - Quick Start Summary
## What You Asked For, What We're Building, When

---

## TL;DR - Your Exact Requirements

### Your 5 GPU Switching Features

```
┌─────────────────────────────────────────────────────────────┐
│ FEATURE 1: Interactive GPU Method Selection                │
│ ┌───────────────────────────────────────────────────────┐  │
│ │ User chooses:                                         │  │
│ │ • Full Passthrough (100% perf, 1 VM)                 │  │
│ │ • SR-IOV (90% perf, multiple VMs)                    │  │
│ │ • vGPU (95% perf, NVIDIA only)                       │  │
│ │ • virtIO (50% perf, unlimited VMs)                   │  │
│ │                                                       │  │
│ │ With menu showing what's available on their hardware │  │
│ └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ FEATURE 2: Multi-GPU with Mixed Modes                      │
│ ┌───────────────────────────────────────────────────────┐  │
│ │ User's scenario:                                      │  │
│ │ GPU 1 (RTX 4090)  → Full Passthrough  → Gaming       │  │
│ │ GPU 2 (RX 7900)   → SR-IOV (4 VFs)    → Work VMs     │  │
│ │ GPU 3 (Arc)       → virtIO            → Research VM  │  │
│ │                                                       │  │
│ │ Each GPU independent, different methods              │  │
│ └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ FEATURE 3: Single GPU Mode Switching                       │
│ ┌───────────────────────────────────────────────────────┐  │
│ │ One GPU, multiple modes:                              │  │
│ │                                                       │  │
│ │ Monday-Friday:  SR-IOV mode   (work with 4 VMs)     │  │
│ │ Weekend:        Full mode     (gaming)               │  │
│ │                                                       │  │
│ │ Switch with menu, optional persistence               │  │
│ └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ FEATURE 4: Per-VM GPU Assignment                           │
│ ┌───────────────────────────────────────────────────────┐  │
│ │ User assigns:                                         │  │
│ │ VM 100 (Clean)    ← GPU 2 SR-IOV VF0                │  │
│ │ VM 101 (Gaming)   ← GPU 1 Full Passthrough          │  │
│ │ VM 102 (Research) ← GPU 2 SR-IOV VF1                │  │
│ │                                                       │  │
│ │ Via config file or interactive GUI                   │  │
│ └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ FEATURE 5: Simple Intuitive Interface                      │
│ ┌───────────────────────────────────────────────────────┐  │
│ │ NOT: confusing command-line                          │  │
│ │ BUT: Menu system or GUI                              │  │
│ │                                                       │  │
│ │ Options:                                              │  │
│ │ • Terminal menu (TUI) - lightweight                  │  │
│ │ • Desktop GUI (GTK) - full-featured                  │  │
│ │ • Web dashboard (optional) - remote access           │  │
│ └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## Your Decisions

### Decision 1: What Interface?
✅ **Your Choice:** All three (TUI + GTK Desktop + Web)
   - **Priority:** TUI and GTK most important
   - **Web:** Optional nice-to-have

### Decision 2: When?
✅ **Your Choice:** Use bash tonight, Python this week
   - **Tonight:** Current gpu-setup.sh for FortMox installation
   - **This Week:** Build Python version in parallel
   - **Next Week:** Deploy Python version when ready

### Decision 3: Keep the bash script?
✅ **Your Choice:** Yes, keep both versions
   - Both coexist
   - Users choose which to use
   - Easy comparison and fallback

### Decision 4: Future features?
✅ **Your Choice:** All of them!
   - Performance monitoring
   - Automated profile selection
   - Grafana/Prometheus integration

---

## Timeline

```
TONIGHT (Session 1)
├─ 20 min: Review this summary
├─ 30 min: Install Proxmox Desktop
├─ 30 min: Install Proxmox Laptop
├─ 45 min: Deploy FortMox (both systems)
├─ 30 min: Create VMs
└─ DONE! Everything working by end of night

THIS WEEK (Phase 1 - 6-8 hours)
├─ Tuesday: Build Python GPU package core
├─ Tuesday: Implement GPU detector
├─ Wednesday: Build TUI menu system
├─ Wednesday: Test multi-GPU setup
└─ Thursday: Documentation complete

NEXT WEEK (Phase 2 - 8-10 hours)
├─ Thursday: Build GTK Desktop GUI
├─ Friday: Create setup wizard
├─ Friday: Integration testing
└─ Weekend: Polish and optimization

LATER (Phase 3 - ongoing)
├─ Web UI (optional)
├─ Performance monitoring
├─ Automated profiles
└─ Grafana integration
```

---

## Implementation Phases

### Phase 0: RIGHT NOW (Tonight)
**What:** Keep everything as-is, use current scripts

**Files used:**
- `scripts/gpu-setup.sh` (existing bash)
- `config/system.yaml`
- VM templates

**Result:** FortMox running with current GPU setup

---

### Phase 1: THIS WEEK
**What:** Build Python GPU tool with TUI

**Files to create:**
```
fortmox_gpu/
├─ __init__.py                   ← Package definition
├─ cli.py                        ← Command entry point
├─ core/
│  ├─ gpu_detector.py           ← Find GPUs
│  ├─ gpu_manager.py            ← Configure GPUs
│  ├─ vm_manager.py             ← Assign to VMs
│  └─ config_parser.py          ← Read/write configs
├─ ui/
│  └─ tui_menu.py               ← Terminal menus
├─ profiles/
│  ├─ manager.py                ← Save/load configs
│  └─ presets.yaml              ← Built-in templates
└─ utils/
   ├─ logger.py
   └─ helpers.py

setup.py                          ← Python installer
requirements.txt                  ← Dependencies
```

**Features:**
- ✓ Detect all GPUs
- ✓ Show available methods
- ✓ Interactive selection menu
- ✓ Multi-GPU configuration
- ✓ VM assignment
- ✓ Profile management
- ✓ Mode switching
- ✓ Validation

**Command:** `sudo fortmox-gpu-setup`

---

### Phase 2: NEXT WEEK
**What:** Add GTK Desktop GUI

**Additional files:**
```
fortmox_gpu/ui/
├─ gui_gtk.py                  ← GTK application
├─ widgets/
│  ├─ gpu_list.py             ← GPU list display
│  ├─ method_selector.py       ← Method selection
│  ├─ vm_assignment.py         ← GPU-VM matrix
│  └─ status_monitor.py        ← Real-time status
├─ dialogs/
│  ├─ setup_wizard.py          ← Step-by-step wizard
│  ├─ gpu_config_dialog.py     ← GPU configuration
│  └─ vm_assignment_dialog.py  ← VM assignment
└─ resources/
   ├─ icons/
   └─ themes/
```

**Features:**
- ✓ Visual GPU list
- ✓ Drag-and-drop assignment
- ✓ Setup wizard (multi-step)
- ✓ Configuration preview
- ✓ Real-time status
- ✓ System tray integration

**Command:** `sudo fortmox-gpu-gui`

---

### Phase 3: LATER (Optional)
**What:** Add Web UI and future features

**Additional files:**
```
fortmox_gpu/
├─ ui/
│  └─ web_ui.py               ← Flask web app
├─ api/
│  ├─ gpu_api.py              ← REST API
│  └─ auth.py                 ← Authentication
├─ monitoring/                 ← FUTURE
│  ├─ gpu_monitor.py          ← Temperature, power
│  ├─ thermal_monitor.py      ← Heat monitoring
│  └─ metrics_exporter.py     ← Prometheus export
├─ automation/                 ← FUTURE
│  ├─ profile_auto_selector.py ← Auto GPU allocation
│  └─ smart_allocator.py      ← Intelligent assignment
└─ web/                        ← Web assets
   ├─ static/
   └─ templates/
```

**Features:**
- ✓ Browser dashboard
- ✓ Remote access
- ✓ Real-time monitoring
- ✓ Temperature tracking
- ✓ Power monitoring
- ✓ Prometheus integration
- ✓ Auto-profile selection

**Commands:** `fortmox-gpu-web` (browser at http://localhost:5000)

---

## How It Works: Example Walkthrough

### Scenario: You have 2 GPUs and want GPU setup

#### Using TUI (This Week)

```bash
$ sudo fortmox-gpu-setup

╔════════════════════════════════════════════╗
║  FortMox GPU Setup Wizard                 ║
║  Using Terminal Interface (TUI)           ║
╚════════════════════════════════════════════╝

Step 1: Detect GPUs
  ✓ Found 2 GPUs:
    [1] NVIDIA RTX 4090 (0000:01:00.0)
    [2] AMD RX 7900 (0000:04:00.0)

Step 2: Configure GPU 1 (NVIDIA RTX 4090)
  Supported methods:
    1. Full Passthrough (100%)
    2. SR-IOV (not supported for NVIDIA)
    3. virtIO (50%)

  Select [1-3]: 1

  ✓ Full Passthrough selected
  Assign to VM: [101 - gaming-vm] [Change] ✓

Step 3: Configure GPU 2 (AMD RX 7900)
  Supported methods:
    1. Full Passthrough (100%)
    2. SR-IOV (90%, 4 VFs)
    3. virtIO (50%)

  Select [1-3]: 2

  ✓ SR-IOV selected (4 virtual functions)
  Assign VMs:
    VM 100 (clean-vm) → VF0 [✓]
    VM 102 (research-vm) → VF1 [✓]
    [+ Add another VM]

Step 4: Review Configuration
  GPU 1: NVIDIA RTX 4090 → Full → VM 101 (gaming)
  GPU 2: AMD RX 7900 → SR-IOV (4) → VM 100 + VM 102

  Apply changes? [y/n]: y

  ✓ Configuration saved to /etc/fortmox/gpu-config.yaml
  ✓ Kernel parameters updated
  ✓ Reboot recommended for persistence

  Reboot now? [y/n]: y
  System rebooting...
```

#### Using Desktop GUI (Next Week)

```
Window appears with:
- List of GPUs (NVIDIA, AMD visible)
- For each GPU:
  - Radio buttons: Full / SR-IOV / virtIO
  - Selected method shows details
- Grid showing VM assignments
  - Drag-drop GPU to VM row
  - Visual validation
- Preview of resulting config
- Apply button
```

#### Using Web UI (Later)

```
Browser: http://localhost:5000/gpu-setup

Dashboard shows:
- Visual GPU cards
- Method selection dropdowns
- VM assignment table (drag-drop capable)
- Real-time status indicators
- Configuration preview pane
- Save/Apply buttons
```

---

## Commands You'll Use

### Tonight (Bash)
```bash
sudo /opt/FortMoxDesktop\ Workstation/scripts/gpu-setup.sh
```

### This Week (Python TUI)
```bash
sudo fortmox-gpu-setup          # Interactive wizard
sudo fortmox-gpu-config list    # Show current config
sudo fortmox-gpu-profile load   # Load saved profile
```

### Next Week (Python GUI)
```bash
sudo fortmox-gpu-gui            # Launch graphical app
sudo fortmox-gpu-config export  # Export to YAML
```

### Later (Monitoring)
```bash
sudo fortmox-gpu-monitor        # Real-time monitoring
sudo fortmox-gpu-metrics export # Export to Prometheus
```

---

## Documents Created

**Analysis & Planning:**
- `docs/planning/gpu-switching-analysis.md` - Current system analysis
- `docs/planning/gpu-enhancement-plan.md` - Feature details
- `docs/planning/implementation-roadmap.md` - Timeline & tasks
- `docs/planning/gpu-tool-quick-start.md` - This file

**Building:**
- `GPU-PYTHON-GUIDE.md` - Python version docs (NEXT)
- `GPU-LEGACY-GUIDE.md` - Bash version docs (NEXT)
- `GPU-MIGRATION-GUIDE.md` - How to upgrade (NEXT)

---

## What Happens Next

### Immediate (Next 10 minutes)
1. You review this summary
2. Confirm everything is correct
3. Ask any questions

### Tonight
1. Use current gpu-setup.sh as-is
2. Install FortMox on Desktop
3. Install FortMox on Laptop
4. Everything working by end of night

### This Week
1. I build Python GPU tool with TUI
2. Test with your FortMox systems
3. Document Python version
4. You try it and give feedback

### Next Week
1. I add GTK Desktop GUI
2. Full integration testing
3. Complete documentation
4. Ready for production use

### Later
1. Optional web UI
2. Performance monitoring
3. Automated profiles
4. Grafana integration

---

## Confirmation Checklist

Before we proceed, please confirm:

- [ ] You understand the 5 GPU switching features
- [ ] Timeline is good (bash tonight, Python this week)
- [ ] Interface choices are correct (TUI + GTK priority)
- [ ] Backward compatibility plan is acceptable
- [ ] Future features (monitoring, automation) are what you want
- [ ] You're ready to proceed with tonight's installation

**If all checked:** Let's build it! 🚀

---

**Status:** Ready for tonight's installation
**Created:** November 17, 2025
**Next Action:** Your confirmation, then begin!

