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
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"strings"
)

// TestUpdateIntegraion tests updating a Rollbar integration.
func (s *Suite) TestUpdateIntegraion() {
	integration := "slack"
	id := 557954
	u := s.client.BaseURL + pathIntegration
	u = strings.ReplaceAll(u, "{integration}", integration)
	serviceAccountID := "123456"
	enabled := false
	showMessageButtons := false
	channel := "#demo"

	rs := responseFromFixture("integration/update.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		integ := SlackIntegration{}
		err := json.NewDecoder(req.Body).Decode(&integ.Settings)
		s.Nil(err)
		s.Equal(serviceAccountID, integ.Settings.ServiceAccountID)
		s.Equal(enabled, integ.Settings.Enabled)
		s.Equal(showMessageButtons, integ.Settings.ShowMessageButtons)
		s.Equal(channel, integ.Settings.Channel)

		return rs, nil
	}

	httpmock.RegisterResponder("PUT", u, r)
	integ, err := s.client.UpdateIntegration(integration, channel, serviceAccountID, enabled, showMessageButtons)
	slackIntegration := integ.(*SlackIntegration)
	s.Nil(err)
	s.Equal(id, slackIntegration.ProjectID)

	s.checkServerErrors("PUT", u, func() error {
		_, err = s.client.UpdateIntegration(integration, channel, serviceAccountID, enabled, showMessageButtons)

		return err
	})
}

// TestReadIntegration tests reading a Rollbar integration.
func (s *Suite) TestReadIntegration() {

	id := 557954
	integration := "slack"
	u := s.client.BaseURL + pathIntegration
	u = strings.ReplaceAll(u, "{id}", strconv.Itoa(id))
	u = strings.ReplaceAll(u, "{integration}", integration)
	serviceAccountID := "123456"
	enabled := false
	showMessageButtons := false
	channel := "#demo"

	// Success
	r := responderFromFixture("integration/read.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	integ, err := s.client.ReadIntegration(integration)
	slackIntegration := integ.(*SlackIntegration)
	s.Nil(err)
	s.Equal(serviceAccountID, slackIntegration.Settings.ServiceAccountID)
	s.Equal(enabled, slackIntegration.Settings.Enabled)
	s.Equal(showMessageButtons, slackIntegration.Settings.ShowMessageButtons)
	s.Equal(channel, slackIntegration.Settings.Channel)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ReadIntegration(integration)
		return err
	})
}
