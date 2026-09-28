> **Archived.** Historical snapshot from early development. Numbers, status claims and "your system" details may be out of date. Current docs: [docs/README.md](../README.md).

# Paranoid Proxmox Workstation - Complete Index

## 📋 Project Overview

A comprehensive, paranoid Proxmox hypervisor setup designed for:
- **Security Research** - Malware analysis, penetration testing
- **Gaming** - Full GPU passthrough with ultra-low latency (Looking Glass)
- **Privacy** - Compartmentalization and encryption
- **Daily Computing** - Isolated, secure workloads
- **All Platforms** - Works on laptops AND desktops equally

**Key Features**: Universal deployment, security hardening, gaming performance, OPNsense firewall MicroVM, Looking Glass support, battery optimization, Wayland compositor, YAML configuration

---

## 📚 Documentation Structure

### 🔵 START HERE
1. **README.md** - Quick overview and getting started
   - Key features overview
   - Quick start commands
   - Project structure
   - Design principles

2. **docs/reference/architecture.md** - Complete system design
   - System architecture diagram
   - Laptop/desktop universal design
   - Compartmentalization strategy
   - GPU passthrough methods
   - Security layers
   - Performance characteristics

3. **docs/reference/quick-reference.md** - Cheat sheet for common tasks
   - Essential commands
   - VM IDs and IP addresses
   - Network setup
   - Troubleshooting quick fixes

### 🔵 FEATURES & PLANNING
4. **docs/reference/features.md** - Complete feature list
   - All requirements met (detailed)
   - Advanced features
   - Configuration options
   - Test scenarios
   - Performance expectations

5. **docs/README.md** - This file
   - Navigation guide
   - File descriptions
   - Reading recommendations

### 🔵 DETAILED SETUP GUIDES

#### Firewall & Network
6. **docs/guides/opnsense-firewall-setup.md** - OPNsense configuration
   - Installation and setup
   - Network bridges
   - Firewall rules
   - VM isolation configuration
   - Traffic monitoring
   - Troubleshooting

#### Display & Graphics
7. **docs/guides/looking-glass-setup.md** - Ultra-low latency gaming
   - Architecture overview
   - Installation on host and guests
   - Configuration options
   - Network access setup
   - Performance tuning
   - Comparison with other display methods

#### Additional Guides (Pending)
- `docs/guides/proxmox-debloat.md` - Removing unnecessary services
- `docs/guides/gpu-passthrough.md` - Detailed GPU setup (all 4 methods)
- `docs/guides/wayland-setup.md` - Wayland compositor configuration
- `docs/guides/power-management.md` - Battery optimization
- `docs/guides/network-isolation.md` - Advanced firewall rules
- `docs/guides/storage-encryption.md` - Encrypted vaults setup
- `docs/guides/security-hardening.md` - Kernel hardening details

---

## ⚙️ Configuration Files

### Main Configuration
```
config/system.yaml
├── system settings (hostname, timezone)
├── hardware (CPU, GPU, IOMMU)
├── wayland (compositor selection)
├── power management (battery, profiles)
├── security (hardening levels)
├── network (isolation modes)
├── storage (encryption)
├── microvm (firewall settings)
└── vm deployment (which VMs to create)
```

**Edit this one file to customize everything!**

### VM Templates
```
config/vm-templates/
├── opnsense-microvm.yaml (firewall - vmid: 50)
├── clean-vm.yaml (general computing - vmid: 100)
├── gaming-vm.yaml (gaming with GPU - vmid: 101)
├── research-vm.yaml (malware analysis - vmid: 102)
└── tools-vm.yaml (network testing - vmid: 103)
```

Each template has:
- CPU/memory specifications
- GPU configuration
- Network setup
- Security hardening
- VM-specific packages
- Firewall rules
- Boot configuration

---

## 🚀 Deployment & Scripts

### Deployment Workflow
```
1. Edit config/system.yaml (customize for your hardware)
   ↓
2. Run ./scripts/config-manager.sh (hardware detection)
   ↓
3. Run ./scripts/deploy.sh (install and configure)
   ↓
4. Run ./scripts/verify.sh (comprehensive testing)
   ↓
5. Access documentation to fine-tune setup
```

