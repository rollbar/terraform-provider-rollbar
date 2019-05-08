package rollbar

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiBaseURL string = "https://api.rollbar.com/api/1/"

// Option : A function to add the base url and other parameters to the client.
type Option func(*Client) error

// Client : A data structure for the rollbar client.
type Client struct {
	APIKey     string
	APIBaseURL string
}

// BaseURL : A function for construction the rul to the api.
func BaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.APIBaseURL = baseURL
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

// NewClient : A function for initiating the client.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	client := Client{
		APIKey:     apiKey,
		APIBaseURL: apiBaseURL,
	}
	if err := client.parseOptions(opts...); err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *Client) makeRequest(req *http.Request) ([]byte, error) {

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
