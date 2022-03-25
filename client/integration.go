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

type SlackIntegration struct {
	ProjectID int `model:"project_id" mapstructure:"project_id" json:"project_id"`
	Settings  struct {
		Channel            string `model:"channel" mapstructure:"channel" json:"channel"`
		Enabled            bool   `model:"enabled" mapstructure:"enabled" json:"enabled"`
		ShowMessageButtons bool   `model:"show_message_buttons" mapstructure:"show_message_buttons" json:"show_message_buttons"`
		ServiceAccountID   string `model:"service_account_id" mapstructure:"service_account_id" json:"service_account_id"`
	} `model:"settings" mapstructure:"settings"`
}

// UpdateIntegration updates a new Rollbar integration.
func (c *RollbarAPIClient) UpdateIntegration(integration, channel, serviceAccountID string, enabled, showMessageButtons bool) (interface{}, error) {
	u := c.BaseURL + pathIntegration
	l := log.With().
		Str("integration", integration).
		Logger()
	l.Debug().Msg("Update integration")
	resp, err := c.Resty.R().
		SetBody(map[string]interface{}{"channel": channel, "service_account_id": serviceAccountID, "enabled": enabled,
			"show_message_buttons": showMessageButtons}).
		SetResult(slackIntegrationResponse{}).
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
	sir := resp.Result().(*slackIntegrationResponse)
	return &sir.Result, nil
}

// ReadIntegration reads a Rollbar integration from the API. If no matching integration is found,
// returns error ErrNotFound.
func (c *RollbarAPIClient) ReadIntegration(integration string) (interface{}, error) {
	u := c.BaseURL + pathIntegration

	l := log.With().
		Str("integration", integration).
		Logger()
	l.Debug().Msg("Reading Integration from API")

	resp, err := c.Resty.R().
		SetResult(slackIntegrationResponse{}).
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
	slr := resp.Result().(*slackIntegrationResponse)
	if slr.Err != 0 {
		l.Warn().Msg("Integration not found")
		return nil, ErrNotFound
	}
	l.Debug().Msg("Integration successfully read")
	return &slr.Result, nil

}

type slackIntegrationResponse struct {
	Err    int              `json:"err"`
	Result SlackIntegration `json:"result"`
}
