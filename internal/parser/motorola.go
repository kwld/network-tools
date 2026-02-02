package parser

import (
	"regexp"
	"strings"

	"github.com/kwld/network-tools/pkg/models"
)

// MotorolaParser handles Motorola devices
type MotorolaParser struct{}

func (p *MotorolaParser) GetVendorName() string {
	return "motorola"
}

func (p *MotorolaParser) GetCommands() []string {
	return []string{
		"show system",
		"show interfaces",
		"show lldp neighbors detail",
	}
}

func (p *MotorolaParser) Parse(outputs map[string]string, device *models.Device) error {
	device.Vendor = "motorola"

	// Parse system info
	if output, ok := outputs["show system"]; ok {
		p.parseSystem(output, device)
	}

	// Parse interfaces
	if output, ok := outputs["show interfaces"]; ok {
		device.Ports = p.parseInterfaces(output)
	}

	// Parse LLDP neighbors
	if output, ok := outputs["show lldp neighbors detail"]; ok {
		p.parseNeighbors(output, device)
	}

	return nil
}

func (p *MotorolaParser) parseSystem(output string, device *models.Device) {
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if match := regexp.MustCompile(`(?i)hostname[:\s]+(.+)`).FindStringSubmatch(line); len(match) > 1 {
			device.Hostname = strings.TrimSpace(match[1])
		}
		
		if match := regexp.MustCompile(`(?i)model[:\s]+(.+)`).FindStringSubmatch(line); len(match) > 1 {
			device.Model = strings.TrimSpace(match[1])
		}
		
		if match := regexp.MustCompile(`(?i)version[:\s]+(.+)`).FindStringSubmatch(line); len(match) > 1 {
			device.Version = strings.TrimSpace(match[1])
		}
	}
}

func (p *MotorolaParser) parseInterfaces(output string) []models.Port {
	var ports []models.Port
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
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

		// Parse status
		for _, field := range fields {
			status := strings.ToLower(field)
			if status == "up" {
				port.Status = "up"
			} else if status == "down" {
				port.Status = "down"
			}
		}

		ports = append(ports, port)
	}

	return ports
}

func (p *MotorolaParser) parseNeighbors(output string, device *models.Device) {
	lines := strings.Split(output, "\n")
	var currentNeighbor *models.Neighbor
	var currentPort string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Local Interface
		if strings.HasPrefix(line, "Local Interface:") {
			if currentNeighbor != nil && currentPort != "" {
				p.attachNeighborToPort(device, currentPort, currentNeighbor)
			}
			currentNeighbor = &models.Neighbor{Protocol: "LLDP"}
			currentPort = strings.TrimSpace(strings.TrimPrefix(line, "Local Interface:"))
		}

		if currentNeighbor == nil {
			continue
		}

		// System Name
		if strings.HasPrefix(line, "System Name:") {
			name := strings.TrimSpace(strings.TrimPrefix(line, "System Name:"))
			currentNeighbor.SystemName = name
			currentNeighbor.DeviceID = name
		}

		// Port ID
		if strings.HasPrefix(line, "Port ID:") {
			currentNeighbor.PortID = strings.TrimSpace(strings.TrimPrefix(line, "Port ID:"))
		}

		// Management Address
		if strings.Contains(line, "Management Address:") {
			if match := regexp.MustCompile(`([0-9.]+)`).FindStringSubmatch(line); len(match) > 0 {
				currentNeighbor.IPAddress = match[0]
			}
		}
	}

	if currentNeighbor != nil && currentPort != "" {
		p.attachNeighborToPort(device, currentPort, currentNeighbor)
	}
}

func (p *MotorolaParser) attachNeighborToPort(device *models.Device, portName string, neighbor *models.Neighbor) {
	for i := range device.Ports {
		if device.Ports[i].Name == portName {
			device.Ports[i].Neighbor = neighbor
			break
		}
	}
}
