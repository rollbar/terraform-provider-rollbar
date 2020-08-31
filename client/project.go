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

	u := c.url
	u.Path = path.Join(u.Path, PathProjectCreate)

	resp, err := c.resty.R().
		SetBody(map[string]interface{}{"name": name}).
		SetResult(ProjectResult{}).
		Post(u.String())
	if err != nil {
		l.Err(err).Msg("Error creating project")
		return nil, err
	}

	pr := resp.Result().(*ProjectResult)
	if pr.Err != 0 {
		l.Error().Msg("Unexpected error creating project")
	}

	return &pr.Result, nil
}

// ReadProject fetches data for the specified Project from the Rollbar API.
func (c *RollbarApiClient) ReadProject(id string) (*Project, error) {
	// NOTE: Since the project ID is ultimately an integer, it seems
	// appropriate that the argument to this function should be an int.
	// However the ID will be represented as a string in the
	// schema.ResourceData, and will be consumed as a string by this function
	// when constructing the URL for the API call.
	//
	// Should this client ever be extracted as a library, it would be
	// appropriate to make argument `id` an integer.  Until then, keeping it
	// as a string eliminates two needless type conversions.
	l := log.With().Str("id", id).Logger()
	l.Debug().Msg("Reading project from API")

	u := c.url
	u.Path = path.Join(u.Path, PathProjectRead)

	resp, err := c.resty.R().
		SetPathParams(map[string]string{
			"id": id,
		}).
		SetResult(ProjectResult{}).
		Get(u.String())
	if err != nil {
		l.Err(err).Msg("Error reading project")
		return nil, err
	}

	pr := resp.Result().(*ProjectResult)
	if pr.Err != 0 {
		l.Error().Msg("Unexpected error reading project")
	}

	return &pr.Result, nil
}
