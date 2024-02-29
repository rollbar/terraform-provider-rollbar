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
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

func init() {
	resource.AddTestSweepers("rollbar_project", &resource.Sweeper{
		Name: "rollbar_project",
		F:    sweepResourceProject,
	})
}

// TestAccProject tests creation and deletion of a Rollbar project.
func (s *AccSuite) TestAccProject() {
	rn := "rollbar_project.foo"

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: s.configResourceProject(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", s.randName),
					s.checkProjectExists(rn, s.randName),
					s.checkProjectInProjectList(rn),
				),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTeamAssignProject tests assigning a team to a project
func (s *AccSuite) TestAccTeamAssignProject() {
	projectResourceName := "rollbar_project.test_project"
	teamName := fmt.Sprintf("%s-team-0", s.randName)
	projectName := s.randName
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s"
		}

		resource "rollbar_project" "test_project" {
			name = "%s"
			team_ids = [rollbar_team.test_team.id]
		}
	`
	config := fmt.Sprintf(tmpl, teamName, projectName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(projectResourceName),
					s.checkProjectTeams(projectResourceName),
				),
			},
		},
	})
}

// TestAccTeamAddProject tests adding a team to a project.
func (s *AccSuite) TestAccProjectAddTeam() {
	team1ResourceName := "rollbar_team.test_team_1"
	team1Name := fmt.Sprintf("%s-team-1", s.randName)
	team2ResourceName := "rollbar_team.test_team_2"
	team2Name := fmt.Sprintf("%s-team-2", s.randName)
	projectResourceName := "rollbar_project.test_project"
	projectName := s.randName

	// language=hcl
	tmpl1 := `
		resource "rollbar_team" "test_team_1" {
			name = "%s"
		}

		resource "rollbar_project" "test_project" {
			name = "%s"
			team_ids = [rollbar_team.test_team_1.id]
		}
	`
	config1 := fmt.Sprintf(tmpl1, team1Name, projectName)

	// language=hcl
	tmpl2 := `
		resource "rollbar_team" "test_team_1" {
			name = "%s"
		}

		resource "rollbar_team" "test_team_2" {
			name = "%s"
		}

		resource "rollbar_project" "test_project" {
			name = "%s"
			team_ids = [
				rollbar_team.test_team_1.id,
				rollbar_team.test_team_2.id,
			]
		}
	`
	config2 := fmt.Sprintf(tmpl2, team1Name, team2Name, projectName)

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair(projectResourceName, "team_ids.0", team1ResourceName, "id"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair(projectResourceName, "team_ids.0", team1ResourceName, "id"),
					resource.TestCheckTypeSetElemAttrPair(projectResourceName, "team_ids.1", team2ResourceName, "id"),
					s.checkProjectTeams(projectResourceName),
				),
			},
		},
	})
}

// TestAccProjectRemoveTeam tests removing a team from a project.
func (s *AccSuite) TestAccProjectRemoveTeam() {
	team1ResourceName := "rollbar_team.test_team_1"
	team1Name := fmt.Sprintf("%s-team-1", s.randName)
	team2ResourceName := "rollbar_team.test_team_2"
	team2Name := fmt.Sprintf("%s-team-2", s.randName)
	projectResourceName := "rollbar_project.test_project"
	projectName := s.randName

	// language=hcl
	tmpl1 := `
		resource "rollbar_team" "test_team_1" {
			name = "%s"
		}

		resource "rollbar_team" "test_team_2" {
			name = "%s"
		}

		resource "rollbar_project" "test_project" {
			name = "%s"
			team_ids = [
				rollbar_team.test_team_1.id,
				rollbar_team.test_team_2.id,
			]
		}
	`
	config1 := fmt.Sprintf(tmpl1, team1Name, team2Name, projectName)

	// language=hcl
	tmpl2 := `
		resource "rollbar_team" "test_team_1" {
			name = "%s"
		}

		resource "rollbar_project" "test_project" {
			name = "%s"
			team_ids = [rollbar_team.test_team_1.id]
		}
	`
	config2 := fmt.Sprintf(tmpl2, team1Name, projectName)

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair(projectResourceName, "team_ids.*", team1ResourceName, "id"),
					resource.TestCheckTypeSetElemAttrPair(projectResourceName, "team_ids.*", team2ResourceName, "id"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair(projectResourceName, "team_ids.0", team1ResourceName, "id"),
					s.checkProjectTeams(projectResourceName),
				),
			},
		},
	})
}

// TestAccProjectDeleteOnAPIBeforeApply tests creating a Rollbar project with
// Terraform; then deleting the project via API before re-applying Terraform
// configuration.
func (s *AccSuite) TestAccProjectDeleteOnAPIBeforeApply() {
	rn := "rollbar_project.test"
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
			name = "%s"
		}
	`
	config := fmt.Sprintf(tmpl, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Initial create
			{
				Config: config,
			},
			// Before running Terraform, delete the project on Rollbar but not in local state
			{
				PreConfig: func() {
					c := client.NewClient(client.DefaultBaseURL, os.Getenv("ROLLBAR_API_KEY"))
					projects, err := c.ListProjects()
					s.Nil(err)
					for _, p := range projects {
						if p.Name == s.randName {
							err = c.DeleteProject(p.ID)
							s.Nil(err)
							log.Info().
								Str("project_name", s.randName).
								Int("project_id", p.ID).
								Msg("Deleted project from API before re-applying Terraform config")
						}
					}
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					s.checkProjectExists(rn, s.randName),
					s.checkProjectInProjectList(rn),
					resource.TestCheckResourceAttr(rn, "name", s.randName),
				),
			},
		},
	})
}

