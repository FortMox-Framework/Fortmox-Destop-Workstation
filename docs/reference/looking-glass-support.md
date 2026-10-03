# Looking Glass Integration Status

> **Alpha status:** Looking Glass is not integrated by the current CLI. The Gaming VM YAML contains Looking Glass, ivshmem, GPU, and SPICE intent, but VM creation currently maps only basic VMID/name/CPU/memory/disk values and a default network. No Looking Glass behavior or performance has been verified by FortMox.

## Intended Use

## What is Looking Glass?

Looking Glass is a shared-memory display project commonly used with GPU-passthrough VMs. This repository does not validate compatibility, setup, latency, or gaming performance.

---

## Looking Glass in FortMox: Current Status

### Settings Declared in the Gaming Template

These values document intended settings only; the current `fortmox vm create` command does not apply them.

**In Gaming VM Template (`config/vm-templates/gaming-vm.yaml`):**

```yaml
# Lines 35-43: Looking Glass configuration
looking_glass:
  enabled: true
  # Looking Glass shared memory size (in MB)
  # This value is declarative; VM creation does not apply or validate it.
  # 32MB minimum, 128MB+ recommended for high resolution
  shared_memory: 256
  # ivshmem device configuration
  ivshmem: true
```

**In Gaming VM Display Config (lines 31-33):**

```yaml
display:
  type: spice                    # Better for gaming than VNC
  memory: 256M                   # SPICE shared memory
```

**Package names listed in the template (not installed by the CLI):**

```yaml
packages_windows:
  - looking-glass-client        # Ultra-low latency remote gaming

packages_linux:
  - looking-glass-client        # Ultra-low latency remote gaming
```

---

## Example Architecture (Not Deployed by FortMox)

### Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Proxmox Host                          │
│                                                         │
│  ┌───────────────────────────────────────────────┐    │
│  │          Gaming VM (VMID: 101)                │    │
│  │                                               │    │
│  │  Guest OS (Windows/Linux)                    │    │
│  │  ├─ GPU Passthrough (full direct access)    │    │
│  │  ├─ Looking Glass Host Service              │    │
│  │  ├─ ivshmem Device (shared memory)          │    │
│  │  └─ SPICE Agent (input/audio)               │    │
│  │                                               │    │
│  │  Performance depends on host and guest setup │    │
│  └───────────────────────────────────────────────┘    │
│              ▲        ▲        ▲                       │
│              │        │        │                       │
│           GPU      Memory     SPICE                    │
│              │        │        │                       │
│  ┌───────────┴────────┴────────┴──────────────┐       │
│  │      Proxmox QEMU/KVM Hypervisor          │       │
│  │      • ivshmem support                    │       │
│  │      • SPICE display server               │       │
│  │      • Shared memory management           │       │
│  └───────────────────────────────────────────┘       │
│              ▲        ▲        ▲                       │
│              │        │        │                       │
│           GPU      Memory    Display                   │
│              │        │        │                       │
│  ┌───────────┴────────┴────────┴──────────────┐       │
│  │    Looking Glass Client (Proxmox Host)    │       │
│  │    • Reads shared memory (ivshmem)        │       │
│  │    • Decodes video frames                 │       │
│  │    • Latency not measured by FortMox        │       │
│  │    • Fullscreen or windowed               │       │
│  └───────────────────────────────────────────┘       │
│              ▲                                        │
│              │ Display output                       │
│              │ (Local or network)                    │
│  ┌───────────┴──────────────────────────────┐        │
│  │         User Interface                   │        │
│  │         (Monitor/Screen)                 │        │
│  │                                          │        │
│  │  Keyboard/Mouse: via SPICE              │        │
│  │  Audio: via SPICE                       │        │
│  │  Video: via Looking Glass               │        │
│  └──────────────────────────────────────────┘        │
│                                                       │
└─────────────────────────────────────────────────────┘
```

### Data Flow

```
1. Gaming in VM
   ↓
2. GPU renders frame to VRAM
   ↓
3. Looking Glass Host copies frame to ivshmem (shared memory)
   ↓
4. Looking Glass Client reads from ivshmem
   ↓
5. Decodes and displays on host screen
   ↓
6. Result: Display output; performance depends on the complete configuration
```

---

## Performance

Latency and rendering performance depend on the GPU, host/guest software, shared-memory configuration, display, and workload. FortMox has no comparative measurements or guarantees.

---

## Current Looking Glass Support in FortMox

### Not Implemented by FortMox

The current VM creation command does not configure GPU passthrough, ivshmem, SPICE, Looking Glass software, or guest packages. The template values are declarative only. Follow current upstream Looking Glass and Proxmox documentation, and validate the entire setup independently.

---

## Manual Setup Boundary

### After FortMox Installation

```bash
# Reviewing a dry-run does not configure Looking Glass or create the full VM setup.
./fortmox vm create gaming --dry-run

