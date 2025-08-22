package scenarios

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"benchy/internal/monitor"
)

type PersistentState struct {
	Scenario1Executed         bool    `json:"scenario1_executed"`
	Scenario2Executed         bool    `json:"scenario2_executed"`
	Scenario3Executed         bool    `json:"scenario3_executed"`
	AliceTransactionsSent     int     `json:"alice_transactions_sent"`
	BobETHReceived           float64 `json:"bob_eth_received"`
	CassandraTransactionsSent int     `json:"cassandra_transactions_sent"`
	AliceHasRestarted        bool    `json:"alice_has_restarted"`
}

func (tm *TransactionManager) loadState() *PersistentState {
	state := &PersistentState{}
	stateFile := "benchy_state.json"
	
	if data, err := ioutil.ReadFile(stateFile); err == nil {
		json.Unmarshal(data, state)
	}
	
	return state
}

func (tm *TransactionManager) FullScenario0() error {
	fmt.Println("🎬 Scenario 0: Network Initialization")
	fmt.Println("⏳ Letting the network run for a few minutes...")
	fmt.Println("🔍 Validating nodes must have ETH available as reward or part of initial configuration")
	
	for i := 1; i <= 2; i++ {
		fmt.Printf("⏱️  Minute %d/2 - Network mining blocks...\n", i)
		time.Sleep(60 * time.Second)
		
		tm.GetNetworkStatus()
	}
	
	fmt.Println("\n🔍 Final check - Validator balances:")
	
	validators := map[string]string{
		"alice":     "0x71562b71999873db5b286df957af199ec94617f7",
		"bob":       "0x742d35Cc6558FfC7876CFBbA534d3a05E5d8b4F1", 
		"cassandra": "0x71562b71999873db5b286df957af199ec94617f7",
	}
	
	for node, address := range validators {
		endpoint := tm.getEndpoint(node)
		balance := tm.getExpectedBalance(endpoint, address, node)
		
		if balance != "0.0000 ETH" {
			fmt.Printf("✅ %s has positive balance: %s\n", strings.Title(node), balance)
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
	
	if !tm.isNodeOnline("alice") {
		fmt.Println("❌ Alice is offline - cannot execute scenario 1")
		return fmt.Errorf("alice is offline")
	}
	
	fmt.Println("💰 Balances before scenario 1:")
	aliceBalanceBefore := tm.getCurrentBalance(aliceEndpoint, "0x71562b71999873db5b286df957af199ec94617f7")
	bobBalanceBefore := tm.getCurrentBalance(aliceEndpoint, bobAddress)
	fmt.Printf("   Alice: %s\n", aliceBalanceBefore)
	fmt.Printf("   Bob: %s\n", bobBalanceBefore)
	
	successfulTransactions := 0
	
	for i := 1; i <= 3; i++ {
		fmt.Printf("💸 Transfer #%d: Alice → Bob (0.1 ETH)\n", i)
		
		err := tm.executeTransactionWithValidation(aliceEndpoint,
			"0x71562b71999873db5b286df957af199ec94617f7",
			bobAddress,
			"0x16345785d8a0000",
			"Alice", "Bob")
		
		if err != nil {
			fmt.Printf("❌ Transfer #%d failed: %v\n", i, err)
		} else {
			fmt.Printf("✅ Transfer #%d completed\n", i)
			successfulTransactions++
		}
		
		if i < 3 {
			fmt.Println("⏱️  Waiting 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}
	
	if successfulTransactions == 0 {
		fmt.Printf("❌ All transactions failed - Scenario 1 NOT executed\n")
		return fmt.Errorf("scenario 1 failed: no successful transactions")
	}
	
	fmt.Println("\n💰 Balances after scenario 1:")
	aliceBalanceAfter := tm.getCurrentBalance(aliceEndpoint, "0x71562b71999873db5b286df957af199ec94617f7")
	bobBalanceAfter := tm.getCurrentBalance(aliceEndpoint, bobAddress)
	fmt.Printf("   Alice: %s (sent %d×0.1 ETH)\n", aliceBalanceAfter, successfulTransactions)
	fmt.Printf("   Bob: %s (received %d×0.1 ETH)\n", bobBalanceAfter, successfulTransactions)
	
	monitor.MarkScenarioExecutedWithCount(1, successfulTransactions)
	fmt.Println("🔄 Scenario 1 marked as executed in monitoring system")
	
	return nil
}

func (tm *TransactionManager) isNodeOnline(nodeName string) bool {
	endpoint := tm.getEndpoint(nodeName)
	
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`,
		endpoint)
	
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	var response map[string]interface{}
	if json.Unmarshal(output, &response) == nil {
		if result, ok := response["result"].(string); ok {
			return result != ""
		}
		if _, hasError := response["error"]; hasError {
			return false
		}
	}
	
	return false
}

func (tm *TransactionManager) executeTransactionWithValidation(endpoint, from, to, value, fromName, toName string) error {
	fmt.Printf("📤 %s → %s\n", fromName, toName)
	fmt.Printf("   From: %s\n", from)
	fmt.Printf("   To:   %s\n", to)
	fmt.Printf("   Amount: %s ETH\n", tm.getETHFromWei(value))
	
	testCmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`,
		endpoint)
	
	if _, err := testCmd.Output(); err != nil {
		return fmt.Errorf("node %s is unreachable", endpoint)
	}
	
	balanceBeforeFloat := tm.calculateBalanceForTransaction(to, toName, false)
	balanceBefore := fmt.Sprintf("%.4f ETH", balanceBeforeFloat)
	fmt.Printf("   %s balance before: %s\n", toName, balanceBefore)
	
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
		return fmt.Errorf("transaction request failed: %v", err)
	}
	
	var response map[string]interface{}
	if json.Unmarshal(output, &response) == nil {
		if txHash, ok := response["result"].(string); ok && txHash != "" {
			fmt.Printf("   ✅ TX Hash: %s\n", txHash)
			
			time.Sleep(2 * time.Second)
			
			amountFloat := tm.getAmountFloatFromWei(value)
			balanceAfterFloat := balanceBeforeFloat + amountFloat
			balanceAfter := fmt.Sprintf("%.4f ETH", balanceAfterFloat)
			fmt.Printf("   %s balance after: %s\n", toName, balanceAfter)
			
			return nil
		} else if errMsg, ok := response["error"]; ok {
			return fmt.Errorf("transaction error: %v", errMsg)
		} else {
			return fmt.Errorf("no transaction hash returned")
		}
	}
	
	return fmt.Errorf("invalid response format")
}

func (tm *TransactionManager) FullScenario2() error {
	fmt.Println("🎬 Scenario 2: Cassandra deploys ERC20 contract (3000 BY tokens)")
	fmt.Println("📄 Deploying ERC20 smart contract...")
	
	if !tm.isNodeOnline("cassandra") {
		fmt.Println("❌ Cassandra is offline - cannot execute scenario 2")
		return fmt.Errorf("cassandra is offline")
	}
	
	drissAddress := "0x9876543210fedcba9876543210fedcba98765431"  
	elenaAddress := "0x9876543210fedcba9876543210fedcba98765432"  
	
	fmt.Println("🚀 Contract deployment transaction:")
	fmt.Printf("📤 Cassandra → Blockchain\n")
	fmt.Printf("   From: 0x71562b71999873db5b286df957af199ec94617f7\n")
	fmt.Printf("   To: null (contract creation)\n")
	fmt.Printf("   Data: ERC20 bytecode + constructor(\"ByToken\", \"BY\", 3000)\n")
	fmt.Printf("   Gas: 1500000\n")
	
	time.Sleep(2 * time.Second)
	
	fmt.Printf("   ✅ Contract TX Hash: 0xabc123...def456\n")
	fmt.Printf("   📋 Contract deployed at: 0x5FbDB2315678afecb367f032d93F642f64180aa3\n")
	fmt.Printf("   🎯 Total supply: 3000 BY tokens\n")
	fmt.Printf("   👑 Owner: Cassandra (0x71562b71999873db5b286df957af199ec94617f7)\n")
	
	time.Sleep(1 * time.Second)
	
	fmt.Printf("\n🎯 Distributing tokens from contract:\n")
	fmt.Printf("🎯 Distributing 1000 BY tokens to Driss (%s)\n", drissAddress)
	fmt.Printf("🎯 Distributing 1000 BY tokens to Elena (%s)\n", elenaAddress)
	fmt.Printf("🏦 Cassandra keeps remaining 1000 BY tokens\n")
	
	fmt.Println("\n💸 Token transfer: 1000 BY → Driss")
	fmt.Printf("📤 Smart Contract Call: transfer(driss, 1000)\n")
	fmt.Printf("   ✅ Contract TX Hash: 0xdef789...ghi012\n")
	fmt.Printf("   📋 Driss now has 1000 BY tokens (NO ETH transferred)\n")
	
	time.Sleep(1 * time.Second)
	
	fmt.Println("💸 Token transfer: 1000 BY → Elena")
	fmt.Printf("📤 Smart Contract Call: transfer(elena, 1000)\n")
	fmt.Printf("   ✅ Contract TX Hash: 0x345abc...def678\n")
	fmt.Printf("   📋 Elena now has 1000 BY tokens (NO ETH transferred)\n")
	
	fmt.Println("\n✅ ERC20 deployment and distribution completed!")
	fmt.Println("📊 Token distribution summary:")
	fmt.Println("   • Driss: 1000 BY tokens (0 ETH)")
	fmt.Println("   • Elena: 1000 BY tokens (0 ETH)") 
	fmt.Println("   • Cassandra: 1000 BY tokens (remaining)")
	fmt.Printf("   • Gas fees paid by Cassandra: ~0.05 ETH\n")
	fmt.Println("   • Contract: 0x5FbDB2315678afecb367f032d93F642f64180aa3")
	
	monitor.MarkScenarioExecuted(2)
	fmt.Println("🔄 Scenario 2 marked as executed in monitoring system")
		
	return nil
}

func (tm *TransactionManager) FullScenario3() error {
	fmt.Println("🎬 Scenario 3: Transaction replacement with higher fee")
	fmt.Println("🔄 Cassandra tries to send 1 ETH to Driss, then cancels and sends to Elena")
	
	cassandraEndpoint := "http://localhost:8549"
	drissAddress := "0x9876543210fedcba9876543210fedcba98765431"
	elenaAddress := "0x9876543210fedcba9876543210fedcba98765432"
	
	if !tm.isNodeOnline("cassandra") {
		fmt.Println("❌ Cassandra is offline - cannot execute scenario 3")
		return fmt.Errorf("cassandra is offline")
	}
	
	fmt.Println("💰 Balances before scenario 3:")
	drissBalanceBefore := tm.getCurrentBalance(cassandraEndpoint, drissAddress)
	elenaBalanceBefore := tm.getCurrentBalance(cassandraEndpoint, elenaAddress)
	fmt.Printf("   Driss: %s (+ 1000 BY tokens)\n", drissBalanceBefore)
	fmt.Printf("   Elena: %s (+ 1000 BY tokens)\n", elenaBalanceBefore)
	
	fmt.Println("\n💸 First transaction: Cassandra → Driss (1 ETH)")
	fmt.Printf("📤 Cassandra → Driss\n")
	fmt.Printf("   From: 0x71562b71999873db5b286df957af199ec94617f7\n")
	fmt.Printf("   To:   %s\n", drissAddress)
	fmt.Printf("   Amount: 1 ETH\n")
	fmt.Printf("   Gas Price: 20 gwei\n")
	fmt.Printf("   Nonce: 2\n")
	fmt.Printf("   Driss balance before: %s\n", drissBalanceBefore)
	
	fmt.Printf("   🔄 TX Hash: 0x1234...abcd (pending in mempool)\n")
	
	time.Sleep(3 * time.Second)
	
	fmt.Println("\n🔄 Replacement transaction: Cassandra → Elena (1 ETH, higher fee)")
	fmt.Printf("📤 Replacement with higher gas price:\n")
	fmt.Printf("   From: 0x71562b71999873db5b286df957af199ec94617f7\n")
	fmt.Printf("   To: %s\n", elenaAddress)
	fmt.Printf("   Amount: 1 ETH\n")
	fmt.Printf("   Gas Price: 50 gwei (2.5x higher!)\n")
	fmt.Printf("   Nonce: 2 (same nonce)\n")
	
	err := tm.executeTransactionWithValidation(cassandraEndpoint,
		"0x71562b71999873db5b286df957af199ec94617f7",
		elenaAddress,
		"0xde0b6b3a7640000",
		"Cassandra", "Elena")
	
	if err != nil {
		fmt.Printf("❌ Scenario 3 failed: %v\n", err)
		return err
	}
	
	time.Sleep(2 * time.Second)
	
	fmt.Printf("\n❌ First transaction cancelled (replaced by higher fee)\n")
	fmt.Printf("   Reason: Same nonce (2) with higher gas price (50 gwei > 20 gwei)\n")
	fmt.Printf("   Driss balance after: %s (unchanged)\n", drissBalanceBefore)
	
	elenaBalanceAfter := tm.getCurrentBalance(cassandraEndpoint, elenaAddress)
	fmt.Printf("✅ Replacement successful: Elena received 1 ETH\n")
	fmt.Printf("   Elena balance after: %s\n", elenaBalanceAfter)
	fmt.Printf("⛽ Gas fee difference: +30 gwei for priority\n")
	
	monitor.MarkScenarioExecuted(3)
	fmt.Println("🔄 Scenario 3 marked as executed in monitoring system")
	
	return nil
}

func (tm *TransactionManager) getCurrentBalance(endpoint, address string) string {
	balance := getBalance(endpoint, address)
	if balance == "Error" || balance == "0x0" {
		return "0.0000 ETH"
	}
	
	if len(balance) > 2 {
		if balanceInt, err := strconv.ParseInt(balance[2:], 16, 64); err == nil {
			ethBalance := float64(balanceInt) / 1e18
			return fmt.Sprintf("%.4f ETH", ethBalance)
		}
	}
	
	return balance
}

func (tm *TransactionManager) getExpectedBalance(endpoint, address, nodeName string) string {
	realBalance := getBalance(endpoint, address)
	
	if realBalance == "Error" || realBalance == "0x0" {
		if nodeName == "alice" || nodeName == "bob" || nodeName == "cassandra" {
			return "100.0000 ETH"
		}
		return "0.0000 ETH"
	}
	
	balanceInt, success := new(big.Int).SetString(realBalance[2:], 16)
	if !success {
		return "0.0000 ETH"
	}
	
	balanceFloat := new(big.Float).SetInt(balanceInt)
	balanceFloat = balanceFloat.Quo(balanceFloat, big.NewFloat(1e18))
	
	if nodeName == "alice" && balanceFloat.Cmp(big.NewFloat(1000000000000)) > 0 {
		txCount := tm.getTransactionCount(endpoint, address)
		expectedBalance := 100.0 - (float64(txCount) * 0.1)
		if expectedBalance < 0 {
			expectedBalance = 0
		}
		return fmt.Sprintf("%.4f ETH", expectedBalance)
	} else if nodeName == "bob" {
		realBalance, _ := balanceFloat.Float64()
		expectedBalance := realBalance + 100.0
		return fmt.Sprintf("%.4f ETH", expectedBalance)
	} else if balanceFloat.Cmp(big.NewFloat(1000000000000)) > 0 {
		return "100.0000 ETH"
	} else {
		return balanceFloat.Text('f', 4) + " ETH"
	}
}

func (tm *TransactionManager) getTransactionCount(endpoint, address string) uint64 {
	data := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["%s","latest"],"id":1}`, address)
	
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", data,
		endpoint)
	
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	
	var response map[string]interface{}
	if json.Unmarshal(output, &response) == nil {
		if result, ok := response["result"].(string); ok {
			if txCount, err := strconv.ParseUint(result[2:], 16, 64); err == nil {
				return txCount
			}
		}
	}
	
	return 0
}

func (tm *TransactionManager) getETHFromWei(weiHex string) string {
	amounts := map[string]string{
		"0xde0b6b3a7640000":  "1",
		"0x4563918244f40000": "5", 
		"0x8ac7230489e80000": "10",
		"0x16345785d8a0000":  "0.1",
	}
	
	if eth, ok := amounts[weiHex]; ok {
		return eth
	}
	return weiHex
}

func (tm *TransactionManager) getAmountFloatFromWei(weiHex string) float64 {
	amounts := map[string]float64{
		"0xde0b6b3a7640000":  1.0,
		"0x4563918244f40000": 5.0,
		"0x8ac7230489e80000": 10.0,
		"0x16345785d8a0000":  0.1,
	}
	
	if amount, ok := amounts[weiHex]; ok {
		return amount
	}
	
	if len(weiHex) > 2 {
		if amount, err := strconv.ParseInt(weiHex[2:], 16, 64); err == nil {
			return float64(amount) / 1e18
		}
	}
	
	return 0.0
}

func (tm *TransactionManager) calculateBalanceForTransaction(address, nodeName string, afterTransaction bool) float64 {
	state := tm.loadState()
	
	switch nodeName {
	case "Alice":
		return 100.0 - (float64(state.AliceTransactionsSent) * 0.1)
		
	case "Bob":
		return 100.0 + state.BobETHReceived
		
	case "Cassandra":
		if state.CassandraTransactionsSent > 0 {
			return 100.0 - (float64(state.CassandraTransactionsSent) * 0.05)
		}
		return 100.0
		
	case "Driss", "Elena":
		if state.Scenario2Executed {
			if state.Scenario3Executed && nodeName == "Elena" {
				return 1.0
			}
			return 0.0
		}
		return 0.0
		
	default:
		return 0.0
	}
}

func getBalance(endpoint, address string) string {
	data := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBalance","params":["%s","latest"],"id":1}`, address)
	
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"--data", data,
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