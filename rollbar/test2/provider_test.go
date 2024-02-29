/*
 * Copyright (c) 2024 Rollbar, Inc.
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

package test2

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rollbar/terraform-provider-rollbar/rollbar"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

const schemaKeyToken = "api_key"
const projectKeyToken = "project_api_key"
const schemaKeyBaseURL = "api_url"

func init() {
	// Setup nice logging
	log.Logger = log.
		With().Caller().
		Logger()
	if os.Getenv("TERRAFORM_PROVIDER_ROLLBAR_DEBUG") == "1" {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

// TestMain connects Terraform sweeper system with Go's testing framework.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// AccSuite is the acceptance testing suite.
type AccSuite struct {
	suite.Suite
	provider  *schema.Provider
	providers map[string]*schema.Provider

	// The following variables are populated before each test by SetupTest():
	randName string // Name of a Rollbar project
	randID   int    // ID of a Rollbar resource
}

// providerConfigure sets up authentication in a Resty HTTP client.
func providerConfigureTest(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	token := d.Get(schemaKeyToken).(string)
	projectToken := d.Get(projectKeyToken).(string)
	baseURL := d.Get(schemaKeyBaseURL).(string)
	c := client.NewTestClient(baseURL, token)
	pc := client.NewTestClient(baseURL, projectToken)
	return map[string]*client.RollbarAPIClient{schemaKeyToken: c, projectKeyToken: pc}, diags
}
func (s *AccSuite) SetupSuite() {
	maxprocs := runtime.GOMAXPROCS(0)
	log.Debug().
		Int("GOMAXPROCS", maxprocs).
		Send()

	// Setup testing
	s.provider = rollbar.Provider()
	s.provider.ConfigureContextFunc = providerConfigureTest
	s.providers = map[string]*schema.Provider{
		"rollbar": s.provider,
	}
}

// preCheck ensures we are ready to run the test
func (s *AccSuite) preCheck() {
	token := os.Getenv("ROLLBAR_API_KEY")
	s.NotEmpty(token, "ROLLBAR_API_KEY must be set for acceptance tests")
	log.Debug().Msg("Passed preflight check")
}

func (s *AccSuite) SetupTest() {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	s.randName = fmt.Sprintf("tf-acc-test-%s", randString)
	s.randID = acctest.RandInt()
}

func TestAccSuite(t *testing.T) {
	suite.Run(t, new(AccSuite))
}

// getResourceIDString returns the ID string of a resource.
func (s *AccSuite) getResourceIDString(ts *terraform.State, resourceName string) (string, error) {
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
func (s *AccSuite) getResourceIDInt(ts *terraform.State, resourceName string) (int, error) {
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
func (s *AccSuite) checkResourceStateSanity(rn string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		_, err := s.getResourceIDString(ts, rn)
		return err
	}
}

// getResourceAttrString returns the string value of a named attribute of a
// Terraform state resource.
func (s *AccSuite) getResourceAttrString(ts *terraform.State, resourceName string, attribute string) (string, error) {
	rootModule := ts.RootModule()
	var rs *terraform.ResourceState
	var ok bool
	if rootModule != nil {
		rs, ok = rootModule.Resources[resourceName]
	}
	if !ok {
		err := fmt.Errorf("can't find resource: %s", resourceName)
		log.Err(err).Send()
		return "", err
	}
	value, ok := rs.Primary.Attributes[attribute]
	if !ok {
		err := fmt.Errorf("can't find attribute: %s", attribute)
		log.Err(err).Send()
		return "", err
	}
	return value, nil
}

// getResourceAttrInt returns the integer value of a named attribute of a
// Terraform state resource.
func (s *AccSuite) getResourceAttrInt(ts *terraform.State, resourceName string, attribute string) (int, error) {
	value, err := s.getResourceAttrString(ts, resourceName, attribute)
	if err != nil {
		return 0, err
	}
	i, err1 := strconv.Atoi(value)
	if err1 != nil {
		log.Err(err1).Send()
		return 0, err1
	}
	return i, nil
}

// client returns the current Rollbar API client
func (s *AccSuite) client() *client.RollbarAPIClient {
	return s.provider.Meta().(map[string]*client.RollbarAPIClient)[schemaKeyToken]
}

// getResourceAttrIntSlice returns value of a named attribute of a Terraform
// state resource as a slice of integers.
func (s *AccSuite) getResourceAttrIntSlice(ts *terraform.State, resourceName string, attribute string) ([]int, error) {
	var value []int
	count, err := s.getResourceAttrInt(ts, resourceName, attribute+".#")
	if err != nil {
		return nil, err
	}
	for i := 0; i < count; i++ {
		elementAttr := fmt.Sprintf("%s.%d", attribute, i)
		element, err := s.getResourceAttrInt(ts, resourceName, elementAttr)
		if err != nil {
			return nil, err
		}
		value = append(value, element)
	}
	return value, nil
}
