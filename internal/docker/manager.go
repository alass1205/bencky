package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type DockerManager struct {
	composeDir string
}

func NewDockerManager() (*DockerManager, error) {
	pwd, _ := os.Getwd()
	composeDir := filepath.Join(pwd, "docker")
	return &DockerManager{composeDir: composeDir}, nil
}

func (dm *DockerManager) LaunchNetwork() error {
	fmt.Println("🚀 Launching REAL Ethereum network with Docker...")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(dm.composeDir); err != nil {
		return fmt.Errorf("failed to change to docker directory: %v", err)
	}

	fmt.Println("🧹 Cleaning up existing containers...")
	cmd := exec.Command("docker-compose", "down", "-v")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: cleanup failed (this is normal if first run): %v\n", err)
	}

	fmt.Println("🔄 Starting network containers...")
	cmd = exec.Command("docker-compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start network: %v", err)
	}

	fmt.Println("⏳ Waiting for nodes to initialize...")
	time.Sleep(15 * time.Second)

	fmt.Println("✅ Network launched successfully!")
	fmt.Println("📍 Nodes accessible at:")
	fmt.Println("  - Alice (Geth):      http://localhost:8545")
	fmt.Println("  - Bob (Nethermind):  http://localhost:8547")
	fmt.Println("  - Cassandra (Geth):  http://localhost:8549")
	fmt.Println("  - Driss (Nethermind): http://localhost:8551")
	fmt.Println("  - Elena (Geth):      http://localhost:8553")

	return nil
}

func (dm *DockerManager) StopContainer(containerName string, duration int) error {
	fmt.Printf("⚠️  Stopping %s for %d seconds...\n", containerName, duration)
	
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(dm.composeDir); err != nil {
		return fmt.Errorf("failed to change to docker directory: %v", err)
	}

	cmd := exec.Command("docker-compose", "stop", containerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}

	fmt.Printf("⏳ Waiting %d seconds...\n", duration)
	time.Sleep(time.Duration(duration) * time.Second)

	fmt.Printf("🔄 Restarting %s...\n", containerName)
	cmd = exec.Command("docker-compose", "start", containerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart container: %v", err)
	}

	fmt.Printf("✅ %s is back online!\n", containerName)
	return nil
}
