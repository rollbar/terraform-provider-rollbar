/*
 * Copyright (c) 2023 Rollbar, Inc.
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
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// Team represents a Rollbar team.
type Team struct {
	ID          int
	AccountID   int `json:"account_id"`
	Name        string
	AccessLevel string `json:"access_level"`
}

// CreateTeam creates a new Rollbar team.
func (c *RollbarAPIClient) CreateTeam(name, level string) (Team, error) {
	c.m.Lock()
	defer c.m.Unlock()
	var t Team
	l := log.With().
		Str("name", name).
		Logger()
	l.Debug().
		Str("access_level", level).
		Msg("Creating new team")

	// Sanity check
	if name == "" {
		return t, fmt.Errorf("name cannot be blank")
	}

	u := c.BaseURL + pathTeamCreate
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
		Int("id", t.ID).
		Msg("Successfully created new team")
	return t, nil
}

// ListTeams lists all Rollbar teams.
func (c *RollbarAPIClient) ListTeams() ([]Team, error) {
	c.m.Lock()
	defer c.m.Unlock()
	log.Debug().Msg("Listing all teams")
	var teams []Team
	u := c.BaseURL + pathTeamList
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
	count := len(teams)
	log.Debug().Int("count", count).Msg("Successfully listed teams")
	return teams, nil
}

// ListCustomTeams lists all custom defined teams, excluding system teams
// "Everyone" and "Owners".
// FIXME: This function needs a better name.
func (c *RollbarAPIClient) ListCustomTeams() ([]Team, error) {
	log.Debug().Msg("Listing custom teams")
	var customTeams []Team
	allTeams, err := c.ListTeams()
	if err != nil {
		log.Err(err).Msg("Error listing custom teams")
		return customTeams, err
	}
	customTeams = filterSystemTeams(allTeams)
	count := len(customTeams)
	log.Debug().Int("count", count).Msg("Successfully listed custom teams")
	return customTeams, nil
}

// ReadTeam reads a Rollbar team from the API. If no matching team is found,
// returns error ErrNotFound.
func (c *RollbarAPIClient) ReadTeam(id int) (Team, error) {
	c.m.Lock()
	defer c.m.Unlock()
	var t Team
	l := log.With().
		Int("id", id).
		Logger()
	l.Debug().Msg("Reading team from API")

	// Sanity check
	if id == 0 {
		return t, fmt.Errorf("id must be non-zero")
	}

	u := c.BaseURL + pathTeamRead
	u = strings.ReplaceAll(u, "{teamID}", strconv.Itoa(id))
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
		l.Err(err).Msg("Error reading team")
		return t, err
	}
	r := resp.Result().(*teamReadResponse)
	t = r.Result
	l.Debug().
		Int("id", t.ID).
		Str("name", t.Name).
		Msg("Successfully read team")
	return t, nil
}

// DeleteTeam deletes a Rollbar team. If no matching team is found, returns
// error ErrNotFound.
func (c *RollbarAPIClient) DeleteTeam(id int) error {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().
		Int("id", id).
		Logger()
	l.Debug().Msg("Deleting team")

	// Sanity check
	if id == 0 {
		return fmt.Errorf("id must be non-zero")
	}

	u := c.BaseURL + pathTeamDelete
	u = strings.ReplaceAll(u, "{teamID}", strconv.Itoa(id))
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
func (c *RollbarAPIClient) AssignUserToTeam(teamID, userID int) error {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().Int("userID", userID).Int("teamID", teamID).Logger()
	l.Debug().Msg("Assigning user to team")
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamID": strconv.Itoa(teamID),
			"userID": strconv.Itoa(userID),
		}).
		SetError(ErrorResult{}).
		Put(c.BaseURL + pathTeamUser)
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

// IsUserAssignedToTeam checks if a user is assigned to a Rollbar team.
func (c *RollbarAPIClient) IsUserAssignedToTeam(teamID, userID int) (bool, error) {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().
		Int("userID", userID).
		Int("teamID", teamID).
		Logger()
	l.Debug().Msg("Checking if user is assigned to team")
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamID": strconv.Itoa(teamID),
			"userID": strconv.Itoa(userID),
		}).
		SetError(ErrorResult{}).
		Get(c.BaseURL + pathTeamUser)
	if err != nil {
		l.Err(err).Msg("Error checking if user is assigned to team")
		return false, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		if resp.StatusCode() == http.StatusNotFound {
			l.Err(err).Msg("User is not assigned to the team")
			return false, nil
		}
		l.Err(err).Msg("Unknown error")
		return false, err
	}
	l.Debug().Msg("User is assigned to the team")
	return true, nil
}

// RemoveUserFromTeam removes a user from a Rollbar team.
func (c *RollbarAPIClient) RemoveUserFromTeam(userID, teamID int) error {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().Int("userID", userID).Int("teamID", teamID).Logger()
	l.Debug().Msg("Removing user from team")
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamID": strconv.Itoa(teamID),
			"userID": strconv.Itoa(userID),
		}).
		SetError(ErrorResult{}).
		Delete(c.BaseURL + pathTeamUser)
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

// FindTeamID finds the ID for a team.
func (c *RollbarAPIClient) FindTeamID(name string) (int, error) {
	l := log.With().
		Str("team_name", name).
		Logger()
	l.Debug().Msg("Finding team ID")
	teams, err := c.ListTeams()
	if err != nil {
		l.Err(err).Send()
		return 0, err
	}
	for _, t := range teams {
		if t.Name == name {
			l.Debug().Int("team_id", t.ID).Msg("Found team ID")
			return t.ID, nil
		}
	}
	l.Debug().Msg("Could not find team ID")
	return 0, ErrNotFound
}

// AssignTeamToProject assigns a Rollbar team to a project.
func (c *RollbarAPIClient) AssignTeamToProject(teamID, projectID int) error {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().
		Int("teamID", teamID).
		Int("projectID", projectID).
		Logger()
	l.Debug().Msg("Assigning team to project")
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamID":    strconv.Itoa(teamID),
			"projectID": strconv.Itoa(projectID),
		}).
		SetError(ErrorResult{}).
		Put(c.BaseURL + pathTeamProject)
	if err != nil {
		l.Err(err).Msg("Error assigning team to project")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Msg("Error assigning team to project")
		return err
	}
	l.Debug().Msg("Successfully assigned team to project")
	return nil
}

// RemoveTeamFromProject removes a Rollbar team from a project.
func (c *RollbarAPIClient) RemoveTeamFromProject(teamID, projectID int) error {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().
		Int("teamID", teamID).
		Int("projectID", projectID).
		Logger()
	l.Debug().Msg("Removing team from project")
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"teamID":    strconv.Itoa(teamID),
			"projectID": strconv.Itoa(projectID),
		}).
		SetError(ErrorResult{}).
		Delete(c.BaseURL + pathTeamProject)
	if err != nil {
		l.Err(err).Msg("Error removing team from project")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Msg("Error removing team from project")
		return err
	}
	l.Debug().Msg("Successfully removed team from project")
	return nil
}

/*
 * Convenience functions
 */

// filterSystemTeams filters out the system teams "Everyone" and "Owners" from a
// list of Rollbar teams.
func filterSystemTeams(teams []Team) []Team {
	customTeams := []Team{}
	for _, t := range teams {
		if t.Name == "Everyone" || t.Name == "Owners" {
			continue
		}
		customTeams = append(customTeams, t)
	}
	return customTeams
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

type teamProjectListResponse struct {
	Err    int
	Result []struct {
		ProjectID int `json:"project_id"`
		TeamID    int `json:"team_id"`
	}
}
