package scenarios

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Comptes pré-définis avec clés privées connues pour demo
var testAccounts = map[string]string{
	"Alice_Test":     "0x742d35Cc6558FfC7876CFBbA534d3a05E5d8b4F1", 
	"Bob_Test":       "0x8ba1f109551bD432803012645Hf20D82DDD618C1",
	"Cassandra_Test": "0x1a2b3c4d5e6f7890abcdef1234567890abcdef12",
	"Driss_Test":     "0x2468ace02468ace02468ace02468ace02468ace0",
	"Elena_Test":     "0x9876543210fedcba9876543210fedcba98765432",
}

func RunRealisticDemo() error {
	fmt.Println("🎬 REALISTIC DEMO: Transactions with 1-10 ETH amounts")
	fmt.Println(strings.Repeat("=", 60))
	
	// Demo 1: Alice envoie 5 ETH à Bob
	fmt.Println("💸 Demo 1: Alice → Bob (5 ETH)")
	sendRealisticTransaction("http://localhost:8545", 
		"0x71562b71999873db5b286df957af199ec94617f7", // Alice dev account (énorme balance)
		"0x742d35Cc6558FfC7876CFBbA534d3a05E5d8b4F1", // Bob test account
		"0x4563918244f40000", // 5 ETH
		"Alice", "Bob")
	
	time.Sleep(3 * time.Second)
	
	// Demo 2: Transfer de 10 ETH vers un nouveau compte  
	fmt.Println("\n💸 Demo 2: Alice → Cassandra (10 ETH)")
	sendRealisticTransaction("http://localhost:8545",
		"0x71562b71999873db5b286df957af199ec94617f7", // Alice dev account
		"0x1a2b3c4d5e6f7890abcdef1234567890abcdef12", // Cassandra test account  
		"0x8ac7230489e80000", // 10 ETH
		"Alice", "Cassandra")
		
	time.Sleep(3 * time.Second)
	
	// Demo 3: Transfer plus petit (1 ETH)
	fmt.Println("\n💸 Demo 3: Alice → Driss (1 ETH)")
	sendRealisticTransaction("http://localhost:8545",
		"0x71562b71999873db5b286df957af199ec94617f7", // Alice dev account
		"0x2468ace02468ace02468ace02468ace02468ace0", // Driss test account
		"0xde0b6b3a7640000", // 1 ETH
		"Alice", "Driss")
	
	fmt.Println("\n🎯 Demo completed! Check balances with 'benchy check-demo-balances'")
	return nil
}

func CheckDemoBalances() error {
	fmt.Println("💰 Demo Account Balances:")
	fmt.Println(strings.Repeat("=", 50))
	
	endpoint := "http://localhost:8545" // Alice's node
	
	fmt.Printf("%-15s %-42s %-15s\n", "Account", "Address", "Balance")
	fmt.Println(strings.Repeat("-", 50))
	
	// Vérifier Alice (sender) 
	aliceBalance := getBalance(endpoint, "0x71562b71999873db5b286df957af199ec94617f7")
	fmt.Printf("%-15s %-42s %-15s\n", "Alice (sender)", "0x71562b71999873db5b286df957af199ec94617f7", formatETHBalance(aliceBalance))
	
	// Vérifier les comptes de destination
	accounts := map[string]string{
		"Bob":       "0x742d35Cc6558FfC7876CFBbA534d3a05E5d8b4F1",
		"Cassandra": "0x1a2b3c4d5e6f7890abcdef1234567890abcdef12", 
		"Driss":     "0x2468ace02468ace02468ace02468ace02468ace0",
		"Elena":     "0x9876543210fedcba9876543210fedcba98765432",
	}
	
	for name, address := range accounts {
		balance := getBalance(endpoint, address)
		fmt.Printf("%-15s %-42s %-15s\n", name, address, formatETHBalance(balance))
	}
	
	fmt.Println(strings.Repeat("=", 50))
	return nil
}

func sendRealisticTransaction(endpoint, from, to, value, fromName, toName string) {
	fmt.Printf("📤 %s → %s\n", fromName, toName)
	fmt.Printf("   From: %s\n", from)
	fmt.Printf("   To:   %s\n", to)
	fmt.Printf("   Amount: %s ETH\n", getETHFromWei(value))
	
	// Vérifier balance avant
	balanceBefore := getBalance(endpoint, to)
	fmt.Printf("   %s balance before: %s\n", toName, formatETHBalance(balanceBefore))
	
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
		fmt.Printf("   ❌ Transaction failed: %v\n", err)
		return
	}
	
	var response map[string]interface{}
	if json.Unmarshal(output, &response) == nil {
		if txHash, ok := response["result"].(string); ok {
			fmt.Printf("   ✅ TX Hash: %s\n", txHash)
			
			// Attendre que la transaction soit minée
			time.Sleep(2 * time.Second)
			
			// Vérifier balance après
			balanceAfter := getBalance(endpoint, to)
			fmt.Printf("   %s balance after: %s\n", toName, formatETHBalance(balanceAfter))
		} else if errMsg, ok := response["error"]; ok {
			fmt.Printf("   ❌ Error: %v\n", errMsg)
		}
	}
}

func getETHFromWei(weiHex string) string {
	amounts := map[string]string{
		"0xde0b6b3a7640000":  "1",
		"0x4563918244f40000": "5", 
		"0x8ac7230489e80000": "10",
	}
	
	if eth, ok := amounts[weiHex]; ok {
		return eth
	}
	return weiHex
}

func formatETHBalance(hexBalance string) string {
	// Conversion rapide des montants connus
	knownBalances := map[string]string{
		"0x0":                 "0.0000 ETH",
		"0xde0b6b3a7640000":  "1.0000 ETH",
		"0x4563918244f40000": "5.0000 ETH", 
		"0x8ac7230489e80000": "10.0000 ETH",
	}
	
	if balance, ok := knownBalances[hexBalance]; ok {
		return balance
	}
	
	// Pour les gros montants, afficher en notation simple
	if len(hexBalance) > 20 {
		return "HUGE ETH"
	}
	
	return hexBalance
}
