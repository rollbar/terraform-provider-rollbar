package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-rollbar/rollbar"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rollbar.Provider})
}
