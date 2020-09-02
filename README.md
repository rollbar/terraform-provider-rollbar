# Terraform provider for Rollbar

## License

This is proprietary software.  **No license whatsoever is granted to this
software without written contract.**


## Status

![Build](https://github.com/jmcvetta/terraform-provider-rollbar/workflows/Build/badge.svg)


## Debugging

Enable writing debug log to `/tmp/terraform-provider-rollbar.log` by setting an
environment variable:

```
export TERRAFORM_PROVIDER_ROLLBAR_DEBUG=1
terraform apply   # or any command that calls the Rollbar provider
```
