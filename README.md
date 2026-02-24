# Chronflow

## Prerequisites

## Quickstart
- Create namespace 
    ```kubectl create namespace chronflow```
- Add postgres credentials to secret storage:
    ``` 
    kubectl create secret generic chronflow-postgres-creds \
    --from-literal=CHRONFLOW_PG_USER=<user> \
    --from-literal=CHRONFLOW_PG_PWD=<password> 
    --namespace=chronflow
    ```
- Setup NATS Jetstream with Helm
    - Add package repository
        ```
        helm repo add nats https://nats-io.github.io/k8s/helm/charts/
        helm repo update
        ```
    - Install the NATS Server with File Storage (Persistence)
        ```
        helm upgrade --install nats nats/nats \
        --namespace chronflow \
        --set config.jetstream.enabled=true \
        --set config.jetstream.fileStore.enabled=true \
        --set config.cluster.enabled=true \
        --wait
        ```
    - Install the JetStream Controller (NACK)
        ```
        helm upgrade --install nack nats/nack \
        --namespace chronflow \
        --set jetstream.nats.url=nats://nats.chronflow.svc.cluster.local:4222 \
        --wait
        ```