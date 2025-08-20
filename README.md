# VM-Chan - Text Analysis Microservice

A production-ready Go microservice that analyzes text sentences and returns word count, vowel count, and consonant count statistics.

## Architecture

This project follows the Clean Architecture pattern with dependency injection and follows Go best practices:

```
vm-chan/
├── cmd/                    # Application entry points
│   └── server/
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── domain/           # Business entities and interfaces
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── service/          # Business logic
│   └── repository/       # Data access layer
├── api/                  # API documentation
├── configs/              # Configuration files
├── docs/                 # Generated documentation
├── deployments/          # Deployment configurations
│   ├── terraform/        # Infrastructure as Code
│   ├── ansible/          # Configuration management
│   └── k8s/             # Kubernetes manifests
├── scripts/              # Build and utility scripts
└── .github/workflows/    # CI/CD pipelines
```

## Features

- RESTful API with OpenAPI/Swagger documentation
- JWT Authentication & Authorization
- Prometheus metrics integration
- Structured logging with Zap
- Configuration management with Viper
- Unit and integration tests
- Docker containerization
- Kubernetes deployment manifests
- Terraform infrastructure provisioning
- Ansible configuration management
- GitHub Actions CI/CD pipeline
- Health checks and graceful shutdown
- CORS and security middleware

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics endpoint

### Authentication
- `POST /auth/login` - Get JWT token

### Text Analysis
- `POST /api/v1/analyze` - Analyze text sentence (requires authentication)

### Documentation
- `GET /swagger/*any` - Interactive API documentation

## Quick Start

### Local Development

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Run tests: `make test`
4. Start the server: `make run`
5. Access Swagger UI: `http://localhost:8080/swagger/index.html`

### Using the API

1. **Login to get JWT token:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

2. **Analyze text with the token:**
```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"sentence": "Hello world!"}'
```

### Docker

```bash
# Build image
make docker-build

# Run container
make docker-run
```

### Kubernetes Deployment

```bash
# Deploy to k3s
make k8s-deploy

# Check status
make k8s-status

# View logs
make k8s-logs
```

## Infrastructure Deployment

### Prerequisites

- AWS CLI configured with appropriate credentials
- Terraform >= 1.7.0
- Ansible >= 2.9
- kubectl configured

### Automated Deployment

```bash
# Set required environment variables
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_KEY_NAME="your-key-name"
export AWS_REGION="us-east-1"

# Deploy everything
./scripts/deploy.sh
```

### Manual Infrastructure Steps

```bash
# Provision EC2 instance with k3s
cd deployments/terraform
terraform init
terraform plan
terraform apply

# Configure k3s with Ansible
cd ../ansible
ansible-playbook -i inventory.ini site.yml

# Deploy application
kubectl apply -f ../k8s/
```

### Cleanup

```bash
# Destroy all infrastructure
./scripts/cleanup.sh
```

## Configuration

The application uses environment variables and YAML configuration:

### Environment Variables
- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: JWT signing secret
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `METRICS_ENABLED`: Enable Prometheus metrics (default: true)

### Configuration File
See `configs/config.yaml` for default configuration values.

## Development

### Available Make Targets

```bash
make help           # Show all available targets
make build          # Build the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Run linting
make security       # Run security scan
make docs           # Generate API documentation
make docker-build   # Build Docker image
make ci             # Run CI pipeline locally
```

## Monitoring

- Prometheus metrics available at `/metrics`
- Health check at `/health`
- Structured JSON logging with configurable levels
- Request/response logging middleware

## Security

- JWT-based authentication with configurable secret
- Input validation and sanitization
- CORS middleware with configurable origins
- Security headers middleware (XSS protection, CSRF, etc.)
- Secure password hashing with bcrypt

## Testing

The project includes comprehensive testing:

- Unit tests for all service layers
- Integration tests for API endpoints
- Test coverage reporting
- Mocked dependencies for isolated testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run security scan
make security

# Run linting
make lint
```

## CI/CD Pipeline

The project includes a complete GitHub Actions pipeline that:

1. Runs tests, linting, and security scans
2. Builds and pushes Docker images to GHCR
3. Provisions AWS infrastructure with Terraform
4. Configures k3s cluster with Ansible
5. Deploys the application to Kubernetes
6. Provides access URLs and management commands

## Project Structure Details

- **cmd/server**: Application entry point and main function
- **internal/config**: Configuration loading and management
- **internal/domain**: Business entities, interfaces, and DTOs
- **internal/handler**: HTTP request handlers
- **internal/middleware**: HTTP middleware (auth, logging, metrics, CORS, security)
- **internal/service**: Business logic implementation
- **internal/repository**: Data access layer (in-memory user store)
- **deployments/terraform**: AWS infrastructure provisioning
- **deployments/ansible**: k3s cluster configuration
- **deployments/k8s**: Kubernetes deployment manifests
- **scripts**: Build, deployment, and cleanup automation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.
