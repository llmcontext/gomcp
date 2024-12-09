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

# Build gomcp-proxy
build:
	@echo "Building gomcp-proxy..."
	@go build -o bin/gomcp-proxy cmd/gomcp-proxy/main.go

.PHONY: test test-coverage fmt vet deps build