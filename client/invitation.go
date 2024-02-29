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
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
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

// ListAllInvitationsPerEmail lists all invitations for all Rollbar teams.
func (c *RollbarAPIClient) ListAllInvitationsPerEmail(email string) (invs []Invitation, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	hasNextPage := true
	page := 1

	l := log.With().
		Str("email", email).
		Logger()
	l.Debug().Msg("Listing invitations")

	for hasNextPage {
		resp, err := c.Resty.R().
			SetResult(invitationListResponse{}).
			SetError(ErrorResult{}).
			SetQueryParams(map[string]string{
				"page":  strconv.Itoa(page),
				"email": email,
			}).
			Get(c.BaseURL + pathInvitations)
		if err != nil {
			l.Err(err).Msg("Error listing invitations")
			return nil, err
		}
		err = errorFromResponse(resp)
		if err != nil {
			l.Err(err).
				Str("status", resp.Status()).
				Msg("Error listing invitations")
			return nil, err
		}
		r := resp.Result().(*invitationListResponse)
		hasNextPage = len(r.Result) > 0
		page++
		invs = append(invs, r.Result...)
	}
	l.Debug().
		Int("invitation_count", len(invs)).
		Msg("Successfully listed invitations")
	return invs, nil
}

// ListInvitations lists all invitations for a Rollbar team.
func (c *RollbarAPIClient) ListInvitations(teamID int) (invs []Invitation, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	hasNextPage := true
	page := 1

	l := log.With().
		Int("teamID", teamID).
		Logger()
	l.Debug().Msg("Listing invitations")

	for hasNextPage {
		resp, err := c.Resty.R().
			SetPathParams(map[string]string{
				"teamID": strconv.Itoa(teamID),
			}).
			SetResult(invitationListResponse{}).
			SetError(ErrorResult{}).
			SetQueryParams(map[string]string{
				"page": strconv.Itoa(page)}).
			Get(c.BaseURL + pathTeamInvitations)
		if err != nil {
			l.Err(err).Msg("Error listing invitations")
			return nil, err
		}
		err = errorFromResponse(resp)
		if err != nil {
			l.Err(err).
				Str("status", resp.Status()).
				Msg("Error listing invitations")
			return nil, err
		}
		r := resp.Result().(*invitationListResponse)
		hasNextPage = len(r.Result) > 0
		page++
		invs = append(invs, r.Result...)
	}
	l.Debug().
		Int("invitation_count", len(invs)).
		Msg("Successfully listed invitations")
	return invs, nil
}

// ListPendingInvitations lists a Rollbar team's pending invitations.
func (c *RollbarAPIClient) ListPendingInvitations(teamID int) ([]Invitation, error) {
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
func (c *RollbarAPIClient) FindPendingInvitations(email string) ([]Invitation, error) {
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
func (c *RollbarAPIClient) CreateInvitation(teamID int, email string) (Invitation, error) {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().
		Int("teamID", teamID).
		Str("email", email).
		Logger()
	l.Debug().Msg("Creating new invitation")

	u := c.BaseURL + pathTeamInvitations
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
func (c *RollbarAPIClient) ReadInvitation(inviteID int) (inv Invitation, err error) {
	c.m.Lock()
	defer c.m.Unlock()
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
func (c *RollbarAPIClient) DeleteInvitation(id int) (err error) {
	return c.CancelInvitation(id)
}

// CancelInvitation cancels a Rollbar team invitation.
func (c *RollbarAPIClient) CancelInvitation(id int) (err error) {
	c.m.Lock()
	defer c.m.Unlock()
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
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		// If the invite has already been canceled, API returns HTTP status '422
		// Unprocessable Entity'.  This is considered success.
		statusUnprocessable := resp.StatusCode() == http.StatusUnprocessableEntity
		alreadyCanceledMsg := strings.Contains(err.Error(), "Invite already canceled")
		if statusUnprocessable && alreadyCanceledMsg {
			l.Debug().Msg("invite already canceled")
			return nil
		}
		l.Err(err).
			Interface("error", err).
			Str("status", resp.Status()).
			Int("status_code", resp.StatusCode()).
			Msg("Error canceling invitation")
		return err
	}
	l.Debug().
		Msg("Successfully canceled invitation")
	return nil
}

// FindInvitations finds all Rollbar team invitations for a given email.
func (c *RollbarAPIClient) FindInvitations(email string) (invs []Invitation, err error) {
	// API converts all invited emails to lowercase.
	// https://github.com/rollbar/terraform-provider-rollbar/issues/139
	email = strings.ToLower(email)
	l := log.With().
		Str("email", email).
		Logger()

	l.Debug().Msg("Finding invitations")
	invs, err = c.ListAllInvitationsPerEmail(email)
	if err != nil && err != ErrNotFound {
		l.Err(err).
			Msg("error finding invitations")
		return invs, err
	}
	if len(invs) == 0 {
		return invs, ErrNotFound
	}
	l.Debug().
		Int("invitation_count", len(invs)).
		Msg("Successfully found invitations")
	return invs, nil
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
