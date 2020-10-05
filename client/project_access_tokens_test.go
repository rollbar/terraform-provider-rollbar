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

	actual, err := s.client.ListProjectAccessTokens(projectID)
	s.Nil(err)
	s.Equal(lpatr.Result, actual)
}

// TestGetProjectAccessTokenByProjectIDAndName tests getting a project access
// token by name.
func (s *ClientTestSuite) TestGetProjectAccessTokenByProjectIDAndName() {
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
	s.Equal(err, ErrNotFound)
}
