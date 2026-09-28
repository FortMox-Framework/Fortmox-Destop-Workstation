# Storage Encryption & Secure Data Management

> **Alpha/manual guide:** The CLI does not create LUKS devices, encrypted vaults, or encryption keys. Performance and security claims below are targets, not FortMox test results.

## Overview

This guide covers implementing LUKS2 encryption for VM storage, ensuring data protection at rest while maintaining performance.

**Target**: LUKS2 encryption with hardware acceleration; actual security and performance depend on the host and configuration.

## Key Principles

1. **Hardware Acceleration**: AES-NI used for wire-speed encryption (no CPU overhead)
2. **Key Management**: Separate encryption keys per vault
3. **Performance**: <1% CPU overhead with hardware acceleration
4. **Flexibility**: Different encryption levels for different threat levels

## Part 1: Preparation

### Check Hardware Acceleration

```bash
# Check for AES-NI support (Intel/AMD)
grep aes /proc/cpuinfo

# Output should include:
# flags : ... aes aesni ...

# Verify cryptography module
cat /proc/crypto | grep aes
```

### Install Tools

```bash
sudo apt update
sudo apt install -y \
    cryptsetup \
    cryptsetup-initramfs \
    dmsetup \
    keyutils \
    libcryptsetup12
```

## Part 2: Create Encrypted Storage

### Identify Storage Device

```bash
# List block devices
lsblk

# Or with detailed info
fdisk -l

# Example output:
# /dev/sdb     <- Secondary disk (will encrypt)
# /dev/sdc     <- Another disk
```

### Create LUKS2 Encrypted Volume

```bash
# Create LUKS2 volume on /dev/sdb
sudo cryptsetup luksFormat --type luks2 --cipher aes-xts-plain64 --key-size 512 /dev/sdb

# Enter passphrase twice
# WARNING: You will overwrite data on /dev/sdb

# Verify creation
sudo cryptsetup luksDump /dev/sdb
# Should show: LUKS2 format
```

### Open Encrypted Volume

```bash
# Open the encrypted volume
sudo cryptsetup luksOpen /dev/sdb encrypted_vault

# Check the mapping
sudo dmsetup ls
# Output should show: encrypted_vault (xxx:yyy)

# Or check /dev/mapper/
ls -la /dev/mapper/ | grep encrypted
```

### Create Filesystem

```bash
# Create ext4 filesystem (supports TRIM for SSDs)
sudo mkfs.ext4 -m 0 /dev/mapper/encrypted_vault

# Create mount point
sudo mkdir -p /mnt/encrypted_vault

# Mount the volume
sudo mount /dev/mapper/encrypted_vault /mnt/encrypted_vault

# Verify
df -h | grep encrypted_vault
```

### Enable Discard/TRIM for SSD

If using SSD, enable TRIM for better performance:

```bash
# Edit /etc/fstab
sudo nano /etc/fstab

# Add discard option:
# /dev/mapper/encrypted_vault /mnt/encrypted_vault ext4 defaults,discard 0 2

# Or remount with discard
sudo mount -o remount,discard /mnt/encrypted_vault

# Verify
sudo tune2fs -l /dev/mapper/encrypted_vault | grep discard
```

### Persistent Decryption

To auto-decrypt on boot:

```bash
# Create keyfile
sudo dd if=/dev/urandom of=/root/.vault_keyfile bs=1024 count=4
sudo chmod 400 /root/.vault_keyfile

# Add key to LUKS
sudo cryptsetup luksAddKey /dev/sdb /root/.vault_keyfile

# Add to /etc/crypttab
sudo nano /etc/crypttab

# Add line:
# encrypted_vault /dev/sdb /root/.vault_keyfile luks,x-systemd.device-timeout=0
```

Update initramfs:

```bash
sudo update-initramfs -u -k all
```

## Part 3: Storage Vault for VMs

### Create VM Storage

On encrypted volume:

```bash
# Create directory for VM disks
sudo mkdir -p /mnt/encrypted_vault/vmdisks

# Change ownership to Proxmox
sudo chown 107:107 /mnt/encrypted_vault/vmdisks

# Verify
ls -la /mnt/encrypted_vault/ | grep vmdisks
```

### Add Storage to Proxmox

In Proxmox Web UI:
**Datacenter → Storage → Add → Directory**

Or via CLI:

```bash
# Create storage directory
pvesm add dir secure-vault \
    --content images,rootdir \
    --path /mnt/encrypted_vault/vmdisks
```

### Use Encrypted Storage for VMs

When creating VMs, select storage:
- **Storage**: secure-vault (encrypted)
- **Format**: qcow2 (supports snapshots)

Or edit VM config:

