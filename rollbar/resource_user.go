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

package rollbar

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},

		Schema: map[string]*schema.Schema{
			// Required
			"email": {
				Description: "The user's email address",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"team_ids": {
				Description: "IDs of the teams to which this user belongs",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			// Computed
			"username": {
				Description: "The user's username",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description: "The ID of the user",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"status": {
				Description: "Status of the user.  Either `invited` or `subscribed`",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// resourceUserCreate creates a new Rollbar user resource.
func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Get("email").(string)
	teamIDs := getTeamIDs(d)
	l := log.With().
		Str("email", email).
		Ints("teamIDs", teamIDs).
		Logger()
	l.Info().Msg("Creating rollbar_user resource")
	d.SetId(email)
	return resourceUserCreateOrUpdate(ctx, d, meta)
}

// resourceUserCreateOrUpdate does the heavy lifting of assigning and/or
// inviting user to specified groups, and removing user from groups no longer
// specified.
func resourceUserCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.RollbarAPIClient)
	email := d.Get("email").(string)
	teamIDs := getTeamIDs(d)
	l := log.With().
		Str("email", email).
		Ints("expected_team_ids", teamIDs).
		Logger()
	l.Debug().Msg("Creating or updating rollbar_user resource")

	// Check if a Rollbar user exists for this email
	userID, err := c.FindUserID(email)
	l = l.With().Int("user_id", userID).Logger()
	switch err {
	case nil:
		l.Debug().Msg("Found existing user")
		mustSet(d, "user_id", userID)
		mustSet(d, "status", "registered")
	case client.ErrNotFound:
		l.Debug().Msg("Existing user not found")
		mustSet(d, "status", "invited")
	default: // Actual error
		l.Err(err).Send()
		return diag.FromErr(err)
	}

	// Teams to which this user SHOULD belong
	teamsExpected := make(map[int]bool)
	for _, id := range teamIDs {
		teamsExpected[id] = true
	}

	teamsCurrent, err := resourceUserCurrentTeams(c, email, userID)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}

	err = resourceUserAddTeams(resourceUserAddRemoveTeamsArgs{
		client:        c,
		userID:        userID,
		email:         email,
		teamsExpected: teamsExpected,
		teamsCurrent:  teamsCurrent,
	})
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}

	err = resourceUserRemoveTeams(resourceUserAddRemoveTeamsArgs{
		client:        c,
		userID:        userID,
		email:         email,
		teamsExpected: teamsExpected,
		teamsCurrent:  teamsCurrent,
	})
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}

	d.SetId(email)
	l.Debug().Msg("Successfully created or updated rollbar_user resource")
	return resourceUserRead(ctx, d, meta)
}

// resourceUserAddRemoveTeamsArgs encapsulates the arguments to
// resourceUserAddTeams and resourceUserRemoveTeams.
type resourceUserAddRemoveTeamsArgs struct {
	client        *client.RollbarAPIClient
	userID        int
	email         string
	teamsExpected map[int]bool
	teamsCurrent  map[int]bool
}

// resourceUserAddTeams adds new team memberships to a Rollbar user, either by
// assigning a registered user to the team or by inviting an email address to
// the team.
func resourceUserAddTeams(args resourceUserAddRemoveTeamsArgs) error {
	l := log.With().
		Int("user_id", args.userID).
		Str("email", args.email).
		Interface("expected_teams", args.teamsExpected).
		Interface("current_teams", args.teamsCurrent).
		Logger()
	errMsg := "Error joining teams"

	// Teams to which this user should be added
	var teamsToJoin []int
	for id := range args.teamsExpected {
		if !args.teamsCurrent[id] {
			teamsToJoin = append(teamsToJoin, id)
		}
	}
	l.Debug().Interface("teams_to_join", teamsToJoin).Msg("Teams to join")

	// Add user to those teams
	for _, teamID := range teamsToJoin {
		l = l.With().Int("teamID", teamID).Logger()
		// If user already exists we can assign to teams without invitation.  If
		// user does not already exist we must send an invitation.
		if args.userID != 0 {
			err := args.client.AssignUserToTeam(teamID, args.userID)
			if err != nil {
				l.Err(err).Msg(errMsg)
				return err
			}
			l.Debug().Msg("Assigned user to team")
		} else {
			inv, err := args.client.CreateInvitation(teamID, args.email)
			if err != nil {
				l.Err(err).Msg(errMsg)
				return err
			}
			l.Debug().
				Int("inviteID", inv.ID).
				Msg("Invited user to team")
		}
	}
	return nil
}

// resourceUserRemoveTeams removes team memberships from a Rollbar user.
func resourceUserRemoveTeams(args resourceUserAddRemoveTeamsArgs) error {
	l := log.With().
		Int("user_id", args.userID).
		Str("email", args.email).
		Interface("expected_teams", args.teamsExpected).
		Interface("current_teams", args.teamsCurrent).
		Logger()
	errMsg := "Error removing user from team"

	// Teams from which this user should be removed
	teamsToLeave := make(map[int]bool)
	for id := range args.teamsCurrent {
		if !args.teamsExpected[id] {
			teamsToLeave[id] = true
		}
	}
	l.Debug().Interface("unwanted_teams", teamsToLeave).Msg("Unwanted teams")

	// Leave teams
	if args.userID != 0 {
		l.Debug().Msg("Removing registered user from teams")
		currentTeams, _ := args.client.ListUserCustomTeams(args.userID)
		for _, t := range currentTeams {
			if teamsToLeave[t.ID] {
				err := args.client.RemoveUserFromTeam(args.userID, t.ID)
				if err != nil {
					l.Err(err).Msg(errMsg)
					return err
				}
			}
		}
	}

	// Cancel invitations
	l.Debug().Msg("Canceling invitations")
	invitations, err := args.client.FindPendingInvitations(args.email)
	if err != nil && err != client.ErrNotFound {
		l.Err(err).Msg(errMsg)
		return err
	}
	for _, inv := range invitations {
		if teamsToLeave[inv.TeamID] {
			err := args.client.CancelInvitation(inv.ID)
			if err != nil {
				l.Err(err).Msg(errMsg)
				return err
			}
		}

	}

	return nil
}

