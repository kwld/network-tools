package config

import (
	"os"
	"testing"
)

func TestLoadDevicesConfig(t *testing.T) {
	// Create a temporary YAML file
	yamlContent := `devices:
  - ip: 192.168.1.1
    ssh:
      username: admin
      password: password123
    snmp:
      community: public
      version: 2c
  - ip: 192.168.1.2
    ssh:
      username: testuser
      key_path: /path/to/key
`
	tmpFile, err := os.CreateTemp("", "devices-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the config
	config, err := LoadDevicesConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Devices) != 2 {
		t.Errorf("Expected 2 devices, got %d", len(config.Devices))
	}

	// Test first device
	device1 := config.Devices[0]
	if device1.IP != "192.168.1.1" {
		t.Errorf("Expected IP 192.168.1.1, got %s", device1.IP)
	}
	if device1.SSH == nil {
		t.Fatal("Expected SSH config, got nil")
	}
	if device1.SSH.Username != "admin" {
		t.Errorf("Expected SSH username admin, got %s", device1.SSH.Username)
	}
	if device1.SSH.Password != "password123" {
		t.Errorf("Expected SSH password password123, got %s", device1.SSH.Password)
	}

	// Test GetDeviceConfig
	foundDevice := config.GetDeviceConfig("192.168.1.1")
	if foundDevice == nil {
		t.Error("Expected to find device 192.168.1.1")
	}
	if foundDevice != nil && foundDevice.IP != "192.168.1.1" {
		t.Errorf("Expected IP 192.168.1.1, got %s", foundDevice.IP)
	}

	// Test non-existent device
	notFound := config.GetDeviceConfig("192.168.1.99")
	if notFound != nil {
		t.Error("Expected nil for non-existent device")
	}

	// Test GetAllIPs
	ips := config.GetAllIPs()
	if len(ips) != 2 {
		t.Errorf("Expected 2 IPs, got %d", len(ips))
	}
}

func TestLoadDevicesConfigNonExistent(t *testing.T) {
	_, err := LoadDevicesConfig("/non/existent/file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadDevicesConfigInvalidYAML(t *testing.T) {
	// Create a temporary file with invalid YAML
	invalidYAML := `this is not: valid: yaml: content`
	tmpFile, err := os.CreateTemp("", "invalid-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(invalidYAML)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	_, err = LoadDevicesConfig(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}
