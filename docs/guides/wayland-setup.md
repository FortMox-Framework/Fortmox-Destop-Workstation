# Wayland Compositor Setup - Hikari/Woodland/dwl

> **Alpha/manual guide:** The CLI does not install or configure a Wayland compositor on the Proxmox host. Compatibility, resource, and latency claims here are not FortMox test results.

## Overview

Wayland is a modern display server protocol that provides better security, performance, and resource efficiency than X11. This guide covers setting up three lightweight Wayland compositors with stacking window manager support.

**Goal**: Lightweight, secure, high-performance desktop environment on Proxmox host

## Why Wayland?

| Aspect | Notes |
|--------|-------|
| Security | Isolation properties depend on compositor, protocols, and applications. |
| Performance and memory | Depend on hardware, compositor, drivers, and workload; not measured by FortMox. |
| GPU and input support | Verify compatibility with the selected compositor and Proxmox host environment. |

## Compositor Comparison

| Compositor | Notes |
|------------|-------|
| **Hikari** | Verify current project status, package availability, and Proxmox compatibility upstream. |
| **Woodland** | Verify project status and platform compatibility upstream. |
| **dwl** | Verify project status and platform compatibility upstream. |

## Prerequisites

### System Requirements

```bash
# Check graphics support
glxinfo | grep "OpenGL version"

# Check for Wayland libraries
pkg-config --exists wayland-client && echo "Wayland support: YES"

# Check for DRM support
ls -la /dev/dri/card0
```

### Kernel Requirements

```bash
# Wayland-capable kernel (all modern kernels)
uname -r

# GPU drivers should be updated
lspci | grep -i "vga\|3d\|display"
```

## Method 1: Hikari (Recommended)

### Overview

Hikari is a compositor project. FortMox does not package, install, or validate it on Proxmox.

**Features:**
- Hybrid tiling + stacking window management
- Multi-monitor support
- Floating window support
- Keyboard-driven with mouse support
- Active development
- Minimal resource usage

### Installation

#### 1.1 Install Dependencies

```bash
# Debian/Ubuntu
sudo apt update
sudo apt install -y \
    wayland \
    libwayland-dev \
    wayland-protocols \
    libwlroots-dev \
    libpixman-1-dev \
    libxcb-composite0-dev \
    libxcb-icccm4-dev \
    libxcb-res0-dev \
    libxcb-xinerama0-dev \
    libxcb-xkb-dev \
    libxkbcommon-dev \
    libxkbcommon-x11-dev \
    libxkbregistry-dev \
    libxkbregistry0 \
    libxkbfile-dev \
    libpango-1.0-0 \
    libpango-1.0-dev \
    fonts-inconsolata \
    fonts-liberation \
    meson \
    ninja-build \
    pkg-config
```

#### 1.2 Build from Source (Optional - Pre-built Available)

```bash
# Add Hikari PPA (Ubuntu)
sudo add-apt-repository ppa:geier/hikari

# Or build from source
git clone https://github.com/hikari-mirror/hikari.git
cd hikari
meson setup build -Dwerror=false
ninja -C build
sudo ninja -C build install
```

#### 1.3 Install Pre-built Package

```bash
# Ubuntu/Debian (simpler method)
sudo apt install hikari

# Verify installation
which hikari
hikari --version
```

### Configuration

#### 1.4 Create Configuration File

```bash
mkdir -p ~/.config/hikari
nano ~/.config/hikari/hikarirc
```

Add default configuration:

