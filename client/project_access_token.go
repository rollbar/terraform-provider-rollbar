package client

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// ProjectAccessToken represents a Rollbar project access token.
type ProjectAccessToken struct {
	Name                    string  `mapstructure:"name"`
	ProjectID               int     `json:"project_id" mapstructure:"project_id"`
	AccessToken             string  `json:"access_token" mapstructure:"access_token"`
	Scopes                  []Scope `mapstructure:"scopes"`
	Status                  Status  `mapstructure:"status"`
	RateLimitWindowSize     *int    `json:"rate_limit_window_size" mapstructure:"rate_limit_window_size"`
	RateLimitWindowCount    *int    `json:"rate_limit_window_count" mapstructure:"rate_limit_window_count"`
	CurRateLimitWindowCount *int    `json:"cur_rate_limit_window_count" mapstructure:"cur_rate_limit_window_count"`
	CurRateLimitWindowStart *int    `json:"cur_rate_limit_window_start" mapstructure:"cur_rate_limit_window_start"`
	DateCreated             int     `json:"date_created" mapstructure:"date_created"`
	DateModified            int     `json:"date_modified" mapstructure:"date_modified"`
}

// Scope represents the scope of a Rollbar project access token.
type Scope string

// Possible values forproject access token scope
const (
	ScopeWrite          = Scope("write")
	ScopeRead           = Scope("read")
	ScopePostServerItem = Scope("post_server_item")
	ScopePostClientItem = Scope("post_client_item")
)

// ProjectAccessTokenArgs encapsulates the required and optional arguments for
// creating and updating Rollbar project access tokens.
type ProjectAccessTokenArgs struct {
	// Required
	ProjectID            int     `json:"-"`
	Name                 string  `json:"name"`
	Scopes               []Scope `json:"scopes"`
	Status               Status  `json:"status"`
	RateLimitWindowSize  uint    `json:"rate_limit_window_size"`
	RateLimitWindowCount uint    `json:"rate_limit_window_count"`
}

// sanityCheck checks that the arguments are sane.
func (args *ProjectAccessTokenArgs) sanityCheck() error {
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
	return nil
}

// ListProjectAccessTokens lists the Rollbar project access tokens for the
// specified Rollbar project.
func (c *RollbarApiClient) ListProjectAccessTokens(projectID int) ([]ProjectAccessToken, error) {
	u := apiUrl + pathPatList

	l := log.With().
		Str("url", u).
		Logger()

	resp, err := c.resty.R().
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
	switch resp.StatusCode() {
	case http.StatusOK:
		pats := resp.Result().(*patListResponse).Result
		return pats, nil
	case http.StatusNotFound:
		l.Warn().Msg("Project not found")
		return nil, ErrNotFound
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return nil, ErrUnauthorized
	default:
		errResp := resp.Error().(*ErrorResult)
		l.Err(errResp).Msg("Unexpected error")
		return nil, errResp
	}
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
		l.Debug().Msg("Found matching project access token")
		if t.AccessToken == token {
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
	//return fmt.Errorf("delete PAT not yet implemented by Rollbar API")
	log.Warn().Msg("Deleting project access tokens not yet implemented by Rollbar API.")
	return nil
}

// CreateProjectAccessToken creates a Rollbar project access token.
func (c *RollbarApiClient) CreateProjectAccessToken(args ProjectAccessTokenArgs) (ProjectAccessToken, error) {
	l := log.With().
		Interface("args", args).
		Logger()
	var pat ProjectAccessToken

	err := args.sanityCheck()
	if err != nil {
		l.Err(err).Msg("Arguments to create project access token failed sanity check.")
		return pat, err
	}

	u := apiUrl + pathPatCreate
	resp, err := c.resty.R().
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
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusCreated:
		// FIXME: currently API returns `200 OK` on successful create; but it
		//  should instead return `201 Created`.
		//  https://github.com/rollbar/terraform-provider-rollbar/issues/8
		r := resp.Result().(*patCreateResponse)
		pat = r.Result
		l.Debug().
			Interface("token", pat).
			Msg("Successfully created new project access token")
		return pat, nil
	case http.StatusUnauthorized:
		l.Warn().Msg("Unauthorized")
		return pat, ErrUnauthorized
	default:
		er := resp.Error().(*ErrorResult)
		l.Error().
			Int("StatusCode", resp.StatusCode()).
			Str("Status", resp.Status()).
			Interface("ErrorResult", er).
			Msg("Error creating project access token")
		return pat, er
	}
}

/*
 * Containers for unmarshalling Rollbar API responses
 */

type patListResponse struct {
	Error  int
	Result []ProjectAccessToken
}

type patCreateResponse struct {
	Error  int
	Result ProjectAccessToken
}
