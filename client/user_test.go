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
	"github.com/jarcoal/httpmock"
	"net/http"
)

// TestListUsers tests listing all Rollbar users.
func (s *Suite) TestListUsers() {
	u := apiUrl + pathUserList

	// Success
	r := responderFromFixture("user/list.json", http.StatusOK)
	httpmock.RegisterResponder("GET", u, r)
	expected := []User{
		{
			Email:    "jason.mcvetta@gmail.com",
			ID:       238101,
			Username: "jmcvetta",
		},
		{
			Email:    "cory@rollbar.com",
			ID:       2,
			Username: "coryvirok",
		},
	}
	actual, err := s.client.ListUsers()
	s.Nil(err)
	s.Subset(actual, expected)
	s.Len(actual, len(expected))

	// Internal Server Error
	r = httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, errResult500)
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ListUsers()
	s.NotNil(err)

	// Server unreachable
	httpmock.Reset()
	_, err = s.client.ListUsers()
	s.NotNil(err)

	// Unauthorized
	r = httpmock.NewJsonResponderOrPanic(http.StatusUnauthorized,
		ErrorResult{Err: 401, Message: "Unauthorized"})
	httpmock.RegisterResponder("GET", u, r)
	_, err = s.client.ListUsers()
	s.Equal(ErrUnauthorized, err)
}
