# Quick Reference Guide

## Essential Commands

### Deployment
```bash
# Load configuration and detect hardware
./scripts/config-manager.sh

# Deploy Proxmox configuration
./scripts/deploy.sh

# Verify installation
./scripts/verify.sh
```

### Proxmox VM Management
```bash
# Start VM
qm start <vmid>

# Stop VM
qm stop <vmid>

# Access console
qm terminal <vmid>

# List VMs
qm list

# Detailed VM info
qm config <vmid>
```

### Common VM IDs
| VM | ID | IP | Purpose |
|----|----|----|---------|
| OPNsense Firewall | 50 | 192.168.100.1 | Firewall/router |
| Clean VM | 100 | 192.168.100.100 | Daily work |
| Gaming VM | 101 | 192.168.100.101 | Gaming/GPU |
| Research VM | 102 | 192.168.100.102 | Malware analysis |
| Tools VM | 103 | 192.168.100.103 | Testing tools |

## Configuration Files

### Main Configuration
**File**: `config/system.yaml`
**Edit**: `vi config/system.yaml`
**Key sections**:
- `system.type` - auto/laptop/desktop
- `gpu.method` - virtio/sr-iov/vgpu/full
- `wayland.compositor` - hikari/woodland/dwl
- `microvm.opnsense_firewall.enabled` - true/false
- `vms.*enabled` - enable/disable VMs

### VM Templates
**Directory**: `config/vm-templates/`
- `opnsense-microvm.yaml` - Firewall configuration
- `clean-vm.yaml` - General computing
- `gaming-vm.yaml` - GPU gaming (Looking Glass)
- `research-vm.yaml` - Malware analysis
- `tools-vm.yaml` - Network testing

## Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Overview & quick start |
| `docs/reference/architecture.md` | Detailed design |
| `docs/reference/features.md` | Complete feature list |
| `docs/reference/quick-reference.md` | This file |
| `docs/guides/looking-glass-setup.md` | Looking Glass config |
| `docs/guides/opnsense-firewall-setup.md` | Firewall setup |
| `docs/guides/proxmox-debloat.md` | (pending) |
| `docs/guides/gpu-passthrough.md` | (pending) |
| `docs/guides/security-hardening.md` | (pending) |

## Default Network

### Firewall Setup
- **WAN Bridge**: `vmbr0` (external network)
- **LAN Bridge**: `vmbr-internal` (192.168.100.0/24)
- **Management Bridge**: `vmbr-mgmt` (192.168.99.0/24)

### IP Addresses
```
OPNsense Firewall: 192.168.100.1  (gateway)
Clean VM:          192.168.100.100
Gaming VM:         192.168.100.101
Research VM:       192.168.100.102 (isolated)
Tools VM:          192.168.100.103
```

## GPU Passthrough

### Check Available GPUs
```bash
lspci | grep -i "vga\|3d"
```

### Automatic Method Selection
System selects best method based on hardware:
1. Full passthrough (if multiple GPUs or dedicated)
2. vGPU (NVIDIA cards with vGPU support)
3. SR-IOV (AMD/Intel with SR-IOV support)
4. virtIO (universal fallback)

### Manual GPU Setup
Configure in `config/system.yaml`:
```yaml
gpu:
  method: full  # or sr-iov, vgpu, virtio
  devices: []   # auto-detect
```

## Looking Glass

### Quick Start
```bash
# On Gaming VM (guest side)
# Install client
apt-get install looking-glass-client  # Linux
# or
choco install looking-glass-client    # Windows

# On Proxmox host (client side)
looking-glass-client -h 127.0.0.1 -p 5900

# Full screen
looking-glass-client -f
```

### Shared Memory Size
- Minimum: 32MB
- Recommended: 256MB
- Large (4K): 512MB+

Configure in VM template:
```yaml
looking_glass:
  shared_memory: 256  # MB
```

## OPNsense Firewall

### Access Web UI
- **URL**: https://192.168.100.1
- **SSH**: ssh root@192.168.100.1 -p 22
- **Default Creds**: Set during installation (CHANGE!)

### Create Firewall Rules
1. Go to: Firewall > Rules > LAN
2. Create rule:
   - Source: VM IP (e.g., 192.168.100.100)
   - Destination: Target (e.g., 192.168.100.101)
   - Action: Block (to isolate)

### Monitor Traffic
1. Go to: Monitoring > Traffic
2. View per-interface statistics
3. Check firewall logs: System > Log Files > Firewall

### Enable Logging
1. Firewall > Settings > Advanced
2. Enable Firewall Logging
3. View logs: System > Log Files > Firewall

## Wayland Compositor

### Available Options
- **Hikari** (recommended) - Tiling + stacking
- **Woodland** - Lightweight alternative
- **dwl** - Minimal (requires compilation)

### Configure in system.yaml
```yaml
wayland:
  compositor: auto  # or hikari/woodland/dwl
  stacking: true    # Enable stacking mode
```

