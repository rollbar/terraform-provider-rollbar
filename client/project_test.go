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
	stringResponse := httpmock.NewStringResponse(200, projectListJsonResponse)
	stringResponse.Header.Add("Content-Type", "application/json")
	r := httpmock.ResponderFromResponse(stringResponse)
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
	stringResponse := httpmock.NewStringResponse(http.StatusOK, projectCreateJsonResponse)
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

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("POST", u, r)
	_, err = s.client.CreateProject(name)
	s.Equal(ErrUnauthorized, err)
}

func (s *Suite) TestReadProject() {
	var expected Project
	gofakeit.Struct(&expected)
	u := apiUrl + pathProjectRead
	u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(expected.Id))

	// Success
	pr := projectResponse{Err: 0, Result: expected}
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

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ReadProject(expected.Id)
	s.Equal(ErrUnauthorized, err)
}

func (s *Suite) TestDeleteProject() {
	delId := gofakeit.Number(0, 1000000)
	urlDel := apiUrl + pathProjectDelete
	urlDel = strings.ReplaceAll(urlDel, "{projectId}", strconv.Itoa(delId))
	urlList := apiUrl + pathProjectList

	// Success
	plr := projectListResponse{}
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

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("DELETE", urlDel, r)
	err = s.client.DeleteProject(delId)
	s.Equal(ErrUnauthorized, err)
}

// Contains several deleted projects to test workaround of API bug.
// See https://github.com/rollbar/terraform-provider-rollbar/issues/23
// language=json
const projectListJsonResponse = `
{
  "err": 0,
  "result": [
    {
      "id": 411692,
      "account_id": 317418,
      "settings_data": {
        "integrations": {
          "jira": {},
          "clubhouse": {},
          "bitbucket": {},
          "github": {},
          "trello": {},
          "slack": {},
          "datadog": {},
          "pagerduty": {},
          "gitlab": {},
          "webhook": {},
          "victorops": {},
          "ciscospark": {},
          "asana": {},
          "pivotal": {},
          "campfire": {},
          "azuredevops": {},
          "sprintly": {},
          "hipchat": {},
          "lightstep": {},
          "email": {},
          "flowdock": {}
        },
        "grouping": {
          "auto_upgrade": true,
          "recent_versions": [
            "5.0.0"
          ]
        }
      },
      "date_created": 1602083693,
      "date_modified": 1602083695,
      "name": null
    },
    {
      "id": 411701,
      "account_id": 317418,
      "settings_data": {
        "integrations": {
          "jira": {},
          "clubhouse": {},
          "bitbucket": {},
          "github": {},
          "trello": {},
          "slack": {},
          "datadog": {},
          "pagerduty": {},
          "gitlab": {},
          "webhook": {},
          "victorops": {},
          "ciscospark": {},
          "asana": {},
          "pivotal": {},
          "campfire": {},
          "azuredevops": {},
          "sprintly": {},
          "hipchat": {},
          "lightstep": {},
          "email": {},
          "flowdock": {}
        },
        "grouping": {
          "auto_upgrade": true,
          "recent_versions": [
            "5.0.0"
          ]
        }
      },
      "date_created": 1602084945,
      "date_modified": 1602085330,
      "name": null
    },
    {
      "id": 411704,
      "account_id": 317418,
      "status": "enabled",
      "settings_data": {
        "grouping": {
          "auto_upgrade": true,
          "recent_versions": [
            "5.0.0"
          ]
        }
      },
      "date_created": 1602085345,
      "date_modified": 1602085345,
      "name": "bar"
    },
    {
      "id": 411703,
      "account_id": 317418,
      "status": "enabled",
      "settings_data": {
        "grouping": {
          "auto_upgrade": true,
          "recent_versions": [
            "5.0.0"
          ]
        }
      },
      "date_created": 1602085340,
      "date_modified": 1602085340,
      "name": "foo"
    }
  ]
}
`

// language=json
const projectCreateJsonResponse = `
{
    "err": 0,
    "result": {
        "account_id": 317418,
        "date_created": 1602086539,
        "date_modified": 1602086539,
        "id": 411708,
        "name": "baz",
        "settings_data": {
            "grouping": {
                "auto_upgrade": true,
                "recent_versions": [
                    "5.0.0"
                ]
            }
        },
        "status": "enabled"
    }
}
`