### Available Scripts
```
scripts/
├── deploy.sh (main deployment)
│   ├── Check dependencies
│   ├── Load configuration
│   ├── Install packages
│   ├── Configure kernel
│   ├── Enable IOMMU
│   └── Setup storage
│
├── verify.sh (comprehensive checks)
│   ├── Proxmox installation
│   ├── Hardware capabilities
│   ├── GPU detection
│   ├── Security settings
│   ├── Network configuration
│   └── Detailed report
│
└── config-manager.sh (configuration handler)
    ├── YAML parsing
    ├── Hardware detection
    ├── Auto-detect laptop vs desktop
    └── Apply configurations
```

---

## 🏗️ System Architecture

### Compartmentalization Strategy
```
Proxmox Host (Debloated, Hardened, Lightweight)
├── Wayland Compositor (Hikari/Woodland/dwl)
├── Kernel (Hardened with IOMMU/VT-d)
├── OPNsense MicroVM (Central Firewall - vmid: 50)
│   └── Routes ALL VM traffic through this
│
└── Isolated VMs (Separate security zones)
    ├── OPNsense Firewall (192.168.100.1)
    │   ├── Clean VM (192.168.100.100)
    │   ├── Gaming VM (192.168.100.101)
    │   ├── Research VM (192.168.100.102)
    │   └── Tools VM (192.168.100.103)
    └── Direct GPU passthrough for gaming
```

### VM Network
```
Internet
  ↓ (WAN vmbr0)
OPNsense Firewall (192.168.100.1)
  ↓ (LAN vmbr-internal)
Virtual Bridge (192.168.100.0/24)
  ├── Clean VM (100) → Internet only
  ├── Gaming VM (101) → Internet only
  ├── Research VM (102) → Air-gapped (no network)
  └── Tools VM (103) → Isolated network (DNS/HTTP/HTTPS)
```

---

## 🔒 Security Features

### Layers of Isolation
1. **Host Hardening** - Kernel parameters, minimal services
2. **IOMMU Isolation** - Hardware device isolation
3. **OPNsense Firewall** - Network-level isolation
4. **VM Hardening** - Per-VM security settings
5. **Compartmentalization** - Separate VMs for each threat level
6. **Encryption** - Hardware-accelerated storage encryption
7. **Air-gapping** - Research VM no network access
8. **Logging** - Audit all critical operations

### VM Security Zones
- **Clean VM**: Standard computing, normal isolation
- **Gaming VM**: GPU access, internet access, blocked from other VMs
- **Research VM**: Isolated analysis, air-gapped, maximum hardening
- **Tools VM**: Isolated testing network, controlled access

---

## 🎮 Gaming Features

### GPU Support (All Methods)
1. **virtIO** - Universal, lower performance
2. **SR-IOV** - Medium performance
3. **vGPU** - NVIDIA high performance
4. **Full Passthrough** - Native gaming performance (recommended)

### Looking Glass (Ultra-Low Latency)
- <5ms latency (comparable to native display)
- Native gaming performance
- Shared memory architecture
- Optional network access
- Works with Windows and Linux

### Performance Optimization
- CPU pinning to prevent side-channels
- Huge pages for memory performance
- Cache mode tuning
- I/O threading
- TRIM/Discard for SSD

---

## 🔋 Laptop & Desktop Support

### Universal Design
- **Auto-detection** - Detects battery for laptop vs desktop
- **Battery Optimization** (Laptop):
  - Dynamic CPU frequency (P-States)
  - Deep C-states for idle
  - GPU power management
  - Adaptive brightness
  - 90%+ battery life of bare metal

- **Full Performance** (Desktop):
  - Maximum CPU clocks
  - Performance mode always
  - All features enabled
  - 95-100% native performance

### Tested On
- ✅ Laptops (with battery)
- ✅ Desktops (AC powered)
- ✅ Both use same configuration

---

## 📊 VM Specifications

| VM | ID | CPU | RAM | Storage | GPU | Purpose |
|----|----|-----|-----|---------|-----|---------|
| OPNsense | 50 | 2 | 1-2GB | 20GB | None | Firewall |
| Clean | 100 | 4 | 4-8GB | 50GB | virtIO | Daily work |
| Gaming | 101 | 6+ | 16GB+ | 100GB | Full PT | Gaming |
| Research | 102 | 4 | 2-4GB | 50GB | virtIO | Malware |
| Tools | 103 | 4 | 4-8GB | 50GB | virtIO | Testing |

