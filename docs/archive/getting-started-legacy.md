# Getting Started with FortMoxDesktop - Quick Overview

## 📋 What is FortMoxDesktop?

A **hardened, security-focused Proxmox setup** with:
- 8 layers of security defense
- 5 pre-configured VMs
- Isolated compartments for different workloads
- Gaming-capable GPU support
- Complete encryption and firewall

**Perfect for**: Security research, malware analysis, gaming, privacy-focused computing

---

## 🎯 Your Goal

Transform your Fedora Atomic host from:
```
Fedora Atomic → Proxmox VE → Deploy FortMoxDesktop
```

Into:
```
Professional security workstation with:
- Central OPNsense firewall
- Isolated Clean VM (daily work)
- Isolated Research VM (analysis)
- Isolated Gaming VM (performance)
- Isolated Tools VM (testing)
```

---

## 📚 Guides You Need

### 1️⃣ **docs/guides/install-proxmox.md** (Start Here!)
**What**: Installing Proxmox VE on your Fedora Atomic host
**Time**: ~2 hours
**Steps**: Download → Create USB → Install → Access web UI
**Status**: Your next step!

### 2️⃣ **docs/guides/deploy-fortmox.md** (After Proxmox)
**What**: Deploying FortMoxDesktop on your Proxmox
**Time**: ~1-2 hours
**Steps**: Copy project → Deploy → Create VMs → Verify
**Status**: After completing Guide 1

### 3️⃣ **Reference Guides** (As Needed)
- `docs/guides/opnsense-firewall-setup.md` - Firewall configuration
- `docs/guides/security-hardening.md` - Security details
- `docs/guides/gpu-passthrough.md` - Gaming setup
- And 6 more specialized guides

---

## 🚀 Quick Start Roadmap

### Week 1: Installation

**Day 1-2: Prepare**
```bash
# 1. Read docs/guides/install-proxmox.md
# 2. Collect network info from your Atomic host
# 3. Download Proxmox ISO
# 4. Create bootable USB
```

**Day 3-4: Install**
```bash
# 1. Boot from USB
# 2. Follow installer (20 minutes)
# 3. Reboot into Proxmox
# 4. Access web UI (https://your-ip:8006)
```

**Day 5: Deploy**
```bash
# 1. SSH into Proxmox
# 2. Copy FortMoxDesktop project
# 3. Run deploy.sh
# 4. Run deploy-vms.sh
# 5. Verify installation
```

### Week 2: Configuration

**Day 1-2: Firewall**
```bash
# 1. Download OPNsense ISO
# 2. Boot OPNsense VM
# 3. Complete OPNsense setup
# 4. Configure firewall rules
```

**Day 3-5: VMs**
```bash
# 1. Download Linux ISO
# 2. Install OS in Clean VM
# 3. Install OS in Research VM
# 4. Install OS in Tools VM
# 5. Test isolation
```

---

## 📖 Documentation Structure

```
FortMoxDesktop Workstation/
├── docs/guides/getting-started.md          ← You are here
├── docs/guides/install-proxmox.md   ← Step 1
├── docs/guides/deploy-fortmox.md    ← Step 2
│
├── README.md                   ← Project overview
├── docs/reference/architecture.md             ← Design details
├── docs/reference/quick-reference.md          ← Common commands
│
├── docs/
│   ├── security-hardening.md
│   ├── network-isolation.md
│   ├── storage-encryption.md
│   ├── gpu-passthrough.md
│   ├── wayland-setup.md
│   ├── power-management.md
│   ├── opnsense-firewall-setup.md
│   └── looking-glass-setup.md
│
├── config/
│   ├── system.yaml
│   └── vm-templates/
│       ├── opnsense-microvm.yaml
│       ├── clean-vm.yaml
│       ├── gaming-vm.yaml
│       ├── research-vm.yaml
│       └── tools-vm.yaml
│
└── scripts/
    ├── deploy.sh
    ├── deploy-vms.sh
    ├── verify.sh
    ├── config-manager.sh
    ├── gpu-setup.sh
    ├── power-management.sh
    └── detect-hardware.sh
```

---

## 🔍 How to Use This Project

### For Installation
1. Read: `docs/guides/install-proxmox.md`
2. Then: `docs/guides/deploy-fortmox.md`
3. Reference: Other guides as needed

