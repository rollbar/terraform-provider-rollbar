Example Terraform Configuration
===============================

This folder contains example Terraform configuration, demonstrating the usage
and capabilities of `terraform-provider-rollbar` at the account level.


Usage Instructions
------------------

The following are step-by-step instructions for demonstrating
`terraform-provider-rollbar` using the example Terraform configuration files.

-----

1. First, change directories to the repo checkout:

   ```shell
   cd terraform-provider-rollbar
   ```

1. In your Rollbar account, under _Account Access Tokens_, create a new token and
   grant it _write access_.

1. Install terraform. On Mac, Brew is the easiest way (assuming you have Brew
   installed).  See [terraform installation
   docs](https://learn.hashicorp.com/tutorials/terraform/install-cli) for
   instructions on all supported platforms.

   ```shell
   brew install terraform
   ```

1. Change directories into example folder.

   ```shell
   cd example
   ```

1. Make your Rollbar API key available as an environment variable.

   ```shell
   export ROLLBAR_API_KEY=<yourNewToken>
   ```

1. Initialize Terraform - the Rollbar provider will be automatically downloaded and installed.

   ```shell
   terraform init
   ```

1. Examine Terraform's plan to create resources on Rollbar:

   ```shell
   terraform plan
   ```

1. Apply the plan - meaning Terraform will create the resources it described in the plan.

   ```shell
   terraform apply  # enter yes if you like the plan
   ```

1. Check changes in Rollbar web UI.


-----

To delete the resources created by Terraform, run:

```shell
terraform destroy  # enter yes to confirm

```
