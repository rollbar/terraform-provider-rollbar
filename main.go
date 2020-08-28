/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/jmcvetta/terraform-provider-rollbar/rollbar"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	// Configure logging
	f, err := os.OpenFile("/tmp/rollbar-terraform.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("Error opening log file")
	}
	defer f.Close()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: f}).
		With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Start plugin
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return rollbar.Provider()
		},
	})
}
