Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.13.x+ (to build the provider plugin)

Building The Provider
---------------------
Clone repository outside your Go path (not to `$GOPATH/src/github.com/babbel/terraform-provider-rollbar`)
or set `GO111MODULE=on` (cf. [Go module documenation](https://github.com/golang/go/wiki/Modules#daily-workflow)).

```sh
$ git clone git@github.com:babbel/terraform-provider-rollbar
```

Enter the directory and build the provider

```sh
$ make build-darwin
```

or

```sh
$ make build-linux
```

Using the provider
----------------------

```hcl
provider "rollbar" {
  api_key = "${var.api_key}"
}

resource "rollbar_user" "test" {
  team_id = 333290
  email   = "test@somewhere.com"
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13.x+ is *required*).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

We cannot have the acceptance tests until rollbar changes/improves their api.
The reason is because creating an invitation doesn't yield an userid.
A user id is created when a user accepts their invitation.

Github Releases
---------------------------
In order to push a release to Github the feature branch has to merged into master and then a tag needs to be created with the version name of the provider e.g. **v0.0.1** and pushed.

```sh
git checkout master
git pull origin master
git tag v<semver> -m "release comment"
git push origin master --tags
```
