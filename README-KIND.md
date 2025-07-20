# Local Development with Kind Cluster

This document explains how to set up a local development environment using Kind (Kubernetes in Docker) through OpenTofu (Terraform compatible).

## Overview

For local development, we use:
- **Kind** for running a local Kubernetes cluster
- **OpenTofu** to provision and manage the Kind cluster
- **Local Docker Registry** to store container images
- **ArgoCD** for GitOps deployment
- **Helmfile** to deploy services to the cluster

## Prerequisites

- Docker
- kubectl
- OpenTofu (or Terraform)
- Helm
- Helmfile
- Task (taskfile.dev)

## Setting up the Local Environment

To set up the complete local development environment, run:

```bash
task setup:local
```

This will:
1. Initialize OpenTofu
2. Plan and apply the Kind cluster creation
3. Set up a local Docker registry
4. Build and push Docker images
5. Deploy services using Helmfile

## Managing the Kind Cluster

### Initialize OpenTofu for Kind

```bash
task kind:init
```

### Plan Kind Cluster Creation

```bash
task kind:plan
```

### Create/Update Kind Cluster

```bash
task kind:apply
```

### Destroy Kind Cluster

```bash
task kind:destroy
```

## Working with Docker Images

### Build All Docker Images

```bash
task docker:build-all
```

### Build Individual Services

```bash
task docker:build-api
task docker:build-scraper
task docker:build-consumer
task docker:build-frontend
```

## Deploying with Helmfile

### Deploy Services to Kind

```bash
task helmfile:develop
```

### View Deployment Differences

```bash
task helmfile:diff
```

## Working with ArgoCD

### Get ArgoCD Admin Password

```bash
task argocd:password
```

### Access ArgoCD UI

```bash
task argocd:port-forward
```

Then open https://localhost:8080 in your browser and login with username "admin" and the password obtained from the previous command.

## Cleaning Up

To destroy the Kind cluster and clean up resources:

```bash
task clean:local
```

## Architecture

The local development environment mirrors the production architecture:

- React frontend
- Go backend API
- Kafka for asynchronous processing
- PostgreSQL for data storage
- Redis for caching

All components are containerized and deployed to the local Kind cluster using the same Helm charts and Kubernetes manifests used in production.
