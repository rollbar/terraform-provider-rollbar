package client

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"strings"
)

var errResult500 = ErrorResult{Err: 500, Message: "Internal Server Error"}

func (s *ClientTestSuite) TestListProjects() {
	u := apiUrl + pathProjectList

	// Success
	stringResponse := httpmock.NewStringResponse(200,
		fixture("projects/list.json"))
	stringResponse.Header.Add("Content-Type", "application/json")
	r := httpmock.ResponderFromResponse(stringResponse)
	httpmock.RegisterResponder("GET", u, r)
	expected := []Project{
		{
			Id:           106671,
			AccountId:    8608,
			DateCreated:  1489139046,
			DateModified: 1549293583,
			Name:         "Client-Config",
			Status:       "enabled",
		},
		{
			Id:           12116,
			AccountId:    8608,
			DateCreated:  1407933922,
			DateModified: 1556814300,
			Name:         "My",
			Status:       "enabled",
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
}

func (s *ClientTestSuite) TestCreateProject() {
	u := apiUrl + pathProjectCreate
	name := gofakeit.HackerNoun()

	// Success
	f := fmt.Sprintf(fixture("projects/read.json"), name)
	// FIXME: The actual Rollbar API sends http.StatusOK; but it
	//  _should_ send http.StatusCreated
	stringResponse := httpmock.NewStringResponse(http.StatusOK, f)
	stringResponse.Header.Add("Content-Type", "application/json")
	r := func(req *http.Request) (*http.Response, error) {
		p := Project{}
		err := json.NewDecoder(req.Body).Decode(&p)
		s.Nil(err)
		s.Equal(name, p.Name)
		return stringResponse, nil
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
}

func (s *ClientTestSuite) TestReadProject() {
	var expected Project
	gofakeit.Struct(&expected)
	u := apiUrl + pathProjectRead
	u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(expected.Id))

	// Success
	pr := ProjectResult{Err: 0, Result: expected}
	responder := httpmock.NewJsonResponderOrPanic(http.StatusOK, pr)
	httpmock.RegisterResponder("GET", u, responder)
	actual, err := s.client.ReadProject(expected.Id)
	s.Nil(err)
	s.Equal(&expected, actual)

	// Not Found
	er := ErrorResult{Err: 404, Message: "Not Found"}
	r := httpmock.NewJsonResponderOrPanic(http.StatusNotFound, er)
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
}

func (s *ClientTestSuite) TestDeleteProject() {
	delId := gofakeit.Number(0, 1000000)
	urlDel := apiUrl + pathProjectDelete
	urlDel = strings.ReplaceAll(urlDel, "{projectId}", strconv.Itoa(delId))
	urlList := apiUrl + pathProjectList

	// Success
	plr := ProjectListResult{}
	for len(plr.Result) < 3 {
		var p Project
		gofakeit.Struct(&p)
		if p.Id != delId {
			plr.Result = append(plr.Result, p)
		}
	}
	listResponder := httpmock.NewJsonResponderOrPanic(http.StatusOK, plr)
	delResponder := httpmock.NewJsonResponderOrPanic(http.StatusOK, nil)
	httpmock.RegisterResponder("GET", urlList, listResponder)
	httpmock.RegisterResponder("DELETE", urlDel, delResponder)
	err := s.client.DeleteProject(delId)
	s.Nil(err)
	projList, err := s.client.ListProjects()
	s.Nil(err)
	for _, proj := range projList {
		s.NotEqual(delId, proj.Id)
	}
	for _, count := range httpmock.GetCallCountInfo() {
		s.Equal(1, count)
	}

	// Project not found
	r := httpmock.NewJsonResponderOrPanic(http.StatusNotFound,
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
}
