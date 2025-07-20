terraform {
  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = "~> 0.2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.2"
    }
  }
  required_version = ">= 1.0.0"
}

provider "kind" {}

provider "kubernetes" {
  config_path = module.kind.kubeconfig_path
}

provider "helm" {
  kubernetes {
    config_path = module.kind.kubeconfig_path
  }
}

provider "docker" {}

module "kind" {
  source = "../modules/kind"

  cluster_name  = var.cluster_name
  registry_name = var.registry_name
  registry_port = var.registry_port
}

# Deploy bitnami helm repo
resource "helm_repository" "bitnami" {
  name = "bitnami"
  url  = "https://charts.bitnami.com/bitnami"

  depends_on = [module.kind]
}

# Deploy ArgoCD helm repo
resource "helm_repository" "argo" {
  name = "argo"
  url  = "https://argoproj.github.io/argo-helm"

  depends_on = [module.kind]
}

# Deploy PostgreSQL
resource "helm_release" "postgres" {
  name       = "postgres"
  repository = helm_repository.bitnami.metadata[0].name
  chart      = "postgresql"

  set {
    name  = "auth.username"
    value = "bikepartsfinder"
  }

  set {
    name  = "auth.password"
    value = "bikepartsfinder"
  }

  set {
    name  = "auth.database"
    value = "bikepartsfinder"
  }

  set {
    name  = "primary.persistence.enabled"
    value = "true"
  }

  set {
    name  = "primary.persistence.size"
    value = "1Gi"
  }

  depends_on = [helm_repository.bitnami]
}

# Create PostgreSQL credentials secret
resource "kubernetes_secret" "postgres_credentials" {
  metadata {
    name = "postgres-credentials"
  }

  data = {
    username = "bikepartsfinder"
    password = "bikepartsfinder"
  }

  depends_on = [module.kind]
}

# Deploy Redis
resource "helm_release" "redis" {
  name       = "redis"
  repository = helm_repository.bitnami.metadata[0].name
  chart      = "redis"

  set {
    name  = "auth.enabled"
    value = "false"
  }

  set {
    name  = "master.persistence.enabled"
    value = "true"
  }

  set {
    name  = "master.persistence.size"
    value = "1Gi"
  }

  depends_on = [helm_repository.bitnami]
}

# Create Redis credentials secret
resource "kubernetes_secret" "redis_credentials" {
  metadata {
    name = "redis-credentials"
  }

  data = {
    password = ""
  }

  depends_on = [module.kind]
}

# Deploy Kafka
resource "helm_release" "kafka" {
  name       = "kafka"
  repository = helm_repository.bitnami.metadata[0].name
  chart      = "kafka"

  set {
    name  = "persistence.enabled"
    value = "true"
  }

  set {
    name  = "persistence.size"
    value = "1Gi"
  }

  set {
    name  = "zookeeper.persistence.enabled"
    value = "true"
  }

  set {
    name  = "zookeeper.persistence.size"
    value = "1Gi"
  }

  depends_on = [helm_repository.bitnami]
}

# Deploy ArgoCD
resource "helm_release" "argocd" {
  name       = "argocd"
  repository = helm_repository.argo.metadata[0].name
  chart      = "argo-cd"
  namespace  = "argocd"
  create_namespace = true

  set {
    name  = "server.service.type"
    value = "NodePort"
  }

  set {
    name  = "server.extraArgs"
    value = "{--insecure}"
  }

  depends_on = [helm_repository.argo]
}
