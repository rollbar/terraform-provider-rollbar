package rollbar

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func (s *Client) ListInvites(team_id int) (*ListInvitesResponse, error) {
	var data ListInvitesResponse

	// Invitation call has pagination.
	// We assume that we wont get to 1000 pages of invites.
	// There's a feature request to expire the invitations after some time.
	// Looping until we get an empty invitations list [].
	for i := 1; i < 1000; i++ {
		page_number := i
		url := fmt.Sprintf("%steam/%d/invites?access_token=%s&page=%d", s.ApiBaseUrl, team_id, s.ApiKey, page_number)
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

		fmt.Printf("%+v", data.Result)
	}

	return &data, nil

}
