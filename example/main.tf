terraform {
  required_providers {
    rollbar = {
      source  = "github.com/jmcvetta/rollbar"
      version = "~> 0.1"
    }
  }
}

variable "rollbar_token" {
  type = string
}

provider "rollbar" {
  token = var.rollbar_token
}


# Returns all projects
//data "rollbar_projects" "all" {}
//output "all_projects" {
//  value = data.rollbar_projects.all.projects
//}


resource "rollbar_project" "foo" {
  name = "Foo"
}

resource "rollbar_project" "bar" {
  name = "Bar"
}