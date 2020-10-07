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
	u := apiUrl + pathTeamCreate
	expected := Team{
		ID:          676974,
		AccountID:   317418,
		Name:        teamName,
		AccessLevel: TeamAccessStandard,
	}
	// FIXME: currently API returns `200 OK` on successful create; but it should
	//  instead return `201 Created`.
	//  https://github.com/rollbar/terraform-provider-rollbar/issues/8
	sr := responseFromFixture("team/create.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		type body struct {
			Name string
		}
		b := body{}
		err := json.NewDecoder(req.Body).Decode(&b)
		s.Nil(err)
		s.Equal(teamName, b.Name)
		return sr, nil
	}
	httpmock.RegisterResponder("POST", u, r)

	// Successful create
	actual, err := s.client.CreateTeam(teamName, TeamAccessStandard)
	s.Nil(err)
	s.Equal(expected, actual)

	// Invalid name
	_, err = s.client.CreateTeam("", TeamAccessStandard)
	s.NotNil(err)

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("POST", u, r)
	_, err = s.client.CreateTeam(teamName, TeamAccessStandard)
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.CreateTeam(teamName, TeamAccessStandard)
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("POST", u, r)
	_, err = s.client.CreateTeam(teamName, TeamAccessStandard)
	s.Equal(ErrUnauthorized, err)
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
			ID:          676974,
			AccountID:   317418,
			Name:        "foobar",
			AccessLevel: TeamAccessStandard,
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

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ListTeams()
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.ListTeams()
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ListTeams()
	s.Equal(ErrUnauthorized, err)
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
		AccessLevel: TeamAccessStandard,
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

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadTeam(teamId)
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.ReadTeam(teamId)
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadTeam(teamId)
	s.Equal(ErrUnauthorized, err)
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

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("DELETE", u, r)
	err = s.client.DeleteTeam(teamId)
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	err = s.client.DeleteTeam(teamId)
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("DELETE", u, r)
	err = s.client.DeleteTeam(teamId)
	s.Equal(ErrUnauthorized, err)
}
