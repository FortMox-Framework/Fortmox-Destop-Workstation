# FortMoxDesktop Installation & Testing Guide

## Important: System Requirements

FortMoxDesktop is designed for **Proxmox VE** systems. This guide covers:
1. Testing on non-Proxmox systems (simulation/validation)
2. Installing Proxmox VE on compatible hardware
3. Deploying FortMoxDesktop after Proxmox installation

---

## Part 1: Understanding Your Current System

### Current Environment
- **OS**: Fedora Linux
- **Kernel**: 6.17.7 (excellent, modern kernel)
- **Status**: NOT a Proxmox system yet

### What You Need to Deploy FortMoxDesktop

**Hardware Requirements:**
- CPU with virtualization support (Intel VT-x or AMD SVM)
- CPU with IOMMU support (Intel VT-d or AMD-Vi) - for GPU passthrough
- Minimum 16GB RAM (32GB+ recommended)
- NVMe SSD 256GB+ (essential for performance)
- Discrete GPU (optional, for gaming VM)

**Proxmox VE Installation:**
- Dedicated system or VM
- Proxmox VE 7.x or 8.x
- Cannot be installed over existing Linux (requires fresh install)

---

## Part 2: Pre-Installation Testing on Fedora

### Test 1: Verify Hardware Capabilities

```bash
# Check virtualization support
grep -E "vmx|svm" /proc/cpuinfo | head -5

# Check IOMMU support (for GPU passthrough)
grep -i "iommu" /proc/cpuinfo

# List all GPUs
lspci | grep -iE "vga|3d|display"

# Check RAM
free -h

# Check storage
lsblk -d -o NAME,SIZE,TYPE,ROTA
```

If you see output for all these, your hardware supports FortMoxDesktop!

### Test 2: Run Hardware Detection Script

```bash
# Make executable if not already
chmod +x /home/devsys/claude/FortMoxDesktop\ Workstation/scripts/detect-hardware.sh

# Run detection
sudo /home/devsys/claude/FortMoxDesktop\ Workstation/scripts/detect-hardware.sh
```

This will show you:
- CPU model and core count
- GPU capabilities
- Memory and storage
- Recommendations for your hardware

### Test 3: Validate YAML Configuration

```bash
# Check if YAML is valid (requires python)
python3 -c "
import yaml
with open('/home/devsys/claude/FortMoxDesktop Workstation/config/system.yaml') as f:
    config = yaml.safe_load(f)
    print('✓ Configuration is valid YAML')
    print(f'System type: {config.get(\"system\", {}).get(\"type\", \"auto\")}')
"
```

### Test 4: Verify Scripts Are Executable

```bash
cd "/home/devsys/claude/FortMoxDesktop Workstation"

# Check all scripts
for script in scripts/*.sh; do
    if [ -x "$script" ]; then
        echo "✓ $(basename $script) - executable"
    else
        echo "✗ $(basename $script) - NOT executable"
    fi
done
```

### Test 5: Verify Documentation

```bash
# Count documentation lines
wc -l docs/*.md

# Check critical files exist
for file in README.md CHANGELOG.md docs/reference/architecture.md; do
    if [ -f "$file" ]; then
        echo "✓ $file"
    fi
done
```

---

## Part 3: Full Installation Path

### Option A: On Current Fedora System (Development/Testing)

You CAN test components on Fedora:

```bash
# 1. Install Proxmox tools for testing (if available)
sudo dnf search proxmox

# 2. Test individual components
sudo ./scripts/detect-hardware.sh          # Works on any Linux
# sudo ./scripts/security-audit.sh         # Can test security concepts
# etc.
```

**Limitations**: Cannot test actual VMs or networking without real Proxmox

### Option B: Install Proxmox VE (Recommended)

This is the proper way to deploy FortMoxDesktop:

#### Step 1: Prepare System for Proxmox

```bash
# Proxmox requires:
# - Dedicated hardware (not dual-boot)
# - UEFI or BIOS boot
# - Internet connection
# - Minimum 8GB RAM, 32GB disk (though we recommend more)
```

#### Step 2: Download Proxmox VE ISO

```bash
# Download Proxmox VE (choose version 7.x or 8.x)
# From: https://www.proxmox.com/en/downloads/category/proxmox-virtual-environment

# Download methods:
wget https://enterprise.proxmox.com/iso/proxmox-ve_8.1-2_amd64.iso
# OR
curl -O https://enterprise.proxmox.com/iso/proxmox-ve_8.1-2_amd64.iso
```

#### Step 3: Create Bootable Media

