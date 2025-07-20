resource "kubernetes_namespace" "argocd" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "argocd"
    }
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = var.argocd_version
  namespace  = kubernetes_namespace.argocd.metadata[0].name

  values = [
    templatefile("${path.module}/values.yaml.tpl", {
      environment                = var.environment
      domain                     = var.domain
      repo_url                   = var.repo_url
      create_ingress             = var.create_ingress
      enable_dex                 = var.enable_dex
      cluster_autoscaler_enabled = var.cluster_autoscaler_enabled
    })
  ]

  depends_on = [
    kubernetes_namespace.argocd
  ]
}

resource "kubernetes_config_map" "helmfile_plugin" {
  metadata {
    name      = "helmfile-plugin"
    namespace = kubernetes_namespace.argocd.metadata[0].name
  }

  data = {
    "plugin.yaml" = <<-EOT
    apiVersion: argoproj.io/v1alpha1
    kind: ConfigManagementPlugin
    metadata:
      name: helmfile
    spec:
      version: v1.0
      init:
        command: [sh, -c]
        args: ["helm plugin install https://github.com/mumoshu/helmfile-diff --version=v3.8.1 || true"]
      generate:
        command: [sh, -c]
        args: ["helmfile --no-color -f helmfile.yaml $$HELMFILE_GLOBAL_OPTIONS --environment $$HELMFILE_ENVIRONMENT $$HELMFILE_SELECTOR template"]
    EOT
  }

  depends_on = [
    helm_release.argocd
  ]
}

resource "kubernetes_manifest" "applicationset_components" {
  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "ApplicationSet"
    metadata = {
      name      = "bike-parts-finder-components"
      namespace = kubernetes_namespace.argocd.metadata[0].name
    }
    spec = {
      generators = [
        {
          matrix = {
            generators = [
              {
                list = {
                  elements = [
                    {
                      env    = "develop"
                      branch = "HEAD"
                    },
                    {
                      env    = "production"
                      branch = "main"
                    }
                  ]
                }
              },
              {
                list = {
                  elements = [
                    {
                      component = "argocd"
                      path      = "infra/argocd/applications/argocd.yaml"
                    },
                    {
                      component = "postgres"
                      path      = "infra/argocd/applications/postgres.yaml"
                    },
                    {
                      component = "redis"
                      path      = "infra/argocd/applications/redis.yaml"
                    },
                    {
                      component = "kafka"
                      path      = "infra/argocd/applications/kafka.yaml"
                    },
                    {
                      component = "monitoring"
                      path      = "infra/argocd/applications/monitoring.yaml"
                    },
                    {
                      component = "backup"
                      path      = "infra/argocd/applications/backup.yaml"
                    },
                    {
                      component = "api"
                      path      = "infra/argocd/applications/api.yaml"
                    },
                    {
                      component = "frontend"
                      path      = "infra/argocd/applications/frontend.yaml"
                    },
                    {
                      component = "scraper"
                      path      = "infra/argocd/applications/scraper.yaml"
                    },
                    {
                      component = "consumer"
                      path      = "infra/argocd/applications/consumer.yaml"
                    }
                  ]
                }
              }
            ]
          }
        }
      ]
      template = {
        metadata = {
          name = "{{env}}-{{component}}"
          labels = {
            environment = "{{env}}"
            component   = "{{component}}"
          }
        }
        spec = {
          project = "default"
          source = {
            repoURL        = var.repo_url
            targetRevision = "{{branch}}"
            path           = "{{path}}"
          }
          destination = {
            server = "https://kubernetes.default.svc"
          }
          syncPolicy = {
            automated = {
              prune     = true
              selfHeal  = true
            }
          }
        }
      }
    }
  }

  depends_on = [
    helm_release.argocd
  ]
}
