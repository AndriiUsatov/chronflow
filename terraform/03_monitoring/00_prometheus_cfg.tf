resource "kubernetes_manifest" "chronflow-api-pod-monitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "chronflow-api-pod-monitor"
      namespace = data.terraform_remote_state.namespaces.outputs.monitoring_namespace
      labels = {
        release = "kube-prometheus-stack"
      }
    }
    spec = {
      namespaceSelector = {
        matchNames = [data.terraform_remote_state.namespaces.outputs.chronflow_namespace]
      }
      selector = {
        matchLabels = {
          app = data.terraform_remote_state.chronflow-deployments-output.outputs.chronflow-api-deployment-pod-label-app
        }
      }
      podMetricsEndpoints = [
        {
          port = "rest"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "chronflow-scheduler-pod-monitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "chronflow-scheduler-pod-monitor"
      namespace = data.terraform_remote_state.namespaces.outputs.monitoring_namespace
      labels = {
        release = "kube-prometheus-stack"
      }
    }
    spec = {
      namespaceSelector = {
        matchNames = [data.terraform_remote_state.namespaces.outputs.chronflow_namespace]
      }
      selector = {
        matchLabels = {
          app = data.terraform_remote_state.chronflow-deployments-output.outputs.chronflow-scheduler-deployment-pod-label-app
        }
      }
      podMetricsEndpoints = [
        {
          port = "port"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "chronflow-janitor-pod-monitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "chronflow-janitor-pod-monitor"
      namespace = data.terraform_remote_state.namespaces.outputs.monitoring_namespace
      labels = {
        release = "kube-prometheus-stack"
      }
    }
    spec = {
      namespaceSelector = {
        matchNames = [data.terraform_remote_state.namespaces.outputs.chronflow_namespace]
      }
      selector = {
        matchLabels = {
          app = data.terraform_remote_state.chronflow-deployments-output.outputs.chronflow-janitor-deployment-pod-label-app
        }
      }
      podMetricsEndpoints = [
        {
          port = "port"
        }
      ]
    }
  }
}

resource "kubernetes_manifest" "chronflow-worker-pod-monitor" {
  manifest = {
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "PodMonitor"
    metadata = {
      name      = "chronflow-worker-pod-monitor"
      namespace = data.terraform_remote_state.namespaces.outputs.monitoring_namespace
      labels = {
        release = "kube-prometheus-stack"
      }
    }
    spec = {
      namespaceSelector = {
        matchNames = [data.terraform_remote_state.namespaces.outputs.chronflow_namespace]
      }
      selector = {
        matchLabels = {
          app = data.terraform_remote_state.chronflow-deployments-output.outputs.chronflow-worker-deployment-pod-label-app
        }
      }
      podMetricsEndpoints = [
        {
          port = "port"
        }
      ]
    }
  }
}
