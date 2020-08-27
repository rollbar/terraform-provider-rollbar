package client

import "github.com/go-resty/resty/v2"

// RollbarApiClient is a client for the Rollbar API
type RollbarApiClient struct {
	resty *resty.Client
}

// NewClient sets up a new Rollbar API client
func NewClient(token string) *RollbarApiClient {
	c := resty.New() // HTTP client
	if token != "" {
		c = c.SetHeader("X-Rollbar-Access-Token", token)
	}
	return &RollbarApiClient{resty: c}
}
