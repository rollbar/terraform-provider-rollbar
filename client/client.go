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
)

const apiUrl = "https://api.rollbar.com"

// RollbarApiClient is a client for the Rollbar API.
type RollbarApiClient struct {
	Resty *resty.Client
}

// NewClient sets up a new Rollbar API client.
func NewClient(token string) *RollbarApiClient {
	log.Debug().Msg("Initializing Rollbar client")

	// New Resty HTTP client
	r := resty.New()

	// Authentication
	if token != "" {
		r = r.SetHeader("X-Rollbar-Access-Token", token)
	} else {
		log.Warn().Msg("Rollbar API token not set")
	}

	// Configure Resty to use Zerolog for logging
	r.SetLogger(restyZeroLogger{log.Logger})

	// Rollbar client
	c := RollbarApiClient{Resty: r}
	return &c
}

// Status represents the enabled or disabled status of an entity.
type Status string

// Possible values for status
const (
	StatusEnabled  = Status("enabled")
	StatusDisabled = Status("disabled")
)
