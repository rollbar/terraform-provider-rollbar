package rollbar

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ProjectAccessToken represents a project access token
type ProjectAccessToken struct {
	ProjectID    int    `json:"project_id"`
	AccessToken  string `json:"access_token"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	DateCreated  int    `json:"date_created"`
	DateModified int    `json:"date_modified"`
}

// listProjectAccessTokensResponse represents the list-project-access-tokens response
type listProjectAccessTokensResponse struct {
	Error  int `json:"err"`
	Result []*ProjectAccessToken
}

// ListProjectAccessTokens lists the projects access tokens for a given project id
// https://docs.rollbar.com/reference#list-all-project-access-tokens
func (c *Client) ListProjectAccessTokens(projectID int) ([]*ProjectAccessToken, error) {
	var data listProjectAccessTokensResponse

	url := fmt.Sprintf("%sproject/%d/access_tokens/?access_token=%s", c.APIBaseURL, projectID, c.APIKey)

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

	tokens := []*ProjectAccessToken{}
	for _, t := range data.Result {
		tokens = append(tokens, t)
	}

	return tokens, nil
}

// GetProjectAccessTokenByProjectIDAndName returns the first project access token from the
// list-projects-tokens call that matches a given name nad project id. If there is no
// matching project, return nil.
func (c *Client) GetProjectAccessTokenByProjectIDAndName(projectID int, name string) (*ProjectAccessToken, error) {
	tokens, err := c.ListProjectAccessTokens(projectID)
	if err != nil {
		return nil, err
	}

	for _, t := range tokens {
		if t.Name == name {
			return t, nil
		}
	}

	return nil, nil
}
