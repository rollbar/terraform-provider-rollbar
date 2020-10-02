`rollbar_project` Resource
=========================

Rollbar Project resource.


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

The following attributes are exported:

* `id` - The ID of the created project.
  * **TODO: is ID exported by default?**
