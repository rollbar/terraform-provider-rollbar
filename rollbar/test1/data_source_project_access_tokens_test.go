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

// TestAccProjectAccessTokensDataSourceNoTokensNoPrefix tests reading project
// access tokens with `rollbar_project_access_tokens` data source, with no
// prefix specified, from a project with zero tokens.
func (s *AccSuite) TestAccProjectAccessTokensDataSourceNoTokensNoPrefix() {
	rn := "data.rollbar_project_access_tokens.test"
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		data "rollbar_project_access_tokens" "test" {
			project_id = rollbar_project.test.id
			depends_on = [rollbar_project.test]
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
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "access_tokens.#", "0"),
				),
			},
		},
	})
}

// TestAccProjectAccessTokensDataSourceTwoTokensNoPrefix tests reading project
// access tokens with `rollbar_project_access_tokens` data source, with no
// prefix specified, from a project with two tokens.
func (s *AccSuite) TestAccProjectAccessTokensDataSourceTwoTokensNoPrefix() {
	rn := "data.rollbar_project_access_tokens.test"
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		resource "rollbar_project_access_token" "test1" {
			name = "test-token-1"
			project_id = rollbar_project.test.id
			scopes = ["read"]
		}

		resource "rollbar_project_access_token" "test2" {
			name = "test-token-2"
			project_id = rollbar_project.test.id
			scopes = ["post_server_item"]
		}

		data "rollbar_project_access_tokens" "test" {
			project_id = rollbar_project.test.id
			depends_on = [
				rollbar_project_access_token.test1,
				rollbar_project_access_token.test2,
			]
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
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "access_tokens.#", "2"),
				),
			},
		},
	})
}

// TestAccProjectAccessTokensDataSourceNoTokensWithPrefix tests reading project
// access tokens with `rollbar_project_access_tokens` data source, with a prefix
// specified, from a project with zero tokens.
func (s *AccSuite) TestAccProjectAccessTokensDataSourceNoTokensWithPrefix() {
	rn := "data.rollbar_project_access_tokens.test"
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		data "rollbar_project_access_tokens" "test" {
			project_id = rollbar_project.test.id
			#depends_on = [rollbar_project.test]
			prefix = "test-"
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
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "access_tokens.#", "0"),
				),
			},
		},
	})
}

// TestAccProjectAccessTokensDataSourceWithTokensWithPrefix tests reading
// project access tokens with `rollbar_project_access_tokens` data source from a
// project with two tokens, with a prefix matching one of those tokens.
func (s *AccSuite) TestAccProjectAccessTokensDataSourceWithTokensWithPrefix() {
	rn := "data.rollbar_project_access_tokens.test"
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name         = "%s"
		}

		resource "rollbar_project_access_token" "test1" {
			name = "foo-token"
			project_id = rollbar_project.test.id
			scopes = ["read"]
		}

		resource "rollbar_project_access_token" "test2" {
			name = "bar-token"
			project_id = rollbar_project.test.id
			scopes = ["post_server_item"]
		}

		data "rollbar_project_access_tokens" "test" {
			project_id = rollbar_project.test.id
			prefix = "foo"
			depends_on = [
				rollbar_project_access_token.test1,
				rollbar_project_access_token.test2,
			]
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
					resource.TestCheckResourceAttrSet(rn, "project_id"),
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "access_tokens.#", "1"),
				),
			},
		},
	})
}
