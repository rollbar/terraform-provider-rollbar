Terraform provider for Rollbar
==============================

The Rollbar provider is used to interact with Rollbar resources.

The provider allows you to manage your Rollbar account's projects, members, and
teams easily. It needs to be configured with the proper credentials before it
can be used.



Status
------

[![Build & Test](https://github.com/rollbar/terraform-provider-rollbar/workflows/Build%20&%20Test/badge.svg)](https://github.com/rollbar/terraform-provider-rollbar/actions)
[![Coverage Status](https://coveralls.io/repos/github/rollbar/terraform-provider-rollbar/badge.svg?branch=merge-jmcvetta)](https://coveralls.io/github/rollbar/terraform-provider-rollbar?branch=merge-jmcvetta)
[![CodeQL](https://github.com/rollbar/terraform-provider-rollbar/workflows/CodeQL/badge.svg)](https://github.com/rollbar/terraform-provider-rollbar/actions?query=workflow%3ACodeQL)
[![ShiftLeft Scan](https://github.com/rollbar/terraform-provider-rollbar/workflows/ShiftLeft%20Scan/badge.svg)](https://github.com/rollbar/terraform-provider-rollbar/actions?query=workflow%3A%22ShiftLeft+Scan%22)
[![Maintainability](https://api.codeclimate.com/v1/badges/c5097d1a11f6f2310089/maintainability)](https://codeclimate.com/github/rollbar/terraform-provider-rollbar/maintainability)



Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x
- [Go](https://golang.org/doc/install) 1.14.x


Usage
-----

[See the docs for usage information.](docs/README.md)


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


### Dev Script

Running `make dev` will:
* Build and install the provider 
* Run `terraform apply` in the `examples` folder with debug logging enabled
* Display the logs on completion.


License
-------

This is Free Software, released under the terms of the [MIT license](LICENSE).


History
-------

Derived from
[jmcvetta/terraform-provider-rollbar-jmcvetta](https://github.com/jmcvetta/terraform-provider-rollbar-jmcvetta)
and
[babbel/terraform-provider-rollbar](https://github.com/babbel/terraform-provider-rollbar)
