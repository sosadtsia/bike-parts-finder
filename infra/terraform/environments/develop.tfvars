# Environment configuration for develop
environment = "develop"
region      = "us-east-2"
prefix      = "use2-develop-"

# Cluster configuration
cluster_name       = "bpf-cluster"
kubernetes_version = "1.30"

# Node groups configuration
node_groups = {
  default = {
    name           = "default"
    instance_types = ["t3.medium"]
    desired_size   = 2
    min_size       = 1
    max_size       = 3
    disk_size      = 50
  }
}

# VPC settings
vpc_cidr = "10.0.0.0/16"

# Tags
tags = {
  environment = "develop"
  project     = "bike-parts-finder"
  managedBy   = "terraform"
}