### For Learning
1. Read: `README.md` (overview)
2. Read: `docs/reference/architecture.md` (how it works)
3. Read: `docs/reference/features.md` (what's included)
4. Read: `docs/reference/quick-reference.md` (common tasks)

### For Reference
- Need to understand security? → `docs/guides/security-hardening.md`
- Need GPU setup? → `docs/guides/gpu-passthrough.md`
- Need network details? → `docs/guides/network-isolation.md`
- Need firewall config? → `docs/guides/opnsense-firewall-setup.md`

---

## ⚡ TL;DR - The Fast Version

If you're in a hurry:

```bash
# 1. Install Proxmox (follow docs/guides/install-proxmox.md)
#    Takes: ~2 hours

# 2. Deploy FortMoxDesktop (follow docs/guides/deploy-fortmox.md)
#    Takes: ~1 hour

# 3. Total time: ~3 hours for working system

# Done!
```

---

## ✅ Success Criteria

After following all guides, you'll have:

**✓ Proxmox VE** installed on your hardware
**✓ FortMoxDesktop** deployed and configured
**✓ 5 VMs** created and ready to use:
  - OPNsense firewall
  - Clean VM (work)
  - Research VM (analysis)
  - Gaming VM (optional)
  - Tools VM (testing)

**✓ Security** features active:
  - Kernel hardening
  - AppArmor MAC
  - Firewall rules
  - Encryption ready
  - Audit logging

**✓ Verification** all systems working properly

---

## 🎯 Next Action Items

### Right Now (5 minutes)
- [ ] Read this file (you're doing it!)
- [ ] Understand the roadmap above

### Today (30 minutes)
- [ ] Read `docs/guides/install-proxmox.md` completely
- [ ] Collect network information from your Atomic host
- [ ] Decide on installation timing

### This Week (2-4 hours)
- [ ] Download Proxmox ISO
- [ ] Create bootable USB
- [ ] Install Proxmox
- [ ] Access Proxmox web UI

### Next Week (1-2 hours)
- [ ] Deploy FortMoxDesktop
- [ ] Create VMs
- [ ] Verify installation

---

## 💡 Pro Tips

### Before You Start
1. **Backup everything** from Fedora Atomic (installation will erase it)
2. **Take screenshots** of current settings
3. **Write down passwords** securely
4. **Have internet access** during installation

### During Installation
1. **Don't rush** - installation takes time
2. **Keep power on** - don't let computer sleep
3. **Write down IP address** - you'll need it
4. **Keep password safe** - you can't recover it!

### After Installation
1. **Verify everything** before adding VMs
2. **Test firewall rules** before production
3. **Backup regularly** - VMs are critical
4. **Monitor resources** - watch CPU/RAM/disk

---

## ❓ Common Questions

**Q: Can I keep Fedora Atomic?**
A: After Proxmox installs, you can create a Fedora VM inside it. Atomic won't exist as the host anymore.

**Q: How long does this take?**
A: ~3-4 hours total for basic setup. Add 1-2 hours for full configuration.

**Q: Do I need to be a Linux expert?**
A: No. Just follow the guides step-by-step. Copy-paste commands work.

**Q: What if something breaks?**
A: You have the project files. Reinstall Proxmox and try again. It's a fresh start.

**Q: Can I test before installing?**
A: Yes, but limited. You can read documentation and validate scripts on current VM.

**Q: Do I need another computer?**
A: Recommended for accessing Proxmox web UI and downloading ISOs. Not absolutely required.

---

## 📞 Getting Help

**If you get stuck:**

1. **Check the relevant guide** for your issue
   - Installation problem? → docs/guides/install-proxmox.md
   - Deployment issue? → docs/guides/deploy-fortmox.md
   - Firewall question? → docs/guides/opnsense-firewall-setup.md

2. **Check troubleshooting section** in that guide

3. **Search documentation** for keywords

4. **Ask detailed questions** if still stuck

---

## 🎓 Learning Resources

This project includes comprehensive documentation on:

- **Proxmox Administration** - System management
- **Virtualization** - VM concepts
- **Security** - Hardening and isolation
- **Networking** - Firewalls and VLANs
- **Storage** - Encryption and backups
- **GPU Passthrough** - Gaming setup
- **Power Management** - Battery and performance


---

## 🚀 Ready to Begin?

### Start Here:

1. **Read** `docs/guides/install-proxmox.md`
   - Takes: 20 minutes reading
   - No action yet, just understanding

2. **Prepare** your system
   - Takes: 30 minutes collecting info
   - Backup data
   - Create USB

3. **Install** Proxmox
   - Takes: ~2 hours
   - Follow installer step-by-step
   - Reboot into Proxmox

4. **Deploy** FortMoxDesktop
   - Takes: ~1 hour
   - Read `docs/guides/deploy-fortmox.md`
   - Run scripts
   - Verify success

5. **Configure** and use
   - Takes: varies
   - Install OSes in VMs
   - Configure firewall
   - Test everything

---

## Summary

| Action | Time | Guide |
|--------|------|-------|
| Understand | 20 min | This file + docs/reference/architecture.md |
| Prepare | 30 min | docs/guides/install-proxmox.md steps 1-5 |
| Install Proxmox | 2 hours | docs/guides/install-proxmox.md steps 6-13 |
| Deploy FortMoxDesktop | 1 hour | docs/guides/deploy-fortmox.md steps 1-11 |
| **Total** | **~4 hours** | **Complete system ready!** |

---

## Next Step

👉 **Read: `docs/guides/install-proxmox.md`**

This is your complete installation guide. Just follow it step-by-step!

Questions? Ask away! 🚀

---

**Last Updated**: 2024-11-17
**Status**: Ready to Deploy
**Difficulty**: Easy (just follow guides)
