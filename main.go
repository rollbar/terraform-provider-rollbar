package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/jmcvetta/terraform-provider-rollbar/rollbar"
	"github.com/rs/zerolog/log"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return rollbar.Provider()
		},
	})
}

func init() {
	// Configure logger to display file and line number
	log.Logger = log.With().Caller().Logger()
}
