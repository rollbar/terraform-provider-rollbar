package client

import (
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

// fixtureResponder creates an httpmock.Responder based on a fixture file
// loaded from folder 'client/fixtures/'.
func fixtureResponder(fixturePath string, status int) httpmock.Responder {
	const fixtureFolder = "fixtures/"
	b, err := ioutil.ReadFile(fixtureFolder + fixturePath) // #nosec
	if err != nil {
		log.Fatal().
			Err(err).
			Str("fixtureFolder", fixtureFolder).
			Str("fixturePath", fixturePath).
			Msg("Error loading fixture")
	}
	s := string(b)
	rs := httpmock.NewStringResponse(status, s)
	rs.Header.Add("Content-Type", "application/json")
	r := httpmock.ResponderFromResponse(rs)
	return r
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
	// Pretty logging
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().
		Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

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
