package rollbar

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccRollbarProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: nil,
	})
}

//func testAccProvider() {
//	return map[string]*schema.Provider{
//		"rollbar":
//	}
//}

func testAccPreCheck(t *testing.T) {
	if token := os.Getenv("HASHICUPS_USERNAME"); token == "" {
		t.Fatal("HASHICUPS_USERNAME must be set for acceptance tests")
	}
}
