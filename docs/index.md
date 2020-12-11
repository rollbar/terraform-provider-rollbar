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

* `api_key` - (Required) Rollbar API authentication token. Value will be
  sourced from environment variable `ROLLBAR_API_KEY` if set.
* `api_url` - (Optional) Base URL for the Rollbar API.  Defaults to
  https://api.rollbar.com.  Value will be sourced from environment variable
  `ROLLBAR_API_URL` if set.


Data Sources
------------

* [`rollbar_project`](data-sources/project.md) - A Rollbar project
* [`rollbar_projects`](data-sources/projects.md) - List all Rollbar
  projects
* [`rollbar_project_access_token`](data-sources/project_access_token.md)
  - An access token belonging to a Rollbar project
* [`rollbar_project_access_tokens`](data-sources/project_access_tokens.md)
  - List all access tokens belonging to a Rollbar project


Resources
---------

* [`rollbar_project`](resources/project.md) - A Rollbar project
* [`rollbar_project_access_token`](resources/project_access_token.md) - A
  Rollbar project access token
* [`rollbar_user`](resources/user.md) - A Rollbar user
* [`rollbar_team`](resources/team.md) - A Rollbar team
