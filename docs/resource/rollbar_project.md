`rollbar_project` Resource
=========================

Rollbar project resource.


Example Usage
-------------

```hcl
# Create a project
resource "rollbar_project" "foo" {
  name         = "Foo"
}
```

Argument Reference
------------------

The following arguments are supported:

* `name` - (Required) The human readable name for the project.


Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the created project.
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