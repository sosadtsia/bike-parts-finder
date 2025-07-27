# ArgoCD Setup for Bike Parts Finder

This directory contains ArgoCD configurations for deploying all components of the Bike Parts Finder application.

## Directory Structure

- `apps/` - Contains all application manifests for individual components
- `app-of-apps.yaml` - Main application that deploys all component applications
- `argocd.yaml` - Self-managed ArgoCD application
- `helmfile-plugin-config.yaml` - ConfigMap for Helmfile plugin configuration

## Setup Instructions

1. Install ArgoCD in your cluster:

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

2. Apply the Helmfile plugin configuration:

```bash
kubectl apply -f helmfile-plugin-config.yaml
```

3. Deploy ArgoCD as a self-managed application:

```bash
kubectl apply -f argocd.yaml
```

4. Deploy all applications using the app-of-apps pattern:

```bash
kubectl apply -f app-of-apps.yaml
```

## Access ArgoCD UI

To access the ArgoCD UI, use port forwarding:

```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

Then visit: https://localhost:8080

Default credentials:
- Username: admin
- Password: Get the auto-generated password with:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

## Components Deployed

The following components are deployed through ArgoCD:

1. ArgoCD (self-managed)
2. API Service
3. Scraper Service
4. Consumer Service
5. PostgreSQL Database
6. Redis Cache
7. Kafka Messaging
8. Ingress NGINX Controller
9. Prometheus Stack (including Grafana)

Each component is defined in its own Application manifest and deployed via the Helmfile plugin.
