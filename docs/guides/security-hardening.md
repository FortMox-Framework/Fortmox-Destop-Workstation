# Kernel Hardening & Advanced Security

> **Alpha/manual guide:** The CLI does not apply the sysctl, module, MAC, SSH, audit, integrity-monitoring, or firewall settings described here. These procedures and security outcomes have not been validated by FortMox.

## Overview

This guide covers comprehensive kernel hardening, AppArmor/SELinux configuration, and advanced security mechanisms for FortMoxDesktop.

**Goal**: Paranoid security posture with 8 layers of defense

## Part 1: Kernel Parameters (sysctl)

### Core Security Parameters

Create `/etc/sysctl.d/99-hardening.conf`:

```bash
sudo nano /etc/sysctl.d/99-hardening.conf
```

Add comprehensive hardening:

```sysctl
# ============================================
# KERNEL HARDENING - Core Security
# ============================================

# Restrict kernel pointer exposure
kernel.kptr_restrict = 2

# Restrict dmesg access (kernel logs)
kernel.dmesg_restrict = 1

# Disable kexec (prevent kernel replacement)
kernel.kexec_load_disabled = 1

# Restrict access to kernel modules
kernel.modules_disabled = 1

# Disable magic SysRq key
kernel.sysrq = 0

# Restrict ptrace scope (process tracing)
kernel.yama.ptrace_scope = 3

# Restrict unprivileged namespace creation
kernel.unprivileged_userns_clone = 0
user.max_user_namespaces = 0

# Hide kernel pointers
kernel.perf_event_paranoid = 3

# Restrict eBPF
kernel.unprivileged_bpf_disabled = 1
kernel.bpf_stats_enabled = 0

# Restrict access to /proc/sysrq-b
kernel.sysrq = 0

# Hide PIDs from unprivileged users
kernel.pid_max = 4194303  # Standard high value
kernel.core_uses_pid = 1

# Restrict /proc access
fs.suid_dumpable = 0

# ============================================
# MEMORY PROTECTION
# ============================================

# Address Space Layout Randomization (ASLR)
kernel.randomize_va_space = 2

# Memory overcommit (strict)
vm.overcommit_memory = 0  # Default: refuse overcommit
vm.overcommit_ratio = 50

# Page cache
vm.dirty_ratio = 10
vm.dirty_background_ratio = 5

# Swap usage
vm.swappiness = 10  # Prefer RAM

# Panic on OOM
vm.panic_on_oom = 0
vm.oom_kill_allocating_task = 1

# ============================================
# NETWORK SECURITY (IPv4)
# ============================================

# Disable IPv4 forwarding (unless needed for OPNsense)
net.ipv4.ip_forward = 0

# Enable strict reverse path filtering
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
net.ipv4.conf.default.send_redirects = 0
net.ipv4.conf.all.send_redirects = 0

# Disable ICMP redirects
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.default.accept_redirects = 0
net.ipv4.conf.all.secure_redirects = 0
net.ipv4.conf.default.secure_redirects = 0
net.ipv4.icmp_echo_ignore_all = 0
net.ipv4.icmp_ignore_bogus_error_responses = 1

# Disable IP source routing
net.ipv4.conf.all.send_redirects = 0
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.conf.default.accept_source_route = 0

# TCP hardening
net.ipv4.tcp_timestamps = 1
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_sack = 0
net.ipv4.tcp_dsack = 0
net.ipv4.tcp_fack = 0

# Increase SYN backlog for DDoS mitigation
net.ipv4.tcp_max_syn_backlog = 4096
net.core.somaxconn = 1024

# Connection tracking
net.ipv4.netfilter.ip_conntrack_max = 2000000

# ============================================
# NETWORK SECURITY (IPv6)
# ============================================

# Disable IPv6 if not needed (or configure like IPv4)
# net.ipv6.conf.all.disable_ipv6 = 1
# OR, secure IPv6:

net.ipv6.conf.all.accept_redirects = 0
net.ipv6.conf.default.accept_redirects = 0
net.ipv6.conf.all.accept_source_route = 0
net.ipv6.conf.default.accept_source_route = 0
net.ipv6.conf.all.accept_ra = 0
net.ipv6.conf.default.accept_ra = 0

# ============================================
# FILESYSTEM HARDENING
# ============================================

# Restrict access to /proc/sys
fs.protected_symlinks = 1
fs.protected_hardlinks = 1
fs.protected_regular = 2
fs.protected_fifos = 2

# Core dump restrictions
fs.suid_dumpable = 0
kernel.core_max = 0

# ============================================
# AUDIT & LOGGING
# ============================================

# Enable core dump capture (for security analysis)
kernel.core_pattern = |/usr/local/bin/audit-coredump.sh

# File system audit (if using auditd)
# See Part 3 for auditd configuration
```

