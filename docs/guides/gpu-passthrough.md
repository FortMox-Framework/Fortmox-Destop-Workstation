# GPU Passthrough & Acceleration

> **Alpha/manual guide:** The CLI only reports a GPU strategy. It does not bind VFIO devices, configure SR-IOV/vGPU, or attach GPUs to VMs. Treat all commands and performance figures here as unverified guidance until tested on supported hardware.

## Overview

This guide describes possible GPU acceleration methods on Proxmox VE, from virtIO to full passthrough; it is not an implementation checklist or compatibility guarantee.

## GPU Methods Comparison

| Method | Availability | FortMox support |
|--------|--------------|-----------------|
| **virtIO** | Depends on the virtual display and guest configuration | Strategy planning only; no VM changes |
| **SR-IOV** | Requires compatible GPU, firmware, drivers, and software | Strategy planning only; no device configuration |
| **vGPU** | Requires compatible hardware and host/guest software; licensing may apply | Strategy planning only; no device configuration |
| **Full Passthrough** | Requires suitable IOMMU groups and host/guest configuration | Strategy planning only; no VFIO binding or VM assignment |

Performance, compatibility, isolation, and latency vary by system. FortMox has not benchmarked or certified any GPU method.

## Prerequisites

### Hardware Requirements

```bash
# Check CPU supports IOMMU
grep -i "flags.*iommu\|flags.*vt-d\|flags.*amd-vi" /proc/cpuinfo

# For Intel:
grep "vmx" /proc/cpuinfo

# For AMD:
grep "svm" /proc/cpuinfo

# List GPU
lspci | grep -i "vga\|3d\|display"

# Check GPU type
lspci -nnk | grep -i "vga\|3d" -A3
```

### BIOS/UEFI Configuration

1. **Enable Virtualization** (VT-x/SVM)
2. **Enable IOMMU** (VT-d for Intel, AMD-Vi for AMD)
3. **Disable Secure Boot** (if using passthrough)
4. **Disable Fast Boot** (for consistency)

## Method 1: virtIO GPU (Universal, Lightweight)

### Use Cases
- Daily work (web, office, coding)
- Research VMs (no games)
- Multiple VMs sharing GPU
- Bandwidth-limited networks

### Setup

#### 1.1 Add GPU Device to VM

```bash
# SSH into Proxmox host
ssh root@proxmox

# Get GPU ID (example: 0000:01:00.0)
lspci | grep VGA

# Edit VM configuration
nano /etc/pve/qemu-server/100.conf  # Change 100 to your VM ID

# Add to configuration:
vga: qxl
```

#### 1.2 Install Guest Drivers

**Windows Guest:**

```bash
# Download QXL drivers from
# https://www.spice-space.org/download/windows/

# Install drivers
# Run as Administrator
# May need to install SPICE guest tools separately
```

**Linux Guest:**

```bash
# Debian/Ubuntu
sudo apt install xserver-xorg-video-qxl

# Arch
sudo pacman -S xf86-video-qxl

# Restart X server
sudo systemctl restart gdm3  # or lightdm, sddm, etc.
```

#### 1.3 Verify

```bash
# In guest VM
glxinfo | grep "direct rendering"
glxgears  # Test graphics
```

### Performance Optimization

Edit `/etc/pve/qemu-server/100.conf`:

```
vga: qxl
# Add more memory for QXL
# virtio-gpu-mem: 256M
scsihw: virtio-scsi-single
bus: pcie
```

## Method 2: SR-IOV GPU Passthrough

### Supported GPUs

**NVIDIA**:
- RTX A series (professional)
- RTX 30/40 series with proper drivers (consumer)
- GeForce RTX with newest drivers

**AMD**:
- RDNA 2+ with SR-IOV support
- Pro Radeon (professional)

### Setup

#### 2.1 Enable IOMMU and SR-IOV in Kernel

Edit `/etc/default/grub`:

