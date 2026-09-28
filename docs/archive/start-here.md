> **Archived.** Historical snapshot from early development. Numbers, status claims and "your system" details may be out of date. Current docs: [docs/README.md](../README.md).

# 🚀 FortMoxDesktop Installation - START HERE

## What To Do Right Now

You have everything you need to install FortMoxDesktop on your Fedora Atomic host.

**Choose your path:**

---

## 📖 Path 1: Quick Start (Recommended)

### Read These 3 Guides In Order

1. **docs/guides/getting-started.md** (5 min read)
   - Overview and roadmap
   - What you'll achieve
   - Timeline and checklist

2. **docs/guides/install-proxmox.md** (15 min read)
   - Step-by-step Proxmox installation
   - Copy-paste commands
   - Troubleshooting included

3. **docs/guides/deploy-fortmox.md** (15 min read)
   - Deploy FortMoxDesktop on Proxmox
   - Create 5 VMs
   - Verify everything works

**Total Time**: ~4 hours (install + setup)

---

## 📋 Path 2: Learn First (If Unsure)

### Read Documentation Before Installing

1. **README.md** - What is FortMoxDesktop?
2. **docs/reference/architecture.md** - How does it work?
3. **docs/guides/getting-started.md** - Installation roadmap
4. **docs/guides/install-proxmox.md** - Follow installation guide

---

## ✅ Your Current Status

```
✓ FortMoxDesktop project: COMPLETE
✓ All scripts: READY
Historical status and documentation counts removed; use the current [Getting Started guide](../guides/getting-started.md).
```

---

## 🎯 What Happens After Installation

```
BEFORE:
  Fedora Atomic (Host)
  └─ Fedora Linux (Guest VM - where you are)

AFTER:
  Proxmox VE (Host)
  ├─ OPNsense Firewall (vmid: 50)
  ├─ Clean VM (vmid: 100)
  ├─ Gaming VM (vmid: 101)
  ├─ Research VM (vmid: 102)
  └─ Tools VM (vmid: 103)

All with 8 layers of security!
```

---

## ⚡ Timeline

| Phase | What | Time | When |
|-------|------|------|------|
| **Prepare** | Read guides, collect info | 1 hour | Today |
| **Install** | Proxmox on Atomic | 2 hours | This week |
| **Deploy** | FortMoxDesktop setup | 1 hour | After Proxmox |
| **Configure** | VMs and firewall | 1-2 hours | Next week |
| **Total** | Full system ready | **4-5 hours** | This month |

---

## 📍 Navigation Guide

### For Installation
→ **Read: docs/guides/getting-started.md** (5 min overview)
→ **Read: docs/guides/install-proxmox.md** (15 min detailed steps)
→ **Read: docs/guides/deploy-fortmox.md** (15 min deploy guide)
→ **Start Installation** (2-3 hours)

