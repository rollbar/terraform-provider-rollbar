terraform {
    required_providers {
        rollbar = {
            source = "github.com/jmcvetta/rollbar"
            version = "~> 0.1"
        }
    }
}

provider "rollbar" {
  token = "ffe236fdbdcb452b9e31ba5af898f46a"
}

# Returns all projects
data "rollbar_projects" "all" {}
output "all_projects" {
  value = data.rollbar_projects.all.projects
}


