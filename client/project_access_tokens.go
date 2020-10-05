package client

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// ProjectAccessToken represents a Rollbar project access token
type ProjectAccessToken struct {
	ProjectID    int    `json:"project_id" fake:"{number1,1000000}"`
	AccessToken  string `json:"access_token"`
	Name         string `json:"name" fake:"{hackernoun}"`
	Status       string `json:"status" fake:"{randomwords:[enabled,disabled]}"`
	DateCreated  int    `json:"date_created"`
	DateModified int    `json:"date_modified"`
}

// listProjectAccessTokensResponse represents the list-project-access-tokens response
type listProjectAccessTokensResponse struct {
	Error  int `json:"err"`
	Result []ProjectAccessToken
}

// ListProjectAccessTokens lists the project's access tokens.
func (c *RollbarApiClient) ListProjectAccessTokens(projectID int) ([]ProjectAccessToken, error) {
	u := apiUrl + pathPATList

	l := log.With().
		Str("url", u).
		Logger()

	resp, err := c.resty.R().
		SetResult(listProjectAccessTokensResponse{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"projectId": strconv.Itoa(projectID),
		}).
		Get(u)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		errResp := resp.Error().(*ErrorResult)
		l.Err(errResp).Send()
		return nil, errResp
	}
	pats := resp.Result().(*listProjectAccessTokensResponse).Result

	return pats, nil
}

// ProjectAccessTokenByName returns the first project access token from the that
// matches a given name. If there is no match it returns error ErrPATNotFound.
func (c *RollbarApiClient) ProjectAccessTokenByName(projectID int, name string) (ProjectAccessToken, error) {
	var pat ProjectAccessToken
	tokens, err := c.ListProjectAccessTokens(projectID)
	if err != nil {
		return pat, err
	}

	for _, t := range tokens {
		if t.Name == name {
			return t, nil
		}
	}

	return pat, ErrNotFound
}
