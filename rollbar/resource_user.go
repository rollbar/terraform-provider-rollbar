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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
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
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			// Computed
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.RollbarApiClient)
	email := d.Get("email").(string)
	teamIDs := d.Get("teams").([]int)
	l := log.With().
		Str("email", email).
		Ints("teamIDs", teamIDs).
		Logger()
	l.Info().Msg("Creating resource rollbar_user")

	// If a user already exists, we assign the user to each team.
	var userExists bool // User already exists?
	userID, err := c.UserIdFromEmail(email)
	switch err {
	default:
		// Error
		l.Err(err).Send()
		return diag.FromErr(err)
	case client.ErrNotFound:
		// User does not yet exist
		l.Debug().Msg("User does not already exist")
	case nil:
		// User already exists
		l := l.With().Int("userID", userID).Logger()
		l.Debug().Msg("User already exists")
		userExists = true
	}

	// Teams to which this user SHOULD belong
	teamsDesired := make(map[int]bool)
	for _, id := range teamIDs {
		teamsDesired[id] = false
	}

	// Teams to which this user currently does belong
	teamsCurrent := make(map[int]bool)
	if userExists { // If user doesn't exist, they don't belong to any teams
		ut, err := c.ListUserTeams(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		for _, t := range ut {
			teamsCurrent[t] = true
		}
	}

	// Teams to which this user should be added
	teamsToJoin := make(map[int]bool)
	for id, _ := range teamsDesired {
		if !teamsCurrent[id] {
			teamsToJoin[id] = true
		}
	}
	for teamID, join := range teamsToJoin {
		l = l.With().Int("teamID", teamID).Logger()
		if !join {
			continue
		}
		// If user already exists we can assign to teams without invitation.  If
		// user does not already exist we must send an invitation.
		if userExists {
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
				Msg("Assigned user to team")
		}
		teamsToJoin[teamID] = false // Task complete
	}

	// Teams from which this user should be removed
	teamsToLeave := make(map[int]bool)
	for id, _ := range teamsCurrent {
		if !teamsDesired[id] {
			teamsToLeave[id] = true
		}
	}

	for teamID, leave := range teamsToLeave {
		if !leave {
			continue
		}
		err := c.RemoveUserFromTeam(teamID, userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		teamsToLeave[teamID] = false // Task complete
	}

	return nil
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var invites []string
	invited := false
	userPresent := false
	email := d.Id()
	teamID := d.Get("team_id").(int)
	l := log.With().
		Str("email", email).
		Int("teamID", teamID).
		Logger()
	l.Debug().Msg("Reading resource user")

	c := meta.(*client.RollbarApiClient)
	listInvites, err := c.ListInvitations(teamID)
	if err != nil {
		l.Err(err).Msg("Error reading resource user")
		return diag.FromErr(err)
	}

	listUsers, err := c.ListUsers()
	if err != nil {
		if err != nil {
			l.Err(err).Msg("Error reading resource user")
			return diag.FromErr(err)
		}
	}

	// This logic is needed so that we can check if the the user already was invited.
	// Check if there's an active invite for the user or the user has already accepted the invite.
	for _, invite := range listInvites {
		// Find the corresponding invite with the provided email.
		if invite.ToEmail == email {
			// Append all the invites into a slice.
			invites := append(invites, invite.Status)

			// Get the last invite (they are sequential).
			lastInv := invites[len(invites)-1]

			// If the invitation is pending that means that the user is invited.
			if lastInv == "pending" {
				invited = true
			}
		}
	}
	// Check if the user is present in the team.
	for _, user := range listUsers {
		if user.Email == email {
			userPresent = true
		}
	}

	if !userPresent && !invited {
		d.SetId("")
		err := fmt.Errorf("no user or invitee found with the email %s", email)
		l.Err(err).Msg("Error reading resource user")
		return diag.FromErr(err)
	}

	err = d.Set("email", email)
	if err != nil {
		l.Err(err).Msg("Error reading resource user")
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Fatal().Msg("Not yet implemented")
	return nil
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Id()
	teamID := d.Get("team_id").(int)
	l := log.With().
		Str("email", email).
		Int("teamID", teamID).
		Logger()
	l.Debug().Msg("Deleting resource user")

	c := meta.(*client.RollbarApiClient)
	id, err := c.UserIdFromEmail(email)
	if err != nil {
		l.Err(err).Msg("Error deleting resource user")
		return diag.FromErr(err)
	}
	err = c.RemoveUserFromTeam(id, teamID)
	if err != nil {
		l.Err(err).Msg("Error deleting resource user")
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceUserImporter(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Import needs to be done with 2 values email and team id which will be split.
	sParts := strings.Split(d.Id(), ":")

	if len(sParts) != 2 {
		return nil, fmt.Errorf("invalid ID specified. Supplied ID must be written as <email>:<team_id>")
	}

	teamIDInt, err := strconv.Atoi(sParts[1])

	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamIDInt)
	d.SetId(sParts[0])

	return []*schema.ResourceData{d}, nil
}

func resourceUserStateId(userID int, email string) string {
	return fmt.Sprintf("%d:%s", userID, email)
}
