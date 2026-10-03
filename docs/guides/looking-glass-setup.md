# Looking Glass Setup Guide

> **Alpha/manual guide:** These instructions are not integrated or validated by the FortMox CLI. VM creation currently does not apply the Gaming template's Looking Glass, ivshmem, SPICE, or GPU settings.

## Overview

**Looking Glass** is an ultra-low latency remote display solution for QEMU/KVM VMs. It provides nearly native gaming performance when accessing a GPU-passthrough VM over the network from your Proxmox host.

### Key Features
- **Shared-memory display**: Uses a configured shared-memory device between host and guest
- **Performance**: Depends on GPU, guest, host, and display configuration; not benchmarked here
- **Remote access**: Access gaming VM from any network device
- **Secure**: Runs locally on same machine or trusted network
- **Cross-platform**: Windows and Linux guest support

## Architecture

```
┌─────────────────────────────────────────────┐
│         Proxmox Host (Laptop/Desktop)        │
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │ Gaming VM (Windows or Linux)        │   │
│  │                                     │   │
│  │  • GPU Passthrough                  │   │
│  │  • Looking Glass Server             │   │
│  │  • Shared Memory (ivshmem)          │   │
│  └─────────────────────────────────────┘   │
│            ▲      ▲      ▲                  │
│            │      │      │                  │
│            GPU   Memory  SPICE              │
│            │      │      │                  │
│  ┌─────────┴──────┴──────┴─────────────┐   │
│  │  Looking Glass Client               │   │
│  │  (SPICE Display)                    │   │
│  │  Shared Memory (ivshmem)            │   │
│  │  Direct Memory Access               │   │
│  └─────────────────────────────────────┘   │
│            ▲                                │
│            │ Ultra-low latency display    │
│            │ (local or network)           │
│  ┌─────────┴──────────────────────────┐   │
│  │  User Interface                    │   │
│  │  (Monitor/Laptop Screen)           │   │
│  └────────────────────────────────────┘   │
│                                             │
└─────────────────────────────────────────────┘
```

## Installation

### 1. Proxmox Host Setup

Ensure the gaming VM template has Looking Glass enabled:

```yaml
looking_glass:
  enabled: true
  shared_memory: 256  # MB
  ivshmem: true
```

### 2. Install Looking Glass Host Module

On the **Proxmox host**, install Looking Glass support:

```bash
# For QEMU/KVM support (usually pre-installed)
apt-get install -y looking-glass-host-qemu

# Verify
dpkg -l | grep looking-glass
```

### 3. Configure Shared Memory

The current VM creation command does not configure the ivshmem device. Any setup must be performed and verified manually:

```xml
<!-- In VM XML configuration -->
<ivshmem name='looking-glass'>
  <model type='ivshmem'/>
  <address type='pci' domain='0x0000' bus='0x00' slot='0x14' function='0x0'/>
</ivshmem>
```

The `gaming-vm.yaml` template contains settings that the current VM creation command does not apply. Shared-memory sizing must be selected and configured manually:
- **32MB minimum** - Low resolution/quality
- **64MB** - 1080p comfortable
- **128MB** - 1440p recommended
- **256MB+** - 4K or high-refresh gaming

## Guest OS Setup

### Windows Guest (Recommended for Gaming)

#### 1. Install Looking Glass Guest Agent

```powershell
# Download from GitHub releases or use SPICE drivers
# https://github.com/looking-glass/looking-glass/releases

# Extract and run installer:
.\looking-glass-host-setup.exe
```

Or install via Windows package manager:

```powershell
choco install looking-glass-client
```

#### 2. Enable ivshmem Device

1. Open Device Manager
2. Look for "PCI Device" with unknown driver
3. Install driver from Looking Glass release package
4. Or run: `devcon install looking-glass-ivshmem.inf`

#### 3. Install SPICE Guest Tools

```powershell
# Enables clipboard, file sharing, input passthrough
choco install spice-guest-tools
```

### Linux Guest

#### 1. Install Looking Glass Host Package

```bash
# Ubuntu/Debian
sudo apt-get install -y looking-glass-client looking-glass-host-qemu

# Fedora/RHEL
sudo dnf install -y looking-glass-client looking-glass-host-qemu
```

#### 2. Load Kernel Module

```bash
# Load ivshmem module
sudo modprobe kvmfr

# Make permanent
echo "kvmfr" | sudo tee -a /etc/modules-load.d/kvmfr.conf
```

#### 3. Create Shared Memory Device

```bash
# Setup huge pages for better performance
echo 256 | sudo tee /proc/sys/vm/nr_hugepages

# Make permanent
echo "vm.nr_hugepages = 256" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

#### 4. Setup Device Access

```bash
# Create udev rule for Looking Glass device
sudo cat > /etc/udev/rules.d/99-looking-glass.rules << 'EOF'
SUBSYSTEM=="kvmfr", GROUP="kvm", MODE="0660"
EOF

# Apply rules
sudo udevadm control --reload
```

## Client-Side Setup

### Install Looking Glass Client

#### On Linux Host

```bash
# Build from source (recommended for latest features)
git clone https://github.com/looking-glass/looking-glass.git
cd looking-glass/client
mkdir build && cd build
cmake -DCMAKE_BUILD_TYPE=Release ..
make
sudo make install

# Or install from package manager
sudo apt-get install -y looking-glass-client  # Ubuntu/Debian
sudo dnf install -y looking-glass-client      # Fedora/RHEL
```

#### On Windows Host

```powershell
# Download and install
# https://github.com/looking-glass/looking-glass/releases/download/B6/looking-glass-B6-windows.zip

