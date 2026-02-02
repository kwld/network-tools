package visualizer

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/kwld/network-tools/internal/mapper"
)

// GraphVizGenerator generates network diagrams using GraphViz
type GraphVizGenerator struct {
	layout string
	format string
}

// NewGraphVizGenerator creates a new GraphViz generator
func NewGraphVizGenerator(layout, format string) *GraphVizGenerator {
	if layout == "" {
		layout = "dot"
	}
	if format == "" {
		format = "svg"
	}
	return &GraphVizGenerator{
		layout: layout,
		format: format,
	}
}

// Generate creates a network diagram from graph data
func (g *GraphVizGenerator) Generate(graph *mapper.Graph, outputPath string) error {
	// Generate DOT source
	dotSource := g.generateDOT(graph)

	// Write to file and render
	return g.renderDOT(dotSource, outputPath)
}

// generateDOT creates DOT format source code
func (g *GraphVizGenerator) generateDOT(graph *mapper.Graph) string {
	var sb strings.Builder

	sb.WriteString("graph NetworkTopology {\n")
	sb.WriteString("  layout=" + g.layout + ";\n")
	sb.WriteString("  node [shape=box, style=filled, fontname=\"Arial\"];\n")
	sb.WriteString("  edge [fontname=\"Arial\", fontsize=10];\n\n")

	// Add nodes
	for _, node := range graph.Nodes {
		color := g.getNodeColor(node.Status)
		label := g.escapeLabel(node.Label)
		
		if node.Vendor != "" || node.Model != "" {
			label = fmt.Sprintf("%s\\n[%s %s]", label, node.Vendor, node.Model)
		}

		sb.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", fillcolor=\"%s\"];\n",
			g.escapeID(node.ID), label, color))
	}

	sb.WriteString("\n")

	// Add edges
	for _, edge := range graph.Edges {
		sb.WriteString(fmt.Sprintf("  \"%s\" -- \"%s\" [label=\"%s\", taillabel=\"%s\", headlabel=\"%s\"];\n",
			g.escapeID(edge.Source),
			g.escapeID(edge.Target),
			edge.Protocol,
			edge.SourcePort,
			edge.TargetPort))
	}

	sb.WriteString("}\n")

	return sb.String()
}

// renderDOT executes GraphViz to render the diagram
func (g *GraphVizGenerator) renderDOT(dotSource, outputPath string) error {
	// Check if dot command is available
	_, err := exec.LookPath("dot")
	if err != nil {
		return fmt.Errorf("graphviz 'dot' command not found: %w", err)
	}

	// Execute dot command
	cmd := exec.Command("dot", "-T"+g.format, "-o", outputPath)
	cmd.Stdin = strings.NewReader(dotSource)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to render diagram: %w, output: %s", err, string(output))
	}

	return nil
}

// getNodeColor returns color based on device status
func (g *GraphVizGenerator) getNodeColor(status string) string {
	switch status {
	case "success":
		return "#90EE90" // Light green
	case "partial":
		return "#FFD700" // Gold
	case "failed":
		return "#FFB6C1" // Light red
	default:
		return "#D3D3D3" // Light gray
	}
}

// escapeLabel escapes special characters in labels
func (g *GraphVizGenerator) escapeLabel(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// escapeID escapes special characters in node IDs
func (g *GraphVizGenerator) escapeID(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// GenerateDOTFile writes DOT source to a file
func (g *GraphVizGenerator) GenerateDOTFile(graph *mapper.Graph, filePath string) error {
	dotSource := g.generateDOT(graph)
	
	// Write to file
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cat > %s", filePath))
	cmd.Stdin = strings.NewReader(dotSource)
	
	return cmd.Run()
}
