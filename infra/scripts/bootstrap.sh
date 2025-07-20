#!/bin/bash
# Bootstrap script for deploying infrastructure and ArgoCD using OpenTofu

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TERRAFORM_DIR="${SCRIPT_DIR}/../terraform"
ENV=${1:-develop}

echo "⚙️  Bootstrapping Bike Parts Finder Infrastructure and ArgoCD (environment: ${ENV})..."

# Check for required tools
check_tool() {
  if ! command -v "$1" &> /dev/null; then
    echo "❌ $1 is required but not installed. Please install it first."
    exit 1
  fi
}

check_tool tofu
check_tool kubectl
check_tool helm

# Create S3 backend if it doesn't exist
setup_backend() {
  local bucket="bike-parts-finder-tfstate-${ENV}"
  local region="us-east-2"
  local table="bike-parts-finder-tfstate-lock-${ENV}"

  echo "🪣 Setting up S3 backend for OpenTofu state..."

  # Check if bucket exists
  if ! aws s3api head-bucket --bucket "${bucket}" --region "${region}" 2>/dev/null; then
    echo "Creating S3 bucket for state: ${bucket}"
    aws s3api create-bucket \
      --bucket "${bucket}" \
      --region "${region}" \
      --create-bucket-configuration LocationConstraint="${region}"

    aws s3api put-bucket-versioning \
      --bucket "${bucket}" \
      --versioning-configuration Status=Enabled \
      --region "${region}"

    aws s3api put-bucket-encryption \
      --bucket "${bucket}" \
      --server-side-encryption-configuration '{
        "Rules": [
          {
            "ApplyServerSideEncryptionByDefault": {
              "SSEAlgorithm": "AES256"
            }
          }
        ]
      }' \
      --region "${region}"
  fi

  # Check if DynamoDB table exists
  if ! aws dynamodb describe-table --table-name "${table}" --region "${region}" 2>/dev/null; then
    echo "Creating DynamoDB table for state locking: ${table}"
    aws dynamodb create-table \
      --table-name "${table}" \
      --attribute-definitions AttributeName=LockID,AttributeType=S \
      --key-schema AttributeName=LockID,KeyType=HASH \
      --billing-mode PAY_PER_REQUEST \
      --region "${region}"
  fi

  # Create backend config file
  cat > "${TERRAFORM_DIR}/backend-${ENV}.conf" <<EOF
bucket         = "${bucket}"
key            = "terraform.tfstate"
region         = "${region}"
dynamodb_table = "${table}"
encrypt        = true
EOF

  echo "✅ Backend configuration created at ${TERRAFORM_DIR}/backend-${ENV}.conf"
}

# Initialize OpenTofu
init_tofu() {
  echo "🔄 Initializing OpenTofu..."
  cd "${TERRAFORM_DIR}"
  tofu init -backend-config="backend-${ENV}.conf"
}

# Create terraform.tfvars
create_tfvars() {
  local tfvars_file="${TERRAFORM_DIR}/terraform.tfvars"

  echo "📝 Creating terraform.tfvars..."
  cat > "${tfvars_file}" <<EOF
environment = "${ENV}"
aws_region  = "us-east-2"
project_name = "bike-parts-finder"
domain_suffix = "${ENV}.example.com"
EOF

  echo "✅ Terraform variables file created at ${tfvars_file}"
}

# Apply OpenTofu configuration
apply_tofu() {
  echo "🚀 Applying OpenTofu configuration..."
  cd "${TERRAFORM_DIR}"
  tofu plan -out=tfplan
  tofu apply -auto-approve tfplan
}

# Get and display ArgoCD initial password
get_argocd_password() {
  echo "🔑 Retrieving ArgoCD admin password..."
  kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
}

# Verify ArgoCD access
verify_argocd() {
  echo "🔍 Verifying ArgoCD deployment..."
  kubectl -n argocd get pods

  echo "📋 ArgoCD applications:"
  kubectl -n argocd get applications
}

# Main execution
setup_backend
init_tofu
create_tfvars
apply_tofu
get_argocd_password
verify_argocd

echo ""
echo "✅ Bootstrap complete!"
echo ""
echo "Access ArgoCD UI:"
echo "  1. Port forward: kubectl port-forward svc/argocd-server -n argocd 8080:443"
echo "  2. Open in browser: https://localhost:8080"
echo "  3. Login with username 'admin' and the password shown above"
echo ""
echo "Next steps:"
echo "  1. Update ArgoCD admin password"
echo "  2. Configure repository credentials if using private repository"
echo "  3. Verify ApplicationSet is creating applications"
echo ""