```bash
# Hikari Configuration File
# ~/.config/hikari/hikarirc

# Workspace setup
workspace 1
workspace 2
workspace 3
workspace 4
workspace 5
workspace 6
workspace 7
workspace 8
workspace 9
workspace 10

# Keyboard layout
keyboard-layout us

# Mouse configuration
mouse-button-modifier mod4  # Super key

# Keybindings
binding mod4+w kill-window
binding mod4+d exec "rofi -show drun"
binding mod4+return exec kitty
binding mod4+e exec thunar
binding mod4+f fullscreen-window
binding mod4+shift+f fullscreen-screen
binding mod4+m maximize-window

# Workspace switching
binding mod4+1 switch-workspace 1
binding mod4+2 switch-workspace 2
binding mod4+3 switch-workspace 3
binding mod4+4 switch-workspace 4
binding mod4+5 switch-workspace 5
binding mod4+6 switch-workspace 6
binding mod4+7 switch-workspace 7
binding mod4+8 switch-workspace 8
binding mod4+9 switch-workspace 9
binding mod4+0 switch-workspace 10

# Window management
binding mod4+h focus-left
binding mod4+j focus-down
binding mod4+k focus-up
binding mod4+l focus-right
binding mod4+shift+h exchange-left
binding mod4+shift+j exchange-down
binding mod4+shift+k exchange-up
binding mod4+shift+l exchange-right

# View/Layout switching
binding mod4+v cycle-view
binding mod4+shift+v cycle-layout
binding mod4+minus decrease-master-count
binding mod4+plus increase-master-count

# System commands
binding mod4+shift+q exit-compositor
binding mod4+shift+r reload-configuration

# Screen dimming/power
binding XF86MonBrightnessUp     exec brightnessctl set +5%
binding XF86MonBrightnessDown   exec brightnessctl set 5%-
binding XF86AudioRaiseVolume    exec pactl set-sink-volume @DEFAULT_SINK@ +5%
binding XF86AudioLowerVolume    exec pactl set-sink-volume @DEFAULT_SINK@ -5%
binding XF86AudioMute           exec pactl set-sink-mute @DEFAULT_SINK@ toggle

# Output configuration (monitor setup)
output eDP-1 scale 1 enable mode 1920x1080@60
```

### Additional Setup

#### 1.5 Install Complementary Tools

```bash
# Terminal
sudo apt install kitty

# Launcher/Menu
sudo apt install rofi

# File manager
sudo apt install thunar

# Status bar (optional)
sudo apt install waybar

# Lock screen
sudo apt install swaylock

# Notification daemon
sudo apt install dunst

# Screenshot
sudo apt install grim slurp wl-clipboard
```

#### 1.6 Create Startup Script

```bash
nano ~/.config/hikari/startup.sh
chmod +x ~/.config/hikari/startup.sh
```

Add:

```bash
#!/bin/bash

# Hikari startup script
# Place in ~/.config/hikari/startup.sh

# Set environment variables
export XDG_CURRENT_DESKTOP=hikari
export QT_QPA_PLATFORM=wayland
export QT_AUTO_SCREEN_SCALE_FACTOR=1
export QT_LOGGING_RULES="*=false"

# Start daemons/services
dunst &
pipewire &
wireplumber &

# Status bar (if using waybar)
# waybar &

# Wallpaper
# swaybg -i /path/to/wallpaper.jpg -m fill &
```

#### 1.7 Enable Session Login

Create `/usr/share/wayland-sessions/hikari.desktop`:

```bash
sudo nano /usr/share/wayland-sessions/hikari.desktop
```

Add:

```ini
[Desktop Entry]
Name=Hikari
Comment=Hybrid tiling Wayland compositor
Exec=/usr/bin/hikari
Type=Application
```

### Usage

```bash
# Start Hikari manually
hikari

# Or select from display manager (GDM, SDDM, etc.)
# Choose "Hikari" from session menu at login

# Keyboard shortcuts
mod4 (Super) + h/j/k/l    : Focus window
mod4 + return            : Open terminal
mod4 + d                 : Application launcher
mod4 + 1-9               : Switch workspace
mod4 + shift + h/j/k/l   : Move window
mod4 + f                 : Toggle fullscreen
mod4 + v                 : Cycle view (tiling ↔ stacking)
```

## Method 2: Woodland (Alternative)

### Installation

```bash
# Debian/Ubuntu
sudo apt install woodland

# Or build from source
git clone https://github.com/fantonangeli/woodland
cd woodland
meson build
ninja -C build
sudo ninja -C build install
```

### Basic Configuration

```bash
mkdir -p ~/.config/woodland
nano ~/.config/woodland/config.toml
```