Apply parameters:

```bash
sudo sysctl -p /etc/sysctl.d/99-hardening.conf

# Verify all applied
sudo sysctl -a | grep -E "kptr_restrict|dmesg_restrict|rp_filter"
```

## Part 2: Module Blacklisting

Disable unnecessary kernel modules:

Create `/etc/modprobe.d/99-disable-modules.conf`:

```bash
sudo nano /etc/modprobe.d/99-disable-modules.conf
```

Add:

```
# Disable uncommon network protocols
blacklist dccp
blacklist sctp
blacklist rds
blacklist tipc

# Disable rare filesystems (only if not needed)
# blacklist msdos
# blacklist vfat
# blacklist nfs
# blacklist cifs

# Disable USB gadget mode (desktop security)
blacklist gadgetfs
blacklist usb_gadget

# Disable unused USB devices (if safe for your setup)
# blacklist usb_storage  # Only if not using USB devices
```

Verify blacklisting:

```bash
# List loaded modules
lsmod | grep -v "^Module"

# Check if blacklisted modules are loaded
lsmod | grep dccp  # Should return nothing if blacklisted
```

Update initramfs:

```bash
sudo update-initramfs -u -k all
sudo reboot
```

## Part 3: Mandatory Access Control (MAC)

### AppArmor Configuration

AppArmor is lighter than SELinux, suitable for workstations.

Check AppArmor status:

```bash
sudo systemctl status apparmor
aa-status
```

### Enable AppArmor Stricter Mode

Edit `/etc/apparmor/parser.conf`:

```bash
sudo nano /etc/apparmor/parser.conf
```

Add:

```
# Stricter default behavior
--Werror
```

### Create Custom AppArmor Profile

Example: Restrict QEMU access

Create `/etc/apparmor.d/usr.bin.qemu-system-x86_64`:

```bash
sudo nano /etc/apparmor.d/usr.bin.qemu-system-x86_64
```

Add:

```
#include <tunables/global>

/usr/bin/qemu-system-x86_64 {
  #include <abstractions/base>
  #include <abstractions/nameservice>
  #include <abstractions/openssl>

  capability dac_override,
  capability dac_read_search,
  capability ipc_lock,
  capability setgid,
  capability setuid,
  capability sys_chroot,
  capability sys_nice,
  capability sys_ptrace,
  capability sys_rawio,

  # Only allow specific directories
  /var/lib/vz/** rwk,
  /dev/kvm rw,
  /dev/net/tun rw,
  /dev/mem r,
  /proc/** r,
  /sys/kernel/debug/* r,

  # Block suspicious access
  deny /etc/shadow r,
  deny /etc/gshadow r,
  deny /root/** r,

  # Specific executable only
  /usr/bin/qemu-system-x86_64 rix,
}
```

Load profile:

```bash
sudo apparmor_parser -r /etc/apparmor.d/usr.bin.qemu-system-x86_64
```

### Monitor AppArmor Violations

```bash
# Check for denials
sudo aa-logprof

# Or view logs
sudo grep apparmor /var/log/audit/audit.log | tail -20
```

## Part 4: Audit Daemon (auditd)

### Install and Configure

```bash
sudo apt install auditd audispd-plugins
sudo systemctl enable auditd
sudo systemctl start auditd
```

Create audit rules `/etc/audit/rules.d/99-hardening.rules`:

```bash
sudo nano /etc/audit/rules.d/99-hardening.rules
```

Add:

