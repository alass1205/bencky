package monitor

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type ContainerStats struct {
	CPUUsage    float64
	MemoryUsage string
	MemoryLimit string
	IsRunning   bool
}

func GetContainerStats(containerName string) (*ContainerStats, error) {
	// Utiliser docker stats pour obtenir les métriques
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", 
		"{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}", containerName)
	
	output, err := cmd.Output()
	if err != nil {
		return &ContainerStats{IsRunning: false}, nil
	}

	parts := strings.Split(strings.TrimSpace(string(output)), "\t")
	if len(parts) < 3 {
		return &ContainerStats{IsRunning: false}, nil
	}

	// Parser le CPU (ex: "1.23%")
	cpuStr := strings.TrimSuffix(parts[0], "%")
	cpu, _ := strconv.ParseFloat(cpuStr, 64)

	// Parser la mémoire (ex: "123.4MiB / 7.789GiB")
	memoryUsage := parts[1]

	return &ContainerStats{
		CPUUsage:    cpu,
		MemoryUsage: memoryUsage,
		IsRunning:   true,
	}, nil
}

func GetDetailedNodeInfo(nodeName string) (map[string]interface{}, error) {
	node, exists := nodeEndpoints[nodeName]
	if !exists {
		return nil, fmt.Errorf("node %s not found", nodeName)
	}

	result := make(map[string]interface{})
	
	// Informations blockchain
	cmd := exec.Command("curl", "-s", "-X", "POST", 
		"-H", "Content-Type: application/json",
		"--data", `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`,
		node)
	
	output, err := cmd.Output()
	if err == nil {
		var response map[string]interface{}
		if json.Unmarshal(output, &response) == nil {
			if blockHex, ok := response["result"].(string); ok {
				if blockNum, err := strconv.ParseInt(blockHex[2:], 16, 64); err == nil {
					result["blockNumber"] = blockNum
				}
			}
		}
	}

	// Statistiques conteneur
	containerName := fmt.Sprintf("benchy-%s", nodeName)
	stats, _ := GetContainerStats(containerName)
	result["stats"] = stats

	return result, nil
}

var nodeEndpoints = map[string]string{
	"alice":     "http://localhost:8545",
	"bob":       "http://localhost:8547", 
	"cassandra": "http://localhost:8549",
	"driss":     "http://localhost:8551",
	"elena":     "http://localhost:8553",
}
