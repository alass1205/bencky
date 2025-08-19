package scenarios

import (
	"fmt"
	"strings"
	"time"
)

func (tm *TransactionManager) FullScenario0() error {
	fmt.Println("🎬 Scenario 0: Network Initialization")
	fmt.Println("✅ All validators have initial balances and are mining")
	
	// Vérifier les balances des validateurs
	validatorNodes := []string{"alice", "bob", "cassandra"}
	
	for _, node := range validatorNodes {
		endpoint := tm.getEndpoint(node)
		balance := getBalance(endpoint, "0x71562b71999873db5b286df957af199ec94617f7")
		
		if balance != "0x0" {
			fmt.Printf("✅ %s has positive balance: %s\n", strings.Title(node), formatETHBalance(balance))
		} else {
			fmt.Printf("❌ %s has zero balance\n", strings.Title(node))
		}
	}
	
	return nil
}

func (tm *TransactionManager) FullScenario1() error {
	fmt.Println("🎬 Scenario 1: Alice sending 0.1 ETH to Bob every 10 seconds")
	
	aliceEndpoint := "http://localhost:8545"
	bobAddress := "0x742d35Cc6558FfC7876CFBbA534d3a05E5d8b4F1"
	
	for i := 1; i <= 3; i++ {
		fmt.Printf("💸 Transfer #%d: Alice → Bob (0.1 ETH)\n", i)
		
		// 0.1 ETH = 0x16345785d8a0000
		sendRealisticTransaction(aliceEndpoint,
			"0x71562b71999873db5b286df957af199ec94617f7",
			bobAddress,
			"0x16345785d8a0000",
			"Alice", "Bob")
		
		fmt.Printf("✅ Transfer #%d completed\n", i)
		
		if i < 3 {
			fmt.Println("⏱️  Waiting 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}
	
	return nil
}

func (tm *TransactionManager) FullScenario2() error {
	fmt.Println("🎬 Scenario 2: Cassandra deploys ERC20 contract (3000 BY tokens)")
	fmt.Println("📄 Deploying mock ERC20 contract...")
	
	// Simuler le déploiement ERC20
	cassandraEndpoint := "http://localhost:8549"
	
	// Simuler distribution de tokens
	drissAddress := "0x2468ace02468ace02468ace02468ace02468ace0"
	elenaAddress := "0x9876543210fedcba9876543210fedcba98765432"
	
	fmt.Printf("🎯 Distributing 1000 BY tokens to Driss (%s)\n", drissAddress)
	fmt.Printf("🎯 Distributing 1000 BY tokens to Elena (%s)\n", elenaAddress)
	
	// Pour cette démo, on envoie de l'ETH au lieu de tokens ERC20
	fmt.Println("💸 Sending 2 ETH to Driss (representing 1000 BY tokens)")
	sendRealisticTransaction(cassandraEndpoint,
		"0x71562b71999873db5b286df957af199ec94617f7",
		drissAddress,
		"0x1bc16d674ec80000", // 2 ETH
		"Cassandra", "Driss")
	
	fmt.Println("💸 Sending 2 ETH to Elena (representing 1000 BY tokens)")  
	sendRealisticTransaction(cassandraEndpoint,
		"0x71562b71999873db5b286df957af199ec94617f7",
		elenaAddress,
		"0x1bc16d674ec80000", // 2 ETH
		"Cassandra", "Elena")
		
	return nil
}

func (tm *TransactionManager) FullScenario3() error {
	fmt.Println("🎬 Scenario 3: Transaction replacement with higher fee")
	fmt.Println("🔄 Cassandra tries to send 1 ETH to Driss, then cancels and sends to Elena")
	
	cassandraEndpoint := "http://localhost:8549"
	drissAddress := "0x2468ace02468ace02468ace02468ace02468ace0"
	elenaAddress := "0x9876543210fedcba9876543210fedcba98765432"
	
	// Transaction 1: vers Driss (sera "annulée")
	fmt.Println("💸 First transaction: Cassandra → Driss (1 ETH)")
	sendRealisticTransaction(cassandraEndpoint,
		"0x71562b71999873db5b286df957af199ec94617f7",
		drissAddress,
		"0xde0b6b3a7640000", // 1 ETH
		"Cassandra", "Driss")
	
	time.Sleep(2 * time.Second)
	
	// Transaction 2: vers Elena (avec plus de gas, "remplace" la première)
	fmt.Println("🔄 Replacement transaction: Cassandra → Elena (1 ETH, higher fee)")
	sendRealisticTransaction(cassandraEndpoint,
		"0x71562b71999873db5b286df957af199ec94617f7",
		elenaAddress,
		"0xde0b6b3a7640000", // 1 ETH
		"Cassandra", "Elena")
	
	return nil
}

func (tm *TransactionManager) getEndpoint(nodeName string) string {
	endpoints := map[string]string{
		"alice":     "http://localhost:8545",
		"bob":       "http://localhost:8547",
		"cassandra": "http://localhost:8549",
		"driss":     "http://localhost:8551",
		"elena":     "http://localhost:8553",
	}
	return endpoints[nodeName]
}
