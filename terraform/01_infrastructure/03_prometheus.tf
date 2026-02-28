resource "helm_release" "kube-prometheus-stack" {
  name       = "kube-prometheus-stack"
  namespace = data.terraform_remote_state.namespaces.outputs.monitoring_namespace
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"

  timeout = 900

  wait = true

  set = [
    {
      name  = "grafana.enabled"
      value = "true"
    },
    {
      name  = "alertmanager.enabled"
      value = "false"
    },
    {
      name  = "grafana.sidecar.dashboards.enabled"
      value = "false"
    },
    {
      name  = "grafana.admin.existingSecret"
      value = "grafana-admin-creds"
    },
    {
      name  = "grafana.admin.userKey"
      value = "ADMIN_USER"
    },
    {
      name  = "grafana.admin.passwordKey"
      value = "ADMIN_PASSWORD"
    },
    {
      name  = "grafana.service.type" 
      value = "LoadBalancer"
    },
    {
      name  = "grafana.service.port" 
      value = tostring(var.grafana_exposed_on_port)
    }
  ]
}
