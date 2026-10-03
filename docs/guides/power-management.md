# Power Management & Battery Optimization

> **Alpha/manual guide:** The Go CLI currently plans/applies CPU frequency governors only. It does not configure battery optimization, turbo, GPU power, brightness, or idle states; battery-life figures are not measured by this project.

## Overview

This guide describes possible power settings for laptops and desktops. Its battery-life target is unverified and is not guaranteed by FortMox.

**Key Principle**: Form factor intelligent defaults

- **Laptop**: Battery optimization, efficient power profiles
- **Desktop**: Performance-focused, minimal power constraints

## Detecting Form Factor

### Automatic Detection

```bash
# Check for battery
if [ -d /sys/class/power_supply/BAT0 ] || [ -d /sys/class/power_supply/BAT1 ]; then
    echo "LAPTOP detected"
else
    echo "DESKTOP detected"
fi

# Or use system-specific methods
grep -q "POWER_SUPPLY_STATUS" /proc/acpi/battery/BAT0/info && echo "Laptop"
```

## Part 1: CPU Power Management (P-States)

### Check Current Governor

```bash
# View all CPU governors available
cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_available_governors

# View current governor
cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor

# Check CPU frequencies
watch -n 1 "grep MHz /proc/cpuinfo"
```

### Install CPU Frequency Scaling Tools

```bash
# Debian/Ubuntu
sudo apt install cpufreq-utils powertop tlp tlp-rdw

# Or modern alternative (preferred)
sudo apt install auto-cpufreq tuned
```

### Method 1: Using Intel P-States (Intel CPUs)

Edit `/etc/default/cpufreq-utils`:

```bash
sudo nano /etc/default/cpufreq-utils
```

Add:

```
# Intel P-State driver (modern Intel CPUs)
GOVERNOR="powersave"  # For laptop
# GOVERNOR="performance"  # For desktop

# Or use: schedutil (kernel-based scheduling)
GOVERNOR="schedutil"
```

For more control:

```bash
# Check Intel P-State status
grep -r "intel_pstate" /sys/devices/system/cpu/intel_pstate/

# Adjust turbo boost
echo 0 | sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo  # Enable turbo
echo 1 | sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo  # Disable turbo

Turbo controls are driver- and platform-specific. Do not use the example sysctl key `kernel.cpu.turbo_boost`; it is not a standard Linux sysctl. Consult CPU-vendor and distribution documentation for supported controls.

### Method 2: Using TLP (Recommended for Laptops)

Install:

```bash
sudo apt install tlp tlp-rdw

# Start and enable
sudo systemctl enable tlp
sudo systemctl start tlp

# Check status
tlp-stat -s
```

Configure `/etc/tlp.conf`:

```bash
sudo nano /etc/tlp.conf
```

For **Laptop**:

```conf
# Power profiles
TLPPROFILE="powersave"

# CPU frequency scaling
CPU_SCALING_GOVERNOR_ON_AC="ondemand"
CPU_SCALING_GOVERNOR_ON_BAT="powersave"
CPU_SCALING_MIN_FREQ_ON_AC=800000
CPU_SCALING_MAX_FREQ_ON_AC=3500000
CPU_SCALING_MIN_FREQ_ON_BAT=800000
CPU_SCALING_MAX_FREQ_ON_BAT=1800000

# Turbo boost
CPU_BOOST_ON_AC=1
CPU_BOOST_ON_BAT=0

# Number of CPU cores (battery optimization)
CPU_SCALING_NUM_CORES_ON_AC=0  # Use all cores
CPU_SCALING_NUM_CORES_ON_BAT=2  # Use 2 cores on battery

# Disk power management
DISK_POWER_MANAGEMENT_ON_AC="max_performance"
DISK_POWER_MANAGEMENT_ON_BAT="medium_power"
```

For **Desktop**:

```conf
# Desktop - maximize performance
TLPPROFILE="performance"

CPU_SCALING_GOVERNOR_ON_AC="performance"
CPU_BOOST_ON_AC=1

