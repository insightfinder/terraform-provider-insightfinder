.PHONY: build install test test-unit test-acc testacc test-coverage fmt lint docs clean

default: build

# Build the provider
build:
	go build -o terraform-provider-insightfinder

# Install the provider locally for testing
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/insightfinder/insightfinder/1.0.0/linux_amd64
	cp terraform-provider-insightfinder ~/.terraform.d/plugins/registry.terraform.io/insightfinder/insightfinder/1.0.0/linux_amd64/

# Run all tests (unit + acceptance)
test:
	@echo "Running unit tests..."
	go test -v -short ./internal/provider ./internal/provider/client
	@echo "\nTo run acceptance tests, use: make testacc"

# Run unit tests only
test-unit:
	go test -v -short ./internal/provider ./internal/provider/client

# Run acceptance tests
test-acc:
	TF_ACC=1 go test -v ./internal/provider -timeout 120m

# Alias for test-acc (for backwards compatibility)
testacc: test-acc

# Run tests with coverage
test-coverage:
	go test -v -cover -coverprofile=coverage.out ./internal/provider ./internal/provider/client
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific resource tests
test-project:
	TF_ACC=1 go test -v ./internal/provider -run TestAccProjectResource -timeout 30m

test-servicenow:
	TF_ACC=1 go test -v ./internal/provider -run TestAccServiceNowResource -timeout 30m

test-jwt:
	TF_ACC=1 go test -v ./internal/provider -run TestAccJWTConfigResource -timeout 30m

test-loglabels:
	TF_ACC=1 go test -v ./internal/provider -run TestAccLogLabelsResource -timeout 30m

test-datasources:
	TF_ACC=1 go test -v ./internal/provider -run TestAccProjectDataSource -timeout 30m
	TF_ACC=1 go test -v ./internal/provider -run TestAccSystemsDataSource -timeout 30m

# Format code
fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

# Lint code
lint:
	golangci-lint run

# Generate documentation
docs:
	go generate

# Clean build artifacts
clean:
	rm -f terraform-provider-insightfinder
	rm -rf dist/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Initialize the provider for local development
dev-init: install
	cd examples/provider && terraform init

# Test the provider locally
dev-test: dev-init
	cd examples/provider && terraform plan
