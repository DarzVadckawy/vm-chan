.PHONY: help build test clean docker deploy docs lint security sonar sonar-local

BINARY_NAME=vm-chan
DOCKER_IMAGE=vm-chan
VERSION?=latest
COVERAGE_FILE=coverage.out
SONAR_SCANNER_VERSION=4.8.0.2856

help:
	@echo "VM-Chan Microservice"
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build:
	@echo "Building $(BINARY_NAME) for local execution..."
	@go mod tidy
	@CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) ./cmd/server

build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@go mod tidy

run: build
	@echo "Starting $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf coverage.out coverage.html
	@go clean

test:
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./test/...

lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2"; \
	fi

security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Run: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

docs:
	@echo "Generating API documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/server/main.go -o ./api; \
	else \
		echo "swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) .
	@docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

docker-run: docker-build
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 $(DOCKER_IMAGE):latest

docker-push: docker-build
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_IMAGE):$(VERSION)
	@docker push $(DOCKER_IMAGE):latest

terraform-init:
	@echo "Initializing Terraform..."
	@cd deployments/terraform && terraform init

terraform-plan:
	@echo "Planning Terraform deployment..."
	@cd deployments/terraform && terraform plan

terraform-apply:
	@echo "Applying Terraform configuration..."
	@cd deployments/terraform && terraform apply

terraform-destroy:
	@echo "Destroying Terraform infrastructure..."
	@cd deployments/terraform && terraform destroy

terraform-output:
	@echo "Getting Terraform outputs..."
	@cd deployments/terraform && terraform output

terraform-ip:
	@echo "Getting EC2 instance public IP..."
	@cd deployments/terraform && terraform output -raw public_ip 2>/dev/null || echo "No Terraform state found. Run 'make terraform-apply' first."

deploy-remote:
	@echo "Deploying to remote AWS infrastructure..."
	@echo "Use one of the following methods:"
	@echo "  1. Push to main branch (triggers GitHub Actions CI/CD)"
	@echo "  2. Run ./scripts/deploy.sh for manual deployment"
	@echo "  3. Use GitHub Actions manually via GitHub UI"

ci: lint security test
	@echo "CI pipeline completed successfully!"

install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed successfully!"
