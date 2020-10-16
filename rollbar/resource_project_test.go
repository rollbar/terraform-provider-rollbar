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

package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
)

// TestAccProject tests creation and deletion of a Rollbar project.
func (s *AccSuite) TestAccProject() {
	rn := "rollbar_project.foo"

	resource.Test(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.configResourceProject(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", s.projectName),
					s.checkProjectExists(rn, s.projectName),
					s.checkProjectInProjectList(rn),
				),
			},
		},
	})
}

func (s *AccSuite) configResourceProject() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "foo" {
		  name         = "%s"
		}
	`
	return fmt.Sprintf(tmpl, s.projectName)
}

// checkProjectExists tests that the newly created project exists
func (s *AccSuite) checkProjectExists(rn string, name string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		s.Nil(err)
		c := s.provider.Meta().(*client.RollbarApiClient)
		proj, err := c.ReadProject(id)
		s.Nil(err)
		s.Equal(name, proj.Name, "project name from API does not match project name in Terraform config")
		return nil
	}
}

// checkProjectInProjectList tests that the newly created project is present in
// the list of all projects.
func (s *AccSuite) checkProjectInProjectList(rn string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		s.Nil(err)
		c := s.provider.Meta().(*client.RollbarApiClient)
		projList, err := c.ListProjects()
		s.Nil(err)
		found := false
		for _, proj := range projList {
			if proj.Id == id {
				found = true
			}
		}
		s.True(found, "Project not found in project list")
		return nil
	}
}
