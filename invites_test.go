package rollbar

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestListInvites(t *testing.T) {
	teardown := setup()
	defer teardown()
	handUrl := fmt.Sprintf("/team/%s/invites/", teamID)

	mux.HandleFunc(handUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check if we've hit the second page and return an empty paginated
		// response thus simulating the apis behaviour.
		splstr := strings.Split(r.URL.RawQuery, "&")
		if splstr[1] == "page=2" {
			fmt.Fprint(w, fixture("teams/empty_invites.json"))
		} else {
			fmt.Fprint(w, fixture("teams/invites.json"))
		}
	})

	// team_id should e int for this particular func
	teamTOI, _ := strconv.Atoi(teamID)
	_, err := client.ListInvites(teamTOI)

	if err != nil {
		t.Fatal(err)
	}
}
