# Deploy FortMoxDesktop on Proxmox - Simple Guide

## Overview

After Proxmox is installed, you'll deploy FortMoxDesktop using **simple copy-paste commands**.

**Time needed**: ~1-2 hours total
- Copy project: 5 minutes
- Run deployment: 30 minutes
- Create VMs: 20 minutes
- Verify: 10 minutes

---

## Prerequisites

✓ Proxmox VE is installed and running
✓ You can access Proxmox web UI (https://your-ip:8006)
✓ You know your root password
✓ FortMoxDesktop project is on this current VM (you have it!)

---

## Step 1: Get FortMoxDesktop Project

**You already have the project!** It's at:
```
/home/devsys/claude/FortMoxDesktop Workstation/
```

You need to copy it to your Proxmox host.

### Option A: Copy via SCP (Easiest)

From your current Fedora Linux VM, copy to Proxmox:

```bash
# Get your Proxmox host IP
# (This is the IP of your Fedora Atomic that now runs Proxmox)

PROXMOX_IP="192.168.1.50"  # Change to YOUR Proxmox IP!

# Copy FortMoxDesktop to Proxmox
scp -r "/home/devsys/claude/FortMoxDesktop Workstation" \
       root@$PROXMOX_IP:/root/

# It will ask for password:
# Enter your Proxmox root password
# (the one you set during installation)

# Wait for copy to complete
# You'll see: "✓ 100% done"
```

### Option B: Manual Copy via Web UI

1. Open Proxmox web UI (https://your-ip:8006)
2. Login as root
3. Create project directory:
   ```
   Left sidebar → Nodes → pve → Shell

   mkdir -p /root/FortMoxDesktop
   ```
4. (This is harder, use Option A if possible)

---

## Step 2: SSH into Proxmox Host

From your Fedora Linux VM:

```bash
# SSH into Proxmox (replace IP with yours)
ssh root@192.168.1.50

# It will ask for password
# Enter your Proxmox root password

# You should see prompt:
# root@proxmox:~#
```

**You are now on your Proxmox host!**

---

## Step 3: Verify FortMoxDesktop is Copied

```bash
# List the copied project
ls -la /root/"FortMoxDesktop Workstation"/

# You should see:
# README.md
# docs/reference/architecture.md
# config/
# docs/
# scripts/
# ... and more files

# If you see these, everything copied correctly!
echo "✓ FortMoxDesktop is ready!"
```

---

## Step 4: Run the Deployment Script

```bash
# Change to project directory
cd "/root/FortMoxDesktop Workstation"

# Run the main deployment script
sudo ./scripts/deploy.sh

# The script will:
# 1. Check system requirements
# 2. Apply kernel hardening
# 3. Configure storage
# 4. Set up networking
# 5. Configure security settings

# Wait for completion (this takes 10-30 minutes)
# Watch the output for any [SUCCESS] or [ERROR] messages
```

**What you'll see**:
```
[INFO] FortMoxDesktop Deployment
[INFO] Checking system...
[✓] Proxmox VE detected
[✓] Kernel 6.x found
[INFO] Applying kernel hardening...
[✓] Kernel parameters applied
[INFO] Setting up storage...
[✓] Storage configured
... (more output)
[SUCCESS] Deployment completed!
```

---

## Step 5: Verify Installation

```bash
# Run verification script
sudo ./scripts/verify.sh

# This will check that everything is working:
# - Kernel hardening applied
# - Security modules loaded
# - Storage configured
# - Networking ready
# - All prerequisites met

# You'll see output like:
# [✓] IOMMU enabled
# [✓] AppArmor running
# [✓] Firewall configured
# ... and many more checks

# If you see mostly [✓], everything is good!
```

---

## Step 6: Create Virtual Networks

Before creating VMs, set up network bridges.

Edit `/etc/network/interfaces`:

```bash
# Open the file with nano
nano /etc/network/interfaces

# You should see something like:
# auto lo
# iface lo inet loopback
#
# auto vmbr0
# iface vmbr0 inet dhcp
#     bridge-ports eno1
#     bridge-stp off
#     bridge-fd 0
```

**Add these lines at the end** (before saving):

```
# Internal network bridges for VMs

auto vmbr-internal
iface vmbr-internal inet static
    address 192.168.100.1
    netmask 255.255.255.0
    bridge-ports none
    bridge-stp off
    bridge-fd 0

auto vmbr100
iface vmbr100 inet static
    address 192.168.100.1
    netmask 255.255.255.0
    bridge-ports none
    bridge-stp off
    bridge-fd 0

auto vmbr101
iface vmbr101 inet static
    address 192.168.101.1
    netmask 255.255.255.0
    bridge-ports none
    bridge-stp off
    bridge-fd 0

auto vmbr102
iface vmbr102 inet static
    address 192.168.102.1
    netmask 255.255.255.0
    bridge-ports none
    bridge-stp off
    bridge-fd 0

auto vmbr103
iface vmbr103 inet static
    address 192.168.103.1
    netmask 255.255.255.0
    bridge-ports none
    bridge-stp off
    bridge-fd 0
```

**Save the file**:
- Press CTRL+X
- Press Y (yes)
- Press ENTER

**Apply network changes**:

```bash
# Restart networking
sudo systemctl restart networking

# Verify bridges were created
ip a | grep "vmbr"

# You should see your new bridges listed
```

---

## Step 7: Deploy Virtual Machines

Now create your VMs:

```bash
# From the project directory:
cd "/root/FortMoxDesktop Workstation"

# Create the VMs
sudo ./scripts/deploy-vms.sh

# This will create:
# - OPNsense Firewall (vmid: 50)
# - Clean VM (vmid: 100)
# - Gaming VM (vmid: 101)
# - Research VM (vmid: 102)
# - Tools VM (vmid: 103)

# You'll see output like:
# [INFO] Creating OPNsense Firewall MicroVM (vmid: 50)...
# [✓] OPNsense VM created
# [INFO] Creating Clean VM (vmid: 100)...
# [✓] Clean VM created
# ... etc

# Wait for all VMs to be created
```

---

## Step 8: Verify VMs Were Created

```bash
# List created VMs
qm list

# You should see:
#       VMID NAME                 STATUS     MEM(MB)    BOOTDISK(GB) PID
#         50 opnsense-microvm     stopped      2048            5.00 0
#        100 clean-vm             stopped      8192           50.00 0
#        101 gaming-vm            stopped     16384          100.00 0
#        102 research-vm          stopped      8192          100.00 0
#        103 tools-vm             stopped      4096           30.00 0

# If you see these, VMs created successfully!
echo "✓ All VMs created!"
```

---

## Step 9: Access Proxmox Web UI

Now you can manage VMs from the web interface:

1. **Open Proxmox web UI**:
   ```
   https://your-proxmox-ip:8006
   ```

2. **Login** as root

3. **In left sidebar**, you'll see:
   ```
   Datacenter
   └─ pve (your node)
      ├─ Domains
      └─ Virtual Machines
         ├─ 50 - opnsense-microvm
         ├─ 100 - clean-vm
         ├─ 101 - gaming-vm
         ├─ 102 - research-vm
         └─ 103 - tools-vm
   ```

---

## Step 10: Next Steps - Configure VMs

### OPNsense Firewall Setup

The OPNsense firewall (vmid: 50) needs special setup:

1. **Download OPNsense ISO**:
   ```bash
   # On Proxmox, download OPNsense
   wget https://mirrors.evowise.com/opnsense/releases/24.1/OPNsense-24.1-dvd-amd64.iso \
        -O /var/lib/vz/template/iso/opnsense-24.1.iso
   ```

2. **Boot OPNsense VM**:
   - In web UI, right-click VM 50
   - Select "Start"
   - Watch console output

3. **Complete OPNsense setup**:
   - Follow docs/guides/opnsense-firewall-setup.md for detailed steps

### Other VMs

For Clean, Gaming, Research, and Tools VMs:

1. **You need guest OS ISOs**:
   - For Linux: Download Ubuntu Server ISO
   - For Windows: Have Windows ISO ready

2. **Upload ISOs to Proxmox**:
   - In web UI: Storage → Proxmox (local)
   - Upload your ISO files

3. **Install OS in each VM**:
   - Start VM
   - Boot from ISO
   - Complete OS installation
   - Install Proxmox guest agent (qemu-guest-agent)

---

## Step 11: Verify Everything Works

```bash
# On Proxmox host, run full verification
cd "/root/FortMoxDesktop Workstation"
sudo ./scripts/verify.sh

# Check for any issues
# Most [✓] items should pass
# Some may show [!] warnings (normal)
```

---

## Security Configuration

After basic setup, apply security hardening:

```bash
# 1. Configure kernel hardening
# Already done by deploy.sh
sysctl -a | grep kernel.kptr_restrict
# Should show: kernel.kptr_restrict = 2

# 2. Verify AppArmor is running
systemctl status apparmor
# Should show: active (running)

# 3. Enable audit logging
systemctl status auditd
# Should show: active (running)

# 4. Verify firewall
ufw status
# Should show active (if UFW enabled)
```

---

## What You've Deployed

🎉 **FortMoxDesktop is now installed with**:

✅ **Security (8 layers)**
- Kernel hardening applied
- AppArmor MAC enabled
- Audit logging running
- IOMMU isolation active
- Encryption ready

✅ **5 Virtual Machines**
- OPNsense firewall (central)
- Clean VM (daily work)
- Gaming VM (high performance)
- Research VM (isolated)
- Tools VM (testing)

✅ **Networking**
- Virtual bridges configured
- VLAN-ready
- Firewall rules applied
- Traffic isolation ready

✅ **Storage**
- Ready for encrypted storage
- LUKS encryption supported
- Snapshots enabled

---

## Troubleshooting

### Problem: Deploy script fails

```bash
# Check Proxmox is running:
systemctl status pveproxy

# Check system requirements:
./scripts/verify.sh

# Run with verbose output:
bash -x ./scripts/deploy.sh 2>&1 | tee deploy.log
```

### Problem: VMs won't start

```bash
# Check VM configuration:
qm config 50  # Check VM 50 config

# Check Proxmox logs:
tail -f /var/log/pve/tasks

# Check for errors:
systemctl status pvestatd
```

### Problem: Network bridges not working

```bash
# Check bridges:
ip a | grep vmbr

# Check interface config:
cat /etc/network/interfaces

# Restart networking:
sudo systemctl restart networking
```

### Problem: Can't SSH into Proxmox

```bash
# Check SSH service:
systemctl status ssh

# Check firewall:
ufw status

# Allow SSH:
ufw allow 22/tcp
```

---

## Complete Setup Checklist

After following all steps, verify:

- [ ] Proxmox installed and accessible
- [ ] FortMoxDesktop project copied to /root/
- [ ] deploy.sh ran successfully
- [ ] verify.sh shows mostly [✓]
- [ ] Virtual network bridges created
- [ ] 5 VMs created successfully
- [ ] VMs visible in web UI
- [ ] OPNsense ISO downloaded
- [ ] Guest OS ISOs ready
- [ ] Security features verified

---

## Next Steps (What to Do Now)

### Immediate

1. ✅ Run this guide to deploy
2. ✅ Verify everything works
3. Get guest OS ISOs:
   ```bash
   # Download Ubuntu Server (for Clean/Research/Tools VMs)
   wget https://releases.ubuntu.com/jammy/ubuntu-22.04.3-live-server-amd64.iso

   # Download Windows (optional, for Gaming VM)
   # (Download from Microsoft directly)
   ```

### This Week

1. Install OPNsense in VM 50
2. Install Linux in VM 100 (Clean)
3. Install Linux in VM 102 (Research)
4. Configure OPNsense firewall
5. Test VM isolation

### Next Week

1. Fine-tune firewall rules
2. Test air-gapping on Research VM
3. Setup encrypted storage (optional)
4. Configure power management (optional)
5. Test full security setup

---

## Summary

| Step | What | Time | Status |
|------|------|------|--------|
| 1 | Copy project | 5 min | ✅ |
| 2 | SSH to Proxmox | 1 min | ✅ |
| 3 | Verify copy | 1 min | ✅ |
| 4 | Run deploy.sh | 30 min | ✅ |
| 5 | Run verify.sh | 10 min | ✅ |
| 6 | Create bridges | 5 min | ✅ |
| 7 | Create VMs | 20 min | ✅ |
| 8 | Verify VMs | 2 min | ✅ |
| 9 | Web UI access | 1 min | ✅ |
| 10 | Next steps | varies | ← You are here |

---

## Congratulations! 🚀

**FortMoxDesktop is deployed on your Proxmox host!**

Your system is now:
- ✅ Secure (8 layers of defense)
- ✅ Isolated (compartmentalized VMs)
- ✅ Performance (minimal overhead)
- ✅ Flexible (easily configurable)

**Need help with next steps?** Refer to the detailed guides in `/root/FortMoxDesktop Workstation/docs/`

---

**Last Updated**: 2024-11-17
**Difficulty**: Easy (just follow steps)
**Duration**: ~2 hours total
