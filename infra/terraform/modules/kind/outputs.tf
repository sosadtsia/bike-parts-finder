output "cluster_name" {
  description = "Name of the created Kind cluster"
  value       = kind_cluster.cluster.name
}

output "kubeconfig_path" {
  description = "Path to the kubeconfig file"
  value       = kind_cluster.cluster.kubeconfig_path
}

output "registry_url" {
  description = "URL of the local registry"
  value       = "localhost:${var.registry_port}"
}
