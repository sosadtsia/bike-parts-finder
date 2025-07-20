server:
  extraArgs:
    - --insecure
  config:
    repositories: |
      - type: git
        url: ${repo_url}
    url: https://argocd-${environment}.${domain}
    application.resourceTrackingMethod: annotation
  rbacConfig:
    policy.default: role:readonly
    policy.csv: |
      g, system:cluster-admins, role:admin
  ingress:
    enabled: ${create_ingress}
    annotations:
      kubernetes.io/ingress.class: alb
      alb.ingress.kubernetes.io/scheme: internet-facing
      alb.ingress.kubernetes.io/target-type: ip
      alb.ingress.kubernetes.io/healthcheck-path: /
    hosts:
      - argocd-${environment}.${domain}

configs:
  secret:
    createSecret: true
  plugins:
    helmfile.yaml: |
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
          args: ["helmfile --no-color -f helmfile.yaml $HELMFILE_GLOBAL_OPTIONS --environment $HELMFILE_ENVIRONMENT $HELMFILE_SELECTOR template"]

dex:
  enabled: ${enable_dex}

repoServer:
  autoscaling:
    enabled: ${cluster_autoscaler_enabled}
    minReplicas: 2

applicationSet:
  enabled: true

global:
  securityContext:
    fsGroup: 65534
