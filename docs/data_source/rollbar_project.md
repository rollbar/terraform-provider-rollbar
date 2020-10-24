`rollbar_project` Data Source
==============================

Use this data source to retrieve information about a Rollbar project.


Example Usage
-------------

To retrieve info about a project:

```hcl
data "rollbar_project" "foobar" {
  name = "foobar"
}

output "project_scopes" {
  value = data.rollbar_project.foobar.scopes
}
```


Argument Reference
------------------

The following arguments are supported:

* `name` - (Required) The human readable name for the project.


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - ID of project
* `account_id` - ID of account that owns the project
* `date_created` - Date the project was created
* `date_modified` - Date the project was last modified
* `status` - Status of the project
