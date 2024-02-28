package test1

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccResourceTeamUser_createInvited tests creating, importing and destroying a new rollbar_team_user
// resource with an invited user.
func (s *AccSuite) TestAccResourceTeamUser_createInvited() {
	rn := "rollbar_team_user.test_team_user"
	email := fmt.Sprintf("terraform-provider-test+%s@rollbar.com", s.randName)
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_team_user" "test_team_user" {
			team_id = rollbar_team.test_team.id
			email = "%s"
		}
	`
	config := fmt.Sprintf(tmpl, s.randName, email)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttrSet(rn, "team_id"),
					resource.TestCheckResourceAttr(rn, "email", email),
					resource.TestCheckResourceAttrSet(rn, "invite_id"),
					resource.TestCheckNoResourceAttr(rn, "user_id"),
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

// TestAccResourceTeamUser_createRegistered tests creating, importing and destroying a new rollbar_team_user
// resource with a registered user.
func (s *AccSuite) TestAccResourceTeamUser_createRegistered() {
	rn := "rollbar_team_user.test_team_user"
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_team_user" "test_team_user" {
			team_id = rollbar_team.test_team.id
			# This email already has an account.  
			# https://github.com/rollbar/terraform-provider-rollbar/issues/91
			email = "terraform-provider-test@rollbar.com"
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
					resource.TestCheckResourceAttrSet(rn, "team_id"),
					resource.TestCheckResourceAttr(rn, "email", "terraform-provider-test@rollbar.com"),
					resource.TestCheckResourceAttr(rn, "invite_id", "0"),
					resource.TestCheckResourceAttrSet(rn, "user_id"),
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
