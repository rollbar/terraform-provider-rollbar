package client

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
})
