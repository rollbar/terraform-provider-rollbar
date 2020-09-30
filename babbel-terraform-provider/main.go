package main

import (
	"github.com/babbel/terraform-provider-rollbar/rollbar"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rollbar.Provider})
}
