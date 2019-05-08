package rollbar

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestListProjectAccessTokens(t *testing.T) {
	teardown := setup()
	defer teardown()

	projectID := 12116
	handURL := fmt.Sprintf("/project/%d/access_tokens/", projectID)

	mux.HandleFunc(handURL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("project_access_tokens/list.json"))
	})

	actual, err := client.ListProjectAccessTokens(projectID)

	if err != nil {
		t.Fatal(err)
	}

	expected := []*ProjectAccessToken{
		{
			ProjectID:    projectID,
			AccessToken:  "access-token-12116-1",
			Name:         "post_client_item",
			Status:       "enabled",
			DateCreated:  1407933922,
			DateModified: 1407933922,
		},
		{
			ProjectID:    projectID,
			AccessToken:  "access-token-12116-2",
			Name:         "post_server_item",
			Status:       "enabled",
			DateCreated:  1407933922,
			DateModified: 1439579817,
		},
		{
			ProjectID:    projectID,
			AccessToken:  "access-token-12116-3",
			Name:         "write",
			Status:       "enabled",
			DateCreated:  1407933922,
			DateModified: 1407933922,
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected response %v, got %v.", expected, actual)
	}
}
