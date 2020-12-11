`rollbar_team` Resource
=========================

Rollbar team resource.


Example Usage
-------------

```hcl
# Create a project
resource "rollbar_team" "foo" {
  name         = "Foo"
  access_level = "standard"
}
```

Argument Reference
------------------

The following arguments are supported:

* `name` - (Required) Human readable name for the team
* `access_level` - (Optional) The team's access level.  Must be "standard",
  "light", or "view". Defaults to "standard".


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `account_id` - ID of account that owns the team


Import
------

Teams can be imported using the team ID, e.g.

```
$ terraform import rollbar_team.foo 689493
```