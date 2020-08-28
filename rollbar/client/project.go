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
	"path"
)

// ListProjects queries API for the list of projects
func (c *RollbarApiClient) ListProjects() ([]Project, error) {
	url := c.url
	url.Path = path.Join(url.Path, PathProjectList)

	resp, err := c.resty.R().
		SetResult(ListProjectsResult{}).
		Get(url.String())
	if err != nil {
		return nil, err
	}

	lpr := resp.Result().(*ListProjectsResult)
	if lpr.Err != 0 {
		log.Error().
			Int("err", lpr.Err).
			Msg("Unexpected error listing projects")
		return nil, err
	}

	return lpr.Result, nil
}

// CreateProject creates a new project
func (c *RollbarApiClient) CreateProject(name string) (*Project, error) {
	l := log.With().Str("name", name).Logger()

	u := c.url
	u.Path = path.Join(u.Path, PathProjectCreate)

	resp, err := c.resty.R().
		SetBody(map[string]interface{}{"name": name}).
		Post(u.String())
	if err != nil {
		l.Err(err).Msg("Error creating project")
		return nil, err
	}

	cpr := resp.Result().(*CreateProjectResult)
	if cpr.Err != 0 {
		l.Error().Msg("Unexpected error creating project")
	}

	return &cpr.Result, nil
}
