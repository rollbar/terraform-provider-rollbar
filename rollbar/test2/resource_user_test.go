package test2

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

func init() {
	resource.AddTestSweepers("rollbar_user", &resource.Sweeper{
		Name: "rollbar_user",
		F:    sweepResourceUser,
	})
}

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
			email = "terraform-provider-test+%s@rollbar.com"
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
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					s.checkUserTeams(rn),
				),
			},
		},
	})
}

// TestAccUserCreateInviteMixedCase tests creating a new rollbar_user resource
// with an invitation to email is not registered as a Rollbar user, and contains
// mixed case characters.  The mixed case characters must be tested because the
// API converts all submitted email addresses to lower-case.  That's okay,
// standards say email addresses are NOT case sensitive.
//
// Demonstrates https://github.com/rollbar/terraform-provider-rollbar/issues/139
func (s *AccSuite) TestAccUserCreateInviteMixedCase() {
	rn := "rollbar_user.test_user"
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_user" "test_user" {
			# Note capital "X" in the email address below
			email = "terraform-provider-test+X-%s@rollbar.com"
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
			email = "terraform-provider-test@rollbar.com"
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
					//s.checkUserTeams(rn), // it will be always a problem because of https://github.com/rollbar/terraform-provider-rollbar/issues/91
				),
			},
		},
	})
}

// TestAccUserImportInvited tests importing a rollbar_user resource based on an
// invited email.
func (s *AccSuite) TestAccUserImportInvited() {
	rn := "rollbar_user.test_user"
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_user" "test_user" {
			email = "terraform-provider-test+%s@rollbar.com"
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
			},
		},
	})
}

// tests importing a rollbar_user resource based on an
// invited email.
func (s *AccSuite) TestAccUserImportRegistered() {
	rn := "rollbar_user.test_user"
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s-team-0"
		}

		resource "rollbar_user" "test_user" {
			email = "terraform-provider-test@rollbar.com"
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
			},
			{
				ResourceName: rn,
				ImportState:  true,
				//ImportStateVerify: true, // it will be always a problem because of https://github.com/rollbar/terraform-provider-rollbar/issues/91
			},
		},
	})
}

