package rollbar

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const apiBaseURL string = "https://api.rollbar.com/api/1/"
const apiScheme string = "https"
const apiHost string = "api.rollbar.com"
const apiBasePath string = "/api/1"

// Option : A function to add the base url and other parameters to the client.
type Option func(*Client) error

// Client : A data structure for the rollbar client.
type Client struct {
	APIKey     string
	APIBaseURL string
	APIKey      string
	APIScheme   string
	APIHost     string
	APIBasePath string
}

// BaseURL : A function for construction the url to the api.
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
		APIKey:      apiKey,
		APIScheme:   apiScheme,
		APIHost:     apiHost,
		APIBasePath: apiBasePath,
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
func (c *Client) url(withAccessToken bool, queryMap map[string]string, pathComponents ...string) string {
	query := url.Values{}
	for key, value := range queryMap {
		query.Add(key, value)
	}
	if withAccessToken {
		query.Add("access_token", c.APIKey)
	}

	components := append([]string{c.APIBasePath}, pathComponents...)
	path := strings.Join(components, "/")

	u := url.URL{
		Scheme:   c.APIScheme,
		Host:     c.APIHost,
		Path:     path,
		RawQuery: query.Encode(),
	}

	return u.String()
}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("The http status code is not '200' and is '%d'", resp.StatusCode)
	}

	return body, nil
}
