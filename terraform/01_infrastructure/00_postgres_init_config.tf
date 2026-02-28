resource "kubernetes_config_map_v1" "postgres-init-script" {
  metadata {
    name = "postgres-init-script"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  data = {
    "01-schema.sql": <<-EOT
        CREATE SCHEMA chronflow;
    EOT

    "02-task-table.sql" = <<-EOT
        CREATE TABLE chronflow.task (
        	id UUID PRIMARY KEY,
        	url text NOT NULL,
        	method text NOT NULL,
        	headers jsonb,
        	body bytea,
        	status smallint NOT NULL DEFAULT 0,
        	run_at timestamptz NOT NULL,
        	created timestamptz NOT NULL DEFAULT current_timestamp,
        	updated timestamptz NOT NULL DEFAULT current_timestamp,
        	retry_count integer DEFAULT 0,
        	error_message TEXT DEFAULT NULL,
            CONSTRAINT check_method CHECK (method IN ('GET', 'POST', 'PUT', 'PATCH', 'DELETE'))
        );
    EOT

    "03-task-table-index.sql" = <<-EOT
        CREATE INDEX idx_dispatch_ready ON chronflow.task (run_at) WHERE status = 0;
        CREATE INDEX idx_task_monitoring ON chronflow.task (status, run_at);
    EOT

    "04-notify-func.sql" = <<-EOT
        CREATE FUNCTION chronflow.notify_new_task() RETURNS trigger AS $$
        BEGIN
          PERFORM pg_notify('new_task_event', '');
          RETURN NEW;
        END;
        $$ LANGUAGE plpgsql;
    EOT

    "05-trigger.sql" = <<-EOT
        CREATE TRIGGER task_inserted
        AFTER INSERT ON chronflow.task
        FOR EACH ROW EXECUTE FUNCTION chronflow.notify_new_task();
    EOT

  }
}

resource "kubernetes_config_map_v1" "chronflow-postgres-config" {
  metadata {
    name = "chronflow-postgres-config"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
  }

  data = {
    CHRONFLOW_PG_URI                = "chronflow-db-service.chronflow.svc.cluster.local"
    CHRONFLOW_PG_PORT               = tostring(var.pg_port)
    CHRONFLOW_PG_TASK_DB            = "chronflow"
    CHRONFLOW_PG_SSL_MODE           = "disable"
    CHRONFLOW_PG_TASK_SCHEMA        = "chronflow"
    CHRONFLOW_PG_TASK_TABLE         = "task"
    CHRONFLOW_PG_NOTIFICATION_EVENT = "new_task_event"
  }
}
