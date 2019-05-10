package rollbar

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ListUsersResponse represents the list users response.
type ListUsersResponse struct {
	Error  int `json:"err"`
	Result struct {
		Users []struct {
			Username string `json:"username"`
			ID       int    `json:"id"`
			Email    string `json:"email"`
		}
	}
}

// InviteResponse represents the list invites response.
type InviteResponse struct {
	Error  int `json:"err"`
	Result struct {
		ID           int    `json:"id"`
		FromUserID   int    `json:"from_user_id"`
		TeamID       int    `json:"team_id"`
		ToEmail      string `json:"to_email"`
		Status       string `json:"status"`
		DateCreated  int    `json:"date_created"`
		DateRedeemed int    `json:"date_redeemed"`
	}
}

// InviteUser sends an invitation to a user.
func (c *Client) InviteUser(teamID int, email string) (*InviteResponse, error) {
	var data InviteResponse

	type requestData struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	}

	reqData := requestData{c.AccessToken, email}
	b, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	bytes, err := c.post(b, "team", strconv.Itoa(teamID), "invites")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// ListUsers : A function for listing the users.
func (c *Client) ListUsers() (*ListUsersResponse, error) {
	var data ListUsersResponse

	bytes, err := c.get("users")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// This response doesn't have pagination so it might break
// in the future.
func (c *Client) getID(email string) (int, error) {
	var userID int

	l, err := c.ListUsers()
	if err != nil {
		return 0, err
	}

	for _, user := range l.Result.Users {
		if user.Email == email {
			userID = user.ID
		}
	}

	return userID, nil
}

// GetUser is fetches one user.
func (c *Client) GetUser(email string) (int, error) {
	userID, err := c.getID(email)
	if err != nil {
		return 0, fmt.Errorf("There was a problem with getting the user id %s", err)
	}
	return userID, nil

}

// RemoveUserTeam is removes a user from a team.
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
