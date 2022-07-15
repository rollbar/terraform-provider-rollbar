`rollbar_service_link` Resource
=========================

Rollbar projects can be configured with service links, dynamically constructed links that use templated fields from your Rollbar items to provide better context.

This resource can manage service links.  See the following api documentation for more details:

* [Rollbar API Service Links](https://explorer.docs.rollbar.com/#tag/Service-Links)

Example Usage
-------------

```hcl
# Set the rollbar provider to manage a project
#
# NOTE: the account access token `api_key` is not required
# for managing project however it should be set if you want
# to use the provider for managing account-level resources
# in addition to project-level resources

provider "rollbar" {
    api_key         = "my-account-access-token"  # optional for this resource
    project_api_key = "my-project-access-token"
}

# Create a service links for a project
#

resource "rollbar_service_link" "service_link" {
  name = "service_link_name"
  template = "https://some-service.xyz/commit/{{code_version}}"
}
```

Argument Reference
------------------

The following arguments are supported:

* `name` - (Required) The name of the service link
* `template` - (Required) The url that contains templated variables referencing an occurrences data. [Examples](https://docs.rollbar.com/docs/service-links)


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the service link
