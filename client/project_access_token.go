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
	"strconv"
)

// ProjectAccessToken represents a Rollbar project access token.
type ProjectAccessToken struct {
	Name                 string  `mapstructure:"name"`
	ProjectID            int     `json:"project_id" mapstructure:"project_id"`
	AccessToken          string  `json:"access_token" mapstructure:"access_token"`
	Scopes               []Scope `mapstructure:"scopes"`
	Status               Status  `mapstructure:"status"`
	RateLimitWindowSize  int     `json:"rate_limit_window_size" mapstructure:"rate_limit_window_size"`
	RateLimitWindowCount int     `json:"rate_limit_window_count" mapstructure:"rate_limit_window_count"`
	DateCreated          int     `json:"date_created" mapstructure:"date_created"`
	DateModified         int     `json:"date_modified" mapstructure:"date_modified"`

	// FIXME: Should we eliminate these fields?
	//  https://github.com/rollbar/terraform-provider-rollbar/issues/52
	//CurRateLimitWindowCount int     `json:"cur_rate_limit_window_count" mapstructure:"cur_rate_limit_window_count"`
	//CurRateLimitWindowStart int     `json:"cur_rate_limit_window_start" mapstructure:"cur_rate_limit_window_start"`
}

// Scope represents the scope of a Rollbar project access token.
type Scope string

// Possible values for project access token scope
const (
	ScopeWrite          = Scope("write")
	ScopeRead           = Scope("read")
	ScopePostServerItem = Scope("post_server_item")
	ScopePostClientItem = Scope("post_client_item")
)

// ProjectAccessTokenUpdateArgs encapsulates arguments for creating a Rollbar
// project access token.
type ProjectAccessTokenCreateArgs struct {
	ProjectID            int     `json:"-"`
	Name                 string  `json:"name"`
	Scopes               []Scope `json:"scopes"`
	Status               Status  `json:"status"`
	RateLimitWindowSize  int     `json:"rate_limit_window_size"`
	RateLimitWindowCount int     `json:"rate_limit_window_count"`
}

