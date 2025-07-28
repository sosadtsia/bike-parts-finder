# Create a local to determine if the Kubernetes provider is available
locals {
  has_kubernetes_provider = var.eks_endpoint != ""
}

resource "kubernetes_namespace" "argocd" {
  metadata {
    name = var.namespace
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = var.argocd_chart_version
  namespace  = kubernetes_namespace.argocd.metadata[0].name
  timeout    = 600

  values = [
    templatefile("${path.module}/templates/values.yaml", {
      domain = var.domain
    })
  ]
}

# ArgoCD ConfigMaps for plugins and RBAC
resource "kubernetes_config_map" "argocd_cm" {
  depends_on = [helm_release.argocd]

  metadata {
    name      = "argocd-cm"
    namespace = kubernetes_namespace.argocd.metadata[0].name
    labels = {
      "app.kubernetes.io/name"    = "argocd-cm"
      "app.kubernetes.io/part-of" = "argocd"
    }
  }

  data = {
    "configManagementPlugins" = <<-EOT
      - name: helmfile
        init:
          command: ["/bin/sh", "-c"]
          args: ["helm plugin install https://github.com/databus23/helm-diff || true && helm repo update"]
        generate:
          command: ["/bin/sh", "-c"]
          args: ["cd $ARGOCD_APP_SOURCE_PATH && helmfile -e default -l name=$HELMFILE_RELEASE template"]
    EOT
  }
}

resource "kubernetes_config_map" "argocd_rbac_cm" {
  depends_on = [helm_release.argocd]

  metadata {
    name      = "argocd-rbac-cm"
    namespace = kubernetes_namespace.argocd.metadata[0].name
    labels = {
      "app.kubernetes.io/name"    = "argocd-rbac-cm"
      "app.kubernetes.io/part-of" = "argocd"
    }
  }

  data = {
    "policy.csv" = <<-EOT
      p, role:admin, applications, *, */*, allow
      p, role:admin, clusters, *, *, allow
      p, role:admin, repositories, *, *, allow
      p, role:admin, logs, *, *, allow
      p, role:admin, exec, *, */*, allow
      g, admin, role:admin
    EOT
  }
}

# Self-managed ArgoCD application
resource "kubectl_manifest" "argocd_application" {
  depends_on = [
    kubernetes_config_map.argocd_cm,
    kubernetes_config_map.argocd_rbac_cm
  ]

  yaml_body = templatefile("${path.module}/templates/argocd-application.yaml", {
    namespace = kubernetes_namespace.argocd.metadata[0].name
    domain    = var.domain
    chart_version = var.argocd_chart_version
  })
}

# Application manifests for components
resource "kubectl_manifest" "component_applications" {
  for_each = var.components

  depends_on = [kubectl_manifest.argocd_application]

  yaml_body = templatefile("${path.module}/templates/application.yaml", {
    name      = each.key
    namespace = kubernetes_namespace.argocd.metadata[0].name
    app_namespace = var.app_namespace
    repo_url  = var.repo_url
    release_name = each.key
  })
}

# App of Apps pattern
resource "kubectl_manifest" "app_of_apps" {
  depends_on = [kubectl_manifest.component_applications]

  yaml_body = templatefile("${path.module}/templates/app-of-apps.yaml", {
    namespace = kubernetes_namespace.argocd.metadata[0].name
    repo_url  = var.repo_url
    app_path  = "infra/argocd/apps"
  })
}
