package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kwld/network-tools/internal/mapper"
	"github.com/kwld/network-tools/internal/scanner"
	"github.com/kwld/network-tools/internal/visualizer"
	"github.com/kwld/network-tools/pkg/models"
)

func main() {
	log.Println("Network Topology Scanner Starting...")

	// Load configuration from environment
	config := loadConfig()

	// Read device IPs from file
	ips, err := readDeviceIPs(config.DevicesFile)
	if err != nil {
		log.Fatalf("Failed to read devices file: %v", err)
	}

	if len(ips) == 0 {
		log.Fatal("No devices to scan")
	}

	log.Printf("Found %d devices to scan", len(ips))

	// Create scanner
	scannerConfig := scanner.Config{
		SSHUsername:   config.SSHUsername,
		SSHPassword:   config.SSHPassword,
		SSHKeyPath:    config.SSHKeyPath,
		SNMPCommunity: config.SNMPCommunity,
		SNMPVersion:   config.SNMPVersion,
		SNMPUsername:  config.SNMPUsername,
		SNMPAuthPass:  config.SNMPAuthPass,
		SNMPPrivPass:  config.SNMPPrivPass,
		Timeout:       time.Duration(config.ScanTimeout) * time.Second,
		MaxWorkers:    config.ConcurrentScans,
	}

	s := scanner.NewScanner(scannerConfig)

	// Scan devices
	log.Println("Scanning devices...")
	devices := s.ScanDevices(ips)
	log.Printf("Scan complete. Processed %d devices", len(devices))

	// Build topology
	log.Println("Building topology...")
	builder := mapper.NewTopologyBuilder()
	topology := builder.BuildTopology(devices)
	topology.ScannedAt = time.Now()

	// Print summary
	printSummary(topology)

	// Ensure output directory exists
	exporter := visualizer.NewExporter(config.OutputDir)
	if err := exporter.EnsureOutputDir(); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Export JSON
	log.Println("Exporting topology data...")
	jsonPath, err := exporter.ExportJSON(topology)
	if err != nil {
		log.Printf("Failed to export JSON: %v", err)
	} else {
		log.Printf("Exported JSON: %s", jsonPath)
	}

	// Export summary
	summaryPath, err := exporter.ExportSummary(topology)
	if err != nil {
		log.Printf("Failed to export summary: %v", err)
	} else {
		log.Printf("Exported summary: %s", summaryPath)
	}

	// Generate visual diagram
	log.Println("Generating network diagram...")
	graph := mapper.BuildGraph(topology)
	
	generator := visualizer.NewGraphVizGenerator(config.DiagramLayout, config.DiagramFormat)
	
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	diagramPath := fmt.Sprintf("%s/topology_%s.%s", config.OutputDir, timestamp, config.DiagramFormat)
	
	if err := generator.Generate(graph, diagramPath); err != nil {
		log.Printf("Failed to generate diagram: %v", err)
		log.Println("Note: GraphViz must be installed for diagram generation")
	} else {
		log.Printf("Generated diagram: %s", diagramPath)
	}

	// Also save DOT file
	dotPath := fmt.Sprintf("%s/topology_%s.dot", config.OutputDir, timestamp)
	if err := generator.GenerateDOTFile(graph, dotPath); err != nil {
		log.Printf("Failed to save DOT file: %v", err)
	} else {
		log.Printf("Saved DOT file: %s", dotPath)
	}

	log.Println("Scanner complete!")
}

// Config holds application configuration
type Config struct {
	SSHUsername     string
	SSHPassword     string
	SSHKeyPath      string
	SNMPCommunity   string
	SNMPVersion     string
	SNMPUsername    string
	SNMPAuthPass    string
	SNMPPrivPass    string
	DevicesFile     string
	OutputDir       string
	ScanTimeout     int
	ConcurrentScans int
	DiagramFormat   string
	DiagramLayout   string
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	return Config{
		SSHUsername:     getEnv("NETWORK_SSH_USERNAME", "admin"),
		SSHPassword:     getEnv("NETWORK_SSH_PASSWORD", ""),
		SSHKeyPath:      getEnv("NETWORK_SSH_KEY_PATH", ""),
		SNMPCommunity:   getEnv("NETWORK_SNMP_COMMUNITY", "public"),
		SNMPVersion:     getEnv("NETWORK_SNMP_VERSION", "2c"),
		SNMPUsername:    getEnv("NETWORK_SNMP_V3_USER", ""),
		SNMPAuthPass:    getEnv("NETWORK_SNMP_V3_AUTH_PASS", ""),
		SNMPPrivPass:    getEnv("NETWORK_SNMP_V3_PRIV_PASS", ""),
		DevicesFile:     getEnv("NETWORK_DEVICES_FILE", "/config/devices.txt"),
		OutputDir:       getEnv("NETWORK_OUTPUT_DIR", "/output"),
		ScanTimeout:     getEnvInt("NETWORK_SCAN_TIMEOUT", 30),
		ConcurrentScans: getEnvInt("NETWORK_CONCURRENT_SCANS", 5),
		DiagramFormat:   getEnv("NETWORK_DIAGRAM_FORMAT", "svg"),
		DiagramLayout:   getEnv("NETWORK_DIAGRAM_LAYOUT", "dot"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// readDeviceIPs reads IP addresses from a file
func readDeviceIPs(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ips = append(ips, line)
	}

	return ips, scanner.Err()
}

// printSummary prints scan results summary
func printSummary(topology *models.Topology) {
	fmt.Println("\n========== Scan Summary ==========")
	fmt.Printf("Total Devices:    %d\n", topology.Summary.TotalDevices)
	fmt.Printf("Success:          %d\n", topology.Summary.SuccessDevices)
	fmt.Printf("Failed:           %d\n", topology.Summary.FailedDevices)
	fmt.Printf("Total Ports:      %d\n", topology.Summary.TotalPorts)
	fmt.Printf("Total Connections: %d\n", topology.Summary.TotalConnections)
	
	if len(topology.Summary.FailedIPs) > 0 {
		fmt.Println("\nFailed IPs:")
		for _, ip := range topology.Summary.FailedIPs {
			fmt.Printf("  - %s\n", ip)
		}
	}
	fmt.Println("==================================\n")
}
