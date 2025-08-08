package monitor

import (
	"os/exec"
	"strconv"
	"strings"
)

type DockerStats struct {
	CPUPercent  float64
	MemoryUsage string
}

func GetDockerStats(containerName string) (*DockerStats, error) {
	cmd := exec.Command("docker", "stats", containerName, "--no-stream", "--format", "{{.CPUPerc}},{{.MemUsage}}")
	output, err := cmd.Output()
	if err != nil {
		return &DockerStats{CPUPercent: 0, MemoryUsage: "N/A"}, err
	}

	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(parts) != 2 {
		return &DockerStats{CPUPercent: 0, MemoryUsage: "N/A"}, nil
	}

	cpuStr := strings.TrimSuffix(parts[0], "%")
	cpuPercent, _ := strconv.ParseFloat(cpuStr, 64)

	return &DockerStats{
		CPUPercent:  cpuPercent,
		MemoryUsage: parts[1],
	}, nil
}
