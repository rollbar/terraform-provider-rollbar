/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package rollbar

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

var testAccProviders map[string]*schema.Provider
var testAccProviderFactories func(providers *[]*schema.Provider) map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var testAccProviderFunc func() *schema.Provider

func init() {
	// Log to console
	log.Logger = log.
		With().Caller().
		Logger()
	if os.Getenv("TERRAFORM_PROVIDER_ROLLBAR_DEBUG") == "1" {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Setup testing
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"rollbar": testAccProvider,
	}
	//testAccProviderFactories = func(providers *[]*schema.Provider) map[string]func() (*schema.Provider, error) {
	//	// this is an SDKV2 compatible hack, the "factory" functions are
	//	// effectively singletons for the lifecycle of a resource.Test
	//	var providerNames = []string{"aws", "awseast", "awswest", "awsalternate", "awsus-east-1", "awsalternateaccountalternateregion", "awsalternateaccountsameregion", "awssameaccountalternateregion", "awsthird"}
	//	var factories = make(map[string]func() (*schema.Provider, error), len(providerNames))
	//	for _, name := range providerNames {
	//		p := Provider()
	//		factories[name] = func() (*schema.Provider, error) { //nolint:unparam
	//			return p, nil
	//		}
	//		*providers = append(*providers, p)
	//	}
	//	return factories
	//}
	testAccProviderFunc = func() *schema.Provider { return testAccProvider }
}
func testAccPreCheck(t *testing.T) {
	// FIXME: Add preflight check for API credentials
	log.Warn().Msg("Need to add preflight check for credentials")

	//if token := os.Getenv("HASHICUPS_USERNAME"); token == "" {
	//	t.Fatal("HASHICUPS_USERNAME must be set for acceptance tests")
	//}
}
