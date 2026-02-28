# Chronflow
A real-time, event-driven HTTP request orchestrator.

ChronFlow is a service built in Go for scheduling and executing resilient HTTP calls. 
It sits between your application and external APIs, ensuring that webhooks and service-to-service requests are delivered reliably, even under heavy load or during downstream outages.

## Key components
- **REST API**: The entry point for users to submit and schedule HTTP requests.
- **Scheduler**: The clock of the system. It ensures tasks move from the database to NATS only when their scheduled time arrives.
- **NATS**: The high-speed messaging backbone that distributes "due" tasks to available workers.
- **Worker**: The execution engine that handles the outbound HTTP logic.
- **Internal gRPC**: The private channel for Workers to report success or failure back to the API.
- **Janitor**: The watchdog that ensures no task remains in a "stuck" state indefinitely.

## Core loop
1. **Ingestion**: The API receives a "Task" (HTTP request blueprint: URL, Method, Payload) via a standard REST interface.
2. **Scheduling**: The Scheduler monitors the database. When a task reaches its execution time, the Scheduler "dispatches" it by publishing a message to NATS.
3. **Execution**: A Worker subscribed to NATS picks up the message and executes the outbound HTTP call to the target destination.
4. **Reporting**: Once the call is finished, the Worker uses gRPC to notify the API of the result. The API then updates the final status in the database.
5. **Recovery**: The Janitor acts as the safety net. It identifies tasks that are stuck (e.g., a Worker crashed before it could report back via gRPC) and resets them for re-execution.

## Prerequisites
- [Golang version >= 1.25.5](https://go.dev/doc/install)
- K8s cluster config
- [kubectl CLI](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [TaskFile](https://taskfile.dev/docs/installation)
- [Terraform](https://developer.hashicorp.com/terraform/install)
- [Protobuf compiler](https://github.com/protocolbuffers/protobuf#protobuf-compiler-installation)
- [Docker](https://docs.docker.com/get-started/)
- Protobuf go & go-grpc plugins
  ```
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```
## Deploy
- Build images
    ```
    task image
    ```
- Create namespaces: **chronflow** and **monitoring**
    ```
    task tf:deploy:namespace
    ```
- Add postgres credentials secrets:
    ``` 
    kubectl create secret generic chronflow-postgres-creds \
    --from-literal=CHRONFLOW_PG_USER=<user> \
    --from-literal=CHRONFLOW_PG_PWD=<password> \ 
    --namespace=chronflow
    ```
- Add Grafana admin credentials secret
    ```
    kubectl create secret generic grafana-admin-creds \ 
    --from-literal=ADMIN_USER=<admin-user> \
    --from-literal=ADMIN_PASSWORD=<admin-password> \ 
    --namespace monitoring
    ```
- **OPTIONAL** - Add metrics server for HPA. Verify: `kubectl top nodes`
    ```
    kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
    
    kubectl patch deployment metrics-server -n kube-system \ 
    --type='json' \ 
    -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
    ```
- Deploy infrastructure
    ```
    task tf:deploy:infrastructure
    ```
- Add Chronflow deployments
    ```
    task tf:deploy:chronflow
    ```
- Monitoring:
    - Set Grafana URL env variable:
        ```
        export TF_VAR_grafana_url="grafana-url"
        ```
    - Set Grafana API token
        ```
        export TF_VAR_grafana_api_token="your-api-token"
        ```
    - Prometheus PodMonitors and Grafana Task Dashboard
        ```
        task tf:deploy:monitoring
        ```
## Quickstart

1. Build images
    ```
    task image
    ```
2. Run `tf:deploy`
    ```
    task tf:deploy \
    PG_USER="<pg_admin_user>" \
    PG_PASSWORD="<pg_admin_password>" \ 
    GRAFANA_USER="<grafana_admin_user>" \
    GRAFANA_PASSWORD="<grafana_admin_password>" \
    GRAFANA_URL="http://<host>:8081"
    ```
## Results
- Task API exposed on port `:80` (Base URL: `/api/v1`)
- To access swagger go to `http://<host>:80/`
- Grafana exposed on port `http://<host>:8081`
