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

// Package client is a client library for accessing the Rollbar API.
package client

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// DefaultBaseURL is the default base URL for the Rollbar API.
const DefaultBaseURL = "https://api.rollbar.com"
const Version = "v1.15.0"

// RollbarAPIClient is a client for the Rollbar API.
type RollbarAPIClient struct {
	BaseURL string // Base URL for Rollbar API
	Resty   *resty.Client

	m sync.Mutex
}

// NewTestClient sets up a new Rollbar API test client.
func NewTestClient(baseURL, token string) *RollbarAPIClient {
	log.Debug().Msg("Initializing Rollbar client")

	// New Resty HTTP client
	r := resty.New()

	// Use default transport - needed for VCR
	r.SetTransport(http.DefaultTransport).
		// set timeout on http client
		SetTimeout(30 * time.Second).
		// Set retry count to 1 (try 2 times before it fails)
		SetRetryCount(1).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	// Authentication
	if token != "" {
		r = r.SetHeaders(map[string]string{
			"X-Rollbar-Access-Token": token,
			"X-Rollbar-Terraform":    "true"})
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
	c := RollbarAPIClient{
		Resty:   r,
		BaseURL: baseURL,
	}
	return &c
}

func (c *RollbarAPIClient) SetHeaderResource(header string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.Resty.SetHeader("X-Rollbar-Terraform-Resource", header)
}

func (c *RollbarAPIClient) SetHeaderDataSource(header string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.Resty.SetHeader("X-Rollbar-Terraform-DataSource", header)
}

// NewClient sets up a new Rollbar API client.
func NewClient(baseURL, token string) *RollbarAPIClient {
	log.Debug().Msg("Initializing Rollbar client")
	now := time.Now().Format(time.RFC3339Nano)
	// New Resty HTTP client
	r := resty.New()

	// Use default transport - needed for VCR
	r.SetTransport(http.DefaultTransport).
		// set timeout on http client
		SetTimeout(30 * time.Second).
		// Set retry count to 4 (try 5 times before it fails)
		SetRetryCount(4).
		SetRetryWaitTime(8 * time.Second).
		SetRetryMaxWaitTime(50 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				if err != nil { // network error
					return true
				}
				return r.StatusCode() == http.StatusNotFound ||
					r.StatusCode() == http.StatusTooManyRequests ||
					r.StatusCode() == http.StatusInternalServerError ||
					r.StatusCode() == http.StatusBadGateway
			})
	// Authentication

	if token != "" {
		r = r.SetHeaders(map[string]string{
			"X-Rollbar-Access-Token":      token,
			"X-Rollbar-Terraform-Version": Version,
			"X-Rollbar-Terraform-Date":    now,
		})
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
	c := RollbarAPIClient{
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
