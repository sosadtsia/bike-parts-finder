resource "kind_cluster" "cluster" {
  name           = var.cluster_name
  wait_for_ready = true

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role = "control-plane"

      kubeadm_config_patches = [
        <<-EOT
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
        EOT
      ]

      extra_port_mappings {
        container_port = 80
        host_port      = 80
        protocol       = "TCP"
      }

      extra_port_mappings {
        container_port = 443
        host_port      = 443
        protocol       = "TCP"
      }
    }

    node {
      role = "worker"
    }

    node {
      role = "worker"
    }
  }
}

# Wait for cluster to be ready
resource "null_resource" "wait_for_cluster" {
  depends_on = [kind_cluster.cluster]

  provisioner "local-exec" {
    command = "kubectl wait --for=condition=Ready nodes --all --timeout=300s"
  }
}

# Deploy Ingress NGINX controller
resource "null_resource" "deploy_ingress_nginx" {
  depends_on = [null_resource.wait_for_cluster]

  provisioner "local-exec" {
    command = "kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml"
  }
}

# Setup local registry
resource "docker_container" "registry" {
  name  = var.registry_name
  image = "registry:3"

  restart = "always"

  ports {
    internal = 5000
    external = var.registry_port
  }
}

# Connect registry to kind network
resource "null_resource" "connect_registry" {
  depends_on = [docker_container.registry, kind_cluster.cluster]

  provisioner "local-exec" {
    command = "docker network connect kind ${var.registry_name} || true"
  }
}

# Create ConfigMap for registry
resource "kubernetes_config_map" "registry_config" {
  depends_on = [kind_cluster.cluster]

  metadata {
    name      = "local-registry-hosting"
    namespace = "kube-public"
  }

  data = {
    "localRegistryHosting.v1" = <<-EOT
      host: "localhost:${var.registry_port}"
      help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
    EOT
  }
}

# Wait for ingress controller to be ready
resource "null_resource" "wait_for_ingress" {
  depends_on = [null_resource.deploy_ingress_nginx]

  provisioner "local-exec" {
    command = "kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s"
  }
}
