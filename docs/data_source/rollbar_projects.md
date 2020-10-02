`rollbar_projects` Data Source
==============================

Use this data source to retrieve information about all Rollbar projects you can
access.


Example Usage
-------------

To retrieve info about all projects:

```hcl
data "rollbar_projects" "all" {}

output "all_projects" {
  value = data.rollbar_projects.all.projects
}
```

Argument Reference
------------------

This data source accepts no arguments.


Attribute Reference
-------------------

* `id` - ID of project
* `name` - Name of project
* `account_id` - ID of account that owns the project
* `date_created` - Date the project was created
* `date_modified` - Date the project was last modified
* `status` - Status of the project
