package mapper

import (
	"strings"

	"github.com/kwld/network-tools/pkg/models"
)

// TopologyBuilder builds network topology from scanned devices
type TopologyBuilder struct{}

// NewTopologyBuilder creates a new topology builder
func NewTopologyBuilder() *TopologyBuilder {
	return &TopologyBuilder{}
}

// BuildTopology creates a topology from device scan results
func (tb *TopologyBuilder) BuildTopology(devices []models.Device) *models.Topology {
	topology := &models.Topology{
		Devices:     devices,
		Connections: tb.extractConnections(devices),
	}
	
	topology.Summary = tb.calculateSummary(topology)
	
	return topology
}

// extractConnections identifies connections between devices
func (tb *TopologyBuilder) extractConnections(devices []models.Device) []models.Connection {
	var connections []models.Connection
	seen := make(map[string]bool)

	// Build device lookup map
	deviceLookup := make(map[string]*models.Device)
	for i := range devices {
		deviceLookup[devices[i].IP] = &devices[i]
		if devices[i].Hostname != "" {
			deviceLookup[strings.ToLower(devices[i].Hostname)] = &devices[i]
		}
	}

	// Iterate through devices and their ports
	for _, device := range devices {
		for _, port := range device.Ports {
			if port.Neighbor == nil {
				continue
			}

			// Try to find the neighbor device
			targetDevice := tb.findTargetDevice(port.Neighbor, deviceLookup)
			if targetDevice == nil {
				continue
			}

			// Create connection
			conn := models.Connection{
				SourceDevice: device.Hostname,
				SourcePort:   port.Name,
				TargetDevice: targetDevice.Hostname,
				TargetPort:   port.Neighbor.PortID,
				Protocol:     port.Neighbor.Protocol,
			}

			// Avoid duplicate connections
			connKey := tb.connectionKey(conn)
			reverseKey := tb.reverseConnectionKey(conn)
			
			if !seen[connKey] && !seen[reverseKey] {
				connections = append(connections, conn)
				seen[connKey] = true
			}
		}
	}

	return connections
}

// findTargetDevice finds a device by neighbor information
func (tb *TopologyBuilder) findTargetDevice(neighbor *models.Neighbor, deviceLookup map[string]*models.Device) *models.Device {
	// Try by IP address
	if neighbor.IPAddress != "" {
		if device, ok := deviceLookup[neighbor.IPAddress]; ok {
			return device
		}
	}

	// Try by system name
	if neighbor.SystemName != "" {
		if device, ok := deviceLookup[strings.ToLower(neighbor.SystemName)]; ok {
			return device
		}
	}

	// Try by device ID
	if neighbor.DeviceID != "" {
		if device, ok := deviceLookup[strings.ToLower(neighbor.DeviceID)]; ok {
			return device
		}
	}

	return nil
}

// connectionKey creates a unique key for a connection
func (tb *TopologyBuilder) connectionKey(conn models.Connection) string {
	return conn.SourceDevice + ":" + conn.SourcePort + "->" + conn.TargetDevice + ":" + conn.TargetPort
}

// reverseConnectionKey creates a key for the reverse direction
func (tb *TopologyBuilder) reverseConnectionKey(conn models.Connection) string {
	return conn.TargetDevice + ":" + conn.TargetPort + "->" + conn.SourceDevice + ":" + conn.SourcePort
}

// calculateSummary generates summary statistics
func (tb *TopologyBuilder) calculateSummary(topology *models.Topology) models.Summary {
	summary := models.Summary{
		TotalDevices:     len(topology.Devices),
		TotalConnections: len(topology.Connections),
	}

	for _, device := range topology.Devices {
		if device.ScanStatus == "success" || device.ScanStatus == "partial" {
			summary.SuccessDevices++
		} else {
			summary.FailedDevices++
			summary.FailedIPs = append(summary.FailedIPs, device.IP)
		}
		summary.TotalPorts += len(device.Ports)
	}

	return summary
}
