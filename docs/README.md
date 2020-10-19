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
  token = var.rollbar_token
}
```

Argument Reference
------------------

The following arguments are supported:

* `token` - (Required) This is the Rollbar authentication token. The value can be
  sourced from the ROLLBAR_TOKEN environment variable.


Data Sources
------------

* [`rollbar_projects`](data_source/rollbar_projects.md) - List all Rollbar
  projects
* [`rollbar_project_access_tokens`](data_source/rollbar_project_access_tokens.md)
  - List all access tokens belonging to a Rollbar


Resources
---------

* [`rollbar_project`](resource/rollbar_project.md) - A Rollbar project
* [`rollbar_project_access_token`](resource/rollbar_project_access_token.md) - A
  Rollbar project access token
