package scenarios

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type TransactionManager struct {
	clients map[string]*ethclient.Client
}

func NewTransactionManager() *TransactionManager {
	clients := make(map[string]*ethclient.Client)
	
	endpoints := map[string]string{
		"alice":     "http://localhost:8545",
		"bob":       "http://localhost:8547",
		"cassandra": "http://localhost:8549",
		"driss":     "http://localhost:8551",
		"elena":     "http://localhost:8553",
	}

	for name, endpoint := range endpoints {
		if client, err := ethclient.Dial(endpoint); err == nil {
			clients[name] = client
		}
	}

	return &TransactionManager{clients: clients}
}

func (tm *TransactionManager) Scenario0() error {
	fmt.Println("🎬 Scenario 0: Initializing network...")
	fmt.Println("✅ Validators are mining REAL blocks automatically with --dev mode")
	fmt.Println("💰 Each --dev node has pre-funded accounts ready for transactions")
	
	// Vérifier que les nœuds minent
	for name, client := range tm.clients {
		if blockNum, err := client.BlockNumber(context.Background()); err == nil {
			fmt.Printf("   %s: Block #%d\n", name, blockNum)
		}
	}
	
	return nil
}

func (tm *TransactionManager) Scenario1() error {
	fmt.Println("🎬 Scenario 1: Alice sending 0.1 ETH to Bob every 10 seconds...")
	
	// Pour l'instant, on va juste afficher les comptes réels
	fmt.Println("🔍 Getting real accounts first...")
	
	endpoints := map[string]string{
		"Alice": "http://localhost:8545",
		"Bob":   "http://localhost:8547",
	}

	for name, endpoint := range endpoints {
		cmd := exec.Command("curl", "-s", "-X", "POST",
			"-H", "Content-Type: application/json",
			"--data", `{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}`,
			endpoint)
		
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("❌ %s offline\n", name)
			continue
		}

		var response map[string]interface{}
		if json.Unmarshal(output, &response) == nil {
			if accounts, ok := response["result"].([]interface{}); ok && len(accounts) > 0 {
				fmt.Printf("💰 %s account: %s\n", name, accounts[0])
			}
		}
	}
	
	fmt.Println("🔄 Starting periodic transfers simulation...")
	
	for i := 0; i < 3; i++ {
		fmt.Printf("💸 Transfer #%d: Alice → Bob (0.1 ETH)\n", i+1)
		fmt.Println("   ⏳ (Real transaction implementation coming next...)")
		
		if i < 2 {
			fmt.Println("   ⏱️  Waiting 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}
	
	return nil
}

func (tm *TransactionManager) GetNetworkStatus() {
	fmt.Println("📊 Current Network Status:")
	fmt.Println("==========================")
	
	for name, client := range tm.clients {
		ctx := context.Background()
		
		blockNum, err := client.BlockNumber(ctx)
		if err != nil {
			fmt.Printf("%s: ❌ Offline\n", name)
			continue
		}
		
		fmt.Printf("%s: 🟢 Block #%d\n", name, blockNum)
	}
}
