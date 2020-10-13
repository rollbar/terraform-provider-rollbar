/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/rollbar"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

// Suite is the acceptance testing suite.
type Suite struct {
	suite.Suite
	provider  *schema.Provider
	providers map[string]*schema.Provider

	// The following variables are populated before each test by SetupTest():
	projectName string // Name of a Rollbar project
}

func (s *Suite) SetupSuite() {
	// Log to console
	log.Logger = log.
		With().Caller().
		Logger()
	if os.Getenv("TERRAFORM_PROVIDER_ROLLBAR_DEBUG") == "1" {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Setup testing
	s.provider = rollbar.Provider()
	s.providers = map[string]*schema.Provider{
		"rollbar": s.provider,
	}
}

// preCheck ensures we are ready to run the test
func (s *Suite) preCheck() {
	token := os.Getenv("ROLLBAR_TOKEN")
	s.NotEmpty(token, "ROLLBAR_TOKEN must be set for acceptance tests")
	log.Debug().Msg("Passed preflight check")
}

func (s *Suite) SetupTest() {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	s.projectName = fmt.Sprintf("tf-acc-test-%s", randString)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
