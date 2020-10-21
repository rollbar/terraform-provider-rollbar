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

// TestCreateInvitation tests creating a Rollbar team invitation.
func (s *Suite) TestCreateInvitation() {
	teamID := 572097
	email := "test@rollbar.com"
	u := apiUrl + pathInvitations
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

// TestReadInvitation tests reading a Rollbar team invitation from the API.
func (s *Suite) TestReadInvitation() {
	id := 153650
	u := apiUrl + pathInvitation
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

func (s *Suite) TestCancelInvitation() {
	invitationId := 153650
	u := apiUrl + pathInvitation
	u = strings.ReplaceAll(u, "{inviteId}", strconv.Itoa(invitationId))

	r := responderFromFixture("invitation/cancel.json", http.StatusOK)
	httpmock.RegisterResponder("DELETE", u, r)
	err := s.client.CancelInvitation(invitationId)
	s.Nil(err)

	// DeleteInvitation is an alias for CancelInvitation.
	err = s.client.DeleteInvitation(invitationId)
	s.Nil(err)

	s.checkServerErrors("DELETE", u, func() error {
		err := s.client.CancelInvitation(invitationId)
		return err
	})

}
