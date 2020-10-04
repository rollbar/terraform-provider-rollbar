package client

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
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
	Result []ProjectAccessToken
}

// ListProjectAccessTokens lists the projects access tokens for a given project id
// https://docs.rollbar.com/reference#list-all-project-access-tokens
func (c *RollbarApiClient) ListProjectAccessTokens(projectID int) ([]ProjectAccessToken, error) {
	//return []ProjectAccessToken{}, nil

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

// GetProjectAccessTokenByProjectIDAndName returns the first project access token from the
// list-projects-tokens call that matches a given name nad project id. If there is no
// matching project, return nil.
func (c *RollbarApiClient) GetProjectAccessTokenByProjectIDAndName(projectID int, name string) (ProjectAccessToken, error) {
	return ProjectAccessToken{}, nil
	/*
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

	*/
}