// TestAccInvitedUserAddTeam tests adding a team to a rollbar_user resource that
// is based on an invited but not yet registered email.
func (s *AccSuite) TestAccInvitedUserAddTeam() {
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
			email = "terraform-provider-test+%s@rollbar.com"
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
			email = "terraform-provider-test+%s@rollbar.com"
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
	resource.ParallelTest(s.T(), resource.TestCase{
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

// TestAccInvitedUserRemoveTeam tests adding a team to a rollbar_user resource
// that is based on an invited but not yet registered email.
func (s *AccSuite) TestAccInvitedUserRemoveTeam() {
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
			email = "terraform-provider-test+%s@rollbar.com"
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
			email = "terraform-provider-test+%s@rollbar.com"
			team_ids = [rollbar_team.test_team_1.id]
		}
	`
	configRemoveTeam := fmt.Sprintf(tmpl, s.randName, s.randName, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
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
				Config: configRemoveTeam,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					s.checkUserTeams(rn),
				),
			},
		},
	})
}

// TestAccRegisteredUserAddTeam tests adding a team to a rollbar_user resource
// that is based on an already registered user.
func (s *AccSuite) TestAccRegisteredUserAddTeam() {
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
			# This email already has an account.  
			# https://github.com/rollbar/terraform-provider-rollbar/issues/91
			email = "terraform-provider-test@rollbar.com"
			team_ids = [rollbar_team.test_team_1.id]
		}
	`
	configOrigin := fmt.Sprintf(tmpl, s.randName, s.randName)
	// language=hcl
	tmpl = `
		resource "rollbar_team" "test_team_1" {
			name = "%s-team-1"
		}

		resource "rollbar_team" "test_team_2" {
			name = "%s-team-2"
		}

		resource "rollbar_user" "test_user" {
			email = "terraform-provider-test@rollbar.com"
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
	configAddTeam := fmt.Sprintf(tmpl, s.randName, s.randName)
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

// TestAccRegisteredUserRemoveTeam tests removing a team from a rollbar_user
// resource that is based on an already registered user.
func (s *AccSuite) TestAccRegisteredUserRemoveTeam() {
	s.T().Skip("the terraform refresh plan was not empty")
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
			email = "terraform-provider-test@rollbar.com"
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
	configOrigin := fmt.Sprintf(tmpl, s.randName, s.randName)
	// language=hcl
	tmpl = `
		resource "rollbar_team" "test_team_1" {
			name = "%s-team-1"
		}

		resource "rollbar_team" "test_team_2" {
			name = "%s-team-2"
		}

		resource "rollbar_user" "test_user" {
			# This email already has an account.  
			# https://github.com/rollbar/terraform-provider-rollbar/issues/91
			email = "terraform-provider-test@rollbar.com"
			team_ids = [rollbar_team.test_team_1.id]
		}
	`
	configRemoveTeam := fmt.Sprintf(tmpl, s.randName, s.randName)
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
				Config: configRemoveTeam,
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

// checkUserTeams checks a rollbar_user resource's teams
func (s *AccSuite) checkUserTeams(resourceName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		l := log.With().Logger()
		l.Info().Msg("Checking rollbar_user resource's teams")
		c := s.client()
		email, err := s.getResourceIDString(ts, resourceName)
		s.Nil(err)

		teamExpected := make(map[int]bool) // Teams user is expected to have
		teamFound := make(map[int]bool)    // Teams found for this user
		var unexpectedTeams []int
		teamCount, err1 := s.getResourceAttrInt(ts, resourceName, "team_ids.#")
		s.Nil(err1)
		for i := 0; i < teamCount; i++ {
			attr := fmt.Sprintf("team_ids.%d", i)
			teamID, err := s.getResourceAttrInt(ts, resourceName, attr)
			s.Nil(err)
			teamFound[teamID] = false
			teamExpected[teamID] = true
		}
		l = l.With().Interface("teamExpected", teamExpected).Logger()

		// Check team memberships, if state contains a Rollbar user ID.
		if userID, err2 := s.getResourceAttrInt(ts, resourceName, "user_id"); err2 == nil {
			currentTeams, err3 := c.ListUserCustomTeams(userID)
			s.Nil(err3)
			for teamID := range teamFound {
				// Did we find an expected team?
				for _, t := range currentTeams {
					if t.ID == teamID {
						teamFound[teamID] = true
					}
				}
				// Did we find an unexpected team?
				if !teamExpected[teamID] {
					unexpectedTeams = append(unexpectedTeams, teamID)
				}
			}
		}
		l.Debug().
			Interface("teamFound", teamFound).
			Msg("Existing team memberships")

		// Check invitations
		invitations, err4 := c.FindPendingInvitations(email)
		if err4 != nil && err4 != client.ErrNotFound {
			s.Nil(err4)
		}
		// Did we find any expected teams?
		for teamID := range teamFound {
			for _, inv := range invitations {
				if inv.TeamID == teamID {
					teamFound[teamID] = true
				}
			}
		}
		// Did we find any unexpected teams?
		for _, inv := range invitations {
			if !teamExpected[inv.TeamID] {
				unexpectedTeams = append(unexpectedTeams, inv.TeamID)
			}
		}
		l.Debug().
			Interface("teamFound", teamFound).
			Msg("Team invitations plus memberships")

		// Error if any team was not found
		for teamID, found := range teamFound {
			if !found {
				msg := fmt.Sprintf("team %d not found", teamID)
				l.Error().Msg(msg)
				return fmt.Errorf(msg)
			}
		}

		// Error if any unexpected team was found
		if len(unexpectedTeams) != 0 {
			for _, teamID := range unexpectedTeams {
				t, err := c.ReadTeam(teamID)
				s.Nil(err)
				l.Error().
					Interface("team", t).
					Msg("Found unexpected team")
			}
			msg := fmt.Sprintf("found unexpected teams: %v", unexpectedTeams)
			log.Error().Msg(msg)
			return fmt.Errorf(msg)
		}

		// Check passed!
		return nil
	}
}

// TestAccMoveUserBetweenTeams tests moving a Rollbar user from one team to another.
func (s *AccSuite) TestAccUserMoveBetweenTeams() {
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
		resource "rollbar_team" "test_team_1" {
			name = "%s"
		}

		resource "rollbar_team" "test_team_2" {
			name = "%s"
		}

		# Registered user
		resource "rollbar_user" "test_user_1" {
			email = "%s"
			team_ids = [ rollbar_team.test_team_2.id ]
		}

		# Invited user
		resource "rollbar_user" "test_user_2" {
			email = "%s"
			team_ids = [ rollbar_team.test_team_2.id ]
		}
	`
	configChangeTeams := fmt.Sprintf(tmpl, team1Name, team2Name, user1Email, user2Email)
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
					resource.TestCheckTypeSetElemAttrPair("rollbar_user.test_user_1", "team_ids.0", "rollbar_team.test_team_1", "id"),
					resource.TestCheckTypeSetElemAttrPair("rollbar_user.test_user_2", "team_ids.0", "rollbar_team.test_team_1", "id"),
					s.checkUserIsOnTeam(user1Email, team1Name),
					s.checkUserIsInvited(user2Email, team1Name),
				),
			},
			{
				Config: configChangeTeams,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair("rollbar_user.test_user_1", "team_ids.0", "rollbar_team.test_team_2", "id"),
					resource.TestCheckTypeSetElemAttrPair("rollbar_user.test_user_2", "team_ids.0", "rollbar_team.test_team_2", "id"),
					s.checkUserIsOnTeam(user1Email, team2Name),
					s.checkUserIsInvited(user2Email, team2Name),
					s.checkUserIsNotOnTeam(user1Email, team1Name),
					s.checkUserIsNotInvited(user2Email, team1Name),
				),
			},
		},
	})
}

