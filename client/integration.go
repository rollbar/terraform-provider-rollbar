/*
 * Copyright (c) 2022 Rollbar, Inc.
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
	"github.com/rs/zerolog/log"
)

const (
	SLACK     string = "slack"
	WEBHOOK   string = "webhook"
	EMAIL     string = "email"
	PAGERDUTY string = "pagerduty"
)

var Integrations = map[string]interface{}{EMAIL: emailIntegrationResponse{}, PAGERDUTY: pagerDutyIntegrationResponse{},
	SLACK: slackIntegrationResponse{}, WEBHOOK: webhookIntegrationResponse{}}

type EmailIntegration struct {
	ProjectID int64 `model:"project_id" mapstructure:"project_id" json:"project_id"`
	Settings  struct {
		Enabled     bool `model:"enabled" mapstructure:"enabled" json:"enabled"`
		ScrubParams bool `model:"scrub_params" mapstructure:"scrub_params" json:"scrub_params"`
	} `model:"settings" mapstructure:"settings"`
}

type PagerDutyIntegration struct {
	ProjectID int64 `model:"project_id" mapstructure:"project_id" json:"project_id"`
	Settings  struct {
		Enabled    bool   `model:"enabled" mapstructure:"enabled" json:"enabled"`
		ServiceKey string `model:"service_key" mapstructure:"service_key" json:"service_key"`
	} `model:"settings" mapstructure:"settings"`
}

type SlackIntegration struct {
	ProjectID int64 `model:"project_id" mapstructure:"project_id" json:"project_id"`
	Settings  struct {
		Channel            string `model:"channel" mapstructure:"channel" json:"channel"`
		Enabled            bool   `model:"enabled" mapstructure:"enabled" json:"enabled"`
		ShowMessageButtons bool   `model:"show_message_buttons" mapstructure:"show_message_buttons" json:"show_message_buttons"`
		ServiceAccountID   string `model:"service_account_id" mapstructure:"service_account_id" json:"service_account_id"`
	} `model:"settings" mapstructure:"settings"`
}

type WebhookIntegration struct {
	ProjectID int64 `model:"project_id" mapstructure:"project_id" json:"project_id"`
	Settings  struct {
		Enabled bool   `model:"enabled" mapstructure:"enabled" json:"enabled"`
		URL     string `model:"url" mapstructure:"url" json:"url"`
	} `model:"settings" mapstructure:"settings"`
}

// UpdateIntegration updates a new Rollbar integration.
func (c *RollbarAPIClient) UpdateIntegration(integration string, bodyMap map[string]interface{}) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathIntegration
	l := log.With().
		Str("integration", integration).
		Logger()
	l.Debug().Msg("Update integration")
	resp, err := c.Resty.R().
		SetBody(bodyMap).
		SetResult(Integrations[integration]).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"integration": integration,
		}).
		Put(u)

	if err != nil {
		l.Err(err).Msg("Error updating integration")
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	l.Debug().Msg("integration successfully updated")
	switch integration {
	case EMAIL:
		return &(resp.Result().(*emailIntegrationResponse)).Result, nil
	case PAGERDUTY:
		return &(resp.Result().(*pagerDutyIntegrationResponse)).Result, nil
	case SLACK:
		return &(resp.Result().(*slackIntegrationResponse)).Result, nil
	case WEBHOOK:
		return &(resp.Result().(*webhookIntegrationResponse)).Result, nil
	}
	return nil, nil
}

// ReadIntegration reads a Rollbar integration from the API. If no matching integration is found,
// returns error ErrNotFound.
func (c *RollbarAPIClient) ReadIntegration(integration string) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathIntegration

	l := log.With().
		Str("integration", integration).
		Logger()
	l.Debug().Msg("Reading Integration from API")

	resp, err := c.Resty.R().
		SetResult(Integrations[integration]).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"integration": integration,
		}).
		Get(u)

	if err != nil {
		l.Err(err).Msg(resp.Status())
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	var errInt int
	switch integration {
	case EMAIL:
		i := resp.Result().(*emailIntegrationResponse)
		errInt = i.Err
	case PAGERDUTY:
		i := resp.Result().(*pagerDutyIntegrationResponse)
		errInt = i.Err
	case SLACK:
		i := resp.Result().(*slackIntegrationResponse)
		errInt = i.Err
	case WEBHOOK:
		i := resp.Result().(*webhookIntegrationResponse)
		errInt = i.Err
	}
	if errInt != 0 {
		l.Warn().Msg("Integration not found")
		return nil, ErrNotFound
	}
	l.Debug().Msg("Integration successfully read")
	switch integration {
	case EMAIL:
		return &(resp.Result().(*emailIntegrationResponse)).Result, nil
	case PAGERDUTY:
		return &(resp.Result().(*pagerDutyIntegrationResponse)).Result, nil
	case SLACK:
		return &(resp.Result().(*slackIntegrationResponse)).Result, nil
	case WEBHOOK:
		return &(resp.Result().(*webhookIntegrationResponse)).Result, nil
	}
	return nil, nil

}

type emailIntegrationResponse struct {
	Err    int              `json:"err"`
	Result EmailIntegration `json:"result"`
}

type pagerDutyIntegrationResponse struct {
	Err    int                  `json:"err"`
	Result PagerDutyIntegration `json:"result"`
}

type slackIntegrationResponse struct {
	Err    int              `json:"err"`
	Result SlackIntegration `json:"result"`
}

type webhookIntegrationResponse struct {
	Err    int                `json:"err"`
	Result WebhookIntegration `json:"result"`
}
