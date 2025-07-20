variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "bike-parts-finder"
}

variable "environment" {
  description = "Environment (develop, production)"
  type        = string
  default     = "develop"
}

variable "aws_region" {
  description = "AWS region to deploy to"
  type        = string
  default     = "us-east-2"
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

variable "db_instance_class" {
  description = "RDS instance type"
  type        = string
  default     = "db.t3.medium"
}

variable "db_allocated_storage" {
  description = "Allocated storage for the database in GB"
  type        = number
  default     = 20
}

variable "db_name" {
  description = "Name of the database"
  type        = string
  default     = "bikepartsfinder"
}

variable "db_username" {
  description = "Master username for the database"
  type        = string
  default     = "postgres" # In production, this should be set via environment variables or secrets
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
