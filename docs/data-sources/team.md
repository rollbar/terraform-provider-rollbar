`rollbar_team` Data Source
==============================

Use this data source to retrieve information about a Rollbar team.


Example Usage
-------------

To retrieve info about a team by name or ID:

```hcl
data "rollbar_team" "foobar" {
  name = "foobar"
}

data "rollbar_team" "example" {
  team_id = 123456
}
```


Argument Reference
------------------

The following arguments are supported:

* `team_id` - (Optional) Rollbar team ID.
* `name` - (Optional) Human readable name for the team. Conflicts with `team_id`.

One of `team_id` or `name` must be specified.

Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the team
* `account_id` - ID of account that owns the team
* `access_level` - Team access level. Will be one of `standard`, `light` or `view`.
