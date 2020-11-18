package client

import "github.com/rs/zerolog/log"

// EmailIntegrationArgs encapsulates the arguments required to update a Rollbar
// notifications email integration.
type EmailIntegrationArgs struct {
	Token                string `json:"-"`
	Enabled              bool   `json:"enabled"`
	IncludeRequestParams bool   `json:"include_request_params"`
}

// UpdateNotificationsEmailIntegration updates settings for Rollbar notifications
// email integration.
func (c *RollbarApiClient) UpdateNotificationsEmailIntegration(args EmailIntegrationArgs) error {
	l := log.With().
		Interface("args", args).
		Logger()
	l.Debug().Msg("Updating email notification integration")

	u := apiUrl + pathNotificationIntegrationEmail
	resp, err := c.Resty.R().
		// Project access token is required, instead of the account access token
		// Resty is configured with by default.
		SetHeader("X-Rollbar-Access-Token", args.Token).
		SetBody(args).
		SetError(ErrorResult{}).
		Put(u)
	if err != nil {
		l.Err(err).Msg("Error updating email notification integration")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Msg("Error updating email notification integration")
		return err
	}
	return nil
}