---

## 🔍 Troubleshooting Quick Links

### Common Issues
| Issue | Solution |
|-------|----------|
| IOMMU not working | Read `docs/reference/architecture.md` (Enable IOMMU section) |
| OPNsense unreachable | Check `docs/reference/quick-reference.md` (Troubleshooting) |
| GPU passthrough failing | See `docs/guides/gpu-passthrough.md` (pending) |
| Research VM reaching internet | Check `docs/guides/opnsense-firewall-setup.md` (Firewall Rules) |
| Looking Glass slow | Read `docs/guides/looking-glass-setup.md` (Performance Tuning) |
| Laptop battery drain | Check `docs/guides/power-management.md` (pending) |
| VM won't start | Run `./scripts/verify.sh` for diagnostics |

---

## ✅ Implementation Checklist

### Before Deployment
- [ ] Read `README.md`
- [ ] Review `docs/reference/architecture.md`
- [ ] Read `docs/reference/features.md`
- [ ] Check hardware requirements
- [ ] Edit `config/system.yaml`

### Deployment
- [ ] Run `./scripts/config-manager.sh`
- [ ] Run `./scripts/deploy.sh`
- [ ] Run `./scripts/verify.sh`
- [ ] Review verification report

### Post-Deployment
- [ ] Read `docs/guides/opnsense-firewall-setup.md`
- [ ] Configure OPNsense firewall rules
- [ ] Setup Looking Glass (if gaming)
- [ ] Enable logging and monitoring
- [ ] Test VM isolation
- [ ] Create backups

