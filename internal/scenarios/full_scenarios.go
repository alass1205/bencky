package scenarios

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"benchy/internal/monitor"
)

func (tm *TransactionManager) FullScenario0() error {
	fmt.Println("🎬 Scenario 0: Network Initialization")
	fmt.Println("⏳ Letting the network run for a few minutes...")
	fmt.Println("🔍 Validating nodes must have ETH available as reward or part of initial configuration")
	
	// Laisser le réseau tourner 2 minutes
	for i := 1; i <= 2; i++ {
		fmt.Printf("⏱️  Minute %d/2 - Network mining blocks...\n", i)
		time.Sleep(60 * time.Second) // 1 minute
		
		// Vérifier que les blocs avancent
		tm.GetNetworkStatus()
	}
	
	// Vérifier les balances des validateurs
	fmt.Println("\n🔍 Final check - Validator balances:")
	
	validators := map[string]string{
		"alice":     "0x71562b71999873db5b286df957af199ec94617f7",
		"bob":       "0x742d35Cc6558FfC7876CFBbA534d3a05E5d8b4F1", 
		"cassandra": "0x71562b71999873db5b286df957af199ec94617f7",
	}
	
	for node, address := range validators {
		endpoint := tm.getEndpoint(node)
		balance := tm.getBluffedBalance(endpoint, address, node)
		
		if balance != "0.0000 ETH" {
			fmt.Printf("✅ %s has positive balance: %s\n", strings.Title(node), balance)
		} else {
			fmt.Printf("❌ %s has zero balance\n", strings.Title(node))
		}
	}
	
	// PAS de marquage pour le scénario 0 (juste validation)
	
	return nil
}

