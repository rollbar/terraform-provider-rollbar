/*
 * Copyright (c) 2020 Rollbar, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package test1

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
)

// TestAccProjectsDataSource tests listing of all projects with
// `rollbar_projects` data source.
func (s *AccSuite) TestAccProjectsDataSource() {
	rn := "data.rollbar_projects.all"

	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.configDataSourceProjects(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "projects.#"),
					s.checkProjectInProjectDataSource(rn),
				),
			},
		},
	})
}

func (s *AccSuite) configDataSourceProjects() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
		
		data "rollbar_projects" "all" {
			depends_on = [rollbar_project.test]
		}
	`
	return fmt.Sprintf(tmpl, s.randName)
}

// checkProjectInProjectDataSource tests that newly created project is in the
// list of all projects returned by data source `rollbar_projects`.
func (s *AccSuite) checkProjectInProjectDataSource(rn string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		// How many projects should we expect in the project list?
		c := s.provider.Meta().(map[string]*client.RollbarAPIClient)[schemaKeyToken]
		pl, err := c.ListProjects()
		s.Nil(err)
		expectedCount := strconv.Itoa(len(pl))
		err = resource.TestCheckResourceAttr(rn, "projects.#", expectedCount)(ts)
		if err != nil {
			return err
		}

		/*

			FIXME: API behavior is not consistent.  We need a different way to check
			 data source is correctly populated.

			// Does our project appear as expected in the data source output?
			//
			// FIXME: This relies on the API always returning projects in ascending
			//  order of ID.  This API behavior is not documented or guaranteed.
			//
			// Construct the name of the TF resource that should represent the
			// newly-created Rollbar project. Indexes begin at 0, so we must
			// subtract one from the total item count to get the index of the final
			// project in the list.
			index := strconv.Itoa(len(pl) - 1)
			projectNameResource := fmt.Sprintf("projects.%s.name", index)
			err = resource.TestCheckResourceAttr(rn, projectNameResource, s.randName)(ts)
			if err != nil {
				return err
			}

		*/
		return nil
	}
}
