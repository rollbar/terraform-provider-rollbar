package rollbar

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

// TestAccRollbarProjectAccessTokensDataSource tests reading a project with
// `rollbar_project` data source.
func TestAccRollbarProjectAccessTokensDataSource(t *testing.T) {

	rn := "data.rollbar_project_access_tokens.test"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-acc-test-%s", randString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRollbarProjectAccessTokensDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					testAccCheckProjectAccessTokensDataSourceExists(rn),

					// By default Rollbar provisions a new project with 4 access tokens.
					resource.TestCheckResourceAttr(rn, "access_tokens.#", "4"),
				),
			},
		},
	})
}

func testAccRollbarProjectAccessTokensDataSourceConfig(projName string) string {
	// language=terraform
	return fmt.Sprintf(`
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
	
		data "rollbar_project_access_tokens" "test" {
			project_id = rollbar_project.test.id
			depends_on = [rollbar_project.test]
		}
	`, projName)
}

func testAccCheckProjectAccessTokensDataSourceExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("can't find project access tokens data source: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("project access tokens data source ID not set")
		}

		return nil
	}
}
