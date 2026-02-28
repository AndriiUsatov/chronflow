data "terraform_remote_state" "namespaces" {
  backend = "local"

  config = {
    path = "../00_init/terraform.tfstate"
  }
}

data "terraform_remote_state" "infrastructureme_out" {
  backend = "local"

  config = {
    path = "../01_infrastructure/terraform.tfstate"
  }
}
