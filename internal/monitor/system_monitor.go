package monitor

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ContainerStats struct {
	CPUUsage    float64
	MemoryUsage string
	MemoryLimit string
	IsRunning   bool
}

// Cache global pour éviter les appels répétés
var (
	statsCache      = make(map[string]*ContainerStats)
	statsCacheMutex sync.RWMutex
	lastCacheUpdate time.Time
	cacheDuration   = 2 * time.Second // Cache valide 2 secondes
)

// Vérifier l'état réel d'un conteneur sans cache
func isContainerRunning(containerName string) bool {
	cmd := exec.Command("docker", "inspect", containerName, "--format", "{{.State.Running}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// Obtenir toutes les stats en une seule commande
func getAllContainerStats() (map[string]*ContainerStats, error) {
	statsCacheMutex.Lock()
	defer statsCacheMutex.Unlock()
	
	// Utiliser le cache si récent
	if time.Since(lastCacheUpdate) < cacheDuration && len(statsCache) > 0 {
		result := make(map[string]*ContainerStats)
		for k, v := range statsCache {
			result[k] = v
		}
		return result, nil
	}
	
	// Obtenir les stats de tous les conteneurs benchy en une commande
	cmd := exec.Command("docker", "stats", "--no-stream", "--format",
		"{{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}", 
		"benchy-alice", "benchy-bob", "benchy-cassandra", "benchy-driss", "benchy-elena")
	
	output, err := cmd.Output()
	if err != nil {
		// Si la commande échoue, essayer individuellement pour les conteneurs en ligne
		return getIndividualStats(), nil
	}
	
	result := make(map[string]*ContainerStats)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}
		
		containerName := parts[0]
		cpuStr := strings.TrimSuffix(parts[1], "%")
		cpu, _ := strconv.ParseFloat(cpuStr, 64)
		memoryUsage := parts[2]
		
		result[containerName] = &ContainerStats{
			CPUUsage:    cpu,
			MemoryUsage: memoryUsage,
			IsRunning:   true,
		}
	}
	
	// Mettre à jour le cache
	statsCache = result
	lastCacheUpdate = time.Now()
	
	return result, nil
}

// Fallback : obtenir les stats individuellement
func getIndividualStats() map[string]*ContainerStats {
	containers := []string{"benchy-alice", "benchy-bob", "benchy-cassandra", "benchy-driss", "benchy-elena"}
	result := make(map[string]*ContainerStats)
	
	// Utiliser des goroutines pour paralléliser
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for _, container := range containers {
		wg.Add(1)
		go func(containerName string) {
			defer wg.Done()
			
			cmd := exec.Command("docker", "stats", "--no-stream", "--format",
				"{{.CPUPerc}}\t{{.MemUsage}}", containerName)
			
			output, err := cmd.Output()
			if err != nil {
				mu.Lock()
				result[containerName] = &ContainerStats{IsRunning: false}
				mu.Unlock()
				return
			}
			
			parts := strings.Split(strings.TrimSpace(string(output)), "\t")
			if len(parts) < 2 {
				mu.Lock()
				result[containerName] = &ContainerStats{IsRunning: false}
				mu.Unlock()
				return
			}
			
			cpuStr := strings.TrimSuffix(parts[0], "%")
			cpu, _ := strconv.ParseFloat(cpuStr, 64)
			memoryUsage := parts[1]
			
			mu.Lock()
			result[containerName] = &ContainerStats{
				CPUUsage:    cpu,
				MemoryUsage: memoryUsage,
				IsRunning:   true,
			}
			mu.Unlock()
		}(container)
	}
	
	wg.Wait()
	return result
}

func GetContainerStats(containerName string) (*ContainerStats, error) {
	// Vérifier d'abord si le conteneur existe vraiment (sans cache)
	if !isContainerRunning(containerName) {
		// Invalider le cache pour ce conteneur
		statsCacheMutex.Lock()
		delete(statsCache, containerName)
		statsCacheMutex.Unlock()
		
		return &ContainerStats{IsRunning: false}, nil
	}
	
	// Si le conteneur tourne, utiliser le cache normal
	allStats, err := getAllContainerStats()
	if err != nil {
		return &ContainerStats{IsRunning: false}, err
	}
	
	if stats, exists := allStats[containerName]; exists {
		return stats, nil
	}
	
	return &ContainerStats{IsRunning: false}, nil
}

// Optimiser aussi cette fonction avec goroutines
func GetDetailedNodeInfo(nodeName string) (map[string]interface{}, error) {
	node, exists := nodeEndpoints[nodeName]
	if !exists {
		return nil, fmt.Errorf("node %s not found", nodeName)
	}
	
	result := make(map[string]interface{})
	
	// Utiliser des channels pour paralléliser
	type blockResult struct {
		blockNumber int64
		err         error
	}
	
	type statsResult struct {
		stats *ContainerStats
		err   error
	}
	
	blockChan := make(chan blockResult, 1)
	statsChan := make(chan statsResult, 1)
	
	// Goroutine pour obtenir le numéro de bloc
	go func() {
		cmd := exec.Command("curl", "-s", "-X", "POST",
			"-H", "Content-Type: application/json",
			"--data", `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`,
			node)
		
		output, err := cmd.Output()
		if err != nil {
			blockChan <- blockResult{0, err}
			return
		}
		
		var response map[string]interface{}
		if json.Unmarshal(output, &response) == nil {
			if blockHex, ok := response["result"].(string); ok {
				if blockNum, err := strconv.ParseInt(blockHex[2:], 16, 64); err == nil {
					blockChan <- blockResult{blockNum, nil}
					return
				}
			}
		}
		blockChan <- blockResult{0, fmt.Errorf("failed to parse block")}
	}()
	
	// Goroutine pour obtenir les stats du conteneur
	go func() {
		containerName := fmt.Sprintf("benchy-%s", nodeName)
		stats, err := GetContainerStats(containerName)
		statsChan <- statsResult{stats, err}
	}()
	
	// Attendre les résultats avec timeout
	timeout := time.After(3 * time.Second)
	
	for i := 0; i < 2; i++ {
		select {
		case blockRes := <-blockChan:
			if blockRes.err == nil {
				result["blockNumber"] = blockRes.blockNumber
			}
		case statsRes := <-statsChan:
			if statsRes.err == nil {
				result["stats"] = statsRes.stats
			}
		case <-timeout:
			return result, fmt.Errorf("timeout getting node info")
		}
	}
	
	return result, nil
}

var nodeEndpoints = map[string]string{
	"alice":     "http://localhost:8545",
	"bob":       "http://localhost:8547",
	"cassandra": "http://localhost:8549",
	"driss":     "http://localhost:8551",
	"elena":     "http://localhost:8553",
}