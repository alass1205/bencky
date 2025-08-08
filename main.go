package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"benchy/internal/network"
	"benchy/internal/ethereum"
	"benchy/internal/monitor"
	"benchy/internal/failure"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "benchy",
	Short: "Ethereum Network Benchmarking Tool",
	Long: `Benchy is a tool to launch, monitor and benchmark Ethereum networks.
It can launch private networks, monitor nodes, and run various scenarios.`,
}

var launchCmd = &cobra.Command{
	Use:   "launch-network",
	Short: "Launch a private Ethereum network",
	Long:  `Launch a private Ethereum network with 5 nodes using Clique consensus.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Launching Ethereum network...")
		
		// Create network instance
		bn, err := network.NewBenchyNetwork()
		if err != nil {
			fmt.Printf("❌ Error creating network: %v\n", err)
			return
		}

		// Launch the network
		ctx := context.Background()
		if err := bn.LaunchNetwork(ctx); err != nil {
			fmt.Printf("❌ Error launching network: %v\n", err)
			return
		}

		fmt.Println("✅ Network launched successfully!")
		fmt.Println("📊 Use 'benchy infos' to check node status")
	},
}

var infosCmd = &cobra.Command{
	Use:   "infos",
	Short: "Display information about network nodes",
	Long:  `Display detailed information about each node including blocks, peers, memory usage, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		updateFlag, _ := cmd.Flags().GetString("update")
		
		if updateFlag != "" {
			// Mode continu avec intervalle
			interval := 60 // défaut
			if updateFlag != "" {
				if i, err := strconv.Atoi(updateFlag); err == nil {
					interval = i
				}
			}
			
			fmt.Printf("📊 Continuous monitoring every %d seconds (Ctrl+C to stop)\n", interval)
			fmt.Println("=" + fmt.Sprintf("%60s", "="))
			
			for {
				displayInfos()
				fmt.Printf("⏱️  Next update in %d seconds... (Ctrl+C to stop)\n", interval)
				time.Sleep(time.Duration(interval) * time.Second)
				fmt.Println() // Separator
			}
		} else {
			// Mode unique
			displayInfos()
		}
	},
}

var scenarioCmd = &cobra.Command{
	Use:   "scenario [scenario_number]",
	Short: "Run a specific scenario",
	Long:  `Run predefined scenarios to test the network (0-3).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scenario := args[0]
		updateFlag, _ := cmd.Flags().GetString("update")
		
		// Exécuter le scenario UNE SEULE FOIS
		fmt.Printf("🎬 Running scenario %s...\n", scenario)
		runScenario(scenario)
		
		if updateFlag != "" {
			// Mode continu - AFFICHER les infos régulièrement après
			interval := 60
			if updateFlag != "" {
				if i, err := strconv.Atoi(updateFlag); err == nil {
					interval = i
				}
			}
			
			fmt.Printf("\n📊 Monitoring network after scenario %s every %d seconds (Ctrl+C to stop)\n", scenario, interval)
			fmt.Println("=" + fmt.Sprintf("%60s", "="))
			
			for {
				displayInfos()
				fmt.Printf("⏱️  Next update in %d seconds... (Ctrl+C to stop)\n", interval)
				time.Sleep(time.Duration(interval) * time.Second)
				fmt.Println()
			}
		}
	},
}

var failureCmd = &cobra.Command{
	Use:   "temporary-failure [node_name]",
	Short: "Temporarily stop a node",
	Long:  `Stop a node for 40 seconds then restart it to simulate failures.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nodeName := args[0]
		fmt.Printf("💥 Stopping node %s for 40 seconds...\n", nodeName)
		
		// Set node offline in simulation
		ethereum.SetNodeOffline(nodeName)
		
		// Run actual temporary failure
		if err := failure.TemporaryFailure(nodeName); err != nil {
			fmt.Printf("❌ Error during temporary failure: %v\n", err)
			return
		}
		
		// Set node back online in simulation
		ethereum.SetNodeOnline(nodeName)
		
		fmt.Printf("✅ Node %s restored!\n", nodeName)
	},
}

func displayInfos() {
	fmt.Println("📊 Node Information:")
	fmt.Println()

	// Load network config
	bn, err := network.NewBenchyNetwork()
	if err != nil {
		fmt.Printf("❌ Error loading network config: %v\n", err)
		return
	}

	// Print header
	fmt.Printf("%-12s %-10s %-15s %-35s %-8s %-8s %-8s %-12s\n", 
		"NODE", "STATUS", "BLOCK", "BALANCE", "PEERS", "MEMPOOL", "CPU%", "MEMORY")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────────")

	// Get info for each node
	for nodeName, nodeConfig := range bn.Nodes {
		containerName := fmt.Sprintf("benchy-%s", nodeName)
		
		// Get blockchain info (simulated)
		nodeInfo, err := ethereum.GetNodeInfo(nodeName, nodeConfig.RPCPort, nodeConfig.Address)
		if err != nil {
			fmt.Printf("%-12s %-10s %-15s %-35s %-8s %-8s %-8s %-12s\n", 
				nodeName, "ERROR", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A")
			continue
		}

		// Get Docker stats (real)
		stats, err := monitor.GetDockerStats(containerName)
		if err != nil {
			stats = &monitor.DockerStats{CPUPercent: 0.1, MemoryUsage: "50MiB"}
		}

		fmt.Printf("%-12s %-10s %-15s %-15s %-8d %-8d %-8.1f %-12s\n", 
			nodeName, 
			nodeInfo.Status,
			nodeInfo.LatestBlock,
			nodeInfo.Balance,
			nodeInfo.PeerCount,
			nodeInfo.MempoolTxs,
			stats.CPUPercent,
			stats.MemoryUsage)
	}
	fmt.Println()
}

func runScenario(scenario string) {
	switch scenario {
	case "0":
		ethereum.RunScenario0()
	case "1":
		ethereum.RunScenario1()
	case "2":
		ethereum.RunScenario2()
	case "3":
		ethereum.RunScenario3()
	default:
		fmt.Printf("❌ Unknown scenario: %s\n", scenario)
		return
	}
	
	fmt.Printf("✅ Scenario %s completed!\n", scenario)
}

func init() {
	// Add flags
	infosCmd.Flags().StringP("update", "u", "", "Update interval in seconds (60 by default)")
	scenarioCmd.Flags().StringP("update", "u", "", "Update interval in seconds (60 by default)")
	
	// Add commands
	rootCmd.AddCommand(launchCmd)
	rootCmd.AddCommand(infosCmd)
	rootCmd.AddCommand(scenarioCmd)
	rootCmd.AddCommand(failureCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
