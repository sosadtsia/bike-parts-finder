# Bike Parts Finder

A cloud-native application for finding bicycle parts by brand, model, and year.

## Overview

Bike Parts Finder helps cyclists find compatible parts for their bicycles by searching across multiple online retailers. It uses web scraping to gather real-time parts data and presents structured results with images when available.

## Architecture

```
User
|
React.js Frontend <--> Go API Backend
|
+-----------+------------+
|            |
PostgreSQL    Redis Cache
^            |
|            |
Kafka Consumer <--- Kafka Topic: scrape_results
^
|
Kafka on EKS
^
|
Scraper (Go) <--- Kafka Topic: scrape_requests
|
External websites (e.g., JensonUSA)
```

## Technology Stack

- **Frontend**: Go WebAssembly with TailwindCSS
- **Backend API**: Go
- **Database**: PostgreSQL on EKS
- **Cache**: Redis on EKS
- **Message Broker**: Apache Kafka on EKS
- **Container Orchestration**: Kubernetes (Amazon EKS)
- **Infrastructure as Code**: OpenTofu (Terraform compatible)
- **Kubernetes Package Management**: Helmfile (declarative Helm chart management)
- **CI/CD & GitOps**: GitHub Actions, ArgoCD
- **Monitoring**: Prometheus, Grafana
- **Backup**: Velero

## Getting Started

### Prerequisites

- [Go](https://golang.org/) (1.22+) - For backend services and WebAssembly frontend
- [Docker](https://www.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [helm](https://helm.sh/) - Required by Helmfile
- [helmfile](https://helmfile.readthedocs.io/) - For managing Helm releases
- [Task](https://taskfile.dev/)
- [OpenTofu](https://opentofu.org/)
- [AWS CLI](https://aws.amazon.com/cli/) (configured)

### Development Environment Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/svosadtsia/bike-parts-finder.git
   cd bike-parts-finder
   ```

2. Deploy local infrastructure (uses local state):
   ```bash
   task bootstrap:local ENV=develop
   ```

3. Port forward ArgoCD UI:
   ```bash
   task argocd:portforward
   ```

   Access ArgoCD at https://localhost:8080 (get password with `task argocd:password`)

4. Start local development servers:
   ```bash
   # Backend API (in one terminal)
   cd cmd/api
   go run main.go

   # Frontend (in another terminal)
   cd web/frontend
   go run -tags js,wasm main.go  # Compile WebAssembly

   # Alternatively, you can use the Makefile
   cd web/frontend
   make build
   ```

5. Access the application at http://localhost:8080

## Project Structure

```
.
├── .github/           # GitHub Actions workflows
├── cmd/               # Go application entrypoints
│   ├── api/           # Backend API service
│   ├── scraper/       # Web scraper service
│   └── consumer/      # Kafka consumer service
├── docs/              # Documentation
│   ├── api.md         # API documentation
│   └── development.md # Development guide
├── pkg/               # Shared Go packages
│   ├── models/        # Data models
│   ├── database/      # Database access
│   ├── cache/         # Redis cache utilities
│   ├── kafka/         # Kafka utilities
│   └── scraping/      # Web scraping logic
├── web/               # Frontend application
│   └── frontend/      # Go WebAssembly application
├── infra/             # Infrastructure as code
│   ├── terraform/     # OpenTofu configuration
│   ├── helmfile/      # Helmfile configuration for Kubernetes resources
│   ├── kubernetes/    # Raw Kubernetes manifests
│   └── argocd/        # ArgoCD configuration
└── Taskfile.yml       # Task definitions
```

## API Documentation

The API includes versioning to ensure backward compatibility. All endpoints are available under:

```
/api/v1/...
```

Health check endpoints are also available:

```
/health      # Basic health check
/health/ready # Readiness check (verifies database and cache connections)
```

See [API Documentation](./docs/api.md) for details on all available endpoints.

## Infrastructure

The infrastructure is deployed on AWS EKS using a GitOps approach with ArgoCD. Kubernetes resources are managed declaratively using Helmfile. See [Infrastructure README](./infra/README.md) for details.

## Deployment

### CI/CD Pipeline

Deployments are automated through GitHub Actions:

1. Infrastructure is provisioned with OpenTofu
2. ArgoCD is bootstrapped on the cluster
3. Helmfile is used to manage Helm chart releases
4. ArgoCD manages application deployments using GitOps

For manual deployments with remote state:
```bash
task bootstrap:remote ENV=production REGION=us-east-2
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
