## ArgoCD Helm Values
global:
  domain: ${domain}

server:
  extraArgs:
    - --insecure
  ingress:
    enabled: ${create_ingress}
    ingressClassName: nginx
    hosts:
      - ${domain}
    annotations:
      kubernetes.io/tls-acme: "true"
    tls:
      - secretName: argocd-tls-secret
        hosts:
          - ${domain}
  config:
    repositories: |
      - type: git
        url: ${repo_url}
        name: app-repo
    resource.customizations: |
      argoproj.io/Application:
        health.lua: |
          hs = {}
          hs.status = "Progressing"
          hs.message = ""
          if obj.status ~= nil then
            if obj.status.health ~= nil then
              hs.status = obj.status.health.status
              if obj.status.health.message ~= nil then
                hs.message = obj.status.health.message
              end
            end
          end
          return hs

controller:
  metrics:
    enabled: true
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 256Mi

dex:
  enabled: ${enable_dex}

redis:
  resources:
    limits:
      cpu: 200m
      memory: 128Mi
    requests:
      cpu: 50m
      memory: 64Mi

repoServer:
  resources:
    limits:
      cpu: 300m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
