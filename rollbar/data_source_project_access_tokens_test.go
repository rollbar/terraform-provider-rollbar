package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccRollbarProjectAccessTokensDataSource tests reading project access
// tokens with `rollbar_project_access_tokens` data source.
func (s *Suite) TestAccRollbarProjectAccessTokensDataSource() {
	rn := "data.rollbar_project_access_tokens.test"

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.configDataSourceProjectAccessTokens(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					s.checkResourceStateSanity(rn),

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
				Config: s.configDataSourceProjectAccessTokens("post"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					s.checkResourceStateSanity(rn),

					// By default Rollbar provisions a new project with 4 access
					// tokens, 2 of whose names beging with "post".
					resource.TestCheckResourceAttr(rn, "access_tokens.#", "2"),
				),
			},
		},
	})
}

// configDataSourceProjectAccessTokens generates Terraform
// configuration for resource `rollbar_project_access_tokens`. If `prefix` is
// not empty, it will be supplied as the `prefix` argument to the data source.
func (s *Suite) configDataSourceProjectAccessTokens(prefix string) string {
	var configPrefix string
	if prefix != "" {
		configPrefix = fmt.Sprintf(`prefix = "%s"`, prefix)
	}
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
	
		data "rollbar_project_access_tokens" "test" {
			project_id = rollbar_project.test.id
			depends_on = [rollbar_project.test]
			%s
		}
	`
	return fmt.Sprintf(tmpl, s.projectName, configPrefix)
}
