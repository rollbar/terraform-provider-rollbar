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
	"net/http"
	"strconv"
	"strings"
)

// Invitation represents an invitation for a user to join a Rollbar team.
type Invitation struct {
	ID           int    `json:"id"`
	FromUserID   int    `json:"from_user_id"`
	TeamID       int    `json:"team_id"`
	ToEmail      string `json:"to_email"`
	Status       string `json:"status"`
	DateCreated  int    `json:"date_created"`
	DateRedeemed int    `json:"date_redeemed"`
}

// ListInvitations lists all invitations for a Rollbar team.
func (c *RollbarApiClient) ListInvitations(teamID int) (invs []Invitation, err error) {
	l := log.With().
		Int("teamID", teamID).
		Logger()
	l.Debug().Msg("Listing invitations")
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamID": strconv.Itoa(teamID),
		}).
		SetResult(invitationListResponse{}).
		SetError(ErrorResult{}).
		Get(c.BaseURL + pathInvitations)
	if err != nil {
		l.Err(err).Msg("Error listing invitations")
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).
			Str("status", resp.Status()).
			Msg("Error listing invitations")
		return
	}
	r := resp.Result().(*invitationListResponse)
	invs = r.Result
	l.Debug().
		Int("invitation_count", len(invs)).
		Msg("Successfully listed invitations")
	return
}

// ListPendingInvitations lists a Rollbar team's pending invitations.
func (c *RollbarApiClient) ListPendingInvitations(teamID int) ([]Invitation, error) {
	l := log.With().Int("teamID", teamID).Logger()
	l.Debug().Msg("Listing pending invitations")
	var pending []Invitation
	all, err := c.ListInvitations(teamID)
	if err != nil {
		l.Err(err).Send()
		return pending, err
	}
	for _, inv := range all {
		if inv.Status == "pending" {
			pending = append(pending, inv)
		}
	}
	l.Debug().
		Int("invitation_count", len(pending)).
		Msg("Successfully listed pending invitations")
	return pending, nil
}

// FindPendingInvitations finds pending Rollbar team invitations for the given
// email.
func (c *RollbarApiClient) FindPendingInvitations(email string) ([]Invitation, error) {
	l := log.With().Str("email", email).Logger()
	l.Debug().Msg("Finding pending invitations")
	var pending []Invitation
	all, err := c.FindInvitations(email)
	if err != nil {
		l.Err(err).Send()
		return pending, err
	}
	for _, inv := range all {
		if inv.Status == "pending" {
			pending = append(pending, inv)
		}
	}
	l.Debug().
		Int("invitation_count", len(pending)).
		Msg("Successfully found pending invitations")
	return pending, nil
}

// CreateInvitation sends a Rollbar team invitation to a user.
func (c *RollbarApiClient) CreateInvitation(teamID int, email string) (Invitation, error) {
	l := log.With().
		Int("teamID", teamID).
		Str("email", email).
		Logger()
	l.Debug().Msg("Creating new invitation")

	u := c.BaseURL + pathInvitations
	var inv Invitation
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamID": strconv.Itoa(teamID),
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
		l.Err(err).Send()
		return inv, err
	}
	r := resp.Result().(*invitationResponse)
	inv = r.Result
	l.Debug().Msg("Successfully created new invitation")
	return inv, nil
}

// ReadInvitation reads a Rollbar team invitation from the API.
func (c *RollbarApiClient) ReadInvitation(inviteID int) (inv Invitation, err error) {
	l := log.With().
		Int("inviteID", inviteID).
		Logger()
	l.Debug().Msg("Reading invitation from Rollbar API")
	u := c.BaseURL + pathInvitation
	u = strings.ReplaceAll(u, "{inviteID}", strconv.Itoa(inviteID))
	resp, err := c.Resty.R().
		SetResult(invitationResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		l.Err(err).Msg("Error reading invitation from API")
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Msg("Error reading invitation from API")
		return
	}
	inv = resp.Result().(*invitationResponse).Result
	l.Debug().
		Interface("invitation", inv).
		Msg("Successfully read invitation from API")
	return
}

// DeleteInvitation is an alias for CancelInvitation.
func (c *RollbarApiClient) DeleteInvitation(id int) (err error) {
	return c.CancelInvitation(id)
}

// CancelInvitation cancels a Rollbar team invitation.
func (c *RollbarApiClient) CancelInvitation(id int) (err error) {
	l := log.With().Int("id", id).Logger()
	l.Debug().Msg("Canceling invitation")

	u := c.BaseURL + pathInvitation
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"inviteID": strconv.Itoa(id),
		}).
		SetError(ErrorResult{}).
		Delete(u)
	if err != nil {
		l.Err(err).Msg("Error canceling invitation")
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		// If the invite has already been canceled, API returns HTTP status '422
		// Unprocessable Entity'.  This is considered success.
		statusUnprocessable := resp.StatusCode() == http.StatusUnprocessableEntity
		alreadyCanceledMsg := strings.Contains(err.Error(), "Invite already cancelled")
		if statusUnprocessable && alreadyCanceledMsg {
			l.Debug().Msg("invite already cancelled")
			return nil
		}
		l.Err(err).
			Interface("error", err).
			Str("status", resp.Status()).
			Int("status_code", resp.StatusCode()).
			Msg("Error canceling invitation")
		return
	}
	l.Debug().
		Msg("Successfully canceled invitation")
	return
}

// FindInvitations finds all Rollbar team invitations for a given email. Note
// this method is quite inefficient, as it must read all invitations for all
// teams.
func (c *RollbarApiClient) FindInvitations(email string) (invs []Invitation, err error) {
	// API converts all invited emails to lowercase.
	// https://github.com/rollbar/terraform-provider-rollbar/issues/139
	email = strings.ToLower(email)
	l := log.With().
		Str("email", email).
		Logger()

	l.Debug().Msg("Finding invitations")
	teams, err := c.ListCustomTeams()
	if err != nil {
		l.Err(err).Send()
		return
	}
	var allInvs []Invitation
	for _, t := range teams {
		teamInvs, err := c.ListInvitations(t.ID)
		// Team may have been deleted by another process after we listed all
		// teams, but before we queried the team for invitations.  Therefore we
		// ignore ErrNotFound.
		// https://github.com/rollbar/terraform-provider-rollbar/issues/88
		if err != nil && err != ErrNotFound {
			l.Err(err).
				Str("team_name", t.Name).
				Msg("error finding invitations")
			return invs, err
		}
		allInvs = append(allInvs, teamInvs...)
	}
	for _, inv := range allInvs {
		if inv.ToEmail == email {
			invs = append(invs, inv)
		}
	}
	if len(invs) == 0 {
		return invs, ErrNotFound
	}
	l.Debug().
		Int("invitation_count", len(invs)).
		Msg("Successfully found invitations")
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
