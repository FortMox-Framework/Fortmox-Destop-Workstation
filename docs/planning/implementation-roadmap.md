> **Planning document.** Describes designs and roadmaps, not necessarily what is implemented. See [docs/README.md](../README.md) for current docs.

# FortMox GPU Tool - Implementation Roadmap
## Your Requirements & Decisions Summary

**Document Date:** November 17, 2025

---

## Your GPU Switching Requirements (Summary)

You need the GPU management system to support:

### ✅ Feature 1: Interactive GPU Method Selection
```
What it does:
- Choose between Full Passthrough, SR-IOV, vGPU, virtIO
- Show available methods based on hardware
- Explain trade-offs (performance vs. sharing)

Implementation: Menu-based or GUI selection
```

### ✅ Feature 2: Multi-GPU Setup with Mixed Modes
```
What it does:
- Configure multiple GPUs independently
- Example: GPU1 → Full (gaming), GPU2 → SR-IOV (work)
- Each GPU gets its own method and VM assignments

Implementation: Configuration wizard, step through each GPU
```

### ✅ Feature 3: Single GPU Mode Switching
```
What it does:
- Have 1 GPU but want to switch between modes
- Example: Full passthrough for gaming, SR-IOV for work
- Switch without permanent lock-in

Implementation: Runtime mode switcher, optional persistence
```

### ✅ Feature 4: Per-VM GPU Assignment
```
What it does:
- Specify which GPU goes to which VM
- Example: VM 100 → GPU 2 VF0, VM 101 → GPU 1 Full
- User controls mapping, not automatic

Implementation: Configuration file or interactive assignment GUI
```

### ✅ Feature 5: Simple Intuitive Interface
```
What it does:
- NO command-line only, but user-friendly
- Menu-driven or graphical
- Visual feedback and help text
- Step-by-step wizards

Implementation: Multiple UI options (TUI, Desktop, Web)
```

---

## Your Implementation Decisions

### Decision 1: GUI Interface Preference
**Your Choice:** All three (TUI + Desktop GUI + Web), BUT TUI & GTK Desktop are most important

**What this means:**
- **TUI (Terminal UI):** Menu system in terminal → lightweight, SSH-friendly
- **GTK Desktop GUI:** Native Linux GUI application → full-featured
- **Web UI:** Browser-based → optional nice-to-have
- **Priority:** GTK Desktop and TUI are core, Web is bonus

**Implementation:**
```
Phase 1 (This Week):
  ✓ TUI menu system (6-8 hours)
  - Command-line menus with color, ASCII art
  - Works over SSH
  - All features available

Phase 2 (Next Week):
  ✓ GTK Desktop GUI (8-10 hours)
  - Native Linux application
  - Drag-and-drop GPU assignment
  - Real-time status
  - File manager integration

Phase 3 (Later):
  □ Web UI (optional - 6-8 hours)
  - Browser dashboard
  - Remote access
  - Mobile-friendly
```

### Decision 2: Deployment Timeline
**Your Choice:** Use bash tonight for FortMox installation, add Python GPU tool this week

**What this means:**
- **Tonight:** Use existing `gpu-setup.sh` as-is
- **This Week:** Build Python version in parallel
- **Next Week:** Deploy Python GUI version when ready
- **Benefit:** No delays to FortMox installation, professional GPU tool by end of week

**Schedule:**
```
TONIGHT (Session 1 - 2 hours):
  - Review FortMox deployment plan
  - Install Proxmox on Desktop
  - Install Proxmox on Laptop
  - Deploy FortMox with current gpu-setup.sh
  - Everything working by night's end

TUESDAY-WEDNESDAY (Phase 1 - 6-8 hours):
  - Build Python GPU package
  - Implement TUI menu system
  - GPU detector core
  - Basic profiles

THURSDAY-FRIDAY (Phase 2 - 8-10 hours):
  - Build GTK Desktop GUI
  - Multi-GPU setup wizard
  - Mode switching
  - Integration testing

WEEKEND & BEYOND (Phase 3 - ongoing):
  - Web UI (optional)
  - Advanced monitoring
  - Automation features
  - Polish and optimization
```

