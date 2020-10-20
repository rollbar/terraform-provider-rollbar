package client

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"strings"
)

func (s *Suite) TestCreateInvitation() {
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
	s.Equal(email, inv.ToEmail)

	s.checkServerErrors("POST", u, func() error {
		_, err = s.client.CreateInvitation(teamID, email)
		return err
	})
}
func (s *Suite) TestReadInvitation() {
	id := 153650
	u := apiUrl + pathInvitationRead
	u = strings.ReplaceAll(u, "{inviteId}", strconv.Itoa(id))

	// Success
	r := responderFromFixture("invitation/read.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	actual, err := s.client.ReadInvitation(id)
	s.Nil(err)
	expected := Invitation{
		DateCreated:  1603192477,
		DateRedeemed: 0,
		FromUserID:   5325,
		ID:           153650,
		Status:       "pending",
		TeamID:       676971,
		ToEmail:      "test@rollbar.com",
	}
	s.Equal(expected, actual)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ReadInvitation(id)
		return err
	})
}
