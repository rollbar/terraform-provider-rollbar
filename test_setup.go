package rollbar

import (
	"net/http"
	"net/http/httptest"
)

const (
	teamID    = "2131231"
	userID    = "1"
	userEmail = "brian@rollbar.com"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient(ApiKey("mockkey"), BaseURL(server.URL+"/"))

	return func() {
		server.Close()
	}
}
