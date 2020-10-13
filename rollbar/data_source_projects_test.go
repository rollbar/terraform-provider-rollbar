package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"strconv"
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
				Config: s.configDataSourceRollbarProjects(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "projects.#"),
					s.checkRollbarProjectInProjectDataSource(rn),
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

// checkRollbarProjectInProjectDataSource tests that newly created project is in
// the list of all projects returned by data source `rollbar_projects`.
func (s *Suite) checkRollbarProjectInProjectDataSource(rn string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		// How many projects should we expect in the project list?
		c := s.provider.Meta().(*client.RollbarApiClient)
		pl, err := c.ListProjects()
		s.Nil(err)
		expectedCount := strconv.Itoa(len(pl))
		err = resource.TestCheckResourceAttr(rn, "projects.#", expectedCount)(ts)
		if err != nil {
			return err
		}

		// Does our project appear as expected in the data source output?
		//
		// FIXME: This relies on the API always returning projects in ascending
		//  order of ID.  This API behavior is not documented or guarnateed.
		//
		// Construct the name of the TF resource that should represent the
		// newly-created Rollbar project. Indexes begin at 0, so we must
		// subtract one from the total item count to get the index of the final
		// project in the list.
		index := strconv.Itoa(len(pl) - 1)
		projectNameResource := fmt.Sprintf("projects.%s.name", index)
		err = resource.TestCheckResourceAttr(rn, projectNameResource, s.projectName)(ts)
		if err != nil {
			return err
		}
		return nil
	}
}
