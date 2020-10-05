package client

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"net/http"
)

var _ = Describe("Project", func() {
	It("Lists all projects", func() {
		u := apiUrl + pathProjectList
		s := fixture("projects/list.json")
		stringResponse := httpmock.NewStringResponse(200, s)
		stringResponse.Header.Add("Content-Type", "application/json")
		responder := httpmock.ResponderFromResponse(stringResponse)
		httpmock.RegisterResponder("GET", u, responder)

		expected := []Project{
			{
				Id:           106671,
				AccountId:    8608,
				DateCreated:  1489139046,
				DateModified: 1549293583,
				Name:         "Client-Config",
				Status:       "enabled",
			},
			{
				Id:           12116,
				AccountId:    8608,
				DateCreated:  1407933922,
				DateModified: 1556814300,
				Name:         "My",
				Status:       "enabled",
			},
		}
		actual, err := c.ListProjects()
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(ContainElements(expected))
		Expect(actual).To(HaveLen(len(expected)))
	})

	When("creating a new project", func() {
		u := apiUrl + pathProjectCreate
		name := gofakeit.HackerNoun()

		Context("and creation succeeds", func() {
			s := fmt.Sprintf(fixture("projects/create.json"), name)
			// FIXME: The actual Rollbar API sends http.StatusOK; but it
			//  _should_ send http.StatusCreated
			stringResponse := httpmock.NewStringResponse(http.StatusOK, s)
			stringResponse.Header.Add("Content-Type", "application/json")
			responder := func(req *http.Request) (*http.Response, error) {
				//proj := make(map[string]string)
				proj := Project{}
				if err := json.NewDecoder(req.Body).Decode(&proj); err != nil {
					return nil, err
				}
				if proj.Name != name {
					msg := "incorrect name sent to API"
					log.Error().
						Str("expected", name).
						Str("actual", proj.Name).
						Msg(msg)
					return nil, fmt.Errorf(msg)
				}
				return stringResponse, nil
			}
			It("is created correctly", func() {
				httpmock.RegisterResponder("POST", u, responder)
				proj, err := c.CreateProject(name)
				Expect(err).NotTo(HaveOccurred())
				Expect(proj.Name).To(Equal(name))
			})
		})

		Context("and creation fails", func() {
			s := fixture("projects/create.json")
			// FIXME: The actual Rollbar API sends http.StatusOK; but it
			//  _should_ send http.StatusCreated
			stringResponse := httpmock.NewStringResponse(http.StatusInternalServerError, s)
			stringResponse.Header.Add("Content-Type", "application/json")
			responder := httpmock.ResponderFromResponse(stringResponse)
			It("handles the error cleanly", func() {
				httpmock.RegisterResponder("POST", u, responder)
				_, err := c.CreateProject(name)
				Expect(err).To(MatchError(&ErrorResult{}))
			})
		})
	})
})
