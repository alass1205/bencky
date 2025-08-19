package monitor

import (
	"fmt"
	"time"
)

func (nm *NetworkMonitor) DisplayNetworkInfoContinuous(updateInterval int) error {
	if updateInterval <= 0 {
		updateInterval = 60 // Default 60 seconds
	}
	
	fmt.Printf("🔄 Starting continuous monitoring (update every %d seconds)...\n", updateInterval)
	fmt.Println("Press Ctrl+C to stop")
	
	for {
		// Clear screen (simple)
		fmt.Print("\033[2J\033[H")
		
		// Display timestamp
		fmt.Printf("🕐 Last update: %s\n\n", time.Now().Format("15:04:05"))
		
		// Display network info
		if err := nm.DisplayNetworkInfo(); err != nil {
			return err
		}
		
		// Wait for next update
		time.Sleep(time.Duration(updateInterval) * time.Second)
	}
}
