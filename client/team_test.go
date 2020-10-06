package client

import (
	"github.com/jarcoal/httpmock"
	"net/http"
)

func (s *Suite) TestCreateTeam() {

	// Setup API mock
	u := apiUrl + pathTeamCreate
	expected := Team{
		ID:          676974,
		AccountID:   317418,
		Name:        "foobar",
		AccessLevel: TeamAccessStandard,
	}
	// FIXME: currently API returns `200 OK` on successful create; but it should
	//  instead return `201 Created`.
	//  https://github.com/rollbar/terraform-provider-rollbar/issues/8
	sr := httpmock.NewStringResponse(http.StatusOK, teamCreateResponse)
	sr.Header.Add("Content-Type", "application/json")
	r := httpmock.ResponderFromResponse(sr)
	httpmock.RegisterResponder("POST", u, r)

	// Successful create
	actual, err := s.client.CreateTeam("foobar", TeamAccessStandard)
	s.Nil(err)
	s.Equal(expected, actual)

	// Invalid name
	_, err = s.client.CreateTeam("", TeamAccessStandard)
	s.NotNil(err)

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("POST", u, r)
	_, err = s.client.CreateTeam("foobar", TeamAccessStandard)
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.CreateTeam("foobar", TeamAccessStandard)
	s.NotNil(err)
}

const teamCreateResponse = `
{
    "err": 0,
    "result": {
        "access_level": "standard",
        "account_id": 317418,
        "id": 676974,
        "name": "foobar"
    }
}
`
