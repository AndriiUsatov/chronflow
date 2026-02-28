resource "kubernetes_config_map_v1" "chronflow-grpc-config" {
  metadata {
    name      = "chronflow-grpc-config"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  data = {
    CHRONFLOW_GRPC_TASK_UPDATE_SERVER_TRANSP = "tcp"
    CHRONFLOW_GRPC_TASK_UPDATE_SERVER_URL    = "dns:///chronflow-api-grpc-service.chronflow.svc.cluster.local"
    CHRONFLOW_GRPC_TASK_UPDATE_SERVER_PORT   = tostring(var.task_grpc_port)
  }
}

resource "kubernetes_config_map_v1" "chronflow-nats-config" {
  metadata {
    name      = "chronflow-nats-config"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  data = {
    CHRONFLOW_NATS_URL                     = "nats://nats.chronflow.svc.cluster.local:4222"
    CHRONFLOW_NATS_TASK_STREAM             = var.nats_task_stream
    CHRONFLOW_NATS_TASK_TO_PROCESS_SUBJECT = var.nats_task_subject
    CHRONFLOW_NATS_TASK_TO_PROCESS_DURABLE = var.nats_task_durable_consumer
  }
}

resource "kubernetes_config_map_v1" "chronflow-api-config" {
  metadata {
    name      = "chronflow-api-config"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  data = {
    "CHRONFLOW_TASK_API_PORT" : tostring(var.task_api_port)
  }
}

resource "kubernetes_config_map_v1" "chronflow-janitor-config" {
  metadata {
    name      = "chronflow-janitor-config"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  data = {
    "CHRONFLOW_JANITOR_HEARTBEAT_PORT" : tostring(var.janitor_hearbeat_port)
  }
}

resource "kubernetes_config_map_v1" "chronflow-scheduler-config" {
  metadata {
    name      = "chronflow-scheduler-config"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  data = {
    "CHRONFLOW_SCHEDULER_PORT" : tostring(var.scheduler_port)
  }
}

resource "kubernetes_config_map_v1" "chronflow-worker-config" {
  metadata {
    name      = "chronflow-worker-config"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  data = {
    "CHRONFLOW_WORKER_PORT" : tostring(var.worker_port)
  }
}

