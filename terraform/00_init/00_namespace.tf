resource "kubernetes_namespace_v1" "chronflow_namespace" {
  metadata {
    name = "chronflow"

    labels = {
      managed-by = "terraform"
    }

    annotations = {
      description = "Namespace for Chronflow app resources"
    }
  }
}

resource "kubernetes_namespace_v1" "monitoring_namespace" {
  metadata {
    name = "monitoring"

    labels = {
      managed-by = "terraform"
    }

    annotations = {
      description = "Namespace for monitoring infrastructure"
    }
  }
}
