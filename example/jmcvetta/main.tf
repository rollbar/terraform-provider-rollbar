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

data "rollbar_project" "foo" {
  name = "Foo"
  depends_on = [rollbar_project.foo]
}

output "project_foo" {
  value = data.rollbar_project.foo
}