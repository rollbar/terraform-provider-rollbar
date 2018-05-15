package rollbar

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiBaseUrl string = "https://api.rollbar.com/api/1/"

var apiKey string

type Option func(*Client) error

type Client struct {
	ApiKey     string
	ApiBaseUrl string
}

func BaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.ApiBaseUrl = baseURL
		return nil
	}
}

func ApiKey(apiKey string) Option {
	return func(c *Client) error {
		c.ApiKey = apiKey
		return nil
	}
}

func (c *Client) parseOptions(opts ...Option) error {
	// Range over each options function and apply it to our API type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		ApiKey:     apiKey,
		ApiBaseUrl: apiBaseUrl,
	}
	if err := client.parseOptions(opts...); err != nil {
		return nil, err
	}
	return client, nil

}

func (s *Client) makeRequest(req *http.Request) ([]byte, error) {

	client := &http.Client{}
	resp, client_err := client.Do(req)

	if client_err != nil {
		return nil, client_err
	}

	defer resp.Body.Close()

	body, read_body_err := ioutil.ReadAll(resp.Body)

	if read_body_err != nil {
		return nil, read_body_err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}
