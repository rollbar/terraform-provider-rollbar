#
# Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
# reserved.
#
# NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
# between author and licensee.
#

#===============================================================================
#
# Convenience script to build & cleanly install provider; run `terraform apply`
# on the example Terraform configs; and cat the debug log output.
#
#===============================================================================

set -e 
set -x

# Clear the screen but not the scrollback buffer
clear -x

# Cleanup last run
rm -vrf .terraform /tmp/rollbar-terraform.log

# Build & install the latest provider
(cd .. && make)

# Initialize terraform
terraform init

# Test the provider
terraform apply

# Print the debug log
cat /tmp/rollbar-terraform.log 
