data "terraform_remote_state" "namespaces" {
  backend = "local"

  config = {
    path = "../00_init/terraform.tfstate"
  }
}