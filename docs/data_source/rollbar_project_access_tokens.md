`rollbar_project_access_tokens` Data Source
===========================================

Use this data source to retrieve information about all project access tokens
belonging to a Rollbar project.


Example Usage
-------------

To retrieve info about all projects:

```hcl-terraform
resource "rollbar_project" "test" {
  name = "foobar"
}

data "rollbar_project_access_tokens" "test" {
  project_id = rollbar_project.test.id
  prefix = "post_item"
}

output "tokens" {
  value = data.rollbar_project_access_tokens.test.access_tokens
}
```

Argument Reference
------------------

* `project_id` - (Required) ID of a Rollbar project
* `prefix` - (Optional) Project name begins with this prefix


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `access_tokens` - An array of Rollbar project access tokens.  Each item in the
  `access_tokens` block consists of the fields documented below.
  
Items in the `access_tokens` block have the following attributes:

* `access_token` - API token
* `project_id` - ID of the project that owns the token
* `cur_rate_limit_window_count` - Number of API hits that occurred in the
  current rate limit window
* `cur_rate_limit_window_start` - Time when the current rate limit window began
* `date_created` - Date the token was created
* `date_modified` - Date the token was last modified
* `name` - Name of the token
* `rate_limit_window_count` - Maximum allowed API hits during a rate limit
  window
* `rate_limit_window_size` - Duration of a rate limit window
* `scopes` - Project access scopes for the token.  Possible values are `read`,
  `write`, `post_server_item`, or `post_client_item`.
* `status` - Status of the token
