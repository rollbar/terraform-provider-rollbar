package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccRollbarProjectAccessTokenDataSource tests reading a project access
// token with `rollbar_project_access_token` data source.
func (s *AccSuite) TestAccRollbarProjectAccessTokenDataSource() {
	rn := "data.rollbar_project_access_token.test"

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.configDataSourceRollbarProjectAccessToken(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttrSet(rn, "access_token"),
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					resource.TestCheckResourceAttrSet(rn, "date_created"),
					resource.TestCheckResourceAttrSet(rn, "date_modified"),
					resource.TestCheckResourceAttr(rn, "name", "post_client_item"),
				),
			},
		},
	})

}

// configDataSourceRollbarProjectAccessTokens generates Terraform configuration
// for resource `rollbar_project_access_tokens`. If `prefix` is not empty, it
// will be supplied as the `prefix` argument to the data source.
func (s *AccSuite) configDataSourceRollbarProjectAccessToken() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
	
		data "rollbar_project_access_token" "test" {
			project_id = rollbar_project.test.id
			name = "post_client_item"
			depends_on = [rollbar_project.test]
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}
