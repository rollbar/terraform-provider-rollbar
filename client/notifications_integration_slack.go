package client

import "github.com/rs/zerolog/log"

type SlackIntegrationArgs struct {
	Enabled            bool
	ServiceAccountID   int `json:"service_account_id"`
	Channel            string
	ShowMessageButtons bool `json:"show_message_buttons"`
}

// UpdateNotificationsSlackIntegration updates settings for Rollbar notifications
// email integration.
func (c *RollbarApiClient) UpdateNotificationsSlackIntegration(args SlackIntegrationArgs) error {
	l := log.With().
		Interface("args", args).
		Logger()
	l.Debug().Msg("Updating notifications Slack integration")

	u := apiUrl + pathNotificationIntegrationSlack
	resp, err := c.Resty.R().
		SetBody(args).
		SetError(ErrorResult{}).
		Put(u)
	if err != nil {
		l.Err(err).Msg("Error updating notifications Slack integration")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Msg("Error updating notifications Slack integration")
		return err
	}
	return nil
}
