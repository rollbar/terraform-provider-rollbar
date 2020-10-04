package client

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

var _ = Describe("Project Access Tokens", func() {

	Context("getting PATs by project ID", func() {
		u := apiUrl + pathPATList
		projectID := 12116
		u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(projectID))

		When("there are no tokens attached to the project", func() {
			It("lists zero tokens", func() {
				s := `{ "err": 0, "result": [] }`
				stringResponse := httpmock.NewStringResponse(200, s)
				stringResponse.Header.Add("Content-Type", "application/json")
				responder := httpmock.ResponderFromResponse(stringResponse)
				httpmock.RegisterResponder("GET", u, responder)

				pats, err := c.ListProjectAccessTokens(projectID) // Project ID doesn't matter

				Expect(err).ToNot(HaveOccurred())
				Expect(pats).To(HaveLen(0))
			})
		})

		When("there are tokens attached to the project", func() {
			It("lists the correct tokens", func() {
				s := fixture("project_access_tokens/list.json")
				stringResponse := httpmock.NewStringResponse(200, s)
				stringResponse.Header.Add("Content-Type", "application/json")
				responder := httpmock.ResponderFromResponse(stringResponse)
				httpmock.RegisterResponder("GET", u, responder)

				expected := []ProjectAccessToken{
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
				actual, err := c.ListProjectAccessTokens(projectID)
				log.Debug().
					Interface("actual", actual).
					Interface("expected", expected).
					Msg("List project access tokens")

				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(HaveLen(3))
				Expect(actual).To(ContainElements(expected))
			})
		})
	})

	Context("getting PAT by project ID and name", func() {

	})

})

/*

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
