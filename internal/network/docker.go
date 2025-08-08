package network

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

type NodeConfig struct {
	Address   string `json:"address"`
	Validator bool   `json:"validator"`
	Client    string `json:"client"`
	RPCPort   int    `json:"rpc_port"`
	P2PPort   int    `json:"p2p_port"`
}

type BenchyNetwork struct {
	Nodes       map[string]NodeConfig
	NetworkName string
	Containers  map[string]string
}

func NewBenchyNetwork() (*BenchyNetwork, error) {
	nodes, err := loadNodeConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to load node configs: %v", err)
	}

	return &BenchyNetwork{
		Nodes:       nodes,
		NetworkName: "benchy-network",
		Containers:  make(map[string]string),
	}, nil
}

func loadNodeConfigs() (map[string]NodeConfig, error) {
	data, err := ioutil.ReadFile("configs/accounts.json")
	if err != nil {
		return nil, err
	}

	var nodes map[string]NodeConfig
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (bn *BenchyNetwork) LaunchNetwork(ctx context.Context) error {
	fmt.Println("🔧 Creating Docker network...")
	cmd := exec.Command("docker", "network", "create", bn.NetworkName)
	cmd.Run()

	fmt.Println("🚀 Launching nodes...")
	for nodeName, nodeConfig := range bn.Nodes {
		fmt.Printf("  Starting %s (%s)...\n", nodeName, nodeConfig.Client)
		
		if nodeConfig.Client == "geth" {
			if err := bn.launchGethNode(nodeName, nodeConfig); err != nil {
				return fmt.Errorf("failed to launch %s: %v", nodeName, err)
			}
		} else if nodeConfig.Client == "nethermind" {
			if err := bn.launchNethermindNode(nodeName, nodeConfig); err != nil {
				return fmt.Errorf("failed to launch %s: %v", nodeName, err)
			}
		}
	}

	fmt.Println("✅ All nodes launched successfully!")
	return nil
}

func (bn *BenchyNetwork) launchGethNode(nodeName string, nodeConfig NodeConfig) error {
	containerName := fmt.Sprintf("benchy-%s", nodeName)
	genesisPath, _ := filepath.Abs("configs/genesis.json")
	
	// Initialize genesis
	fmt.Printf("    Initializing genesis for %s...\n", nodeName)
	initArgs := []string{
		"run", "--rm",
		"-v", fmt.Sprintf("%s:/genesis.json", genesisPath),
		"-v", fmt.Sprintf("benchy-%s-data:/data", nodeName),
		"ethereum/client-go:latest",
		"init", "/genesis.json",
	}
	
	exec.Command("docker", initArgs...).Run()
	
	// Start node
	args := []string{
		"run", "-d",
		"--name", containerName,
		"--network", bn.NetworkName,
		"-p", fmt.Sprintf("%d:8545", nodeConfig.RPCPort),
		"-p", fmt.Sprintf("%d:30303", nodeConfig.P2PPort),
		"-v", fmt.Sprintf("benchy-%s-data:/data", nodeName),
		"ethereum/client-go:latest",
		"--datadir", "/data",
		"--networkid", "1337",
		"--http", "--http.addr", "0.0.0.0", "--http.port", "8545",
		"--http.api", "admin,debug,web3,eth,txpool,miner,net",
		"--http.corsdomain", "*",
		"--port", "30303",
		"--allow-insecure-unlock",
		"--syncmode", "full",
	}

	// CORRECTION: Supprimer --miner.threads, utiliser seulement --mine
	if nodeConfig.Validator {
		args = append(args, "--mine", "--miner.etherbase", nodeConfig.Address)
	}

	exec.Command("docker", args...).Run()
	fmt.Printf("    ✅ %s container started\n", nodeName)
	return nil
}

func (bn *BenchyNetwork) launchNethermindNode(nodeName string, nodeConfig NodeConfig) error {
	containerName := fmt.Sprintf("benchy-%s", nodeName)
	
	// CORRECTION: Utiliser l'image Nethermind
	args := []string{
		"run", "-d",
		"--name", containerName,
		"--network", bn.NetworkName,
		"-p", fmt.Sprintf("%d:8545", nodeConfig.RPCPort),
		"-p", fmt.Sprintf("%d:30303", nodeConfig.P2PPort),
		"-v", fmt.Sprintf("benchy-%s-data:/data", nodeName),
		"nethermind/nethermind:latest",
		"--datadir", "/data",
		"--JsonRpc.Enabled", "true",
		"--JsonRpc.Host", "0.0.0.0",
		"--JsonRpc.Port", "8545",
		"--Network.NetworkId", "1337",
		"--Network.P2PPort", "30303",
	}

	if nodeConfig.Validator {
		args = append(args, "--Init.ChainSpecPath", "null", "--Mining.Enabled", "true")
	}

	exec.Command("docker", args...).Run()
	fmt.Printf("    ✅ %s container started\n", nodeName)
	return nil
}
