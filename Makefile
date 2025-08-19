.PHONY: build run clean test install deps

# Variables
BINARY_NAME=benchy
BUILD_DIR=./bin

# Construction du projet
build:
	@echo "🔨 Building Benchy..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/benchy
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Installation des dépendances
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download

# Exécution
run:
	@go run ./cmd/benchy $(ARGS)

# Installation globale
install:
	@echo "📦 Installing Benchy globally..."
	@go install ./cmd/benchy

# Tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Nettoyage
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Afficher l'aide
help:
	@echo "Benchy Makefile Commands:"
	@echo "  make deps     - Install dependencies"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Run with 'go run'"
	@echo "  make install  - Install globally"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build files"
