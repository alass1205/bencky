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

func (dm *DockerManager) CleanNetwork() error {
	fmt.Println("🧹 Cleaning up existing containers and persistent state...")
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(dm.composeDir); err != nil {
		return fmt.Errorf("failed to change to docker directory: %v", err)
	}

	// Nettoyer Docker
	cmd := exec.Command("docker-compose", "down", "-v")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: cleanup failed (this is normal if first run): %v\n", err)
	}

	// Supprimer le fichier d'état
	stateFile := filepath.Join(originalDir, "benchy_state.json")
	if err := os.Remove(stateFile); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: failed to remove state file: %v\n", err)
	} else if err == nil {
		fmt.Println("🗑️  Removed persistent state file")
	}

	return nil
}

func (dm *DockerManager) LaunchNetwork() error {
	fmt.Println("🚀 Launching REAL Ethereum network with Docker...")
	
	// Nettoyer complètement avant de lancer
	if err := dm.CleanNetwork(); err != nil {
		return err
	}
	
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(dm.composeDir); err != nil {
		return fmt.Errorf("failed to change to docker directory: %v", err)
	}

	fmt.Println("🔄 Starting network containers...")
	cmd := exec.Command("docker-compose", "up", "-d")
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

	// Stop container
	cmd := exec.Command("docker-compose", "stop", containerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}

	// Countdown visuel avec monitoring en parallèle
	fmt.Printf("📊 Monitor with 'benchy infos' in another terminal to see %s as 🔴 OFF\n", containerName)
	for i := duration; i > 0; i-- {
		fmt.Printf("\r⏳ Restarting in %d seconds...", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Print("\n")

	// Restart container
	fmt.Printf("🔄 Restarting %s...\n", containerName)
	cmd = exec.Command("docker-compose", "start", containerName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart container: %v", err)
	}

	fmt.Printf("✅ %s is back online! Run 'benchy infos' to confirm.\n", containerName)
	return nil
}