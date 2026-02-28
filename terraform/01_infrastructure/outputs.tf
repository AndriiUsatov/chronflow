output "chronflow-postgres-config-name" {
  value = kubernetes_config_map_v1.chronflow-postgres-config.metadata[0].name
}