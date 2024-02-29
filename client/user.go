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
	"strconv"

	"github.com/rs/zerolog/log"
)

// User represents a Rollbar user.
type User struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	Email    string `json:"email"`

	// https://github.com/rollbar/terraform-provider-rollbar/issues/65
	// EmailEnabled bool   `json:"email_enabled"`
}

// ListUsers lists all Rollbar users.
func (c *RollbarAPIClient) ListUsers(email string) (users []User, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	log.Debug().Msg("Listing users with email: " + email)
	u := c.BaseURL + pathUsers
	resp, err := c.Resty.R().
		SetResult(userListResponse{}).
		SetError(ErrorResult{}).
		SetQueryParam("email", email).
		Get(u)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return
	}
	users = resp.Result().(*userListResponse).Result.Users
	count := len(users)
	log.Debug().
		Int("count", count).
		Msg("Successfully listed users")
	return
}

// ListTestUsers is used only for testing purposes
func (c *RollbarAPIClient) ListTestUsers() (users []User, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	log.Debug().Msg("Listing users")
	u := c.BaseURL + pathUsers
	resp, err := c.Resty.R().
		SetResult(userListResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		log.Err(err).Msg("Error listing users")
		return
	}
	users = resp.Result().(*userListResponse).Result.Users
	count := len(users)
	log.Debug().
		Int("count", count).
		Msg("Successfully listed users")
	return
}

// ReadUser reads a Rollbar user from the API.
func (c *RollbarAPIClient) ReadUser(id int) (user User, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().Int("id", id).Logger()
	l.Debug().Msg("Reading user from API")
	u := c.BaseURL + pathUser
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{"userID": strconv.Itoa(id)}).
		SetResult(userReadResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		log.Err(err).Msg("Error reading user from API")
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		log.Err(err).Msg("Error reading user from API")
		return
	}
	user = resp.Result().(*userReadResponse).Result
	log.Debug().
		Interface("user", user).
		Msg("Successfully read user from API")
	return
}

// FindUserID finds the user ID for a given email.
func (c *RollbarAPIClient) FindUserID(email string) (int, error) {
	l := log.With().Str("email", email).Logger()
	l.Debug().Msg("Getting user ID from email")
	users, err := c.ListUsers(email)
	if err != nil {
		l.Err(err).Msg("Error getting user ID from email")
		return 0, err
	}
	if len(users) > 0 {
		if users[0].Email == email {
			l.Debug().Int("user_id", users[0].ID).Msg("Found user")
			return users[0].ID, nil
		}
	}
	l.Debug().Msg("No user found")
	return 0, ErrNotFound
}

// ListUserTeams lists a Rollbar user's teams.
func (c *RollbarAPIClient) ListUserTeams(userID int) (teams []Team, err error) {
	c.m.Lock()
	defer c.m.Unlock()
	l := log.With().Int("userID", userID).Logger()
	l.Debug().Msg("Reading teams for Rollbar user")
	u := c.BaseURL + pathUserTeams
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{"userID": strconv.Itoa(userID)}).
		SetResult(userTeamListResponse{}).
		SetError(ErrorResult{}).
		Get(u)
	if err != nil {
		log.Err(err).Msg("Error reading Rollbar user's teams from API")
		return
	}
	err = errorFromResponse(resp)
	if err != nil {
		log.Err(err).Msg("Error reading Rollbar user's teams from API")
		return
	}
	teams = resp.Result().(*userTeamListResponse).Result.Teams
	log.Debug().
		Interface("teams", teams).
		Msg("Successfully read Rollbar user's teams from API")
	return
}

// ListUserCustomTeams lists a Rollbar user's custom defined teams, excluding
// system teams "Everyone" and "Owners".
func (c *RollbarAPIClient) ListUserCustomTeams(userID int) (teams []Team, err error) {
	teams, err = c.ListUserTeams(userID)
	teams = filterSystemTeams(teams)
	return
}

/*
 * Containers for unmarshalling Rollbar API responses
 */

type userListResponse struct {
	Error  int `json:"err"`
	Result struct {
		Users []User `json:"users"`
	} `json:"result"`
}

type userReadResponse struct {
	Error  int  `json:"err"`
	Result User `json:"result"`
}

type userTeamListResponse struct {
	Error  int `json:"err"`
	Result struct {
		Teams []Team `json:"teams"`
	} `json:"result"`
}
