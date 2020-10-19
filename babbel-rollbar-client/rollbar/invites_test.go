package rollbar

import (
	"fmt"
	"net/http"
	"testing"
)

func TestListInvites(t *testing.T) {
	teardown := setup()
	defer teardown()

	var fixtureName string
	vars, err := vars("users.json")

	if err != nil {
		t.Fatal(err)
	}

	teamID := vars.TeamID
	handURL := fmt.Sprintf("/team/%d/invites/", teamID)

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
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

	_, err = client.ListInvites(teamID)

	if err != nil {
		t.Fatal(err)
	}
}
