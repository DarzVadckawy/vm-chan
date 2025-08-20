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
├── pkg/                  # Public library code
├── api/                  # API documentation
├── deployments/          # Deployment configurations
│   ├── terraform/        # Infrastructure as Code
│   ├── ansible/          # Configuration management
│   └── k8s/             # Kubernetes manifests
├── scripts/              # Build and utility scripts
├── test/                 # Additional test files
└── .github/workflows/    # CI/CD pipelines
```

## Features

- ✅ RESTful API with OpenAPI/Swagger documentation
- ✅ JWT Authentication & Authorization
- ✅ Prometheus metrics integration
- ✅ Structured logging with Zap
- ✅ Configuration management with Viper
- ✅ Unit and integration tests
- ✅ Docker containerization
- ✅ Kubernetes deployment manifests
- ✅ Terraform infrastructure provisioning
- ✅ Ansible configuration management
- ✅ GitHub Actions CI/CD pipeline
- ✅ SonarQube integration
- ✅ Health checks and graceful shutdown

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Authentication
- `POST /auth/login` - Get JWT token

### Text Analysis
- `POST /api/v1/analyze` - Analyze text sentence (requires authentication)

## Quick Start

### Local Development

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Run tests: `make test`
4. Start the server: `make run`
5. Access Swagger UI: `http://localhost:8080/swagger/index.html`

### Docker

```bash
# Build image
docker build -t vm-chan:latest .

# Run container
docker run -p 8080:8080 vm-chan:latest
```

### Kubernetes Deployment

```bash
# Deploy to k3s
kubectl apply -f deployments/k8s/
```

## Infrastructure Deployment

### Prerequisites

- AWS CLI configured
- Terraform >= 1.0
- Ansible >= 2.9

### Deploy Infrastructure

```bash
# Provision EC2 instance
cd deployments/terraform
terraform init
terraform plan
terraform apply

# Configure k3s
cd ../ansible
ansible-playbook -i inventory.yml site.yml
```

## Configuration

The application uses environment variables for configuration:

- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: JWT signing secret
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `METRICS_ENABLED`: Enable Prometheus metrics (default: true)

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration
```

## Monitoring

- Prometheus metrics available at `/metrics`
- Health check at `/health`
- Structured JSON logging

## Security

- JWT-based authentication
- Input validation and sanitization
- Rate limiting middleware
- Secure headers middleware

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.