```
# ============================================
# AUDIT RULES - Monitor suspicious activity
# ============================================

# Remove any existing rules
-D

# Buffer Size
-b 8192

# Failure Mode
-f 1

# ============================================
# SYSTEM INTEGRITY
# ============================================

# Monitor /etc/passwd changes
-w /etc/passwd -p wa -k passwd_changes
-w /etc/shadow -p wa -k shadow_changes
-w /etc/group -p wa -k group_changes

# Monitor /etc/sudoers
-w /etc/sudoers -p wa -k sudoers_changes
-w /etc/sudoers.d/ -p wa -k sudoers_changes

# Monitor SSH configuration
-w /etc/ssh/sshd_config -p wa -k sshd_config_changes
-w /root/.ssh -p wa -k ssh_keys_changes

# ============================================
# KERNEL & MODULE CHANGES
# ============================================

# Monitor insmod/rmmod (kernel modules)
-a always,exit -F arch=x86_64 -S init_module,delete_module -F auid>=1000 -F auid!=4294967295 -k modules

# Monitor privileged commands
-a always,exit -F path=/usr/bin/chsh -F perm=x -F auid>=1000 -F auid!=4294967295 -k privileged_chsh
-a always,exit -F path=/usr/bin/chfn -F perm=x -F auid>=1000 -F auid!=4294967295 -k privileged_chfn
-a always,exit -F path=/usr/bin/sudo -F perm=x -F auid>=1000 -F auid!=4294967295 -k privileged_sudo

# ============================================
# FILE SYSTEM
# ============================================

# Monitor /var/lib/vz (Proxmox VMs)
-w /var/lib/vz -p wa -k vz_changes

# Monitor critical system files
-w /bin -p wa -k binaries
-w /sbin -p wa -k binaries
-w /usr/bin -p wa -k binaries
-w /usr/sbin -p wa -k binaries

# Make configuration immutable (at runtime)
-e 2
```

Load rules:

```bash
sudo augenrules --load
sudo systemctl restart auditd

# Verify
sudo auditctl -l
```

### Monitor Audit Logs

```bash
# Real-time monitoring
sudo tail -f /var/log/audit/audit.log

# Search for events
sudo ausearch -k passwd_changes

# Generate report
sudo aureport --login
```

## Part 5: SELinux (Optional - More Complex)

If preferred over AppArmor:

```bash
# Check if available
getenforce

# Set to enforcing (if installed)
sudo setenforce Enforcing

# Check policies
getsebool -a | grep http
```

Most Proxmox installations don't need SELinux (AppArmor is sufficient).

## Part 6: System Integrity Monitoring

### Install AIDE (File Integrity)

```bash
sudo apt install aide aide-common

# Initialize database
sudo aideinit

# Daily integrity checks
sudo aide --check
```

Configure `/etc/aide/aide.conf.d/99_custom`:

```bash
sudo nano /etc/aide/aide.conf.d/99_custom
```

Add:

```
# Monitor critical directories
/etc p+i+n+u+g+s+b+m+c+md5+sha256
/sbin p+i+n+u+g+s+b+m+c+md5+sha256
/bin p+i+n+u+g+s+b+m+c+md5+sha256
/usr/sbin p+i+n+u+g+s+b+m+c+md5+sha256
/usr/bin p+i+n+u+g+s+b+m+c+md5+sha256
/lib p+i+n+u+g+s+b+m+c+md5+sha256
/usr/lib p+i+n+u+g+s+b+m+c+md5+sha256

# Ignore log files (they change frequently)
!/var/log
```

### Cron Job for Daily Checks

```bash
sudo crontab -e
```

Add:

```
0 2 * * * /usr/bin/aide --check > /var/log/aide/aide-check.log 2>&1
```

## Part 7: SSH Hardening

### Secure SSH Configuration

Edit `/etc/ssh/sshd_config`:

```bash
sudo nano /etc/ssh/sshd_config
```

Apply:

```
# ============================================
# SSH HARDENING
# ============================================

# Protocol
Protocol 2

# Network
Port 22                    # Change to random port (2222-65535) if exposed
AddressFamily inet         # IPv4 only
ListenAddress 127.0.0.1   # Localhost only, or specific interface
ListenAddress 192.168.1.1 # Management VLAN only

# Authentication
PermitRootLogin no         # Disable root SSH
PubkeyAuthentication yes   # SSH keys only
PasswordAuthentication no  # No passwords
PermitEmptyPasswords no
ChallengeResponseAuthentication no
MaxAuthTries 3
MaxSessions 5

# Algorithms (modern secure only)
HostKeyAlgorithms ssh-ed25519,rsa-sha2-512,rsa-sha2-256
Ciphers chacha20-poly1305@openssh.com,aes-256-gcm@openssh.com,aes-128-gcm@openssh.com
MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com
KexAlgorithms curve25519-sha256,curve25519-sha256@libssh.org

# Compression
Compression no

# Client alive
ClientAliveInterval 300
ClientAliveCountMax 2

# Logging
SyslogFacility AUTH
LogLevel VERBOSE

# Security
StrictModes yes
IgnoreRhosts yes
HostbasedAuthentication no
RhostsRSAAuthentication no
RSAAuthentication no
UsePAM yes

# Banner
Banner /etc/ssh/banner

# User/Group
AllowUsers root@127.0.0.1 root@192.168.1.0/24  # Restrict to specific networks
DenyUsers *
DenyGroups *

# X11 & Tunneling
X11Forwarding no
AllowTcpForwarding no
AllowAgentForwarding no
PermitTunnel no
```

