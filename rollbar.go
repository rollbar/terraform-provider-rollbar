package rollbar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiBaseUrl string = "https://api.rollbar.com/api/1/"

type Client struct {
	ApiKey string
}

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
type ListInvitesResponse struct {
	Error  int `json:err`
	Result []struct {
		Id           int    `json:"id"`
		FromUserId   int    `json:"from_user_id"`
		TeamId       int    `json:"team_id"`
		ToEmail      string `json:"to_email"`
		Status       string `json:"status"`
		DateCreated  int    `json:"date_created"`
		DateRedeemed int    `json:"date_redeemed"`
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

// type GetResponse struct {
// 	Error  int `json:err`
// 	Result struct {
// 		Username     string `json:"username"`
// 		Id           int    `json:"id"`
// 		Email        string `json:"email"`
// 		EmailEnabled int    `json:"email_enabled"`
// 	}
// }

func NewClient(apikey string) *Client {
	return &Client{
		ApiKey: apikey,
	}
}

func (s *Client) makeRequest(req *http.Request) ([]byte, error) {

	client := &http.Client{}
	resp, client_err := client.Do(req)

	if client_err != nil {
		return nil, client_err
	}

	defer resp.Body.Close()

	body, read_body_err := ioutil.ReadAll(resp.Body)

	if read_body_err != nil {
		return nil, read_body_err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

// func (s *Client) GetUser(user_id int) (*GetResponse, error) {
// 	var data GetResponse

// 	url := fmt.Sprintf("%suser/%d?access_token=%s", apiBaseUrl, user_id, s.ApiKey)
// 	req, new_request_err := http.NewRequest("GET", url, nil)

// 	if new_request_err != nil {
// 		return nil, new_request_err
// 	}

// 	bytes, make_request_err := s.makeRequest(req)

// 	if make_request_err != nil {
// 		return nil, make_request_err
// 	}

// 	unmarshal_error := json.Unmarshal(bytes, &data)

// 	if unmarshal_error != nil {
// 		return nil, unmarshal_error
// 	}

// 	return &data, nil
// }

func (s *Client) InviteUser(team_id int, email string) (*InviteResponse, error) {
	var data InviteResponse

	type requestData struct {
		Access_token string `json:"access_token"`
		Email        string `json:"email"`
	}

	url := fmt.Sprintf("%steam/%d/invites", apiBaseUrl, team_id)
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

	url := fmt.Sprintf("%susers?access_token=%s", apiBaseUrl, s.ApiKey)
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

	url := fmt.Sprintf("%steam/%d/user/%d?access_token=%s", apiBaseUrl, team_id, user_id, s.ApiKey)
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

func (s *Client) ListInvites(team_id int) (*ListInvitesResponse, error) {
	var data ListInvitesResponse

	// Invitation call has pagination.
	// We assume that we wont get to 1000 pages of invites.
	// There's a feature request to expire the invitations after some time.
	// Looping until we get an empty invitations list [].
	for i := 1; i < 1000; i++ {
		page_number := i
		url := fmt.Sprintf("%steam/%d/invites?access_token=%s&page=%d", apiBaseUrl, team_id, s.ApiKey, page_number)
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

		if len(data.Result) == 0 {
			return &data, nil
		}

		fmt.Printf("%s", data.Result)
	}

	return &data, nil

}