### Decision 3: Backwards Compatibility
**Your Choice:** Keep bash script unchanged, add Python as alternative

**What this means:**
- Both `gpu-setup.sh` (bash) and Python version exist
- Users can choose which to use
- No forced migration
- Easy to compare implementations

**File structure:**
```
scripts/
├─ gpu-setup.sh                  (existing bash - unchanged)
├─ gpu-setup.sh.bak             (backup of original)
└─ gpu-setup                     (new Python entry point)

fortmox_gpu/                      (new Python package)
├─ __init__.py
├─ cli.py
├─ core/
├─ ui/
└─ profiles/

To use Python version: sudo fortmox-gpu-setup
To use bash version:   sudo ./scripts/gpu-setup.sh
```

---

## Your Future Enhancement Features

**Your Choice:** All of the following (for future versions)

### Future Feature 1: Performance Monitoring
```
Monitor GPU health and usage:
- Temperature monitoring
- Power consumption tracking
- GPU utilization percentage
- Memory usage per VM
- Real-time dashboards
- Historical graphs
```

### Future Feature 2: Automated Profile Selection
```
Intelligent GPU allocation:
- Auto-detect running VMs
- Automatically choose best GPU mode
- Switch profiles based on VM activity
- Smart resource allocation
- Power optimization
```

### Future Feature 3: Monitoring Integration
```
Export to external dashboards:
- Prometheus metrics export
- Grafana dashboard integration
- InfluxDB support
- Custom alert thresholds
- Historical data storage
```

**How we'll build this:**
```
fortmox_gpu/
├─ core/
│  └─ [existing GPU management]
│
├─ monitoring/  (NEW - Future)
│  ├─ gpu_monitor.py           (GPU metrics collection)
│  ├─ thermal_monitor.py       (Temperature tracking)
│  ├─ power_monitor.py         (Power usage)
│  └─ metrics_exporter.py      (Prometheus/Grafana export)
│
├─ automation/  (NEW - Future)
│  ├─ profile_auto_selector.py (Auto-detect running VMs)
│  ├─ smart_allocator.py       (Intelligent assignment)
│  └─ scheduler.py             (Time-based switching)
│
└─ integrations/  (NEW - Future)
   ├─ prometheus.py            (Prometheus exporter)
   ├─ grafana.py               (Grafana dashboard)
   ├─ influxdb.py              (InfluxDB connector)
   └─ custom_hooks.py          (Custom integrations)
```

---

## Building Strategy

### Why We Build This Way

**1. Keep Bash Script Tonight**
- No risk to installation
- Known, tested solution
- Can compare implementations later
- Users can choose preferred version

**2. Build Python in Parallel**
- Professional, feature-rich tool
- Extensible for future features
- Better code organization
- Easier maintenance and updates

**3. Phase Implementation**
- TUI first (lowest dependency, most useful)
- Desktop GUI next (higher quality)
- Web UI later (optional enhancement)
- Future features integrated naturally

**4. Modular Architecture**
- GPU detection module reusable
- UI modules can be swapped
- Easy to add monitoring later
- Plugin system for extensibility

---

## Implementation Details

### Phase 1: TUI + Core (This Week)

```
Build timeline: 6-8 hours

Files to create:
├─ fortmox_gpu/__init__.py         (package init)
├─ fortmox_gpu/cli.py              (entry point)
├─ fortmox_gpu/core/
│  ├─ gpu_detector.py              (detect GPUs & capabilities)
│  ├─ gpu_manager.py               (configure GPUs)
│  ├─ vm_manager.py                (manage VM assignments)
│  └─ config_parser.py             (YAML config handling)
├─ fortmox_gpu/ui/
│  └─ tui_menu.py                  (terminal menu system)
├─ fortmox_gpu/profiles/
│  ├─ manager.py                   (save/load profiles)
│  └─ presets.yaml                 (built-in profiles)
├─ setup.py                        (Python package setup)
└─ requirements.txt                (dependencies)

Features implemented:
✓ Detect all GPUs
✓ Show available methods
✓ Interactive method selection menu
✓ Multi-GPU configuration
✓ VM-GPU assignment
✓ Profile save/load
✓ Basic mode switching
✓ Configuration validation
```

