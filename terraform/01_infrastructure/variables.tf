variable "pg_port" {
    type = number
    description = "Postgres port"
    default = 5432
}

variable "grafana_exposed_on_port" {
    type = number
    default = 8081
}