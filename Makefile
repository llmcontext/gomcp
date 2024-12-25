# Define all binary names
BINARIES := gomcp gomcp-proxy

# Build directory
BUILD_DIR=./bin

.PHONY: build test-coverage fmt vet deps

all: build

# Build all binaries
build:
	@echo "Building binaries..."
	@for binary in $(BINARIES); do \
		echo "Building $$binary..."; \
		go build -o $(BUILD_DIR)/$$binary cmd/$$binary/*.go; \
	done

# Install binaries to /usr/local/bin
install: build
	@echo "Installing binaries to /usr/local/bin..."
	@for binary in $(BINARIES); do \
		cp $(BUILD_DIR)/$$binary /usr/local/bin/; \
	done
	@echo "Installation complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# staticcheck
staticcheck:
	@echo "Running staticcheck..."
	@staticcheck ./...
