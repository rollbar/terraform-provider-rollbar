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

func NewClient(apiKey string, opts ...Option) (*Client, error) {
	client := &Client{
		// ApiKey:     apiKey,
		ApiBaseUrl: apiBaseUrl,
	}
	if err := client.parseOptions(opts...); err != nil {
		return nil, err
	}
	return client, nil

}

func (s *Client) makeRequest(req *http.Request) ([]byte, error) {

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("The http status code is not '200' and is '%d'", resp.StatusCode)
	}

	return body, nil
}
