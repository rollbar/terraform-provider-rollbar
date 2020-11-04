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
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"strings"
)

func (s *Suite) TestCreateTeam() {
	// Setup API mock
	teamName := "foobar"
	accessLevel := "standard"
	u := apiUrl + pathTeamCreate
	expected := Team{
		ID:          676974,
		AccountID:   317418,
		Name:        teamName,
		AccessLevel: accessLevel,
	}
	// FIXME: currently API returns `200 OK` on successful create; but it should
	//  instead return `201 Created`.
	//  https://github.com/rollbar/terraform-provider-rollbar/issues/8
	sr := responseFromFixture("team/create.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		b := make(map[string]interface{})
		err := json.NewDecoder(req.Body).Decode(&b)
		s.Nil(err)
		s.Equal(teamName, b["name"])
		s.Equal(accessLevel, b["access_level"])
		return sr, nil
	}
	httpmock.RegisterResponder("POST", u, r)

	// Successful create
	actual, err := s.client.CreateTeam(teamName, "standard")
	s.Nil(err)
	s.Equal(expected, actual)

	// Invalid name
	_, err = s.client.CreateTeam("", "standard")
	s.NotNil(err)

	s.checkServerErrors("POST", u, func() error {
		_, err = s.client.CreateTeam(teamName, "standard")
		return err
	})
}

func (s *Suite) TestListTeams() {
	// Setup API mock
	u := apiUrl + pathTeamList
	expected := []Team{
		{
			AccessLevel: "everyone",
			AccountID:   317418,
			ID:          662037,
			Name:        "Everyone",
		},
		{
			ID:          676971,
			AccountID:   317418,
			Name:        "my-test-team",
			AccessLevel: "standard",
		},
		{
			AccessLevel: "owner",
			AccountID:   317418,
			ID:          662036,
			Name:        "Owners",
		},
	}
	r := responderFromFixture("team/list.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)

	// Successful list
	actual, err := s.client.ListTeams()
	s.Nil(err)
	s.Equal(expected, actual)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ListTeams()
		return err
	})
}

func (s *Suite) TestReadTeam() {
	// Setup API mock
	teamId := 676974
	u := apiUrl + pathTeamRead
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(teamId))
	expected := Team{
		ID:          676974,
		AccountID:   317418,
		Name:        "foobar",
		AccessLevel: "standard",
	}
	r := responderFromFixture("team/read.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)

	// Successful create
	actual, err := s.client.ReadTeam(teamId)
	s.Nil(err)
	s.Equal(expected, actual)

	// Invalid ID
	_, err = s.client.ReadTeam(0)
	s.NotNil(err)

	// FIXME: Workaround API bug
	//  https://github.com/rollbar/terraform-provider-rollbar/issues/79
	r = httpmock.NewJsonResponderOrPanic(http.StatusForbidden,
		ErrorResult{Err: 1, Message: "Team not found in this account."})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadTeam(teamId)
	s.Equal(ErrNotFound, err)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ReadTeam(teamId)
		return err
	})
}

func (s *Suite) TestDeleteTeam() {
	// Setup API mock
	teamId := 676974
	u := apiUrl + pathTeamDelete
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(teamId))
	r := responderFromFixture("team/delete.json", http.StatusOK)
	httpmock.RegisterResponder("DELETE", u, r)

	// Successful delete
	err := s.client.DeleteTeam(teamId)
	s.Nil(err)

	// Invalid ID
	err = s.client.DeleteTeam(0)
	s.NotNil(err)

	s.checkServerErrors("DELETE", u, func() error {
		return s.client.DeleteTeam(teamId)
	})
}

// TestAssignUserToTeam tests assigning a user to a Rollbar team.
func (s *Suite) TestAssignUserToTeam() {
	teamID := 676971
	userID := 238101

	// Successful assignment
	u := apiUrl + pathTeamUser
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(teamID))
	u = strings.ReplaceAll(u, "{userId}", strconv.Itoa(userID))
	r := responderFromFixture("team/assign_user.json", http.StatusOK)
	httpmock.RegisterResponder("PUT", u, r)
	err := s.client.AssignUserToTeam(teamID, userID)
	s.Nil(err)

	s.checkServerErrors("PUT", u, func() error {
		err = s.client.AssignUserToTeam(teamID, userID) // non-existent user
		return err
	})

	// API returns status 403 when the team or user is not found.
	u = apiUrl + pathTeamUser
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(teamID))
	u = strings.ReplaceAll(u, "{userId}", "0")
	r = responderFromFixture("team/assign_user_not_found.json", http.StatusForbidden)
	httpmock.RegisterResponder("PUT", u, r)
	err = s.client.AssignUserToTeam(teamID, 0) // non-existent user
	s.Equal(ErrNotFound, err)
}

// TestRemoveUserFromTeam tests removing a user from a Rollbar team.
func (s *Suite) TestRemoveUserFromTeam() {
	teamID := 676971
	userID := 238101

	// Successful assignment
	u := apiUrl + pathTeamUser
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(teamID))
	u = strings.ReplaceAll(u, "{userId}", strconv.Itoa(userID))
	r := responderFromFixture("team/remove_user.json", http.StatusOK)
	httpmock.RegisterResponder("DELETE", u, r)
	err := s.client.RemoveUserFromTeam(teamID, userID)
	s.Nil(err)

	s.checkServerErrors("DELETE", u, func() error {
		err = s.client.RemoveUserFromTeam(teamID, userID) // non-existent user
		return err
	})

	// API returns status 422 when the team or user is not found.
	u = apiUrl + pathTeamUser
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(teamID))
	u = strings.ReplaceAll(u, "{userId}", "0")
	r = responderFromFixture("team/remove_user_not_found.json", http.StatusUnprocessableEntity)
	httpmock.RegisterResponder("DELETE", u, r)
	err = s.client.RemoveUserFromTeam(teamID, 0) // non-existent user
	s.Equal(ErrNotFound, err)
}

// TestListCustomTeams tests listing custom defined teams.
func (s *Suite) TestListCustomTeams() {
	u := apiUrl + pathTeamList
	expected := []Team{
		{
			ID:          676971,
			AccountID:   317418,
			Name:        "my-test-team",
			AccessLevel: "standard",
		},
	}
	r := responderFromFixture("team/list.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)

	actual, err := s.client.ListCustomTeams()
	s.Nil(err)
	s.Equal(expected, actual)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ListCustomTeams()
		return err
	})
}