```bash
# On Linux, create bootable USB:
# 1. Plug in USB drive (8GB+)
# 2. Identify device name:
lsblk

# 3. Write ISO to USB:
sudo dd if=proxmox-ve_8.1-2_amd64.iso of=/dev/sdX bs=4M status=progress

# Replace /dev/sdX with your USB device (e.g., /dev/sdb)
# WARNING: This will erase the USB drive!

# Verify write:
sudo sync
```

#### Step 4: Boot Proxmox Installer

```bash
# 1. Reboot your system
# 2. Boot from USB (usually F12, F2, or DEL during boot)
# 3. Select "Proxmox VE (Graphical)"
# 4. Follow installer:
#    - Select target disk (will format!)
#    - Set hostname, domain, password
#    - Configure network
#    - Review and install
```

#### Step 5: Post-Installation Setup

After Proxmox is installed:

```bash
# 1. Log into web UI
#    https://your-proxmox-ip:8006
#    User: root
#    Password: (from installer)

# 2. Copy FortMoxDesktop files to Proxmox host
scp -r "FortMoxDesktop Workstation" root@proxmox-ip:/root/

# 3. SSH into Proxmox host
ssh root@proxmox-ip

# 4. Run FortMoxDesktop deployment
cd "/root/FortMoxDesktop Workstation"
sudo ./scripts/deploy.sh
```

---

## Part 4: Testing Strategy (Before Full Deployment)

### Phase 1: Validation (No Changes to System)

```bash
# 1. Check syntax of all scripts
bash -n scripts/*.sh

# 2. Validate configuration
python3 << 'EOF'
import yaml
with open('config/system.yaml') as f:
    yaml.safe_load(f)
print("✓ Configuration valid")
EOF

# 3. Run verification script (read-only)
# Note: Some checks require root, but won't make changes
sudo ./scripts/verify.sh
```

### Phase 2: Dry-Run (Read-Only Checks)

```bash
# Check what deploy.sh would do (without making changes)
# Edit deploy.sh temporarily to add debug mode:
bash -x ./scripts/deploy.sh 2>&1 | head -100
```

### Phase 3: Incremental Deployment

If you proceed with Proxmox installation:

```bash
# 1. Deploy networking first
sudo ./scripts/deploy.sh --step network

# 2. Deploy storage
sudo ./scripts/deploy.sh --step storage

# 3. Deploy security
sudo ./scripts/deploy.sh --step security

# 4. Create VMs
sudo ./scripts/deploy-vms.sh
```

---

## Part 5: What You Can Test NOW (On Fedora)

### Test A: Configuration Parsing

```bash
cd "/home/devsys/claude/FortMoxDesktop Workstation"

# Test config manager script
# This doesn't require Proxmox
./scripts/config-manager.sh config/system.yaml
```

### Test B: Documentation Validation

```bash
# Check all documentation is present
ls -la docs/

# Verify critical guides exist
cat docs/guides/proxmox-debloat.md | wc -l
cat docs/guides/gpu-passthrough.md | wc -l
cat docs/guides/security-hardening.md | wc -l
```

### Test C: Hardware Detection (Works on Any Linux)

```bash
# This gives hardware recommendations
sudo ./scripts/detect-hardware.sh
```

### Test D: Script Syntax Check

```bash
# Validate all bash scripts
for script in scripts/*.sh; do
    bash -n "$script" && echo "✓ $(basename $script)"
done
```

---

## Part 6: Migration Path: Fedora → Proxmox

### If You Want to Keep Fedora

You can run Fedora as a **VM inside Proxmox**:

```
┌─────────────────────────────────┐
│    Proxmox VE (Bare Metal)      │
│                                 │
│  ┌──────────────────────────┐   │
│  │  Fedora VM (your current │   │
│  │  system runs here as VM) │   │
│  └──────────────────────────┘   │
│                                 │
│  ┌──────────────────────────┐   │
│  │  Gaming VM (GPU passth)  │   │
│  └──────────────────────────┘   │
│                                 │
│  ┌──────────────────────────┐   │
│  │  Research VM (isolated)  │   │
│  └──────────────────────────┘   │
│                                 │
└─────────────────────────────────┘
```

Steps:
1. Install Proxmox on bare metal
2. Export current Fedora system (backup)
3. Create Fedora VM in Proxmox
4. Restore data to VM
5. Run FortMoxDesktop deployment

---

## Part 7: Installation Checklist

### Pre-Installation

