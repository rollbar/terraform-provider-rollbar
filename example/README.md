Example Terraform Configuration
===============================

This folder contains example Terraform configuration, demonstrating the usage
and capabilities of `terraform-provider-rollbar`.


-----

Clone entire sales-engineering repo

cd TerraForm dir

in your account under 'Account Access Tokens' create a new token and grant it 'write access' you don't have to do this but i did in this example

Then run

cd demo

install terraform (using brew is the easiest way assuming you have brew installed)

brew install terraform

export ROLLBAR_API_KEY=<yourNewToken>

terraform init

terraform plan

terraform apply (enter yes if you like the plan)

See in RB UI

To delete

terraform destroy (enter yes to confirm)