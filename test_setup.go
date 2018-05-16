package rollbar

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

const (
	teamID    = "2131231"
	userID    = "1"
	userEmail = "brian@rollbar.com"
	fixPath   = "testdata/fixtures/"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient(ApiKey("mockapikey"), BaseURL(server.URL+"/"))

	return func() {
		server.Close()
	}
}

func fixture(path string) string {
	b, err := ioutil.ReadFile(fixPath + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}
