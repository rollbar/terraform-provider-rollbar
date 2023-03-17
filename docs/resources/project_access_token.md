`rollbar_project_access_token` Resource
=========================

Rollbar project access token resource.


Example Usage
-------------

```hcl
# Create a project
resource "rollbar_project" "foo" {
  name         = "Foo"
}

# Create an access token for the project
resource "rollbar_project_access_token" "bar" {
  name = "bar"
  project_id = rollbar_project.foo.id
  scopes = ["read", "post_server_item"]
  
  depends_on = [rollbar_project.foo]
}
```

Argument Reference
------------------

The following arguments are supported:

* `name` - (Required) The human readable name for the token.
* `project_Id` - (Required) ID of the Rollbar project to which this token
  belongs.
* `scopes` - (Required) List of access [scopes](https://docs.rollbar.com/#section/Authentication/Project-access-tokens) 
  granted to the token.  Possible values are `read`, `write`,
  `post_server_item`, and `post_client_server`.
* `status` - (Optional) Status of the token.  Possible values are `enabled` 
  and `disabled`.
* `rate_limit_window_count` - (Optional) Total number of calls allowed within
  the rate limit window
* `rate_limit_window_size` - (Optional) Total number of seconds that makes up
  the rate limit window


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `access_token` - Access token for Rollbar API
* `date_created` - Date the project was created
* `date_modified` - Date the project was last modified
* `cur_rate_limit_window_count` - Count of calls in the current window
* `cur_rate_limit_window_start` - Time when the current window began


Import
------

Projects can be imported using a combination of the `project_id` and
`access_token` joined by a `/`, e.g.

```
$ terraform import rollbar_project_access_token.baz 411703/d19f7ada16534b1c94e91d9da3dbae5a
```
