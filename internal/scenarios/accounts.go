package scenarios

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func GetRealAccounts() error {
	fmt.Println("🔍 Getting REAL accounts from each node...")
	fmt.Println("=" + strings.Repeat("=", 60))

	endpoints := map[string]string{
		"Alice":     "http://localhost:8545",
		"Bob":       "http://localhost:8547",
		"Cassandra": "http://localhost:8549",
		"Driss":     "http://localhost:8551",
		"Elena":     "http://localhost:8553",
	}

	for name, endpoint := range endpoints {
		fmt.Printf("%-12s: ", name)
		
		// Obtenir les comptes
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
		if json.Unmarshal(output, &response) != nil {
			fmt.Printf("❌ Parse error\n")
			continue
		}

		if accounts, ok := response["result"].([]interface{}); ok && len(accounts) > 0 {
			address := accounts[0].(string)
			
			// Obtenir le solde
			balanceCmd := exec.Command("curl", "-s", "-X", "POST",
				"-H", "Content-Type: application/json",
				"--data", fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBalance","params":["%s","latest"],"id":1}`, address),
				endpoint)
			
			balanceOutput, _ := balanceCmd.Output()
			var balanceResponse map[string]interface{}
			
			balance := "0x0"
			if json.Unmarshal(balanceOutput, &balanceResponse) == nil {
				if bal, ok := balanceResponse["result"].(string); ok {
					balance = bal
				}
			}
			
			fmt.Printf("🟢 %s (Balance: %s)\n", address, balance)
		} else {
			fmt.Printf("❌ No accounts\n")
		}
	}
	
	fmt.Println("=" + strings.Repeat("=", 60))
	return nil
}
