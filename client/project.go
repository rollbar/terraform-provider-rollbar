/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package client

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// Project represents a Rollbar project.
type Project struct {
	Id           int    `model:"id" fake:"{number:1,1000000}"`
	Name         string `model:"name" fake:"{hackernoun}"`
	AccountId    int    `json:"account_id" model:"account_id" fake:"{number:1,1000000}"`
	DateCreated  int    `json:"date_created" model:"date_created" fake:"{number:1,1000000}"`
	DateModified int    `json:"date_modified" model:"date_modified" fake:"{number:1,1000000}"`
	Status       string `model:"status" fake:"{randomstring:[enabled,disabled]}"`
}

// FIXME: finish implementing the entire set of Project fields
/*
	SettingsData struct {
		Grouping struct {
			AutoUpgrade    bool     `json:"auto_upgrade"`
			RecentVersions []string `json:"recent_versions"`
		} `json:"grouping"`
		Integrations struct {
			Asana       interface{} `json:"asana"`
			AzureDevops interface{} `json:"azuredevops"`
			Bitbucket   interface{} `json:"bitbucket"`
				//"campfire": {},
				//"ciscospark": {},
				//"clubhouse": {},
				//"datadog": {},
				//"email": {
				//	"enabled": true
				//},
				//"flowdock": {},
				//"github": {},
				//"gitlab": {},
				//"hipchat": {},
				//"jira": {},
				//"lightstep": {},
				//"pagerduty": {},
				//"pivotal": {},
				//"slack": {},
				//"sprintly": {},
				//"trello": {},
				//"victorops": {},
				//"webhook": {}
		} `json:"integrations"`
		TimeFormat string `json:"time_format"`
		Timezone   string `json:"timezone"`
	} `json:"settings_data"`
*/

// ListProjects lists all Rollbar projects.
func (c *RollbarApiClient) ListProjects() ([]Project, error) {
	u := apiUrl + pathProjectList
	l := log.With().
		Str("url", u).
		Logger()

	resp, err := c.resty.R().
		SetResult(projectListResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return nil, ErrUnauthorized
	case http.StatusOK:
		l.Debug().Msg("Successfully listed projects")
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return nil, er
	}

	lpr := resp.Result().(*projectListResponse)
	// FIXME: After deleting a project through the API, it still shows up in
	//  the list of projects returned by the API - only with its name set to
	//  nil. This seemingly undesirable behavior should be fixed on the API
	//  side. We work around it by removing any result with an empty name.
	cleaned := make([]Project, 0)
	for _, proj := range lpr.Result {
		if proj.Name != "" {
			cleaned = append(cleaned, proj)
		}
	}

	return cleaned, nil
}

// CreateProject creates a new Rollbar project.
func (c *RollbarApiClient) CreateProject(name string) (*Project, error) {
	u := apiUrl + pathProjectCreate
	l := log.With().
		Str("name", name).
		Str("url", u).
		Logger()
	l.Debug().Msg("Creating new project")

	resp, err := c.resty.R().
		SetBody(map[string]interface{}{"name": name}).
		SetResult(projectResponse{}).
		SetError(ErrorResult{}).
		Post(u)
	if err != nil {
		l.Err(err).Msg("Error creating project")
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return nil, ErrUnauthorized
	case http.StatusOK:
		l.Debug().Msg("Project successfully created")
		pr := resp.Result().(*projectResponse)
		return &pr.Result, nil
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return nil, er
	}
}

// ReadProject a Rollbar project from the API. If no matching project is found,
// returns error ErrNotFound.
func (c *RollbarApiClient) ReadProject(projectId int) (*Project, error) {
	u := apiUrl + pathProjectRead

	l := log.With().
		Int("projectId", projectId).
		Str("url", u).
		Logger()
	l.Debug().Msg("Reading project from API")

	resp, err := c.resty.R().
		SetResult(projectResponse{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"projectId": strconv.Itoa(projectId),
		}).
		Get(u)
	if err != nil {
		l.Err(err).Msg("Error reading project")
		return nil, err
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		pr := resp.Result().(*projectResponse)
		l.Debug().Msg("Project successfully read")
		return &pr.Result, nil
	case http.StatusNotFound:
		l.Warn().Msg("Project not found")
		return nil, ErrNotFound
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return nil, ErrUnauthorized
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return nil, er
	}
}

// DeleteProject deletes a Rollbar project. If no matching project is found,
// returns error ErrNotFound.
func (c *RollbarApiClient) DeleteProject(projectId int) error {
	u := apiUrl + pathProjectDelete
	l := log.With().
		Int("projectId", projectId).
		Str("url", u).
		Logger()
	l.Debug().Msg("Deleting project")

	resp, err := c.resty.R().
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"projectId": strconv.Itoa(projectId),
		}).
		Delete(u)
	if err != nil {
		l.Err(err).Msg("Error deleting project")
		return err
	}
	l.Debug().Bytes("body", resp.Body()).Msg("Response body")
	switch resp.StatusCode() {
	case http.StatusNotFound:
		l.Warn().Msg("Project not found")
		return ErrNotFound
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return ErrUnauthorized
	case http.StatusOK:
		l.Debug().Msg("Project successfully deleted")
		return nil
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return er
	}
}

/*
 * Containers for unmarshalling API responses
 */

type projectListResponse struct {
	Err    int       `json:"err"`
	Result []Project `json:"result"`
}

type projectResponse struct {
	Err    int     `json:"err"`
	Result Project `json:"result"`
}
