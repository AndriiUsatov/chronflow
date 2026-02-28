data "terraform_remote_state" "namespaces" {
  backend = "local"

  config = {
    path = "../00_init/terraform.tfstate"
  }
}

data "terraform_remote_state" "chronflow-deployments-output" {
  backend = "local"

  config = {
    path = "../02_chronflow/terraform.tfstate"
  }
}
