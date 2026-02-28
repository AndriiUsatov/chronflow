output "chronflow-api-deployment-pod-label-app" {
  value = kubernetes_deployment_v1.chronflow-api-deployment.spec[0].template[0].metadata[0].labels["app"]
}

output "chronflow-scheduler-deployment-pod-label-app" {
  value = kubernetes_deployment_v1.chronflow-scheduler-deployment.spec[0].template[0].metadata[0].labels["app"]
}

output "chronflow-janitor-deployment-pod-label-app" {
  value = kubernetes_deployment_v1.chronflow-janitor-deployment.spec[0].template[0].metadata[0].labels["app"]
}

output "chronflow-worker-deployment-pod-label-app" {
  value = kubernetes_deployment_v1.chronflow-worker-deployment.spec[0].template[0].metadata[0].labels["app"]
}