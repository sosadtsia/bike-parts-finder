output "argocd_namespace" {
  description = "The namespace where ArgoCD is installed"
  value       = kubernetes_namespace.argocd.metadata[0].name
}

output "argocd_server_service_name" {
  description = "The name of the ArgoCD server service"
  value       = "${helm_release.argocd.name}-server"
}

output "argocd_url" {
  description = "The URL to access ArgoCD"
  value       = var.create_ingress ? "https://argocd-${var.environment}.${var.domain}" : "Use port-forward to access ArgoCD UI"
}
