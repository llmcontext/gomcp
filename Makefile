# Binary name
BINARY := gomcp

# Build directory
BUILD_DIR=./bin

.PHONY: build test-coverage fmt vet deps

all: build

# Build all binaries
build:
	@echo "Building $(BINARY)..."
	@go build -o $(BUILD_DIR)/$(BINARY) cmd/*.go; \

# Install binaries to /usr/local/bin
install: build
	@echo "Installing $(BINARY) to /usr/local/bin..."
	@cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/; \
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

inspector: build
	@echo "Running inspector..."
	npx @modelcontextprotocol/inspector $(BUILD_DIR)/$(BINARY) --debug