```bash
nano /etc/default/grub
```

Modify `GRUB_CMDLINE_LINUX_DEFAULT`:

```
# Intel
GRUB_CMDLINE_LINUX_DEFAULT="quiet intel_iommu=on iommu=pt"

# AMD
GRUB_CMDLINE_LINUX_DEFAULT="quiet amd_iommu=on iommu=pt"
```

Update GRUB:

```bash
sudo update-grub
sudo reboot
```

#### 2.2 Enable SR-IOV on GPU

```bash
# Find GPU
lspci -nnk | grep VGA

# Example: 02:00.0 VGA compatible controller
# Note the domain:bus:slot.function (02:00.0)

# Check if SR-IOV capable
lspci -vs 02:00.0 | grep SR-IOV

# Enable SR-IOV (enable N virtual functions)
echo 4 > /sys/class/pci_bus/0000:02/device/sriov_numvfs

# To persist across reboot, create:
nano /etc/modprobe.d/99-sriov.conf
```

Add to `/etc/modprobe.d/99-sriov.conf`:

```
# NVIDIA SR-IOV setup
# Uncomment if using NVIDIA SR-IOV
# options nvidia sriov_capable=1
# options nvidia NVreg_EnablePCIeP2P=1

# AMD SR-IOV (usually auto-enabled)
```

#### 2.3 Identify Virtual Functions

```bash
# After enabling, check virtual functions
lspci | grep "Virtual Function"

# You should see something like:
# 02:00.1 3D controller: ...
# 02:00.2 3D controller: ...
# etc.

# Get IDs
lspci -nn | grep "Virtual Function"
```

#### 2.4 Configure VM

Edit `/etc/pve/qemu-server/101.conf` (Gaming VM):

```
# Standard settings
cores: 6
memory: 16384
balloon: 0

# Add SR-IOV GPU device
# Use the virtual function address
hostpci0: 02:00.1
```

#### 2.5 Install Guest Drivers

```bash
# Download NVIDIA drivers that support SR-IOV
# Or AMD drivers with SR-IOV support

# Windows: Run installer from GPU manufacturer
# Linux:
sudo apt install nvidia-driver-550  # or appropriate version
sudo reboot
```

#### 2.6 Verify

```bash
# In VM
nvidia-smi  # or radeon equivalent
```

## Method 3: NVIDIA vGPU Passthrough

### Requirements

- NVIDIA GPU with vGPU support (RTX/Tesla line)
- NVIDIA vGPU license (enterprise paid)
- NVIDIA KVM driver

### Setup (Enterprise, Beyond Scope)

This method requires licensing. See NVIDIA documentation for:
- License server setup
- Driver installation with vGPU
- VM configuration

## Method 4: Full GPU Passthrough

### Use Cases

- Gaming workloads where a separately configured passthrough setup is appropriate; FortMox provides no performance guarantee.
- AI/ML training
- Heavy compute workloads
- Only 1 VM can use GPU at a time

### Supported GPUs

Nearly all discrete GPUs:
- NVIDIA GeForce/RTX
- AMD Radeon
- Intel Arc
- Professional GPUs

### Setup

#### 4.1 Prepare Host Kernel

Edit `/etc/default/grub`:

```bash
nano /etc/default/grub
```

Update `GRUB_CMDLINE_LINUX_DEFAULT`:

```
# Intel
GRUB_CMDLINE_LINUX_DEFAULT="quiet intel_iommu=on iommu=pt vfio-pci.ids=10de:1234,10de:1235 initcall_blacklist=sysfb_init"

# AMD
GRUB_CMDLINE_LINUX_DEFAULT="quiet amd_iommu=on iommu=pt vfio-pci.ids=10de:1234,10de:1235 initcall_blacklist=sysfb_init"

# Replace 10de:1234 with your GPU's vendor:device ID
```

Get GPU vendor:device IDs:

```bash
lspci -nn | grep -i "vga\|3d"
# Output example: 02:00.0 VGA compatible controller [0300]: NVIDIA Corporation ... [10de:2503]
```

#### 4.2 Update Initramfs

```bash
nano /etc/modprobe.d/99-vfio-blacklist.conf
```

Add:

```
# Blacklist nvidia/nouveau drivers on boot
blacklist nouveau
blacklist amdgpu
blacklist nvidia

# VFIO takes priority
options vfio-pci ids=10de:2503,10de:1234
```

Update kernel and reboot:

```bash
sudo update-grub
sudo update-initramfs -u -k all
sudo reboot
```

#### 4.3 Verify VFIO Binding

After reboot:

```bash
# Check if GPU is bound to VFIO
lspci -ks 02:00  # Replace with your GPU address
# Should show: Kernel driver in use: vfio-pci

# Or check via sysfs
ls -la /sys/bus/pci/devices/0000:02:00.0/driver/
# Should contain symlink to vfio-pci
```

#### 4.4 Configure VM

Edit `/etc/pve/qemu-server/101.conf`:

```
# VM Settings
cores: 8
memory: 16384
balloon: 0
machine: pc-i440fx-8.0

# CPU pinning (optional but recommended)
cpuunits: 1024
numa: 1
hostpci0: 02:00,pcie=on,x-vga=on,rombar=1

# For dual GPU setups or additional GPUs:
# hostpci1: 03:00,pcie=on

# VBIOS ROM (optional - for better compatibility)
# Download VBIOS from TechPowerUp GPU database
# hostpci0: 02:00,romfile=MyGPU.rom,pcie=on,x-vga=on
```

**Important Flags Explained:**

- `x-vga=on`: GPU is primary display (can use VGA output)
- `pcie=on`: Use PCIe slot (better for modern GPUs)
- `rombar=1`: Enable ROM bar (firmware/BIOS)
- `romfile`: Custom VBIOS (optional, helps with compatibility)

#### 4.5 Install Host GPU Drivers

**For host console display** (if you have a second GPU for host):

```bash
# NVIDIA (if host needs display)
sudo apt install nvidia-driver-550
sudo apt install nvidia-utils

# AMD
sudo apt install amdgpu-core
```

#### 4.6 Install Guest Drivers

**Windows Guest:**

```
1. Boot Windows VM
2. Device Manager may show unknown device
3. Download NVIDIA/AMD drivers for your GPU
4. Run installer as Administrator
5. Reboot guest
```

**Linux Guest:**

```bash
# Ubuntu/Debian NVIDIA
sudo apt update
sudo apt install nvidia-driver-550
sudo reboot

# Fedora NVIDIA
sudo dnf install nvidia-driver nvidia-utils

# AMD (open source)
sudo apt install amdgpu

# AMD (proprietary)
# Download from https://www.amd.com/en/support
```

#### 4.7 Verify GPU Passthrough

**From Host:**

```bash
# Confirm VFIO binding
lspci -ks 02:00
# Output: Kernel driver in use: vfio-pci

# Monitor GPU access
watch -n 1 "lspci -ks 02:00"
```

**From Guest:**

```bash
# Windows (PowerShell)
Get-WmiObject Win32_VideoController

# Linux
lspci -k | grep -i "vga\|3d" -A3
nvidia-smi  # or equivalent for AMD
glxgears    # Test 3D graphics
```

### Performance Tuning

#### CPU Pinning (Gaming VMs)

In VM config `/etc/pve/qemu-server/101.conf`:

```
# Pin specific CPU cores to VM
# Host has 16 cores (0-15), pin cores 8-15 to gaming VM
cores: 8
cpus: 8
# Advanced: Manual pinning
hostpci0: 02:00,pcie=on,x-vga=on

# Use numactl to bind
numactl --cpunodebind=1 --membind=1 qm start 101
```

#### Memory Optimization