# 2. Install Guest OS in Gaming VM
# Boot VM, install Windows 11 or Ubuntu

# 3. Install Looking Glass in Guest
# For Windows: Follow looking-glass-setup.md instructions
# For Linux: apt-get install looking-glass-client

# 4. Launch Looking Glass Client on Proxmox Host
looking-glass-client spice://localhost:5900
# or
looking-glass-client -h [VM_IP] -p 5900
```

---

## Step-by-Step: Enable Looking Glass on Gaming VM

### On Proxmox Host

```bash
# 1. SSH into Proxmox host
ssh root@proxmox-desktop

# 2. Verify ivshmem device
pveam list                    # List available devices
lsmod | grep ivshmem         # Check kernel module

# 3. Gaming VM should already have ivshmem enabled
# (configured in gaming-vm.yaml)

# 4. Check VM config
cat /etc/pve/qemu-server/101.conf | grep -i "ivshmem\|spice"

# Should show:
# spice: ...
# ivshmem: ...
```

### On Gaming VM (Windows)

```powershell
# 1. Download Looking Glass Host
# https://github.com/looking-glass/looking-glass/releases
# Download: looking-glass-host-setup.exe

# 2. Run as Administrator
.\looking-glass-host-setup.exe

# 3. Follow wizard, install Looking Glass Host service

# 4. Verify service is running
Get-Service -Name "Looking Glass Host"
# Should show: Running
```

### On Gaming VM (Linux - Ubuntu/Fedora)

```bash
# 1. Install Looking Glass Client
sudo apt-get install looking-glass-client    # Debian/Ubuntu
# OR
sudo dnf install looking-glass-client        # Fedora/RHEL

# 2. Start Looking Glass Host
# For Ubuntu/Debian:
sudo systemctl start looking-glass-host
sudo systemctl enable looking-glass-host

# 3. Verify
looking-glass-client -h <gaming-vm-ip>
```

### On Proxmox Host

```bash
# 1. Install Looking Glass Client
apt-get install looking-glass-client

# 2. Get Gaming VM IP (from inside the VM)
# Then run:
looking-glass-client spice://[GAMING_VM_IP]:5900

# 3. Alternatively, use localhost if accessing locally
looking-glass-client spice://localhost:5900

# 4. Fullscreen mode
# Press: Ctrl+Alt+F to toggle fullscreen
```

---

## Shared Memory Size Recommendations

```yaml
looking_glass:
  shared_memory: XXX  # Adjust based on your needs

Recommendations:
├─ 32MB   : 720p/60fps minimal
├─ 64MB   : 1080p/60fps comfortable
├─ 128MB  : 1440p/144fps recommended (default for many)
├─ 256MB  : 4K/60fps or 1440p/240fps (FortMox default)
└─ 512MB+ : High resolution + low latency priority
```

**In FortMox Gaming VM:** 256MB (suitable for high-res/high-refresh gaming)

---

## Performance Tuning

### Guest OS Optimization (Linux)

```bash
# Install performance tools
sudo apt-get install looking-glass-client libvulkan1 mesa-vulkan-drivers

# Run Looking Glass with optimization flags
looking-glass-client \
  --fps=0 \
  --spice-disable-audio \
  --disable-help-text \
  -f 60              # Cap at 60 FPS (adjust as needed)
```

### Guest OS Optimization (Windows)

```powershell
# In Looking Glass Host settings:
# Settings → Performance
# - Enable Performance Mode
# - Disable Cursor Rendering (use host cursor)
# - Enable Memory Pooling
# - Disable Frame Serialization
```

### Proxmox Host Optimization

```bash
# Edit VM config
nano /etc/pve/qemu-server/101.conf

# Add optimizations:
# CPU: host or high passthrough
# Memory: no ballooning (balloon: 0)
# Disk cache: writethrough or writeback
```

---

## Troubleshooting Looking Glass

### Issue: "ivshmem device not found"

**Solution:**
```bash
# Check if device is enabled in VM config
cat /etc/pve/qemu-server/101.conf | grep ivshmem

# If missing, add it:
# ivshmem: looking-glass:256
```

### Issue: High latency

**Solution:**
1. Independently verify GPU assignment and guest driver status; FortMox defines no performance threshold.
2. Increase shared memory size
3. Use local display (not network)
4. Check CPU load on Proxmox host

### Issue: "Looking Glass client won't connect"

**Solution:**
```bash
# Verify host service running
# On Windows:
Get-Service "Looking Glass Host"

# On Linux:
sudo systemctl status looking-glass-host

# Verify SPICE port is open
netstat -tlnp | grep -i spice

