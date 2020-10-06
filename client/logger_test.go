package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"net/http"
)

// TestRestyZeroLogger tests Resty trace logging using Zerolog as the logger.
func (s *Suite) TestRestyZeroLogger() {
	// For the most part we're just testing that nothing blows up.  No panics
	// means the test is passing.

	s.client.resty.EnableTrace()

	u := apiUrl + pathProjectList

	// Debug log
	s.client.resty.SetDebug(true)
	lpr := projectListResponse{}
	rOk := httpmock.NewJsonResponderOrPanic(http.StatusOK, lpr)
	httpmock.RegisterResponder("GET", u, rOk)
	_, err := s.client.ListProjects()
	s.Nil(err)

	// Warn log
	f := func(*resty.RequestLog) error {
		return nil
	}
	s.client.resty.OnRequestLog(f)
	// Calling OnRequestLog twice triggers a message to warn log
	s.client.resty.OnRequestLog(f)

	// Error log
	s.client.resty.SetProxy("not_a_valid_url") // Invalid URL triggers message to error log
}
