resource "kubernetes_deployment_v1" "chronflow-api-deployment" {
  metadata {
    name      = "chronflow-api-deployment"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    labels = {
      app = "chronflow-api"
    }
  }
  spec {
    selector {
      match_labels = {
        app = "chronflow-api"
      }
    }
    template {
      metadata {
        labels = {
          app = "chronflow-api"
        }
      }
      spec {
        container {
          name              = "chronflow-api"
          image             = "chronflow-api:latest"
          image_pull_policy = "Always"
          env_from {
            config_map_ref {
              name = data.terraform_remote_state.infrastructureme_out.outputs.chronflow-postgres-config-name
            }
          }
          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.chronflow-api-config.metadata[0].name
            }
          }
          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.chronflow-grpc-config.metadata[0].name
            }
          }
          env_from {
            secret_ref {
              name = "chronflow-postgres-creds"
            }
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
          port {
            name           = "rest"
            container_port = var.task_api_port
          }
          port {
            name           = "grpc"
            container_port = var.task_grpc_port
          }
        }
      }
    }
  }
}

resource "kubernetes_horizontal_pod_autoscaler_v2" "chronflow-api-hpa" {
  metadata {
    name      = "chronflow-api-hpa"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }
  spec {
    scale_target_ref {
      api_version = "apps/v1"
      kind        = "Deployment"
      name        = kubernetes_deployment_v1.chronflow-api-deployment.metadata[0].name
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

resource "kubernetes_service_v1" "chronflow-api-http-service" {
  metadata {
    name      = "chronflow-api-http-service"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    labels = {
      app = "chronflow-api"
    }
  }

  spec {
    type = "LoadBalancer"
    selector = {
      app = "chronflow-api"
    }
    port {
      name        = "rest-api"
      protocol    = "TCP"
      port        = 80
      target_port = "rest"
    }
  }
}


resource "kubernetes_service_v1" "chronflow-api-grpc-service" {
  metadata {
    name      = "chronflow-api-grpc-service"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    labels = {
      app = "chronflow-api"
    }
  }

  spec {
    type       = "ClusterIP"
    cluster_ip = "None"
    selector = {
      app = "chronflow-api"
    }
    port {
      name        = "grpc-api"
      protocol    = "TCP"
      port        = var.task_grpc_port
      target_port = "grpc"
    }
  }
}
