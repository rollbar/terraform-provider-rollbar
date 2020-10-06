package client

import (
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

// fixture loads a JSON file from the fixtures folder and returns it as a string
func fixture(path string) string {
	const fixPath = "testdata/fixtures/"
	b, err := ioutil.ReadFile(fixPath + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

/*
 * Testify setup
 */

// Suite is a Testify test suite for the Rollbar API client
type Suite struct {
	suite.Suite
	client *RollbarApiClient
}

func (s *Suite) SetupSuite() {
	// Seed gofakeit random generator
	gofakeit.Seed(0) // Setting seed to 0 will use time.Now().UnixNano()

	// Setup RollbarApiClient and enable mocking
	c, err := NewClient("fakeTokenString")
	s.Nil(err)
	httpmock.ActivateNonDefault(c.GetHttpClient())
	s.client = c
}

func (s *Suite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (s *Suite) BeforeTest() {
	httpmock.Reset()
}

// TestRollbarClientTestSuite connects the Testify test suite to the 'go test'
// built-in testing framework.
func TestRollbarClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
