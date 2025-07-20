variable "cluster_name" {
  description = "Name of the Kind cluster"
  type        = string
  default     = "bike-parts-finder"
}

variable "registry_name" {
  description = "Name of the local Docker registry"
  type        = string
  default     = "kind-registry"
}

variable "registry_port" {
  description = "Port for the local Docker registry"
  type        = number
  default     = 5000
}