// resourceUserCurrentTeams returns user's current team memberships.
func resourceUserCurrentTeams(c *client.RollbarAPIClient, email string, userID int) (currentTeams map[int]bool, err error) {
	l := log.With().
		Str("email", email).
		Int("user_id", userID).
		Logger()
	currentTeams = make(map[int]bool)

	// Registered user team memberships
	if userID != 0 {
		var teams []client.Team
		teams, err = c.ListUserTeams(userID)
		if err != nil && err != client.ErrNotFound {
			l.Err(err).Send()
			return
		}
		for _, t := range teams {
			currentTeams[t.ID] = true
		}
	}

	// Teams to which email has been invited
	var invitations []client.Invitation
	invitations, err = c.FindPendingInvitations(email)
	if err != nil && err != client.ErrNotFound {
		l.Err(err).Send()
		return
	}
	for _, inv := range invitations {
		currentTeams[inv.TeamID] = true
	}

	l.Debug().
		Interface("current_teams", currentTeams).
		Msg("Current teams")
	return currentTeams, nil
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Id()
	userID := d.Get("user_id").(int)
	l := log.With().
		Str("email", email).
		Int("userID", userID).
		Logger()
	l.Info().Msg("Reading rollbar_user resource")
	c := meta.(*client.RollbarAPIClient)
	var err error

	// If user ID is not in state, try to query it from Rollbar
	if userID == 0 {
		userID, err = c.FindUserID(email)
		switch err {
		case nil:
			l = log.With().
				Str("email", email).
				Int("userID", userID).
				Logger()
			l.Debug().Msg("Found registered user")
		case client.ErrNotFound:
			l.Debug().Msg("No registered user found")
		default:
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	}

	// If no user ID was found, user has been invited but not yet registered.
	if userID == 0 {
		mustSet(d, "status", "invited")
	} else {
		mustSet(d, "status", "registered")
	}

	currentTeams, err := resourceUserCurrentTeams(c, email, userID)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	var teamIDs []int
	for teamID := range currentTeams {
		teamIDs = append(teamIDs, teamID)
	}
	mustSet(d, "team_ids", teamIDs)

	l.Debug().Msg("Successfully read rollbar_user resource")
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Get("email").(string)
	teamIDs := getTeamIDs(d)
	l := log.With().
		Str("email", email).
		Ints("teamIDs", teamIDs).
		Logger()
	l.Info().Msg("Updating rollbar_user resource")
	return resourceUserCreateOrUpdate(ctx, d, meta)
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Id()
	l := log.With().
		Str("email", email).
		Logger()
	l.Info().Msg("Deleting rollbar_user resource")
	c := meta.(*client.RollbarAPIClient)

	// Try to get user ID
	userID := d.Get("user_id").(int)
	if userID == 0 {
		userID, _ = c.FindUserID(email)
	}

	teamsCurrent, err := resourceUserCurrentTeams(c, email, userID)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	teamsExpected := make(map[int]bool) // Empty
	err = resourceUserRemoveTeams(resourceUserAddRemoveTeamsArgs{
		client:        c,
		email:         email,
		userID:        userID,
		teamsCurrent:  teamsCurrent,
		teamsExpected: teamsExpected,
	})
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}

	d.SetId("")

	l.Debug().Msg("Successfully deleted rollbar_user resource")
	return nil
}

// getTeamIDs gets team IDs for a resource.
func getTeamIDs(d *schema.ResourceData) []int {
	set := d.Get("team_ids").(*schema.Set)
	teamIDs := make([]int, set.Len())
	for i, teamID := range set.List() {
		teamIDs[i] = teamID.(int)
	}
	return teamIDs
}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	email := d.Id()
	mustSet(d, "email", email)
	teamIDsSet := d.Get("team_ids").(*schema.Set)
	l := log.With().
		Str("email", email).
		Interface("team_ids", teamIDsSet.List()).
		Logger()
	l.Info().Msg("Importing rollbar_user resource")

	var teamIDs []int
	c := meta.(*client.RollbarAPIClient)

	invitations, err := c.FindInvitations(email)
	if err != nil && err != client.ErrNotFound {
		l.Err(err).Send()
		return nil, err
	}
	if len(invitations) > 0 {
		mustSet(d, "status", "invited")
	}
	for _, inv := range invitations {
		teamIDs = append(teamIDs, inv.TeamID)
	}

	userID, err := c.FindUserID(email)
	if err == nil {
		mustSet(d, "user_id", userID)
		mustSet(d, "status", "registered")
		teams, err := c.ListUserTeams(userID)
		if err != nil {
			l.Err(err).Send()
			return nil, err
		}
		for _, t := range teams {
			teamIDs = append(teamIDs, t.ID)
		}
	}

	mustSet(d, "team_ids", teamIDs)

	return []*schema.ResourceData{d}, nil
}
