resource "kubernetes_manifest" "task-stream" {
  manifest = {
    apiVersion = "jetstream.nats.io/v1beta2"
    kind       = "Stream"
    metadata = {
      name      = "task-stream"
      namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    }
    spec = {
      name     = var.nats_task_stream
      subjects = [var.nats_task_subject]
      storage  = "file"
      maxAge   = "12h"
    }
  }
}

resource "kubernetes_manifest" "tasks-consumer" {
  manifest = {
    apiVersion = "jetstream.nats.io/v1beta2"
    kind       = "Consumer"
    metadata = {
      name      = "tasks-consumer"
      namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    }
    spec = {
      streamName  = var.nats_task_stream
      durableName = var.nats_task_durable_consumer
      deliverPolicy = "all"
      maxDeliver = 20
      ackPolicy = "explicit"
      filterSubject = var.nats_task_subject
    }
  }
  depends_on = [ kubernetes_manifest.task-stream ]
}
