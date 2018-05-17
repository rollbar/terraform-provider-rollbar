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

	vars, err := vars("users.json")

	if err != nil {
		t.Fatal(err)
	}

	teamID := vars.TeamID
	userID := vars.UserID
	userEmail := vars.UserEmail

	handUrl := fmt.Sprintf("/team/%d/user/%d", teamID, userID)
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

	err = client.RemoveUserTeam(userEmail, teamID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInviteUser(t *testing.T) {
	teardown := setup()
	defer teardown()

	vars, err := vars("users.json")

	if err != nil {
		t.Fatal(err)
	}

	teamID := vars.TeamID
	userEmail := vars.UserEmail
	handUrl := fmt.Sprintf("/team/%d/invites", teamID)

	mux.HandleFunc(handUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("teams/invite_user.json"))
	})
	_, err = client.InviteUser(teamID, userEmail)

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

	vars, err := vars("users.json")

	if err != nil {
		t.Fatal(err)
	}

	userID := vars.UserID
	userEmail := vars.UserEmail
	handUrl := "/users"

	mux.HandleFunc(handUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/users.json"))
	})
	userID, err = client.getID(userEmail)

	if err != nil {
		t.Fatal(err)
	}
	// This will fail if we get anything rather than a int.
	strconv.Itoa(userID)
}
