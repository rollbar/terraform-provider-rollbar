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
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rs/zerolog/log"
)

// TestAccProjectDataSource tests reading a project with `rollbar_project` data
// source.
func (s *AccSuite) TestAccProjectDataSource() {
	rn := "data.rollbar_project.test"

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Debug().Msg("Testing data source rollbar_project with invalid project name")
				},
				Config:      s.configDataSourceProjectNotFound(),
				ExpectError: regexp.MustCompile("no project with the name"),
			},
			{
				PreConfig: func() {
					log.Debug().Msg("Testing data source rollbar_project")
				},
				Config: s.configDataSourceProject(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", s.randName),
					resource.TestCheckResourceAttrSet(rn, "id"),
					resource.TestCheckResourceAttrSet(rn, "account_id"),
					resource.TestCheckResourceAttrSet(rn, "date_created"),
					resource.TestCheckResourceAttrSet(rn, "date_modified"),
					resource.TestCheckResourceAttr(rn, "status", "enabled"),
				),
			},
		},
	})
}

func (s *AccSuite) configDataSourceProject() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
		
		data "rollbar_project" "test" {
			name = "%s"
			depends_on = [rollbar_project.test]
		}
	`
	return fmt.Sprintf(tmpl, s.randName, s.randName)
}

func (s *AccSuite) configDataSourceProjectNotFound() string {
	// language=hcl
	tmpl := `
		data "rollbar_project" "test" {
			name = "%s"
		}
	`
	return fmt.Sprintf(tmpl, s.randName)
}
