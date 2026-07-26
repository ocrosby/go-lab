.PHONY: help build test test-unit test-integration test-coverage clean docker-build docker-test swagger lint

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build targets
build: ## Build the application binary
	go build -o bin/api cmd/api/main.go

swagger: ## Generate Swagger documentation
	swag init -g cmd/api/main.go

# Test targets
test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	go test -v -race -short ./internal/... ./pkg/...

test-integration: ## Run integration tests
	go test -v -race -run TestIntegration ./...

test-coverage: ## Run tests with coverage report
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

test-coverage-cli: ## Show coverage summary in terminal
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Quality targets
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy go modules
	go mod tidy

# Docker targets
docker-build: ## Build Docker image
	docker build -t api:latest .

docker-test: ## Run tests in Docker
	docker build -f Dockerfile.test --target test-runner -t api:test .
	docker run --rm api:test

docker-coverage: ## Generate coverage report in Docker
	docker build -f Dockerfile.test --target coverage -t api:coverage .
	docker run --rm api:coverage

# Development targets
dev: swagger ## Run in development mode
	go run cmd/api/main.go server

dev-watch: ## Run with file watching (requires air)
	air

# Clean targets
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker image prune -f

# CI targets
ci: lint test-coverage ## Run CI pipeline locally

# Generate mock files
generate-mocks: ## Generate mock files
	go generate ./...

# Database targets (for future use)
migrate-up: ## Run database migrations up
	@echo "No migrations defined yet"

migrate-down: ## Run database migrations down
	@echo "No migrations defined yet"

# Quick setup for new developers
setup: ## Setup development environment
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	$(MAKE) swagger

# Health check
health: ## Check if application dependencies are available
	@echo "Checking Go version..."
	@go version
	@echo "Checking Docker..."
	@docker --version
	@echo "Checking make..."
	@make --version
	@echo "All dependencies are available!"