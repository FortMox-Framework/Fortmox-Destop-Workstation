# Advanced Network Isolation & Firewall Rules

> **Alpha/manual guide:** FortMox does not currently create the bridges, firewall rules, routing, VLANs, IDS/IPS, or isolation described here. The topology below is a target design and has not been validated by the CLI.

## Overview

This guide covers advanced network isolation techniques beyond basic OPNsense firewall rules, including VLAN segmentation, traffic shaping, IDS/IPS integration, and zero-trust firewall rules.

**Target**: Network compartmentalization with defense-in-depth

## Network Topology Review

```
┌─────────────────────────────────────────────────┐
│          Proxmox Host (Hypervisor)             │
│                                                │
│  ┌──────────────────────────────────────────┐ │
│  │    OPNsense MicroVM (vmid: 50)          │ │
│  │    - Central firewall/router             │ │
│  │    - Traffic enforcement                 │ │
│  │    - Logging & monitoring                │ │
│  └──────────────────────────────────────────┘ │
│           │                                   │
│  ┌────────┼────────────────────────┐        │
│  │        │                        │         │
│  ▼        ▼                        ▼         │
│┌────────┬──────────┐    ┌──────────┬───────┐│
││ Clean  │ Gaming   │    │ Research │ Tools ││
││ VM(100)│ VM(101)  │    │ VM(102)  │VM(103)││
│└────────┴──────────┘    └──────────┴───────┘│
│                                                │
│  Virtual Bridges:                            │
│  - vmbr0 (WAN/external)                      │
│  - vmbr-internal (VM communication via OPN)  │
│  - vmbr-isolated (Research VM - air-gapped)  │
│                                                │
└─────────────────────────────────────────────────┘
```

## Part 1: VLAN Configuration

### VLAN Setup

Create multiple VLANs for different threat levels:

```
VLAN 100: Clean VM (normal internet)
VLAN 101: Gaming VM (normal internet)
VLAN 102: Research VM (no internet - air-gapped)
VLAN 103: Tools VM (isolated test network)
VLAN 200: Management (Proxmox admin access)
```

### Proxmox VLAN Configuration

Edit `/etc/network/interfaces`:

```bash
sudo nano /etc/network/interfaces
```

Add VLAN definitions:

```
# Physical interface
auto enp0s3
iface enp0s3 inet manual

# Bridge for WAN (external)
auto vmbr0
iface vmbr0 inet dhcp
    bridge-ports enp0s3
    bridge-stp off
    bridge-fd 0

# VLAN bridges for VMs
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

Restart networking:

```bash
sudo systemctl restart networking
```

### Assign VMs to VLANs

Edit `/etc/pve/qemu-server/100.conf` (Clean VM):

```
net0: virtio=XX:XX:XX:XX:XX:XX,bridge=vmbr100
```

Edit `/etc/pve/qemu-server/101.conf` (Gaming VM):

```
net0: virtio=XX:XX:XX:XX:XX:XX,bridge=vmbr101
```

Edit `/etc/pve/qemu-server/102.conf` (Research VM):

```
net0: virtio=XX:XX:XX:XX:XX:XX,bridge=vmbr102
```

## Part 2: OPNsense Advanced Firewall Rules

### Zero-Trust Firewall Philosophy

**Rule**: Block all traffic by default, allow only what's explicitly needed

### Access OPNsense Web UI

```bash
# From host
firefox https://192.168.100.1
# Default: admin / opnsense

