package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccRollbarProjectDataSource tests reading a project with
// `rollbar_project` data source.
func (s *Suite) TestAccRollbarProjectDataSource() {
	rn := "data.rollbar_project.test"

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.testAccRollbarProjectDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", s.projectName),
					resource.TestCheckResourceAttrSet(rn, "id"),
					resource.TestCheckResourceAttrSet(rn, "account_id"),
					resource.TestCheckResourceAttrSet(rn, "date_created"),
				),
			},
		},
	})
}

func (s *Suite) testAccRollbarProjectDataSourceConfig() string {
	return fmt.Sprintf(`
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
		
		data "rollbar_project" "test" {
			name = "%s"
			depends_on = [rollbar_project.test]
		}
	`, s.projectName, s.projectName)
}
