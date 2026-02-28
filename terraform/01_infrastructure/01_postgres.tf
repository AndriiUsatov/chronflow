resource "kubernetes_service_v1" "chronflow-db-service" {
  metadata {
    name = "chronflow-db-service"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  spec {
    type = "ClusterIP"
    cluster_ip = "None"
    selector = {
        app = "postgres"
    }
    port {
      port = var.pg_port
    }
  }

}

resource "kubernetes_stateful_set_v1" "chronflow-postgres" {
  metadata {
    name = "postgres"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  spec {
    service_name = kubernetes_service_v1.chronflow-db-service.metadata[0].name
    replicas = "1"
    selector {
      match_labels = {
        app: "postgres"
      }
    }
    template {
      metadata {
        labels = {
          app: "postgres"
        }
      }
      spec {
        container {
          name = "postgres"
          image = "postgres:17-alpine"
          env {
            name = "POSTGRES_DB"
            value_from {
                config_map_key_ref {
                  name = "chronflow-postgres-config"
                  key = "CHRONFLOW_PG_TASK_DB"
                }
            } 
          }
          env {
            name = "POSTGRES_USER"
            value_from {
                secret_key_ref {
                    name =  "chronflow-postgres-creds"
                    key = "CHRONFLOW_PG_USER"
                }
            } 
          }
          env {
            name = "POSTGRES_PASSWORD"
            value_from {
                secret_key_ref {
                    name =  "chronflow-postgres-creds"
                    key = "CHRONFLOW_PG_PWD"
                }
            } 
          }
          port {
            container_port = var.pg_port
          }
          volume_mount {
            name = "postgres-data"
            mount_path = "/var/lib/postgresql/data"
          }
          volume_mount {
            name = "init-scripts"
            mount_path = "/docker-entrypoint-initdb.d"
          }
        }
        volume {
          name = "init-scripts"
          config_map {
            name = "postgres-init-script"
          }
        }
      }
    }
    volume_claim_template {
      metadata {
        name = "postgres-data"
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            "storage" = "1Gi"
          }
        }
      }
    }
  }
  depends_on = [ kubernetes_config_map_v1.postgres-init-script, kubernetes_config_map_v1.chronflow-postgres-config ]
}
