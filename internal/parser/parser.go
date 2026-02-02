package parser

import "github.com/kwld/network-tools/pkg/models"

// Parser defines the interface for vendor-specific parsers
type Parser interface {
	// GetCommands returns the list of commands to execute on the device
	GetCommands() []string
	
	// Parse parses the command outputs and populates the device struct
	Parse(outputs map[string]string, device *models.Device) error
	
	// GetVendorName returns the vendor name this parser handles
	GetVendorName() string
}

// GetParser returns the appropriate parser for a given vendor
func GetParser(vendor string) Parser {
	switch vendor {
	case "mikrotik":
		return &MikrotikParser{}
	case "cisco":
		return &CiscoParser{}
	case "motorola":
		return &MotorolaParser{}
	default:
		return &GenericParser{}
	}
}
