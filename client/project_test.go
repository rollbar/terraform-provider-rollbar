package client

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {
	It("Lists all projects", func() {
		u := apiUrl + pathProjectList
		httpmock.RegisterResponder("GET", u, httpmock.NewStringResponder(200, fixture("projects/list.json")))
		expected := []*Project{
			{
				Id:           12112,
				AccountId:    8608,
				DateCreated:  1407933721,
				DateModified: 1457475137,
				Name:         "",
			},
			{
				Id:           106671,
				AccountId:    8608,
				DateCreated:  1489139046,
				DateModified: 1549293583,
				Name:         "Client-Config",
			},
			{
				Id:           12116,
				AccountId:    8608,
				DateCreated:  1407933922,
				DateModified: 1556814300,
				Name:         "My",
			},
		}
		actual, err := c.ListProjects()
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(Equal(expected))
	})
})

/*
func TestListProjects(t *testing.T) {
	teardown := setup()
	defer teardown()

	handURL := "/projects"

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("projects/list.json"))
	})

	response, err := c.ListProjects()

	if err != nil {
		t.Fatal(err)
	}

	expected := []*Project{
		{
			ID:           12112,
			AccountID:    8608,
			DateCreated:  1407933721,
			DateModified: 1457475137,
			Name:         "",
		},
		{
			ID:           106671,
			AccountID:    8608,
			DateCreated:  1489139046,
			DateModified: 1549293583,
			Name:         "Client-Config",
		},
		{
			ID:           12116,
			AccountID:    8608,
			DateCreated:  1407933922,
			DateModified: 1556814300,
			Name:         "My",
		},
	}

	if !reflect.DeepEqual(response, expected) {
		t.Errorf("expected response %v, got %v.", response, expected)
	}
}

func TestGetProjectByName(t *testing.T) {
	teardown := setup()
	defer teardown()

	handURL := "/projects"

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("projects/list.json"))
	})

	examples := []struct {
		name     string
		expected *Project
	}{
		{
			name:     "ProjectDoesNotExist",
			expected: nil,
		},
		{
			name: "Client-Config",
			expected: &Project{
				ID:           106671,
				AccountID:    8608,
				DateCreated:  1489139046,
				DateModified: 1549293583,
				Name:         "Client-Config",
			},
		},
	}

	for _, example := range examples {
		actual, err := c.GetProjectByName(example.name)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, example.expected) {
			t.Errorf("expected project %v, got %v.", example.expected, actual)
		}
	}
}


 */