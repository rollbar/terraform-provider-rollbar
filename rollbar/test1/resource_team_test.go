package test1

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("rollbar_team", &resource.Sweeper{
		Name: "rollbar_team",
		F:    sweepResourceTeam,
	})
}

// TestAccTeamInvalidName tests failure when team name is an empty string.
func (s *AccSuite) TestAccTeamInvalidName() {
	// language=hcl
	config := `
		resource "rollbar_team" "test" {
			name = ""
		}
	`
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Invalid name - failure expected
			{
				Config:      config,
				ExpectError: regexp.MustCompile("name cannot be blank"),
			},
		},
	})
}

// TestAccTeamCreate tests creating a Rollbar team.
func (s *AccSuite) TestAccTeamCreate() {
	rn := "rollbar_team.test"
	teamName := fmt.Sprintf("%s-team-0", s.randName)
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
			name = "%s"
		}
	`
	config := fmt.Sprintf(tmpl, teamName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", teamName),
					s.checkTeam(rn, teamName, "standard"),
				),
			},
		},
	})
}

// TestAccTeamUpdateAccessLevel tests updating the access level on a Rollbar
// team.
func (s *AccSuite) TestAccTeamUpdateAccessLevel() {
	rn := "rollbar_team.test"
	teamName := fmt.Sprintf("%s-team-0", s.randName)
	// language=hcl
	tmpl1 := `
		resource "rollbar_team" "test" {
			name = "%s"
		}
	`
	config1 := fmt.Sprintf(tmpl1, teamName)
	// language=hcl
	tmpl2 := `
		resource "rollbar_team" "test" {
			name = "%s"
			access_level = "light"
		}
	`
	config2 := fmt.Sprintf(tmpl2, teamName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Initial create
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", teamName),
					s.checkTeam(rn, teamName, "standard"),
				),
			},
			// Update team access level
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", teamName),
					s.checkTeam(rn, teamName, "light"),
				),
			},
		},
	})
}

// TestAccTeamUpdateName tests updating the name of a Rollbar team.
func (s *AccSuite) TestAccTeamUpdateName() {
	rn := "rollbar_team.test"
	teamName1 := fmt.Sprintf("%s-team-1", s.randName)
	teamName2 := fmt.Sprintf("%s-team-2", s.randName)
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
			name = "%s"
		}
	`
	config1 := fmt.Sprintf(tmpl, teamName1)
	config2 := fmt.Sprintf(tmpl, teamName2)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Initial create
			{
				Config: config1,
			},
			// Update team name
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", teamName2),
				),
			},
		},
	})
}

// TestAccTeamImport tests importing a Rollbar team into Terraform.
func (s *AccSuite) TestAccTeamImport() {
	rn := "rollbar_team.test"
	teamName1 := fmt.Sprintf("%s-team-1", s.randName)
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
			name = "%s"
		}
	`
	config1 := fmt.Sprintf(tmpl, teamName1)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Initial create
			{
				Config: config1,
			},
			// Import the team
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTeamDeleteOnAPIBeforeApply tests creating a Rollbar team with
// Terraform, then deleting the team via API, before again applying Terraform
// config.
// FIXME: This code used to pass reliably, but no longer does.   Why?
//
//	https://github.com/rollbar/terraform-provider-rollbar/issues/154
func (s *AccSuite) TestAccTeamDeleteOnAPIBeforeApply() {
	rn := "rollbar_team.test"
	teamName1 := fmt.Sprintf("%s-team-1", s.randName)
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
			name = "%s"
		}
	`
	config1 := fmt.Sprintf(tmpl, teamName1)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Initial create
			{
				Config: config1,
			},
			// Before running Terraform, delete the team on Rollbar but not in local state
			{
				PreConfig: func() {
					c := client.NewClient(client.DefaultBaseURL, os.Getenv("ROLLBAR_API_KEY"))
					teams, err := c.ListCustomTeams()
					s.Nil(err)
					for _, t := range teams {
						if t.Name == teamName1 {
							err = c.DeleteTeam(t.ID)
							s.Nil(err)
							log.Info().
								Str("team_name", teamName1).
								Int("team_id", t.ID).
								Msg("Deleted team from API before re-applying Terraform config")
						}
					}
				},
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					s.checkTeam(rn, teamName1, "standard"),
					resource.TestCheckResourceAttr(rn, "name", teamName1),
				),
			},
		},
	})
}

// checkTeam checks that the newly created team exists and has correct
// attributes.
func (s *AccSuite) checkTeam(rn, teamName, accessLevel string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		s.Nil(err)
		c := s.provider.Meta().(map[string]*client.RollbarAPIClient)[schemaKeyToken]
		t, err := c.ReadTeam(id)
		s.Nil(err)
		s.Equal(teamName, t.Name, "team name from API does not match team name in Terraform config")
		s.Equal(accessLevel, t.AccessLevel)
		return nil
	}
}

