resource "kubernetes_deployment_v1" "chronflow-scheduler-deployment" {
  metadata {
    name      = "chronflow-scheduler-deployment"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    labels = {
      app = "chronflow-scheduler"
    }
  }

  spec {
    selector {
      match_labels = {
        app = "chronflow-scheduler"
      }
    }
    template {
      metadata {
        labels = {
          app = "chronflow-scheduler"
        }
      }
      spec {
        container {
          name              = "chronflow-scheduler"
          image             = "chronflow-scheduler:latest"
          image_pull_policy = "Always"
          env_from {
            config_map_ref {
              name = data.terraform_remote_state.infrastructureme_out.outputs.chronflow-postgres-config-name
            }
          }
          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.chronflow-nats-config.metadata[0].name
            }
          }
          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.chronflow-scheduler-config.metadata[0].name
            }
          }
          env_from {
            secret_ref {
              name = "chronflow-postgres-creds"
            }
          }
          port {
            name           = "port"
            container_port = var.scheduler_port
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

resource "kubernetes_horizontal_pod_autoscaler_v2" "chronflow-scheduler-hpa" {
  metadata {
    name      = "chronflow-scheduler-hpa"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }
  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment_v1.chronflow-scheduler-deployment.metadata[0].name
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
