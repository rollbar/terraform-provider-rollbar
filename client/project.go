package client

import (
	"github.com/rs/zerolog/log"
	"path"
)

// ListProjects queries API for the list of projects
func (c *RollbarApiClient) ListProjects() ([]Project, error) {
	url := c.url
	url.Path = path.Join(url.Path, PathListProjects)

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
	}

	return lpr.Result, nil
}
