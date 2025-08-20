.PHONY: help build test clean docker deploy docs lint security

# Variables
BINARY_NAME=vm-chan
DOCKER_IMAGE=vm-chan
VERSION?=latest
COVERAGE_FILE=coverage.out

# Default target
help: ## Display this help message
	@echo "VM-Chan Microservice"
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
build: ## Build the application for the local OS
	@echo "Building $(BINARY_NAME) for local execution..."
	@go mod tidy
	@CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) ./cmd/server

build-linux: ## Build the application for Linux (for Docker)
	@echo "Building $(BINARY_NAME) for Linux..."
	@go mod tidy
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(BINARY_NAME) ./cmd/server

run: build ## Run the application locally
	@echo "Starting $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf coverage.out coverage.html
	@go clean

# Testing
test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v -tags=integration ./test/...

# Code Quality
lint: ## Run linting
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2"; \
	fi

security: ## Run security scan
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Run: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Documentation
docs: ## Generate API documentation
	@echo "Generating API documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/server/main.go -o ./api; \
	else \
		echo "swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) .
	@docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

docker-run: docker-build ## Run Docker container
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 $(DOCKER_IMAGE):latest

docker-push: docker-build ## Push Docker image
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_IMAGE):$(VERSION)
	@docker push $(DOCKER_IMAGE):latest

# Infrastructure
terraform-init: ## Initialize Terraform
	@echo "Initializing Terraform..."
	@cd deployments/terraform && terraform init

terraform-plan: ## Plan Terraform deployment
	@echo "Planning Terraform deployment..."
	@cd deployments/terraform && terraform plan

terraform-apply: ## Apply Terraform configuration
	@echo "Applying Terraform configuration..."
	@cd deployments/terraform && terraform apply

terraform-destroy: ## Destroy Terraform infrastructure
	@echo "Destroying Terraform infrastructure..."
	@cd deployments/terraform && terraform destroy

# Kubernetes
k8s-deploy: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	@kubectl apply -f deployments/k8s/

k8s-status: ## Check Kubernetes deployment status
	@echo "Checking deployment status..."
	@kubectl get pods -n vm-chan
	@kubectl get services -n vm-chan

k8s-logs: ## Get application logs
	@echo "Getting application logs..."
	@kubectl logs -f deployment/vm-chan -n vm-chan

k8s-delete: ## Delete Kubernetes deployment
	@echo "Deleting Kubernetes deployment..."
	@kubectl delete -f deployments/k8s/

# Complete workflows
ci: lint security test ## Run CI pipeline locally
	@echo "CI pipeline completed successfully!"

cd: docker-build k8s-deploy ## Run CD pipeline locally
	@echo "CD pipeline completed successfully!"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed successfully!"
