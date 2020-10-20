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
)

// Invite represents an invitation for a user to join a Rollbar team.
type Invite struct {
	ID           int    `json:"id"`
	FromUserID   int    `json:"from_user_id"`
	TeamID       int    `json:"team_id"`
	Email        string `json:"to_email"`
	Status       string `json:"status"`
	DateCreated  int    `json:"date_created"`
	DateRedeemed int    `json:"date_redeemed"`
}

// CreateInvite sends an invitation to a user.
func (c *RollbarApiClient) CreateInvite(teamID int, email string) (Invite, error) {
	l := log.With().
		Int("teamID", teamID).
		Str("email", email).
		Logger()
	l.Debug().Msg("Creating new invite")

	u := apiUrl + pathInviteCreate
	var inv Invite
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamId": strconv.Itoa(teamID),
		}).
		SetBody(map[string]string{
			"email": email,
		}).
		SetResult(inviteCreateResponse{}).
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
	r := resp.Result().(*inviteCreateResponse)
	inv = r.Result
	return inv, nil
}

func (c *RollbarApiClient) ReadInvite(inviteID int) (inv Invite, err error) {
	l := log.With().
		Int("inviteID", inviteID).
		Logger()
	l.Debug().Msg("Reading invitation from Rollbar API")
	return
}

/*
 * Containers for unmarshalling Rollbar API responses
 */

type inviteCreateResponse struct {
	Error  int    `json:"err"`
	Result Invite `json:"result"`
}
