# Bike Parts Finder Infrastructure

This directory contains the infrastructure-as-code configuration for the Bike Parts Finder application.

## Architecture

The infrastructure is deployed on AWS EKS using a GitOps approach with ArgoCD. The main components include:

- PostgreSQL (on EKS)
- Redis (on EKS)
- Kafka (on EKS)
- ArgoCD (on EKS)
- Velero for backups
- Prometheus and Grafana for monitoring

## Prerequisites

- [Go Task](https://taskfile.dev/) - Task runner
- [OpenTofu](https://opentofu.org/) - Infrastructure as Code
- [kubectl](https://kubernetes.io/docs/tasks/tools/) - Kubernetes CLI
- [helm](https://helm.sh/) - Kubernetes package manager
- [AWS CLI](https://aws.amazon.com/cli/) - AWS Command Line Interface

## Local Development

For local development, the infrastructure can be provisioned with a local state file:

```bash
# Deploy with local state in develop environment
task bootstrap:local ENV=develop

# Port forward ArgoCD server to localhost
task argocd:portforward

# Get ArgoCD password
task argocd:password

# Destroy infrastructure
task destroy ENV=develop
```

## CI/CD Deployment

For CI/CD pipelines, the infrastructure is provisioned with remote state in S3:

```bash
# Deploy with remote state in production environment
task bootstrap:remote ENV=production REGION=us-east-2

# List ArgoCD applications
task argocd:apps
```

## Available Tasks

Run `task --list` to see all available tasks:

```
task: Available tasks for this project:
* argocd:apps:          Show ArgoCD applications
* argocd:password:      Get ArgoCD admin password
* argocd:portforward:   Port forward ArgoCD server to localhost:8080
* apply:                Apply infrastructure changes
* bootstrap:local:      Bootstrap ArgoCD locally with local state
* bootstrap:remote:     Bootstrap ArgoCD with remote state (for CI/CD)
* check:deps:           Check for required dependencies
* destroy:              Destroy infrastructure
* init:                 Initialize OpenTofu with local or remote backend
* plan:                 Plan infrastructure changes
* set:vars:             Set environment-specific variables
* setup:s3:             Setup S3 bucket for remote state (only used when USE_REMOTE=true)
```
