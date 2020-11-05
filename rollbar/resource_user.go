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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"team_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			// Computed
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceUserCreate creates a new Rollbar user resource.
func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Get("email").(string)
	teamIDs := getUserTeamIDs(d)
	l := log.With().
		Str("email", email).
		Ints("teamIDs", teamIDs).
		Logger()
	l.Info().Msg("Creating rollbar_user resource")
	d.SetId(email)
	mustSet(d, "status", "invited")
	return resourceUserCreateOrUpdate(ctx, d, meta)
}

// resourceUserCreateOrUpdate does the heavy lifting of assigning and/or
// inviting user to specified groups, and removing user from groups no longer
// specified.
func resourceUserCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.RollbarApiClient)
	email := d.Get("email").(string)
	teamIDs := getUserTeamIDs(d)
	l := log.With().
		Str("email", email).
		Ints("expected_team_ids", teamIDs).
		Logger()
	l.Debug().Msg("Creating or updating rollbar_user resource")

	// Check if a Rollbar user exists for this email
	userID, err := c.FindUserID(email)
	switch err {
	case nil:
		l = l.With().Int("user_id", userID).Logger()
		l.Debug().Int("id", userID).Msg("Found existing user")
		mustSet(d, "user_id", userID)
	case client.ErrNotFound:
		l.Debug().Int("id", userID).Msg("Existing user not found")
	default: // Actual error
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	l = l.With().Int("user_id", userID).Logger()

	// Teams to which this user SHOULD belong
	teamsExpected := make(map[int]bool)
	for _, id := range teamIDs {
		teamsExpected[id] = true
	}

	// Teams to which this user currently belongs
	teamsCurrent := make(map[int]bool)
	if userID != 0 { // If user doesn't exist, they don't belong to any teams
		teams, err := c.ListUserCustomTeams(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		for _, t := range teams {
			teamsCurrent[t.ID] = true
		}
	}
	l.Debug().Interface("current_teams", teamsCurrent).Msg("Current teams")

	// Teams to which this user should be added
	var teamsToJoin []int
	for id, _ := range teamsExpected {
		if !teamsCurrent[id] {
			teamsToJoin = append(teamsToJoin, id)
		}
	}
	l.Debug().Interface("teams_to_join", teamsToJoin).Msg("Teams to join")
	// Add user to those teams
	for _, teamID := range teamsToJoin {
		l = l.With().Int("teamID", teamID).Logger()
		// If user already exists we can assign to teams without invitation.  If
		// user does not already exist we must send an invitation.
		if userID != 0 {
			err = c.AssignUserToTeam(teamID, userID)
			if err != nil {
				l.Err(err).Send()
				return diag.FromErr(err)
			}
			l.Debug().Msg("Assigned user to team")
		} else {
			inv, err := c.CreateInvitation(teamID, email)
			if err != nil {
				l.Err(err).Send()
				return diag.FromErr(err)
			}
			l.Debug().
				Int("inviteID", inv.ID).
				Msg("Invited user to team")
		}
	}

	// Teams from which this user should be removed
	var teamsToRemove []int
	for id, _ := range teamsCurrent {
		if !teamsExpected[id] {
			teamsToRemove = append(teamsToRemove, id)
		}
	}
	l.Debug().Ints("teams_to_remove", teamsToRemove).Msg("Teams to leave")
	// Remove user from those teams
	if userID != 0 {
		l.Debug().Msg("Removing user from teams")
		for _, teamID := range teamsToRemove {
			err := c.RemoveUserFromTeam(userID, teamID)
			if err != nil {
				l.Err(err).Send()
				return diag.FromErr(err)
			}
		}
	}

	// Invitations which should be cancelled
	var invitationsToCancel []int
	invitations, err := c.FindPendingInvitations(email)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	for _, inv := range invitations {
		if !teamsExpected[inv.TeamID] {
			invitationsToCancel = append(invitationsToCancel, inv.ID)
		}
	}
	l.Debug().
		Int("count", len(invitationsToCancel)).
		Msg("Canceling invitations")
	// Cancel those invitations
	for _, invitationID := range invitationsToCancel {
		err := c.CancelInvitation(invitationID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	}

	d.SetId(email)
	l.Debug().Msg("Successfully created or updated rollbar_user resource")
	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Id()
	userID := d.Get("user_id").(int)
	l := log.With().
		Str("email", email).
		Int("userID", userID).
		Logger()
	l.Info().Msg("Reading rollbar_user resource")
	c := meta.(*client.RollbarApiClient)
	var err error

	// If user ID is not in state, try to query it from Rollbar
	if userID == 0 {
		userID, err = c.FindUserID(email)
		if err != client.ErrNotFound && err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	}

	// If Rollbar user already exists, list user's teamIDs
	var teamIDs []int
	if userID != 0 {
		mustSet(d, "status", "registered")
		u, err := c.ReadUser(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		mustSet(d, "username", u.Username)
		teams, err := c.ListUserCustomTeams(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		for _, t := range teams {
			teamIDs = append(teamIDs, t.ID)
		}
	}

	// Add pending invitations to team IDs
	invitations, err := c.FindPendingInvitations(email)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	for _, inv := range invitations {
		teamIDs = append(teamIDs, inv.TeamID)
	}

	mustSet(d, "team_ids", teamIDs)

	l.Debug().Msg("Successfully read rollbar_user resource")
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Get("email").(string)
	teamIDs := getUserTeamIDs(d)
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
	c := meta.(*client.RollbarApiClient)

	// Try to get user ID
	userID := d.Get("user_id").(int)
	if userID == 0 {
		userID, _ = c.FindUserID(email)
	}

	// If user ID is known, remove user from teams
	if userID != 0 {
		teams, err := c.ListUserCustomTeams(userID)
		if err != nil && err != client.ErrNotFound {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		for _, t := range teams {
			err := c.RemoveUserFromTeam(userID, t.ID)
			if err != nil {
				l.Err(err).Send()
				return diag.FromErr(err)
			}
		}
	}

	// Cancel user's invitations
	invitations, err := c.FindPendingInvitations(email)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	for _, inv := range invitations {
		err := c.CancelInvitation(inv.ID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	l.Debug().Msg("Successfully deleted rollbar_user resource")
	return nil
}

// getUserTeamIDs gets the team IDs for a rollbar_user resource.
func getUserTeamIDs(d *schema.ResourceData) []int {
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
	c := meta.(*client.RollbarApiClient)

	invitations, err := c.FindInvitations(email)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	mustSet(d, "status", "invited")
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
