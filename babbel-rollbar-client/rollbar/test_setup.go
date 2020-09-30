package rollbar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
)

const (
	fixPath = "testdata/fixtures/"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

// VarsUsers : A data structure for the mock response.
type VarsUsers struct {
	TeamID    int    `json:"TeamID"`
	UserID    int    `json:"UserID"`
	UserEmail string `json:"UserEmail"`
}

func setup() func() {
	const mockAPIKey = "mockApiKey"

	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient(mockAPIKey, BaseURL(server.URL+"/"))

	return func() {
		server.Close()
	}
}

// BaseURL returns a function for configuring the url inside a Client.
func BaseURL(baseURL string) Option {
	return func(c *Client) error {
		var apiURL = &url.URL{}
		apiURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.BasePath = apiURL.Path
		if apiURL.Path == "/" {
			c.BasePath = ""
		}

		c.Host = apiURL.Host
		c.Scheme = apiURL.Scheme

		return nil
	}
}

func fixture(path string) string {
	b, err := ioutil.ReadFile(fixPath + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Another function is used for reading the fixture files
// because the other returns a string
// and a case should be formed and 2 values should be returned
// which makes it a bit complex.
func vars(fName string) (*VarsUsers, error) {
	var data VarsUsers

	fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/variables/%s", fixPath, fName))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileBytes, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
