# OPNsense MicroVM Firewall Setup Guide

> **Alpha/manual guide:** The CLI can generate a basic VM creation command, but does not attach WAN/LAN interfaces, configure bridges, route traffic through OPNsense, or install/configure OPNsense. Do not assume VM traffic is firewalled by following the VM plan alone.

## Overview

OPNsense is a lightweight, open-source firewall/router that runs as a MicroVM on your Proxmox host. All VM traffic flows through OPNsense, enabling fine-grained traffic control and isolation between VMs beyond traditional VLANs.

### Key Benefits

- **Complete VM Isolation**: Block direct communication between VMs
- **Centralized Traffic Control**: All traffic flows through single point
- **Deep Monitoring**: Track every packet sent/received by each VM
- **Flexible Compartmentalization**: Create custom rules per VM
- **Lightweight**: MicroVM uses minimal resources (1-2GB RAM, 2 vCPU)
- **Advanced Features**: IDS/IPS, logging, traffic shaping, NAT

## Architecture

```
┌──────────────────────────────────────────────────┐
│           Proxmox Host Hypervisor                 │
│                                                  │
│  ┌────────────────────────────────────────────┐ │
│  │  OPNsense MicroVM Firewall (vmid: 50)      │ │
│  │                                            │ │
│  │  WAN (vmbr0) ←→ Routing ←→ LAN (vmbr-internal) │ │
│  │                    ↓                       │ │
│  │              Firewall Rules                │ │
│  │              Traffic Logging               │ │
│  │              IDS/IPS                       │ │
│  │              Traffic Shaping               │ │
│  │              NAT/Routing                   │ │
│  └────────────────────────────────────────────┘ │
│                         │                       │
│      ┌──────────────────┼──────────────────┐    │
│      │                  │                  │    │
│  ┌───▼────┐      ┌──────▼──────┐  ┌──────▼──┐ │
│  │ Clean   │      │   Gaming    │  │Research │ │
│  │ VM (100)│      │   VM (101)  │  │ VM (102)│ │
│  │         │      │             │  │         │ │
│  │192.168. │      │ 192.168.    │  │192.168. │ │
│  │100.100  │      │ 100.101     │  │100.102  │ │
│  └─────────┘      └─────────────┘  └─────────┘ │
│      ↓                  ↓                ↓      │
│  ┌────────────────────────────────────────────┐ │
│  │  Virtual Bridge (vmbr-internal)            │ │
│  │  192.168.100.0/24                         │ │
│  │  Gateway: 192.168.100.1 (OPNsense)        │ │
│  └────────────────────────────────────────────┘ │
│                                                  │
│  All traffic between VMs goes through OPNsense! │
│                                                  │
└──────────────────────────────────────────────────┘
       │
       │ Only allowed traffic
       │ exits to external network
       │
    Internet
```

## Installation

### 1. Download OPNsense ISO

```bash
# On Proxmox host
cd /var/lib/vz/template/iso/
wget https://mirrors.evowise.com/opnsense/releases/24.1/OPNsense-24.1-dvd-amd64.iso

# Or alternative mirror
wget https://ftp.opnsense.org/releases/24.1/OPNsense-24.1-dvd-amd64.iso
```

### 2. Create OPNsense MicroVM

Via Proxmox Web UI or CLI:

```bash
# Using template from config/vm-templates/opnsense-microvm.yaml
qm create 50 \
  --name opnsense-microvm \
  --node <hostname> \
  --memory 1024 \
  --cores 2 \
  --sockets 1 \
  --machine q35 \
  --bios ovmf

# Add network interfaces
qm set 50 -net0 virtio,bridge=vmbr0  # WAN
qm set 50 -net1 virtio,bridge=vmbr-internal  # LAN
qm set 50 -net2 virtio,bridge=vmbr-mgmt  # Management (optional)

# Add storage
qm set 50 -scsi0 local-lvm:100,cache=writethrough

# Set CD-ROM for OPNsense ISO
qm set 50 -ide2 /var/lib/vz/template/iso/OPNsense-24.1-dvd-amd64.iso,media=cdrom
```

### 3. Boot and Install OPNsense

```bash
# Start VM
qm start 50

# Connect to console
qm terminal 50
```

Follow OPNsense installation:

1. **Select Install**: "Install OPNsense"
2. **Keyboard**: Select your keyboard layout
3. **Partitioning**: Use automatic partitioning
4. **Network**:
   - WAN (em0): DHCP or static IP
   - LAN (em1): Static IP `192.168.100.1` / `255.255.255.0`
   - Management (em2): Static IP `192.168.99.1` / `255.255.255.0` (optional)
