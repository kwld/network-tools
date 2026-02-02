# Network Topology Mapping System

A comprehensive network topology mapping system that scans network devices (switches, routers) and creates visual maps of connections per port. Built with Go and designed for easy deployment with Docker Compose.

## Features

### Core Functionality
- **Multi-Device Scanning**: Scan multiple network devices concurrently
- **Multi-Protocol Support**: SSH, SNMPv2c, and SNMPv3
- **Vendor Support**: Mikrotik, Cisco, Motorola, and generic devices with LLDP/CDP
- **Connection Discovery**: Automatically map physical connections via LLDP and CDP
- **Visual Diagrams**: Generate network topology diagrams using GraphViz
- **Data Export**: JSON export of complete topology data
- **Concurrent Scanning**: Configurable parallel device scanning for faster results

### Supported Protocols
- **SSH**: Primary method for detailed device information
- **SNMP v2c**: Fallback for devices without SSH access
- **SNMP v3**: Secure SNMP with authentication and encryption
- **LLDP**: Link Layer Discovery Protocol (IEEE 802.1AB)
- **CDP**: Cisco Discovery Protocol

### Supported Vendors
- **Mikrotik**: RouterOS devices
- **Cisco**: IOS, IOS-XE devices
- **Motorola**: Network switches
- **Generic**: Any device supporting LLDP via SNMP

## Prerequisites

- **Go 1.20+** (for local development)
- **Docker & Docker Compose** (for containerized deployment)
- **GraphViz** (automatically included in Docker image)
- Network access to target devices
- Valid credentials for SSH/SNMP access

## Installation

### Option 1: Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/kwld/network-tools.git
cd network-tools
```

2. Create environment configuration:
```bash
cp config/config.example.env .env
```

3. Edit `.env` with your credentials:
```bash
NETWORK_SSH_USERNAME=admin
NETWORK_SSH_PASSWORD=your_password
NETWORK_SNMP_COMMUNITY=public
# ... other settings
```

4. Add device IPs to `config/devices.txt`:
```
192.168.1.1
192.168.1.2
192.168.1.10
```

5. Run the scanner:
```bash
docker-compose up --build
```

### Option 2: Local Development

1. Install Go 1.20+

2. Install GraphViz:
```bash
# Ubuntu/Debian
sudo apt-get install graphviz

# macOS
brew install graphviz

# Alpine Linux
apk add graphviz
```

3. Clone and build:
```bash
git clone https://github.com/kwld/network-tools.git
cd network-tools
go mod download
make build
```

4. Configure environment variables and run:
```bash
export NETWORK_SSH_USERNAME=admin
export NETWORK_SSH_PASSWORD=password
# ... other variables
make run
```

## Configuration

### Environment Variables

All configuration is done via environment variables. See `config/config.example.env` for a complete list.

#### Authentication
- `NETWORK_SSH_USERNAME`: SSH username (default: admin)
- `NETWORK_SSH_PASSWORD`: SSH password
- `NETWORK_SSH_KEY_PATH`: Path to SSH private key (optional)

#### SNMP Configuration
- `NETWORK_SNMP_COMMUNITY`: SNMPv2c community string (default: public)
- `NETWORK_SNMP_VERSION`: SNMP version: 2c or 3 (default: 2c)
- `NETWORK_SNMP_V3_USER`: SNMPv3 username
- `NETWORK_SNMP_V3_AUTH_PASS`: SNMPv3 authentication password
- `NETWORK_SNMP_V3_PRIV_PASS`: SNMPv3 privacy password

#### Scanner Settings
- `NETWORK_DEVICES_FILE`: Path to device list file (default: /config/devices.txt)
- `NETWORK_OUTPUT_DIR`: Output directory (default: /output)
- `NETWORK_SCAN_TIMEOUT`: Timeout per device in seconds (default: 30)
- `NETWORK_CONCURRENT_SCANS`: Number of concurrent scans (default: 5)

#### Visualization
- `NETWORK_DIAGRAM_FORMAT`: Output format: svg, png, pdf (default: svg)
- `NETWORK_DIAGRAM_LAYOUT`: GraphViz layout: dot, neato, fdp, circo (default: dot)

### Device List File

Create a file with one IP address per line:

```
# Example devices.txt
192.168.1.1
192.168.1.2
192.168.1.10

# Comments are allowed
# 10.0.0.1
```

## Usage

### Using Docker Compose

```bash
# Build and run
docker-compose up --build

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f

# Clean up
docker-compose down
```

### Using Makefile

```bash
# Build the scanner
make build

# Run locally
make run

# Build Docker image
make docker-build

# Run with Docker
make docker-run

# Clean up
make clean
```

### Direct Execution

```bash
# Build
go build -o scanner ./cmd/scanner