### For Learning
→ **Read: README.md** (project overview)
→ **Read: docs/reference/architecture.md** (system design)
→ **Read: docs/reference/quick-reference.md** (common tasks)
→ **Read: docs/** (detailed guides)

### For Reference
→ **Check specific docs/**.md** files as needed during use

---

## 🔑 Key Points

### What You're Installing
- **Proxmox VE** - Professional hypervisor on your hardware
- **FortMoxDesktop** - Hardened security configuration
- **5 VMs** - Pre-configured for different uses

### What This Enables
- ✅ Isolated malware analysis (Research VM)
- ✅ Gaming with native GPU (Gaming VM)
- ✅ Daily work in isolation (Clean VM)
- ✅ Network testing in sandbox (Tools VM)
- ✅ Central firewall control (OPNsense)

### Security You Get
- ✅ 8 layers of defense
- ✅ Zero-trust firewall
- ✅ Complete encryption
- ✅ Air-gapping capability
- ✅ Audit logging

---

## ⚠️ Important Before Starting

**Backup Important Data**
- Fedora Atomic will be REPLACED
- Back up anything you need from current system
- Installation will erase the 1TB disk

**Prepare Information**
- Write down your network settings
- Choose and write down root password
- Plan IP address assignment

**Have These Ready**
- USB drive (8GB+)
- Another computer for downloads
- Internet connection
- ~4 hours of uninterrupted time

---

## 🚦 Next Steps (In Order)

### Step 1: This Moment
```
✓ You're reading this
✓ You understand the overview
→ Next: Read docs/guides/getting-started.md
```

### Step 2: Today/Tomorrow
```
→ Read docs/guides/install-proxmox.md completely
→ Collect network information
→ Download Proxmox ISO
→ Create bootable USB
```

### Step 3: Installation Day
```
→ Boot USB
→ Follow installer (20 minutes)
→ Access Proxmox web UI
→ Verify it's working
```

### Step 4: Deployment
```
→ Copy FortMoxDesktop project
→ Run deploy.sh
→ Create 5 VMs
→ Verify success
```

### Step 5: Configuration
```
→ Install OPNsense
→ Install guest OSes
→ Configure firewall
→ Test everything
```

---

## 📖 Complete Guide List

```
START HERE → You are reading this!

INSTALLATION GUIDES:
├─ docs/guides/getting-started.md ............... Overview & roadmap
├─ docs/guides/install-proxmox.md ....... Proxmox installation
└─ docs/guides/deploy-fortmox.md ........ FortMoxDesktop deployment

REFERENCE GUIDES:
├─ README.md ........................ Project overview
├─ docs/reference/architecture.md ................. System design
├─ docs/reference/quick-reference.md .............. Common commands
├─ docs/reference/features.md ..................... Complete features
└─ docs/README.md ........................ Topic index

DETAILED DOCUMENTATION:
├─ docs/guides/opnsense-firewall-setup.md
├─ docs/guides/security-hardening.md
├─ docs/guides/gpu-passthrough.md
├─ docs/guides/network-isolation.md
├─ docs/guides/storage-encryption.md
├─ docs/guides/power-management.md
├─ docs/guides/wayland-setup.md
├─ docs/guides/proxmox-debloat.md
└─ docs/guides/looking-glass-setup.md

CONFIGURATION FILES:
├─ config/system.yaml
└─ config/vm-templates/*.yaml

AUTOMATION SCRIPTS:
├─ scripts/deploy.sh
├─ scripts/deploy-vms.sh
├─ scripts/verify.sh
├─ scripts/config-manager.sh
├─ scripts/gpu-setup.sh
├─ scripts/power-management.sh
└─ scripts/detect-hardware.sh
```

---

## 🎯 Your Success Path

```
Read docs/guides/getting-started.md
           ↓
Read docs/guides/install-proxmox.md
           ↓
Install Proxmox (~2 hours)
           ↓
Read docs/guides/deploy-fortmox.md
           ↓
Deploy FortMoxDesktop (~1 hour)
           ↓
✅ SYSTEM COMPLETE AND WORKING!
```

---

## ✨ What You'll Have

After completing all steps:

```
Professional Security Workstation:

✓ Proxmox hypervisor on bare metal
✓ 5 isolated virtual machines
✓ Central firewall controlling all traffic
✓ 8 layers of security defense
✓ Encrypted storage ready
✓ Complete audit logging
✓ Air-gapped research capability
✓ Gaming VM with GPU support
✓ Network isolation and testing
✓ Professional documentation
```

---

## 🆘 If You Get Stuck

1. **Check the relevant guide**
   - Installation problem? → docs/guides/install-proxmox.md
   - Deployment issue? → docs/guides/deploy-fortmox.md
   - Firewall question? → docs/guides/opnsense-firewall-setup.md

2. **Check troubleshooting section** in that guide

3. **Search documentation** for your keyword

4. **Read docs/reference/architecture.md** to understand the system

---

## Project Statistics

Historical file and line counts are omitted because the repository has since been reorganized. See the current [documentation index](../README.md).

---

## 🎓 What You'll Learn

By following this guide, you'll learn:
- How to install Proxmox VE
- How to configure a hypervisor
- How to set up VMs
- How to isolate workloads
- How to configure firewalls
- Security hardening best practices
- Professional virtualization setup

**Valuable skills!**

---

## ⏰ Time Estimates

| Task | Time | Difficulty |
|------|------|------------|
| Read overview guides | 30 min | Easy |
| Download & prepare USB | 20 min | Easy |
| Install Proxmox | 20 min | Easy |
| Deploy FortMoxDesktop | 30 min | Easy |
| Verify & test | 15 min | Easy |
| **Total** | **~2.5 hours** | **All Easy!** |

---

## ✅ Checklist Before Starting

- [ ] Read this file
- [ ] Read docs/guides/getting-started.md
- [ ] Backed up important data
- [ ] Collected network information
- [ ] Downloaded Proxmox ISO (or ready to)
- [ ] Have USB drive
- [ ] Have time for installation
- [ ] Ready to proceed

---

## 🚀 Ready?

### Next Action: Read docs/guides/getting-started.md

That file will guide you through the complete process.

**Time to invest**: ~4 hours
**Result**: Professional security workstation
**Difficulty**: Easy (just follow steps)

---

## Final Thought

Everything you need is documented. Every step is explained. Every command is provided.

**You've got this!** 🎉

---

**Last Updated**: 2024-11-17
**Status**: Ready for Installation
**Support**: All guides included in this project

Let's go! → Read docs/guides/getting-started.md next