# Or use Windows package manager
choco install looking-glass-client
```

#### On macOS

```bash
brew install looking-glass-client
```

### Running Looking Glass Client

#### Basic Usage

```bash
# Connect to local SPICE server
looking-glass-client

# Connect with custom host/port
looking-glass-client -h 127.0.0.1 -p 5900

# Full screen mode
looking-glass-client -f

# Specify frame buffer size
looking-glass-client -s 256M
```

#### Common Options

```bash
looking-glass-client \
  --host=127.0.0.1 \
  --port=5900 \
  --framebuffer-memory=256 \
  --fullscreen=yes \
  --mouse-speed=100 \
  --vsync=yes
```

#### Configuration File

Create `~/.looking-glass/client.ini`:

```ini
[spice]
host=127.0.0.1
port=5900

[app]
fullscreen=yes
mouse-speed=100
vsync=yes
fps=0  # Unlimited

[display]
rotation=0

[performance]
memoryEfficient=no
```

## Network Access (Optional)

For accessing the gaming VM from another machine on the network:

### SPICE Network Configuration

1. **Expose SPICE port** on Proxmox host:
   ```bash
   # Forward port 5900 to gaming VM SPICE
   # (Configure via Proxmox UI or manually)
   ```

2. **Connect from remote machine**:
   ```bash
   looking-glass-client -h <proxmox-host-ip> -p 5900
   ```

### Secure Remote Access (VPN/SSH)

For security-conscious deployments:

```bash
# Use SSH tunnel instead of direct SPICE
ssh -L 5900:127.0.0.1:5900 user@proxmox-host

# Then connect client
looking-glass-client -h 127.0.0.1 -p 5900
```

## Performance Tuning

### Shared Memory Size

Larger = Better quality but more RAM usage:

```bash
# View current allocation
cat /proc/meminfo | grep Huge

# Allocate huge pages
echo 512 | sudo tee /proc/sys/vm/nr_hugepages
```

### VM Configuration

For optimal performance:

```yaml
# In gaming-vm.yaml
looking_glass:
  enabled: true
  shared_memory: 256-512  # Higher = better quality

specs:
  memory:
    max: 16384  # Adequate for gaming + Looking Glass
```

### CPU Pinning

Ensure Looking Glass thread runs on dedicated core:

```yaml
security:
  pin_vcpu: true
  cpu_cores_pinned: [0, 1, 2, 3, 4, 5]  # Pin 6 cores
```

### Network Optimization

For network-based access:

```bash
# Increase TCP buffers
sudo sysctl -w net.core.rmem_max=134217728
sudo sysctl -w net.core.wmem_max=134217728
```

## Troubleshooting

### "SPICE server not ready" Error

1. **Verify VM is running**: Check Proxmox interface
2. **Check SPICE listener**: `netstat -tuln | grep 5900`
3. **Restart VM**: Power off and on

```bash
# Check VM status
qm status 101  # or your VM ID

# Check SPICE port
sudo netstat -tuln | grep LISTEN
```

### Low Frame Rate

1. **Increase shared memory**: 128MB → 256MB or 512MB
2. **Enable huge pages**: See above
3. **Reduce VM load**: Close other applications
4. **Check CPU pinning**: Ensure VCPUs aren't over-subscribed

### High Latency

1. **Check network**: Use local connection, not remote
2. **Verify SPICE binding**: Should be localhost or host IP
3. **Disable vsync**: May help in some cases
4. **Check GPU load**: Monitor with nvidia-smi or radeon-top

### Graphics Glitches/Artifacts

1. **Update Looking Glass**: Both client and host
2. **Verify driver compatibility**: GPU driver version
3. **Check frame buffer**: May need larger shared memory
4. **Restart SPICE server**: Restart VM

## Comparison: Display Methods

| Method | Notes |
|--------|-------|
| **Looking Glass, SPICE, VNC, RDP, physical display** | Latency, performance, and security depend on configuration and workload; FortMox has not compared these methods. |

## Security Considerations

### Local Use (Recommended)

Looking Glass is designed for local network use:
- Ultra-low latency
- Minimal network overhead
- Suitable for same physical machine

```bash
# Local use (safest)
looking-glass-client -h 127.0.0.1
```

### Remote Access (Caution)

For remote access, use VPN/SSH:

```bash
# VPN-based access (better)
ssh -L 5900:127.0.0.1:5900 user@proxmox-host
looking-glass-client

# TLS support
looking-glass-client --tls -h remote.host.com
```

### Disable SPICE in Production

If not using Looking Glass, disable SPICE:

```yaml
display:
  type: none  # Headless
  # Only use GPU directly or HDMI output
```

## Further Resources

- **Official Documentation**: https://looking-glass.io/
- **GitHub Repository**: https://github.com/looking-glass/looking-glass
- **Community Forum**: https://forum.looking-glass.io/

## Integration with Proxmox

### Proxmox UI Support

Looking Glass can be integrated with Proxmox UI:

1. Install Looking Glass console plugin
2. Access directly from VM console dropdown
3. No need for external client

```bash
# Install Proxmox plugin (if available)
apt-get install -y proxmox-looking-glass
```

### Automation

Enable Looking Glass in deployment:

```bash
# In deployment script
qm set <vmid> -vga none  # Disable default VGA
qm set <vmid> -ivshmem "looking-glass" # Add ivshmem device
```

---

**Note**: Looking Glass is ideal for gaming VMs on laptops and desktops where ultra-low latency is important. For other VMs, SPICE or VNC may be sufficient.
