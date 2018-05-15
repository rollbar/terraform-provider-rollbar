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

func (s *Client) InviteUser(team_id int, email string) (*InviteResponse, error) {
	var data InviteResponse

	type requestData struct {
		Access_token string `json:"access_token"`
		Email        string `json:"email"`
	}

	url := fmt.Sprintf("%steam/%d/invites", s.ApiBaseUrl, team_id)
	r := requestData{s.ApiKey, email}
	b, unmarshal_error := json.Marshal(r)

	if unmarshal_error != nil {
		return nil, unmarshal_error
	}

	req, new_request_err := http.NewRequest("POST", url, bytes.NewBuffer(b))

	if new_request_err != nil {
		return nil, new_request_err
	}

	bytes, make_request_err := s.makeRequest(req)

	if make_request_err != nil {
		return nil, make_request_err
	}

	unmarshal_error = json.Unmarshal(bytes, &data)

	if unmarshal_error != nil {
		return nil, unmarshal_error
	}
	return &data, nil
}

func (s *Client) ListUsers() (*ListUsersResponse, error) {
	var data ListUsersResponse

	url := fmt.Sprintf("%susers?access_token=%s", s.ApiBaseUrl, s.ApiKey)
	req, new_request_err := http.NewRequest("GET", url, nil)

	if new_request_err != nil {
		return nil, new_request_err
	}

	bytes, make_request_err := s.makeRequest(req)

	if make_request_err != nil {
		return nil, make_request_err
	}

	unmarshal_error := json.Unmarshal(bytes, &data)

	if unmarshal_error != nil {
		return nil, unmarshal_error
	}

	return &data, nil
}

func (s *Client) getId(email string) (int, error) {
	var user_id int

	l, list_err := s.ListUsers()

	if list_err != nil {
		return 0, list_err
	}

	for _, user := range l.Result.Users {
		if user.Email == email {
			user_id = user.Id
		}

	}

	return user_id, nil
}

func (s *Client) RemoveUserTeam(email string, team_id int) error {
	user_id, get_id_err := s.getId(email)

	if get_id_err != nil {
		return get_id_err
	}

	url := fmt.Sprintf("%steam/%d/user/%d?access_token=%s", s.ApiBaseUrl, team_id, user_id, s.ApiKey)
	req, new_request_err := http.NewRequest("DELETE", url, nil)

	if new_request_err != nil {
		return new_request_err
	}

	_, request_err := s.makeRequest(req)

	if request_err != nil {
		return request_err
	}

	return nil
}