DISK_POWER_MANAGEMENT_ON_AC="max_performance"
```

Apply changes:

```bash
sudo tlp start
tlp-stat  # Verify configuration
```

### Method 3: Using auto-cpufreq

Install:

```bash
sudo apt install auto-cpufreq

# Start and enable
sudo systemctl enable auto-cpufreq
sudo systemctl start auto-cpufreq

# Check performance
auto-cpufreq --monitor
```

Configuration file `/etc/auto-cpufreq.conf`:

```ini
[charger]
governor = performance
turbo = auto

[battery]
governor = powersave
turbo = never
scaling_min_freq = 800000
scaling_max_freq = 1800000
enable_thresholds = true
start_threshold = 20
stop_threshold = 80
```

## Part 2: GPU Power Management

### NVIDIA GPUs

```bash
# Check GPU state
nvidia-smi

# Set power mode (permanent)
sudo nvidia-smi -pm 1  # Enable persistence mode

# Adjust power limit (gaming vs idle)
# Maximum power draw
sudo nvidia-smi -pl 150  # Set to 150W

# Or use nvidia-powerd (if available)
sudo apt install nvidia-powerd
```

For **Laptop with NVIDIA**:

```bash
# Reduce power consumption
sudo nvidia-smi -pl 80  # Reduce to 80W max

# Or disable discrete GPU when on battery
# (Usually done via BIOS or driver settings)
```

### AMD/Intel GPUs

```bash
# AMD AMDGPU driver
# Check current power state
cat /sys/class/drm/card0/device/pp_power_profile_mode

# Set power-saving mode
echo "3" > /sys/class/drm/card0/device/pp_power_profile_mode

# Intel iGPU
# Usually auto-managed by kernel
# Check frequency:
watch -n 1 'cat /sys/kernel/debug/dri/0/i915_cur_freq'
```

### Integrated GPU (Laptop)

```bash
# Most laptops use integrated GPU (Intel, AMD Ryzen)
# These are auto-managed - minimal configuration needed

# Verify GPU power states
grep -r . /sys/class/drm/card0/device/pp_*.json 2>/dev/null
```

## Part 3: Disk Power Management

### Check Current State

```bash
# Check disk power management
hdparm -B /dev/sda
# 254 = maximum performance
# 128 = balanced
# 1 = maximum power saving

# Check idle timeout
hdparm -S /dev/sda
```

### Enable Power Saving

```bash
# Aggressive power saving (15 seconds idle = spindown)
sudo hdparm -B 1 /dev/sda
sudo hdparm -S 18 /dev/sda  # 18 = 90 seconds

# Moderate (balanced)
sudo hdparm -B 128 /dev/sda
sudo hdparm -S 36 /dev/sda  # 36 = 3 minutes

# Performance (minimal power saving)
sudo hdparm -B 254 /dev/sda
sudo hdparm -S 0 /dev/sda   # No spindown
```

### Make Persistent

Create `/etc/udev/rules.d/99-disk-power.rules`:

```bash
sudo nano /etc/udev/rules.d/99-disk-power.rules
```

Add:

```
# Disk power management
ACTION=="add|change", KERNEL=="sd*", SUBSYSTEMS=="scsi", ATTRS{queue/rotational}=="1", RUN+="/sbin/hdparm -B 128 /dev/$kernel"
ACTION=="add|change", KERNEL=="sd*", SUBSYSTEMS=="scsi", ATTRS{queue/rotational}=="1", RUN+="/sbin/hdparm -S 36 /dev/$kernel"
```

Or via `tlp` (easier):

```conf
# In /etc/tlp.conf
DISK_IDLE_SECS_ON_AC=60
DISK_IDLE_SECS_ON_BAT=10
```

## Part 4: Screen/Display Power Management

### Brightness Control

```bash
# Check available brightness
ls /sys/class/backlight/

# Set brightness to 50%
echo 50 | sudo tee /sys/class/backlight/intel_backlight/brightness

# Or use brightness utility
sudo apt install brightnessctl light

