package models

import "time"

// Device represents a network device (switch, router, etc.)
type Device struct {
	IP          string            `json:"ip"`
	Hostname    string            `json:"hostname"`
	Vendor      string            `json:"vendor"`
	Model       string            `json:"model"`
	Version     string            `json:"version"`
	Ports       []Port            `json:"ports"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	LastScanned time.Time         `json:"last_scanned"`
	ScanStatus  string            `json:"scan_status"` // success, failed, partial
	Error       string            `json:"error,omitempty"`
}

// Port represents a network port on a device
type Port struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`           // up, down, disabled
	Speed       string `json:"speed,omitempty"`  // 1G, 10G, etc.
	MACAddress  string `json:"mac_address,omitempty"`
	Neighbor    *Neighbor `json:"neighbor,omitempty"`
}

// Neighbor represents a discovered neighbor via LLDP/CDP
type Neighbor struct {
	DeviceID    string `json:"device_id"`     // Hostname or device identifier
	PortID      string `json:"port_id"`       // Remote port identifier
	SystemName  string `json:"system_name,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
	Platform    string `json:"platform,omitempty"`
	Capabilities string `json:"capabilities,omitempty"`
	Protocol    string `json:"protocol"` // LLDP or CDP
}

// Connection represents a link between two devices
type Connection struct {
	SourceDevice string `json:"source_device"`
	SourcePort   string `json:"source_port"`
	TargetDevice string `json:"target_device"`
	TargetPort   string `json:"target_port"`
	Protocol     string `json:"protocol"` // LLDP or CDP
}

// Topology represents the complete network topology
type Topology struct {
	Devices     []Device     `json:"devices"`
	Connections []Connection `json:"connections"`
	ScannedAt   time.Time    `json:"scanned_at"`
	Summary     Summary      `json:"summary"`
}

// Summary provides statistics about the scan
type Summary struct {
	TotalDevices   int      `json:"total_devices"`
	SuccessDevices int      `json:"success_devices"`
	FailedDevices  int      `json:"failed_devices"`
	TotalPorts     int      `json:"total_ports"`
	TotalConnections int    `json:"total_connections"`
	FailedIPs      []string `json:"failed_ips,omitempty"`
}