- [ ] Verified hardware supports virtualization (VT-x/SVM)
- [ ] Verified hardware supports IOMMU (VT-d/AMD-Vi) - optional for GPU
- [ ] Have Proxmox VE ISO downloaded
- [ ] Have bootable USB drive prepared (8GB+)
- [ ] Backed up any data on target system
- [ ] Have another computer to access web UI
- [ ] Have internet connection during installation

### Installation

- [ ] Boot from Proxmox ISO
- [ ] Complete Proxmox installer
- [ ] Set root password
- [ ] Configure network (static IP recommended)
- [ ] Installation completes (~15 min)
- [ ] System reboots into Proxmox

### Post-Installation

- [ ] Access Proxmox web UI (https://ip:8006)
- [ ] Log in with root + password
- [ ] Copy FortMoxDesktop files to /root/
- [ ] SSH into Proxmox host
- [ ] Run: sudo ./scripts/deploy.sh
- [ ] Run: sudo ./scripts/verify.sh
- [ ] Follow docs/guides/opnsense-firewall-setup.md

### Verification

- [ ] All kernel hardening applied
- [ ] OPNsense firewall VM created
- [ ] Clean VM created
- [ ] Gaming VM created (if desired)
- [ ] Research VM created (air-gapped)
- [ ] Tools VM created
- [ ] Firewall rules configured
- [ ] All VMs start without errors

---

## Part 8: Troubleshooting Installation

### Problem: "No virtualization support detected"

```bash
# Check CPU flags
grep -E "vmx|svm" /proc/cpuinfo

# If empty:
# 1. Enable in BIOS/UEFI (VT-x or SVM)
# 2. Reboot
# 3. Check again
```

### Problem: "Can't download Proxmox ISO"

```bash
# Try alternative mirror
wget https://mirrors.kernel.org/proxmox/iso/proxmox-ve_8.1-2_amd64.iso

# Or download manually from:
# https://www.proxmox.com/en/downloads
```

### Problem: "Deploy script fails"

```bash
# Check Proxmox is actually installed
pveversion

# Check you're root
whoami  # Should output: root

# Run with verbose output
bash -x ./scripts/deploy.sh 2>&1 | tee deploy.log
```

### Problem: "VMs don't start"

```bash
# Check Proxmox is running
systemctl status pveproxy

# Check VM config syntax
qm config 100  # Check VM 100

# View error logs
tail -f /var/log/pve/tasks
```

---

## Part 9: Quick Decision Tree

```
Do you have Proxmox VE installed?
├─ YES → Go to Part 8 (Deploy FortMoxDesktop)
└─ NO → Do you want to install Proxmox?
    ├─ YES → Go to Part 3 (Install Proxmox)
    └─ NO → You can only test on Fedora (Part 5)

For testing on Fedora only:
  ✓ Can validate scripts
  ✓ Can read documentation
  ✓ Can run detect-hardware.sh
  ✗ Cannot create VMs
  ✗ Cannot test networking
  ✗ Cannot test firewall
```

---

## Part 10: Next Steps

### Immediate (Today)

1. **Run hardware detection**
   ```bash
   sudo /home/devsys/claude/FortMoxDesktop\ Workstation/scripts/detect-hardware.sh
   ```

2. **Review hardware compatibility**
   - Check if CPU supports VT-x/SVM
   - Check if IOMMU is available
   - Note RAM and storage amounts

3. **Decide: Install Proxmox or test only?**

### If Installing Proxmox

4. **Prepare system**
   - Backup data
   - Create bootable USB
   - Plan network configuration

5. **Install Proxmox VE**
   - Boot from ISO
   - Follow installer (~30 minutes)

6. **Deploy FortMoxDesktop**
   - Copy files to Proxmox
   - Run deploy.sh
   - Run verify.sh

### If Testing on Fedora Only

4. **Run validation tests**
   ```bash
   bash -n scripts/*.sh  # Check syntax
   ```

5. **Read documentation**
   - Understand architecture
   - Plan VM strategy
   - Review security hardening

---

## Summary

**Current Status**: FortMoxDesktop project is complete and ready for deployment on Proxmox

**Your Path**:
1. Verify hardware supports virtualization
2. Install Proxmox VE on dedicated system
3. Deploy FortMoxDesktop using provided scripts
4. Follow documentation guides for specific features

**Estimated Time**:
- Hardware check: 5 minutes
- Proxmox installation: 30 minutes
- FortMoxDesktop deployment: 20 minutes
- Firewall setup: 30 minutes
- **Total: ~2 hours**

Ready to proceed? Let me know your hardware details and I can help you get started! 🚀

---

**Last Updated**: 2024-11-17
**Difficulty**: Intermediate (requires Linux experience)
**Support**: Refer to specific docs/*.md guides for detailed steps
