Example Terraform Configuration
===============================

This folder contains example Terraform configuration, demonstrating the usage
and capabilities of `terraform-provider-rollbar`.


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


export ROLLBAR_API_KEY=<yourNewToken>

terraform init

terraform plan

terraform apply (enter yes if you like the plan)

See in RB UI

To delete

terraform destroy (enter yes to confirm)