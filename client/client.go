/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

/*
 * Package client is a client library for accessing the Rollbar API.
 */
package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/http"
)

const apiUrl = "https://api.rollbar.com"

// RollbarApiClient is a client for the Rollbar API
type RollbarApiClient struct {
	resty *resty.Client
}

// NewClient sets up a new Rollbar API c
func NewClient(token string) (*RollbarApiClient, error) {
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
	c := RollbarApiClient{resty: r}
	return &c, nil
}

// GetHttpClient allows access to the underlying http.Client object, for use
// with mocking tools.
func (c *RollbarApiClient) GetHttpClient() *http.Client {
	return c.resty.GetClient()
}
