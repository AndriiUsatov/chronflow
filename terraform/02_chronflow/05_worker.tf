resource "kubernetes_deployment_v1" "chronflow-worker-deployment" {
  metadata {
    name      = "chronflow-worker-deployment"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    labels = {
      app = "chronflow-worker"
    }
  }
  spec {
    selector {
      match_labels = {
        app = "chronflow-worker"
      }
    }
    template {
      metadata {
        labels = {
          app = "chronflow-worker"
        }
      }
      spec {
        container {
          name              = "chronflow-worker"
          image             = "chronflow-worker:latest"
          image_pull_policy = "Always"
          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.chronflow-grpc-config.metadata[0].name
            }
          }
          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.chronflow-nats-config.metadata[0].name
            }
          }
          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.chronflow-worker-config.metadata[0].name
            }
          }
          port {
            name           = "port"
            container_port = var.worker_port
          }
          resources {
            requests = {
              memory = "128Mi"
              cpu    = "250m"
            }
            limits = {
              memory = "256Mi"
              cpu    = "500m"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_horizontal_pod_autoscaler_v2" "chronflow-worker-hpa" {
  metadata {
    name      = "chronflow-worker-hpa"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }
  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment_v1.chronflow-worker-deployment.metadata[0].name
    }
    min_replicas = 1
    max_replicas = 3
    metric {
      type = "Resource"
      resource {
        name = "cpu"
        target {
          type                = "Utilization"
          average_utilization = 50
        }
      }
    }
  }
}
