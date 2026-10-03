# Install Proxmox VE

> **Important:** FortMox does not install Proxmox and this repository does not certify a Proxmox host install. Use the current [official Proxmox VE download and installation documentation](https://www.proxmox.com/en/downloads). Version-specific download commands have been retired; do not use a pinned installer image from this page.

## Overview

Install Proxmox VE only on hardware you intend to dedicate to the hypervisor. Installing it erases the selected target disk. Fedora Atomic is not a supported Proxmox host or a prerequisite for these steps.

**Time needed**: ~2 hours total
- Download + USB creation: 20 minutes
- Installation: 20 minutes
- Initial setup: 20 minutes

---

## What You'll Need

**Physical items**:
- [ ] USB drive (8GB or larger)
- [ ] Another computer (for downloading and creating USB)
- [ ] Your Fedora Atomic computer

**Information to collect**:
- [ ] Network IP address (write down your current IP)
- [ ] Network gateway (usually 192.168.1.1)
- [ ] Network DNS (usually 8.8.8.8 and 8.8.4.4)
- [ ] Create a ROOT PASSWORD (write it down!)

---

## Step 1: Backup Important Data (CRITICAL!)

**WARNING**: Installation will erase your Fedora Atomic system!

```bash
# If you have important data, back it up NOW
# This means:
# - Copy important files to USB or external drive
# - Take screenshots of important settings
# - Export any configurations
# - Save any work in progress

# Quick backup to USB:
# 1. Plug in USB drive
# 2. Copy files manually using file manager
# 3. Eject safely
```

---

## Step 2: Collect Your Network Information

You'll need this during installation. Get these details:

```bash
# Run these commands on your Fedora Atomic system:

# Current IP address:
ip addr show | grep "inet " | head -2

# Example output:
# inet 192.168.1.100/24

# Default gateway:
ip route | grep default

# Example output:
# default via 192.168.1.1 dev ...

# DNS servers:
cat /etc/resolv.conf | grep nameserver

# Example output:
# nameserver 8.8.8.8
# nameserver 8.8.4.4
```

**Write these down**:
```
My Network Info:
  IP Address: _______________
  Gateway: _______________
  DNS 1: _______________
  DNS 2: _______________
  Root Password (new): _______________
```

---

## Step 3: Download Proxmox ISO

Do this on another computer with internet and a USB drive handy.

### Option A: Download from Proxmox (Recommended)

```bash
# On another Linux computer:

# Create a downloads directory
mkdir -p ~/proxmox-download
cd ~/proxmox-download

# Download the current Proxmox VE ISO from the official downloads page linked above.

# Wait for download (1-2 GB, so 5-15 minutes depending on speed)
# You'll see progress bar
```

### Option B: Download from Browser (Windows/Mac)

1. Go to: https://www.proxmox.com/en/downloads/category/proxmox-virtual-environment
2. Select the current Proxmox VE release.
3. Download the ISO file
4. Save to computer

---

## Step 4: Create Bootable USB Drive

**On a separate computer** (not the Proxmox target host)

### On Linux:

```bash
# 1. Plug in USB drive (8GB+)

# 2. Identify USB device
lsblk

# Look for your USB in the output, like:
# sdb                      8:0    1    14.4G  0 disk
# └─sdb1                   8:1    1    14.4G  0 part /media/user/USB

# In this example, device is /dev/sdb (NOT /dev/sdb1)

# 3. Unmount the USB (if it auto-mounted)
sudo umount /dev/sdb1

# 4. Write ISO to USB (REPLACE sdb with YOUR device!)
sudo dd if=~/proxmox-download/<current-proxmox-ve-iso>.iso of=/dev/sdb bs=4M status=progress

# Wait for completion (you'll see: "copied X bytes")

# 5. Sync to ensure everything is written
sudo sync

# 6. Eject USB
sudo eject /dev/sdb

echo "✓ USB bootable drive created!"
```

### On Windows:

1. Download and install: https://rufus.ie/
2. Open Rufus
3. Select USB drive from dropdown
4. Select Proxmox ISO file
5. Click START
6. Wait for completion

### On Mac:

1. Open Terminal
2. Find USB device:
   ```bash
   diskutil list
   ```
3. Unmount:
   ```bash
   diskutil unmountDisk /dev/disk2  # Replace with your disk
   ```
4. Write ISO:
   ```bash
   sudo dd if=~/Downloads/<current-proxmox-ve-iso>.iso of=/dev/rdisk2 bs=4m
   ```
5. Eject:
   ```bash
   diskutil eject /dev/disk2
   ```

---

## Step 5: Prepare the Dedicated Proxmox Target Host

**On your Fedora Atomic computer:**

```bash
# 1. Save any open work
# 2. Close all applications
# 3. Shut down cleanly (don't force off)
# 4. Keep power plugged in (important!)

sudo shutdown -h now

# Or use the GUI:
# Settings → Power → Shut Down
```

---

## Step 6: Boot from USB and Start Installation

**On your Fedora Atomic computer:**

```bash
# 1. Insert USB drive into computer

# 2. Turn on computer

# 3. Immediately press boot key (usually one of these):
#    - F12 (Dell, Lenovo)
#    - F2 (HP, Asus)
#    - DEL (some boards)
#    - ESC (Apple)
#    - Try F12 first!

# 4. You should see a boot menu with your USB listed

# 5. Select USB drive with arrow keys
#    It might say "UEFI: SanDisk" or similar

# 6. Press ENTER

# 7. You'll see Proxmox boot options
```

---

## Step 7: Select Proxmox Installation

**When you see the Proxmox boot screen:**

```
Current Proxmox VE installer
├─ Proxmox VE (Graphical)      ← SELECT THIS
├─ Proxmox VE (Terminal UI)
└─ Boot without installing
```

**Use arrow keys to select** "Proxmox VE (Graphical)" **and press ENTER**

Wait for graphical installer to load (30-60 seconds)

---

## Step 8: Follow the Installer

### Screen 1: Welcome
```
"Welcome to Proxmox Virtual Environment"

[Accept License] → Click or press TAB then ENTER
```

### Screen 2: Select Target Disk
```
"Please select the disk where Proxmox VE should be installed"

Select your 1TB SSD:
  [ ] QEMU HARDDISK (1000 GB)  ← SELECT THIS
  [ ] USB Drive (8 GB)

WARNING: This will ERASE the disk!

[Next] → Click or press ENTER
```

### Screen 3: Set Location & Time
```
"Please select the location and timezone"

Country: [United States ▼]
Timezone: [America/New_York ▼]  (or your timezone)
Keyboard: [en-US]

[Next] → Click or ENTER
```

### Screen 4: Set Root Password
```
"Please choose a strong root password"

Password: [________________]
Confirm:  [________________]

Email: [admin@example.com]  (any email is fine)

IMPORTANT: Write down your password!
(You'll need it to log in later)

[Next] → Click or ENTER
```

### Screen 5: Configure Network
```
"Management Interface"

This is where you set network settings

Management Interface: [eth0 ▼]  (leave as is)

IP Address (DHCP):  ○ [selected]
IP Address (Static): ○

RECOMMENDED: Keep DHCP selected for now
(Easier for first-time setup)

[Next] → Click or ENTER
```

### Screen 6: Review Summary
```
"Review your settings before installation"

Disk: /dev/sda (1000 GB)
Password: ••••••••••
Network: DHCP on eth0
...

Everything looks good?

[Install] → Click or ENTER

Installation starts!
```

---

## Step 9: Wait for Installation

**The installer will:**

1. Format the disk (1-2 minutes)
2. Extract files (3-5 minutes)
3. Install packages (10-15 minutes)
4. Create boot config (1-2 minutes)

You'll see a progress bar. **Just wait, don't turn off!**

```
Installing Proxmox Virtual Environment
[████████████████░░░░░░░░░░░░] 45%

Do not turn off computer!
```

---

## Step 10: Reboot

When installation completes:

```
"Installation complete!

[Reboot] → Click or ENTER
```

**The system will reboot automatically** (5-30 seconds)

```
Proxmox VE
Proxmox VE (recovery mode)
Boot existing system
```

It will default to "Proxmox VE" which is correct.

Wait for it to boot (30-60 seconds)

---

## Step 11: First Login

After reboot, you'll see a login screen:

```
Proxmox VE GNU/Linux tty1

proxmox login: _
```

**DO NOT LOG IN HERE** (you don't need to)

Instead, use the web interface (next step)

---

## Step 12: Access Proxmox Web UI

**On ANY computer on your network:**

1. **Find your Proxmox IP address**:
   ```bash
   # On Proxmox terminal, login as root (if needed):
   # username: root
   # password: (your password from Step 4)

   # Then type:
   ip addr show | grep "inet"

   # You'll see something like:
   # inet 192.168.1.50/24

   # Your Proxmox IP is: 192.168.1.50
   ```

2. **Open web browser** on another computer:
   ```
   Type in address bar:
   https://192.168.1.50:8006

   (Replace 192.168.1.50 with YOUR IP)
   ```

3. **Accept security warning** (self-signed certificate):
   ```
   Your browser will show a warning about "untrusted certificate"

   Click "Advanced"
   Click "Accept risk"

   This is normal for first login
   ```

4. **Log in to Proxmox**:
   ```
   Username: root
   Password: (your password from Step 4)
   Realm: Linux PAM

   Click [Login]
   ```

5. **You're in!** 🎉
   ```
   You'll see the Proxmox dashboard with:
   - Node information
   - System resources
   - VM creation options
   ```

---

## Step 13: Verify Installation

In the Proxmox web UI:

```
Left sidebar → Click "pve" (your node name)

You should see:
✓ Status: OK (green)
✓ CPU: 8 cores
✓ Memory: 16 GB
✓ Disk: 1 TB
```

---

## Success! 🎉

Proxmox is now installed and running on the dedicated target host.

**Next step**: Deploy FortMoxDesktop (see next guide)

---

## Troubleshooting

### Problem: Can't access web UI

```bash
# Check if Proxmox is running:
# On Proxmox terminal:

systemctl status pveproxy

# Should show:
# pveproxy.service - loaded active running
```

### Problem: Forgot IP address

```bash
# On Proxmox terminal:

ip addr show

# Look for your network interface (probably eth0)
# and inet address
```

### Problem: Proxmox didn't boot

```bash
# USB drive didn't work?
# Try:
1. Remake bootable USB
2. Try different USB port
3. Try F2 instead of F12 for boot menu
```

### Problem: Installation failed

```bash
# Reboot and try again:
1. Press power button to turn off
2. Wait 10 seconds
3. Turn back on
4. Insert USB
5. Try installation again
```

---

## Next Steps

Once Proxmox is installed:

1. **Verify it's working** (Step 13 above)
2. **Copy FortMoxDesktop project** to Proxmox
3. **Deploy FortMoxDesktop** (see next guide)
4. **Create VMs** (automated)
5. **Configure firewall** (follow docs)

---

## Quick Summary

| Step | What | Time |
|------|------|------|
| 1 | Backup data | 10 min |
| 2 | Collect network info | 5 min |
| 3 | Download ISO | 10 min |
| 4 | Create USB | 5 min |
| 5 | Shut down Atomic | 2 min |
| 6-7 | Boot and start installer | 2 min |
| 8-10 | Follow installer screens | 20 min |
| 10 | Installation process | 20 min |
| 11-12 | Reboot and web UI | 5 min |
| 13 | Verify | 2 min |
| **Total** | | **~2 hours** |

---

## Your Success Checklist

- [ ] Backed up important data
- [ ] Downloaded Proxmox ISO
- [ ] Created bootable USB
- [ ] Shut down Fedora Atomic
- [ ] Booted from USB
- [ ] Selected Proxmox installer
- [ ] Selected target disk
- [ ] Set root password (SAVED IT!)
- [ ] Configured network
- [ ] Started installation
- [ ] Installation completed
- [ ] System rebooted
- [ ] Accessed web UI
- [ ] Verified Proxmox running

---

## That's It!

Proxmox is now running on the dedicated target host.

Next: **Deploy FortMoxDesktop** (see docs/guides/deploy-fortmox.md)

Questions? Just ask! 🚀

---

**Difficulty**: Easy (just follow steps)
**Support**: Refer to troubleshooting section if stuck
