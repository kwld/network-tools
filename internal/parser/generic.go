package parser

import (
	"github.com/kwld/network-tools/pkg/models"
)

// GenericParser handles generic devices with LLDP support
type GenericParser struct{}

func (p *GenericParser) GetVendorName() string {
	return "generic"
}

func (p *GenericParser) GetCommands() []string {
	return []string{
		"show system",
		"show interfaces",
		"show lldp neighbors",
	}
}

func (p *GenericParser) Parse(outputs map[string]string, device *models.Device) error {
	device.Vendor = "generic"
	
	// Generic parser provides minimal information
	// Most data will come from SNMP fallback
	
	return nil
}
