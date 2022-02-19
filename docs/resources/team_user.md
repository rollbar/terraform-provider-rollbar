`rollbar_team_user` Resource
=========================

Manage a single Rollbar team member assignment. If a registered Rollbar user exists, they will be
assigned to the team, otherwise they will be invited to join.


Example Usage
-------------

```hcl
# Create a team
resource "rollbar_team" "developers" {
  name = "developers"
}

# Create a project
resource "rollbar_team_user" "foo" {
  team_id = rollbar_team.developers.id
  email   = "some_dev@company.com"
}
```

!> **NOTE** When using this resource in conjunction with `rollbar_user` resource it is advisable to add the following `lifecycle` argument to prevent the teams being unassigned on subsequent runs:

```hcl
resource "rollbar_team" "devs" {
  //...
  lifecycle {
     ignore_changes = [team_ids]
  }
}
```

Argument Reference
------------------

The following arguments are supported:

* `team_id` - (Required) ID of the team to which this user belongs
* `email` - (Required) The user's email address


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `status` - Status of the user. Either `invited` or `registered`
* `user_id` - The ID of the user if status is `registered`
* `invite_id` - Invitation ID if status is `invited`

Import
------

Resource can be imported using the team ID and email address separated by a comma e.g.

```
$ terraform import rollbar_team_user.foo 689493,some_dev@company.com
```