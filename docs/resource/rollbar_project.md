`rollbar_project` Resource
=========================

Rollbar project resource.


Example Usage
-------------

```hcl
# Create a team
resource "rollbar_team" "foo" {
  name = "foo"
}

# Create a project and assign the team
resource "rollbar_project" "bar" {
  name         = "Bar"
  team_ids = [rollbar_team.foo.id]
}
```

Argument Reference
------------------

The following arguments are supported:

* `name` - (Required) Human readable name for the project
* `team_ids` - (Optional) IDs of teams assigned to the project


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the project
* `account_id` - ID of account that owns the project
* `date_created` - Date the project was created
* `date_modified` - Date the project was last modified
* `status` - Status of the project


Import
------

Projects can be imported using the project ID, e.g.

```
$ terraform import rollbar_project.foo 411703
```