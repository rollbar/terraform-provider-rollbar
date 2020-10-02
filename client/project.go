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

// ListProjects queries API for the list of projects
func (c *RollbarApiClient) ListProjects() ([]Project, error) {
	u := apiUrl + pathProjectList
	l := log.With().
		Str("url", u).
		Logger()

	resp, err := c.resty.R().
		SetResult(ProjectListResult{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		errResp := resp.Error().(*ErrorResult)
		l.Err(errResp).Send()
		return nil, errResp
	}
	lpr := resp.Result().(*ProjectListResult)

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

// CreateProject creates a new project
func (c *RollbarApiClient) CreateProject(name string) (*Project, error) {
	u := apiUrl + pathProjectCreate
	l := log.With().
		Str("name", name).
		Str("url", u).
		Logger()
	l.Debug().Msg("Creating new project")

	resp, err := c.resty.R().
		SetBody(map[string]interface{}{"name": name}).
		SetResult(ProjectResult{}).
		SetError(ErrorResult{}).
		Post(u)
	if err != nil {
		l.Err(err).Msg("Error creating project")
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return nil, er
	}

	pr := resp.Result().(*ProjectResult)
	if pr.Err != 0 {
		l.Error().Msg("Unexpected error creating project")
	}

	return &pr.Result, nil
}

// ReadProject fetches data for the specified Project from the Rollbar API.
func (c *RollbarApiClient) ReadProject(projectId int) (*Project, error) {
	u := apiUrl + pathProjectRead

	l := log.With().
		Int("projectId", projectId).
		Str("url", u).
		Logger()
	l.Debug().Msg("Reading project from API")

	resp, err := c.resty.R().
		SetResult(ProjectResult{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"projectId": strconv.Itoa(projectId),
		}).
		Get(u)
	if err != nil {
		l.Err(err).Msg("Error reading project")
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return nil, er
	}

	pr := resp.Result().(*ProjectResult)
	if pr.Err != 0 {
		l.Error().Msg("Unexpected error reading project")
	}

	l.Debug().Msg("Project successfully read")
	return &pr.Result, nil
}

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
	if resp.StatusCode() != http.StatusOK {
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating a project")
		return er
	}

	l.Debug().Msg("Project successfully deleted")
	return nil
}
