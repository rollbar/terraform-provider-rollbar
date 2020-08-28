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

set -x

clear -x  # Clear the screen but not the scrollback buffer

# Build & cleanly install the latest provider
(cd .. && make) \
    && rm -vrf .terraform /tmp/rollbar-terraform.log \
    && terraform init \
    && terraform apply

# Print the debug log
cat /tmp/rollbar-terraform.log 
