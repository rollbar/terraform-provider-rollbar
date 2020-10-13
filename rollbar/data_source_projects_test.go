package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"strconv"
)

// TestAccRollbarProjectsDataSource tests listing of all projects with
// `rollbar_projects` data source.
func (s *Suite) TestAccRollbarProjectsDataSource() {
	rn := "data.rollbar_projects.all"

	// How many projects should we expect in the project list, after our TF
	// config creates one more project?  There's no guarantee that the Rollbar
	// account will have either zero or non-zero project count when test is run.
	c := s.provider.Meta().(*client.RollbarApiClient)
	pl, err := c.ListProjects()
	s.Nil(err)
	currentProjectCount := len(pl)
	expectedCount := strconv.Itoa(currentProjectCount + 1)

	// Construct the name of the TF resource that should represent the
	// newly-created Rollbar project.
	//
	// FIXME: This relies on the API always returning projects in ascending
	//  order of ID.  This API behavior is not documented or guarnateed.
	index := strconv.Itoa(currentProjectCount) // index counting begins at 0; so no need to add 1 to project count
	projectNameResource := fmt.Sprintf("projects.%s.name", index)

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.configDataSourceRollbarProjects(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "projects.#"),
					resource.TestCheckResourceAttr(rn, "projects.#", expectedCount),
					resource.TestCheckResourceAttr(rn, projectNameResource, s.projectName),
				),
			},
		},
	})
}

func (s *Suite) configDataSourceRollbarProjects() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
		
		data "rollbar_projects" "all" {
			depends_on = [rollbar_project.test]
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}
