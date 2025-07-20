#!/bin/bash
set -e

# Create namespaces
kubectl create namespace databases --dry-run=client -o yaml | kubectl apply -f -

# Add Helm repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

# Install PostgreSQL
echo "Installing PostgreSQL..."
helm install postgres bitnami/postgresql \
  --set auth.username=bikepartsfinder \
  --set auth.password=bikepartsfinder \
  --set auth.database=bikepartsfinder \
  --set primary.persistence.enabled=true \
  --set primary.persistence.size=1Gi

# Create PostgreSQL credentials secret in default namespace
kubectl create secret generic postgres-credentials \
  --from-literal=username=bikepartsfinder \
  --from-literal=password=bikepartsfinder \
  --dry-run=client -o yaml | kubectl apply -f -

# Install Redis
echo "Installing Redis..."
helm install redis bitnami/redis \
  --set auth.enabled=false \
  --set master.persistence.enabled=true \
  --set master.persistence.size=1Gi

# Create Redis credentials secret in default namespace
kubectl create secret generic redis-credentials \
  --from-literal=password="" \
  --dry-run=client -o yaml | kubectl apply -f -

# Install Kafka
echo "Installing Kafka..."
helm install kafka bitnami/kafka \
  --set persistence.enabled=true \
  --set persistence.size=1Gi \
  --set zookeeper.persistence.enabled=true \
  --set zookeeper.persistence.size=1Gi

# Install ArgoCD
echo "Installing ArgoCD..."
helm install argocd argo/argo-cd \
  --namespace argocd \
  --create-namespace \
  --set server.service.type=NodePort \
  --set server.extraArgs="{--insecure}"

# Wait for ArgoCD to become ready
echo "Waiting for ArgoCD to become ready..."
kubectl wait --namespace argocd \
  --for=condition=available deployment/argocd-server \
  --timeout=120s

# Get the ArgoCD password
echo "ArgoCD admin password:"
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
echo ""

# Set up port forwarding instructions
echo ""
echo "To access ArgoCD UI, run:"
echo "kubectl port-forward svc/argocd-server -n argocd 8080:443"
echo "Then open https://localhost:8080 in your browser"
echo "Login with username: admin and the password shown above"

echo ""
echo "All dependencies installed successfully!"
