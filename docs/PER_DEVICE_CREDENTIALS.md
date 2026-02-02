# Per-Device Credentials Examples

This document provides examples of using the per-device credentials feature.

## Basic Setup

### 1. Using Simple Text File (Same Credentials for All)

**config/devices.txt:**
```
192.168.1.1
192.168.1.2
192.168.1.3
```

**.env:**
```bash
NETWORK_SSH_USERNAME=admin
NETWORK_SSH_PASSWORD=password123
NETWORK_SNMP_COMMUNITY=public
```

All devices will use the same credentials from the `.env` file.

---

## Advanced: Per-Device Credentials

### 2. Using YAML Configuration (Different Credentials per Device)

**config/devices.yaml:**
```yaml
devices:
  # Corporate switch with standard credentials
  - ip: 192.168.1.1
    ssh:
      username: admin
      password: Corp@2024
    snmp:
      community: corporate
      version: 2c

  # Data center switch with SSH key authentication
  - ip: 192.168.1.2
    ssh:
      username: dcadmin
      key_path: /config/keys/datacenter_key
    snmp:
      community: datacenter
      version: 2c

  # Legacy device with SNMP only
  - ip: 192.168.1.3
    snmp:
      community: legacy_ro
      version: 2c

  # High-security device with SNMPv3
  - ip: 10.0.0.1
    ssh:
      username: secadmin
      password: SecureP@ss2024
    snmp:
      version: 3
      username: snmpv3user
      auth_password: Auth1234!
      priv_password: Priv5678!
```

**.env:**
```bash
# Default credentials (used as fallback)
NETWORK_SSH_USERNAME=admin
NETWORK_SSH_PASSWORD=defaultpass
NETWORK_SNMP_COMMUNITY=public

# YAML config file location
NETWORK_DEVICES_YAML=/config/devices.yaml
```

---

## Real-World Scenarios

### Scenario 1: Mixed Vendor Environment

```yaml
devices:
  # Mikrotik switches
  - ip: 192.168.1.10
    ssh:
      username: mikrotik-admin
      password: mt_password
    snmp:
      community: mikrotik_snmp
      version: 2c

  - ip: 192.168.1.11
    ssh:
      username: mikrotik-admin
      password: mt_password
    snmp:
      community: mikrotik_snmp
      version: 2c

  # Cisco switches
  - ip: 192.168.1.20
    ssh:
      username: cisco-admin
      password: cisco_password
    snmp:
      community: cisco_snmp
      version: 2c

  - ip: 192.168.1.21
    ssh:
      username: cisco-admin
      password: cisco_password
    snmp:
      community: cisco_snmp
      version: 2c
```

### Scenario 2: SSH Keys for Production, Passwords for Lab

```yaml
devices:
  # Production devices with SSH keys
  - ip: 10.0.0.1
    ssh:
      username: prod-admin
      key_path: /config/keys/production_key
    snmp:
      community: prod_snmp
      version: 2c

  - ip: 10.0.0.2
    ssh:
      username: prod-admin
      key_path: /config/keys/production_key
    snmp:
      community: prod_snmp
      version: 2c

  # Lab devices with passwords
  - ip: 192.168.100.1
    ssh:
      username: lab-admin
      password: lab123
    snmp:
      community: public
      version: 2c

  - ip: 192.168.100.2
    ssh:
      username: lab-admin
      password: lab123
    snmp:
      community: public
      version: 2c
```

### Scenario 3: SNMP-Only Monitoring Devices

```yaml
devices:
  # Legacy switches without SSH
  - ip: 192.168.1.100
    snmp:
      community: legacy_readonly
      version: 2c

  - ip: 192.168.1.101
    snmp:
      community: legacy_readonly
      version: 2c

  # Modern switches with SNMPv3
  - ip: 192.168.1.110
    snmp:
      version: 3
      username: monitoring
      auth_password: MonAuth@123
      priv_password: MonPriv@456

  - ip: 192.168.1.111
    snmp:
      version: 3
      username: monitoring
      auth_password: MonAuth@123
      priv_password: MonPriv@456
```

---

## Running the Scanner

Once your `devices.yaml` is configured:

```bash
# With Docker Compose
docker-compose up

# Locally
export NETWORK_DEVICES_YAML=/path/to/devices.yaml
./scanner
```

The scanner will:
1. Load the YAML configuration
2. Use per-device credentials for each IP
3. Extract device model information during scanning
4. Generate topology with complete device details

---

## Expected Output

With proper credentials configured, you'll see model information in the JSON output:

```json
{
  "devices": [
    {
      "ip": "192.168.1.1",
      "hostname": "core-switch-01",
      "vendor": "mikrotik",
      "model": "CRS328-24P-4S+",
      "version": "7.11.2 (stable)",
      "scan_status": "success"
    },
    {
      "ip": "192.168.1.2",
      "hostname": "access-switch-01",
      "vendor": "cisco",
      "model": "WS-C2960X-48FPS-L",
      "version": "15.2(7)E3",
      "scan_status": "success"
    }
  ]
}
```

---

## Troubleshooting

### Issue: Credentials not working for specific device

**Check:**
1. Verify the IP address in `devices.yaml` matches exactly
2. Test credentials manually: `ssh admin@192.168.1.1`
3. Check scanner logs for authentication errors

### Issue: YAML not being loaded

**Check:**
1. Verify file path: `echo $NETWORK_DEVICES_YAML`
2. Ensure YAML syntax is correct: `yamllint config/devices.yaml`
3. Check scanner logs: "Loading device configuration from YAML"

### Issue: Model information not extracted

**Possible causes:**
- Device responded via SNMP (limited info)
- Vendor not fully supported (uses generic parser)
- SSH authentication failed (fell back to SNMP)

**Solutions:**
- Enable SSH access on the device
- Configure correct SSH credentials
- Check if device supports LLDP/CDP
