Terraform provider for Rollbar
==============================

The Rollbar provider allows Terraform to control resources on
[Rollbar.com](https://rollbar.com), the Continuous Code Improvement Platform.

The provider allows you to manage your Rollbar projects, tokens, users, and
teams with ease. It must be configured with the proper credentials before it can
be used.


Status
------

[![Build & Test](https://github.com/rollbar/terraform-provider-rollbar/workflows/Build%20&%20Test/badge.svg)](https://github.com/rollbar/terraform-provider-rollbar/actions)
[![Coverage Status](https://coveralls.io/repos/github/rollbar/terraform-provider-rollbar/badge.svg)](https://coveralls.io/github/rollbar/terraform-provider-rollbar)
[![CodeQL](https://github.com/rollbar/terraform-provider-rollbar/workflows/CodeQL/badge.svg)](https://github.com/rollbar/terraform-provider-rollbar/actions?query=workflow%3ACodeQL)
[![ShiftLeft Scan](https://github.com/rollbar/terraform-provider-rollbar/workflows/ShiftLeft%20Scan/badge.svg)](https://github.com/rollbar/terraform-provider-rollbar/actions?query=workflow%3A%22ShiftLeft+Scan%22)
[![Maintainability](https://api.codeclimate.com/v1/badges/c5097d1a11f6f2310089/maintainability)](https://codeclimate.com/github/rollbar/terraform-provider-rollbar/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/c5097d1a11f6f2310089/test_coverage)](https://codeclimate.com/github/rollbar/terraform-provider-rollbar/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/rollbar/terraform-provider-rollbar)](https://goreportcard.com/report/github.com/rollbar/terraform-provider-rollbar)


Usage
-----

Example usage:

```hcl
# Install Rollbar provider from Terraform Registry
terraform {
  required_providers {
    rollbar = {
      source  = "rollbar/rollbar"
    }
  }
}

# Configure the Rollbar provider
provider "rollbar" {
  api_key = "YOUR_API_KEY"
}

# Create a team
resource "rollbar_team" "frontend" {
  name = "frontend-team"
}

# Create a project and assign the team
resource "rollbar_project" "frontend" {
  name         = "react-frontend"
  team_ids = [rollbar_team.frontend.id]
}
```

See [the docs](docs/index.md) for detailed usage information.


Debugging
---------

Enable writing debug log to `/tmp/terraform-provider-rollbar.log` by setting an
environment variable:

```
export TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1
terraform apply   # or any command that calls the Rollbar provider
```

This is necessary because Terraform providers aren’t _actually_ plugins - they
don’t get loaded into the running Terraform process.  Rather a provider is a
stand alone program that is started as a child processes and communicates with
Terraform via gRPC.  Anything that child process writes to stdout/stderr is
lost.  So if we want debug logging we must write to a file.


Development
-----------

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x, 0.13.x, or
  0.14.x
- [Go](https://golang.org/doc/install) 1.14.x, 1.15.x

See [`Quick Tests` workflow](.github/workflows/test.yml) for details of version compatibility testing.



### Building locally

Run `make build` to build the provider from source.

Run `make install` to build from source then install the provider locally for
usage by Terraform.

Folder [`./example`](./example) contains example Terraform configuration. To
use this config with a locally built provider, copy file `provider.tf.local` to
overwrite file `provider.tf`. Then run `terraform init` to initialize Terraform
with the provider.

Running `make plan`, `make apply`, or `make destroy` will:
* Build the provider from your working directory, and install for local
  Terraform.
* Setup Terraform configuration to use the freshly built provider
* Run `terraform <plan|apply|destroy>` in the `examples` folder with debug
  logging enabled.
* Display the logs on completion.


### Testing

This provider includes both unit tests, and acceptance tests that run with a
live Rollbar account.

To run the unit tests:

```shell
$ make test
```

To run the acceptance tests:

```shell
$ export TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1  # Enable debug output
$ make testacc
```


### Continuous Delivery

We use [semantic-release](semantic-release) for continuous delivery. When a PR
is merged into `master` the [`Semantic Release`
workflow](.github/workflows/release.yml) is triggered.  If relevant
[Conventional Commits](https://www.conventionalcommits.org/) annotations are
found in the Git log, semantic-release creates a new release and calls
[GoReleaser](https://goreleaser.com/) to build the binaries.

> This way no human is directly involved in the release process and the releases
are guaranteed to be [unromantic and
unsentimental](http://sentimentalversioning.org/).
>
> _-- from the semantic-release docs_

After a new release has been created by semantic-release, it may take up to 10
minutes for the provider binaries to be built and attached by GoReleaser.


### Terraform Versions

Several Makefile targets build the provider inside a Docker container, then
test it against different versions of Terraform. Environment variable
`ROLLBAR_API_KEY` must be set.

* `make terraform012` - Terraform 0.12.x
* `make terraform013` - Terraform 0.13.x
* `make terraform014` - Terraform 0.14.x


License
-------

This is Free Software, released under the terms of the [MIT license](LICENSE).


History
-------

Derived from
[babbel/terraform-provider-rollbar](https://github.com/babbel/terraform-provider-rollbar)
and
[jmcvetta/terraform-provider-rollbar](https://github.com/jmcvetta/terraform-provider-rollbar)


[latest-release]: https://github.com/rollbar/terraform-provider-rollbar/releases/latest
[requiring-providers]: https://www.terraform.io/docs/configuration/provider-requirements.html#requiring-providers
[semantic-release]: https://github.com/semantic-release/semantic-release
[pub-to-registry]: https://github.com/rollbar/terraform-provider-rollbar/issues/153