package rollbar

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestListInvites(t *testing.T) {
	var fixtureName string
	teardown := setup()
	defer teardown()
	handUrl := fmt.Sprintf("/team/%s/invites/", teamID)

	mux.HandleFunc(handUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		page := r.URL.Query().Get("pages")

		if page == "0" || page == "1" {
			fixtureName = "invites"
		} else {
			fixtureName = "empty_invites"
		}

		fmt.Fprint(w, fixture(fmt.Sprintf("teams/%s.json", fixtureName)))

	})

	// team_id should e int for this particular func
	teamTOI, _ := strconv.Atoi(teamID)
	_, err := client.ListInvites(teamTOI)

	if err != nil {
		t.Fatal(err)
	}
}
