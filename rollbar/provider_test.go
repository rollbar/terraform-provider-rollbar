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

package rollbar

import (
	"fmt"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

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

	// VCR Mode
	vcr      bool // Is VCR mode enabled?
	recorder *recorder.Recorder
}

func (s *AccSuite) SetupSuite() {
	// Setup testing
	s.provider = Provider()
	s.providers = map[string]*schema.Provider{
		"rollbar": s.provider,
	}
	useVCR := os.Getenv("TF_ACC_VCR")
	if useVCR == "1" {
		s.vcr = true
		log.Warn().Msg("VCR mode enabled")
		runtime.GOMAXPROCS(1)
	}

	// Log GOMAXPROCS
	maxprocs := runtime.GOMAXPROCS(0)
	log.Debug().
		Int("GOMAXPROCS", maxprocs).
		Send()

}

// preCheck ensures we are ready to run the test
func (s *AccSuite) preCheck() {
	token := os.Getenv("ROLLBAR_API_KEY")
	s.NotEmpty(token, "ROLLBAR_API_KEY must be set for acceptance tests")
	log.Debug().Msg("Passed preflight check")
}

func (s *AccSuite) SetupTest() {
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	if s.vcr {
		// For VCR mode we need deterministic names
		randString = "vcr"
		cassetteName := fmt.Sprintf("vcr/%s", s.T().Name())
		log.Warn().
			Str("cassette_name", cassetteName).
			Msg("using VCR cassette")
		r, err := recorder.New(cassetteName)
		s.Nil(err)
		s.recorder = r
		r.AddFilter(vcrFilterHeaders)
		s.provider = ProviderWithTransport(s.recorder)
		s.providers = map[string]*schema.Provider{
			"rollbar": s.provider,
		}
		switch r.Mode() {
		case recorder.ModeRecording:
			log.Debug().Msg("VCR recording mode")
		case recorder.ModeReplaying:
			log.Debug().Msg("VCR replay mode")
		case recorder.ModeDisabled:
			log.Debug().Msg("VCR disabled")
		}
	}
	s.randName = fmt.Sprintf("tf-acc-test-%s", randString)
}

func (s *AccSuite) TearDownTest() {
	if s.vcr {
		s.recorder.Stop()
	}
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
	rs, ok := ts.RootModule().Resources[resourceName]
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
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Err(err).Send()
		return 0, err
	}
	return i, nil
}

// client returns the current Rollbar API client
func (s *AccSuite) client() *client.RollbarAPIClient {
	return s.provider.Meta().(*client.RollbarAPIClient)
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

// parallelTestVCR runs a resource.TestCase using s.parallelTestVCR( unless
// VCR mode is enabled, whereupon it runs the test with resource.Test.
func (s *AccSuite) parallelTestVCR(t *testing.T, c resource.TestCase) {
	if s.vcr {
		switch s.recorder.Mode() {
		case recorder.ModeRecording:
			log.Warn().Msg("VCR Mode Recording")
			resource.Test(t, c)
		case recorder.ModeReplaying:
			log.Warn().Msg("VCR Mode Replaying")
			resource.Test(t, c)
		case recorder.ModeDisabled: // Will this ever actually happen?
			log.Warn().Msg("VCR Mode Disabled")
			resource.ParallelTest(t, c)
		default:
			s.Fail("WTF?")
		}
	} else {
		log.Warn().Msg("No VCR Mode")
		resource.ParallelTest(t, c)
	}
}

// vcrFilterHeaders removes unnecessary headers from VCR recordings.
func vcrFilterHeaders(i *cassette.Interaction) error {
	delete(i.Request.Headers, "X-Rollbar-Access-Token")
	delete(i.Request.Headers, "User-Agent")
	for key := range i.Response.Headers {
		deleteHeader := false
		if strings.HasPrefix(key, "X-") {
			deleteHeader = true
		}
		if strings.HasPrefix(key, "Access-Control-") {
			deleteHeader = true
		}
		switch key {
		case "Alt-Svc", "Content-Length", "Etag", "Server", "Via":
			deleteHeader = true
		}
		if deleteHeader {
			delete(i.Response.Headers, key)
		}
	}
	return nil
}
