package client

import (
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"strings"
)

// TestListProjectAccessTokens tests listing  project access tokens.
func (s *Suite) TestListProjectAccessTokens() {
	projectID := 12116
	u := apiUrl + pathPATList
	u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(projectID))

	var lpatr listProjectAccessTokensResponse
	gofakeit.Struct(&lpatr)
	r := httpmock.NewJsonResponderOrPanic(http.StatusOK, lpatr)
	httpmock.RegisterResponder("GET", u, r)

	// Valid project ID
	actual, err := s.client.ListProjectAccessTokens(projectID)
	s.Nil(err)
	s.Equal(lpatr.Result, actual)

	// Unreachable server
	httpmock.Reset()
	_, err = s.client.ListProjectAccessTokens(projectID)
	s.NotNil(err)
	s.NotEqual(ErrNotFound, err)
}

// TestProjectAccessTokenByName tests getting a project access token by name.
func (s *Suite) TestProjectAccessTokenByName() {
	projectID := 12116
	u := apiUrl + pathPATList
	u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(projectID))

	var lpatr listProjectAccessTokensResponse
	gofakeit.Struct(&lpatr)
	r := httpmock.NewJsonResponderOrPanic(http.StatusOK, lpatr)
	httpmock.RegisterResponder("GET", u, r)

	// PAT with name exists
	actual := lpatr.Result[0]
	expected, err := s.client.ReadProjectAccessToken(projectID, actual.Name)
	s.Nil(err)
	s.Equal(expected, actual)

	// PAT with name does not exist
	_, err = s.client.ReadProjectAccessToken(projectID, "this-name-does-not-exist")
	s.Equal(ErrNotFound, err)

	// Project ID not found
	r = httpmock.NewJsonResponderOrPanic(http.StatusNotFound, ErrorResult{Err: 404, Message: "Not Found"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadProjectAccessToken(projectID, "this-name-does-not-exist")
	s.Equal(ErrNotFound, err)

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError,
		ErrorResult{Err: 500, Message: "Internal Server Error"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadProjectAccessToken(projectID, "this-name-does-not-exist")
	s.NotNil(err)
	s.NotEqual(ErrNotFound, err)
}

/*
 * Actual recorded responses from Rollbar API
 */

const patListJsonResponse = `
{
    "err": 0,
    "result": [
        {
            "access_token": "80f235b890c34ca49bcea692c2b90421",
            "cur_rate_limit_window_count": null,
            "cur_rate_limit_window_start": null,
            "date_created": 1601982124,
            "date_modified": 1601982124,
            "name": "post_client_item",
            "project_id": 411334,
            "rate_limit_window_count": null,
            "rate_limit_window_size": null,
            "scopes": [
                "post_client_item"
            ],
            "status": "enabled"
        },
        {
            "access_token": "8d4b7e0e6a1a498db82cffd1eda93376",
            "cur_rate_limit_window_count": null,
            "cur_rate_limit_window_start": null,
            "date_created": 1601982124,
            "date_modified": 1601982124,
            "name": "post_server_item",
            "project_id": 411334,
            "rate_limit_window_count": null,
            "rate_limit_window_size": null,
            "scopes": [
                "post_server_item"
            ],
            "status": "enabled"
        },
        {
            "access_token": "90b2521327a647f9aa80ef6d84427485",
            "cur_rate_limit_window_count": null,
            "cur_rate_limit_window_start": null,
            "date_created": 1601982124,
            "date_modified": 1601982124,
            "name": "read",
            "project_id": 411334,
            "rate_limit_window_count": null,
            "rate_limit_window_size": null,
            "scopes": [
                "read"
            ],
            "status": "enabled"
        },
        {
            "access_token": "d6d4b456f72048dfb8a933afe3ac66f6",
            "cur_rate_limit_window_count": null,
            "cur_rate_limit_window_start": null,
            "date_created": 1601982124,
            "date_modified": 1601982124,
            "name": "write",
            "project_id": 411334,
            "rate_limit_window_count": null,
            "rate_limit_window_size": null,
            "scopes": [
                "write"
            ],
            "status": "enabled"
        }
    ]
}
`

const patUpdateJsonResponse = `
{
    "err": 0
}
`
