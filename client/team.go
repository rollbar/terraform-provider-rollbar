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
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
)

// Team represents a Rollbar team.
type Team struct {
	ID          int
	AccountID   int `json:"account_id"`
	Name        string
	AccessLevel TeamAccessLevel `json:"access_level"`
}

// TeamAccessLevel represents the Rollbar team access level.
type TeamAccessLevel string

// Possible values for team access level
const (
	TeamAccessStandard = TeamAccessLevel("standard")
	TeamAccessLight    = TeamAccessLevel("light")
	TeamAccessView     = TeamAccessLevel("view")
)

// CreateTeam creates a new Rollbar team.
func (c *RollbarApiClient) CreateTeam(name string, level TeamAccessLevel) (Team, error) {
	var t Team
	l := log.With().
		Str("name", name).
		Str("access_level", string(level)).
		Logger()
	l.Debug().Msg("Creating new team")

	// Sanity check
	if name == "" {
		return t, fmt.Errorf("name cannot be blank")
	}

	u := apiUrl + pathTeamCreate
	resp, err := c.resty.R().
		SetBody(map[string]interface{}{"name": name}).
		SetResult(teamCreateResponse{}).
		SetError(ErrorResult{}).
		Post(u)
	if err != nil {
		l.Err(err).Msg("Error creating team")
		return t, err
	}
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusCreated:
		// FIXME: currently API returns `200 OK` on successful create; but it
		//  should instead return `201 Created`.
		//  https://github.com/rollbar/terraform-provider-rollbar/issues/8
		r := resp.Result().(*teamCreateResponse)
		t = r.Result
		l.Debug().
			Interface("team", t).
			Msg("Successfully created new team")
		return t, nil
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return t, ErrUnauthorized
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating team")
		return t, er
	}
}

// ListTeams lists all Rollbar teams.
func (c *RollbarApiClient) ListTeams() ([]Team, error) {
	log.Debug().Msg("Listing all teams")
	var teams []Team
	u := apiUrl + pathTeamList
	resp, err := c.resty.R().
		SetResult(teamListResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		log.Err(err).Msg("Error listing teams")
		return teams, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		r := resp.Result().(*teamListResponse)
		teams = r.Result
		log.Debug().
			Interface("teams", teams).
			Msg("Successfully listed teams")
		return teams, nil
	case http.StatusUnauthorized:
		log.Warn().Msg("Unauthorized")
		return teams, ErrUnauthorized
	default:
		er := resp.Error().(*ErrorResult)
		log.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error listing teams")
		return teams, er
	}
}

// ReadTeam reads a Rollbar team from the API. If no matching team is found,
// returns error ErrNotFound.
func (c *RollbarApiClient) ReadTeam(id int) (Team, error) {
	var t Team
	l := log.With().
		Int("id", id).
		Logger()
	l.Debug().Msg("Reading team")

	// Sanity check
	if id == 0 {
		return t, fmt.Errorf("id must be non-zero")
	}

	u := apiUrl + pathTeamRead
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(id))
	resp, err := c.resty.R().
		SetResult(teamReadResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		l.Err(err).Msg("Error creating team")
		return t, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		r := resp.Result().(*teamReadResponse)
		t = r.Result
		l.Debug().
			Interface("team", t).
			Msg("Successfully read team")
		return t, nil
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return t, ErrUnauthorized
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error reading team")
		return t, er
	}

}

// DeleteTeam deletes a Rollbar team. If no matching team is found, returns
// error ErrNotFound.
func (c *RollbarApiClient) DeleteTeam(id int) error {
	l := log.With().
		Int("id", id).
		Logger()
	l.Debug().Msg("Deleting team")

	// Sanity check
	if id == 0 {
		return fmt.Errorf("id must be non-zero")
	}

	u := apiUrl + pathTeamDelete
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(id))
	resp, err := c.resty.R().
		SetError(ErrorResult{}).
		Delete(u)
	if err != nil {
		l.Err(err).Msg("Error deleting team")
		return err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		l.Debug().Msg("Successfully deleted team")
		return nil
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return ErrUnauthorized
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error deleting team")
		return er
	}

}

/*
 * Containers for unmarshalling API responses
 */

type teamCreateResponse struct {
	Err    int
	Result Team
}

type teamListResponse struct {
	Err    int
	Result []Team
}

type teamReadResponse struct {
	Err    int
	Result Team
}