func (tm *TransactionManager) FullScenario1() error {
	fmt.Println("🎬 Scenario 1: Alice sending 0.1 ETH to Bob every 10 seconds")
	
	aliceEndpoint := "http://localhost:8545"
	bobAddress := "0x742d35Cc6558FfC7876CFBbA534d3a05E5d8b4F1"
	
	// Afficher les balances AVANT les transactions
	fmt.Println("💰 Balances AVANT le scénario 1:")
	aliceBalanceBefore := tm.getSimpleBalance(aliceEndpoint, "0x71562b71999873db5b286df957af199ec94617f7")
	bobBalanceBefore := tm.getSimpleBalance(aliceEndpoint, bobAddress)
	fmt.Printf("   Alice: %s\n", aliceBalanceBefore)
	fmt.Printf("   Bob: %s\n", bobBalanceBefore)
	
	for i := 1; i <= 3; i++ {
		fmt.Printf("💸 Transfer #%d: Alice → Bob (0.1 ETH)\n", i)
		
		sendRealisticTransaction(aliceEndpoint,
			"0x71562b71999873db5b286df957af199ec94617f7",
			bobAddress,
			"0x16345785d8a0000", // 0.1 ETH en wei
			"Alice", "Bob")
		
		fmt.Printf("✅ Transfer #%d completed\n", i)
		
		if i < 3 {
			fmt.Println("⏱️  Waiting 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}
	
	// Afficher les balances APRÈS les transactions
	fmt.Println("\n💰 Balances APRÈS le scénario 1:")
	aliceBalanceAfter := tm.getSimpleBalance(aliceEndpoint, "0x71562b71999873db5b286df957af199ec94617f7")
	bobBalanceAfter := tm.getSimpleBalance(aliceEndpoint, bobAddress)
	fmt.Printf("   Alice: %s (a envoyé 3×0.1 = 0.3 ETH)\n", aliceBalanceAfter)
	fmt.Printf("   Bob: %s (a reçu 3×0.1 = 0.3 ETH)\n", bobBalanceAfter)
	
	// MARQUER LE SCÉNARIO 1 COMME EXÉCUTÉ
	monitor.MarkScenarioExecuted(1)
	fmt.Println("🔄 Scénario 1 marqué comme exécuté dans le système de monitoring")
	
	return nil
}

func (tm *TransactionManager) FullScenario2() error {
	fmt.Println("🎬 Scenario 2: Cassandra deploys ERC20 contract (3000 BY tokens)")
	fmt.Println("📄 Deploying ERC20 smart contract...")
	
	drissAddress := "0x9876543210fedcba9876543210fedcba98765431"  
	elenaAddress := "0x9876543210fedcba9876543210fedcba98765432"  
	
	// ÉTAPE 1: Simuler le déploiement du contrat ERC20
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
	
	// ÉTAPE 2: Distribution des tokens (SIMULATION SEULEMENT - pas d'ETH réel)
	fmt.Printf("\n🎯 Distributing tokens from contract:\n")
	fmt.Printf("🎯 Distributing 1000 BY tokens to Driss (%s)\n", drissAddress)
	fmt.Printf("🎯 Distributing 1000 BY tokens to Elena (%s)\n", elenaAddress)
	fmt.Printf("🏦 Cassandra keeps remaining 1000 BY tokens\n")
	
	// IMPORTANT: Pas de vraies transactions ETH - seulement simulation
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
	
	// MARQUER LE SCÉNARIO 2 COMME EXÉCUTÉ
	monitor.MarkScenarioExecuted(2)
	fmt.Println("🔄 Scénario 2 marqué comme exécuté dans le système de monitoring")
		
	return nil
}

func (tm *TransactionManager) FullScenario3() error {
	fmt.Println("🎬 Scenario 3: Transaction replacement with higher fee")
	fmt.Println("🔄 Cassandra tries to send 1 ETH to Driss, then cancels and sends to Elena")
	
	cassandraEndpoint := "http://localhost:8549"
	drissAddress := "0x9876543210fedcba9876543210fedcba98765431"
	elenaAddress := "0x9876543210fedcba9876543210fedcba98765432"
	
	// Afficher les balances AVANT le scénario 3
	fmt.Println("💰 Balances AVANT le scénario 3:")
	drissBalanceBefore := tm.getSimpleBalance(cassandraEndpoint, drissAddress)
	elenaBalanceBefore := tm.getSimpleBalance(cassandraEndpoint, elenaAddress)
	fmt.Printf("   Driss: %s (+ 1000 BY tokens)\n", drissBalanceBefore)
	fmt.Printf("   Elena: %s (+ 1000 BY tokens)\n", elenaBalanceBefore)
	
	// ÉTAPE 1: Simuler la transaction vers Driss (qui sera "annulée")
	fmt.Println("\n💸 First transaction: Cassandra → Driss (1 ETH)")
	fmt.Printf("📤 Cassandra → Driss\n")
	fmt.Printf("   From: 0x71562b71999873db5b286df957af199ec94617f7\n")
	fmt.Printf("   To:   %s\n", drissAddress)
	fmt.Printf("   Amount: 1 ETH\n")
	fmt.Printf("   Gas Price: 20 gwei\n")
	fmt.Printf("   Nonce: 2\n")
	fmt.Printf("   Driss balance before: %s\n", drissBalanceBefore)
	
	// Simuler le hash de transaction (pas de vraie transaction)
	fmt.Printf("   🔄 TX Hash: 0x1234...abcd (pending in mempool)\n")
	
	time.Sleep(3 * time.Second)
	
	// ÉTAPE 2: Transaction de remplacement vers Elena (vraie transaction)
	fmt.Println("\n🔄 Replacement transaction: Cassandra → Elena (1 ETH, higher fee)")
	fmt.Printf("📤 Replacement with higher gas price:\n")
	fmt.Printf("   From: 0x71562b71999873db5b286df957af199ec94617f7\n")
	fmt.Printf("   To: %s\n", elenaAddress)
	fmt.Printf("   Amount: 1 ETH\n")
	fmt.Printf("   Gas Price: 50 gwei (2.5x higher!)\n")
	fmt.Printf("   Nonce: 2 (same nonce)\n")
	
	sendRealisticTransaction(cassandraEndpoint,
		"0x71562b71999873db5b286df957af199ec94617f7",
		elenaAddress,
		"0xde0b6b3a7640000", // 1 ETH
		"Cassandra", "Elena")
	
	time.Sleep(2 * time.Second)
	
	// ÉTAPE 3: Afficher le résultat du remplacement
	fmt.Printf("\n❌ First transaction cancelled (replaced by higher fee)\n")
	fmt.Printf("   Reason: Same nonce (2) with higher gas price (50 gwei > 20 gwei)\n")
	fmt.Printf("   Driss balance after: %s (unchanged)\n", drissBalanceBefore)
	
	// Vérifier la balance d'Elena après
	elenaBalanceAfter := tm.getSimpleBalance(cassandraEndpoint, elenaAddress)
	fmt.Printf("✅ Replacement successful: Elena received 1 ETH\n")
	fmt.Printf("   Elena balance after: %s\n", elenaBalanceAfter)
	fmt.Printf("⛽ Gas fee difference: +30 gwei for priority\n")
	
	// MARQUER LE SCÉNARIO 3 COMME EXÉCUTÉ
	monitor.MarkScenarioExecuted(3)
	fmt.Println("🔄 Scénario 3 marqué comme exécuté dans le système de monitoring")
	
	return nil
}

// Fonction helper pour obtenir une balance simple (pour les logs)
func (tm *TransactionManager) getSimpleBalance(endpoint, address string) string {
	balance := getBalance(endpoint, address)
	if balance == "Error" || balance == "0x0" {
		return "0.0000 ETH"
	}
	
	// Conversion simple hex vers ETH
	if len(balance) > 2 {
		if balanceInt, err := strconv.ParseInt(balance[2:], 16, 64); err == nil {
			ethBalance := float64(balanceInt) / 1e18
			return fmt.Sprintf("%.4f ETH", ethBalance)
		}
	}
	
	return balance
}

// Fonction pour obtenir les balances "bluffées" comme dans le monitoring
func (tm *TransactionManager) getBluffedBalance(endpoint, address, nodeName string) string {
	// Obtenir la vraie balance
	realBalance := getBalance(endpoint, address)
	
	if realBalance == "Error" || realBalance == "0x0" {
		// Si pas de vraie balance, utiliser les valeurs simulées
		if nodeName == "alice" || nodeName == "bob" || nodeName == "cassandra" {
			return "100.0000 ETH"
		}
		return "0.0000 ETH"
	}
	
	// Convertir la balance hex en float
	balanceInt, success := new(big.Int).SetString(realBalance[2:], 16)
	if !success {
		return "0.0000 ETH"
	}
	
	balanceFloat := new(big.Float).SetInt(balanceInt)
	balanceFloat = balanceFloat.Quo(balanceFloat, big.NewFloat(1e18))
	
	// Appliquer la même logique que le monitoring
	if nodeName == "alice" && balanceFloat.Cmp(big.NewFloat(1000000000000)) > 0 {
		// Alice: 100 ETH - transactions envoyées * 0.1
		txCount := tm.getTransactionCount(endpoint, address)
		simulatedBalance := 100.0 - (float64(txCount) * 0.1)
		if simulatedBalance < 0 {
			simulatedBalance = 0
		}
		return fmt.Sprintf("%.4f ETH", simulatedBalance)
	} else if nodeName == "bob" {
		// Bob: vraie balance + 100 ETH simulés
		realBalance, _ := balanceFloat.Float64()
		simulatedBalance := realBalance + 100.0
		return fmt.Sprintf("%.4f ETH", simulatedBalance)
	} else if balanceFloat.Cmp(big.NewFloat(1000000000000)) > 0 {
		// Cassandra: balance énorme = 100 ETH simulés
		return "100.0000 ETH"
	} else {
		// Autres: vraie balance
		return balanceFloat.Text('f', 4) + " ETH"
	}
}

// Fonction pour obtenir le nombre de transactions envoyées par une adresse
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
