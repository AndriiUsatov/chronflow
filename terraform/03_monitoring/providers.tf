terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "3.0.1"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "4.25.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "grafana" {
  url = var.grafana_url
  auth = var.grafana_api_token
}