/*
 * Convenience functions
 */

func (s *AccSuite) configResourceProject() string {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "foo" {
		  name         = "%s"
		}
	`
	return fmt.Sprintf(tmpl, s.randName)
}

// checkProjectExists tests that the newly created project exists
func (s *AccSuite) checkProjectExists(rn string, name string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		s.Nil(err)
		c := s.provider.Meta().(map[string]*client.RollbarAPIClient)[schemaKeyToken]
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
		c := s.provider.Meta().(map[string]*client.RollbarAPIClient)[schemaKeyToken]
		projList, err := c.ListProjects()
		s.Nil(err)
		found := false
		for _, proj := range projList {
			if proj.ID == id {
				found = true
			}
		}
		s.True(found, "Project not found in project list")
		return nil
	}
}

// sweepResourceProject cleans up orphaned projects created by failed acceptance
// test runs.
func sweepResourceProject(_ string) error {
	log.Info().Msg("Cleaning up Rollbar projects from acceptance test runs.")

	c := client.NewClient(client.DefaultBaseURL, os.Getenv("ROLLBAR_API_KEY"))
	projects, err := c.ListProjects()
	if err != nil {
		log.Err(err).Send()
		return err
	}

	counter := 0
	for _, p := range projects {
		l := log.With().
			Str("name", p.Name).
			Int("id", p.ID).
			Logger()
		if strings.HasPrefix(p.Name, "tf-acc-test-") {
			err = c.DeleteProject(p.ID)
			if err != nil {
				l.Err(err).Send()
				return err
			}
			counter++
			l.Debug().Msg("Deleted project")
		}
	}

	log.Info().Int("count", counter).Msg("Projects cleanup complete")
	return nil
}

// checkProjectTeams checks that the project is assigned to the correct teams.
func (s *AccSuite) checkProjectTeams(projectResourceName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		l := log.With().Logger()
		l.Info().Msg("Checking rollbar_project resource's teams")
		projectID, err := s.getResourceIDInt(ts, projectResourceName)
		s.Nil(err)
		expected, err := s.getResourceAttrIntSlice(ts, projectResourceName, "team_ids")
		s.Nil(err)
		actual, err := s.client().FindProjectTeamIDs(projectID)
		s.Nil(err)
		s.ElementsMatch(expected, actual)
		return nil
	}
}
