resource "helm_release" "nats" {
    name = "nats"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    repository = "https://nats-io.github.io/k8s/helm/charts/"
    chart = "nats"
    set = [ 
        {
            name = "config.jetstream.enabled"
            value = "true"
        } ,
        {
            name = "config.jetstream.fileStore.enabled"
            value = "true"
        },
        {
            name = "config.cluster.enabled"
            value = "true"
        }
    ]
}

resource "helm_release" "nack" {
    name = "nack"
    namespace = data.terraform_remote_state.namespaces.outputs.chronflow_namespace
    repository = "https://nats-io.github.io/k8s/helm/charts/"
    chart = "nack"
    set = [ 
        {
            name = "jetstream.nats.url"
            value = "nats://nats.chronflow.svc.cluster.local:4222"
        } 
    ]
    depends_on = [ helm_release.nats ]
}

