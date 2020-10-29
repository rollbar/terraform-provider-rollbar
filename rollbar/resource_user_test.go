package rollbar_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"strings"
)

// TestAccUser tests CRUD operations for a Rollbar user.
func (s *AccSuite) TestAccUser() {
	rn := "rollbar_user.test_user"
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			//// Invalid name - failure expected
			//{
			//	Config:      s.configResourceUserInvalidname(),
			//	ExpectError: regexp.MustCompile("name cannot be blank"),
			//},
			//
			// Initial create
			{
				Config: s.configResourceUser(),
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					s.checkUser(rn),
				),
			},
		},
	})
}

func (s *AccSuite) configResourceUser() string {
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_user" "test_user" {
			email = "jason.mcvetta+rollbar-tf-acc-test@gmail.com"
			team_ids = [rollbar_team.test_team.id]
		}
	`
	return fmt.Sprintf(tmpl, s.randName)
}

func (s *AccSuite) checkUser(resourceName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		c := s.client()
		email, err := s.getResourceIDString(ts, resourceName)
		s.Nil(err)
		teamIDsString, err := s.getResourceAttrString(ts, resourceName, "teams")
		s.Nil(err)

		teamFound := make(map[int]bool)
		components := strings.Split(teamIDsString, ",")
		for _, teamIDstring := range components {
			teamID, err := strconv.Atoi(teamIDstring)
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
				return fmt.Errorf("team %d not found", teamID)
			}
		}

		// Test passed!
		return nil
	}
}
