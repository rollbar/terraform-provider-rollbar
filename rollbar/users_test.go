package rollbar

import (
	"fmt"
	"net/http"
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

	handURL := fmt.Sprintf("/team/%d/user/%d", teamID, userID)
	handURLGet := "/users/"

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/remove_user.json"))
	})

	mux.HandleFunc(handURLGet, func(w http.ResponseWriter, r *http.Request) {
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
	handURL := fmt.Sprintf("/team/%d/invites", teamID)

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
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

	handURL := "/users"

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
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
	handURL := "/users"

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/users.json"))
	})

	userID, err = client.getID(userEmail)

	if err != nil {
		t.Fatal(err)
	}

	varType := fmt.Sprintf("%T", userID)

	if varType != "int" {
		t.Errorf("Expected 'integer', got: '%T'", userID)
	}
}

func GetUser(t *testing.T) {

	teardown := setup()
	defer teardown()

	vars, err := vars("users.json")

	if err != nil {
		t.Fatal(err)
	}

	userEmail := vars.UserEmail
	userID, err := client.getID(userEmail)
	varType := fmt.Sprintf("%T", userID)

	if varType != "int" {
		t.Errorf("Expected 'integer', got: '%T'", userID)
	}
}