```
# Disable memory balloon for gaming
balloon: 0

# Use huge pages for better performance
hugepages: 1024  # 1GB = 512 x 2MB pages

# Memory tuning
memory: 16384
```

#### Cache Optimization

Add to VM config:

```
cache: writeback  # Better for gaming (be careful with data loss)
# or cache: writethrough  # Safer but slower
```

### Troubleshooting

#### VM Won't Start

```bash
# Check VFIO binding
lspci -ks 02:00

# Check Proxmox task log
tail -f /var/log/pve/tasks

# Verify GPU isn't still in use by host
ps aux | grep -i nvidia
ps aux | grep -i amd
```

#### Black Screen on Boot

```
1. Check if VBIOS needs to be extracted:
   - Download from TechPowerUp GPU BIOS database
   - Add romfile=/path/to/bios.rom to hostpci config

2. Try booting with -DKMS for nouveau blacklist

3. Check VM console output for errors
```

#### Low Performance

```
1. Verify GPU is fully assigned (not shared)
2. Check CPU pinning is correct
3. Enable huge pages
4. Disable CPU migration

5. Monitor thermal throttling:
   nvidia-smi dmon
```

## Choosing the Right Method

```
User: "I want web, office, coding"
→ Use: virtIO (simple, works everywhere)

User: "I want gaming and heavy work"
→ Use: Full GPU Passthrough (best performance)

User: "I want multiple VMs with some GPU access"
→ Use: SR-IOV (if supported by your GPU)

User: "I want enterprise features"
→ Use: vGPU (if you have licensing)
```

## Configuration Template

### Gaming VM (100% GPU Performance)

```yaml
# config/vm-templates/gaming-vm.yaml
gpu_method: "auto"  # Will use full passthrough if available

# In /etc/pve/qemu-server/101.conf:
cores: 8
memory: 16384
balloon: 0
machine: pc-i440fx-8.0
hostpci0: 02:00,pcie=on,x-vga=on
cpu: host
hugepages: 1024
cache: writeback
```

### Clean VM (Light GPU Usage)

```yaml
# config/vm-templates/clean-vm.yaml
gpu_method: "virtio"

# In /etc/pve/qemu-server/100.conf:
cores: 4
memory: 8192
vga: qxl
scsihw: virtio-scsi-single
```

### Research VM (No GPU)

```yaml
# config/vm-templates/research-vm.yaml
gpu_method: "none"

# In /etc/pve/qemu-server/102.conf:
cores: 4
memory: 8192
vga: std  # Standard VGA, no acceleration
```

## Security Considerations

### Full GPU Passthrough Isolation

```bash
# IOMMU groups keep devices together
for d in /sys/kernel/iommu_groups/*/devices/*; do
    n=${d#*/iommu_groups/}; n=${n%/devices/*}
    printf 'IOMMU Group %s ' "$n"
    lspci -nns "${d##*/}"
done

# Ideally GPU is in its own IOMMU group for clean isolation
```

### Access Control

```bash
# Restrict GPU VM access from other VMs
# Via firewall rules in OPNsense (see docs/guides/opnsense-firewall-setup.md)

# Monitor GPU usage
watch -n 1 nvidia-smi  # In host or guest
```

## Monitoring GPU Performance

### Real-time Monitoring

```bash
# NVIDIA GPU stats
nvidia-smi dmon
# Output: memory, GPU utilization, temperature

# AMD GPU stats
rocm-smi --showuse
rocm-smi --showtemp
```

### Gaming Benchmarks

```bash
# Linux
glxgears
# Should show high FPS (60+)

# Windows (3DMark, GFXBench, UnigineHaven)
# Download from benchmarking sites
```

## Summary

- **virtIO, SR-IOV, vGPU, and full passthrough**: Availability and performance depend on the hardware and software stack; FortMox does not configure or benchmark these methods.

Choose a method only after checking current hardware/vendor documentation and independently validating the host and guest configuration.

---

**CLI status**: Strategy planning only; GPU configuration is not implemented.
