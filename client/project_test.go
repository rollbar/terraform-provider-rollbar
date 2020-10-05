package client

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
)

var errResult = ErrorResult{Err: 500, Message: "Internal Server Error"}

var _ = Describe("Project", func() {
	When("list succeeds", func() {
		u := apiUrl + pathProjectList

		It("lists all projects", func() {
			s := fixture("projects/list.json")
			stringResponse := httpmock.NewStringResponse(200, s)
			stringResponse.Header.Add("Content-Type", "application/json")
			responder := httpmock.ResponderFromResponse(stringResponse)
			httpmock.RegisterResponder("GET", u, responder)

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
			actual, err := c.ListProjects()
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(ContainElements(expected))
			Expect(actual).To(HaveLen(len(expected)))
		})

		When("list fails", func() {
			responder := httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult)
			It("handles the error", func() {
				httpmock.RegisterResponder("GET", u, responder)
				_, err := c.ListProjects()
				Expect(err).To(MatchError(&errResult))
			})
		})
	})

	When("creating a new project", func() {
		u := apiUrl + pathProjectCreate
		name := gofakeit.HackerNoun()

		Context("and creation succeeds", func() {
			s := fmt.Sprintf(fixture("projects/read.json"), name)
			// FIXME: The actual Rollbar API sends http.StatusOK; but it
			//  _should_ send http.StatusCreated
			stringResponse := httpmock.NewStringResponse(http.StatusOK, s)
			stringResponse.Header.Add("Content-Type", "application/json")
			responder := func(req *http.Request) (*http.Response, error) {
				//proj := make(map[string]string)
				proj := Project{}
				if err := json.NewDecoder(req.Body).Decode(&proj); err != nil {
					return nil, err
				}
				if proj.Name != name {
					msg := "incorrect name sent to API"
					log.Error().
						Str("expected", name).
						Str("actual", proj.Name).
						Msg(msg)
					return nil, fmt.Errorf(msg)
				}
				return stringResponse, nil
			}
			It("is created correctly", func() {
				httpmock.RegisterResponder("POST", u, responder)
				proj, err := c.CreateProject(name)
				Expect(err).NotTo(HaveOccurred())
				Expect(proj.Name).To(Equal(name))
			})
		})

		Context("and creation fails", func() {
			s := fixture("projects/read.json")
			// FIXME: The actual Rollbar API sends http.StatusOK; but it
			//  _should_ send http.StatusCreated
			stringResponse := httpmock.NewStringResponse(http.StatusInternalServerError, s)
			stringResponse.Header.Add("Content-Type", "application/json")
			responder := httpmock.ResponderFromResponse(stringResponse)
			It("handles the error cleanly", func() {
				httpmock.RegisterResponder("POST", u, responder)
				_, err := c.CreateProject(name)
				Expect(err).To(MatchError(&ErrorResult{}))
			})
		})
	})

	When("reading a project", func() {
		var expected Project
		gofakeit.Struct(&expected)
		u := apiUrl + pathProjectRead
		u = strings.ReplaceAll(u, "{projectId}", strconv.Itoa(expected.Id))

		Context("and read succeeds", func() {
			It("has the exepected properties", func() {
				pr := ProjectResult{Err: 0, Result: expected}
				responder := httpmock.NewJsonResponderOrPanic(http.StatusOK, pr)
				httpmock.RegisterResponder("GET", u, responder)
				actual, err := c.ReadProject(expected.Id)
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(Equal(&expected))
			})
		})
		Context("and read fails", func() {
			It("handles the error", func() {
				responder := httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult)
				httpmock.RegisterResponder("GET", u, responder)
				_, err := c.ReadProject(expected.Id)
				Expect(err).To(MatchError(&errResult))
			})
		})

	})

	When("deleting a project", func() {
		delId := gofakeit.Number(0, 1000000)
		urlDel := apiUrl + pathProjectDelete
		urlDel = strings.ReplaceAll(urlDel, "{projectId}", strconv.Itoa(delId))
		urlList := apiUrl + pathProjectList

		Context("and delete succeeds", func() {
			It("is not included in project list", func() {
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
				err := c.DeleteProject(delId)
				Expect(err).NotTo(HaveOccurred())
				projList, err := c.ListProjects()
				Expect(err).NotTo(HaveOccurred())
				for _, proj := range projList {
					Expect(proj.Id).NotTo(Equal(delId))
				}
				for _, count := range httpmock.GetCallCountInfo() {
					Expect(count).To(Equal(1))
				}
			})
		})

		Context("and delete fails", func() {
			It("handles the error", func() {

			})
		})

	})

})
