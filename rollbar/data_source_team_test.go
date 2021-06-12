/*
 * Copyright (c) 2021 Rollbar, Inc.
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
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rs/zerolog/log"
)

// TestTeamDataSource tests reading a team with
// `rollbar_team` data source.
func (s *AccSuite) TestTeamDataSource() {
	rn := "data.rollbar_team.test"

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					log.Debug().Msg("Testing data source rollbar_team with invalid team name")
				},
				Config:      s.configDataSourceTeamNotFound(),
				ExpectError: regexp.MustCompile("no team with the name"),
			},
			{
				PreConfig: func() {
					log.Debug().Msg("Testing data source rollbar_team")
				},
				Config: s.configDataSourceTeam(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", s.randName),
					resource.TestCheckResourceAttrSet(rn, "id"),
					resource.TestCheckResourceAttrSet(rn, "account_id"),
					resource.TestCheckResourceAttrSet(rn, "access_level"),
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

		data "rollbar_team" "test" {
			name = "%s"
			depends_on = [rollbar_team.test]
		}
	`
	return fmt.Sprintf(tmpl, s.randName, s.randName)
}

func (s *AccSuite) configDataSourceTeamNotFound() string {
	// language=hcl
	tmpl := `
		data "rollbar_team" "test" {
			name = "%s"
		}
	`
	return fmt.Sprintf(tmpl, s.randName)
}
