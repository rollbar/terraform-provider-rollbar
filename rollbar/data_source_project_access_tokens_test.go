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

package rollbar

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccProjectAccessTokensDataSourceNoPrefix tests reading project access
// tokens with `rollbar_project_access_tokens` data source, with no prefix
// specified.
func (s *AccSuite) TestAccProjectAccessTokensDataSourceNoPrefix() {
	rn := "data.rollbar_project_access_tokens.test"
	s.parallelTestVCR(s.T(), resource.TestCase{
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
}

// TestAccProjectAccessTokensDataSourceWithPrefix tests reading project access
// tokens with `rollbar_project_access_tokens` data source, with a specific
// prefix.
func (s *AccSuite) TestAccProjectAccessTokensDataSourceWithPrefix() {
	rn := "data.rollbar_project_access_tokens.test"
	s.parallelTestVCR(s.T(), resource.TestCase{
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

// configDataSourceProjectAccessTokens generates Terraform configuration for
// resource `rollbar_project_access_tokens`. If `prefix` is not empty, it will
// be supplied as the `prefix` argument to the data source.
func (s *AccSuite) configDataSourceProjectAccessTokens(prefix string) string {
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
	return fmt.Sprintf(tmpl, s.randName, configPrefix)
}
