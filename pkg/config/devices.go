package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DeviceConfig represents configuration for a single device
type DeviceConfig struct {
	IP   string      `yaml:"ip"`
	SSH  *SSHConfig  `yaml:"ssh,omitempty"`
	SNMP *SNMPConfig `yaml:"snmp,omitempty"`
}

// SSHConfig holds SSH credentials for a device
type SSHConfig struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	KeyPath  string `yaml:"key_path,omitempty"`
}

// SNMPConfig holds SNMP credentials for a device
type SNMPConfig struct {
	Community    string `yaml:"community,omitempty"`
	Version      string `yaml:"version,omitempty"`
	Username     string `yaml:"username,omitempty"`
	AuthPassword string `yaml:"auth_password,omitempty"`
	PrivPassword string `yaml:"priv_password,omitempty"`
}

// DevicesConfig represents the YAML configuration file structure
type DevicesConfig struct {
	Devices []DeviceConfig `yaml:"devices"`
}

// LoadDevicesConfig loads device configurations from a YAML file
func LoadDevicesConfig(filepath string) (*DevicesConfig, error) {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", filepath)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config DevicesConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return &config, nil
}

// GetDeviceConfig returns the configuration for a specific IP address
func (dc *DevicesConfig) GetDeviceConfig(ip string) *DeviceConfig {
	for i := range dc.Devices {
		if dc.Devices[i].IP == ip {
			return &dc.Devices[i]
		}
	}
	return nil
}

// GetAllIPs returns a list of all IP addresses in the configuration
func (dc *DevicesConfig) GetAllIPs() []string {
	ips := make([]string, 0, len(dc.Devices))
	for _, device := range dc.Devices {
		ips = append(ips, device.IP)
	}
	return ips
}