Example config:

```toml
[general]
mod_key = "Super_L"
terminal = "kitty"
launcher = "rofi -show drun"

[colors]
background = "#1e1e2e"
accent = "#89dceb"

[input]
kb_layout = "us"
mouse_accel = 1.2
```

## Method 3: dwl (Minimal Fallback)

### Installation

```bash
# Build minimal dwm-like Wayland compositor
git clone https://github.com/djpohly/dwl
cd dwl

# Edit config.h if needed
nano config.h

# Build
make

# Install
sudo make install
```

### Configuration

Edit `config.h` before building:

```c
/* dwl configuration */

static const char *tags[] = { "1", "2", "3", "4", "5", "6", "7", "8", "9" };

static const unsigned int nmaster = 1;      /* clients in master area */
static const float mfact = 0.55;            /* master factor [0.05..0.95] */

static const struct keybind binds[] = {
    /* modifier      key             action          argument */
    { MODKEY,        XKB_KEY_d,      spawn,          SHCMD("rofi -show drun") },
    { MODKEY,        XKB_KEY_Return, spawn,          SHCMD("kitty") },
    { MODKEY,        XKB_KEY_q,      quit,           {0} },
    /* ... more keybindings ... */
};
```

## Multi-Monitor Configuration

### Detecting Monitors

```bash
# List connected outputs
wayland ls  # or similar compositor-specific command

# For Hikari
hikari -c  # Check outputs in config preview
```

### Configuring Outputs (Hikari Example)

```bash
# Edit ~/.config/hikari/hikarirc
output HDMI-A-1 enable mode 1920x1080@60 pos 0 0
output DP-2 enable mode 2560x1440@144 pos 1920 0
output eDP-1 disable  # Disable laptop screen when docked
```

### Getting Output Names

```bash
# In Wayland session, check logs
journalctl --user-unit sway -n 20 | grep "output"

# Or use wlr-randr (if available)
sudo apt install wlr-randr
wlr-randr

# Output example:
# HDMI-A-1 (connected)
# 1920x1080 @ 60.00 Hz (preferred, current)
```

## Font and Theme Configuration

### Install Fonts

```bash
# Monospace fonts for terminal
sudo apt install fonts-inconsolata fonts-liberation fonts-dejavu

# Or modern alternatives
sudo apt install fonts-fira-code fonts-source-code-pro

# Copy to local fonts
mkdir -p ~/.local/share/fonts
cp /path/to/font.ttf ~/.local/share/fonts/
fc-cache -fv
```

### GTK Theme (for Wayland compatibility)

```bash
sudo apt install lxappearance
lxappearance  # Launch GUI theme selector

# Or manually configure
mkdir -p ~/.config/gtk-3.0
nano ~/.config/gtk-3.0/settings.ini
```

Add:

```ini
[Settings]
gtk-application-prefer-dark-theme = true
gtk-theme-name = Adwaita-dark
gtk-icon-theme-name = Adwaita
gtk-font-name = Inconsolata 11
```

## Integration with Other Tools

### Screen Locking

```bash
# Install swaylock
sudo apt install swaylock

# Bind to key (in Hikari config)
binding mod4+alt+l exec swaylock -c 000000

# Or with swayidle for automatic locking
sudo apt install swayidle

# Create ~/.config/swayidle/config
exec swayidle -w before-sleep 'swaylock -c 000000'
```

### Screenshot Support

```bash
# Install screenshot tools
sudo apt install grim slurp wl-clipboard

# Bind keys (Hikari example)
binding Print exec grim - | wl-copy
binding mod4+shift+s exec grim -g "$(slurp)" - | wl-copy
```

### Looking Glass Integration

For Gaming VM with Looking Glass + Wayland:

```bash
# Install Looking Glass client
# (See docs/guides/looking-glass-setup.md)

# Configure for Wayland
# Some versions may need: looking-glass-client -m 9  # Use SDL backend
```

## Performance Tuning

### Enable Hardware Acceleration