### Start Compositor
```bash
# Manual start (if not auto-starting)
hikari

# With configuration file
hikari -c ~/.config/hikari/hikari.conf
```

## Power Management

### Check Current State
```bash
# CPU frequency
cat /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor

# Power profile
powerprofilesctl get

# Battery status (laptop)
upower -i /org/freedesktop/UPower/devices/battery_BAT0
```

### Change Power Profile
```bash
# Set balanced (laptop)
powerprofilesctl set balanced

# Set performance (desktop)
powerprofilesctl set performance

# Or via system config
vi config/system.yaml
power:
  profile: balanced  # or performance/powersave
```

## Troubleshooting

### OPNsense Not Responding
```bash
# Check if VM is running
qm status 50

# Restart VM
qm stop 50
qm start 50

# Check network
ping 192.168.100.1
```

### VMs Cannot Reach Internet
```bash
# Check VM default gateway
ip route  # Should show 192.168.100.1 as default

# Check OPNsense WAN
ssh root@192.168.100.1
netstat -rn  # Check default route
ifconfig em0  # Check WAN interface
```

### GPU Passthrough Not Working
```bash
# Verify IOMMU
grep iommu /proc/cmdline  # Should show iommu=pt intel_iommu=on (Intel)

# Check if device exists
lspci | grep -i "vga\|3d"

# Reboot if just enabled IOMMU
reboot
```

### Looking Glass Slow/Laggy
```bash
# Check shared memory
cat /proc/meminfo | grep Huge

# Increase if needed
echo 512 | sudo tee /proc/sys/vm/nr_hugepages

# Verify ivshmem in VM
ps aux | grep qemu  # Should show -ivshmem
```

### Research VM Isolation Not Working
```bash
# Verify in OPNsense firewall rules
# 1. Block rule: 192.168.100.102 -> any destination
# 2. Log blocked packets enabled
# 3. Check firewall logs for blocked packets

# Check VM routing
ssh user@192.168.100.102
ip route  # Should only show gateway (192.168.100.1)

# Test blocking
ping 8.8.8.8  # Should timeout/fail
```

## Security Checklist

- [ ] Change OPNsense default password
- [ ] Update Proxmox kernel (if IOMMU enabled)
- [ ] Configure firewall rules for VM isolation
- [ ] Enable firewall logging
- [ ] Test Research VM isolation (can't reach internet)
- [ ] Configure automatic backups
- [ ] Review kernel hardening (sysctl)
- [ ] Enable audit logging
- [ ] Set strong root password
- [ ] Disable unnecessary services

## Performance Tuning

### For Gaming
```yaml
# In gaming-vm.yaml
specs:
  cpu:
    cores: 6  # More cores = better
  memory:
    max: 16384  # More RAM = better

gpu:
  method: full  # Full passthrough

security:
  pin_vcpu: true  # Reduce latency
```

### For Laptop Battery
```yaml
# In system.yaml
power:
  profile: balanced
  gpu_dynamic_power: true
  adaptive_brightness: true
  idle_states: deep  # Deeper C-states
```

### For Research VM
```yaml
# No GPU needed
specs:
  cpu:
    cores: 2  # Minimal
  memory:
    max: 4096  # Minimal
```

## Common Issues

| Issue | Solution |
|-------|----------|
| High latency on Gaming VM | Check CPU pinning, increase vCPU cores |
| OPNsense web UI slow | Restart lighttpd: `service lighttpd restart` |
| VM won't start | Check logs: `pvestatd` journal |
| GPU not detected | Enable IOMMU and reboot |
| Research VM reaches internet | Check firewall rules block 192.168.100.102 |
| Battery drain laptop | Lower power profile, disable GPU dynamic power |

## Useful Links

- **Proxmox Docs**: https://pve.proxmox.com/wiki/Main_Page
- **OPNsense Docs**: https://docs.opnsense.org/
- **Looking Glass**: https://looking-glass.io/
- **Wayland Compositors**:
  - Hikari: https://github.com/hikari-no-yume/hikari
  - Woodland: https://github.com/free-deskto/woodland
  - dwl: https://github.com/djpohly/dwl

## File Locations

```
Proxmox Config: /etc/pve/
VM Storage:    /var/lib/vz/
Sysctl Config: /etc/sysctl.d/
Network:       /etc/network/interfaces
Logs:          /var/log/
VM Configs:    /etc/pve/nodes/<hostname>/qemu-server/
```

## Next Steps

1. Read `README.md` for overview
2. Edit `config/system.yaml` for your hardware
3. Run `./scripts/deploy.sh`
4. Run `./scripts/verify.sh`
5. Read `docs/guides/opnsense-firewall-setup.md` for firewall config
6. Read `docs/guides/looking-glass-setup.md` for gaming setup
7. Create test VMs
8. Enable monitoring and logging

---

For detailed information, see the full documentation in `docs/` directory.
