package rollbar

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const (
	team_id    = "2131231"
	user_id    = "1"
	user_email = "brian@rollbar.com"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient(ApiKey("mock_api_key"), BaseURL(server.URL+"/"))

	return func() {
		server.Close()
	}
}
func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestRemoveUserTeam(t *testing.T) {
	teardown := setup()
	defer teardown()
	handler_url := "/team/" + team_id + "/user/" + user_id
	handler_url_get := "/users/"

	mux.HandleFunc(handler_url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/remove_user.json"))
	})

	mux.HandleFunc(handler_url_get, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/users.json"))
	})

	teamtoi, _ := strconv.Atoi(team_id)
	err := client.RemoveUserTeam(user_email, teamtoi)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInviteUser(t *testing.T) {
	teardown := setup()
	defer teardown()
	handler_url := "/team/" + team_id + "/invites"

	mux.HandleFunc(handler_url, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("teams/invite_user.json"))
	})
	teamtoi, _ := strconv.Atoi(team_id)
	_, err := client.InviteUser(teamtoi, user_email)

	if err != nil {
		t.Fatal(err)
	}
}

func TestListUsers(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/users.json"))
	})
	// Email is taken from the ./testdata/fixtures/users/users.json
	user_id, err := client.getId(user_email)

	if err != nil {
		t.Fatal(err)
	}
	// This will fail if we get anything rather than a number
	strconv.Itoa(user_id)
}
