package visualizer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kwld/network-tools/pkg/models"
)

// Exporter handles exporting topology data to various formats
type Exporter struct {
	outputDir string
}

// NewExporter creates a new exporter
func NewExporter(outputDir string) *Exporter {
	return &Exporter{outputDir: outputDir}
}

// ExportJSON exports topology to JSON format
func (e *Exporter) ExportJSON(topology *models.Topology) (string, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("topology_%s.json", timestamp)
	filePath := filepath.Join(e.outputDir, filename)

	data, err := json.MarshalIndent(topology, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write JSON file: %w", err)
	}

	return filePath, nil
}

// ExportSummary exports a text summary of the scan
func (e *Exporter) ExportSummary(topology *models.Topology) (string, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("summary_%s.txt", timestamp)
	filePath := filepath.Join(e.outputDir, filename)

	summary := fmt.Sprintf(`Network Topology Scan Summary
=============================
Scan Time: %s

Devices:
  Total: %d
  Success: %d
  Failed: %d
  
Ports: %d
Connections: %d

Failed IPs:
`,
		topology.ScannedAt.Format(time.RFC3339),
		topology.Summary.TotalDevices,
		topology.Summary.SuccessDevices,
		topology.Summary.FailedDevices,
		topology.Summary.TotalPorts,
		topology.Summary.TotalConnections,
	)

	for _, ip := range topology.Summary.FailedIPs {
		summary += fmt.Sprintf("  - %s\n", ip)
	}

	err := os.WriteFile(filePath, []byte(summary), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write summary file: %w", err)
	}

	return filePath, nil
}

// EnsureOutputDir ensures the output directory exists
func (e *Exporter) EnsureOutputDir() error {
	return os.MkdirAll(e.outputDir, 0755)
}
