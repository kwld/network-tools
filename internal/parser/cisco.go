package parser

import (
	"regexp"
	"strings"

	"github.com/kwld/network-tools/pkg/models"
)

// CiscoParser handles Cisco devices
type CiscoParser struct{}

func (p *CiscoParser) GetVendorName() string {
	return "cisco"
}

func (p *CiscoParser) GetCommands() []string {
	return []string{
		"show version",
		"show interfaces status",
		"show cdp neighbors detail",
		"show lldp neighbors detail",
	}
}

func (p *CiscoParser) Parse(outputs map[string]string, device *models.Device) error {
	device.Vendor = "cisco"

	// Parse show version
	if output, ok := outputs["show version"]; ok {
		p.parseVersion(output, device)
	}

	// Parse interfaces
	if output, ok := outputs["show interfaces status"]; ok {
		device.Ports = p.parseInterfaces(output)
	}

	// Parse CDP neighbors
	if output, ok := outputs["show cdp neighbors detail"]; ok {
		p.parseNeighbors(output, device, "CDP")
	}

	// Parse LLDP neighbors
	if output, ok := outputs["show lldp neighbors detail"]; ok {
		p.parseNeighbors(output, device, "LLDP")
	}

	return nil
}

func (p *CiscoParser) parseVersion(output string, device *models.Device) {
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Extract hostname
		if strings.HasPrefix(line, "hostname") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				device.Hostname = parts[1]
			}
		}
		
		// Extract version
		if strings.Contains(strings.ToLower(line), "version") && strings.Contains(line, "Version") {
			if match := regexp.MustCompile(`Version\s+([^\s,]+)`).FindStringSubmatch(line); len(match) > 1 {
				device.Version = match[1]
			}
		}
		
		// Extract model
		if strings.Contains(line, "cisco") && strings.Contains(line, "bytes of memory") {
			if match := regexp.MustCompile(`cisco\s+([^\s]+)`).FindStringSubmatch(line); len(match) > 1 {
				device.Model = match[1]
			}
		}
	}
}

func (p *CiscoParser) parseInterfaces(output string) []models.Port {
	var ports []models.Port
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Port") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		port := models.Port{
			Name:   fields[0],
			Status: "unknown",
		}

		if len(fields) > 1 {
			status := strings.ToLower(fields[1])
			if strings.Contains(status, "connected") {
				port.Status = "up"
			} else if strings.Contains(status, "notconnect") || strings.Contains(status, "disabled") {
				port.Status = "down"
			}
		}

		if len(fields) > 3 {
			port.Speed = fields[3]
		}

		ports = append(ports, port)
	}

	return ports
}

func (p *CiscoParser) parseNeighbors(output string, device *models.Device, protocol string) {
	lines := strings.Split(output, "\n")
	var currentNeighbor *models.Neighbor
	var currentPort string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Device ID
		if strings.HasPrefix(line, "Device ID:") {
			if currentNeighbor != nil && currentPort != "" {
				p.attachNeighborToPort(device, currentPort, currentNeighbor)
			}
			currentNeighbor = &models.Neighbor{Protocol: protocol}
			currentNeighbor.DeviceID = strings.TrimSpace(strings.TrimPrefix(line, "Device ID:"))
			currentPort = ""
		}

		if currentNeighbor == nil {
			continue
		}

		// Interface (local port)
		if strings.Contains(line, "Interface:") {
			if match := regexp.MustCompile(`Interface:\s*([^,]+)`).FindStringSubmatch(line); len(match) > 1 {
				currentPort = strings.TrimSpace(match[1])
			}
			// Remote port
			if match := regexp.MustCompile(`Port ID.*:\s*(.+)`).FindStringSubmatch(line); len(match) > 1 {
				currentNeighbor.PortID = strings.TrimSpace(match[1])
			}
		}

		// IP address
		if strings.Contains(line, "IP address:") {
			if match := regexp.MustCompile(`IP address:\s*([0-9.]+)`).FindStringSubmatch(line); len(match) > 1 {
				currentNeighbor.IPAddress = match[1]
			}
		}

		// Platform
		if strings.HasPrefix(line, "Platform:") {
			currentNeighbor.Platform = strings.TrimSpace(strings.TrimPrefix(line, "Platform:"))
		}

		// System Name
		if strings.HasPrefix(line, "System Name:") {
			currentNeighbor.SystemName = strings.TrimSpace(strings.TrimPrefix(line, "System Name:"))
		}
	}

	if currentNeighbor != nil && currentPort != "" {
		p.attachNeighborToPort(device, currentPort, currentNeighbor)
	}
}

func (p *CiscoParser) attachNeighborToPort(device *models.Device, portName string, neighbor *models.Neighbor) {
	for i := range device.Ports {
		if device.Ports[i].Name == portName {
			device.Ports[i].Neighbor = neighbor
			break
		}
	}
}
