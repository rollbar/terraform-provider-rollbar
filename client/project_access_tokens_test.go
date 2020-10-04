package client

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project Access Tokens", func() {
	u := apiUrl + pathPATList

	When("There are no tokens attached to the project", func() {
		It("lists zero tokens", func() {
			s := `{ "err": 0, "result": [] }`
			stringResponse := httpmock.NewStringResponse(200, s)
			stringResponse.Header.Add("Content-Type", "application/json")
			responder := httpmock.ResponderFromResponse(stringResponse)
			httpmock.RegisterResponder("GET", u, responder)

			pats, err := c.ListProjectAccessTokens(0) // Project ID doesn't matter

			Expect(err).ToNot(HaveOccurred())
			Expect(pats).To(HaveLen(0))
		})
	})

	When("There are tokens attached to the project", func() {
		s := fixture("project_access_tokens/list.json")
		stringResponse := httpmock.NewStringResponse(200, s)
		stringResponse.Header.Add("Content-Type", "application/json")
		responder := httpmock.ResponderFromResponse(stringResponse)
		httpmock.RegisterResponder("GET", u, responder)

	})
})

/*
func TestListProjectAccessTokens(t *testing.T) {
	teardown := setup()
	defer teardown()

	projectID := 12116
	handURL := fmt.Sprintf("/project/%d/access_tokens/", projectID)

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("project_access_tokens/list.json"))
	})

	actual, err := client.ListProjectAccessTokens(projectID)

	if err != nil {
		t.Fatal(err)
	}

	expected := []*ProjectAccessToken{
		{
			ProjectID:    projectID,
			AccessToken:  "access-token-12116-1",
			Name:         "post_client_item",
			Status:       "enabled",
			DateCreated:  1407933922,
			DateModified: 1407933922,
		},
		{
			ProjectID:    projectID,
			AccessToken:  "access-token-12116-2",
			Name:         "post_server_item",
			Status:       "enabled",
			DateCreated:  1407933922,
			DateModified: 1439579817,
		},
		{
			ProjectID:    projectID,
			AccessToken:  "access-token-12116-3",
			Name:         "write",
			Status:       "enabled",
			DateCreated:  1407933922,
			DateModified: 1407933922,
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected response %v, got %v.", expected, actual)
	}
}

func TestGetProjectAccessTokenByProjectIDAndName(t *testing.T) {
	teardown := setup()
	defer teardown()

	projectID := 12116
	handURL := fmt.Sprintf("/project/%d/access_tokens/", projectID)

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
