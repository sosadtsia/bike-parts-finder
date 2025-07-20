terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.10"
    }
  }

  # Configure remote state backend
  backend "s3" {
    # These will be filled via environment-specific backend config
    # bucket         = "bike-parts-finder-tfstate"
    # key            = "terraform.tfstate"
    # region         = "us-west-2"
    # dynamodb_table = "bike-parts-finder-tfstate-lock"
    # encrypt        = true
  }
}

provider "aws" {
  region = var.aws_region

  # Use default tags for all resources
  default_tags {
    tags = {
      Project     = "bike-parts-finder"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}

# EKS Cluster Module
module "eks" {
  source          = "./modules/eks"
  cluster_name    = "${var.project_name}-${var.environment}"
  cluster_version = var.cluster_version
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnets

  # Node groups configuration
  node_groups = {
    core = {
      name           = "core"
      instance_types = ["t3.medium"]
      min_size       = 2
      max_size       = 3
      desired_size   = 2
      disk_size      = 50
    }

    workloads = {
      name           = "workloads"
      instance_types = ["t3.large"]
      min_size       = 1
      max_size       = 5
      desired_size   = 2
      disk_size      = 100
    }
  }

  # Enable IRSA (IAM Roles for Service Accounts)
  enable_irsa = true

  # Tags
  tags = {
    Environment = var.environment
  }
}

# VPC Module
module "vpc" {
  source               = "./modules/vpc"
  vpc_name             = "${var.project_name}-${var.environment}"
  vpc_cidr             = var.vpc_cidr
  azs                  = var.availability_zones
  private_subnets      = var.private_subnets
  public_subnets       = var.public_subnets
  enable_nat_gateway   = true
  single_nat_gateway   = var.environment == "develop" # Use single NAT gateway for develop environment
  enable_dns_hostnames = true
}

# Security groups
module "security_groups" {
  source      = "./modules/security"
  vpc_id      = module.vpc.vpc_id
  environment = var.environment
}

# S3 bucket for backups
module "backups" {
  source      = "./modules/backups"
  environment = var.environment
  project     = var.project_name
}

# KMS for encryption
module "kms" {
  source      = "./modules/kms"
  environment = var.environment
  project     = var.project_name
}

# ArgoCD Installation
module "argocd" {
  source = "./modules/argocd"

  environment                = var.environment
  domain                     = "${var.project_name}.${var.domain_suffix}"
  repo_url                   = var.repo_url
  create_ingress             = var.create_argocd_ingress
  enable_dex                 = var.enable_argocd_dex
  cluster_autoscaler_enabled = true # We're using node groups with autoscaling

  depends_on = [
    module.eks
  ]
}

# Outputs
output "cluster_name" {
  value = module.eks.cluster_name
}

output "cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "private_subnets" {
  value = module.vpc.private_subnets
}

output "public_subnets" {
  value = module.vpc.public_subnets
}

output "backup_bucket_name" {
  value = module.backups.bucket_name
}
