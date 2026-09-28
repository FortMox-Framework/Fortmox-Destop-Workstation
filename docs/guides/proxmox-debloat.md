# Proxmox VE Debloat & Lightweight Hardening

> **Alpha/manual guide:** FortMox does not remove Proxmox services or packages. Review every manual change against Proxmox requirements; the performance and compatibility claims below are not verified by this project.

## Overview

Proxmox VE comes with many features enabled by default that may not be necessary for a focused security workstation. This guide covers removing unnecessary services, reducing memory footprint, and streamlining the installation while maintaining core functionality.

**Goal**: Achieve <3% CPU overhead from hypervisor while maintaining all critical features.

## Philosophy

- **Remove**: Services not needed for your use case
- **Keep**: All security and VM management features
- **Optimize**: Kernel parameters for performance and security
- **Monitor**: Ensure nothing critical is broken

## Pre-Debloat Checklist

Before removing services, document your current state:

```bash
# Check running services
systemctl list-units --type=service --all | grep running

# Check installed packages
dpkg -l | grep -E "proxmox|pve" | wc -l

# Check memory usage
free -h

# Check installed Python packages (often bloated)
pip list | wc -l
```

## Step 1: Identify Unnecessary Services

### List All Running Services

```bash
systemctl list-units --type=service --state=running
```

### Common Unnecessary Services to Disable

