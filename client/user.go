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
	"github.com/rs/zerolog/log"
	"strconv"
)

// User represents a Rollbar user.
type User struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	Email    string `json:"email"`

	// https://github.com/rollbar/terraform-provider-rollbar/issues/65
	//EmailEnabled bool   `json:"email_enabled"`
}

// ListUsers lists all Rollbar users.
func (c *RollbarApiClient) ListUsers() (users []User, err error) {
	log.Debug().Msg("Listing users")
	u := apiUrl + pathUsers
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
	log.Debug().
		Interface("users", users).
		Msg("Successfully listed users")
	return
}

// ReadUser reads a Rollbar user from the API.
func (c *RollbarApiClient) ReadUser(id int) (user User, err error) {
	l := log.With().Int("id", id).Logger()
	l.Debug().Msg("Reading user from API")
	u := apiUrl + pathUser
	resp, err := c.Resty.R().
		SetPathParams(map[string]string{"userId": strconv.Itoa(id)}).
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

// UserIdFromEmail finds the user ID for a given email.  WARNING: this is a
// potentially slow call.  Don't repeat it unnecessarily.
func (c *RollbarApiClient) UserIdFromEmail(email string) (int, error) {
	l := log.With().Str("email", email).Logger()
	l.Debug().Msg("Getting user ID from email")
	users, err := c.ListUsers()
	if err != nil {
		l.Err(err).Msg("Error getting user ID from email")
		return 0, err
	}
	for _, u := range users {
		if u.Email == email {
			return u.ID, nil
		}
	}
	return 0, ErrNotFound
}

/*
// This response doesn't have pagination so it might break
// in the future.
func (c *Client) getID(email string) (int, error) {
	var userID int

	l, err := c.ListUsers()
	if err != nil {
		return 0, err
	}

	for _, user := range l.Result.Users {
		if user.ToEmail == email {
			userID = user.ID
		}
	}

	return userID, nil
}


// GetUser fetches one user.
func (c *Client) GetUser(email string) (int, error) {
	userID, err := c.getID(email)
	if err != nil {
		return 0, fmt.Errorf("There was a problem with getting the user id %s", err)
	}
	return userID, nil

}

// RemoveUserTeam removes a user from a team.
func (c *Client) RemoveUserTeam(email string, teamID int) error {
	userID, err := c.GetUser(email)
	if err != nil {
		return err
	}

	err = c.delete("team", strconv.Itoa(teamID), "user", strconv.Itoa(userID))
	if err != nil {
		return err
	}

	return nil
}
*/

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
