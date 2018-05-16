package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ListUsersResponse struct {
	Error  int `json:err`
	Result struct {
		Users []struct {
			Username string `json:"username"`
			Id       int    `json:"id"`
			Email    string `json:"email"`
		}
	}
}

type InviteResponse struct {
	Error  int `json:err`
	Result struct {
		Id           int    `json:"id"`
		FromUserId   int    `json:"from_user_id"`
		TeamId       int    `json:"team_id"`
		ToEmail      string `json:"to_email"`
		Status       string `json:"status"`
		DateCreated  int    `json:"date_created"`
		DateRedeemed int    `json:"date_redeemed"`
	}
}

func (s *Client) InviteUser(teamID int, email string) (*InviteResponse, error) {
	var data InviteResponse

	type requestData struct {
		accessToken string `json:"access_token"`
		Email       string `json:"email"`
	}

	url := fmt.Sprintf("%steam/%d/invites", s.ApiBaseUrl, teamID)
	reqData := requestData{s.ApiKey, email}
	b, err := json.Marshal(reqData)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))

	if err != nil {
		return nil, err
	}

	bytes, err := s.makeRequest(req)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Client) ListUsers() (*ListUsersResponse, error) {
	var data ListUsersResponse

	url := fmt.Sprintf("%susers?access_token=%s", s.ApiBaseUrl, s.ApiKey)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	bytes, err := s.makeRequest(req)

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
func (s *Client) getId(email string) (int, error) {
	var userID int

	l, err := s.ListUsers()

	if err != nil {
		return 0, err
	}

	for _, user := range l.Result.Users {
		if user.Email == email {
			userID = user.Id
		}

	}

	return userID, nil
}

func (s *Client) RemoveUserTeam(email string, teamID int) error {
	userID, err := s.getId(email)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%steam/%d/user/%d?access_token=%s", s.ApiBaseUrl, teamID, userID, s.ApiKey)
	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		return err
	}

	_, err = s.makeRequest(req)

	if err != nil {
		return err
	}

	return nil
}
