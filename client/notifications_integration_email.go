package client

import "github.com/rs/zerolog/log"

// UpdateNotificationsEmailIntegration updates settings for Rollbar notifications
// email integration.
func (c *RollbarApiClient) UpdateNotificationsEmailIntegration(enabled, includeRequestParams bool) error {
	l := log.With().
		Bool("enabled", enabled).
		Bool("include_request_params", includeRequestParams).
		Logger()
	l.Debug().Msg("Updating email notification integration")

	u := apiUrl + pathNotificationIntegrationEmail
	resp, err := c.Resty.R().
		SetBody(map[string]interface{}{
			"enabled":                enabled,
			"include_request_params": includeRequestParams,
		}).
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
