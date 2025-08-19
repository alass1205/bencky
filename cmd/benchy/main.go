package main

import (
	"fmt"
	"os"

	"benchy/internal/docker"
	"benchy/internal/monitor"
	"benchy/internal/scenarios"

	"github.com/spf13/cobra"
)

var dockerManager *docker.DockerManager
var networkMonitor *monitor.NetworkMonitor
var transactionManager *scenarios.TransactionManager

var updateInterval int

var rootCmd = &cobra.Command{
	Use:   "benchy",
	Short: "Benchy - REAL Ethereum Network Benchmarking Tool",
	Long:  `A tool to launch, monitor and benchmark REAL Ethereum networks with multiple clients.`,
}

var launchCmd = &cobra.Command{
	Use:   "launch-network",
	Short: "Launch the REAL Ethereum network with 5 nodes",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dockerManager.LaunchNetwork(); err != nil {
			fmt.Printf("❌ Failed to launch network: %v\n", err)
			os.Exit(1)
		}
	},
}

var infosCmd = &cobra.Command{
	Use:   "infos",
	Short: "Display information about REAL network nodes",
	Run: func(cmd *cobra.Command, args []string) {
		if updateInterval > 0 {
			if err := networkMonitor.DisplayNetworkInfoContinuous(updateInterval); err != nil {
				fmt.Printf("❌ Failed to display continuous info: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := networkMonitor.DisplayNetworkInfo(); err != nil {
				fmt.Printf("❌ Failed to get network info: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

var scenarioCmd = &cobra.Command{
	Use:   "scenario [number]",
	Short: "Run predefined scenarios on REAL network",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scenario := args[0]
		fmt.Printf("🎬 Running scenario %s on REAL network...\n", scenario)
		
		var err error
		switch scenario {
		case "0":
			err = transactionManager.FullScenario0()
		case "1":
			err = transactionManager.FullScenario1()
		case "2":
			err = transactionManager.FullScenario2()
		case "3":
			err = transactionManager.FullScenario3()
		default:
			fmt.Printf("❌ Unknown scenario: %s\n", scenario)
			return
		}
		
		if err != nil {
			fmt.Printf("❌ Scenario failed: %v\n", err)
		}
	},
}

var failureCmd = &cobra.Command{
	Use:   "temporary-failure [node]",
	Short: "Simulate temporary node failure",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		node := args[0]
		if err := dockerManager.StopContainer(node, 40); err != nil {
			fmt.Printf("❌ Failed to simulate failure: %v\n", err)
			os.Exit(1)
		}
	},
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run realistic transaction demo",
	Run: func(cmd *cobra.Command, args []string) {
		if err := scenarios.RunRealisticDemo(); err != nil {
			fmt.Printf("❌ Demo failed: %v\n", err)
		}
	},
}

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Show REAL accounts and balances",
	Run: func(cmd *cobra.Command, args []string) {
		if err := scenarios.GetRealAccounts(); err != nil {
			fmt.Printf("❌ Failed to get accounts: %v\n", err)
		}
	},
}

func init() {
	var err error
	
	dockerManager, err = docker.NewDockerManager()
	if err != nil {
		fmt.Printf("❌ Failed to initialize Docker manager: %v\n", err)
		os.Exit(1)
	}

	networkMonitor = monitor.NewNetworkMonitor()
	transactionManager = scenarios.NewTransactionManager()

	rootCmd.PersistentFlags().IntVarP(&updateInterval, "update", "u", 0, "Update interval in seconds")

	rootCmd.AddCommand(launchCmd)
	rootCmd.AddCommand(infosCmd)
	rootCmd.AddCommand(scenarioCmd)
	rootCmd.AddCommand(failureCmd)
	rootCmd.AddCommand(demoCmd)
	rootCmd.AddCommand(accountsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
