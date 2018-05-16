package rollbar

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestRemoveUserTeam(t *testing.T) {
	teardown := setup()
	defer teardown()
	handUrl := fmt.Sprintf("/team/%s/user/%s", teamID, userID)
	handUrlGet := "/users/"

	mux.HandleFunc(handUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/remove_user.json"))
	})

	mux.HandleFunc(handUrlGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/users.json"))
	})

	teamTOI, _ := strconv.Atoi(teamID)
	err := client.RemoveUserTeam(userEmail, teamTOI)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInviteUser(t *testing.T) {
	teardown := setup()
	defer teardown()
	handUrl := fmt.Sprintf("/team/%s/invites", teamID)

	mux.HandleFunc(handUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("teams/invite_user.json"))
	})
	teamTOI, _ := strconv.Atoi(teamID)
	_, err := client.InviteUser(teamTOI, userEmail)

	if err != nil {
		t.Fatal(err)
	}
}

func TestListUsers(t *testing.T) {
	teardown := setup()
	defer teardown()
	handUrl := "/users"

	mux.HandleFunc(handUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/users.json"))
	})
	_, err := client.ListUsers()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetId(t *testing.T) {
	teardown := setup()
	defer teardown()
	handUrl := "/users"
	mux.HandleFunc(handUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/users.json"))
	})
	// Email is taken from the ./testdata/fixtures/users/users.json
	userID, err := client.getId(userEmail)

	if err != nil {
		t.Fatal(err)
	}
	// This will fail if we get anything rather than a number
	strconv.Itoa(userID)
}
