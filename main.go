/*
 * Copyright (c) 2020 Rollbar, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rollbar/terraform-provider-rollbar/rollbar"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	// Configure logging
	//log.Print("sdfsdf")
	//if os.Getenv("TERRAFORM_PROVIDER_ROLLBAR_DEBUG") == "1" {
	p := "/tmp/terraform-provider-rollbar.log"
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	//f := os.Stdout
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error opening log file")
	}
	defer f.Close() // #nosec
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: f}).
		With().Caller().
		Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	//}

	// Serve the plugin
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rollbar.Provider,
		// func() *schema.Provider {
		//	return rollbar.Provider()
		//},
	})
}