```bash
# In /etc/pve/qemu-server/100.conf
virtio0: secure-vault:100/vm-100-disk-0.qcow2,size=50G
```

## Part 4: Multiple Encryption Keys

For compartmentalization, use separate keys:

```bash
# Clean VM vault
sudo cryptsetup luksFormat --type luks2 /dev/sdc
sudo cryptsetup luksOpen /dev/sdc encrypted_clean
sudo mkfs.ext4 /dev/mapper/encrypted_clean
sudo mount /dev/mapper/encrypted_clean /mnt/clean

# Gaming VM vault
sudo cryptsetup luksFormat --type luks2 /dev/sdd
sudo cryptsetup luksOpen /dev/sdd encrypted_gaming
sudo mkfs.ext4 /dev/mapper/encrypted_gaming
sudo mount /dev/mapper/encrypted_gaming /mnt/gaming

# Research VM vault (extra paranoid)
sudo cryptsetup luksFormat --type luks2 --key-size 512 --pbkdf pbkdf2 /dev/sde
sudo cryptsetup luksOpen /dev/sde encrypted_research
sudo mkfs.ext4 /dev/mapper/encrypted_research
sudo mount /dev/mapper/encrypted_research /mnt/research
```

Benefits:
- Compromise of one key doesn't affect others
- Can use different passphrases per vault
- Can destroy one key (shred) without affecting others

## Part 5: Key Management

### Backup Encryption Keys

```bash
# Export LUKS header (for emergency recovery)
sudo cryptsetup luksDump /dev/sdb > /root/vault_header_backup.txt

# Keep in secure location
# Copy to external encrypted storage
```

### Key Rotation (Advanced)

```bash
# Change passphrase
sudo cryptsetup luksChangeKey /dev/sdb

# Add new key
sudo cryptsetup luksAddKey /dev/sdb /path/to/new/keyfile

# Remove old key
sudo cryptsetup luksRemoveKey /dev/sdb
```

### Secure Key Destruction

```bash
# Securely delete key file (multiple passes)
sudo shred -vfz -n 3 /root/.vault_keyfile

# Or using wipe
sudo apt install wipe
sudo wipe /root/.vault_keyfile
```

## Part 6: Performance Optimization

### Encryption Parameters

For gaming/performance-critical:

```bash
# Use XTS mode (faster than CBC)
--cipher aes-xts-plain64

# Reduce key-size if acceptable for use case
--key-size 256  # Still very secure, slightly faster

# Use parallel processing
--pbkdf pbkdf2  # Faster than default argon2
```

### Monitor Performance

```bash
# Check encryption overhead
iostat -x 1 10  # Before and after encryption comparison

# CPU usage
top -p $(pgrep cryptd)

# Should see <1% CPU with AES-NI
```

### Enable Caching

For encrypted VMs, enable disk caching:

```bash
# In /etc/pve/qemu-server/100.conf
virtio0: secure-vault:100/vm-100-disk-0.qcow2,cache=writeback

# WARNING: writeback cache can lose data on crash
# Use cache=writethrough for safer (slightly slower) option
```

## Part 7: Backup & Recovery

### Backup Encrypted VM

```bash
# Create snapshot (even on encrypted storage)
qm snapshot 100 backup-snapshot-$(date +%Y%m%d)

# Export VM
qm exportvmconfig 100 > vm-100-config.txt

# Backup encrypted volume
# Option 1: While mounted
sudo tar -czf /backup/clean-vm-backup.tar.gz /mnt/clean

# Option 2: Via cryptsetup
sudo cryptsetup luksDump /dev/sdc > /backup/clean-luks-header.txt
```

### Restore from Backup

```bash
# Create new encrypted volume
sudo cryptsetup luksFormat --type luks2 /dev/sdc
sudo cryptsetup luksOpen /dev/sdc encrypted_clean
sudo mkfs.ext4 /dev/mapper/encrypted_clean

# Restore data
sudo tar -xzf /backup/clean-vm-backup.tar.gz -C /

# Or restore from VM snapshot
qm importdisk 100 /backup/vm-disk.qcow2 secure-vault
```

## Part 8: LUKS Emergency Access

### Recover Access (Lost Passphrase)

If you lose a passphrase but have keyfile:

```bash
# Use keyfile to open
sudo cryptsetup luksOpen /dev/sdb --key-file=/path/to/keyfile encrypted_vault

# Add new passphrase
sudo cryptsetup luksAddKey /dev/sdb
# Enter existing key (from file via --key-file), then new passphrase
```

### Create Recovery Key

```bash
# Add second passphrase/keyfile for emergency access
sudo cryptsetup luksAddKey /dev/sdb /root/emergency.key

# Store securely (separate location from main key)
sudo shred -vfz -n 3 /root/emergency.key  # When no longer needed
```

