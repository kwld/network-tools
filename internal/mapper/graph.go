package mapper

import (
	"fmt"

	"github.com/kwld/network-tools/pkg/models"
)

// Graph represents a network topology graph
type Graph struct {
	Nodes []Node
	Edges []Edge
}

// Node represents a device in the graph
type Node struct {
	ID       string
	Label    string
	IP       string
	Vendor   string
	Model    string
	Status   string
}

// Edge represents a connection between devices
type Edge struct {
	Source      string
	Target      string
	SourcePort  string
	TargetPort  string
	Protocol    string
}

// BuildGraph creates a graph structure from topology
func BuildGraph(topology *models.Topology) *Graph {
	graph := &Graph{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}

	// Create nodes
	for _, device := range topology.Devices {
		nodeID := device.Hostname
		if nodeID == "" {
			nodeID = device.IP
		}

		node := Node{
			ID:     nodeID,
			Label:  fmt.Sprintf("%s\\n%s", nodeID, device.IP),
			IP:     device.IP,
			Vendor: device.Vendor,
			Model:  device.Model,
			Status: device.ScanStatus,
		}
		graph.Nodes = append(graph.Nodes, node)
	}

	// Create edges
	for _, conn := range topology.Connections {
		edge := Edge{
			Source:      conn.SourceDevice,
			Target:      conn.TargetDevice,
			SourcePort:  conn.SourcePort,
			TargetPort:  conn.TargetPort,
			Protocol:    conn.Protocol,
		}
		graph.Edges = append(graph.Edges, edge)
	}

	return graph
}
