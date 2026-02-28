resource "kubernetes_deployment_v1" "chronflow-janitor-deployment" {
  metadata {
    name      = "chronflow-janitor-deployment"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    labels = {
      app = "chronflow-janitor"
    }
  }
  spec {
    replicas = "1"
    strategy {
      type = "Recreate"
    }
    selector {
      match_labels = {
        app = "chronflow-janitor"
      }
    }
    template {
      metadata {
        labels = {
          app = "chronflow-janitor"
        }
      }
      spec {
        container {
          name              = "chronflow-janitor"
          image             = "chronflow-janitor:latest"
          image_pull_policy = "Always"
          env_from {
            config_map_ref {
              name = data.terraform_remote_state.infrastructureme_out.outputs.chronflow-postgres-config-name
            }
          }
          env_from {
            config_map_ref {
              name = kubernetes_config_map_v1.chronflow-janitor-config.metadata[0].name
            }
          }
          env_from {
            secret_ref {
              name = "chronflow-postgres-creds"
            }
          }
          port {
            name           = "port"
            container_port = var.janitor_hearbeat_port
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
          liveness_probe {
            http_get {
              path = "/heartbeat"
              port = var.janitor_hearbeat_port
            }
            initial_delay_seconds = 15
            period_seconds = 20
          }
        }
      }
    }
  }
}
