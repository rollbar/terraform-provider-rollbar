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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/rollbar"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
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

func TestAccSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// getResourceIDString returns the ID string of a resource.
func (s *Suite) getResourceIDString(ts *terraform.State, resourceName string) (string, error) {
	var id string
	rs, ok := ts.RootModule().Resources[resourceName]
	if !ok {
		return id, fmt.Errorf("can't find resource: %s", resourceName)
	}

	if rs.Primary.ID == "" {
		return id, fmt.Errorf("resource ID not set")
	}
	return rs.Primary.ID, nil
}

// getResourceIDInt returns the ID of a resource as an integer.
func (s *Suite) getResourceIDInt(ts *terraform.State, resourceName string) (int, error) {
	var id int
	idString, err := s.getResourceIDString(ts, resourceName)
	if err != nil {
		return id, err
	}
	id, err = strconv.Atoi(idString)
	if err != nil {
		return id, err
	}
	return id, nil
}

// checkResourceStateSanity checks that the resource is present in the Terraform
// state, and that its ID is set.
func (s *Suite) checkResourceStateSanity(rn string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		_, err := s.getResourceIDString(ts, rn)
		return err
	}
}
