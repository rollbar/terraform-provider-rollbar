package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccRollbarProjectAccessTokensDataSource tests reading project access
// tokens with `rollbar_project_access_tokens` data source.
func (s *Suite) TestAccRollbarProjectAccessTokensDataSource() {

	rn := "data.rollbar_project_access_tokens.test"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-acc-test-%s", randString)

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRollbarProjectAccessTokensDataSourceConfig(name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					testAccCheckProjectAccessTokensDataSourceExists(rn),

					// By default Rollbar provisions a new project with 4 access
					// tokens.
					resource.TestCheckResourceAttr(rn, "access_tokens.#", "4"),
				),
			},
		},
	})

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRollbarProjectAccessTokensDataSourceConfig(name, "post"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					testAccCheckProjectAccessTokensDataSourceExists(rn),

					// By default Rollbar provisions a new project with 4 access
					// tokens, 2 of whose names beging with "post".
					resource.TestCheckResourceAttr(rn, "access_tokens.#", "2"),
				),
			},
		},
	})
}

// testAccRollbarProjectAccessTokensDataSourceConfig generates Terraform
// configuration for resource `rollbar_project_access_tokens`. If `prefix` is
// not empty, it will be supplied as the `prefix` argument to the data source.
func testAccRollbarProjectAccessTokensDataSourceConfig(projName string, prefix string) string {
	var configPrefix string
	if prefix != "" {
		configPrefix = fmt.Sprintf(`prefix = "%s"`, prefix)
	}
	// language=terraform
	return fmt.Sprintf(`
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
	
		data "rollbar_project_access_tokens" "test" {
			project_id = rollbar_project.test.id
			depends_on = [rollbar_project.test]
			%s
		}
	`, projName, configPrefix)
}

// testAccCheckProjectAccessTokensDataSourceExists checks that the data source's
// ID is set in Terraform state.
// FIXME: This is repeated all over the place.  We need a DRY solution
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
