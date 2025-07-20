variable "namespace" {
  description = "The namespace where ArgoCD will be installed"
  type        = string
  default     = "argocd"
}

variable "argocd_version" {
  description = "Version of the ArgoCD Helm chart"
  type        = string
  default     = "5.51.4"
}

variable "environment" {
  description = "Environment (develop, production)"
  type        = string
  default     = "develop"
}

variable "domain" {
  description = "The domain name for ArgoCD ingress"
  type        = string
  default     = "bikepartsfinder.example.com"
}

variable "repo_url" {
  description = "The Git repository URL containing the application manifests"
  type        = string
  default     = "https://github.com/svosadtsia/bike-parts-finder.git"
}

variable "create_ingress" {
  description = "Whether to create an ingress for ArgoCD"
  type        = bool
  default     = true
}

variable "enable_dex" {
  description = "Whether to enable Dex for authentication"
  type        = bool
  default     = false
}

variable "cluster_autoscaler_enabled" {
  description = "Whether the EKS cluster has the cluster autoscaler enabled"
  type        = bool
  default     = true
}