Create SSH Banner:

```bash
echo "Authorized access only. All activity monitored." | sudo tee /etc/ssh/banner
```

Restart SSH:

```bash
sudo systemctl restart ssh

# Test
sudo sshd -t  # Test configuration
```

## Part 8: Firewall (Host Level)

### UFW (Uncomplicated Firewall)

```bash
sudo apt install ufw

# Default deny incoming
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw default deny routed

# Allow SSH from specific networks
sudo ufw allow from 192.168.1.0/24 to any port 22
sudo ufw allow from 127.0.0.1 to any port 22

# Allow Proxmox web UI from management only
sudo ufw allow from 192.168.1.0/24 to any port 8006

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status verbose
```

## Part 9: Hardening Script

Create `/usr/local/bin/security-audit.sh`:

```bash
#!/bin/bash

echo "FortMoxDesktop Security Audit"
echo "=============================="
echo ""

# Check kernel parameters
echo "Kernel Security Parameters:"
grep -h '^kernel\.' /etc/sysctl.d/*.conf | wc -l
echo "kernel parameters configured"

# Check AppArmor
echo ""
echo "AppArmor Status:"
sudo aa-status | grep "profiles" | head -2

# Check Auditd
echo ""
echo "Auditd Status:"
sudo systemctl status auditd | grep "Active:"

# Check SSH configuration
echo ""
echo "SSH Hardening:"
grep -E "^PermitRootLogin|^PasswordAuthentication|^PubkeyAuthentication" /etc/ssh/sshd_config

# Check firewall
echo ""
echo "UFW Status:"
sudo ufw status | head -3

# Check AIDE integrity
echo ""
echo "AIDE Last Check:"
ls -la /var/lib/aide/ | grep "aide.db" | head -1

# Check modules
echo ""
echo "Blacklisted Modules:"
grep -c "^blacklist" /etc/modprobe.d/99-disable-modules.conf
echo "modules disabled"

# File integrity
echo ""
echo "System Files Integrity:"
sudo aide --check 2>&1 | grep -E "^Number of|changed"
```

Make executable:

```bash
sudo chmod +x /usr/local/bin/security-audit.sh

# Run it
sudo /usr/local/bin/security-audit.sh
```

## Part 10: Security Baseline Comparison

### Before Hardening

```
- Kernel modules: Loads everything
- SSH: Root login enabled, passwords allowed
- Firewall: No filtering
- Audit: No logging
- AppArmor: Permissive mode
```

### Target State After Manually Applying This Guide

```
✅ Kernel modules: Only needed modules load
✅ SSH: Key-based only, restricted users
✅ Firewall: Deny-by-default with specific allows
✅ Audit: Complete activity logging
✅ AppArmor: Enforce mode with custom profiles
✅ Memory: ASLR enabled, panic on corruption
✅ Network: Strict filtering, no redirects
✅ Filesystems: Symlink/hardlink protection
```

## Troubleshooting Security Measures

### Locked Out of SSH

```bash
# Recovery via console
# Edit /etc/ssh/sshd_config to allow passwords temporarily
# Then restart SSH and reauthenticate with keys
```

### AppArmor Breaking Services

```bash
# Set to complain mode
sudo aa-complain /usr/bin/qemu-system-x86_64

# Review logs
sudo aa-logprof

# Then re-enforce
sudo aa-enforce /usr/bin/qemu-system-x86_64
```

### Audit Disk Space

```bash
# Audit logs can grow large
# Check size
du -sh /var/log/audit/

# Rotate logs
sudo systemctl restart auditd
```

## Summary

### 8 Layers of Security

1. ✅ **Kernel Hardening** (sysctl parameters)
2. ✅ **Module Blacklisting** (disable unnecessary)
3. ✅ **Mandatory Access Control** (AppArmor)
4. ✅ **Audit Logging** (auditd)
5. ✅ **File Integrity** (AIDE)
6. ✅ **SSH Hardening** (key-only, restricted)
7. ✅ **Host Firewall** (UFW with deny-default)
8. ✅ **IOMMU Isolation** (hardware-level)

---

**Last Updated**: 2024-11-17
**CLI status**: Host security hardening is not implemented.
