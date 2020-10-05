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
	switch resp.StatusCode() {
	case http.StatusOK:
		pats := resp.Result().(*listProjectAccessTokensResponse).Result
		return pats, nil
	case http.StatusNotFound:
		l.Warn().Msg("Project not found")
		return nil, ErrNotFound
	default:
		errResp := resp.Error().(*ErrorResult)
		l.Err(errResp).Msg("Unexpected error")
		return nil, errResp
	}
}

// ProjectAccessTokenByName returns the first project access token from the that
// matches a given name. If there is no match it returns error ErrPATNotFound.
func (c *RollbarApiClient) ProjectAccessTokenByName(projectID int, name string) (ProjectAccessToken, error) {
	l := log.With().
		Int("projectID", projectID).
		Str("name", name).
		Logger()
	l.Debug().Msg("Get project access token by name")

	var pat ProjectAccessToken
	tokens, err := c.ListProjectAccessTokens(projectID)
	if err != nil {
		l.Err(err).
			Msg("Error listing project access tokens")
		return pat, err
	}

	for _, t := range tokens {
		l.Debug().Msg("Found project access token with matching name")
		if t.Name == name {
			return t, nil
		}
	}

	l.Warn().Msg("Could not find project access token with matching name")
	return pat, ErrNotFound
}
