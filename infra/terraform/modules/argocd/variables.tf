variable "namespace" {
  description = "Namespace to deploy ArgoCD"
  type        = string
  default     = "argocd"
}

variable "domain" {
  description = "Domain name for ArgoCD ingress"
  type        = string
  default     = "example.com"
}

variable "argocd_chart_version" {
  description = "Version of the ArgoCD Helm chart"
  type        = string
  default     = "5.36.0"
}

variable "repo_url" {
  description = "Git repository URL for application source"
  type        = string
}

variable "app_namespace" {
  description = "Namespace for application deployments"
  type        = string
  default     = "bike-parts-finder"
}

variable "components" {
  description = "Map of components to deploy"
  type        = map(any)
  default = {
    "api"                 = {},
    "scraper"             = {},
    "consumer"            = {},
    "postgres"            = {},
    "redis"               = {},
    "kafka"               = {},
    "ingress-nginx"       = {},
    "kube-prometheus-stack" = {}
  }
}
