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

/*
 * Package client is a client library for accessing the Rollbar API.
 */
package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

// DefaultBaseURL is the default base URL for the Rollbar API.
const DefaultBaseURL = "https://api.rollbar.com"

// RollbarApiClient is a client for the Rollbar API.
type RollbarApiClient struct {
	BaseURL string // Base URL for Rollbar API
	Resty   *resty.Client
}

// NewClient sets up a new Rollbar API client.
func NewClient(baseURL, token string) *RollbarApiClient {
	log.Debug().Msg("Initializing Rollbar client")

	// New Resty HTTP client
	r := resty.New()

	// Use default transport - needed for VCR
	r.SetTransport(http.DefaultTransport)

	// Authentication
	if token != "" {
		r = r.SetHeader("X-Rollbar-Access-Token", token)
	} else {
		log.Warn().Msg("Rollbar API token not set")
	}

	// Authentication
	if baseURL == "" {
		log.Error().Msg("Rollbar API base URL not set")
	}

	// Configure Resty to use Zerolog for logging
	r.SetLogger(restyZeroLogger{log.Logger})

	// Rollbar client
	c := RollbarApiClient{
		Resty:   r,
		BaseURL: baseURL,
	}
	return &c
}

// Status represents the enabled or disabled status of an entity.
type Status string

// Possible values for status
const (
	StatusEnabled  = Status("enabled")
	StatusDisabled = Status("disabled")
)

// errorFromResponse interprets the status code of Resty response, returning nil
// on success or an appropriate error code
func errorFromResponse(resp *resty.Response) error {
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusCreated:
		return nil
	case http.StatusForbidden:
		// Workaround for known API bug:
		// https://github.com/rollbar/terraform-provider-rollbar/issues/79
		er := resp.Error().(*ErrorResult)
		if strings.Contains(er.Message, "not found in this account") {
			log.Debug().
				Str("status", resp.Status()).
				Int("status_code", resp.StatusCode()).
				Str("error_message", er.Message).
				Int("error_code", er.Err).
				Msg("Workaround entity not found API bug")
			return ErrNotFound
		}
		fallthrough
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrNotFound
	default:
		er := resp.Error().(*ErrorResult)
		log.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Send()
		return er
	}
}
