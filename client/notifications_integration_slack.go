package client

import "github.com/rs/zerolog/log"

type SlackIntegrationArgs struct {
	Token              string `json:"-"`
	Enabled            bool   `json:"enabled"`
	ServiceAccountID   int    `json:"service_account_id"`
	Channel            string `json:"channel"`
	ShowMessageButtons bool   `json:"show_message_buttons,omitempty"`
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
		SetHeader("X-Rollbar-Access-Token", args.Token).
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
