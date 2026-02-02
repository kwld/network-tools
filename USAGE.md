# Network Tools - Quick Start Guide

## Quick Setup (5 minutes)

### 1. Clone the repository
```bash
git clone https://github.com/kwld/network-tools.git
cd network-tools
```

### 2. Configure credentials
```bash
cp .env.example .env
# Edit .env with your network device credentials
nano .env
```

### 3. Add devices to scan
```bash
# Edit config/devices.txt and add your device IPs (one per line)
nano config/devices.txt
```

Example `config/devices.txt`:
```
192.168.1.1
192.168.1.2
192.168.1.10
10.0.0.1
```

### 4. Run the scanner

#### Option A: Docker Compose (Recommended)
```bash
docker compose up --build
```

Results will be in `./output/` directory:
- `topology_*.json` - Complete topology data
- `topology_*.svg` - Visual network diagram
- `topology_*.dot` - GraphViz source
- `summary_*.txt` - Scan summary

#### Option B: Local Development
```bash
# Install Go 1.20+ and GraphViz
sudo apt-get install graphviz  # Ubuntu/Debian
brew install graphviz          # macOS

# Build and run
make build
make run
```

## Configuration Options

Edit `.env` to customize:

### Authentication
```bash
NETWORK_SSH_USERNAME=admin
NETWORK_SSH_PASSWORD=yourpassword
# NETWORK_SSH_KEY_PATH=/path/to/key  # For key-based auth
```

### SNMP Settings
```bash
NETWORK_SNMP_COMMUNITY=public
NETWORK_SNMP_VERSION=2c  # or "3" for SNMPv3
```

### Scanner Performance
```bash
NETWORK_SCAN_TIMEOUT=30         # Seconds per device
NETWORK_CONCURRENT_SCANS=5      # Parallel scans
```

### Diagram Format
```bash
NETWORK_DIAGRAM_FORMAT=svg      # Options: svg, png, pdf
NETWORK_DIAGRAM_LAYOUT=dot      # Options: dot, neato, fdp, circo
```

## Example Workflows

### Basic Network Scan
```bash
# 1. Add device IPs to config/devices.txt
echo "192.168.1.1" > config/devices.txt
echo "192.168.1.2" >> config/devices.txt

# 2. Set credentials in .env
export NETWORK_SSH_USERNAME=admin
export NETWORK_SSH_PASSWORD=password

# 3. Run scan
docker compose up

# 4. View results
ls -lh output/
cat output/summary_*.txt
```

### Large Network Scan (100+ devices)
```bash
# Increase concurrency and timeout
cat > .env << 'END'
NETWORK_SSH_USERNAME=admin
NETWORK_SSH_PASSWORD=password
NETWORK_SCAN_TIMEOUT=60
NETWORK_CONCURRENT_SCANS=20
NETWORK_DEVICES_FILE=/config/devices.txt
NETWORK_OUTPUT_DIR=/output
END

# Run scan
docker compose up
```

### Using SSH Key Authentication
```bash
# Copy SSH key to config directory
cp ~/.ssh/id_rsa config/network_key
chmod 600 config/network_key

# Update docker-compose.yml to mount the key
# Then configure .env
cat > .env << 'END'
NETWORK_SSH_USERNAME=admin
NETWORK_SSH_KEY_PATH=/config/network_key
END

docker compose up
```

## Troubleshooting

### "Connection timeout"
- Verify network connectivity: `ping 192.168.1.1`
- Increase timeout: `NETWORK_SCAN_TIMEOUT=60`
- Check firewall rules

### "Authentication failed"
- Verify credentials are correct
- Try SSH manually: `ssh admin@192.168.1.1`
- Check if SSH is enabled on device

### "GraphViz not found"
- When running locally, install GraphViz:
  - Ubuntu/Debian: `sudo apt-get install graphviz`
  - macOS: `brew install graphviz`
  - Windows: Download from graphviz.org
- Docker image includes GraphViz automatically

### No devices found
- Check `config/devices.txt` exists and has IPs
- Verify file path in .env: `NETWORK_DEVICES_FILE=/config/devices.txt`

## Output Files Explained

### topology_*.json
Complete machine-readable topology data including:
- Device information (IP, hostname, vendor, model)
- Port details (name, status, speed)
- Neighbor relationships (LLDP/CDP)
- Connection mappings

### topology_*.svg
Visual network diagram showing:
- Devices as colored boxes (green=success, yellow=partial, red=failed)
- Connections as lines between devices
- Port labels on each connection
- Device information (hostname, IP, model)

### topology_*.dot
GraphViz DOT source file for customization:
- Can be edited manually
- Regenerate diagram: `dot -Tsvg topology.dot -o custom.svg`

### summary_*.txt
Human-readable summary including:
- Total devices scanned
- Success/failure counts
- Port and connection statistics
- List of failed IPs

## Advanced Usage

### Custom GraphViz Layout
```bash
# Try different layouts for better visualization
NETWORK_DIAGRAM_LAYOUT=neato  # Force-directed layout
NETWORK_DIAGRAM_LAYOUT=fdp    # Force-directed with springs
NETWORK_DIAGRAM_LAYOUT=circo  # Circular layout
```

### Export to PNG or PDF
```bash
NETWORK_DIAGRAM_FORMAT=png
# or
NETWORK_DIAGRAM_FORMAT=pdf
```

### Scheduled Scans (with cron)
```bash
# Add to crontab
0 2 * * * cd /path/to/network-tools && docker compose up > /var/log/network-scan.log 2>&1
```

### Compare Topologies
```bash
# Run scans at different times
docker compose up  # Creates topology_2024-01-01_10-00-00.json
# ... wait ...
docker compose up  # Creates topology_2024-01-01_12-00-00.json

# Compare with jq
jq '.summary' output/topology_*.json
```

## Security Best Practices

1. **Protect .env file**
   ```bash
   chmod 600 .env
   ```

2. **Use SSH keys instead of passwords**
   ```bash
   NETWORK_SSH_KEY_PATH=/config/network_key
   ```

3. **Prefer SNMPv3 over v2c**
   ```bash
   NETWORK_SNMP_VERSION=3
   NETWORK_SNMP_V3_USER=secureuser
   NETWORK_SNMP_V3_AUTH_PASS=authpass
   NETWORK_SNMP_V3_PRIV_PASS=privpass
   ```

4. **Never commit .env to git**
   - Already in .gitignore
   - Use .env.example as template

## Support

- GitHub Issues: https://github.com/kwld/network-tools/issues
- Documentation: README.md
- Examples: This guide

## Next Steps

After your first successful scan:
1. Review the generated topology diagram
2. Check JSON data for accuracy
3. Adjust scanner settings for your network size
4. Set up scheduled scans if needed
5. Integrate with monitoring tools using JSON output
