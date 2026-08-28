`rollbar_service_link` Resource
=========================

Rollbar projects can be configured with service links, dynamically constructed links that use templated fields from your Rollbar items to provide better context.

This resource can manage service links.  See the following api documentation for more details:

* [Rollbar API Service Links](https://docs.rollbar.com/reference/service-links)

Example Usage
-------------

**Option 1: Using a project access token**

```hcl
provider "rollbar" {
    api_key         = "my-account-access-token"  # optional for this resource
    project_api_key = "my-project-access-token"
}

resource "rollbar_service_link" "service_link" {
  name = "service_link_name"
  template = "https://some-service.xyz/commit/{{code_version}}"
}
```

**Option 2: Using an account access token with `project_id`**

If you prefer not to configure a separate project access token, you can use the account access token and specify `project_id` on each resource:

```hcl
provider "rollbar" {
    api_key = "my-account-access-token"
}

resource "rollbar_service_link" "service_link" {
  project_id = 123456
  name       = "service_link_name"
  template   = "https://some-service.xyz/commit/{{code_version}}"
}
```

Argument Reference
------------------

The following arguments are supported:

* `project_id` - (Optional) The Rollbar project ID. When set, the account access token (`api_key`) is used with this project ID instead of requiring a separate project access token (`project_api_key`).
* `name` - (Required) The name of the service link
* `template` - (Required) The url that contains templated variables referencing an occurrences data. [Examples](https://docs.rollbar.com/docs/service-links)


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the service link
