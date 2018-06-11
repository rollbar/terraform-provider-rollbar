package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (c *Client) InviteUser(teamID int, email string) (*InviteResponse, error) {
	var data InviteResponse

	type requestData struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	}

	url := fmt.Sprintf("%steam/%d/invites", c.ApiBaseUrl, teamID)
	reqData := requestData{c.ApiKey, email}
	b, err := json.Marshal(reqData)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))

	if err != nil {
		return nil, err
	}

	bytes, err := c.makeRequest(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) ListUsers() (*ListUsersResponse, error) {
	var data ListUsersResponse

	url := fmt.Sprintf("%susers?access_token=%s", c.ApiBaseUrl, c.ApiKey)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	bytes, err := c.makeRequest(req)

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

func (c *Client) GetUser(email string) (int, error) {
	userID, err := c.getID(email)
	if err != nil {
		return 0, fmt.Errorf("There was a problem with getting the user id %s", err)
	}
	return userID, nil

}

func (c *Client) RemoveUserTeam(email string, teamID int) error {
	userID, err := c.GetUser(email)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%steam/%d/user/%d?access_token=%s", c.ApiBaseUrl, teamID, userID, c.ApiKey)
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}

	_, err = c.makeRequest(req)

	if err != nil {
		return err
	}

	return nil
}
