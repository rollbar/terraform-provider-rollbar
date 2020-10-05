package client

import (
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"strings"
)

// TestListProjectAccessTokens tests listing  project access tokens.
func (s *ClientTestSuite) TestListProjectAccessTokens() {
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
func (s *ClientTestSuite) TestProjectAccessTokenByName() {
	projectID := 12116
	u := apiUrl + pathPATList
	u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(projectID))

	var lpatr listProjectAccessTokensResponse
	gofakeit.Struct(&lpatr)
	r := httpmock.NewJsonResponderOrPanic(http.StatusOK, lpatr)
	httpmock.RegisterResponder("GET", u, r)

	// PAT with name exists
	actual := lpatr.Result[0]
	expected, err := s.client.ProjectAccessTokenByName(projectID, actual.Name)
	s.Nil(err)
	s.Equal(expected, actual)

	// PAT with name does not exist
	_, err = s.client.ProjectAccessTokenByName(projectID, "this-name-does-not-exist")
	s.Equal(ErrNotFound, err)

	// Project ID not found
	r = httpmock.NewJsonResponderOrPanic(http.StatusNotFound, ErrorResult{Err: 404, Message: "Not Found"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ProjectAccessTokenByName(projectID, "this-name-does-not-exist")
	s.Equal(ErrNotFound, err)

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError,
		ErrorResult{Err: 500, Message: "Internal Server Error"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ProjectAccessTokenByName(projectID, "this-name-does-not-exist")
	s.NotNil(err)
	s.NotEqual(ErrNotFound, err)
}
