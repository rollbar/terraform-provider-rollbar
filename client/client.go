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

// RollbarApiClient is a client for the Rollbar API
type RollbarApiClient struct {
	resty *resty.Client
	url   string // API base URL
}

// NewClient sets up a new Rollbar API client
func NewClient(apiUrl string, token string) (*RollbarApiClient, error) {
	log.Debug().
		Str("rollbar_api_url", apiUrl).
		Str("rollbar_token", token).
		Msg("Initializing Rollbar client")

	//// URL
	//u, err := url.Parse(apiUrl)
	//if err != nil {
	//	log.Err(err).Msg("Could not parse Rollbar API URL")
	//	return nil, err
	//}

	// Authentication
	r := resty.New() // HTTP c
	if token != "" {
		r = r.SetHeader("X-Rollbar-Access-Token", token)
	}
	c := RollbarApiClient{resty: r, url: apiUrl}

	return &c, nil
}