### Security Hardening
- [ ] Change OPNsense default password
- [ ] Configure firewall rules
- [ ] Enable audit logging
- [ ] Set up monitoring
- [ ] Review kernel parameters
- [ ] Test isolation (verify VMs can't communicate)

---

## 📖 Reading Recommendations

### For Security-Focused Users
1. `docs/reference/architecture.md` - Understand compartmentalization
2. `docs/guides/opnsense-firewall-setup.md` - Configure firewall
3. `docs/guides/security-hardening.md` (pending) - Kernel hardening
4. `docs/reference/features.md` - See all security features

### For Gaming Users
1. `README.md` - Quick overview
2. `docs/guides/looking-glass-setup.md` - Setup gaming display
3. `config/vm-templates/gaming-vm.yaml` - Review settings
4. `docs/reference/quick-reference.md` - Common commands

### For Laptop Users
1. `README.md` - Overview
2. `docs/reference/architecture.md` - Laptop-specific design
3. `docs/guides/power-management.md` (pending) - Battery optimization
4. `docs/reference/quick-reference.md` - Troubleshooting

### For Developers/Advanced
1. `docs/reference/architecture.md` - Full system design
2. All `docs/*.md` files
3. `config/*.yaml` files - Understand configuration format
4. `scripts/*.sh` - Review automation scripts

---

## 🎯 Use Cases

### 1. Malware Analysis
- Use Research VM (vmid: 102)
- Air-gapped from network
- See: `docs/reference/architecture.md` (Compartmentalization Strategy)
- Follow: `docs/guides/opnsense-firewall-setup.md` (Isolation Rules)

### 2. Gaming
- Use Gaming VM (vmid: 101)
- Enable full GPU passthrough
- Setup Looking Glass
- See: `docs/guides/looking-glass-setup.md`

### 3. Daily Computing
- Use Clean VM (vmid: 100)
- Normal internet access
- Daily applications
- Standard security

### 4. Penetration Testing
- Use Tools VM (vmid: 103)
- Isolated testing network
- Can probe other VMs safely
- Monitor traffic via OPNsense

### 5. Development
- Use separate VMs per project
- Encrypted storage per VM
- Git repos isolated
- See: `config/system.yaml` (VM templates)

---

## 📞 Support Resources

### Official Documentation
- Proxmox: https://pve.proxmox.com/wiki/Main_Page
- OPNsense: https://docs.opnsense.org/
- Looking Glass: https://looking-glass.io/

### Key Repositories
- Proxmox Git: https://git.proxmox.com/
- OPNsense GitHub: https://github.com/opnsense
- Looking Glass GitHub: https://github.com/looking-glass/looking-glass

### Community
- Proxmox Forum: https://forum.proxmox.com/
- OPNsense Forum: https://forum.opnsense.org/
- Looking Glass Forum: https://forum.looking-glass.io/

---

## 📝 File Summary

### Core Documentation (Read These)
| File | Lines | Purpose |
|------|-------|---------|
| README.md | 70 | Overview + quick start |
| docs/reference/architecture.md | 300+ | System design |
| docs/reference/features.md | 400+ | Complete feature list |
| docs/reference/quick-reference.md | 300+ | Cheat sheet |
| docs/README.md | This file | Navigation guide |

### Setup Guides
| File | Purpose |
|------|---------|
| docs/guides/opnsense-firewall-setup.md | Firewall + network isolation |
| docs/guides/looking-glass-setup.md | Gaming display setup |
| docs/guides/proxmox-debloat.md (pending) | Remove unnecessary services |
| docs/guides/gpu-passthrough.md (pending) | GPU configuration |
| docs/guides/wayland-setup.md (pending) | Wayland compositor |
| docs/guides/power-management.md (pending) | Battery optimization |
| docs/guides/security-hardening.md (pending) | Kernel hardening |

### Configuration Files
| File | Lines | Purpose |
|------|-------|---------|
| config/system.yaml | 350+ | Main configuration |
| config/vm-templates/opnsense-microvm.yaml | 250+ | Firewall VM |
| config/vm-templates/clean-vm.yaml | 200+ | General computing VM |
| config/vm-templates/gaming-vm.yaml | 250+ | Gaming VM |
| config/vm-templates/research-vm.yaml | 250+ | Malware analysis VM |
| config/vm-templates/tools-vm.yaml | 200+ | Testing tools VM |

### Deployment Scripts
| File | Lines | Purpose |
|------|-------|---------|
| scripts/deploy.sh | 300+ | Main deployment |
| scripts/verify.sh | 400+ | Comprehensive testing |
| scripts/config-manager.sh | 250+ | Configuration handler |

### Total Project
- **15 files** created
- **3000+ lines** of documentation
- **1000+ lines** of configuration
- **1000+ lines** of automation scripts

---

## 🎓 Learning Path

### Beginner (Just starting)
1. Read `README.md` (5 minutes)
2. Read `docs/reference/architecture.md` (20 minutes)
3. Edit `config/system.yaml` (10 minutes)
4. Run deployment scripts (20 minutes)

### Intermediate (Understanding the system)
1. Read `docs/reference/features.md` (30 minutes)
2. Read `docs/guides/opnsense-firewall-setup.md` (30 minutes)
3. Configure firewall rules (30 minutes)
4. Test VM isolation (20 minutes)

### Advanced (Optimizing)
1. Read all `docs/*.md` files (2 hours)
2. Review `config/*.yaml` templates (30 minutes)
3. Customize scripts (1 hour)
4. Performance tuning (1+ hours)

---

## 🚀 Next Steps

1. **START HERE**: Read `README.md` (5 min)
2. **UNDERSTAND**: Read `docs/reference/architecture.md` (20 min)
3. **CONFIGURE**: Edit `config/system.yaml` (10 min)
4. **DEPLOY**: Run `./scripts/deploy.sh` (20 min)
5. **VERIFY**: Run `./scripts/verify.sh` (10 min)
6. **CONFIGURE FIREWALL**: Follow `docs/guides/opnsense-firewall-setup.md` (30 min)
7. **SETUP GAMING**: Follow `docs/guides/looking-glass-setup.md` (20 min if gaming)
8. **TEST**: Create test VMs and verify isolation (30 min)
9. **HARDEN**: Follow additional guides for full setup (2+ hours)
10. **MONITOR**: Enable logging and monitoring (ongoing)

---

## 📌 Important Notes

- **All files are customizable** via `config/system.yaml`
- **Security is paramount** - Review all firewall rules
- **Test isolation before production use** of Research VM
- **Backup your configuration** regularly
- **Keep everything updated** - Proxmox and OPNsense patches
- **Monitor logs** - Critical for detecting issues
- **Read before executing** - Understand what each script does

---

**Version**: 1.0
**Last Updated**: 2024-11-16
**Status**: Alpha (Core features) + Pending Guides

For questions or issues, refer to the relevant documentation file or troubleshooting section.
