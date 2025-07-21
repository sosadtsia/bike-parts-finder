terraform {
  backend "s3" {
    # Backend configuration will be provided via CLI arguments
    # using -backend-config parameters
  }
}

provider "aws" {
  region = var.aws_region
}

locals {
  cluster_name = var.cluster_name != null ? var.cluster_name : "${var.prefix}${var.project_name}-cluster"
  vpc_name     = "${var.prefix}${var.project_name}-vpc"

  # Update capitalization to match the tag policy requirements
  tags = var.tags != null ? var.tags : {
    Environment = var.environment  # This should use the correct capitalization
    Project     = var.project_name
    ManagedBy   = "terraform"
  }
}

# Create VPC for EKS
module "vpc" {
  source = "./modules/vpc"

  vpc_name        = local.vpc_name
  vpc_cidr        = var.vpc_cidr
  azs             = var.availability_zones
  private_subnets = var.private_subnets
  public_subnets  = var.public_subnets
}

# Create EKS Cluster
module "eks" {
  source = "./modules/eks"

  cluster_name     = local.cluster_name
  cluster_version  = var.kubernetes_version != null ? var.kubernetes_version : var.cluster_version
  vpc_id           = module.vpc.vpc_id
  subnet_ids       = module.vpc.private_subnets
  node_groups      = var.node_groups
  tags             = local.tags
}

