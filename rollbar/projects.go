package rollbar

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Project represents a project
type Project struct {
	AccountID    int    `json:"account_id"`
	ID           int    `json:"id"`
	DateCreated  int    `json:"date_created"`
	DateModified int    `json:"date_modified"`
	Name         string `json:"name"`
}

// listProjectsResponse represents the list projects response
type listProjectsResponse struct {
	Error  int `json:"err"`
	Result []*Project
}

// ListProjects lists the projects for this API Key
// https://docs.rollbar.com/reference#list-all-projects
func (c *Client) ListProjects() ([]*Project, error) {
	var data listProjectsResponse

	url := fmt.Sprintf("%sprojects?access_token=%s", c.APIBaseURL, c.APIKey)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	bytes, err := c.makeRequest(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return nil, err
	}

	projects := []*Project{}
	for _, p := range data.Result {
		projects = append(projects, p)
	}

	return projects, nil
}
