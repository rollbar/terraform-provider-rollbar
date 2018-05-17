package rollbar

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ListInvitesResponse struct {
	Error  int `json:err`
	Result []struct {
		ID           int    `json:"id"`
		FromUserID   int    `json:"from_user_id"`
		TeamID       int    `json:"team_id"`
		ToEmail      string `json:"to_email"`
		Status       string `json:"status"`
		DateCreated  int    `json:"date_created"`
		DateRedeemed int    `json:"date_redeemed"`
	}
}

func (c *Client) ListInvites(teamID int) (*ListInvitesResponse, error) {
	var data ListInvitesResponse

	// Invitation call has pagination.
	// There's a feature request to expire the invitations after some time.
	// Looping until we get an empty invitations list [].
	// Page=0 and page=1 return the same result.
	for i := 1; ; i++ {
		pageNum := i
		url := fmt.Sprintf("%steam/%d/invites?access_token=%s&page=%d", c.ApiBaseUrl, teamID, c.ApiKey, pageNum)
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

		if len(data.Result) == 0 {
			break
		}
	}
	return &data, nil
}