# Or SSH into OPNsense
ssh root@192.168.100.1
```

### Create Firewall Rules

**Principle**:
1. Allow internal services (DNS, DHCP from OPNsense)
2. Allow specific inter-VLAN traffic (Clean ↔ OPNsense)
3. Block all other inter-VLAN traffic
4. Allow external internet only from approved VLANs
5. Block Research VM from all external traffic (air-gapped)

### Rules Configuration (Web UI)

Navigate to: **Firewall → Rules → LAN**

#### Rule Set for Clean VM (VLAN 100)

```
1. Allow | VLAN100 to OPNsense DNS (TCP/UDP 53)
2. Allow | VLAN100 to OPNsense DHCP (UDP 67/68)
3. Allow | VLAN100 to WAN (any protocol)
4. Block | VLAN100 to VLAN101 (any) - Block Gaming VM
5. Block | VLAN100 to VLAN102 (any) - Block Research VM
6. Block | VLAN100 to VLAN103 (any) - Block Tools VM
7. Block | VLAN100 to VLAN200 (any) - Block Management
```

#### Rule Set for Gaming VM (VLAN 101)

```
1. Allow | VLAN101 to OPNsense DNS (TCP/UDP 53)
2. Allow | VLAN101 to OPNsense DHCP (UDP 67/68)
3. Allow | VLAN101 to WAN (any protocol)
4. Block | VLAN101 to VLAN100 (any) - Block Clean VM
5. Block | VLAN101 to VLAN102 (any) - Block Research VM
6. Block | VLAN101 to VLAN103 (any) - Block Tools VM
```

#### Rule Set for Research VM (VLAN 102) - Air-Gapped

```
1. Allow | VLAN102 to OPNsense DNS (TCP/UDP 53) - OPTIONAL
2. Allow | VLAN102 to OPNsense DHCP (UDP 67/68)
3. Block | VLAN102 to WAN (ALL) - NO INTERNET
4. Block | VLAN102 to VLAN100 (any)
5. Block | VLAN102 to VLAN101 (any)
6. Block | VLAN102 to VLAN103 (any)
7. Block | VLAN102 to VLAN200 (any)
```

#### Rule Set for Tools VM (VLAN 103)

```
1. Allow | VLAN103 to OPNsense DNS (TCP/UDP 53)
2. Allow | VLAN103 to OPNsense DHCP (UDP 67/68)
3. Allow | VLAN103 to VLAN100 HTTP (TCP 80)
4. Allow | VLAN103 to VLAN100 HTTPS (TCP 443)
5. Allow | VLAN103 to VLAN101 HTTP (TCP 80) - Test Gaming VM
6. Block | VLAN103 to WAN (ANY) - Tools stays internal
7. Block | All else
```

### Via SSH/CLI

```bash
# SSH into OPNsense
ssh root@192.168.100.1

# List current rules
pfctl -sr

# Add rule blocking Clean-to-Gaming traffic
echo "block drop in from <VLAN100> to <VLAN101>" | pfctl -f -

# Save configuration
# Rules are auto-saved in web UI
```

## Part 3: Network Monitoring & Logging

### Enable NetFlow Monitoring

In OPNsense Web UI:
**System → Settings → Logging** → Enable NetFlow

### View Traffic Logs

**Firewall → Insights → Traffic Graph**

Or via CLI:

```bash
# Real-time firewall logs
tail -f /var/log/filter.log

# Analyze traffic
# Show all blocked traffic
grep "block" /var/log/filter.log | tail -20

# Show connections by VLAN
grep "VLAN100" /var/log/filter.log
```

### Deep Packet Inspection (DPI) with Suricata

Install in OPNsense:
**System → Firmware → Plugins → os-suricata**

Configure:
**Services → IDS/IPS (Suricata)**

Rules:
- Enable ET Open ruleset (free IDS rules)
- Alert on suspicious patterns
- Log all matches

## Part 4: Traffic Shaping (QoS)

### Use Case

Prevent single VM from consuming all bandwidth:

```
Total: 100 Mbps
- Clean VM: 30 Mbps max
- Gaming VM: 50 Mbps max
- Research VM: 5 Mbps (if internet allowed)
- Tools VM: 15 Mbps
```

### Configure Traffic Shaping

In OPNsense:
**Firewall → Traffic Shaping → Queues**

Create queue:
```
Name: Gaming_Queue
Interface: LAN
Bandwidth: 50 Mbps
Priority: High
```

Assign rules to queue:
**Firewall → Rules** → Edit rule → Queue selection

## Part 5: Air-Gapping Research VM

### Complete Network Isolation

For maximum isolation:

```bash
# Remove network device from Research VM entirely
# Edit /etc/pve/qemu-server/102.conf

# Remove or comment out:
# net0: virtio=...

# VM now has NO network access at all
# Can only interact via USB (for transferring samples)
```

### USB Device Isolation

For transferring malware samples safely:

```bash
# Pass USB devices to Research VM only
# In /etc/pve/qemu-server/102.conf
usb0: host=1234:5678  # USB device ID
```

**WARNING**: Ensure USB devices have proper IOMMU isolation

## Part 6: DMZ (Demilitarized Zone) Setup

For isolated testing:

```
Create VLAN 104 for "DMZ"
- Accessible from Tools VM
- Can be accessed FROM WAN (optional)
- Isolated from all other VMs
- Heavily monitored
```

Configure:

```bash
# In OPNsense, create new firewall rules
# VLAN104 ← VLAN103 (Tools can access DMZ)
# VLAN104 ← WAN (Optional: allow inbound testing)
# All other access: BLOCKED
```

## Part 7: VPN Integration

### Allow VPN-only Traffic

For research VMs that need controlled network:

```bash
# Only allow traffic through VPN tunnel
# In OPNsense:
# 1. Configure OpenVPN client
# 2. Route VLAN102 traffic through VPN
# 3. Block all non-VPN traffic from VLAN102
```

### Setup VPN Rule

```bash
# Kill switch: If VPN drops, block all traffic
# In OPNsense Firewall Rules:
# Block VLAN102 to WAN unless VPN_interface is up
```

## Part 8: Network Monitoring Tools

### From Host

```bash
# Monitor bridge traffic
brctl show  # Show bridges
tcpdump -i vmbr100  # Capture VLAN100 traffic
```

### From Proxmox CLI

```bash
# VM network statistics
qm status 100  # Check VM status

