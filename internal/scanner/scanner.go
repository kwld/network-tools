package scanner

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kwld/network-tools/internal/parser"
	"github.com/kwld/network-tools/pkg/models"
)

// Scanner manages the network device scanning process
type Scanner struct {
	sshConfig  SSHConfig
	snmpConfig SNMPConfig
	timeout    time.Duration
	maxWorkers int
}

// Config holds scanner configuration
type Config struct {
	SSHUsername    string
	SSHPassword    string
	SSHKeyPath     string
	SNMPCommunity  string
	SNMPVersion    string
	SNMPUsername   string
	SNMPAuthPass   string
	SNMPPrivPass   string
	Timeout        time.Duration
	MaxWorkers     int
}

// NewScanner creates a new scanner instance
func NewScanner(config Config) *Scanner {
	return &Scanner{
		sshConfig: SSHConfig{
			Port:     22,
			Username: config.SSHUsername,
			Password: config.SSHPassword,
			KeyPath:  config.SSHKeyPath,
			Timeout:  config.Timeout,
		},
		snmpConfig: SNMPConfig{
			Port:      161,
			Version:   config.SNMPVersion,
			Community: config.SNMPCommunity,
			Username:  config.SNMPUsername,
			AuthPass:  config.SNMPAuthPass,
			PrivPass:  config.SNMPPrivPass,
			Timeout:   config.Timeout,
		},
		timeout:    config.Timeout,
		maxWorkers: config.MaxWorkers,
	}
}

// ScanDevices scans multiple devices concurrently
func (s *Scanner) ScanDevices(ips []string) []models.Device {
	var wg sync.WaitGroup
	deviceChan := make(chan models.Device, len(ips))
	semaphore := make(chan struct{}, s.maxWorkers)

	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			device := s.ScanDevice(ip)
			deviceChan <- device
		}(ip)
	}

	wg.Wait()
	close(deviceChan)

	var devices []models.Device
	for device := range deviceChan {
		devices = append(devices, device)
	}

	return devices
}

// ScanDevice scans a single device
func (s *Scanner) ScanDevice(ip string) models.Device {
	device := models.Device{
		IP:          ip,
		LastScanned: time.Now(),
		ScanStatus:  "failed",
	}

	log.Printf("Scanning device: %s", ip)

	// Try SSH first
	sshConfig := s.sshConfig
	sshConfig.Host = ip
	
	sshClient, err := NewSSHClient(sshConfig)
	if err == nil {
		defer sshClient.Close()
		
		vendor, _ := sshClient.DetectVendor()
		device.Vendor = vendor
		
		if err := s.scanViaSSH(sshClient, &device); err == nil {
			device.ScanStatus = "success"
			log.Printf("Successfully scanned %s via SSH (vendor: %s)", ip, vendor)
			return device
		}
		log.Printf("SSH scan partial for %s: %v", ip, err)
	}

	// Fallback to SNMP
	snmpConfig := s.snmpConfig
	snmpConfig.Host = ip
	
	snmpClient, err := NewSNMPClient(snmpConfig)
	if err == nil {
		defer snmpClient.Close()
		
		if err := s.scanViaSNMP(snmpClient, &device); err == nil {
			device.ScanStatus = "partial"
			log.Printf("Scanned %s via SNMP (partial data)", ip)
			return device
		}
		log.Printf("SNMP scan failed for %s: %v", ip, err)
	}

	device.Error = "Failed to connect via SSH or SNMP"
	log.Printf("Failed to scan %s", ip)
	return device
}

// scanViaSSH scans device using SSH
func (s *Scanner) scanViaSSH(client *SSHClient, device *models.Device) error {
	// Get appropriate parser based on vendor
	p := parser.GetParser(device.Vendor)
	if p == nil {
		return fmt.Errorf("no parser available for vendor: %s", device.Vendor)
	}

	// Get commands to execute
	commands := p.GetCommands()
	
	// Execute commands
	outputs, err := client.ExecuteCommands(commands)
	if err != nil {
		return err
	}

	// Parse outputs
	return p.Parse(outputs, device)
}

// scanViaSNMP scans device using SNMP
func (s *Scanner) scanViaSNMP(client *SNMPClient, device *models.Device) error {
	// Get system info
	sysInfo, err := client.GetSystemInfo()
	if err != nil {
		return err
	}

	if sysDesc, ok := sysInfo["sysDescr"]; ok {
		device.Version = sysDesc
	}
	if sysName, ok := sysInfo["sysName"]; ok {
		device.Hostname = sysName
	}

	// Get interface info
	interfaces, err := client.GetInterfaceInfo()
	if err == nil {
		for _, iface := range interfaces {
			if name, ok := iface["name"].(string); ok {
				port := models.Port{
					Name:   name,
					Status: "unknown",
				}
				device.Ports = append(device.Ports, port)
			}
		}
	}

	return nil
}
