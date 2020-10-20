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
	u := apiUrl + pathInvitationCreate
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
		rs := responseFromFixture("invitation/create.json", http.StatusOK)
		return rs, nil
	}
	httpmock.RegisterResponder("POST", u, r)
	inv, err := s.client.CreateInvitation(teamID, email)
	s.Nil(err)
	s.Equal(email, inv.Email)

	s.checkServerErrors("POST", u, func() error {
		_, err = s.client.CreateInvitation(teamID, email)
		return err
	})
}
