/*
 * Copyright (c) 2020 Rollbar, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package client

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func (s *Suite) TestListProjects() {
	u := s.client.BaseURL + pathProjectList

	// Success
	r := responderFromFixture("project/list.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	expected := []Project{
		{
			ID:           411704,
			Name:         "bar",
			AccountID:    317418,
			Status:       "enabled",
			DateCreated:  1602085345,
			DateModified: 1602085345,
		},
		{
			ID:           411703,
			Name:         "foo",
			AccountID:    317418,
			Status:       "enabled",
			DateCreated:  1602085340,
			DateModified: 1602085340,
		},
	}
	actual, err := s.client.ListProjects()
	s.Nil(err)
	s.Len(actual, len(expected))
	s.ElementsMatch(expected, actual)

	s.checkServerErrors("GET", u, func() error {
		_, err = s.client.ListProjects()
		return err
	})
}

func (s *Suite) TestCreateProject() {
	u := s.client.BaseURL + pathProjectCreate
	name := "baz"

	// Success
	// FIXME: The actual Rollbar API sends http.StatusOK; but it
	//  _should_ send http.StatusCreated
	rs := responseFromFixture("project/create.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		p := Project{}
		err := json.NewDecoder(req.Body).Decode(&p)
		s.Nil(err)
		s.Equal(name, p.Name)
		return rs, nil
	}
	httpmock.RegisterResponder("POST", u, r)
	proj, err := s.client.CreateProject(name)
	s.Nil(err)
	s.Equal(name, proj.Name)

	s.checkServerErrors("POST", u, func() error {
		_, err = s.client.CreateProject(name)
		return err
	})
}

func (s *Suite) TestReadProject() {
	expected := Project{
		AccountID:    317418,
		DateCreated:  1602086539,
		DateModified: 1602086539,
		ID:           411708,
		Name:         "baz",
		Status:       "enabled",
	}
	u := s.client.BaseURL + pathProjectRead
	u = strings.ReplaceAll(u, "{projectID}", strconv.Itoa(expected.ID))

	// Success
	r := responderFromFixture("project/read.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	actual, err := s.client.ReadProject(expected.ID)
	s.Nil(err)
	s.Equal(&expected, actual)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ReadProject(expected.ID)
		return err
	})

	// Deleted project API bug
	// FIXME: https://github.com/rollbar/terraform-provider-rollbar/issues/23
	r = responderFromFixture("project/read_deleted.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadProject(expected.ID)
	s.Equal(ErrNotFound, err)
}

func (s *Suite) TestDeleteProject() {
	delID := gofakeit.Number(0, 1000000)
	urlDel := s.client.BaseURL + pathProjectDelete
	urlDel = strings.ReplaceAll(urlDel, "{projectID}", strconv.Itoa(delID))

	// Success
	r := responderFromFixture("project/delete.json", http.StatusOK)
	httpmock.RegisterResponder("DELETE", urlDel, r)
	err := s.client.DeleteProject(delID)
	s.Nil(err)

	s.checkServerErrors("DELETE", urlDel, func() error {
		return s.client.DeleteProject(delID)
	})
}
func (s *Suite) TestFindProjectTeamIDs() {
	var u string // URL
	var r httpmock.Responder
	projectID := 423092
	teamID := 689492
	expected := []int{teamID}

	// Mock list team projects
	u = s.client.BaseURL + pathTeamProjects
	u = strings.ReplaceAll(u, "{teamID}", strconv.Itoa(teamID))
	r = responderFromFixture("team/list_projects_689492.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	u = s.client.BaseURL + pathTeamProjects
	u = strings.ReplaceAll(u, "{teamID}", "689493")
	r = responderFromFixture("team/list_projects_689493.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)

	// Mock list teams
	u = s.client.BaseURL + pathTeamList
	r = responderFromFixture("team/list_2.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)

	actual, err := s.client.FindProjectTeamIDs(projectID)
	s.Nil(err)
	s.Equal(expected, actual)

	expectedCallCount :=
		map[string]int{
			"GET https://api.rollbar.com/api/1/team/689492/projects": 1,
			"GET https://api.rollbar.com/api/1/team/689493/projects": 1,
			"GET https://api.rollbar.com/api/1/teams":                1,
		}
	actualCallCount := httpmock.GetCallCountInfo()
	for call, count := range expectedCallCount {
		s.Equal(count, actualCallCount[call])
	}

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.FindProjectTeamIDs(projectID)
		return err
	})
}

func TestUpdateProjectTeams(t *testing.T) {
	// Setup go-vcr
	httpmock.Deactivate()
	r, err := recorder.New("fixtures/vcr/update_project_teams")
	assert.Nil(t, err)
	defer func() {
		err := r.Stop()
		if err != nil {
			panic(err)
		}
	}()
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Rollbar-Access-Token")
		return nil
	})

	c := NewClient(DefaultBaseURL, os.Getenv("ROLLBAR_API_KEY"))
	c.Resty.GetClient().Transport = r

	prefix := "tf-acc-test-updateprojectteams"
	projectName := prefix
	team0Name := prefix + "-0"
	team1Name := prefix + "-1"
	team2Name := prefix + "-2"

	project, err := c.CreateProject(projectName)
	assert.Nil(t, err)
	team0, err := c.CreateTeam(team0Name, "standard")
	assert.Nil(t, err)
	team1, err := c.CreateTeam(team1Name, "standard")
	assert.Nil(t, err)
	team2, err := c.CreateTeam(team2Name, "standard")
	assert.Nil(t, err)
	err = c.AssignTeamToProject(team0.ID, project.ID)
	assert.Nil(t, err)
	err = c.AssignTeamToProject(team1.ID, project.ID)
	assert.Nil(t, err)

	expectedTeamIDs := []int{team1.ID, team2.ID}
	err = c.UpdateProjectTeams(project.ID, expectedTeamIDs)
	assert.Nil(t, err)
	actualTeamIDs, err := c.FindProjectTeamIDs(project.ID)
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedTeamIDs, actualTeamIDs)

	// Bad project ID
	err = c.UpdateProjectTeams(0, expectedTeamIDs)
	assert.NotNil(t, err)
	// Bad team ID
	err = c.UpdateProjectTeams(project.ID, []int{0})
	assert.NotNil(t, err)

	// Cleanup
	for _, teamID := range []int{team0.ID, team1.ID, team2.ID} {
		err = c.DeleteTeam(teamID)
		assert.Nil(t, err)
	}
	err = c.DeleteProject(project.ID)
	assert.Nil(t, err)
}