# Run with environment variables
export NETWORK_SSH_USERNAME=admin
export NETWORK_SSH_PASSWORD=password
./scanner
```

## Output

The scanner generates several output files in the configured output directory:

### Generated Files

1. **topology_YYYY-MM-DD_HH-MM-SS.json**: Complete topology data in JSON format
2. **topology_YYYY-MM-DD_HH-MM-SS.svg**: Visual network diagram (or png/pdf)
3. **topology_YYYY-MM-DD_HH-MM-SS.dot**: GraphViz DOT source file
4. **summary_YYYY-MM-DD_HH-MM-SS.txt**: Text summary of scan results

### JSON Structure

```json
{
  "devices": [
    {
      "ip": "192.168.1.1",
      "hostname": "switch01",
      "vendor": "mikrotik",
      "model": "CRS328",
      "ports": [
        {
          "name": "ether1",
          "status": "up",
          "speed": "1Gbps",
          "neighbor": {
            "device_id": "switch02",
            "port_id": "ether2",
            "protocol": "LLDP"
          }
        }
      ]
    }
  ],
  "connections": [
    {
      "source_device": "switch01",
      "source_port": "ether1",
      "target_device": "switch02",
      "target_port": "ether2",
      "protocol": "LLDP"
    }
  ],
  "summary": {
    "total_devices": 5,
    "success_devices": 4,
    "failed_devices": 1,
    "total_connections": 8
  }
}
```

## Architecture

### Project Structure

```
network-tools/
├── cmd/scanner/          # Main application
├── internal/
│   ├── scanner/         # Device scanning logic
│   ├── parser/          # Vendor-specific parsers
│   ├── mapper/          # Topology mapping
│   └── visualizer/      # Diagram generation
├── pkg/models/          # Data models
├── config/              # Configuration files
└── output/              # Generated output (gitignored)
```

### How It Works

1. **Device Scanning**: Reads IP list and scans devices concurrently
2. **Protocol Detection**: Tries SSH first, falls back to SNMP
3. **Vendor Detection**: Identifies device vendor from responses
4. **Data Parsing**: Uses vendor-specific parsers to extract port and neighbor info
5. **Topology Building**: Maps connections between devices using LLDP/CDP data
6. **Visualization**: Generates network diagrams using GraphViz
7. **Export**: Saves all data to JSON and visual formats

## Troubleshooting

### Common Issues

**GraphViz not found**
```
Error: graphviz 'dot' command not found
```
Solution: Install GraphViz or use Docker image which includes it.

**Connection timeout**
```
Error: failed to connect: i/o timeout
```
Solution: 
- Verify network connectivity to devices
- Increase `NETWORK_SCAN_TIMEOUT`
- Check firewall rules

**Authentication failed**
```
Error: failed to connect: ssh: handshake failed
```
Solution:
- Verify SSH credentials in `.env`
- Try SSH key authentication with `NETWORK_SSH_KEY_PATH`
- Check device SSH configuration

**No SNMP response**
```
Error: failed to connect: context deadline exceeded
```
Solution:
- Verify SNMP community string
- Check SNMP version (v2c vs v3)
- Ensure SNMP is enabled on devices

### Debug Mode

For verbose logging, run with:
```bash
docker-compose up 2>&1 | tee scan.log
```

## Security Considerations

### Best Practices

1. **Credentials**: Never commit credentials to version control
2. **SSH Keys**: Use SSH key authentication when possible
3. **SNMPv3**: Prefer SNMPv3 with encryption over v2c
4. **File Permissions**: Ensure `.env` has restricted permissions (600)
5. **Network Security**: Run scanner from secure management network

### Secure Credential Storage

```bash
# Set restrictive permissions on .env
chmod 600 .env

# Use Docker secrets (optional)
docker secret create ssh_password password.txt
```

## Testing

### With Virtual Devices

For testing, you can use:
- **GNS3**: Network simulation software
- **EVE-NG**: Emulated Virtual Environment
- **ContainerLab**: Container-based network labs

### Example Test Setup

```bash
# Create test devices list
cat > config/devices.txt << EOF
# Test devices
192.168.56.10
192.168.56.11
EOF

# Run scan
docker-compose up
```

## Performance Tuning

### Concurrent Scanning

Adjust parallel scans based on your system:
```bash
NETWORK_CONCURRENT_SCANS=10  # More concurrent scans
NETWORK_SCAN_TIMEOUT=60      # Longer timeout for slow devices
```

### Large Networks

For networks with 100+ devices:
- Increase `NETWORK_CONCURRENT_SCANS` to 10-20
- Use faster protocols (SSH over SNMP)
- Split into multiple device files and run parallel scanners

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-vendor`)
3. Commit your changes (`git commit -am 'Add support for new vendor'`)
4. Push to the branch (`git push origin feature/new-vendor`)
5. Create a Pull Request

### Adding New Vendor Support

To add support for a new vendor:

1. Create a new parser in `internal/parser/`:
```go
type NewVendorParser struct{}

func (p *NewVendorParser) GetCommands() []string {
    return []string{"show version", "show interfaces"}
}

func (p *NewVendorParser) Parse(outputs map[string]string, device *models.Device) error {
    // Parse vendor-specific output
    return nil
}
```

2. Register the parser in `internal/parser/parser.go`

## License

MIT License - see LICENSE file for details

## Acknowledgments

- GraphViz for visualization
- GoSNMP library for SNMP support
- Go SSH library for SSH connections

## Support

For issues and questions:
- GitHub Issues: https://github.com/kwld/network-tools/issues
- Documentation: This README

## Roadmap

Future enhancements:
- [ ] Web UI for configuration and results
- [ ] Real-time monitoring and alerting
- [ ] Historical topology comparison
- [ ] API endpoints for integration
- [ ] Support for more vendors (Juniper, Arista, etc.)
- [ ] Database storage option
- [ ] Automated scheduled scans
