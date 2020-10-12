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

	// Internal Server Error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ListProjects()
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.ListProjects()
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ListProjects()
	s.Equal(ErrUnauthorized, err)
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

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("POST", u, r)
	_, err = s.client.CreateProject(name)
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.CreateProject(name)
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("POST", u, r)
	_, err = s.client.CreateProject(name)
	s.Equal(ErrUnauthorized, err)
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

	// Not Found
	er := ErrorResult{Err: 404, Message: "Not Found"}
	r = httpmock.NewJsonResponderOrPanic(http.StatusNotFound, er)
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadProject(expected.Id)
	s.Equal(ErrNotFound, err)

	// Internal server error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadProject(expected.Id)
	s.Equal(err, &errResult500)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.ReadProject(expected.Id)
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadProject(expected.Id)
	s.Equal(ErrUnauthorized, err)

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

	// Project not found
	r = httpmock.NewJsonResponderOrPanic(http.StatusNotFound,
		ErrorResult{Err: 404, Message: "Not Found"})
	httpmock.RegisterResponder("DELETE", urlDel, r)
	err = s.client.DeleteProject(delId)
	s.Equal(ErrNotFound, err)

	// Internal Server Error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("DELETE", urlDel, r)
	err = s.client.DeleteProject(delId)
	s.Equal(&errResult500, err)

	// Server unreachable
	httpmock.Reset()
	err = s.client.DeleteProject(delId)
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("DELETE", urlDel, r)
	err = s.client.DeleteProject(delId)
	s.Equal(ErrUnauthorized, err)
}
