terraform {
  required_providers {
    rollbar = {
      source  = "github.com/rollbar/rollbar"
      version = "~> 0.1"
    }
  }
}

variable "rollbar_token" {
  type = string
}

provider "rollbar" {
  version = "0.2.0"

  api_key = var.rollbar_token
}

data "rollbar_project" "example-spa" {
  name = "example"
}

data "rollbar_project_access_token" "example-spa-post-client-item" {
  project_id = data.rollbar_project.example-spa.id
  name       = "post_client_item"
}

data "rollbar_project_access_token" "example-spa-post-server-item" {
  project_id = data.rollbar_project.example-spa.id
  name       = "post_server_item"
}