```bash
# Set environment variables
export LIBVA_DRIVER_NAME=iHD  # Intel
export LIBVA_DRIVER_NAME=radeonsi  # AMD
export LIBVA_DRIVER_NAME=vdpau  # NVIDIA (legacy)

# Or add to ~/.profile
echo 'export LIBVA_DRIVER_NAME=iHD' >> ~/.profile
```

### Monitor Scaling

For High-DPI Displays (4K, high res laptops):

```bash
# In Hikari config
output eDP-1 scale 2 mode 2560x1440@60

# Or via GDM settings (if using GNOME login)
gsettings set org.gnome.desktop.interface scaling-factor 2
```

### Reducing Resource Usage

```bash
# Disable animations
# Monitor scaling
# Disable unnecessary daemons

# Check resource usage
ps aux | grep -i wayland
free -h
```

## Troubleshooting

### Black Screen on Startup

```bash
# Check logs
journalctl --user-unit wayland -n 50

# Try different input method
# Switch to TTY: Ctrl+Alt+F2
# Try: hikari -c  # Run with validation

# Check if GPU drivers are loaded
lsmod | grep -i nouveau
lsmod | grep -i amd
```

### Keyboard Not Working

```bash
# Reload keyboard layout
# Restart compositor

# Check input devices
cat /dev/input/event*

# Use fallback ASCII layout
setxkbmap us
```

### Mouse Issues

```bash
# Check mouse settings in config
# Disable mouse acceleration in settings

# Test mouse:
# Create test window and move it
```

### Application Crashes

```bash
# Run with debug output
MOZ_LOG=all hikari 2>&1 | tee hikari.log

# Check specific application compatibility
# Some X11 apps may need XWayland wrapper
```

### GPU-related Issues

```bash
# Force specific GPU backend
export LIBGL_ALWAYS_INDIRECT=1  # Force software rendering
export LIBVA_DRIVER_NAME=iHD  # Force Intel HD Graphics
export MESA_DEBUG=fps  # Show FPS counter
```

## Performance Monitoring

### Watch Compositor Performance

```bash
# In terminal (within Wayland session)
watch -n 1 'ps aux | grep -i wayland | head -5'

# Memory usage
free -h

# CPU usage
top -p $(pgrep hikari | head -1)
```

### FPS Counter for Gaming

If using Looking Glass:

```bash
# Looking Glass client with FPS counter
looking-glass-client -m 9 -f
# Press F1 to show overlay with FPS
```

## Security Notes

### Wayland Security Benefits

1. **Input Isolation**: Applications can't spy on keyboard input to other apps
2. **Clipboard Isolation**: Better clipboard access control
3. **Screen Recording**: Only authorized apps can record screen
4. **Modern Security Model**: Built on 2010s+ security principles

### Recommended Configuration

```bash
# In ~/.config/hikari/hikarirc
# Restrict window access
# Don't run untrusted apps with full access
# Use workspaces to separate sensitive work
```

## Desktop Environment Summary

| Feature | Hikari | Woodland | dwl |
|---------|--------|----------|-----|
| Configuration | Text file | TOML/Config | C header |
| Ease of Use | Medium | Easy | Hard |
| Features | Excellent | Good | Minimal |
| Performance | Excellent | Excellent | Excellent |
| Support | Active | Active | Maintained |
| **Recommendation** | ✅ Best | Good | Fallback |

## Quick Start

```bash
# 1. Install Hikari (recommended)
sudo apt install hikari kitty rofi

# 2. Create config
mkdir -p ~/.config/hikari
# Copy config from above

# 3. Login to Hikari session (or run: hikari)

# 4. Try shortcuts
# mod4 (Super/Windows key) + return  : Terminal
# mod4 + d                           : Application launcher
# mod4 + h/j/k/l                     : Navigation
```

## Integration with Proxmox

When using this as host desktop environment:

```bash
# Start Proxmox management services
sudo systemctl start pveproxy
sudo systemctl start pvedaemon

# Access web UI from Firefox (Wayland-compatible)
firefox https://localhost:8006
```

---

**CLI status**: Wayland installation and configuration are not implemented.
