package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
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
	teamName0 := fmt.Sprintf("%s-team-0", s.projectName)
	teamName1 := fmt.Sprintf("%s-team-0", s.projectName)

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
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
					teams, err := c.ListTeams()
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
	teams, err := c.ListTeams()
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
