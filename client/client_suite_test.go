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

package client

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

func loadFixture(fixturePath string) string {
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
	return s
}

// responseFromFixture creates an http.Response based on a fixture file loaded
// from folder 'client/fixtures/'.
func responseFromFixture(fixturePath string, status int) *http.Response {
	s := loadFixture(fixturePath)
	rs := httpmock.NewStringResponse(status, s)
	rs.Header.Add("Content-Type", "application/json")
	return rs
}

// responderFromFixture creates an httpmock.Responder based on a fixture file
// loaded from folder 'client/fixtures/'.
func responderFromFixture(fixturePath string, status int) httpmock.Responder {
	rs := responseFromFixture(fixturePath, status)
	r := httpmock.ResponderFromResponse(rs)
	return r
}

/*
 * Testify setup
 */

// Suite is a Testify test suite for the Rollbar API client
type Suite struct {
	suite.Suite
	client *RollbarAPIClient
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

	// Setup RollbarAPIClient and enable mocking
	c := NewTestClient(DefaultBaseURL, "fakeTokenString")

	httpmock.ActivateNonDefault(c.Resty.GetClient())
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

// checkServerErrors check correct handling of various API error responses
func (s *Suite) checkServerErrors(mockMethod, mockUrl string, testFunc func() error) {
	// Not Found
	r := httpmock.NewJsonResponderOrPanic(http.StatusNotFound,
		ErrorResult{Err: 404, Message: "Not Found"})
	httpmock.RegisterResponder(mockMethod, mockUrl, r)
	err := testFunc()
	s.Equal(ErrNotFound, err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder(mockMethod, mockUrl, r)
	err = testFunc()
	s.Equal(ErrUnauthorized, err)

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError,
		ErrorResult{Err: 500, Message: "Internal Server Error"})
	httpmock.RegisterResponder(mockMethod, mockUrl, r)
	err = testFunc()
	s.NotNil(err)
	s.NotEqual(ErrNotFound, err)

	// Unreachable server
	httpmock.Reset()
	err = testFunc()
	s.NotNil(err)
}

// checkServerErrorsWithQuery check correct handling of various API error responses (with query)
func (s *Suite) checkServerErrorsWithQuery(mockMethod, mockUrl string, expectedQuery map[string]string, testFunc func() error) {
	// Not Found
	r := httpmock.NewJsonResponderOrPanic(http.StatusNotFound,
		ErrorResult{Err: 404, Message: "Not Found"})
	httpmock.RegisterResponderWithQuery(mockMethod, mockUrl, expectedQuery, r)
	err := testFunc()
	s.Equal(ErrNotFound, err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponderWithQuery(mockMethod, mockUrl, expectedQuery, r)
	err = testFunc()
	s.Equal(ErrUnauthorized, err)

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError,
		ErrorResult{Err: 500, Message: "Internal Server Error"})
	httpmock.RegisterResponderWithQuery(mockMethod, mockUrl, expectedQuery, r)
	err = testFunc()
	s.NotNil(err)
	s.NotEqual(ErrNotFound, err)

	// Unreachable server
	httpmock.Reset()
	err = testFunc()
	s.NotNil(err)
}
