/*
 * Copyright (c) 2020 Rollbar, Inc.
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
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"strings"
)

var errResult500 = ErrorResult{Err: 500, Message: "Internal Server Error"}

func (s *Suite) TestListProjects() {
	u := apiUrl + pathProjectList

	// Success
	r := responderFromFixture("project/list.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	expected := []Project{
		{
			Id:           411704,
			Name:         "bar",
			AccountId:    317418,
			Status:       "enabled",
			DateCreated:  1602085345,
			DateModified: 1602085345,
		},
		{
			Id:           411703,
			Name:         "foo",
			AccountId:    317418,
			Status:       "enabled",
			DateCreated:  1602085340,
			DateModified: 1602085340,
		},
	}
	actual, err := s.client.ListProjects()
	s.Nil(err)
	s.Subset(actual, expected)
	s.Len(actual, len(expected))

	s.checkServerErrors("GET", u, func() error {
		_, err = s.client.ListProjects()
		return err
	})
}

func (s *Suite) TestCreateProject() {
	u := apiUrl + pathProjectCreate
	name := "baz"

	// Success
	// FIXME: The actual Rollbar API sends http.StatusOK; but it
	//  _should_ send http.StatusCreated
	rs := responseFromFixture("project/create.json", http.StatusOK)
	r := func(req *http.Request) (*http.Response, error) {
		p := Project{}
		err := json.NewDecoder(req.Body).Decode(&p)
		s.Nil(err)
		s.Equal(name, p.Name)
		return rs, nil
	}
	httpmock.RegisterResponder("POST", u, r)
	proj, err := s.client.CreateProject(name)
	s.Nil(err)
	s.Equal(name, proj.Name)

	s.checkServerErrors("POST", u, func() error {
		_, err = s.client.CreateProject(name)
		return err
	})
}

func (s *Suite) TestReadProject() {
	expected := Project{
		AccountId:    317418,
		DateCreated:  1602086539,
		DateModified: 1602086539,
		Id:           411708,
		Name:         "baz",
		Status:       "enabled",
	}
	u := apiUrl + pathProjectRead
	u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(expected.Id))

	// Success
	r := responderFromFixture("project/read.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	actual, err := s.client.ReadProject(expected.Id)
	s.Nil(err)
	s.Equal(&expected, actual)

	s.checkServerErrors("GET", u, func() error {
		_, err := s.client.ReadProject(expected.Id)
		return err
	})

	// Deleted project API bug
	// FIXME: https://github.com/rollbar/terraform-provider-rollbar/issues/23
	r = responderFromFixture("project/read_deleted.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadProject(expected.Id)
	s.Equal(ErrNotFound, err)
}

func (s *Suite) TestDeleteProject() {
	delId := gofakeit.Number(0, 1000000)
	urlDel := apiUrl + pathProjectDelete
	urlDel = strings.ReplaceAll(urlDel, "{projectId}", strconv.Itoa(delId))

	// Success
	r := responderFromFixture("project/delete.json", http.StatusOK)
	httpmock.RegisterResponder("DELETE", urlDel, r)
	err := s.client.DeleteProject(delId)
	s.Nil(err)

	s.checkServerErrors("DELETE", urlDel, func() error {
		return s.client.DeleteProject(delId)
	})
}
