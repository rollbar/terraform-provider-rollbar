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
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

// Invitation represents an invitation for a user to join a Rollbar team.
type Invitation struct {
	ID           int    `json:"id"`
	FromUserID   int    `json:"from_user_id"`
	TeamID       int    `json:"team_id"`
	Email        string `json:"to_email"`
	Status       string `json:"status"`
	DateCreated  int    `json:"date_created"`
	DateRedeemed int    `json:"date_redeemed"`
}

// CreateInvitation sends an invitation to a user.
func (c *RollbarApiClient) CreateInvitation(teamID int, email string) (Invitation, error) {
	l := log.With().
		Int("teamID", teamID).
		Str("email", email).
		Logger()
	l.Debug().Msg("Creating new invitation")

	u := apiUrl + pathInvitationCreate
	var inv Invitation
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamId": strconv.Itoa(teamID),
		}).
		SetBody(map[string]string{
			"email": email,
		}).
		SetResult(invitationResponse{}).
		SetError(ErrorResult{}).
		Post(u)
	if err != nil {
		l.Err(err).Msg("Error creating invitation")
		return inv, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		return inv, err
	}
	r := resp.Result().(*invitationResponse)
	inv = r.Result
	return inv, nil
}

func (c *RollbarApiClient) ReadInvitation(inviteID int) (inv Invitation, err error) {
	l := log.With().
		Int("inviteID", inviteID).
		Logger()
	l.Debug().Msg("Reading invitation from Rollbar API")
	u := apiUrl + pathInvitationRead
	u = strings.ReplaceAll(u, "{inviteId}", strconv.Itoa(inviteID))
	resp, err := c.Resty.R().
		SetResult(invitationResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	err = errorFromResponse(resp)
	if err != nil {
		return
	}
	inv = resp.Result().(invitationResponse).Result
	return
}

/*
 * Containers for unmarshalling Rollbar API responses
 */

type invitationResponse struct {
	Error  int        `json:"err"`
	Result Invitation `json:"result"`
}

type invitationListResponse struct {
	Error  int          `json:"err"`
	Result []Invitation `json:"result"`
}