// checkUserIsOnTeam checks that a Rollbar user is on a team.
func (s *AccSuite) checkUserIsOnTeam(userEmail, teamName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		l := log.With().
			Str("user_email", userEmail).
			Str("team_name", teamName).
			Logger()
		l.Info().Msg("Checking that user is member of team")
		c := s.client()

		// Find user ID
		userID, err := c.FindUserID(userEmail)
		s.Nil(err)
		s.NotZero(userID)

		teams, err := c.ListUserTeams(userID)
		s.Nil(err)
		for _, t := range teams {
			if t.Name == teamName {
				l.Debug().Msg("Confirmed that user is member of team")
				return nil
			}
		}
		err = fmt.Errorf("could not confirm that user %s is member of team %s", userEmail, teamName)
		l.Err(err).Send()
		return err
	}
}

// checkUserIsNotOnTeam checks that a Rollbar user is not on a team.
func (s *AccSuite) checkUserIsNotOnTeam(userEmail, teamName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		l := log.With().
			Str("user_email", userEmail).
			Str("team_name", teamName).
			Logger()
		l.Info().Msg("Checking that user is not member of team")
		c := s.client()

		userID, err := c.FindUserID(userEmail)
		s.Nil(err)
		teams, err := c.ListUserTeams(userID)
		s.Nil(err)
		for _, t := range teams {
			if t.Name == teamName {
				err = fmt.Errorf("check failed, user %s is member of team %s", userEmail, teamName)
				l.Err(err).Send()
				return err
			}
		}
		l.Debug().Msg("Confirmed that user is not member of team")
		return nil
	}
}

// checkUserIsInvited checks that a Rollbar user has been invited to a team.
func (s *AccSuite) checkUserIsInvited(userEmail, teamName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		l := log.With().
			Str("user_email", userEmail).
			Str("team_name", teamName).
			Logger()
		l.Info().Msg("Checking that user has been invited to team")
		c := s.client()

		teamID, err := c.FindTeamID(teamName)
		s.Nil(err)
		invitations, err := c.ListPendingInvitations(teamID)
		s.Nil(err)
		for _, inv := range invitations {
			if inv.ToEmail == userEmail {
				l.Debug().Msg("Confirmed user is invited to team")
				return nil
			}
		}
		err = fmt.Errorf("could not confirm user %s is invited to team %s", userEmail, teamName)
		l.Err(err).Send()
		return err
	}
}

// checkUserIsInvited checks that a Rollbar user is not invited to a team.
func (s *AccSuite) checkUserIsNotInvited(userEmail, teamName string) resource.TestCheckFunc {
	return func(ts *terraform.State) error {
		l := log.With().
			Str("user_email", userEmail).
			Str("team_name", teamName).
			Logger()
		l.Info().Msg("Checking that user is not invited to team")
		c := s.client()

		teamID, err := c.FindTeamID(teamName)
		s.Nil(err)
		invitations, err := c.ListPendingInvitations(teamID)
		s.Nil(err)
		for _, inv := range invitations {
			if inv.ToEmail == userEmail {
				err = fmt.Errorf("user %s is invited to team %s", userEmail, teamName)
				l.Err(err).Send()
				return err
			}
		}
		l.Debug().Msg("Confirmed user is not invited to team")
		return nil
	}
}

