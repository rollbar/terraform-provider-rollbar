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
	"path"
	"strconv"
)

// ListProjects queries API for the list of projects
func (c *RollbarApiClient) ListProjects() ([]Project, error) {
	url := c.url
	url.Path = path.Join(url.Path, PathProjectList)

	resp, err := c.resty.R().
		SetResult(ProjectListResult{}).
		Get(url.String())
	if err != nil {
		return nil, err
	}

	lpr := resp.Result().(*ProjectListResult)
	if lpr.Err != 0 {
		log.Error().
			Int("err", lpr.Err).
			Msg("Unexpected error listing projects")
		return nil, err
	}

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
	l := log.With().Str("name", name).Logger()
	l.Debug().Msg("Creating new project")
	//var u url.URL

	u := *c.url
	u.Path = path.Join(u.Path, PathProjectCreate)
	l = l.With().Str("path", u.Path).Logger()

	resp, err := c.resty.R().
		SetBody(map[string]interface{}{"name": name}).
		SetResult(ProjectResult{}).
		SetError(ErrorResult{}).
		Post(u.String())
	l.Debug().Bytes("body", resp.Body()).Msg("Response Body")
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
func (c *RollbarApiClient) ReadProject(id int) (*Project, error) {
	l := log.With().Int("id", id).Logger()
	l.Debug().Msg("Reading project from API")

	u := *c.url
	u.Path = path.Join(u.Path, PathProjectRead)
	l = l.With().Str("path", u.Path).Logger()

	c.resty.SetDebug(true)
	rzl := RestyZeroLogger{l}
	c.resty.SetLogger(rzl)

	resp, err := c.resty.R().
		SetResult(ProjectResult{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"project_id": strconv.Itoa(id),
		}).
		Get(u.String())
	if err != nil {
		l.Err(err).Msg("Error reading project")
		return nil, err
	}
	l.Debug().Bytes("body", resp.Body()).Msg("Response Body")
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
	l.Debug().Interface("ProjectResult", pr).Send()
	if pr.Err != 0 {
		l.Error().Msg("Unexpected error reading project")
	}

	return &pr.Result, nil
}