# IP address assignment
ip a show vmbr100

# Route information
ip route show
```

### From OPNsense

**Dashboard → Traffic Graph** - Real-time visualization
**Firewall → Insights → Traffic** - Detailed analysis
**Services → IDS/IPS → Alerts** - Security events

## Part 9: Testing & Verification

### Test Rule Enforcement

From Clean VM:

```bash
# Try to reach Gaming VM (should fail)
ping 192.168.101.100  # Should timeout/fail
# Connection refused

# Try to reach OPNsense (should work)
ping 192.168.100.1  # Should succeed
# 64 bytes from 192.168.100.1: ...

# Try internet (should work)
ping 8.8.8.8  # Should succeed (via OPNsense gateway)
```

### Test Research VM Isolation

```bash
# From Research VM
ping 8.8.8.8  # Should fail (air-gapped)
ping 192.168.100.100  # Should fail (isolated)
ping 192.168.100.1  # May work (DNS/DHCP) or fail

# Check routing table
ip route  # Should only show local network
```

### Test Traffic Logging

```bash
# Trigger blocked traffic
# From Clean VM try to reach Gaming VM
traceroute 192.168.101.100

# Check logs in OPNsense
# Should see "BLOCKED" entry in traffic logs
tail -f /var/log/filter.log | grep "block"
```

## Part 10: Monitoring & Alerting

### Enable Alerts on Suspicious Activity

In OPNsense **Services → IDS/IPS (Suricata)**:

```
- Alert on port scans
- Alert on malware signatures (if using ET Pro)
- Alert on unusual protocols
- Log all dropped packets
```

### Create Custom Rules

For unusual network patterns:

```bash
# Alert if Research VM tries to phone home
# (This shouldn't happen if properly isolated)

# Create Suricata rule:
alert ip from $HOME_NET to any (msg:"Research VM Escape Attempt"; classtype:policy-violation; sid:1000001;)
```

## Troubleshooting Network Issues

### VMs Can't Reach Each Other (Expected)

This is normal! They're isolated. To allow controlled communication:

```bash
# Temporarily create specific rule
# In OPNsense: Allow VLAN100 to VLAN101 on specific port
# Example: Allow Clean VM to ping Gaming VM
```

### VM Can't Access Internet

```bash
# Check:
1. Is OPNsense WAN connection active?
   # In OPNsense: System → Interfaces → WAN
2. Is firewall rule allowing to WAN?
3. Can OPNsense itself reach internet?
   # ssh root@192.168.100.1
   # ping 8.8.8.8
```

### High Latency Through OPNsense

```bash
# Check traffic shaping
# Reduce queue size

# Check CPU usage in OPNsense
# May need to increase cores/memory
```

## Performance Impact

Network isolation through OPNsense adds:
- **Latency**: +2-5ms per hop
- **CPU**: 5-15% (minimal MicroVM)
- **Memory**: 1GB baseline
- **Throughput**: No significant loss (most networks are I/O limited)

## Security Summary

### Defense-in-Depth Network Model

1. **Hardware Isolation**: IOMMU groups
2. **VM Isolation**: Each VM is separate QEMU process
3. **Network Isolation**: VLAN segmentation
4. **Firewall**: Zero-trust rules (OPNsense)
5. **Monitoring**: NetFlow + Suricata IDS
6. **Air-Gapping**: Research VM has no network
7. **Traffic Shaping**: Prevent DoS from VMs
8. **Logging**: All traffic logged for audit

### Intended Result (Not Enforced by the Current CLI)

Compromise of one VM:
- ✅ Cannot directly affect other VMs
- ✅ Cannot escape to network (firewall blocks it)
- ✅ Cannot do distributed attacks (isolated)
- ✅ All traffic is logged (forensic analysis possible)

---

**Last Updated**: 2024-11-17
**CLI status**: Network and firewall deployment/testing are not implemented.
