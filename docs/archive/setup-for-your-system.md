> **Archived.** Historical snapshot from early development. Numbers, status claims and "your system" details may be out of date. Current docs: [docs/README.md](../README.md).

# FortMoxDesktop Setup Plan - For Your System

## Your Hardware Profile

```
System: Fedora Atomic Linux (Development VM)
CPU: Ryzen 7 8845HS (8 cores / 16 threads) ✓ Excellent
RAM: 16GB (13GB usable) ✓ Sufficient for testing
Storage: 1TB SSD ✓ Perfect for encrypted VMs
GPU: Integrated Radeon (6 cores) ✓ Supports passthrough
Display: Touchscreen ✓ Good for testing UI

Virtualization: AMD-V capable ✓
IOMMU: AMD-Vi capable ✓
```

## Current Situation

**What You Have**:
- ✓ FortMoxDesktop complete project in `/home/devsys/claude/FortMoxDesktop Workstation`
- ✓ Development VM with good specs
- ✓ All scripts and documentation ready
- ✓ Hardware that COULD run Proxmox

**What You DON'T Have**:
- ✗ Proxmox VE installed (currently running Fedora Atomic)
- ✗ Access to bare metal (you're in a VM)
- ✗ Ability to install Proxmox directly

## Two Paths Forward

### Path A: Test & Validate on Fedora (What You Can Do Now)

**Advantages**:
- No system changes needed
- Can test scripts and documentation
- Can learn the system
- Risk-free

**What You CAN Do**:
✓ Review all documentation
✓ Validate YAML configurations
✓ Test hardware detection script
✓ Understand architecture & design
✓ Plan your future Proxmox deployment

**Commands to Run Now**:

```bash
cd "/home/devsys/claude/FortMoxDesktop Workstation"

# 1. Review the project
cat README.md                    # Overview
cat docs/reference/quick-reference.md          # Common commands
cat docs/reference/architecture.md             # How it all works

# 2. Check your hardware
sudo ./scripts/detect-hardware.sh

# 3. Validate configurations
python3 -c "
import yaml
with open('config/system.yaml') as f:
    config = yaml.safe_load(f)
print('✓ Configuration is valid YAML')
"

# 4. Check script syntax
for script in scripts/*.sh; do
    bash -n "$script" && echo "✓ $(basename $script) is valid"
done

# 5. Read the detailed guides
cat docs/guides/security-hardening.md
cat docs/guides/gpu-passthrough.md
cat docs/guides/power-management.md
```

### Path B: Deploy on Bare Metal Proxmox (Future)

When you have access to dedicated hardware:

1. **Get hardware that meets requirements**
   - Ryzen 7 or better (you already have this in your laptop!)
   - 16GB+ RAM ✓ You have this
   - 1TB SSD ✓ You have this
   - IOMMU support ✓ Your CPU has it

2. **Install Proxmox VE**
   - Download ISO (free)
   - Create bootable USB
   - Install on bare metal (replaces existing OS)
   - Takes ~30 minutes

3. **Deploy FortMoxDesktop**
   ```bash
   # Copy project to Proxmox host
   scp -r "FortMoxDesktop Workstation" root@proxmox-host:/root/

   # SSH in and deploy
   ssh root@proxmox-host
   cd "FortMoxDesktop Workstation"
   sudo ./scripts/deploy.sh
   sudo ./scripts/verify.sh
   ```

---

## Immediate Actions (Path A - Testing)

### Step 1: Explore the Project

```bash
cd "/home/devsys/claude/FortMoxDesktop Workstation"

# Show project structure
echo "=== PROJECT FILES ===" && find . -type f \( -name "*.md" -o -name "*.sh" \) | head -20

# Count everything
echo "=== STATISTICS ===" && \
  echo "Documentation files: $(find docs -name '*.md' | wc -l)" && \
  echo "Config files: $(find config -name '*.yaml' | wc -l)" && \
  echo "Scripts: $(find scripts -name '*.sh' | wc -l)" && \
  echo "Total lines: $(wc -l docs/*.md config/*.yaml scripts/*.sh *.md 2>/dev/null | tail -1 | awk '{print $1}')"
```

### Step 2: Understand the Security Model

This is the most important part - FortMoxDesktop uses **8 layers of defense**:

```
Layer 1: IOMMU (Hardware)      ← Your CPU supports this
   ↓
Layer 2: Kernel Hardening     ← Read: docs/guides/security-hardening.md
   ↓
Layer 3: AppArmor MAC         ← Mandatory Access Control
   ↓
Layer 4: Network Isolation    ← Read: docs/guides/network-isolation.md
   ↓
Layer 5: OPNsense Firewall    ← Read: docs/guides/opnsense-firewall-setup.md
   ↓
Layer 6: VLAN Segmentation    ← Per-VM isolation
   ↓
Layer 7: Storage Encryption   ← Read: docs/guides/storage-encryption.md
   ↓
Layer 8: Audit Logging        ← Complete activity tracking
```

### Step 3: Review Architecture

Read these in order:

1. **README.md** (5 min) - What is FortMoxDesktop?
2. **docs/reference/architecture.md** (20 min) - How does it work?
3. **docs/reference/features.md** (15 min) - What features exist?
4. **docs/reference/quick-reference.md** (10 min) - Common commands

### Step 4: Study Security Details

Your interest areas (based on project):

- **Proxmox Hardening**: `docs/guides/proxmox-debloat.md`
- **Network Security**: `docs/guides/network-isolation.md`
- **Storage Security**: `docs/guides/storage-encryption.md`
- **System Hardening**: `docs/guides/security-hardening.md`
- **GPU Gaming**: `docs/guides/gpu-passthrough.md`

### Step 5: Plan Your Future Proxmox Setup

Based on YOUR hardware:

```
Your Specs: Ryzen 7 8845HS, 16GB RAM, 1TB SSD

Recommended VM Setup:
┌─────────────────────────────────────────────────┐
│         Proxmox Host (2 cores reserved)         │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │ OPNsense Firewall MicroVM (vmid: 50)     │  │
│  │ 2 cores | 2GB RAM | <5GB disk            │  │
│  │ → All traffic routes through this        │  │
│  └──────────────────────────────────────────┘  │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │ Clean VM (vmid: 100)                     │  │
│  │ 4 cores | 6GB RAM | 50GB disk            │  │
│  │ → Daily work, web, office                │  │
│  └──────────────────────────────────────────┘  │
│                                                 │
│  ┌──────────────────────────────────────────┐  │
│  │ Research VM (vmid: 102)                  │  │
│  │ 4 cores | 4GB RAM | 100GB disk           │  │
│  │ → Malware analysis (air-gapped)          │  │
│  └──────────────────────────────────────────┘  │
│                                                 │
│  Gaming VM optional (needs GPU passthrough)    │
│                                                 │
└─────────────────────────────────────────────────┘

Total: 12 cores (of 16), 12GB RAM (of 13GB), 150GB+ disk (of 1TB)
Result: ✓ All VMs run simultaneously
        ✓ Host has reserved resources
        ✓ Comfortable headroom for everything
```

---

## Testing Plan: What to Validate

### Test 1: Configuration Correctness

```bash
# Make sure all config files are valid YAML
python3 << 'EOF'
import yaml
import os

config_dir = "config"
for filename in os.listdir(config_dir):
    if filename.endswith(".yaml"):
        filepath = os.path.join(config_dir, filename)
        try:
            with open(filepath) as f:
                yaml.safe_load(f)
            print(f"✓ {filename}")
        except Exception as e:
            print(f"✗ {filename}: {e}")
EOF
```

### Test 2: Script Validity

```bash
# Validate all bash scripts
for script in scripts/*.sh; do
    if bash -n "$script" 2>/dev/null; then
        echo "✓ $(basename $script) - syntax OK"
    else
        echo "✗ $(basename $script) - syntax ERROR"
        bash -n "$script"  # Show error
    fi
done
```

### Test 3: Documentation Completeness

```bash
# Check all referenced docs exist
grep -h "docs/" docs/*.md *.md | grep -oE "docs/[a-z-]+\.md" | sort -u
```

### Test 4: Script Functionality (Read-Only)

```bash
# Run scripts in read-only mode where possible

# Hardware detection
echo "=== HARDWARE DETECTION ==="
sudo ./scripts/detect-hardware.sh 2>&1 | head -50

# Configuration parsing
echo "=== CONFIG VALIDATION ==="
./scripts/config-manager.sh config/system.yaml 2>&1 | head -20
```

---

## Documentation Reading Order

For learning FortMoxDesktop:

**Priority 1 (Must Read)**:
1. `README.md` - 5 min - What is this?
2. `docs/reference/architecture.md` - 20 min - How does it work?
3. `docs/reference/quick-reference.md` - 10 min - Common tasks

**Priority 2 (Should Read)**:
1. `docs/guides/security-hardening.md` - 30 min - How is it secure?
2. `docs/guides/network-isolation.md` - 30 min - VM isolation explained
3. `docs/guides/proxmox-debloat.md` - 20 min - Lightweight setup

**Priority 3 (As Needed)**:
1. `docs/guides/gpu-passthrough.md` - Gaming setup
2. `docs/guides/storage-encryption.md` - Data protection
3. `docs/guides/power-management.md` - Battery/performance
4. `docs/guides/wayland-setup.md` - Desktop environment
5. `docs/guides/opnsense-firewall-setup.md` - Firewall config

**Advanced (Reference)**:
- `docs/reference/features.md` - Feature breakdown
- `docs/archive/project-discussion-log.md` - Why these choices?
- `CHANGELOG.md` - What's implemented
- `docs/README.md` - Topic index

---

## What You Should Do RIGHT NOW

### Immediate (5 minutes)

```bash
cd "/home/devsys/claude/FortMoxDesktop Workstation"

# Read the overview
cat README.md | less

# See what's in the project
ls -la
echo ""
ls docs/
echo ""
ls config/
echo ""
ls scripts/
```

### Short Term (30 minutes)

```bash
# Run validation tests
echo "=== SYNTAX CHECK ==="
for s in scripts/*.sh; do bash -n "$s" && echo "✓ $(basename $s)"; done

# Check configuration
echo ""
echo "=== CONFIG CHECK ==="
python3 -c "import yaml; f=open('config/system.yaml'); yaml.safe_load(f); print('✓ Config is valid YAML')"

# Read key documentation
cat docs/reference/architecture.md | less
```

### Medium Term (2 hours)

```bash
# Deep dive into security
cat docs/guides/security-hardening.md | less

# Understand network isolation
cat docs/guides/network-isolation.md | less

# Learn about firewall
cat docs/guides/opnsense-firewall-setup.md | less
```

### Long Term (Plan for Future)

When you have bare metal Proxmox:

1. Install Proxmox VE on your hardware
2. Copy FortMoxDesktop to Proxmox
3. Run deployment scripts
4. Follow guides for configuration

---

## Your Next Milestone: Testing on Proxmox

### When You Get Proxmox Access

```bash
# On Proxmox host:
cd "FortMoxDesktop Workstation"

# 1. Validate everything still works
sudo ./scripts/verify.sh

# 2. Deploy the system
sudo ./scripts/deploy.sh

# 3. Create VMs
sudo ./scripts/deploy-vms.sh

# 4. Setup firewall
# Follow: docs/guides/opnsense-firewall-setup.md

# 5. Verify everything
sudo ./scripts/verify.sh
```

Time estimate: 2-3 hours total

---

## Resources for Your Learning

### Your Hardware Supports:

✅ **Virtualization** (AMD-V - Ryzen has it)
- Can run VMs with near-native performance
- Full CPU features available to guests

✅ **IOMMU** (AMD-Vi - Ryzen has it)
- GPU passthrough for gaming
- Direct device access with isolation

✅ **AES-NI** (Ryzen has it)
- Encrypted storage with zero overhead
- Wire-speed encryption

✅ **Memory** (13GB usable is good)
- Can run 3-4 VMs simultaneously
- 4GB per VM is comfortable

✅ **Storage** (1TB SSD is excellent)
- Multiple encrypted vaults
- Snapshots and backups
- Room for multiple systems

---

## Recommended Next Step

**Choose one**:

### Option 1: Learn Now (Best for understanding)
```
Read documentation → Validate scripts → Understand design
(Prepares you for actual deployment)
```

### Option 2: Deploy when possible (Best for hands-on)
```
Get bare metal → Install Proxmox → Deploy FortMoxDesktop
(Jump straight to production)
```

### Option 3: Hybrid (Best balanced)
```
Learn now (read docs) → Validate scripts → Get Proxmox
→ Deploy with confidence
```

---

## Summary for Your System

| Aspect | Your Hardware | Requirement | Status |
|--------|---------------|-------------|--------|
| CPU | Ryzen 7 8845HS | VT-x/SVM | ✅ Excellent |
| Cores | 8c/16t | 4+ | ✅ Great |
| RAM | 16GB | 16GB | ✅ Perfect |
| Storage | 1TB SSD | 256GB+ | ✅ Excellent |
| IOMMU | AMD-Vi | Optional | ✅ Supported |
| GPU | Radeon iGPU | Optional | ✅ Works |
| Display | Touchscreen | Optional | ✅ Good |

**Verdict**: Your hardware is PERFECT for FortMoxDesktop! 🚀

---

## Final Checklist

- [ ] Read README.md & docs/reference/architecture.md
- [ ] Understand the 8 security layers
- [ ] Review your VM layout plan
- [ ] Validate scripts run without errors
- [ ] Check configuration is valid YAML
- [ ] Plan next steps (learning vs deploying)
- [ ] Bookmark key documentation

---

**Everything is ready for you to get started!**

What would you like to do first?
1. **Learn** - Read the documentation deep dive
2. **Test** - Run validation scripts
3. **Deploy** - Plan installation on Proxmox
4. **Ask Questions** - About any part of the system

Let me know! 🎯

---

**Last Updated**: 2024-11-17
**For**: Your System (Fedora Atomic, Ryzen 7 8845HS, 16GB RAM)
**Status**: Ready to Use
