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
	"strconv"

	"github.com/rs/zerolog/log"
)

type ServiceLink struct {
	ID       int    `model:"id" mapstructure:"id"`
	Name     string `model:"name" mapstructure:"name"`
	Template string `model:"template" mapstructure:"template"`
}

// CreateServiceLink creates a new Rollbar service_link.
func (c *RollbarAPIClient) CreateServiceLink(name, template string) (*ServiceLink, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathServiceLinkCreate
	l := log.With().
		Str("name", name).
		Logger()
	l.Debug().Msg("Creating new service link")

	resp, err := c.Resty.R().
		SetBody(map[string]string{"name": name, "template": template}).
		SetResult(serviceLinkResponse{}).
		SetError(ErrorResult{}).
		Post(u)

	if err != nil {
		l.Err(err).Msg("Error creating service link")
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	l.Debug().Msg("Service Link successfully created")
	slr := resp.Result().(*serviceLinkResponse)
	return &slr.Result, nil

}

// UpdateServiceLink updates a Rollbar service link.
func (c *RollbarAPIClient) UpdateServiceLink(id int, name, template string) (*ServiceLink, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathServiceLinkReadOrDeleteOrUpdate
	l := log.With().
		Str("name", name).
		Logger()
	l.Debug().Msg("Updating service link")

	resp, err := c.Resty.R().
		SetBody(map[string]interface{}{"name": name, "template": template}).
		SetResult(serviceLinkResponse{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"id": strconv.Itoa(id),
		}).
		Put(u)
	if err != nil {
		l.Err(err).Msg("Error updating service link")
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	l.Debug().Msg("Service Link successfully updated")
	slr := resp.Result().(*serviceLinkResponse)
	return &slr.Result, nil

}

// ReadServiceLink reads a Rollbar service link from the API. If no matching service link is found,
// returns error ErrNotFound.
func (c *RollbarAPIClient) ReadServiceLink(id int) (*ServiceLink, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathServiceLinkReadOrDeleteOrUpdate

	l := log.With().
		Int("id", id).
		Logger()
	l.Debug().Msg("Reading Service Link from API")

	resp, err := c.Resty.R().
		SetResult(serviceLinkResponse{}).
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"id": strconv.Itoa(id),
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
	slr := resp.Result().(*serviceLinkResponse)
	if slr.Err != 0 {
		l.Warn().Msg("Service Link not found")
		return nil, ErrNotFound
	}
	l.Debug().Msg("Service Link successfully read")
	return &slr.Result, nil

}

// DeleteServiceLink deletes a Rollbar service_link. If no matching service link is found,
// returns error ErrNotFound.
func (c *RollbarAPIClient) DeleteServiceLink(id int) error {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathServiceLinkReadOrDeleteOrUpdate
	l := log.With().
		Int("id", id).
		Logger()
	l.Debug().Msg("Deleting Service Link")

	resp, err := c.Resty.R().
		SetError(ErrorResult{}).
		SetPathParams(map[string]string{
			"id": strconv.Itoa(id),
		}).
		Delete(u)
	if err != nil {
		l.Err(err).Msg("Error deleting Service Link")
		return err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return err
	}
	l.Debug().Msg("Service Link successfully deleted")
	return nil
}

func (c *RollbarAPIClient) ListSerivceLinks() ([]ServiceLink, error) {
	c.m.Lock()
	defer c.m.Unlock()
	u := c.BaseURL + pathServiceLinkCreate

	l := log.With().
		Logger()
	l.Debug().Msg("Reading service links from API")

	resp, err := c.Resty.R().
		SetResult(serviceLinksResponse{}).
		SetError(ErrorResult{}).
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
	sl := resp.Result().(*serviceLinksResponse)
	if sl.Err != 0 {
		l.Warn().Msg("Service link not found")
		return nil, ErrNotFound
	}
	l.Debug().Msg("Service link successfully read")
	return sl.Result, nil
}

type serviceLinkResponse struct {
	Err    int         `json:"err"`
	Result ServiceLink `json:"result"`
}
type serviceLinksResponse struct {
	Err    int           `json:"err"`
	Result []ServiceLink `json:"result"`
}
