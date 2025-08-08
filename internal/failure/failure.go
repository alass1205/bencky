package failure

import (
	"fmt"
	"os/exec"
	"time"
)

func TemporaryFailure(nodeName string) error {
	containerName := fmt.Sprintf("benchy-%s", nodeName)
	
	fmt.Printf("💥 Stopping container %s...\n", containerName)
	
	// Stop the container
	stopCmd := exec.Command("docker", "stop", containerName)
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}
	
	fmt.Println("🕒 Container stopped. Waiting 40 seconds...")
	
	// Wait 40 seconds
	for i := 40; i > 0; i-- {
		fmt.Printf("   ⏱️  Restarting in %d seconds...\n", i)
		time.Sleep(1 * time.Second)
	}
	
	fmt.Printf("🔄 Restarting container %s...\n", containerName)
	
	// Restart the container
	startCmd := exec.Command("docker", "start", containerName)
	if err := startCmd.Run(); err != nil {
		return fmt.Errorf("failed to restart container: %v", err)
	}
	
	fmt.Println("✅ Container restarted successfully!")
	fmt.Println("📊 Use 'benchy infos' to check if node is back online")
	
	return nil
}
