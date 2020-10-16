terraform {
  required_providers {
    rollbar = {
      source  = "github.com/rollbar/rollbar"
      version = "~> 0.2"
    }
  }
}

variable "rollbar_token" {
  type = string
}

provider "rollbar" {
  api_key = var.rollbar_token
}

/*
 * Adapted from @jtsaito
 */

data "rollbar_project" "test" {
  name = rollbar_project.test.name
  depends_on = [rollbar_project.test]
}

data "rollbar_project_access_token" "post_client_item" {
  project_id = data.rollbar_project.test.id
  name       = "post_client_item"
}

data "rollbar_project_access_token" "post_server_item" {
  project_id = data.rollbar_project.test.id
  name       = "post_server_item"
}


/*
 * Added by @jmcvetta
 */

resource "rollbar_project" "test" {
  name = "tf-acc-test-syntax-compatibility"
}

resource "rollbar_project_access_token" "test" {
  name = "test-token"
  project_id = rollbar_project.test.id
  scopes = ["post_client_item"]
  depends_on = [rollbar_project.test]
}

data "rollbar_projects" "all" {}

data "rollbar_project_access_tokens" "test" {
  project_id = rollbar_project.test.id
  prefix = "post_"
}