These can safely be disabled on a focused workstation (verify they're not needed first):

```bash
# Mail services (unless you need local mail)
sudo systemctl disable postfix
sudo systemctl stop postfix

# Enterprise reporting (not needed for personal setup)
sudo systemctl disable pvestatd  # If you don't need stats collection

# Optional: Disable unused hypervisor features
# (Only if you're 100% sure you won't use them)
# sudo systemctl disable pveproxy  # Only if you don't use web UI
# sudo systemctl disable pvedaemon  # Only if you don't use daemon features
```

**WARNING**: Do NOT disable these - they are essential:
- pvestatd (statistics)
- pvedaemon (cluster/guest management)
- pveproxy (web UI)
- qemu (VM execution)
- lxc (container support, if using)

## Step 2: Remove Unnecessary Packages

### List Installed Packages

```bash
dpkg -l | grep "proxmox\|pve" | less
```

### Packages Safe to Remove

```bash
# Backup/HA features (if not using)
# sudo apt remove pve-ha-manager  # Cluster high availability
# sudo apt remove pve-ha-resources

# Subscription utilities (optional)
# sudo apt remove pve-subscription-manager

# Unnecessary dependencies (CAREFUL - check first)
# sudo apt autoremove --dry-run  # See what would be removed first
# sudo apt autoremove

# Optional: Remove spell checkers, unnecessary languages, etc.
sudo apt remove aspell aspell-en
```

### Keep These Packages

Essential for VM management:
- `qemu-server` - KVM/QEMU support
- `lxc` - LXC container support (if using)
- `proxmox-ve` - Meta package
- `proxmox-kernel*` - Optimized kernel
- `libvirt*` - Virtualization library
- `cgroup-lite` - Resource limiting

## Step 3: Optimize Kernel Parameters

### Security-Focused Hardening

Edit `/etc/sysctl.d/99-hardening.conf`:

```bash
sudo nano /etc/sysctl.d/99-hardening.conf
```

Add/modify these parameters:

```sysctl
# Disable kernel module loading (security)
kernel.modules_disabled = 1

# Restrict kernel pointer exposure
kernel.kptr_restrict = 2

# Restrict dmesg access
kernel.dmesg_restrict = 1

# Restrict ptrace access
kernel.yama.ptrace_scope = 2

# Restrict access to /proc/sysrq-b
kernel.sysrq = 0

# Hide PID from non-root users
kernel.pid_max = 65535
fs.suid_dumpable = 0

# Increase PID wrap-around limit (prevent PID reuse attacks)
kernel.pid_max = 4194303

# Increase inotify limits (for VMs with many file watchers)
fs.inotify.max_user_watches = 524288

# Increase connection limits
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048

# Performance: Increase file descriptors
fs.file-max = 2097152

# Performance: Increase port range
net.ipv4.ip_local_port_range = 10000 65535

# Security: Disable ICMP redirects
net.ipv4.conf.all.send_redirects = 0
net.ipv4.conf.default.send_redirects = 0
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.default.accept_redirects = 0

# Security: Prevent IP spoofing
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1

# Performance: Enable TCP Fast Open
net.ipv4.tcp_fastopen = 3
```

Apply changes:

```bash
sudo sysctl -p /etc/sysctl.d/99-hardening.conf
sudo sysctl -a | grep "kernel\." | grep -v "^#"
```

### VM-Specific Kernel Parameters

For VMs (Gaming especially), in `/etc/sysctl.d/99-vm-tuning.conf`:

```sysctl
# Increase memory overcommit slightly (for VM density)
# WARNING: Only if you know your usage
vm.overcommit_memory = 1  # or 0 (default) for strict

# Reduce swappiness (prefer RAM over swap)
vm.swappiness = 10

# Increase dirty page writeback
vm.dirty_ratio = 10
vm.dirty_background_ratio = 5

# Increase max file handles
fs.file-max = 4000000

# Performance: Disable CPU frequency scaling if you want max performance
# (only if on desktop, not laptop)
# This is controlled by BIOS/UEFI mostly
```

## Step 4: Disable Unnecessary Kernel Modules

### List Loaded Modules

```bash
lsmod | head -20
```

### Disable Unnecessary Modules

Create `/etc/modprobe.d/99-blacklist-bloat.conf`:

```bash
sudo nano /etc/modprobe.d/99-blacklist-bloat.conf
```

Add modules to blacklist (only if you don't need them):

```
# Disable unused network modules
# blacklist ppp
# blacklist pppox
# blacklist pppoe

# Disable unused filesystem support
# blacklist cifs       # SMB support
# blacklist nfs        # NFS support
# blacklist ntfs
# blacklist vfat

# Disable unused drivers
# blacklist floppy
# blacklist parport
```

**BE CAREFUL**: Only blacklist if you're sure they're not needed!

## Step 5: Optimize Memory Usage

### Check Memory Usage

```bash
# Current memory usage
free -h

# Per-process memory
ps aux --sort=-%mem | head -20

# Proxmox-specific memory
systemctl status pvedaemon
systemctl status pvestatd
```

### Reduce Background Services Memory

```bash
# Check which daemons are using memory
pgrep -f pve | while read pid; do
    echo "PID: $pid $(ps -p $pid -o comm=) $(ps -p $pid -o rss= | awk '{print $1/1024 " MB"}')"
done
```

### Disable rsync daemon if not needed

```bash
sudo systemctl disable rsync
sudo systemctl stop rsync
```

## Step 6: Optimize Disk I/O

### Check Current I/O Scheduler

```bash
# For each storage device
cat /sys/block/sda/queue/scheduler
```

### Set Optimal I/O Scheduler

Edit `/etc/default/grub` and add to `GRUB_CMDLINE_LINUX_DEFAULT`:

```
# For SSDs: noop or none (lowest overhead)
elevator=none

# For HDDs: cfq or bfq (better for mixed workloads)
elevator=bfq
```

Then update GRUB:

```bash
sudo update-grub
sudo reboot
```

### Persistent I/O Scheduler Configuration

Create `/etc/udev/rules.d/60-scheduler.rules`:

```bash
sudo nano /etc/udev/rules.d/60-scheduler.rules
```

Add:

```
# Set scheduler for NVMe
ACTION=="add|change", KERNEL=="nvme[0-9]*", ATTR{queue/scheduler}="none"

# Set scheduler for SATA/SSD
ACTION=="add|change", KERNEL=="sd*", ATTR{queue/scheduler}="bfq"

# Set scheduler for VM storage (virtio)
ACTION=="add|change", KERNEL=="vd*", ATTR{queue/scheduler}="none"
```

Reload:

```bash
sudo udevadm control --reload-rules
sudo udevadm trigger
```

## Step 7: Clean Unnecessary Files

### Remove Unnecessary Packages

```bash
# Check what would be removed
sudo apt autoremove --dry-run

# Actually remove (careful!)
sudo apt autoremove -y

# Remove package cache
sudo apt clean
sudo apt autoclean

# Check for orphaned config files
sudo dpkg --list |grep "^rc" | awk '{print $2}' | xargs sudo apt remove --purge -y
```

### Reduce Log Files

```bash
# Check journal size
journalctl --disk-usage

# Limit journal size to 1GB
sudo bash -c 'echo "SystemMaxUse=1G" >> /etc/systemd/journald.conf'

# Apply changes
sudo systemctl restart systemd-journald

# Or, clean old logs
sudo journalctl --vacuum=30d  # Keep only last 30 days
```

## Step 8: Optimize Network Stack

### Enable IP Forwarding (for OPNsense and VMs)

```bash
echo "net.ipv4.ip_forward = 1" | sudo tee -a /etc/sysctl.d/99-forwarding.conf
sudo sysctl -p /etc/sysctl.d/99-forwarding.conf
```

### Optimize Bridge Settings

Create `/etc/sysctl.d/99-bridge.conf`:

```sysctl
# Bridge netfilter settings (for firewall)
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1

# Disable bridge STP (Spanning Tree Protocol) if not needed
# net.bridge.bridge-stp = 0
```

## Step 9: IOMMU Configuration (for GPU Passthrough)

### Enable IOMMU in BIOS

- Intel: Enable "VT-d" or "IOMMUv1"
- AMD: Enable "IOMMU" or "AMD-Vi"

### Enable in Kernel

Edit `/etc/default/grub`:

```bash
sudo nano /etc/default/grub
```

Find `GRUB_CMDLINE_LINUX_DEFAULT` and add:

```
# For Intel
intel_iommu=on iommu=pt

# For AMD
amd_iommu=on iommu=pt
```

For UEFI systems, edit `/etc/grub.d/40_custom` instead.

Update GRUB:

```bash
sudo update-grub
sudo reboot
```

### Verify IOMMU is Enabled

```bash
# Check IOMMU groups
for d in /sys/kernel/iommu_groups/*/devices/*; do
    n=${d#*/iommu_groups/}; n=${n%/devices/*}
    printf 'IOMMU Group %s ' "$n"
    lspci -nns "${d##*/}"
done
```

## Step 10: Performance Monitoring

### Create Monitoring Script

Create `/usr/local/bin/proxmox-stats.sh`:

```bash
#!/bin/bash

echo "=== Proxmox Performance Stats ==="
echo ""
echo "CPU Usage:"
top -bn1 | grep "Cpu(s)" | awk '{print "User: "$2" System: "$4" Idle: "$8}'

echo ""
echo "Memory Usage:"
free -h | grep "Mem"

echo ""
echo "Load Average:"
uptime | awk -F'load average:' '{print $2}'

echo ""
echo "Disk I/O:"
iostat -x 1 2 | tail -n 3

echo ""
echo "Network Traffic:"
ifstat -i vmbr0,vmbr-internal 1 1
```

Make it executable:

```bash
sudo chmod +x /usr/local/bin/proxmox-stats.sh
```

Run it:

```bash
/usr/local/bin/proxmox-stats.sh
```

## Step 11: Verification

### Check Everything Still Works

```bash
# Check cluster status
pvecm status

# Check services
systemctl status pvestatd
systemctl status pvedaemon

# List VMs
qm list

# Test VM creation
qm create 999 --name test-vm --memory 1024

# Verify networking
ip a
ip r
```

### Performance Baseline

```bash
# CPU overhead test (before/after)
cat /proc/cpuinfo | grep "cpu MHz"

# Memory footprint
ps aux | awk '{sum += $6} END {print "Total: " sum/1024 " MB"}'

# Check running processes
systemctl list-units --type=service --state=running | wc -l
```

## Post-Debloat Checklist

- [ ] All services started correctly
- [ ] VMs can still be created
- [ ] Networking works (VMs have internet)
- [ ] OPNsense firewall still functions
- [ ] GPU passthrough still works (if used)
- [ ] CPU overhead is <3%
- [ ] Memory footprint is smaller
- [ ] System boots normally

## Monitoring After Debloat

### Watch for Issues

```bash
# Check systemd failures
systemctl --failed

# Check system logs
journalctl -xe | tail -50

# Check Proxmox-specific logs
tail -f /var/log/pve/tasks
```

## Reverting Changes

If something breaks:

```bash
# Re-enable a service
sudo systemctl enable postfix

# Re-install a package
sudo apt install proxmox-backup-client

# Restore kernel parameters from backup
sudo cp /etc/sysctl.d/99-hardening.conf /etc/sysctl.d/99-hardening.conf.bak
# Edit file to revert changes
sudo sysctl -p /etc/sysctl.d/99-hardening.conf
```

## Security Note

Debloating reduces attack surface by removing unused code. However:

1. **Keep security updates**: Always apply Proxmox and kernel updates
2. **Monitor performance**: Watch for issues in logs
3. **Test thoroughly**: Verify nothing critical was removed
4. **Document changes**: Keep notes of what you disabled
5. **Maintain backups**: Always have a recovery plan

## Summary

A debloated Proxmox installation achieves:
- ✅ <3% CPU overhead
- ✅ Reduced memory footprint
- ✅ Smaller attack surface
- ✅ Full VM functionality maintained
- ✅ All security features intact

The key is removing what you don't need while keeping everything critical for security research and VM management.

---

**Last Updated**: 2024-11-17
**CLI status**: Proxmox service/package debloating is not implemented.
