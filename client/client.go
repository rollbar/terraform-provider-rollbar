package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/url"
)

// RollbarApiClient is a client for the Rollbar API
type RollbarApiClient struct {
	r   *resty.Client
	url *url.URL
}

// NewClient sets up a new Rollbar API client
func NewClient(apiUrl string, token string) (*RollbarApiClient, error) {
	log.Debug().
		Str("rollbar_api_url", apiUrl).
		Str("rollbar_token", token).
		Msg("Initializing Rollbar client")

	// URL
	u, err := url.Parse(apiUrl)
	if err != nil {
		log.Err(err).Msg("Could not parse Rollbar API URL")
		return nil, err
	}

	// Authentication
	r := resty.New() // HTTP c
	if token != "" {
		r = r.SetHeader("X-Rollbar-Access-Token", token)
	}
	c := RollbarApiClient{r: r, url: u}

	return &c, nil
}