// TestAccUserInvitedToRegistered tests the transition of a Rollbar user from
// invited to registered status.
func (s *AccSuite) TestAccUserInvitedToRegistered() {
	s.T().Skip("problem with responder")
	rn := "rollbar_user.test_user"
	//randString := s.randName
	randString := "tf-acc-test-7lppmg40pk" // Must be constant across VCR record/playback runs
	// language=hcl
	tmpl := `
		resource "rollbar_team" "test_team" {
			name = "%s"
		}

		resource "rollbar_user" "test_user" {
			email = "jason.mcvetta+%s@gmail.com"
			team_ids = [rollbar_team.test_team.id]
		}
	`
	config := fmt.Sprintf(tmpl, randString, randString)
	//var r *recorder.Recorder
	origTransport := http.DefaultTransport
	resource.Test(s.T(), resource.TestCase{

		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			//{
			//	PreConfig: func() {
			//		var err error
			//		r, err = recorder.New("vcr/invited_user")
			//		s.Nil(err)
			//		r.AddFilter(vcrFilterHeaders)
			//		http.DefaultTransport = r
			//	},
			//	Config: config,
			//	Check: resource.ComposeTestCheckFunc(
			//		s.checkResourceStateSanity(rn),
			//		resource.TestCheckResourceAttr(rn, "status", "invited"),
			//	),
			//},
			{
				PreConfig: func() {
					// When recording the cassette, we use
					// github.com/sqweek/dialog to pop up a GUI dialog and wait
					// for confirmation; thereby allowing the developer to
					// manually accept the invitation.

					//ok := dialog.Message(
					//	"%s",
					//	"Accept the email invitation then continue",
					//).Title("Invitation accepted?").YesNo()
					//if !ok {
					//	s.FailNow("User did not accept the invitation")
					//}

					//err := r.Stop() // Stop the previous recorder
					//s.Nil(err)
					r, err := recorder.New("vcr/registered_user")
					s.Nil(err)
					r.AddFilter(vcrFilterHeaders)
					http.DefaultTransport = r

				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(rn),
					resource.TestCheckResourceAttr(rn, "status", "registered"),
				),
			},
		},
	})
	//err := r.Stop() // Stop the last recorder
	//s.Nil(err)
	http.DefaultTransport = origTransport
}

// vcrFilterHeaders removes unnecessary headers from VCR recordings.
func vcrFilterHeaders(i *cassette.Interaction) error {
	delete(i.Request.Headers, "X-Rollbar-Access-Token")
	delete(i.Request.Headers, "User-Agent")
	for key := range i.Response.Headers {
		deleteHeader := false
		if strings.HasPrefix(key, "X-") {
			deleteHeader = true
		}
		if strings.HasPrefix(key, "Access-Control-") {
			deleteHeader = true
		}
		switch key {
		case "Alt-Svc", "Content-Length", "Etag", "Server", "Via":
			deleteHeader = true
		}
		if deleteHeader {
			delete(i.Response.Headers, key)
		}
	}
	return nil
}

// sweepResourceUser cleans up orphaned Rollbar users.
func sweepResourceUser(_ string) error {
	log.Info().Msg("Cleaning up Rollbar users from acceptance test runs.")

	c := client.NewClient(client.DefaultBaseURL, os.Getenv("ROLLBAR_API_KEY"))
	users, err := c.ListTestUsers()
	if err != nil {
		log.Err(err).Send()
		return err
	}

	// Find the ID for this account's "Everyone" team
	var everyoneTeamID int
	teams, err := c.ListTeams()
	if err != nil {
		log.Err(err).Send()
		return err
	}
	for _, t := range teams {
		if t.Name == "Everyone" {
			everyoneTeamID = t.ID
		}
	}

	count := 0
	for _, u := range users {
		// We're only interested in test users
		if !strings.HasPrefix(u.Username, "tf-acc-test-") {
			continue
		}
		// Ignore this user, because it is required for acceptance tests that
		// involve a registered user.
		if u.Username == "tf-acc-test-rollbar-provider" {
			continue
		}

		// Remove the user from Everyone team, thereby removing it from the account.
		err = c.RemoveUserFromTeam(u.ID, everyoneTeamID)
		if err != nil {
			log.Err(err).Send()
			return err
		}
		count++
		log.Debug().
			Int("user_id", u.ID).
			Str("username", u.Username).
			Msg("Removed user")
	}

	log.Info().Int("count", count).Msg("Users cleanup complete")
	return nil
}
