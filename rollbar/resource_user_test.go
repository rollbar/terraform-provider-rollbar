package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rs/zerolog/log"
	"regexp"
)

// TestAccUserCreateInvite tests creating a new rollbar_user resource with an
// invitation to email is not registered as a Rollbar user.
func (s *AccSuite) TestAccUserCreateInvite() {
	rn := "rollbar_user.test_user"
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_user" "test_user" {
			email = "jason.mcvetta+%s@gmail.com"
			team_ids = [rollbar_team.test_team.id]
		}
	`
	config := fmt.Sprintf(tmpl, s.randName, s.randName)
	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					s.checkUserTeams(rn),
				),
			},
		},
	})
}

// TestAccUserCreateAssign tests creating a new rollbar_user resource by
// assigning an already-registered Rollbar user to the team.
// FIXME: https://github.com/rollbar/terraform-provider-rollbar/issues/91
func (s *AccSuite) TestAccUserCreateAssign() {
	rn := "rollbar_user.test_user"
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_user" "test_user" {
			# This email already has an account.  
			# https://github.com/rollbar/terraform-provider-rollbar/issues/91
			email = "jason.mcvetta+tf-acc-test-rollbar-provider@gmail.com"
			team_ids = [rollbar_team.test_team.id]
		}
	`
	config := fmt.Sprintf(tmpl, s.randName)
	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					s.checkUserTeams(rn),
				),
			},
		},
	})
}

// TestAccUserImport tests importing a rollbar_user resource.
func (s *AccSuite) DontTestAccUserImport() {
	rn := "rollbar_user.test_user"
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_user" "test_user" {
			email = "jason.mcvetta+%s@gmail.com"
			team_ids = [rollbar_team.test_team.id]
		}
	`
	config := fmt.Sprintf(tmpl, s.randName, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: func(iss []*terraform.InstanceState) error {
					log.Warn().Interface("iss", iss).Send()
					return nil
				},
			},
		},
	})
}

// TestAccUserAddTeam tests adding a team to an existing rollbar_user resource.
func (s *AccSuite) TestAccUserAddTeam() {
	rn := "rollbar_user.test_user"
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team_1" {
			name = "%s-team-1"
		}

		resource "rollbar_team" "test_team_2" {
			name = "%s-team-2"
		}

		resource "rollbar_user" "test_user" {
			email = "jason.mcvetta+rollbar-tf-acc-test-%s@gmail.com"
			team_ids = [rollbar_team.test_team_1.id]
		}
	`
	configOrigin := fmt.Sprintf(tmpl, s.randName, s.randName, s.randName)
	// language=hcl
	tmpl = `
		resource "rollbar_team" "test_team_1" {
			name = "%s-team-1"
		}

		resource "rollbar_team" "test_team_2" {
			name = "%s-team-2"
		}

		resource "rollbar_user" "test_user" {
			email = "jason.mcvetta+rollbar-tf-acc-test-%s@gmail.com"
			team_ids = [
				rollbar_team.test_team_1.id,
				rollbar_team.test_team_2.id,
			]
			depends_on = [
				rollbar_team.test_team_1, 
				rollbar_team.test_team_2
			]
		}
	`
	configAddTeam := fmt.Sprintf(tmpl, s.randName, s.randName, s.randName)
	resource.Test(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: configOrigin,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
				),
			},
			{
				Config: configAddTeam,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					s.checkUserTeams(rn),
				),
			},
		},
	})
}

// TestAccUserInvalidConfig tests invalid config when trying to create a
// rollbar_user resource.
func (s *AccSuite) TestAccUserInvalidConfig() {
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_user" "test_user" {
			email = ""
			team_ids = [rollbar_team.test_team.id]
		}
	`
	config := fmt.Sprintf(tmpl, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Email must be supplied"),
			},
		},
	})
}

// checkUserTeams checks team memberships on a rollbar_user resource.
func (s *AccSuite) checkUserTeams(resourceName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		c := s.client()
		email, err := s.getResourceIDString(ts, resourceName)
		s.Nil(err)

		teamFound := make(map[int]bool)
		teamCount, err := s.getResourceAttrInt(ts, resourceName, "team_ids.#")
		s.Nil(err)
		for i := 0; i < teamCount; i++ {
			attr := fmt.Sprintf("team_ids.%d", i)
			teamID, err := s.getResourceAttrInt(ts, resourceName, attr)
			s.Nil(err)
			teamFound[teamID] = false
		}

		// If state contains a Rollbar user ID, check the users teams
		if userID, err := s.getResourceAttrInt(ts, resourceName, "user_id"); err == nil {
			foundTeamIDs, err := c.ListUserTeams(userID)
			s.Nil(err)
			for teamID, _ := range teamFound {
				for _, id := range foundTeamIDs {
					if id == teamID {
						teamFound[teamID] = true
					}
				}
			}
		}

		// If we are expecting team IDs that were not found, check the user's
		// invitations.
		remaining := 0
		for _, found := range teamFound {
			if !found {
				remaining++
			}
		}
		if remaining > 0 {
			invitations, err := c.FindInvitations(email)
			s.Nil(err)
			for teamID, _ := range teamFound {
				for _, inv := range invitations {
					if inv.TeamID == teamID {
						teamFound[teamID] = true
					}
				}
			}
		}

		// Error if any team was not found
		for teamID, found := range teamFound {
			if !found {
				msg := fmt.Sprintf("team %d not found", teamID)
				log.Debug().Msg(msg)
				return fmt.Errorf(msg)
			}
		}

		// Test passed!
		return nil
	}
}
