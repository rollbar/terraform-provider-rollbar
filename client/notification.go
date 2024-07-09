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
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Notification struct {
	ID      int                    `model:"id" mapstructure:"id"`
	Status  string                 `model:"status" mapstructure:"status"`
	Action  string                 `model:"action" mapstructure:"action"`
	Trigger string                 `model:"trigger" mapstructure:"trigger"`
	Channel string                 `model:"channel" mapstructure:"channel"`
	Filters []interface{}          `model:"filters" mapstructure:"filters"`
	Config  map[string]interface{} `model:"config" mapstructure:"config"`
}

// CreateNotification creates a new Rollbar notification.
func (c *RollbarAPIClient) CreateNotification(channel string, filters, trigger, config interface{}, status string) (*Notification, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathNotificationCreate
	u = strings.ReplaceAll(u, "{channel}", channel)
	l := log.With().
		Str("channel", channel).
		Logger()
	l.Debug().Msg("Creating new notification")

	resp, err := c.Resty.R().
		SetBody([]map[string]interface{}{{"filters": filters, "trigger": trigger, "config": config, "status": status}}).
		SetResult(notificationsResponse{}).
		SetError(ErrorResult{}).
		Post(u)
	if err != nil {
		l.Err(err).Msg("Error creating notification")
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	l.Debug().Msg("Notification successfully created")
	nr := resp.Result().(*notificationsResponse)
	return &nr.Result[0], nil

}

// UpdateNotification updates a Rollbar notification.
func (c *RollbarAPIClient) UpdateNotification(notificationID int, channel string, filters, trigger, config interface{}, status string) (*Notification, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathNotificationReadOrDeleteOrUpdate
	l := log.With().
		Str("channel", channel).
		Logger()
	l.Debug().Msg("Updating notification")

	resp, err := c.Resty.R().
		SetBody(map[string]interface{}{"filters": filters, "trigger": trigger, "config": config, "status": status}).
		SetResult(notificationResponse{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"notificationID": strconv.Itoa(notificationID),
			"channel":        channel,
		}).
		Put(u)
	if err != nil {
		l.Err(err).Msg("Error updating notification")
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	l.Debug().Msg("Notification successfully updated")
	nr := resp.Result().(*notificationResponse)
	return &nr.Result, nil

}

// ReadNotification reads a Rollbar notification from the API. If no matching notification is found,
// returns error ErrNotFound.
func (c *RollbarAPIClient) ReadNotification(notificationID int, channel string) (*Notification, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathNotificationReadOrDeleteOrUpdate

	l := log.With().
		Int("notificationID", notificationID).
		Logger()
	l.Debug().Msg("Reading notification from API")

	resp, err := c.Resty.R().
		SetResult(notificationResponse{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"notificationID": strconv.Itoa(notificationID),
			"channel":        channel,
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
	nr := resp.Result().(*notificationResponse)
	if nr.Err != 0 {
		l.Warn().Msg("Notification not found")
		return nil, ErrNotFound
	}
	l.Debug().Msg("Notification successfully read")
	return &nr.Result, nil

}

// DeleteNotification deletes a Rollbar notification. If no matching notification is found,
// returns error ErrNotFound.
func (c *RollbarAPIClient) DeleteNotification(notificationID int, channel string) error {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathNotificationReadOrDeleteOrUpdate
	l := log.With().
		Int("notificationID", notificationID).
		Logger()
	l.Debug().Msg("Deleting notification")

	resp, err := c.Resty.R().
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"notificationID": strconv.Itoa(notificationID),
			"channel":        channel,
		}).
		Delete(u)
	if err != nil {
		l.Err(err).Msg("Error deleting notification")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return err
	}
	l.Debug().Msg("Notifications successfully deleted")
	return nil
}

func (c *RollbarAPIClient) ListNotifications(channel string) ([]Notification, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathNotificationCreate

	l := log.With().
		Logger()
	l.Debug().Msg("Reading notifications from API")

	resp, err := c.Resty.R().
		SetResult(notificationsResponse{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"channel": channel,
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
	nr := resp.Result().(*notificationsResponse)
	if nr.Err != 0 {
		l.Warn().Msg("Notification not found")
		return nil, ErrNotFound
	}
	l.Debug().Msg("Notification successfully read")
	return nr.Result, nil
}

type notificationResponse struct {
	Err    int          `json:"err"`
	Result Notification `json:"result"`
}

type notificationsResponse struct {
	Err    int            `json:"err"`
	Result []Notification `json:"result"`
}
