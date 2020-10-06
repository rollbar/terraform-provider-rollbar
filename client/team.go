package client

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Team struct {
	ID          int             `json:"id"`
	AccountID   int             `json:"account_id"`
	Name        string          `json:"name"`
	AccessLevel TeamAccessLevel `json:"access_level"`
}

type TeamAccessLevel string

const (
	TeamAccessStandard = TeamAccessLevel("standard")
	TeamAccessLight    = TeamAccessLevel("light")
	TeamAccessView     = TeamAccessLevel("view")
)

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
		l.Debug().
			Interface("team", t).
			Msg("Successfully created new team")
		tcr := resp.Result().(*teamCreateResult)
		t = tcr.Result
		return t, nil
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return t, er
	}
}

type teamCreateResult struct {
	Err    int  `json:"err"`
	Result Team `json:"result"`
}
