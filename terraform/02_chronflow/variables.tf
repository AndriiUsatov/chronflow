variable "task_grpc_port" {
  type = number
  description = "Task API GRPC port"
  default = 50051
}

variable "task_api_port" {
  type = number
  description = "Task API port"
  default = 8080
}

variable "janitor_hearbeat_port" {
    type = number
    description = "Port number for janitor heatbeat"
    default = 8080
}

variable "scheduler_port" {
    type = number
    description = "Scheduler metrics port for prometheus"
    default = 8080
}

variable "worker_port" {
    type = number
    description = "Worker metrics port for prometheus"
    default = 8080
}

variable "nats_task_stream" {
  type = string
  description = "NATS task stream"
  default = "TASKS"
}

variable "nats_task_durable_consumer" {
    type = string
    description = "NATS task durable consumer group"
    default = "TASK_CONS"
}

variable "nats_task_subject" {
    type = string
    description = "NATS task subject"
    default = "tasks.process" 
}