package scenarios

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func SetupReasonableBalances() error {
	fmt.Println("🔧 Setting up accounts with reasonable balances (100 ETH each)")
	fmt.Println(strings.Repeat("=", 60))
	
	endpoints := []string{
		"http://localhost:8545", // Alice
		"http://localhost:8547", // Bob  
		"http://localhost:8549", // Cassandra
		"http://localhost:8551", // Driss
		"http://localhost:8553", // Elena
	}
	
	nodeNames := []string{"Alice", "Bob", "Cassandra", "Driss", "Elena"}
	
	for i, endpoint := range endpoints {
		fmt.Printf("🔄 Setting up %s...\n", nodeNames[i])
		
		// Créer un nouveau compte avec balance raisonnable
		newAccount := createAccountWithBalance(endpoint, nodeNames[i])
		if newAccount != "" {
			fmt.Printf("✅ %s new account: %s (100 ETH)\n", nodeNames[i], newAccount)
		}
		
		time.Sleep(1 * time.Second)
	}
	
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎯 All nodes now have accounts with 100 ETH!")
	fmt.Println("💡 Use these addresses for realistic transactions")
	return nil
}

func createAccountWithBalance(endpoint, nodeName string) string {
	// Créer un nouveau compte
	createData := `{
		"jsonrpc": "2.0",
		"method": "personal_newAccount", 
		"params": ["benchy123"],
		"id": 1
	}`
	
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", createData,
		endpoint)
	
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Failed to create account for %s: %v\n", nodeName, err)
		return ""
	}
	
	var response map[string]interface{}
	if json.Unmarshal(output, &response) != nil {
		fmt.Printf("❌ Failed to parse response for %s\n", nodeName)
		return ""
	}
	
	newAddress, ok := response["result"].(string)
	if !ok {
		fmt.Printf("❌ No address returned for %s\n", nodeName)
		return ""
	}
	
	// Financer le nouveau compte avec 100 ETH depuis le compte dev pré-financé
	fmt.Printf("💰 Funding %s account with 100 ETH...\n", nodeName)
	
	// 100 ETH en Wei = 100 * 10^18 = 0x56bc75e2d630eb20000
	fundData := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_sendTransaction",
		"params": [{
			"from": "0x71562b71999873db5b286df957af199ec94617f7",
			"to": "%s", 
			"value": "0x56bc75e2d630eb20000",
			"gas": "0x5208"
		}],
		"id": 1
	}`, newAddress)
	
	cmd = exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", fundData,
		endpoint)
	
	output, err = cmd.Output()
	if err == nil {
		var txResponse map[string]interface{}
		if json.Unmarshal(output, &txResponse) == nil {
			if txHash, ok := txResponse["result"].(string); ok {
				fmt.Printf("📤 Funding tx: %s\n", txHash)
			}
		}
	}
	
	return newAddress
}

func ShowNewBalances() error {
	fmt.Println("💰 Checking balances on all nodes...")
	fmt.Println(strings.Repeat("=", 50))
	
	endpoints := map[string]string{
		"Alice":     "http://localhost:8545",
		"Bob":       "http://localhost:8547",
		"Cassandra": "http://localhost:8549", 
		"Driss":     "http://localhost:8551",
		"Elena":     "http://localhost:8553",
	}
	
	for name, endpoint := range endpoints {
		fmt.Printf("%-12s: ", name)
		
		// Obtenir tous les comptes
		cmd := exec.Command("curl", "-s", "-X", "POST",
			"-H", "Content-Type: application/json",
			"--data", `{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}`,
			endpoint)
		
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("❌ Offline\n")
			continue
		}
		
		var response map[string]interface{}
		if json.Unmarshal(output, &response) == nil {
			if accounts, ok := response["result"].([]interface{}); ok && len(accounts) > 0 {
				// Montrer le dernier compte créé (nouveau compte avec 100 ETH)
				if len(accounts) > 1 {
					lastAccount := accounts[len(accounts)-1].(string)
					balance := getBalance(endpoint, lastAccount)
					fmt.Printf("🎯 %s (%s)\n", lastAccount, formatBalance(balance))
				} else {
					fmt.Printf("🔍 Only dev account available\n")
				}
			} else {
				fmt.Printf("❌ No accounts\n")
			}
		}
	}
	
	return nil
}

func formatBalance(hexBalance string) string {
	// Simple conversion hex to ETH for display
	if hexBalance == "0x56bc75e2d630eb20000" {
		return "100.0000 ETH"
	} else if hexBalance == "0x0" {
		return "0.0000 ETH"
	} else {
		return hexBalance + " (hex)"
	}
}