### Phase 2: GTK Desktop GUI (Next Week)

```
Build timeline: 8-10 hours

Files to create:
├─ fortmox_gpu/ui/gui_gtk.py       (GTK interface)
├─ fortmox_gpu/ui/widgets/
│  ├─ gpu_list.py                  (GPU list widget)
│  ├─ method_selector.py           (method selection)
│  ├─ vm_assignment.py             (GPU-VM matrix)
│  └─ status_monitor.py            (status display)
├─ fortmox_gpu/ui/dialogs/
│  ├─ setup_wizard.py              (multi-step wizard)
│  ├─ gpu_config_dialog.py         (GPU configuration)
│  └─ vm_assignment_dialog.py      (VM assignment)
└─ fortmox_gpu/resources/
   ├─ icons/                       (application icons)
   └─ themes/                      (UI themes)

Features implemented:
✓ Visual GPU list with status
✓ Drag-and-drop GPU assignment
✓ Setup wizard (multi-step)
✓ Real-time configuration preview
✓ Easy method switching
✓ Configuration import/export
✓ Help documentation
✓ System tray integration
```

### Phase 3: Web UI + Future Features

```
Build timeline: 6-8 hours + ongoing

Files to create:
├─ fortmox_gpu/ui/web_ui.py        (Flask app)
├─ fortmox_gpu/api/
│  ├─ gpu_api.py                   (REST API)
│  └─ auth.py                      (authentication)
├─ fortmox_gpu/monitoring/          (FUTURE)
│  ├─ gpu_monitor.py
│  ├─ thermal_monitor.py
│  └─ metrics_exporter.py
├─ fortmox_gpu/automation/          (FUTURE)
│  ├─ profile_auto_selector.py
│  └─ smart_allocator.py
└─ web/                             (Web UI assets)
   ├─ static/
   │  ├─ css/
   │  ├─ js/
   │  └─ images/
   └─ templates/
      ├─ dashboard.html
      ├─ config.html
      └─ status.html

Features implemented:
✓ Browser-based dashboard
✓ Real-time status monitoring
✓ Configuration management
✓ Remote GPU monitoring
✓ Mobile-friendly design
□ Performance metrics (FUTURE)
□ Auto-profile selection (FUTURE)
□ Prometheus integration (FUTURE)
```

---

## File Reference Guide

### Tonight: Files You'll Use
```
scripts/
└─ gpu-setup.sh                    ← Use this tonight for FortMox

config/
├─ system.yaml                     ← Main FortMox config
└─ vm-templates/                   ← VM definitions
   ├─ gaming-vm.yaml
   ├─ clean-vm.yaml
   ├─ research-vm.yaml
   └─ opnsense-microvm.yaml
```

### This Week: Files We'll Create
```
fortmox_gpu/                       ← NEW Python package
├─ __init__.py                     (package definition)
├─ cli.py                          (command entry point)
├─ core/                           (core functionality)
│  ├─ gpu_detector.py              (GPU detection)
│  ├─ gpu_manager.py               (GPU management)
│  ├─ vm_manager.py                (VM GPU assignment)
│  └─ config_parser.py             (YAML parsing)
├─ ui/                             (user interfaces)
│  ├─ tui_menu.py                  (Terminal UI - Phase 1)
│  ├─ gui_gtk.py                   (Desktop GUI - Phase 2)
│  └─ web_ui.py                    (Web UI - Phase 3)
├─ profiles/                       (configuration profiles)
│  ├─ manager.py
│  └─ presets.yaml
├─ utils/                          (utilities)
│  ├─ logger.py
│  └─ helpers.py
└─ monitoring/                     (monitoring - FUTURE)

Documentation:
├─ docs/planning/gpu-switching-analysis.md       (what we analyzed)
├─ docs/planning/gpu-enhancement-plan.md         (how we'll build)
├─ docs/planning/implementation-roadmap.md       (this file)
├─ GPU-PYTHON-GUIDE.md             (NEW - Python version docs)
├─ GPU-LEGACY-GUIDE.md             (NEW - Bash version docs)
└─ GPU-MIGRATION-GUIDE.md          (NEW - How to upgrade)

Installation:
├─ setup.py                        (NEW - Python packaging)
└─ requirements.txt                (NEW - dependencies)
```

