output "chronflow_namespace" {
  value = kubernetes_namespace_v1.chronflow_namespace.metadata[0].name
}

output "monitoring_namespace" {
  value = kubernetes_namespace_v1.monitoring_namespace.metadata[0].name
}