#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PROJECT_NAME="vm-chan"
IMAGE_NAME="vm-chan"
VERSION=${1:-latest}

echo -e "${GREEN}VM-Chan Build Script${NC}"
echo "Version: $VERSION"

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

build_app() {
    print_status "Building Go application..."
    go mod tidy
    go build -o bin/vm-chan ./cmd/server
    print_status "Application built successfully"
}

run_tests() {
    print_status "Running tests..."
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    print_status "Tests completed"
}

build_docker() {
    print_status "Building Docker image..."
    docker build -t $IMAGE_NAME:$VERSION .
    docker tag $IMAGE_NAME:$VERSION $IMAGE_NAME:latest
    print_status "Docker image built: $IMAGE_NAME:$VERSION"
}

generate_docs() {
    print_status "Generating API documentation..."
    if command -v swag &> /dev/null; then
        swag init -g cmd/server/main.go -o ./api
        print_status "Swagger documentation generated"
    else
        print_warning "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"
    fi
}

deploy_k3s() {
    print_status "Deploying to k3s..."

    # Check if kubectl is configured and cluster is accessible
    if ! kubectl cluster-info &>/dev/null; then
        print_warning "No Kubernetes cluster configured or accessible"
        print_warning "Skipping kubectl deployment. Use remote deployment via scripts/deploy.sh instead"
        return 0
    fi

    # Check if k3s is available locally
    if command -v k3s &> /dev/null; then
        docker save $IMAGE_NAME:latest | sudo k3s ctr images import -
    else
        print_warning "k3s not found locally, skipping image import"
    fi

    # Apply manifests with validation disabled
    kubectl apply -f deployments/k8s/ --validate=false

    # Wait for deployment only if we successfully applied
    if kubectl get deployment vm-chan -n vmchan &>/dev/null; then
        kubectl wait --for=condition=available --timeout=300s deployment/vm-chan -n vmchan
    else
        print_warning "Deployment not found, skipping wait"
    fi

    print_status "Deployment completed"
}

case "${1:-all}" in
    "test")
        run_tests
        ;;
    "build")
        build_app
        ;;
    "docker")
        build_docker
        ;;
    "docs")
        generate_docs
        ;;
    "deploy")
        deploy_k3s
        ;;
    "all")
        run_tests
        build_app
        generate_docs
        build_docker
        print_status "Build process completed successfully!"
        ;;
    *)
        echo "Usage: $0 {test|build|docker|docs|deploy|all}"
        echo "  test   - Run tests"
        echo "  build  - Build Go application"
        echo "  docker - Build Docker image"
        echo "  docs   - Generate API documentation"
        echo "  deploy - Deploy to k3s"
        echo "  all    - Run all build steps"
        exit 1
        ;;
esac
