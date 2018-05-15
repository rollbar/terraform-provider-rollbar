package rollbar

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestListInvites(t *testing.T) {
	teardown := setup()
	defer teardown()
	handler_url := "/team/" + team_id + "/invites/"

	mux.HandleFunc(handler_url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("teams/invites.json"))
	})

	teamtoi, _ := strconv.Atoi(team_id)
	_, err := client.ListInvites(teamtoi)

	if err != nil {
		t.Fatal(err)
	}

}
