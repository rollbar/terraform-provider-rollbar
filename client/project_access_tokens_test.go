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

/*
// TestGetProjectAccessTokenByProjectIDAndName tests getting a project access
// token by name.
func (s *ClientTestSuite) TestGetProjectAccessTokenByProjectIDAndName() {
	projectID := 12116
	u := apiUrl + pathPATList
	u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(projectID))

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("project_access_tokens/list.json"))
	})

	examples := []struct {
		name      string
		projectID int
		expected  *ProjectAccessToken
	}{
		{
			name:      "ProjectDoesNotExist",
			projectID: projectID,
			expected:  nil,
		},
		{
			name:      "write",
			projectID: projectID,
			expected: &ProjectAccessToken{
				ProjectID:    projectID,
				AccessToken:  "access-token-12116-3",
				Name:         "write",
				Status:       "enabled",
				DateCreated:  1407933922,
				DateModified: 1407933922,
			},
		},
	}

	for _, example := range examples {
		actual, err := client.GetProjectAccessTokenByProjectIDAndName(example.projectID, example.name)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, example.expected) {
			t.Errorf("expected project %v, got %v.", example.expected, actual)
		}
	}
}


*/
