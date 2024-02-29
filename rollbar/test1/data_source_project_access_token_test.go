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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccProjectAccessTokenDataSource tests reading a project access token with
// `rollbar_project_access_token` data source.
func (s *AccSuite) TestAccProjectAccessTokenDataSource() {
	rn := "data.rollbar_project_access_token.test"
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		resource "rollbar_project_access_token" "test" {
			name = "test-token"
			project_id = rollbar_project.test.id
			scopes = ["read"]
		}

		data "rollbar_project_access_token" "test" {
			project_id = rollbar_project.test.id
			name = "test-token"
			depends_on = [rollbar_project_access_token.test]
		}
	`
	config := fmt.Sprintf(tmpl, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttrSet(rn, "access_token"),
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					resource.TestCheckResourceAttrSet(rn, "date_created"),
					resource.TestCheckResourceAttrSet(rn, "date_modified"),
					resource.TestCheckResourceAttr(rn, "name", "test-token"),
				),
			},
		},
	})
}
