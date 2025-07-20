#!/bin/bash
# Script to bootstrap ArgoCD on a Kubernetes cluster

set -e

echo "⚙️  Bootstrapping ArgoCD for Bike Parts Finder..."

# Create namespace
echo "📁 Creating ArgoCD namespace..."
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -

# Install ArgoCD
echo "🚀 Installing ArgoCD..."
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.8.4/manifests/install.yaml

# Wait for ArgoCD to be ready
echo "⏳ Waiting for ArgoCD to be ready..."
kubectl wait --for=condition=available deployment/argocd-server -n argocd --timeout=300s

# Add Helmfile plugin
echo "🔌 Adding Helmfile plugin..."
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: helmfile-plugin
  namespace: argocd
data:
  plugin.yaml: |
    apiVersion: argoproj.io/v1alpha1
    kind: ConfigManagementPlugin
    metadata:
      name: helmfile
    spec:
      version: v1.0
      init:
        command: [sh, -c]
        args: ["helm plugin install https://github.com/mumoshu/helmfile-diff --version=v3.8.1 || true"]
      generate:
        command: [sh, -c]
        args: ["helmfile --no-color -f helmfile.yaml \$HELMFILE_GLOBAL_OPTIONS --environment \$HELMFILE_ENVIRONMENT \$HELMFILE_SELECTOR template"]
EOF

# Apply the ApplicationSet
echo "📋 Applying ApplicationSet..."
kubectl apply -f ../argocd/applicationsets/all-components.yaml

# Get initial admin password
echo "🔑 Retrieving initial admin password..."
ARGOCD_PASSWORD=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)

echo ""
echo "✅ ArgoCD bootstrap complete!"
echo ""
echo "Access ArgoCD UI:"
echo "  1. Port forward: kubectl port-forward svc/argocd-server -n argocd 8080:443"
echo "  2. Open in browser: https://localhost:8080"
echo "  3. Login with:"
echo "     Username: admin"
echo "     Password: $ARGOCD_PASSWORD"
echo ""
echo "Important: After login, update the password using:"
echo "  argocd account update-password"
echo ""
