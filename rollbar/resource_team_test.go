package rollbar

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("rollbar_team", &resource.Sweeper{
		Name: "rollbar_team",
		F:    sweepResourceTeam,
	})
}

// TestAccTeam tests CRUD operations for a Rollbar team.
func (s *AccSuite) TestAccTeam() {
	rn := "rollbar_team.test"
	teamName0 := fmt.Sprintf("%s-team-0", s.randName)
	teamName1 := fmt.Sprintf("%s-team-0", s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Invalid name - failure expected
			{
				Config:      s.configResourceTeamInvalidname(),
				ExpectError: regexp.MustCompile("name cannot be blank"),
			},
			// Initial create
			{
				Config: s.configResourceTeam(teamName0),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", teamName0),
					s.checkTeam(rn, teamName0, "standard"),
				),
			},
			// Update team access level
			{
				Config: s.configResourceTeamUpdateAccessLevel(teamName0),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", teamName0),
					s.checkTeam(rn, teamName0, "light"),
				),
			},
			// Update team name
			{
				Config: s.configResourceTeamUpdateTeamName(teamName1),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", teamName1),
					s.checkTeam(rn, teamName1, "light"),
				),
			},
			// Before running Terraform, delete the team on Rollbar but not in local state
			{
				PreConfig: func() {
					c := client.NewClient(os.Getenv("ROLLBAR_API_KEY"))
					teams, err := c.ListCustomTeams()
					s.Nil(err)
					for _, t := range teams {
						if t.Name == teamName1 {
							err = c.DeleteTeam(t.ID)
							s.Nil(err)
						}
					}
				},
				Config: s.configResourceTeamUpdateTeamName(teamName1),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "name", teamName1),
					s.checkTeam(rn, teamName1, "light"),
				),
			},
			// Import a team
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func (s *AccSuite) configResourceTeamInvalidname() string {
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
			name = ""
		}
	`
	return tmpl
}

func (s *AccSuite) configResourceTeam(teamName string) string {
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
			name = "%s"
		}
	`
	return fmt.Sprintf(tmpl, teamName)
}

func (s *AccSuite) configResourceTeamUpdateAccessLevel(teamName string) string {
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
			name = "%s"
			access_level = "light"
		}
	`
	return fmt.Sprintf(tmpl, teamName)
}

func (s *AccSuite) configResourceTeamUpdateTeamName(teamName string) string {
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test" {
			name = "%s"
			access_level = "light"
		}
	`
	return fmt.Sprintf(tmpl, teamName)
}

// checkTeam checks that the newly created team exists and has correct
// attributes.
func (s *AccSuite) checkTeam(rn, teamName, accessLevel string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		id, err := s.getResourceIDInt(ts, rn)
		s.Nil(err)
		c := s.provider.Meta().(*client.RollbarApiClient)
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

	c := client.NewClient(os.Getenv("ROLLBAR_API_KEY"))
	teams, err := c.ListCustomTeams()
	if err != nil {
		log.Err(err).Send()
		return err
	}

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
			l.Info().Msg("Deleted team")
		}
	}

	log.Info().Msg("Teams cleanup complete")
	return nil
}

// TestAccUserRemoveTeamWithUsers tests deleting a Rollbar team that has a
// non-zero count of users.
func (s *AccSuite) TestAccDeleteTeamWithUsers() {
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
