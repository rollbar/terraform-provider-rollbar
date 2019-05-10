package rollbar

import (
	"bytes"
	"fmt"
	"io"
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

// BaseURL returns a function for configuring the url inside a Client.
func BaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.APIBaseURL = baseURL
		var apiURL = &url.URL{}
		apiURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.APIBasePath = apiURL.Path
		if apiURL.Path == "/" {
			c.APIBasePath = ""
		}

		c.APIHost = apiURL.Host
		c.APIScheme = apiURL.Scheme

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

func (c *Client) get(pathComponents ...string) ([]byte, error) {
	return c.getWithQueryParams(map[string]string{}, pathComponents...)
}

func (c *Client) getWithQueryParams(queryParams map[string]string, pathComponents ...string) ([]byte, error) {
	url := c.url(true, queryParams, pathComponents...)

	bytes, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (c *Client) post(data []byte, pathComponents ...string) ([]byte, error) {
	url := c.url(false, map[string]string{}, pathComponents...)
	body := bytes.NewBuffer(data)

	bytes, err := c.makeRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (c *Client) delete(pathComponents ...string) error {
	url := c.url(true, map[string]string{}, pathComponents...)

	_, err := c.makeRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	return nil
}

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

func (c *Client) makeRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("The http status code is not '200' and is '%d'", resp.StatusCode)
	}

	return responseBody, nil
}
