package rollbar

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
func TestListUsers(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("users/users.json"))
	})
	response, err := client.ListUsers()
	fmt.Println(response)
	if err != nil {
		t.Fatal(err)
	}
}
