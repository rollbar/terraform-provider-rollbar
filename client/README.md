Rollbar Management API Client
=============================

A client for the Rollbar API, providing access to project & team management
functionality. Used internally by the Rollbar Terraform provider.  Does **not**
provide error reporting functionality. If you want to use Rollbar to collect
errors from your application, this is the wrong client.

See [`rollbar-go`](https://github.com/rollbar/rollbar-go) for the official Go
client for Rollbar, which *does* provide support for error reporting.

[![Documentation](https://godoc.org/github.com/rollbar/terraform-provider-rollbar?status.svg)](http://godoc.org/github.com/rollbar/terraform-provider-rollbar)