## Part 9: Encryption for Different VM Types

### Clean VM (Standard Encryption)

```bash
# Fast encryption for daily work
cryptsetup luksFormat --cipher aes-xts-plain64 --key-size 256 /dev/sdc

# High-security passphrase (16+ characters)
```

### Gaming VM (Performance-Optimized)

```bash
# Fast encryption (less overhead)
cryptsetup luksFormat --cipher aes-xts-plain64 --key-size 256 --pbkdf pbkdf2 /dev/sdd

# Writeback cache for performance
# In Proxmox: cache=writeback
```

### Research VM (Paranoid Encryption)

```bash
# Maximum security
cryptsetup luksFormat --cipher aes-xts-plain64 --key-size 512 /dev/sde

# Slow key derivation (harder to brute force)
# Default: argon2id (memory-hard function)

# Conservative cache settings
# In Proxmox: cache=writethrough
```

## Part 10: Monitoring Encrypted Storage

### Check Encryption Status

```bash
# List all encrypted volumes
sudo cryptsetup status /dev/mapper/encrypted_vault

# Show details
sudo cryptsetup luksDump /dev/sdb

# Check mount status
df -h | grep encrypted
mount | grep encrypted
```

### Monitor Performance

```bash
# Real-time I/O stats
iotop

# Per-device encryption overhead
iostat -dxm 1 | grep sdb  # Watch encrypted device

# Should show <5% additional latency
```

### Automated Decryption Check

Create `/usr/local/bin/check-encrypted.sh`:

```bash
#!/bin/bash

echo "Checking encrypted volumes..."

for volume in encrypted_vault encrypted_clean encrypted_gaming encrypted_research; do
    if [ -e /dev/mapper/$volume ]; then
        echo "✓ $volume: Mounted"
        df /dev/mapper/$volume 2>/dev/null | tail -1
    else
        echo "✗ $volume: NOT mounted"
    fi
done
```

Run it:

```bash
sudo /usr/local/bin/check-encrypted.sh
```

## Part 11: Security Audit

### Verify Encryption Strength

```bash
# Check LUKS version
sudo cryptsetup luksDump /dev/sdb | grep "Version:"
# Should be: 2

# Check cipher algorithm
sudo cryptsetup luksDump /dev/sdb | grep "Cipher name:"
# Should be: aes-xts-plain64

# Check key size
sudo cryptsetup luksDump /dev/sdb | grep "Key size:"
# Should be: 256 bits (or 512 for paranoid)
```

### Performance Verification

```bash
# Benchmark unencrypted
dd if=/dev/zero of=/tmp/test.img bs=1M count=1000

# Benchmark encrypted
dd if=/dev/zero of=/mnt/encrypted_vault/test.img bs=1M count=1000

# Compare speeds - overhead should be <5% with AES-NI
```

## Troubleshooting

### Mount Failures

```bash
# Check if device is recognized
lsblk | grep sdb

# Try manual open
sudo cryptsetup luksOpen /dev/sdb encrypted_vault

# Check mounting
sudo mount /dev/mapper/encrypted_vault /mnt/test

# View errors
dmesg | tail -20
journalctl -xe | tail -20
```

### Performance Issues

```bash
# Check if AES-NI is being used
grep aes-ni /proc/cpuinfo

# If not present, encryption will be slow
# Enable in BIOS/UEFI if available

# Monitor CPU usage
top -p $(pgrep cryptd)
```

### Passphrase Problems

```bash
# Test passphrase
sudo cryptsetup luksOpen /dev/sdb --test-passphrase

# Verify with key file
sudo cryptsetup luksOpen /dev/sdb --key-file=/root/.vault_keyfile

# Check number of key slots
sudo cryptsetup luksDump /dev/sdb | grep "Key Slot"
```

## Summary

### Encryption Strength

| Layer | Method | Protection |
|-------|--------|-----------|
| Storage | LUKS2 AES-256 | Military-grade |
| Transport | Hardware AES-NI | Wire-speed, no overhead |
| Keys | Separate per vault | Compartmentalization |
| Backups | Encrypted separately | Data at rest protection |

### Performance Impact

- **CPU**: <1% with AES-NI
- **Throughput**: <5% reduction
- **Latency**: <1ms additional
- **Result**: No practical performance loss

### Security Posture

✅ Data encrypted at rest
✅ Hardware-accelerated (no performance cost)
✅ Separate keys per VM compartment
✅ Can destroy keys individually
✅ Supports snapshots and backups

---

**Last Updated**: 2024-11-17
**CLI status**: Storage encryption and vault management are not implemented.