// sanityCheck checks that the arguments are sane.
func (args *ProjectAccessTokenCreateArgs) sanityCheck() error {
	l := log.With().
		Interface("args", args).
		Logger()
	if args.ProjectID <= 0 {
		err := fmt.Errorf("project ID cannot be blank")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	if args.Name == "" {
		err := fmt.Errorf("name cannot be blank")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	if len(args.Scopes) < 1 {
		err := fmt.Errorf("at least one scope must be specified")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	for _, s := range args.Scopes {
		switch s {
		case ScopeRead, ScopeWrite, ScopePostClientItem, ScopePostServerItem:
			// Passed sanity check
		default:
			// FIXME: Default switch case needs test coverage.
			//  https://github.com/rollbar/terraform-provider-rollbar/issues/39
			err := fmt.Errorf("invalid scope")
			l.Err(err).Msg("Failed sanity check")
			return err
		}
	}
	switch args.Status {
	case StatusEnabled, StatusDisabled:
		// Passed sanity check
	default:
		// FIXME: Default switch case needs test coverage.
		//  https://github.com/rollbar/terraform-provider-rollbar/issues/39
		err := fmt.Errorf("invalid status")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	if args.RateLimitWindowCount < 0 {
		err := fmt.Errorf("rate limit window count must be zero or greater")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	if args.RateLimitWindowSize < 0 {
		err := fmt.Errorf("rate limit window size must be zero or greater")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	return nil // Sanity check passed
}

// ProjectAccessTokenUpdateArgs encapsulates the required and optional arguments
// for creating a Rollbar project access token.
//
// Currently not all attributes can be updated.
//  https://github.com/rollbar/terraform-provider-rollbar/issues/41
type ProjectAccessTokenUpdateArgs struct {
	ProjectID            int    `json:"-"`
	AccessToken          string `json:"-"`
	RateLimitWindowSize  int    `json:"rate_limit_window_size"`
	RateLimitWindowCount int    `json:"rate_limit_window_count"`
}

// sanityCheck checks that the arguments are sane.
func (args *ProjectAccessTokenUpdateArgs) sanityCheck() error {
	l := log.With().
		Interface("args", args).
		Logger()
	if args.ProjectID <= 0 {
		err := fmt.Errorf("project ID cannot be blank")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	if args.AccessToken == "" {
		err := fmt.Errorf("access token cannot be blank")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	if args.RateLimitWindowCount < 0 {
		err := fmt.Errorf("rate limit window count must be zero or greater")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	if args.RateLimitWindowSize < 0 {
		err := fmt.Errorf("rate limit window size must be zero or greater")
		l.Err(err).Msg("Failed sanity check")
		return err
	}
	return nil // Sanity check passed
}

// ListProjectAccessTokens lists the Rollbar project access tokens for the
// specified Rollbar project.
func (c *RollbarApiClient) ListProjectAccessTokens(projectID int) ([]ProjectAccessToken, error) {
	l := log.With().
		Int("projectID", projectID).
		Logger()
	l.Debug().Msg("Listing project access tokens")

	u := apiUrl + pathPatList
	resp, err := c.Resty.R().
		SetResult(patListResponse{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"projectId": strconv.Itoa(projectID),
		}).
		Get(u)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	pats := resp.Result().(*patListResponse).Result
	return pats, nil
}

// ReadProjectAccessToken reads a Rollbar project access token from the API.  It
// returns the first token that matches `name`. If no matching token is found,
// returns error ErrNotFound.
func (c *RollbarApiClient) ReadProjectAccessToken(projectID int, token string) (ProjectAccessToken, error) {
	l := log.With().
		Int("projectID", projectID).
		Str("token", token).
		Logger()
	l.Debug().Msg("Reading project access token")

	var pat ProjectAccessToken
	tokens, err := c.ListProjectAccessTokens(projectID)
	if err != nil {
		l.Err(err).
			Msg("Error listing project access tokens")
		return pat, err
	}

	for _, t := range tokens {
		if t.AccessToken == token {
			l.Debug().
				Interface("token", t).
				Msg("Found matching project access token")
			return t, nil
		}
	}

	l.Warn().Msg("Could not find matching project access token")
	return pat, ErrNotFound
}

// ReadProjectAccessTokenByName reads a Rollbar project access token from the
// API.  It returns the first token that matches `name`. If no matching token is
// found, returns error ErrNotFound.
func (c *RollbarApiClient) ReadProjectAccessTokenByName(projectID int, name string) (ProjectAccessToken, error) {
	l := log.With().
		Int("projectID", projectID).
		Str("name", name).
		Logger()
	l.Debug().Msg("Reading project access token")

	var pat ProjectAccessToken
	tokens, err := c.ListProjectAccessTokens(projectID)
	if err != nil {
		l.Err(err).
			Msg("Error reading project access token")
		return pat, err
	}

	for _, t := range tokens {
		l.Debug().Msg("Found project access token with matching name")
		if t.Name == name {
			return t, nil
		}
	}

	l.Warn().Msg("Could not find project access token with matching name")
	return pat, ErrNotFound
}

func (c *RollbarApiClient) DeleteProjectAccessToken(projectID int, token string) error {
	// FIXME: Implement this functionality when the API is ready!
	//  https://github.com/rollbar/terraform-provider-rollbar/issues/12
	log.Warn().Msg("Deleting project access tokens not yet implemented by Rollbar API.")
	return nil
}

// CreateProjectAccessToken creates a Rollbar project access token.
func (c *RollbarApiClient) CreateProjectAccessToken(args ProjectAccessTokenCreateArgs) (ProjectAccessToken, error) {
	l := log.With().
		Interface("args", args).
		Logger()
	l.Debug().Msg("Creating new project access token")
	var pat ProjectAccessToken

	err := args.sanityCheck()
	if err != nil {
		l.Err(err).Msg("Arguments to create project access token failed sanity check.")
		return pat, err
	}

	u := apiUrl + pathPatCreate
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"projectId": strconv.Itoa(args.ProjectID),
		}).
		SetBody(args).
		SetResult(patCreateResponse{}).
		SetError(ErrorResult{}).
		Post(u)
	if err != nil {
		l.Err(err).Msg("Error creating project access token")
		return pat, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return pat, err
	}
	r := resp.Result().(*patCreateResponse)
	pat = r.Result
	l.Debug().
		Interface("token", pat).
		Msg("Successfully created new project access token")
	return pat, nil
}

// UpdateProjectAccessToken updates a Rollbar project access token.
func (c *RollbarApiClient) UpdateProjectAccessToken(args ProjectAccessTokenUpdateArgs) error {
	l := log.With().
		Interface("args", args).
		Logger()
	l.Debug().Msg("Updating project access token")

	err := args.sanityCheck()
	if err != nil {
		l.Err(err).Msg("Arguments to update project access token failed sanity check.")
		return err
	}

	u := apiUrl + pathPatUpdate
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{
			"projectId":   strconv.Itoa(args.ProjectID),
			"accessToken": args.AccessToken,
		}).
		SetBody(args).
		SetResult(patUpdateResponse{}).
		SetError(ErrorResult{}).
		Patch(u)
	if err != nil {
		l.Err(err).Msg("Error updating project access token")
		return err
	}
	return errorFromResponse(resp)
}

/*
 * Containers for unmarshalling Rollbar API responses
 */

type patListResponse struct {
	Error  int `json:"err"`
	Result []ProjectAccessToken
}

type patCreateResponse struct {
	Error  int `json:"err"`
	Result ProjectAccessToken
}

type patUpdateResponse struct {
	Error int `json:"err"`
}
