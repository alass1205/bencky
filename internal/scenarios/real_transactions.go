package scenarios

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func SendRealETH() error {
	fmt.Println("💸 Sending REAL ETH transactions on multiple chains!")
	fmt.Println(strings.Repeat("=", 60))
	
	// Test Alice → Alice (même nœud = ça marche toujours)
	fmt.Println("🔄 Test 1: Alice → Alice (self-transfer on Alice's chain)")
	testTransaction("http://localhost:8545", 
		"0x71562b71999873db5b286df957af199ec94617f7", // from
		"0x71562b71999873db5b286df957af199ec94617f7", // to (same)
		"0x16345785d8a0000") // 0.1 ETH
	
	fmt.Println("\n🔄 Test 2: Bob → Bob (self-transfer on Bob's chain)")
	testTransaction("http://localhost:8547",
		"0x71562b71999873db5b286df957af199ec94617f7", // Bob's account
		"0x71562b71999873db5b286df957af199ec94617f7", // to same
		"0x16345785d8a0000") // 0.1 ETH
	
	fmt.Println("\n🔄 Test 3: Creating new account on Alice's chain")
	createAndFundAccount("http://localhost:8545")
	
	return nil
}

func testTransaction(endpoint, from, to, value string) {
	fmt.Printf("📍 Endpoint: %s\n", endpoint)
	
	// Vérifier balance avant
	balanceBefore := getBalance(endpoint, from)
	fmt.Printf("💰 Balance before: %s\n", balanceBefore)
	
	// Envoyer transaction
	transactionData := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_sendTransaction", 
		"params": [{
			"from": "%s",
			"to": "%s",
			"value": "%s",
			"gas": "0x5208"
		}],
		"id": 1
	}`, from, to, value)
	
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", transactionData,
		endpoint)
	
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	
	var response map[string]interface{}
	if json.Unmarshal(output, &response) == nil {
		if txHash, ok := response["result"].(string); ok {
			fmt.Printf("✅ Transaction: %s\n", txHash)
			
			// Attendre un peu
			time.Sleep(2 * time.Second)
			
			// Vérifier balance après
			balanceAfter := getBalance(endpoint, from)
			fmt.Printf("💰 Balance after: %s\n", balanceAfter)
		} else if errMsg, ok := response["error"]; ok {
			fmt.Printf("❌ Error: %v\n", errMsg)
		}
	}
}

func createAndFundAccount(endpoint string) {
	// Créer un nouveau compte
	fmt.Println("🆕 Creating new account...")
	
	createData := `{
		"jsonrpc": "2.0",
		"method": "personal_newAccount",
		"params": [""],
		"id": 1
	}`
	
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", createData,
		endpoint)
	
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Failed to create account: %v\n", err)
		return
	}
	
	var response map[string]interface{}
	if json.Unmarshal(output, &response) == nil {
		if newAddress, ok := response["result"].(string); ok {
			fmt.Printf("🎯 New account: %s\n", newAddress)
			
			// Envoyer des fonds au nouveau compte
			fmt.Println("💸 Funding new account with 5 ETH...")
			testTransaction(endpoint, 
				"0x71562b71999873db5b286df957af199ec94617f7", 
				newAddress, 
				"0x4563918244f40000") // 5 ETH
		}
	}
}

func getBalance(endpoint, address string) string {
	balanceData := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_getBalance",
		"params": ["%s", "latest"],
		"id": 1
	}`, address)
	
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", balanceData,
		endpoint)
	
	output, err := cmd.Output()
	if err != nil {
		return "Error"
	}
	
	var response map[string]interface{}
	if json.Unmarshal(output, &response) == nil {
		if balance, ok := response["result"].(string); ok {
			return balance
		}
	}
	
	return "0x0"
}
