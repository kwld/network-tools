package parser

import (
	"regexp"
	"strings"

	"github.com/kwld/network-tools/pkg/models"
)

// Pre-compiled regex patterns for performance
var (
	hostnameRegex = regexp.MustCompile(`(?i)hostname[:\s]+([^\s]+)`)
	modelRegex    = regexp.MustCompile(`(?i)model[:\s]+(.+)`)
	versionRegex  = regexp.MustCompile(`(?i)version[:\s]+([^\s,]+)`)
)

// GenericParser handles generic devices with LLDP support
type GenericParser struct{}

func (p *GenericParser) GetVendorName() string {
	return "generic"
}

func (p *GenericParser) GetCommands() []string {
	return []string{
		"show system",
		"show version",
		"show interfaces",
		"show lldp neighbors",
	}
}

func (p *GenericParser) Parse(outputs map[string]string, device *models.Device) error {
	device.Vendor = "generic"
	
	// Try to extract basic info from common commands
	if output, ok := outputs["show system"]; ok {
		p.parseSystemInfo(output, device)
	}
	
	if output, ok := outputs["show version"]; ok {
		p.parseVersionInfo(output, device)
	}
	
	return nil
}

// parseSystemInfo tries to extract hostname and model from system output
func (p *GenericParser) parseSystemInfo(output string, device *models.Device) {
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Try to find hostname
		if match := hostnameRegex.FindStringSubmatch(line); len(match) > 1 {
			device.Hostname = strings.TrimSpace(match[1])
		}
		
		// Try to find model
		if match := modelRegex.FindStringSubmatch(line); len(match) > 1 {
			device.Model = strings.TrimSpace(match[1])
		}
	}
}

// parseVersionInfo tries to extract version info
func (p *GenericParser) parseVersionInfo(output string, device *models.Device) {
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Try to find version
		if match := versionRegex.FindStringSubmatch(line); len(match) > 1 {
			device.Version = strings.TrimSpace(match[1])
		}
		
		// Try to find model if not already set
		if device.Model == "" {
			if match := modelRegex.FindStringSubmatch(line); len(match) > 1 {
				device.Model = strings.TrimSpace(match[1])
			}
		}
	}
}
