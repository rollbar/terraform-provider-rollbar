`rollbar_project_access_token` Data Source
===========================================

Use this data source to retrieve information about a project access token
belonging to a Rollbar project.


Example Usage
-------------

To retrieve info about a token:

```hcl
resource "rollbar_project" "test" {
  name = "foobar"
}

data "rollbar_project_access_token" "test" {
  project_id = rollbar_project.test.id
  name = "post_item_client"
  depends_on = [rollbar_project.test]
}

output "token" {
  value = data.rollbar_project_access_tokens.test
}
```

Argument Reference
------------------

* `project_id` - (Required) ID of a Rollbar project
* `name` - (Required) Name of the token


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `access_token` - API token
* `project_id` - ID of the project that owns the token
* `cur_rate_limit_window_count` - Number of API hits that occurred in the
  current rate limit window
* `cur_rate_limit_window_start` - Time when the current rate limit window began
* `date_created` - Date the token was created
* `date_modified` - Date the token was last modified
* `rate_limit_window_count` - Maximum allowed API hits during a rate limit
  window
* `rate_limit_window_size` - Duration of a rate limit window
* `scopes` - Project access scopes for the token.  Possible values are `read`,
  `write`, `post_server_item`, or `post_client_item`.
* `status` - Status of the token
