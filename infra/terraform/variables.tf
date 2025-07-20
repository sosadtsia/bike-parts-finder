variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "bike-parts-finder"
}

variable "environment" {
  description = "The environment to deploy to"
  type        = string
  default     = "develop"
}

variable "aws_region" {
  description = "AWS region to deploy to"
  type        = string
  default     = "us-east-2"
}

variable "region" {
  description = "AWS region to deploy to (alias for aws_region for compatibility with tfvars)"
  type        = string
  default     = "us-east-2"
}

variable "prefix" {
  description = "The prefix used for deployment purposes, affects all names of all resources"
  type        = string
  default     = "use2-develop-"
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
  default     = null
}

variable "kubernetes_version" {
  description = "Kubernetes version to use for the EKS cluster"
  type        = string
  default     = null
}

variable "cluster_version" {
  description = "Kubernetes version to use for EKS cluster"
  type        = string
  default     = "1.28"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones to use"
  type        = list(string)
  default     = ["us-east-2a", "us-east-2b", "us-east-2c"]
}

variable "private_subnets" {
  description = "List of private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "public_subnets" {
  description = "List of public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = null
}

variable "node_groups" {
  description = "Map of EKS node group configurations"
  type        = map(object({
    name           = string
    instance_types = list(string)
    min_size       = number
    max_size       = number
    desired_size   = number
    disk_size      = number
  }))
  default = {
    default = {
      name           = "default"
      instance_types = ["t3.medium"]
      min_size       = 1
      max_size       = 3
      desired_size   = 2
      disk_size      = 50
    }
  }
}

variable "velero_bucket_name" {
  description = "Name of the S3 bucket for Velero backups"
  type        = string
  default     = null # Will be auto-generated if not specified
}

variable "domain_suffix" {
  description = "Domain suffix for the project (e.g., example.com)"
  type        = string
  default     = "example.com"
}

variable "repo_url" {
  description = "Git repository URL for the project"
  type        = string
  default     = "https://github.com/svosadtsia/bike-parts-finder.git"
}

variable "create_argocd_ingress" {
  description = "Whether to create an ingress for ArgoCD"
  type        = bool
  default     = true
}

variable "enable_argocd_dex" {
  description = "Whether to enable Dex authentication for ArgoCD"
  type        = bool
  default     = false
}

variable "deploy_argocd" {
  description = "Whether to deploy ArgoCD Kubernetes resources"
  type        = bool
  default     = false
}
