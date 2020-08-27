package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetCoffees - Returns list of coffees (no auth required)
func (c *RollbarApiClient) GetProjects() ([]Project, error) {
	u := c.url + "foo"
	u := url.Parse("foobar")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/coffees", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	coffees := []Coffee{}
	err = json.Unmarshal(body, &coffees)
	if err != nil {
		return nil, err
	}

	return coffees, nil
}
