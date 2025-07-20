output "cluster_name" {
  description = "Name of the created Kind cluster"
  value       = module.kind.cluster_name
}

output "kubeconfig_path" {
  description = "Path to the kubeconfig file"
  value       = module.kind.kubeconfig_path
}

output "registry_url" {
  description = "URL of the local registry"
  value       = module.kind.registry_url
}

output "argocd_admin_password" {
  description = "ArgoCD admin password"
  value       = "Run the command: kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d"
  sensitive   = true
}
