output "cluster_id" {
  description = "The name/id of the EKS cluster"
  value       = aws_eks_cluster.this.id
}

output "cluster_arn" {
  description = "The Amazon Resource Name (ARN) of the EKS cluster"
  value       = aws_eks_cluster.this.arn
}

output "cluster_endpoint" {
  description = "The endpoint for the Kubernetes API server"
  value       = aws_eks_cluster.this.endpoint
}

output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = aws_eks_cluster.this.certificate_authority[0].data
}

output "cluster_name" {
  description = "The name of the EKS cluster"
  value       = aws_eks_cluster.this.name
}

output "oidc_provider_arn" {
  description = "The ARN of the OIDC Provider if IRSA is enabled"
  value       = var.enable_irsa ? aws_iam_openid_connect_provider.this[0].arn : ""
}

output "oidc_provider_url" {
  description = "The URL of the OIDC Provider if IRSA is enabled"
  value       = var.enable_irsa ? aws_eks_cluster.this.identity[0].oidc[0].issuer : ""
}

output "cluster_security_group_id" {
  description = "Security group ID attached to the EKS cluster"
  value       = aws_security_group.cluster.id
}

output "node_security_group_id" {
  description = "Security group ID attached to the EKS node groups"
  value       = aws_security_group.cluster.id
}

output "node_groups" {
  description = "Outputs from EKS node groups"
  value = {
    for k, v in aws_eks_node_group.this : k => {
      id         = v.id
      arn        = v.arn
      status     = v.status
      node_group = v.node_group_name
    }
  }
}
