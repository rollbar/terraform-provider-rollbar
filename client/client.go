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
)

const apiUrl = "https://api.rollbar.com"

// RollbarApiClient is a client for the Rollbar API
type RollbarApiClient struct {
	resty *resty.Client
}

// NewClient sets up a new Rollbar API client
func NewClient(token string) (*RollbarApiClient, error) {
	log.Debug().Msg("Initializing Rollbar client")

	// Authentication
	r := resty.New() // HTTP c
	if token != "" {
		r = r.SetHeader("X-Rollbar-Access-Token", token)
	} else {
		log.Warn().Msg("Rollbar API token not set")
	}
	c := RollbarApiClient{resty: r}

	return &c, nil
}
