Rollbar Provider
================

The Rollbar provider is used to interact with [Rollbar](https://rollbar.com)
resources.

The provider allows you to manage your Rollbar account's projects, members, and
teams easily. It needs to be configured with the proper credentials before it
can be used.


Example Usage
-------------

```hcl
provider "rollbar" {
  api_key = var.rollbar_token
}
```

Argument Reference
------------------

The following arguments are supported:

* `api_key` - (Required) This is the Rollbar authentication token. The value can
  be sourced from the ROLLBAR_TOKEN environment variable.


Data Sources
------------

* [`rollbar_project`](data_source/rollbar_project.md) - A Rollbar project
* [`rollbar_projects`](data_source/rollbar_projects.md) - List all Rollbar
  projects
* [`rollbar_project_access_token`](data_source/rollbar_project_access_token.md)
  - An access token belonging to a Rollbar project
* [`rollbar_project_access_tokens`](data_source/rollbar_project_access_tokens.md)
  - List all access tokens belonging to a Rollbar project


Resources
---------

* [`rollbar_project`](resource/rollbar_project.md) - A Rollbar project
* [`rollbar_project_access_token`](resource/rollbar_project_access_token.md) - A
  Rollbar project access token
* [`rollbar_user`](resource/rollbar_user.md) - A Rollbar user
