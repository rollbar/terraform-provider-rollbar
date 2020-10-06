package client

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
)

// Team represents a Rollbar team
type Team struct {
	ID          int             `json:"id"`
	AccountID   int             `json:"account_id"`
	Name        string          `json:"name"`
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
		SetResult(teamCreateResult{}).
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
		r := resp.Result().(*teamCreateResult)
		t = r.Result
		l.Debug().
			Interface("team", t).
			Msg("Successfully created new team")
		return t, nil
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

// ListTeams lists all Rollbar teams
func (c *RollbarApiClient) ListTeams() ([]Team, error) {
	var teams []Team
	u := apiUrl + pathTeamList
	resp, err := c.resty.R().
		SetResult(teamListResult{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		log.Err(err).Msg("Error listing teams")
		return teams, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		r := resp.Result().(*teamListResult)
		teams = r.Result
		log.Debug().
			Interface("teams", teams).
			Msg("Successfully listed teams")
		return teams, nil
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

// ReadTeam reads a Rollbar team from the API.
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
		SetResult(teamReadResult{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		l.Err(err).Msg("Error creating team")
		return t, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		r := resp.Result().(*teamReadResult)
		t = r.Result
		l.Debug().
			Interface("team", t).
			Msg("Successfully read team")
		return t, nil
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

/*
 * Containers for unmarshalling API results
 */

type teamCreateResult struct {
	Err    int
	Result Team
}

type teamListResult struct {
	Err    int
	Result []Team
}

type teamReadResult struct {
	Err    int
	Result Team
}