// sweepResourceTeam cleans up orphaned Rollbar teams.
func sweepResourceTeam(_ string) error {
	log.Info().Msg("Cleaning up Rollbar teams from acceptance test runs.")

	c := client.NewClient(client.DefaultBaseURL, os.Getenv("ROLLBAR_API_KEY"))
	teams, err := c.ListCustomTeams()
	if err != nil {
		log.Err(err).Send()
		return err
	}

	count := 0
	for _, t := range teams {
		l := log.With().
			Str("name", t.Name).
			Int("id", t.ID).
			Logger()
		if strings.HasPrefix(t.Name, "tf-acc-test-") {
			err = c.DeleteTeam(t.ID)
			if err != nil {
				l.Err(err).Send()
				return err
			}
			count++
			l.Debug().Msg("Deleted team")
		}
	}

	log.Info().Int("count", count).Msg("Teams cleanup complete")
	return nil
}

// TestAccTeamDeleteTeamWithUsers tests deleting a Rollbar team that has a
// non-zero count of users.
func (s *AccSuite) TestAccTeamDeleteTeamWithUsers() {
	s.T().Skip("Root object was present, but now absent")
	team1Name := fmt.Sprintf("%s-team-1", s.randName)
	team2Name := fmt.Sprintf("%s-team-2", s.randName)
	user1Email := "terraform-provider-test@rollbar.com"
	user2Email := fmt.Sprintf("terraform-provider-test+%s@rollbar.com", s.randName)
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team_1" {
			name = "%s"
		}

		resource "rollbar_team" "test_team_2" {
			name = "%s"
		}

		# Registered user
		resource "rollbar_user" "test_user_1" {
			email = "%s"
			team_ids = [ rollbar_team.test_team_1.id ]
		}

		# Invited user
		resource "rollbar_user" "test_user_2" {
			email = "%s"
			team_ids = [ rollbar_team.test_team_1.id ]
		}
	`
	configOrigin := fmt.Sprintf(tmpl, team1Name, team2Name, user1Email, user2Email)
	// language=hcl
	tmpl = `
		resource "rollbar_team" "test_team_2" {
			name = "%s"
		}
	`
	configRemoveTeam := fmt.Sprintf(tmpl, team2Name)
	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: configOrigin,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rollbar_team.test_team_1", "name", team1Name),
					resource.TestCheckResourceAttr("rollbar_team.test_team_2", "name", team2Name),
					resource.TestCheckTypeSetElemAttrPair("rollbar_user.test_user_1", "team_ids.*", "rollbar_team.test_team_1", "id"),
					resource.TestCheckTypeSetElemAttrPair("rollbar_user.test_user_2", "team_ids.*", "rollbar_team.test_team_1", "id"),
				),
			},
			{
				Config: configRemoveTeam,
				Check: resource.ComposeTestCheckFunc(
					s.checkTeamIsDeleted(team1Name),
					resource.TestCheckResourceAttr("rollbar_team.test_team_2", "name", team2Name),
				),
			},
		},
	})
}

func (s *AccSuite) checkTeamIsDeleted(teamName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		l := log.With().Str("team_name", teamName).Logger()
		l.Info().Msg("Checking that team is deleted")
		c := s.client()
		teams, err := c.ListCustomTeams()
		s.Nil(err)
		for _, t := range teams {
			if t.Name == teamName {
				err := fmt.Errorf("team was NOT deleted: %s", teamName)
				l.Err(err).Send()
				return err
			}
		}
		l.Debug().Msg("Confirmed that team is deleted")
		return nil
	}
}

// TestTeamValidateAccessLevel tests validation of argument `access_level` on a
// `rollbar_team` resource.
func TestTeamValidateAccessLevel(t *testing.T) {
	p := cty.Path{} // placeholder
	validAccessLevels := []string{
		"standard",
		"light",
		"view",
	}
	for _, level := range validAccessLevels {
		d := resourceTeamValidateAccessLevel(level, p)
		assert.Nil(t, d)
	}
	d := resourceTeamValidateAccessLevel("invalid-level", p)
	assert.Len(t, d, 1)
	assert.IsType(t, diag.Diagnostic{}, d[0])
}

func resourceTeamValidateAccessLevel(v interface{}, p cty.Path) diag.Diagnostics {
	s := v.(string)
	switch s {
	case "standard", "light", "view":
		return nil
	default:
		summary := fmt.Sprintf(`Invalid access_level: %q`, s)
		d := diag.Diagnostic{
			Severity:      diag.Error,
			AttributePath: p,
			Summary:       summary,
			Detail:        `Must be "standard", "light", or "view"`,
		}
		return diag.Diagnostics{d}
	}
}