### Installation Commands (After Building)
```bash
# Install Python GPU tool
cd /path/to/FortMoxDesktop\ Workstation
pip install -e .

# Use it
fortmox-gpu-setup          # Interactive TUI setup
fortmox-gpu-gui            # Launch GTK GUI
fortmox-gpu-config         # Configuration tool
fortmox-gpu-switch         # Mode switcher
fortmox-gpu-profile        # Profile manager
```

---

## Success Criteria

### Tonight (Bash Version)
- [ ] FortMox installs successfully
- [ ] Current gpu-setup.sh works without errors
- [ ] Desktop Proxmox is running
- [ ] Laptop Proxmox is running
- [ ] VMs can be created

### This Week (Python Phase 1)
- [ ] Python GPU package created
- [ ] TUI menu system working
- [ ] GPU detector finds all GPUs
- [ ] Multi-GPU configuration works
- [ ] Profile save/load functions
- [ ] Mode switching available
- [ ] Documentation complete

### Next Week (Python Phase 2)
- [ ] GTK Desktop GUI launches
- [ ] GPU setup wizard functional
- [ ] Drag-and-drop assignment works
- [ ] Real-time status display
- [ ] Configuration preview/export
- [ ] Integration tests pass

### Future Phases
- [ ] Web UI accessible
- [ ] Monitoring data collected
- [ ] Prometheus integration working
- [ ] Automated profile selection
- [ ] Performance metrics displayed

---

## Quick Reference: Commands

### Once Everything is Built

```bash
# Use Python GPU tool (new, this week)
sudo fortmox-gpu-setup      # Start interactive setup wizard
sudo fortmox-gpu-gui        # Launch GTK graphical interface

# Use bash GPU tool (current, tonight)
sudo ./scripts/gpu-setup.sh # Use existing bash version

# Configuration
fortmox-gpu-config list     # List all GPUs
fortmox-gpu-config show     # Show current config
fortmox-gpu-config reset    # Reset to defaults

# Profiles
fortmox-gpu-profile list    # List saved profiles
fortmox-gpu-profile load    # Load a profile
fortmox-gpu-profile save    # Save current config

# Mode switching
fortmox-gpu-switch          # Interactive mode switcher
fortmox-gpu-switch full     # Switch to full passthrough
fortmox-gpu-switch sr-iov   # Switch to SR-IOV mode
fortmox-gpu-switch virtio   # Switch to virtual GPU

# Monitoring (future)
fortmox-gpu-monitor         # Real-time GPU monitoring
fortmox-gpu-metrics export  # Export metrics to Prometheus
```

---

## Next Steps

### Right Now (Next 30 minutes)
1. Review this roadmap document
2. Confirm everything matches your needs
3. Ask any clarification questions

### Tonight (Next 2-3 hours)
1. Install Proxmox on Desktop
2. Install Proxmox on Laptop
3. Deploy FortMox using current scripts
4. Verify everything working

### This Week (Phase 1 - Tuesday/Wednesday)
1. Build Python GPU package
2. Implement TUI menu system
3. Create GPU detector and manager
4. Test with your systems
5. Document Python version

### Next Week (Phase 2 - Thursday/Friday)
1. Build GTK Desktop GUI
2. Create setup wizard
3. Integration with TUI
4. Comprehensive testing
5. Final documentation

---

## Questions to Verify

Before we proceed, please confirm:

1. ✅ Do these features match what you need?
2. ✅ Timeline okay (bash tonight, Python this week)?
3. ✅ TUI + GTK Desktop are the priority UIs?
4. ✅ Keep bash script unchanged as fallback?
5. ✅ Want monitoring/automation for future versions?

If everything looks good, we proceed with:
- **Tonight:** FortMox installation with current gpu-setup.sh
- **This Week:** Build Python GPU tool with TUI + GTK GUI
- **Future:** Add monitoring, automation, and web UI

---

**Document Status:** Ready for Implementation
**Created:** November 17, 2025
**Next Action:** Confirmation from you, then start building!

