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

		// TODO: Import functionality
		//Importer: &schema.ResourceImporter{
		//	StateContext: schema.ImportStatePassthroughContext,
		//},

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
	err := d.Set("status", "invited")
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	return resourceUserCreateOrUpdate(ctx, d, meta)
}

// resourceUserCreateOrUpdate does the heavy lifting of assigning and/or
// inviting user to specified groups, and removing user from groups no longer
// specified.
func resourceUserCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.RollbarApiClient)
	es := errSetter{d: d}
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
		es.Set("user_id", userID)
	case client.ErrNotFound:
		l.Debug().Int("id", userID).Msg("Existing user not found")
	default: // Actual error
		l.Err(err).Send()
		return diag.FromErr(err)
	}

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
	teamsToJoin := make(map[int]bool)
	for id, _ := range teamsExpected {
		if !teamsCurrent[id] {
			teamsToJoin[id] = true
		}
	}
	l.Debug().Interface("teams_to_join", teamsToJoin).Msg("Teams to join")
	// Join those teams
	for teamID, _ := range teamsToJoin {
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
		teamsToJoin[teamID] = false // Task complete
	}

	// FIXME: problem is here!  Maybe.  Maybe problem is use of Assign rather than Invite for all team adds

	// Teams from which this user should be removed
	teamsToLeave := make(map[int]bool)
	for id, _ := range teamsCurrent {
		if !teamsExpected[id] {
			teamsToLeave[id] = true
		}
	}
	l.Debug().Interface("teams", teamsToLeave).Msg("Teams to leave")
	// Leave those teams
	for teamID, leave := range teamsToLeave {
		if !leave {
			continue
		}
		err := c.RemoveUserFromTeam(userID, teamID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		teamsToLeave[teamID] = false // Task complete
	}

	if es.err != nil {
		l.Err(es.err).Msg("Error setting state value")
		return diag.FromErr(err)
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
	es := errSetter{d: d}
	//teamIDs := make(map[int]bool)
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
		es.Set("status", "registered")
		u, err := c.ReadUser(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		es.Set("username", u.Username)
		teams, err := c.ListUserCustomTeams(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		for _, t := range teams {
			//teamIDs[t.ID] = true
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
		//teamIDs[inv.TeamID] = true
		teamIDs = append(teamIDs, inv.TeamID)
	}

	//// Flatten the teamIDs map and set state
	//var tids []int
	//for id, _ := range teamIDs {
	//	tids = append(tids, id)
	//}
	//es.Set("team_ids", tids)
	es.Set("team_ids", teamIDs)

	if es.err != nil {
		l.Err(es.err).Msg("Error setting state")
		return diag.FromErr(err)
	}
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
