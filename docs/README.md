Rollbar Provider
================

The Rollbar provider is used to interact with Rollbar resources.

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

WIP

Resources
---------

* [`rollbar_project`](resources/rollbar_project.md) - A Rollbar project