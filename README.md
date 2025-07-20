# Bike Parts Finder

A cloud-native application for finding bicycle parts by brand, model, and year.

## Overview

Bike Parts Finder helps cyclists find compatible parts for their bicycles by searching across multiple online retailers. It uses web scraping to gather real-time parts data and presents structured results with images when available.

## Architecture

```
User
|
React Frontend <--> Go API Backend
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

- **Frontend**: React.js (built and managed with Node.js tooling)
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

- [Go](https://golang.org/) (1.21+) - For backend services
- [Node.js](https://nodejs.org/) (18+) - For React frontend development
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
   npm install
   npm start
   ```

## Project Structure

```
.
├── .github/           # GitHub Actions workflows
├── cmd/               # Go application entrypoints
│   ├── api/           # Backend API service
│   ├── scraper/       # Web scraper service
│   └── consumer/      # Kafka consumer service
├── pkg/               # Shared Go packages
│   ├── models/        # Data models
│   ├── database/      # Database access
│   ├── cache/         # Redis cache utilities
│   ├── kafka/         # Kafka utilities
│   └── scraping/      # Web scraping logic
├── web/               # Frontend application
│   └── frontend/      # React application (uses Node.js toolchain)
├── infra/             # Infrastructure as code
│   ├── terraform/     # OpenTofu configuration
│   ├── helmfile/      # Helmfile configuration for Kubernetes resources
│   ├── kubernetes/    # Raw Kubernetes manifests
│   └── argocd/        # ArgoCD configuration
└── Taskfile.yml       # Task definitions
```

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