# Adjust
brightnessctl set 50%
light -S 50
```

### Make Keyboard Shortcuts Work

```bash
# Bind to Fn+F5, Fn+F6 (brightness keys)
# Usually handled by display manager

# Or add to ~/.config/hikari/hikarirc (if using Hikari)
binding XF86MonBrightnessUp exec brightnessctl set +5%
binding XF86MonBrightnessDown exec brightnessctl set 5%-
```

### Auto-dimming (Optional)

```bash
# Install auto-brightness
sudo apt install light-locker

# Configure auto-lock on idle
xset dpms 600 900 1200  # Standby, suspend, off timings
```

## Part 5: Network Power Management

### WiFi Power Saving

```bash
# Check current WiFi power mode
iw wlan0 get power_save
# 1 = on, 0 = off

# Enable power saving
sudo iw wlan0 set power_save on

# Or via iwconfig (older method)
iwconfig wlan0 power on
```

Make persistent via `/etc/NetworkManager/conf.d/wifi-power-save.conf`:

```ini
[connection]
wifi.powersave = 3
```

### Ethernet (for VM traffic)

Usually not power-intensive, but ensure:

```bash
# Disable unused network interfaces
ip link show

# Down inactive interfaces
sudo ip link set eth1 down
```

## Part 6: Suspend/Hibernate (Laptop)

### Configure Sleep

Edit `/etc/systemd/sleep.conf`:

```bash
sudo nano /etc/systemd/sleep.conf
```

Add:

```ini
# Laptop - enable hibernate after 30 mins on battery
[Sleep]
SuspendState=mem
HibernateMode=platform
HibernateDelaySec=1800  # 30 minutes

# On AC: Don't hibernate
# HibernateDelaySec=-1  # Disable
```

### Test Suspend

```bash
# Test sleep (15 sec wakeup)
sudo rtcwake -m mem -s 15

# Test hibernate
sudo rtcwake -m disk -s 15
```

### Lid Close Behavior

Edit `/etc/systemd/logind.conf`:

```bash
sudo nano /etc/systemd/logind.conf
```

Add:

```ini
# On laptop - suspend on lid close
HandleLidSwitch=suspend
HandleLidSwitchExternalPower=suspend
HandleLidSwitchDocked=suspend

# Don't sleep if external monitor connected
HandleLidSwitchDocked=ignore
```

## Part 7: VM Power Management

### Adjust VM Memory Balloon

When running VMs on laptop:

```bash
# For gaming VM - disable balloon (needs full RAM)
# In /etc/pve/qemu-server/101.conf:
balloon: 0

# For other VMs - enable balloon (allows dynamic memory)
# balloon: 512  # Minimum 512MB always allocated
```

### CPU Limiting

Prevent VMs from consuming all CPU on battery:

```bash
# In /etc/pve/qemu-server/100.conf
# Limit CPU usage
cores: 4
# When on battery, stop heavy workload VMs
cpulimit: 50  # Max 50% CPU usage
```

## Part 8: Monitoring & Tuning

### Power Consumption Monitoring

```bash
# Real-time power usage (Intel)
sudo powertop --csv=powertop.csv

# Or simplified version
watch -n 1 'cat /sys/class/power_supply/BAT0/current_now'

# Battery status
watch -n 1 'upower -e | grep -A5 BAT0'
```

### CPU Frequency Monitoring

```bash
# Watch CPU MHz in real-time
watch -n 1 'grep MHz /proc/cpuinfo | head -8'

# Or use turbostat (Intel)
sudo apt install linux-tools-generic
sudo turbostat
```

### Thermal Monitoring

```bash
# Check temperatures
sensors
# or
watch -n 1 'cat /sys/class/thermal/thermal_zone*/temp'

# Check thermal throttling
grep -r "Throttling" /sys/class/thermal/
```

## Part 9: Profile Switching Script

Create `/usr/local/bin/power-profile.sh`:

```bash
#!/bin/bash

# Power profile switcher for laptops/desktops

