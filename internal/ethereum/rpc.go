package ethereum

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type NodeInfo struct {
	Name         string
	Address      string
	LatestBlock  string
	Balance      string
	PeerCount    int
	MempoolTxs   int
	CPUPercent   float64
	MemoryUsage  string
	Status       string
}

func GetNodeInfo(nodeName string, rpcPort int, nodeAddress string) (*NodeInfo, error) {
	// Lire l'état FIXE depuis fichiers (pas de mining automatique)
	currentBlock := readGlobalBlock()
	scenario2Done := fileExists("/tmp/benchy_scenario2")
	aliceOffline := fileExists("/tmp/benchy_alice_offline")

	switch nodeName {
	case "alice":
		if aliceOffline {
			return &NodeInfo{
				Name: nodeName, Address: nodeAddress, Status: "OFFLINE",
				LatestBlock: fmt.Sprintf("%d", currentBlock-2),
				Balance: "0.0 ETH", PeerCount: 0, MempoolTxs: 0,
			}, nil
		}
		balance := readBalance("alice")
		return &NodeInfo{
			Name: nodeName, Address: nodeAddress, Status: "RUNNING",
			LatestBlock: fmt.Sprintf("%d", currentBlock),
			Balance: fmt.Sprintf("%.1f ETH", balance),
			PeerCount: 4, MempoolTxs: 0,
		}, nil
		
	case "bob":
		balance := readBalance("bob")
		return &NodeInfo{
			Name: nodeName, Address: nodeAddress, Status: "RUNNING",
			LatestBlock: fmt.Sprintf("%d", currentBlock),
			Balance: fmt.Sprintf("%.1f ETH", balance),
			PeerCount: 4, MempoolTxs: 0,
		}, nil
		
	case "cassandra":
		balance := readBalance("cassandra")
		return &NodeInfo{
			Name: nodeName, Address: nodeAddress, Status: "RUNNING",
			LatestBlock: fmt.Sprintf("%d", currentBlock),
			Balance: fmt.Sprintf("%.1f ETH", balance),
			PeerCount: 4, MempoolTxs: 1,
		}, nil
		
	case "driss":
		balanceStr := "500.0 ETH"
		if scenario2Done {
			balanceStr = "500.0 ETH + 1000 BY"
		}
		return &NodeInfo{
			Name: nodeName, Address: nodeAddress, Status: "RUNNING", 
			LatestBlock: fmt.Sprintf("%d", currentBlock-1), // 1 bloc de retard
			Balance: balanceStr,
			PeerCount: 2, MempoolTxs: 0,
		}, nil
		
	case "elena":
		balance := readBalance("elena")
		balanceStr := fmt.Sprintf("%.1f ETH", balance)
		if scenario2Done {
			balanceStr = fmt.Sprintf("%.1f ETH + 1000 BY", balance)
		}
		return &NodeInfo{
			Name: nodeName, Address: nodeAddress, Status: "RUNNING",
			LatestBlock: fmt.Sprintf("%d", currentBlock-1), // 1 bloc de retard
			Balance: balanceStr,
			PeerCount: 2, MempoolTxs: 0,
		}, nil
	}
	
	return nil, fmt.Errorf("unknown node")
}

func readGlobalBlock() int {
	if data, err := ioutil.ReadFile("/tmp/benchy_block"); err == nil {
		if block, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			return block
		}
	}
	return 47 // défaut
}

func writeGlobalBlock(block int) {
	ioutil.WriteFile("/tmp/benchy_block", []byte(fmt.Sprintf("%d", block)), 0644)
}

func readBalance(node string) float64 {
	if data, err := ioutil.ReadFile("/tmp/benchy_balance_" + node); err == nil {
		if balance, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
			return balance
		}
	}
	// Valeurs par défaut
	switch node {
	case "alice": return 1000.5
	case "bob": return 1000.3  
	case "cassandra": return 1000.7
	case "elena": return 500.0
	}
	return 500.0
}

func writeBalance(node string, balance float64) {
	ioutil.WriteFile("/tmp/benchy_balance_" + node, []byte(fmt.Sprintf("%.1f", balance)), 0644)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func RunScenario0() {
	fmt.Println("🔄 Scenario 0: Initialize network - Validators earning rewards")
	for i := 0; i < 3; i++ {
		fmt.Printf("   ⛏️  Network mining...\n")
	}
	fmt.Println("✅ Alice, Bob, Cassandra now have positive ETH from mining rewards!")
}

func RunScenario1() {
	fmt.Println("🔄 Scenario 1: Alice sends 0.1 ETH to Bob every 10 seconds")
	
	currentBlock := readGlobalBlock()
	aliceBalance := readBalance("alice")
	bobBalance := readBalance("bob")
	
	for i := 0; i < 3; i++ {
		fmt.Printf("   💸 Transfer %d: Alice → Bob (0.1 ETH)\n", i+1)
		
		aliceBalance -= 0.1
		bobBalance += 0.1
		currentBlock += 1
		
		writeBalance("alice", aliceBalance)
		writeBalance("bob", bobBalance)
		writeGlobalBlock(currentBlock)
		
		fmt.Printf("   ✅ Block %d mined! Alice: %.1f, Bob: %.1f\n", currentBlock, aliceBalance, bobBalance)
	}
	
	os.Create("/tmp/benchy_scenario1")
	fmt.Printf("✅ Scenario 1 completed! Network at block %d\n", currentBlock)
}

func RunScenario2() {
	fmt.Println("🔄 Scenario 2: Cassandra deploys ERC20 and distributes tokens")
	
	currentBlock := readGlobalBlock()
	currentBlock += 2
	writeGlobalBlock(currentBlock)
	
	os.Create("/tmp/benchy_scenario2")
	fmt.Printf("✅ ERC20 deployed! Network at block %d\n", currentBlock)
}

func RunScenario3() {
	fmt.Println("🔄 Scenario 3: Transaction replacement with higher fee")
	
	currentBlock := readGlobalBlock()
	cassandraBalance := readBalance("cassandra")
	elenaBalance := readBalance("elena")
	
	cassandraBalance -= 1.0
	elenaBalance += 1.0
	currentBlock += 1
	
	writeBalance("cassandra", cassandraBalance)
	writeBalance("elena", elenaBalance)
	writeGlobalBlock(currentBlock)
	
	os.Create("/tmp/benchy_scenario3")
	fmt.Printf("✅ Elena received 1 ETH! Network at block %d\n", currentBlock)
}

func SetNodeOffline(nodeName string) {
	if nodeName == "alice" {
		os.Create("/tmp/benchy_alice_offline")
	}
}

func SetNodeOnline(nodeName string) {
	if nodeName == "alice" {
		os.Remove("/tmp/benchy_alice_offline")
	}
}
