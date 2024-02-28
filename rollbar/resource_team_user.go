/*
 * Copyright (c) 2022 Rollbar, Inc.
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
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

func resourceTeamUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamUserCreate,
		ReadContext:   resourceTeamUserRead,
		UpdateContext: nil, // resourceTeamUserUpdate,
		DeleteContext: resourceTeamUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
			"team_id": {
				Description: "ID of the team to which this user belongs",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"email": {
				Description: "The user's email address",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			// Computed
			"status": {
				Description: "Status of the user. Either `invited` or `registered`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description: "The ID of the user",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"invite_id": {
				Description: "Invitation ID if status is `invited`",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func teamUserID(teamID int, email string) string {
	return fmt.Sprintf("%d%s%s", teamID, ComplexImportSeparator, email)
}

func teamUserFromID(id string) (teamID int, s string, err error) {

	if !strings.Contains(id, ComplexImportSeparator) {
		return 0, s, fmt.Errorf("resource ID missing delimiter (%s)", ComplexImportSeparator)
	}
	l := log.With().
		Str("id", id).
		Logger()
	values := strings.SplitN(id, ComplexImportSeparator, 2)
	l = l.With().Strs("values", values).Logger()
	l.Debug().Msg("Converting userID to int")
	teamID, err = strconv.Atoi(values[0])
	if err != nil {
		return 0, "", fmt.Errorf("unable to parse team ID")
	}
	return teamID, values[1], nil
}

func resourceTeamUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	teamID := d.Get("team_id").(int)
	email := d.Get("email").(string)
	l := log.With().
		Str("email", email).
		Int("team_id", teamID).
		Logger()
	l.Info().Msg("Creating rollbar_team_user resource")

	// Check if a Rollbar user exists for this email
	c.SetHeaderResource(rollbarTeamUser)
	userID, err := c.FindUserID(email)

	l = l.With().Int("user_id", userID).Logger()
	switch err {
	case nil: // User Found, assign them to the team
		l.Debug().Msg("Found existing user")
		mustSet(d, "user_id", userID)
		mustSet(d, "status", "registered")
		er := c.AssignUserToTeam(teamID, userID)
		if er != nil {
			l.Err(er).Msg("error assigning user to team")
			return diag.FromErr(er)
		}
		mustSet(d, "invite_id", 0)
		l.Debug().Msg("Assigned user to team")
	case client.ErrNotFound: // User not found, send an invitation
		l.Debug().Msg("Existing user not found")
		mustSet(d, "status", "invited")
		inv, er := c.CreateInvitation(teamID, email)
		if er != nil {
			l.Err(er).Msg("error assigning user to team")
			return diag.FromErr(er)
		}
		l.Debug().
			Int("inviteID", inv.ID).
			Msg("Invited user to team")
		mustSet(d, "invite_id", inv.ID)
	default: // Actual error
		l.Err(err).Send()
		return diag.FromErr(err)
	}

	d.SetId(teamUserID(teamID, email))
	l.Debug().Msg("Successfully created or updated rollbar_team_user resource")
	return resourceTeamUserRead(ctx, d, meta)
}

func resourceTeamUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	teamID, email, err := teamUserFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	userID := d.Get("user_id").(int)
	l := log.With().
		Str("email", email).
		Int("user_id", userID).
		Int("team_id", teamID).
		Logger()
	l.Info().Msg("Reading rollbar_team_user resource")
	c := meta.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	c.SetHeaderResource(rollbarTeamUser)

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
			mustSet(d, "user_id", userID)
			mustSet(d, "status", "registered")
		case client.ErrNotFound:
			l.Debug().Msg("No registered user found")
			mustSet(d, "status", "invited")
		default:
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	}

	if userID != 0 {
		// Check if user is assigned to the team
		assigned, err := c.IsUserAssignedToTeam(teamID, userID)
		if err != nil {
			l.Err(err).Msg("Error checking if user is assigned to team.")
			return diag.FromErr(err)
		}
		if assigned {
			mustSet(d, "team_id", teamID)
		} else {
			d.SetId("")
		}
		_ = d.Set("invite_id", nil)
	} else {
		// Check if user is invited to the team
		invitations, err := c.ListPendingInvitations(teamID)
		if err != nil {
			l.Err(err).Msg("Error checking if user has pending invitation.")
			return diag.FromErr(err)
		}
		var invite client.Invitation
		for _, i := range invitations {
			if i.ToEmail == email {
				invite = i
			}
		}
		mustSet(d, "invite_id", invite.ID)
	}
	// Ensure team_id and email are set, they may be missing when importing.
	mustSet(d, "team_id", teamID)
	mustSet(d, "email", email)

	l.Debug().Msg("Successfully read rollbar_user resource")
	return nil
}

func resourceTeamUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Id()
	teamID := d.Get("team_id").(int)
	l := log.With().
		Str("email", email).
		Int("team_id", teamID).
		Logger()
	l.Info().Msg("Deleting rollbar_team_user resource")
	c := meta.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	c.SetHeaderResource(rollbarTeamUser)

	userID := d.Get("user_id").(int)
	if userID == 0 {
		// Cancel invitation
		inviteID := d.Get("invite_id").(int)
		err := c.CancelInvitation(inviteID)
		if err != client.ErrNotFound {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	} else {
		// Remove user from team
		err := c.RemoveUserFromTeam(userID, teamID)
		if err != nil {
			if err != client.ErrNotFound {
				l.Err(err).Send()
				return diag.FromErr(err)
			}
		}
	}

	d.SetId("")

	l.Debug().Msg("Successfully deleted rollbar_team_user resource")
	return nil
}
