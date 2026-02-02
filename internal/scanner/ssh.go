package scanner

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHClient represents an SSH connection to a network device
type SSHClient struct {
	client  *ssh.Client
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
	stderr  io.Reader
}

// SSHConfig holds SSH connection configuration
type SSHConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	KeyPath  string
	Timeout  time.Duration
}

// NewSSHClient creates a new SSH client connection
func NewSSHClient(config SSHConfig) (*SSHClient, error) {
	var authMethods []ssh.AuthMethod

	// Try key-based auth first if key path is provided
	if config.KeyPath != "" {
		key, err := os.ReadFile(config.KeyPath)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
	}

	// Add password auth
	if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication methods available")
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: In production, verify host keys
		Timeout:         config.Timeout,
	}

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &SSHClient{client: client}, nil
}

// ExecuteCommand executes a command on the remote device
func (c *SSHClient) ExecuteCommand(command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("command execution failed: %w", err)
	}

	return string(output), nil
}

// ExecuteCommands executes multiple commands sequentially
func (c *SSHClient) ExecuteCommands(commands []string) (map[string]string, error) {
	results := make(map[string]string)

	for _, cmd := range commands {
		output, err := c.ExecuteCommand(cmd)
		if err != nil {
			// Continue with other commands even if one fails
			results[cmd] = fmt.Sprintf("ERROR: %v", err)
		} else {
			results[cmd] = output
		}
	}

	return results, nil
}

// DetectVendor attempts to detect the device vendor from SSH banner or commands
func (c *SSHClient) DetectVendor() (string, error) {
	// Try to get system information
	commands := []string{
		"/system resource print",  // Mikrotik
		"show version",             // Cisco, Motorola
		"show system",              // Generic
	}

	for _, cmd := range commands {
		output, err := c.ExecuteCommand(cmd)
		if err == nil && output != "" {
			output = strings.ToLower(output)
			if strings.Contains(output, "mikrotik") || strings.Contains(output, "routeros") {
				return "mikrotik", nil
			}
			if strings.Contains(output, "cisco") {
				return "cisco", nil
			}
			if strings.Contains(output, "motorola") {
				return "motorola", nil
			}
		}
	}

	return "generic", nil
}

// Close closes the SSH connection
func (c *SSHClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
