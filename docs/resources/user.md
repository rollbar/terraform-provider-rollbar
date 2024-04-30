`rollbar_user` Resource
=========================

Rollbar user resource.


Example Usage
-------------

```hcl
# Create a team
resource "rollbar_team" "developers" {
  name = "developers"
}

# Create a user as a member of that team
resource "rollbar_user" "some_dev" {
  email = "some_dev@company.com"
  team_ids = [rollbar_team.developers.id]
}
```

Argument Reference
------------------

The following arguments are supported:
* `email` - (Required) The user's email address
* `team_ids` - (Required) IDs of the teams to which this user belongs


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `username` - The user's username
* `user_id` - The ID of the user
* `status` - Status of the user.  Either `invited` or `subscribed`


Import
------

Users can be imported using their user email address, e.g.

```
$ terraform import rollbar_user.some_dev some_dev@company.com
```
