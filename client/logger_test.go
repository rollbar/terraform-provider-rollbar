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
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"net/http"
)

// TestRestyZeroLogger tests Resty trace logging using Zerolog as the logger.
func (s *Suite) TestRestyZeroLogger() {
	// For the most part we're just testing that nothing blows up.  No panics
	// means the test is passing.

	s.client.Resty.EnableTrace()

	u := s.client.BaseURL + pathProjectList

	// Debug log
	s.client.Resty.SetDebug(true)
	lpr := projectListResponse{}
	rOk := httpmock.NewJsonResponderOrPanic(http.StatusOK, lpr)
	httpmock.RegisterResponder("GET", u, rOk)
	_, err := s.client.ListProjects()
	s.Nil(err)

	// Warn log
	f := func(*resty.RequestLog) error {
		return nil
	}
	s.client.Resty.OnRequestLog(f)
	// Calling OnRequestLog twice triggers a message to warn log
	s.client.Resty.OnRequestLog(f)

	// Error log
	s.client.Resty.SetProxy("not_a_valid_url") // Invalid URL triggers message to error log

	// Clean up
	s.client.Resty.DisableTrace()
	s.client.Resty.SetDebug(false)
}
