/*
 * Copyright (c) 2024 Rollbar, Inc.
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
	"net/http"
	"strconv"
	"strings"

	"github.com/jarcoal/httpmock"
)

// TestCreateNotification tests creating a Rollbar notification.
func (s *Suite) TestCreateNotification() {
	channel := "email"
	id := 5127954
	u := s.client.BaseURL + pathNotificationCreate
	u = strings.ReplaceAll(u, "{channel}", channel)
	action := "send_email"
	trigger := "new_item"
	status := "enabled"
	filters := []map[string]interface{}{}
	config := map[string]interface{}{}

	rs := responseFromFixture("notification/create.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		n := []Notification{}
		err := json.NewDecoder(req.Body).Decode(&n)
		s.Nil(err)
		s.Equal(trigger, n[0].Trigger)

		return rs, nil
	}

	httpmock.RegisterResponder("POST", u, r)
	notification, err := s.client.CreateNotification(channel, filters, trigger, config, status)
	s.Nil(err)
	s.Equal(trigger, notification.Trigger)
	s.Equal(action, notification.Action)
	s.Equal(id, notification.ID)
	s.Equal(status, notification.Status)

	s.checkServerErrors("POST", u, func() error {
		_, err = s.client.CreateNotification(channel, filters, trigger, config, status)
		return err
	})
}

// TestUpdateNotification tests updating a Rollbar notification.
func (s *Suite) TestUpdateNotification() {
	id := 5127954
	channel := "email"
	u := s.client.BaseURL + pathNotificationReadOrDeleteOrUpdate
	u = strings.ReplaceAll(u, "{notificationID}", strconv.Itoa(id))
	u = strings.ReplaceAll(u, "{channel}", channel)

	action := "send_email"
	trigger := "new_item"
	status := "disabled"
	filters := []map[string]interface{}{}
	config := map[string]interface{}{}

	rs := responseFromFixture("notification/update.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		n := Notification{}
		err := json.NewDecoder(req.Body).Decode(&n)
		s.Nil(err)
		s.Equal(trigger, n.Trigger)

		return rs, nil
	}

	httpmock.RegisterResponder("PUT", u, r)
	notification, err := s.client.UpdateNotification(id, channel, filters, trigger, config, status)
	s.Nil(err)
	s.Equal(trigger, notification.Trigger)
	s.Equal(action, notification.Action)
	s.Equal(id, notification.ID)
	s.Equal(status, notification.Status)

	s.checkServerErrors("PUT", u, func() error {
		_, err = s.client.UpdateNotification(id, channel, filters, trigger, config, status)
		return err
	})
}

// TestReadNotification tests reading a Rollbar notification.
func (s *Suite) TestReadNotification() {

	id := 5127954
	channel := "email"
	u := s.client.BaseURL + pathNotificationReadOrDeleteOrUpdate
	u = strings.ReplaceAll(u, "{notificationID}", strconv.Itoa(id))
	u = strings.ReplaceAll(u, "{channel}", channel)

	// Success
	r := responderFromFixture("notification/read.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	n, err := s.client.ReadNotification(id, channel)
	s.Nil(err)
	s.Equal("new_item", n.Trigger)
	s.Equal("disabled", n.Status)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ReadNotification(id, channel)
		return err
	})

	// Try to read a deleted notification
	r = responderFromFixture("project/read_deleted.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	n, err = s.client.ReadNotification(id, channel)
	s.Equal(ErrNotFound, err)
	s.Nil(n)
}

// TestDeleteNotification tests deleting a Rollbar notification.
func (s *Suite) TestDeleteNotification() {
	id := 5127954
	channel := "email"
	u := s.client.BaseURL + pathNotificationReadOrDeleteOrUpdate
	u = strings.ReplaceAll(u, "{notificationID}", strconv.Itoa(id))
	u = strings.ReplaceAll(u, "{channel}", channel)

	// Success
	r := responderFromFixture("project/delete.json", http.StatusOK)
	httpmock.RegisterResponder("DELETE", u, r)
	err := s.client.DeleteNotification(id, channel)
	s.Nil(err)

	s.checkServerErrors("DELETE", u, func() error {
		return s.client.DeleteNotification(id, channel)
	})
}
