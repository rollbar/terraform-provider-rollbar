package rollbar

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiBaseUrl string = "https://api.rollbar.com/api/1/"

type Client struct {
	ApiKey string
}

func NewClient(apikey string) *Client {
	return &Client{
		ApiKey: apikey,
	}
}

func (s *Client) makeRequest(req *http.Request) ([]byte, error) {

	client := &http.Client{}
	resp, client_err := client.Do(req)

	if client_err != nil {
		return nil, client_err
	}

	defer resp.Body.Close()

	body, read_body_err := ioutil.ReadAll(resp.Body)

	if read_body_err != nil {
		return nil, read_body_err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}