# Test connection
looking-glass-client spice://localhost:5900 -v
```

### Issue: "No video/black screen"

**Solution:**
1. Verify GPU passthrough is working
2. Check guest drivers installed
3. Restart Looking Glass service
4. Check for driver crashes in Event Viewer (Windows)

---

## Looking Glass Integration with Python GPU Tool

### Future Enhancement (Phase 3+)

When we build the Python GPU tool, we'll add Looking Glass management:

```python
# New module: fortmox_gpu/looking_glass/

class LookingGlassManager:
    """Manage Looking Glass configuration and connection"""

    def auto_detect_gaming_vm(self):
        """Find gaming VM with GPU passthrough"""
        # Detect Gaming VM (VMID 101)
        # Verify GPU passthrough enabled
        # Check ivshmem configured

    def get_connection_string(self, vm_id):
        """Generate Looking Glass connection string"""
        # Get VM IP or use localhost
        # Generate SPICE URI
        # Return connection command

    def optimize_performance(self, vm_id):
        """Auto-optimize Looking Glass settings"""
        # Detect resolution of host display
        # Suggest shared memory size
        # Set recommended FPS cap
        # Configure codec

    def launch_client(self, vm_id, fullscreen=False):
        """Launch Looking Glass client"""
        # Build command with optimizations
        # Auto-detect available codec
        # Launch with appropriate flags
```

### Not Available in FortMox

There are no `fortmox-gpu-lg-*` commands, GUI controls, automatic launch, or metrics integration in the current project.

### GUI Integration (Phase 2+)

The GTK Desktop GUI will include:
- Looking Glass status indicator
- One-click launch button
- Performance metrics display
- Connection management
- Configuration panel

---

## Performance

FortMox has not measured Looking Glass latency, frame rate, CPU usage, or resolution limits. Results depend on the GPU, host and guest software, shared-memory settings, and workload.

---

## Current CLI Integration

Looking Glass configuration is not applied by FortMox. This reference and the linked setup guide are manual guidance; they do not certify the Gaming VM template, guest setup, SPICE, ivshmem, latency, or guest OS compatibility.

### Possible Future Work (Not Committed)

- Automatic launch from GUI
- Performance auto-optimization
- Real-time latency monitoring
- Web dashboard integration
- Prometheus metrics export
- Mobile client support

---

## Commands Summary

### Manual Setup

Use current upstream Looking Glass and Proxmox documentation. The commands in older copies of this reference are not a validated or complete installation procedure.

### Configuration Files

```bash
# Example Proxmox VM path; not created or auto-configured by FortMox
/etc/pve/qemu-server/101.conf                # Gaming VM (vmid: 101)

# FortMox template values are declarative only
config/vm-templates/gaming-vm.yaml           # Shared memory size, ivshmem
```

---

## Documentation Available

**In FortMox:**
- `docs/guides/looking-glass-setup.md` (453 lines)
  - Complete setup guide
  - Windows installation
  - Linux installation
  - Troubleshooting
  - Performance tuning

**Online:**
- GitHub: https://github.com/looking-glass/looking-glass
- Website: https://looking-glass.io
- Documentation: https://looking-glass.io/docs

---

## Your Looking Glass Setup for Tonight

### Pre-deployment Check

```bash
# Verify Proxmox supports Looking Glass
dpkg -l | grep qemu-system
# Should be version 5.2+ (supports ivshmem)
```

### After FortMox Deploy

```bash
# 1. Gaming VM will auto-have ivshmem configured
# 2. Install Guest OS (Windows 11 or Ubuntu)
# 3. Install Looking Glass Host in guest
# 4. Launch Looking Glass Client from Proxmox host
# 5. Play games with <5ms latency!
```

---

## Summary: Looking Glass in FortMox

| Feature | Status | Details |
|---------|--------|---------|
| **ivshmem Support** | ✅ Ready | 256MB default |
| **Gaming VM Template** | ✅ Pre-configured | VMID 101 |
| **SPICE Display** | ✅ Enabled | Audio + input |
| **Host Package** | ✅ Documentted | Easy install |
| **Guest Package** | ✅ Included | Windows + Linux |
| **Setup Guide** | ✅ Complete | 453 lines |
| **Auto-launch** | 🔜 Planned | Phase 3+ |
| **GUI Integration** | 🔜 Planned | Phase 2+ |
| **Monitoring** | 🔜 Planned | Phase 3+ |

---

## Next Steps

### Tonight
1. Deploy Gaming VM with FortMox
2. Install Guest OS
3. Install Looking Glass Host
4. Test Looking Glass Client
5. Enjoy ultra-low latency gaming!

### This Week
Continue with Python GPU tool development (doesn't affect Looking Glass)

### Next Week
Add Looking Glass management to Python GUI

---

**Looking Glass status:** Not integrated by the current CLI. See [TODO.md](../../TODO.md) for the implementation and validation work.

