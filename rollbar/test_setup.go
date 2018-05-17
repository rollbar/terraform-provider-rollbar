package rollbar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

const (
	fixPath = "testdata/fixtures/"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

type VarsUsers struct {
	TeamID    int    `json:"TeamID"`
	UserID    int    `json:"UserID"`
	UserEmail string `json:"UserEmail"`
}

func setup() func() {
	const mockApiKey = "mockApiKey"

	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient(mockApiKey, BaseURL(server.URL+"/"))

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

// Another function is used for reading the fixture files
// because the other returns a string
// and a case should be formed and 2 values should be returned
// which makes it a bit complex.
func vars(fName string) (*VarsUsers, error) {
	var data VarsUsers

	fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/variables/%s", fixPath, fName))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileBytes, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
