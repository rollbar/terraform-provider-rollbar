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
	AccessLevel string `json:"access_level"`
}

// CreateTeam creates a new Rollbar team.
func (c *RollbarApiClient) CreateTeam(name string, level string) (Team, error) {
	var t Team
	l := log.With().
		Str("name", name).
		Str("access_level", level).
		Logger()
	l.Info().Msg("Creating new team")

	// Sanity check
	if name == "" {
		return t, fmt.Errorf("name cannot be blank")
	}

	u := apiUrl + pathTeamCreate
	resp, err := c.Resty.R().
		SetBody(map[string]interface{}{
			"name":         name,
			"access_level": level,
		}).
		SetResult(teamCreateResponse{}).
		SetError(ErrorResult{}).
		Post(u)
	if err != nil {
		l.Err(err).Msg("Error creating team")
		return t, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Msg("Error creating team")
		return t, err
	}
	r := resp.Result().(*teamCreateResponse)
	t = r.Result
	l.Debug().
		Interface("team", t).
		Msg("Successfully created new team")
	return t, nil
}

// ListTeams lists all Rollbar teams.
func (c *RollbarApiClient) ListTeams() ([]Team, error) {
	log.Info().Msg("Listing all teams")
	var teams []Team
	u := apiUrl + pathTeamList
	resp, err := c.Resty.R().
		SetResult(teamListResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		log.Err(err).Msg("Error listing teams")
		return teams, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		log.Err(err).Msg("Error listing teams")
		return teams, err
	}
	r := resp.Result().(*teamListResponse)
	teams = r.Result
	log.Debug().Msg("Successfully listed teams")
	return teams, nil
}

// ReadTeam reads a Rollbar team from the API. If no matching team is found,
// returns error ErrNotFound.
func (c *RollbarApiClient) ReadTeam(id int) (Team, error) {
	var t Team
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading Rollbar team from API")

	// Sanity check
	if id == 0 {
		return t, fmt.Errorf("id must be non-zero")
	}

	u := apiUrl + pathTeamRead
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(id))
	resp, err := c.Resty.R().
		SetResult(teamReadResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		l.Err(err).Msg("Error reading team")
		return t, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		// FIXME: Workaround API bug
		//  https://github.com/rollbar/terraform-provider-rollbar/issues/79
		statusForbidden := resp.StatusCode() == http.StatusForbidden
		msgNotFound := strings.Contains(err.Error(), "Team not found in this account")
		if statusForbidden && msgNotFound {
			return t, ErrNotFound
		}
		l.Err(err).Msg("Error reading team")
		return t, err
	}
	r := resp.Result().(*teamReadResponse)
	t = r.Result
	l.Debug().
		Interface("team", t).
		Msg("Successfully read team")
	return t, nil
}

// DeleteTeam deletes a Rollbar team. If no matching team is found, returns
// error ErrNotFound.
func (c *RollbarApiClient) DeleteTeam(id int) error {
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Deleting team")

	// Sanity check
	if id == 0 {
		return fmt.Errorf("id must be non-zero")
	}

	u := apiUrl + pathTeamDelete
	u = strings.ReplaceAll(u, "{teamId}", strconv.Itoa(id))
	resp, err := c.Resty.R().
		SetError(ErrorResult{}).
		Delete(u)
	if err != nil {
		l.Err(err).Msg("Error deleting team")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Msg("Error deleting team")
		return err
	}
	l.Debug().Msg("Successfully deleted team")
	return nil
}

// AssignUserToTeam assigns a user to a Rollbar team.
func (c *RollbarApiClient) AssignUserToTeam(teamID, userID int) error {
	l := log.With().Int("userID", userID).Int("teamID", teamID).Logger()
	l.Info().Msg("Assigning user to team")
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamId": strconv.Itoa(teamID),
			"userId": strconv.Itoa(userID),
		}).
		SetError(ErrorResult{}).
		Put(apiUrl + pathTeamUser)
	if err != nil {
		l.Err(err).Msg("Error assigning user to team")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		// API returns status `403 Forbidden` on invalid user to team assignment
		// https://github.com/rollbar/terraform-provider-rollbar/issues/66
		if resp.StatusCode() == http.StatusForbidden {
			l.Err(err).Msg("Team or user not found")
			return ErrNotFound
		}
		l.Err(err).Msg("Error assigning user to team")
		return err
	}
	l.Debug().Msg("Successfully assigned user to team")
	return nil
}

// RemoveUserFromTeam removes a user from a Rollbar team.
func (c *RollbarApiClient) RemoveUserFromTeam(teamID, userID int) error {
	l := log.With().Int("userID", userID).Int("teamID", teamID).Logger()
	l.Info().Msg("Removing user from team")
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamId": strconv.Itoa(teamID),
			"userId": strconv.Itoa(userID),
		}).
		SetError(ErrorResult{}).
		Delete(apiUrl + pathTeamUser)
	if err != nil {
		l.Err(err).Msg("Error removing user from team")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		// API returns status `422 Unprocessable Entity` on invalid user to team
		// assignment.
		// https://github.com/rollbar/terraform-provider-rollbar/issues/66
		if resp.StatusCode() == http.StatusUnprocessableEntity {
			l.Err(err).Msg("Team or user not found")
			return ErrNotFound
		}
		l.Err(err).Msg("Error removing user from team")
		return err
	}
	l.Debug().Msg("Successfully removed user from team")
	return nil

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
