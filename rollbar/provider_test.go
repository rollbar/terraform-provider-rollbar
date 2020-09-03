package rollbar

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func testAccProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"rollbar": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if token := os.Getenv("HASHICUPS_USERNAME"); token == "" {
		t.Fatal("HASHICUPS_USERNAME must be set for acceptance tests")
	}
}
