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

func (s *AccSuite) TestAccTeamDataSource() {
	rn_name := "data.rollbar_team.test_name"
	rn_id := "data.rollbar_team.test_id"

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Debug().Msg("Testing data source rollbar_team with invalid name")
				},
				Config:      s.configDataSourceTeamNotFoundByName(),
				ExpectError: regexp.MustCompile("Team not found by name"),
			},
			{
				PreConfig: func() {
					log.Debug().Msg("Testing data source rollbar_team with invalid ID")
				},
				Config:      s.configDataSourceTeamNotFoundById(),
				ExpectError: regexp.MustCompile("Team not found by ID"),
			},
			{
				PreConfig: func() {
					log.Debug().Msg("Testing data source rollbar_team")
				},
				Config: s.configDataSourceTeam(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn_name, "name", s.randName),
					resource.TestCheckResourceAttr(rn_id, "name", s.randName),
					resource.TestCheckResourceAttrSet(rn_name, "id"),
					resource.TestCheckResourceAttrSet(rn_name, "account_id"),
					resource.TestCheckResourceAttrSet(rn_name, "access_level"),
					resource.TestCheckResourceAttrSet(rn_id, "id"),
					resource.TestCheckResourceAttrSet(rn_id, "account_id"),
					resource.TestCheckResourceAttrSet(rn_id, "access_level"),
				),
			},
		},
	})
}

func (s *AccSuite) configDataSourceTeam() string {
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
		  name         = "%s"
		}

		data "rollbar_team" "test_name" {
			name = "%s"
			depends_on = [rollbar_team.test]
		}

        data "rollbar_team" "test_id" {
            team_id = rollbar_team.test.id
        }
	`
	return fmt.Sprintf(tmpl, s.randName, s.randName)
}

func (s *AccSuite) configDataSourceTeamNotFoundByName() string {
	// language=hcl
	tmpl := `
		data "rollbar_team" "test" {
			name = "%s"
		}
	`
	return fmt.Sprintf(tmpl, s.randName)
}

func (s *AccSuite) configDataSourceTeamNotFoundById() string {
	// language=hcl
	tmpl := `
		data "rollbar_team" "test" {
			team_id = %d
		}
	`
	return fmt.Sprintf(tmpl, s.randID)
}
