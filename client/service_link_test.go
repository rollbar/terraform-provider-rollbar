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

// TestCreateServiceLink tests creating a Rollbar service_link.
func (s *Suite) TestCreateServiceLink() {

	id := 5127954
	u := s.client.BaseURL + pathServiceLinkCreate
	name := "some_name"
	template := "some_template.{{aaa}}"

	rs := responseFromFixture("service_link/create.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		sl := ServiceLink{}
		err := json.NewDecoder(req.Body).Decode(&sl)
		s.Nil(err)
		s.Equal(name, sl.Name)
		s.Equal(template, sl.Template)

		return rs, nil
	}

	httpmock.RegisterResponder("POST", u, r)
	serviceLink, err := s.client.CreateServiceLink(name, template)
	s.Nil(err)
	s.Equal(name, serviceLink.Name)
	s.Equal(template, serviceLink.Template)
	s.Equal(id, serviceLink.ID)

	s.checkServerErrors("POST", u, func() error {
		_, err = s.client.CreateServiceLink(name, template)
		return err
	})
}

// TestUpdateServiceLink tests updating a Rollbar service_link.
func (s *Suite) TestUpdateServiceLink() {
	id := 5127954
	u := s.client.BaseURL + pathServiceLinkReadOrDeleteOrUpdate
	u = strings.ReplaceAll(u, "{id}", strconv.Itoa(id))
	name := "some_name"
	template := "some_template.{{aaa}}"

	rs := responseFromFixture("service_link/update.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		sl := ServiceLink{}
		err := json.NewDecoder(req.Body).Decode(&sl)
		s.Nil(err)
		s.Equal(name, sl.Name)
		s.Equal(template, sl.Template)

		return rs, nil
	}

	httpmock.RegisterResponder("PUT", u, r)
	serviceLink, err := s.client.UpdateServiceLink(id, name, template)
	s.Nil(err)
	s.Equal(name, serviceLink.Name)
	s.Equal(template, serviceLink.Template)
	s.Equal(id, serviceLink.ID)

	s.checkServerErrors("PUT", u, func() error {
		_, err = s.client.UpdateServiceLink(id, name, template)
		return err
	})
}

// TestReadServiceLink tests reading a Rollbar service_link.
func (s *Suite) TestReadServiceLink() {

	id := 5127954
	u := s.client.BaseURL + pathServiceLinkReadOrDeleteOrUpdate
	u = strings.ReplaceAll(u, "{id}", strconv.Itoa(id))
	name := "some_name"
	template := "some_template.{{aaa}}"

	// Success
	r := responderFromFixture("service_link/read.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	serviceLink, err := s.client.ReadServiceLink(id)
	s.Nil(err)
	s.Equal(name, serviceLink.Name)
	s.Equal(template, serviceLink.Template)
	s.Equal(id, serviceLink.ID)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ReadServiceLink(id)
		return err
	})

	// Try to read a deleted notification
	r = responderFromFixture("service_link/read_deleted.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	serviceLink, err = s.client.ReadServiceLink(id)
	s.Equal(ErrNotFound, err)
	s.Nil(serviceLink)
}

// TestDeleteServiceLink tests deleting a Rollbar service_link.
func (s *Suite) TestDeleteServiceLink() {
	id := 5127954
	u := s.client.BaseURL + pathServiceLinkReadOrDeleteOrUpdate
	u = strings.ReplaceAll(u, "{id}", strconv.Itoa(id))

	// Success
	r := responderFromFixture("service_link/delete.json", http.StatusOK)
	httpmock.RegisterResponder("DELETE", u, r)
	err := s.client.DeleteServiceLink(id)
	s.Nil(err)

	s.checkServerErrors("DELETE", u, func() error {
		return s.client.DeleteServiceLink(id)
	})
}
