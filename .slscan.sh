#!/bin/sh
# This script invokes Shiftleft Scan using docker run command
{ # Prevent execution if this script was only partially downloaded
scan() {
    docker run --rm -e "WORKSPACE=$(pwd)" -e GITHUB_TOKEN -e "SCAN_AUTO_BUILD=true" -v "$(pwd):/app" shiftleft/scan scan $*
}
scan
} # End of wrapping