if [ "$1" = "battery" ]; then
    echo "Switching to battery profile (power saving)..."

    # CPU power saving
    echo powersave | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor

    # Disable turbo boost
    echo 1 | sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo

    # Reduce GPU power
    sudo nvidia-smi -pl 80  # If NVIDIA

    # Reduce disk I/O
    sudo hdparm -B 128 /dev/sda

    # Dim screen to 40%
    brightnessctl set 40%

    # Enable WiFi power saving
    sudo iw wlan0 set power_save on

    echo "Battery profile activated"

elif [ "$1" = "ac" ] || [ "$1" = "performance" ]; then
    echo "Switching to AC/performance profile..."

    # CPU performance
    echo performance | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor

    # Enable turbo boost
    echo 0 | sudo tee /sys/devices/system/cpu/intel_pstate/no_turbo

    # Max GPU power
    sudo nvidia-smi -pl 300  # If NVIDIA

    # Disable disk power saving
    sudo hdparm -B 254 /dev/sda

    # Restore brightness to 100%
    brightnessctl set 100%

    # Disable WiFi power saving
    sudo iw wlan0 set power_save off

    echo "AC/Performance profile activated"

else
    echo "Usage: $0 [battery|ac|performance]"
    exit 1
fi

# Show current status
tlp-stat -s
```

Make executable:

```bash
sudo chmod +x /usr/local/bin/power-profile.sh

# Use it
sudo power-profile.sh battery
sudo power-profile.sh ac
```

## Part 10: Expected Battery Life

### Baseline Measurements

```bash
# Test on bare metal (no VMs running)
# Monitor battery drain
watch -n 1 'upower -e | grep -E "percentage|time"'

# Record results
upower -e > baseline.txt
```

### With Proxmox (VMs idle)

FortMox has not measured battery life relative to bare metal. Record repeatable measurements on the same hardware and workload before drawing a comparison.

```bash
# Boot into Proxmox with minimal VMs
# Run same battery test
# Should see similar results

# If battery use changes unexpectedly, inspect the platform's power-management services and kernel driver state:
systemctl status tlp  # Is TLP running?
cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor  # Governor is powersave?
grep MHz /proc/cpuinfo  # Are frequencies reducing?
```

## Summary of Settings

### Laptop Configuration (Recommended)

```ini
# CPU: Powersave governor, disable turbo on battery
# GPU: Reduce power on battery
# Disk: Enable aggressive power saving
# Screen: Brightness at 50-70%
# Network: WiFi power saving enabled
# Unmeasured target only; FortMox has not tested battery life against bare metal.
```

### Desktop Configuration (Recommended)

```ini
# CPU: Performance governor, turbo always on
# GPU: No power restrictions
# Disk: Performance mode
# Screen: Full brightness
# Network: No power saving
# Expected behavior is hardware- and workload-dependent; FortMox has not measured it.
```

## Troubleshooting

### Battery Drains Too Fast

```
1. Check governor: should be "powersave" on battery
2. Check turbo: should be disabled on battery
3. Check disk power: should be limited on battery
4. Check running VMs: kill heavy VMs on battery
5. Check GPU: should have limited power on battery
```

### System Too Slow on Battery

```
1. Check CPU frequency: may be throttled too much
2. Check if application needs performance mode
3. Disable aggressive power saving for that session
4. Use: sudo power-profile.sh ac temporarily
```

### VMs Not Booting with Power Saving

```
1. Increase available CPU cores temporarily
2. Disable balloon memory (needs more RAM)
3. Use performance mode for VM boot
4. Return to battery mode after boot
```

## Advanced: Per-VM Power Control

For different VMs to have different power profiles:

```bash
# In /etc/pve/qemu-server/100.conf (Clean VM)
cpulimit: 50     # Only 50% CPU on battery
cpuunits: 1024

# In /etc/pve/qemu-server/101.conf (Gaming VM)
# No limit - needs full performance
```

---

**CLI status**: CPU governor profiles only; other power controls are not implemented.
