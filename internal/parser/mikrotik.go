package parser

import (
	"regexp"
	"strings"

	"github.com/kwld/network-tools/pkg/models"
)

// MikrotikParser handles Mikrotik RouterOS devices
type MikrotikParser struct{}

func (p *MikrotikParser) GetVendorName() string {
	return "mikrotik"
}

func (p *MikrotikParser) GetCommands() []string {
	return []string{
		"/system identity print",
		"/system resource print",
		"/interface print detail",
		"/interface ethernet print detail",
		"/ip neighbor print detail",
	}
}

func (p *MikrotikParser) Parse(outputs map[string]string, device *models.Device) error {
	// Parse system identity
	if output, ok := outputs["/system identity print"]; ok {
		if match := regexp.MustCompile(`name:\s*(.+)`).FindStringSubmatch(output); len(match) > 1 {
			device.Hostname = strings.TrimSpace(match[1])
		}
	}

	// Parse system resource for model/version
	if output, ok := outputs["/system resource print"]; ok {
		if match := regexp.MustCompile(`board-name:\s*(.+)`).FindStringSubmatch(output); len(match) > 1 {
			device.Model = strings.TrimSpace(match[1])
		}
		if match := regexp.MustCompile(`version:\s*(.+)`).FindStringSubmatch(output); len(match) > 1 {
			device.Version = strings.TrimSpace(match[1])
		}
	}

	device.Vendor = "mikrotik"

	// Parse interfaces
	if output, ok := outputs["/interface print detail"]; ok {
		device.Ports = p.parseInterfaces(output)
	}

	// Parse neighbors (LLDP/discovery)
	if output, ok := outputs["/ip neighbor print detail"]; ok {
		p.parseNeighbors(output, device)
	}

	return nil
}

func (p *MikrotikParser) parseInterfaces(output string) []models.Port {
	var ports []models.Port
	
	// Split by interface entries (each starts with a number and 'R' or 'X')
	lines := strings.Split(output, "\n")
	var currentPort *models.Port
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// New interface entry
		if regexp.MustCompile(`^\d+\s+[RX]`).MatchString(line) {
			if currentPort != nil {
				ports = append(ports, *currentPort)
			}
			currentPort = &models.Port{Status: "unknown"}
			
			// Extract name
			if match := regexp.MustCompile(`name="([^"]+)"`).FindStringSubmatch(line); len(match) > 1 {
				currentPort.Name = match[1]
			}
		}

		if currentPort == nil {
			continue
		}

		// Parse status
		if strings.Contains(line, "running=yes") || strings.Contains(line, "disabled=no") {
			currentPort.Status = "up"
		} else if strings.Contains(line, "disabled=yes") {
			currentPort.Status = "disabled"
		} else if strings.Contains(line, "running=no") {
			currentPort.Status = "down"
		}

		// Parse speed
		if match := regexp.MustCompile(`rate="([^"]+)"`).FindStringSubmatch(line); len(match) > 1 {
			currentPort.Speed = match[1]
		}

		// Parse MAC
		if match := regexp.MustCompile(`mac-address=([0-9A-Fa-f:]+)`).FindStringSubmatch(line); len(match) > 1 {
			currentPort.MACAddress = match[1]
		}
	}

	if currentPort != nil {
		ports = append(ports, *currentPort)
	}

	return ports
}

func (p *MikrotikParser) parseNeighbors(output string, device *models.Device) {
	lines := strings.Split(output, "\n")
	var currentNeighbor *models.Neighbor
	var currentInterface string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// New neighbor entry
		if regexp.MustCompile(`^\d+`).MatchString(line) {
			if currentNeighbor != nil && currentInterface != "" {
				p.attachNeighborToPort(device, currentInterface, currentNeighbor)
			}
			currentNeighbor = &models.Neighbor{Protocol: "LLDP"}
			currentInterface = ""
		}

		if currentNeighbor == nil {
			continue
		}

		// Parse interface
		if match := regexp.MustCompile(`interface=([^\s]+)`).FindStringSubmatch(line); len(match) > 1 {
			currentInterface = match[1]
		}

		// Parse identity/system-name
		if match := regexp.MustCompile(`identity="([^"]+)"`).FindStringSubmatch(line); len(match) > 1 {
			currentNeighbor.DeviceID = match[1]
			currentNeighbor.SystemName = match[1]
		}

		// Parse interface-name (remote port)
		if match := regexp.MustCompile(`interface-name="([^"]+)"`).FindStringSubmatch(line); len(match) > 1 {
			currentNeighbor.PortID = match[1]
		}

		// Parse address
		if match := regexp.MustCompile(`address=([0-9.]+)`).FindStringSubmatch(line); len(match) > 1 {
			currentNeighbor.IPAddress = match[1]
		}
	}

	if currentNeighbor != nil && currentInterface != "" {
		p.attachNeighborToPort(device, currentInterface, currentNeighbor)
	}
}

func (p *MikrotikParser) attachNeighborToPort(device *models.Device, interfaceName string, neighbor *models.Neighbor) {
	for i := range device.Ports {
		if device.Ports[i].Name == interfaceName {
			device.Ports[i].Neighbor = neighbor
			break
		}
	}
}
