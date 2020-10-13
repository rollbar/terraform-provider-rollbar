package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccRollbarProjectsDataSource tests listing of all projects with
// `rollbar_projects` data source.
func (s *Suite) TestAccRollbarProjectsDataSource() {
	rn := "data.rollbar_projects.all"

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.testAccRollbarProjectsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "projects.#"),
					resource.TestCheckResourceAttr(rn, "projects.#", "1"),
					resource.TestCheckResourceAttr(rn, "projects.0.name", s.projectName),
				),
			},
		},
	})
}

func (s *Suite) testAccRollbarProjectsDataSourceConfig() string {
	return fmt.Sprintf(`
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
		
		data "rollbar_projects" "all" {
			depends_on = [rollbar_project.test]
		}
	`, s.projectName)
}
