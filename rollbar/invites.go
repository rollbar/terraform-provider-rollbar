package rollbar

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ListInvitesResultResponse struct {
	ID           int    `json:"id"`
	FromUserID   int    `json:"from_user_id"`
	TeamID       int    `json:"team_id"`
	ToEmail      string `json:"to_email"`
	Status       string `json:"status"`
	DateCreated  int    `json:"date_created"`
	DateRedeemed int    `json:"date_redeemed"`
}

type ListInvitesResponse struct {
	Error  int `json:"err"`
	Result []ListInvitesResultResponse
}

func (c *Client) ListInvites(teamID int) ([]ListInvitesResultResponse, error) {
	var invites []ListInvitesResultResponse
	// Invitation call has pagination.
	// There's a feature request to expire the invitations after some time.
	// Looping until we get an empty invitations list [].
	// Page=0 and page=1 return the same result.
	for i := 1; ; i++ {
		var data ListInvitesResponse
		var dataResponses []ListInvitesResultResponse

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
		dataResponses = data.Result
		invites = append(invites, dataResponses...)

		if err != nil {
			return nil, err
		}

		if len(data.Result) == 0 {
			break
		}
	}
	return invites, nil
}