5. **Root Password**: Set strong password (change default immediately!)

### 4. Post-Installation Configuration

After reboot, access OPNsense:

```bash
# Via SSH
ssh root@192.168.100.1 -p 22

# Or via web UI
# Browser: https://192.168.100.1
# (Accept self-signed certificate)
```

## Network Configuration

### Virtual Bridges

Create bridges on Proxmox host:

```bash
# Edit /etc/network/interfaces
nano /etc/network/interfaces
```

Add bridges for VM networks:

```bash
# Internal VM bridge (LAN)
auto vmbr-internal
iface vmbr-internal inet manual
    bridge_ports none
    bridge_stp off
    bridge_fd 0

# Management bridge (optional)
auto vmbr-mgmt
iface vmbr-mgmt inet manual
    bridge_ports none
    bridge_stp off
    bridge_fd 0

# External traffic bridge (WAN)
auto vmbr0
iface vmbr0 inet manual
    bridge_ports <physical-nic>  # Your physical NIC name
    bridge_stp off
    bridge_fd 0
```

Apply changes:

```bash
systemctl restart networking
# Or reboot
reboot
```

### OPNsense Network Setup

Access OPNsense web UI (https://192.168.100.1):

1. **System > Settings > General**:
   - Hostname: `opnsense-fw`
   - Domain: `vm.local`
   - DNS: `1.1.1.1`, `1.0.0.1`

2. **Interfaces > Assignments**:
   - WAN: `em0` (DHCP or static)
   - LAN: `em1` (Static 192.168.100.1/24)
   - OPT1: `em2` (Management 192.168.99.1/24)

3. **Interfaces > LAN**:
   - Enable DHCP Server
   - DHCP Range: `192.168.100.100` - `192.168.100.200`
   - DNS Servers: `192.168.100.1`

## Firewall Rules

### Default Policy

Set default firewall policies:

1. **Firewall > Settings > General**:
   - Default Deny: Enable
   - Drop packets by default: Yes

### Create VM Isolation Rules

In OPNsense Web UI, go to **Firewall > Rules > LAN**:

#### Block Clean VM ↔ Gaming VM

```
Source: 192.168.100.100 (Clean VM)
Destination: 192.168.100.101 (Gaming VM)
Protocol: Any
Action: Block
Direction: in/out
```

#### Block Research VM (Complete Isolation)

```
Source: 192.168.100.102 (Research VM)
Destination: !192.168.100.1 (Anywhere except OPNsense)
Protocol: Any
Action: Block
```

#### Allow Tools VM DNS Only

```
Source: 192.168.100.103 (Tools VM)
Destination: Any
Port: 53 (DNS)
Protocol: UDP
Action: Pass
```

#### Allow Standard Internet

```
Source: LAN
Destination: Any
Protocol: Any
Action: Pass
```

### Traffic Logging

Enable comprehensive logging:

1. **Firewall > Settings > Advanced**:
   - Enable Firewall Logging
   - Log Matches: Yes
   - Log Passes: Yes (for troubleshooting)

2. **System > Log Files > Firewall**:
   - View real-time firewall logs

## VM Configuration for Firewall

### Update VM Default Routes

Each VM must route through OPNsense (192.168.100.1):

#### Linux VMs

Edit `/etc/netplan/00-installer-config.yaml`:

```yaml
network:
  version: 2
  ethernets:
    eth0:
      dhcp4: true
      dhcp4-overrides:
        use-routes: true
      routes:
        - to: default
          via: 192.168.100.1
      gateway4: 192.168.100.1
      nameservers:
        addresses:
          - 192.168.100.1
          - 1.1.1.1
```

Apply:

```bash
sudo netplan apply
```

#### Windows VMs

Set default gateway in Network Settings:

```powershell
netsh int ipv4 set address name="Ethernet" static 192.168.100.100 255.255.255.0 192.168.100.1 1

# Verify
route print
```

## Traffic Monitoring

### View Traffic Statistics

In OPNsense Web UI:

1. **Monitoring > Traffic**:
   - Real-time bandwidth per interface
   - Per-VM traffic (if configured)

2. **Monitoring > Firewall**:
   - View blocked/passed packets
   - See which rules are triggering

### Enable NetFlow/sFlow

For advanced monitoring:

1. **System > Diagnostics > Connection**:
   - Enable NetFlow export (optional)
   - Configure external collector

### Intrusion Detection (Optional)

Enable Suricata IDS/IPS:

1. **Services > Intrusion Detection**:
   - Enable Suricata
   - Download rulesets
   - Configure alert rules

## Advanced Features

### Traffic Shaping

Limit bandwidth per VM:

1. **Firewall > Traffic Shaping**:
   - Create queues for each VM
   - Set bandwidth limits per VLAN/IP

```bash
# Example: Limit Gaming VM to 1Gbps
Source: 192.168.100.101
Queue: 1000Mbit
```

### VPN Access (Optional)

Add VPN for secure remote access:

1. **VPN > OpenVPN > Server**:
   - Configure OpenVPN server
   - Create client certs
   - Allow trusted clients only

### Backup Configuration

Backup OPNsense config regularly:

1. **System > Configuration Backup**:
   - Download backup regularly
   - Store securely
   - Test restoration procedure

## MicroVM Benefits

### Resource Efficiency

OPNsense MicroVM uses minimal resources:

| Resource | Allocation | Typical Usage |
|----------|-----------|---|
| **CPU** | 2 cores | 1-5% (idle) |
| **RAM** | 1-2 GB | 300-600 MB |
| **Storage** | 20 GB | 5-8 GB |
| **Network** | 3 interfaces | <1 Mbps (monitoring) |

### Boot Performance

MicroVM boots in seconds:

```bash
# Start OPNsense
qm start 50
# Ready for traffic in ~30 seconds
```

### Minimal Attack Surface

- Only essential services running
- Smaller disk footprint
- Fewer packages to update
- Single purpose: routing/firewalling

## Troubleshooting

### Cannot Access OPNsense Web UI

```bash
# Check if OPNsense is running
qm status 50

# Check network connectivity
ping 192.168.100.1

# Check port 443
netstat -tuln | grep 443

# Restart web UI
ssh root@192.168.100.1
service lighttpd restart
```

### VMs Cannot Reach Internet

1. **Check WAN interface**: `ifconfig -a` on OPNsense
2. **Verify routing**: `netstat -rn` on OPNsense
3. **Check NAT rules**: Firewall > NAT > Outbound
4. **Default gateway**: Ensure VMs use 192.168.100.1

### VMs Cannot Communicate (Even Allowed)

1. **Check firewall rules**: View blocked packets in logs
2. **Verify VM routes**: `route` or `netstat -rn` on VM
3. **Check ARP**: Might need to clear ARP cache
4. **Restart networking**: On affected VMs

### High Latency/Slow Traffic

1. **Check CPU usage**: OPNsense might be overwhelmed
2. **Enable traffic shaping**: Control bandwidth
3. **Increase VM specs**: Add more cores if needed
4. **Check for bottlenecks**: Monitor interface bandwidth

## Security Best Practices

### Initial Setup

1. **Change default password immediately**
2. **Disable SSH password auth**: Use keys only
3. **Enable HTTPS**: Self-signed cert is fine (for local)
4. **Backup configuration**: Regularly save backups

### Ongoing Security

1. **Keep OPNsense updated**: System > Firmware > Check for Updates
2. **Review firewall logs**: Check for suspicious patterns
3. **Monitor bandwidth**: Detect unusual activity
4. **Test isolation rules**: Verify VMs can't communicate unexpectedly
5. **Audit allowed rules**: Only whitelist necessary traffic

### Network Segmentation

```
Research VM:    192.168.100.102/32 (complete isolation)
Tools VM:       192.168.100.103/32 (restricted to DNS/HTTP/HTTPS)
Gaming VM:      192.168.100.101/32 (blocked from other VMs)
Clean VM:       192.168.100.100/32 (standard internet access)
```

## Integration with Proxmox Automation

### Automatic VM Network Setup

Script to automatically configure VM networks:

```bash
#!/bin/bash
# setup-vm-routes.sh

vm_ips=(
    "192.168.100.100:clean"
    "192.168.100.101:gaming"
    "192.168.100.102:research"
    "192.168.100.103:tools"
)

for entry in "${vm_ips[@]}"; do
    ip="${entry%:*}"
    name="${entry#*:}"

    # Add static IP and routing rules
    # (Configure on each VM)
done
```

## Further Resources

- **OPNsense Docs**: https://docs.opnsense.org/
- **GitHub**: https://github.com/opnsense/core
- **Community Forum**: https://forum.opnsense.org/
- **Firewall Guide**: https://docs.opnsense.org/manual/firewall.html

---

**Next Steps**:
1. Create OPNsense MicroVM in Proxmox
2. Configure firewall rules for VM isolation
3. Update all VMs to route through OPNsense
4. Monitor traffic and adjust rules as needed
5. Enable logging for security research
