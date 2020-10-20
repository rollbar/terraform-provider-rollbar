package client

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"strings"
)

func (s *Suite) TestCreateInvite() {
	teamID := 572097
	email := "test@rollbar.com"
	u := apiUrl + pathInviteCreate
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(teamID))

	// Success
	// FIXME: The actual Rollbar API sends http.StatusOK; but it
	//  _should_ send http.StatusCreated
	r := func(req *http.Request) (*http.Response, error) {
		m := make(map[string]string)
		err := json.NewDecoder(req.Body).Decode(&m)
		s.Nil(err)
		s.Contains(m, "email")
		s.Equal(email, m["email"])
		rs := responseFromFixture("invite/create.json", http.StatusOK)
		return rs, nil
	}
	httpmock.RegisterResponder("POST", u, r)
	inv, err := s.client.CreateInvite(teamID, email)
	s.Nil(err)
	s.Equal(email, inv.Email)

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("POST", u, r)
	_, err = s.client.CreateInvite(teamID, email)
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.CreateInvite(teamID, email)
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("POST", u, r)
	_, err = s.client.CreateInvite(teamID, email)
	s.Equal(ErrUnauthorized, err)
}
