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
	teamIDs := d.Get("teams").([]int)
	l := log.With().
		Str("email", email).
		Ints("teamIDs", teamIDs).
		Logger()
	l.Info().Msg("Creating resource rollbar_user")
	u := userStateID{
		Email: email,
	}
	d.SetId(u.String())
	return resourceUserCreateOrUpdate(ctx, d, meta)
}

// resourceUserCreateOrUpdate does the heavy lifting of assigning and/or
// inviting user to specified groups, and removing user from groups no longer
// specified.
func resourceUserCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.RollbarApiClient)
	u, err := userStateIDFromString(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	teamIDs := d.Get("teams").([]int)
	l := log.With().
		Interface("userStateID", u).
		Ints("teamIDs", teamIDs).
		Logger()

	// Check if a Rollbar user ID is known in TF state; or if not, whether a
	// Rollbar user already exists.
	var userID int
	if u.UserID != 0 {
		userID = u.UserID
	} else {
		userID, err = c.UserIdFromEmail(u.Email)
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
			u.UserID = userID
			l := l.With().Int("userID", userID).Logger()
			l.Debug().Msg("User already exists")
		}
	}
	d.SetId(u.String()) // In case this email now has a user ID that wasn't known before.

	// Teams to which this user SHOULD belong
	teamsDesired := make(map[int]bool)
	for _, id := range teamIDs {
		teamsDesired[id] = false
	}

	// Teams to which this user currently belongs
	teamsCurrent := make(map[int]bool)
	if userID != 0 { // If user doesn't exist, they don't belong to any teams
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
	// Join those teams
	for teamID, join := range teamsToJoin {
		l = l.With().Int("teamID", teamID).Logger()
		if !join {
			continue
		}
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
			inv, err := c.CreateInvitation(teamID, u.Email)
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
	// Leave those teams
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

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	usi, err := userStateIDFromString(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	l := log.With().
		Interface("userStateID", usi).
		Logger()
	l.Debug().Msg("Reading resource user")
	c := meta.(*client.RollbarApiClient)
	//es := errSetter{d: d}

	// If user ID did not exist last time, check to see if it exists now.  I.e.,
	// has the invited email registered to become a Rollbar user.
	if usi.UserID == 0 {
		userID, err := c.UserIdFromEmail(usi.Email)
		switch err {
		case client.ErrNotFound: // Do nothing
		case nil:
			// Newly registered Rollbar user
			usi.UserID = userID
			d.SetId(usi.String())
		default: // Other errors
			return diag.FromErr(err)
		}
	}

	// If we still don't have a UserID, we will look at invitations to get team memberships.
	if usi.UserID == 0 {
		invitations, err := c.FindInvitations(usi.Email)
		switch err {
		case client.ErrNotFound:
		}
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Get("email").(string)
	teamIDs := d.Get("teams").([]int)
	l := log.With().
		Str("email", email).
		Ints("teamIDs", teamIDs).
		Logger()
	l.Info().Msg("Creating resource rollbar_user")
	return resourceUserCreateOrUpdate(ctx, d, meta)
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := fmt.Errorf("not yet implemented")
	log.Err(err).Send()
	return diag.FromErr(err)

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

// userStateID represents the data used to identify a Rollbar user in Terraform
// state. As the prospective user progresses through the stages from invitation
// through becoming a registered user - so do the prospective user's API-side
// identifier change from Email to UserID.
type userStateID struct {
	UserID int
	Email  string
}

func (u userStateID) String() string {
	return fmt.Sprintf("%d:%s", u.UserID, u.Email)
}

func userStateIDFromString(s string) (u userStateID, err error) {
	components := strings.Split(s, ":")
	if len(components) != 2 {
		err = fmt.Errorf("invalid user state ID string: %s", s)
		log.Err(err).Send()
		return
	}
	userID, err := strconv.Atoi(components[0])
	if err != nil {
		log.Err(err).Send()
		return
	}
	email := components[1]
	u = userStateID{
		UserID: userID,
		Email:  email,
	}
	return
}
