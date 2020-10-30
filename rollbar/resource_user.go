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
				Type:     schema.TypeList,
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
	teamIDs := getValueAsIntSlice(d, "team_ids")
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
	teamIDs := getValueAsIntSlice(d, "team_ids")
	l := log.With().
		Str("email", email).
		Ints("teamIDs", teamIDs).
		Logger()
	l.Debug().Msg("Creating or updating rollbar_user resource")

	// Check if a Rollbar user exists for this email
	userID, err := c.FindUserID(email)
	switch err {
	case nil:
		es.Set("user_id", userID)
	case client.ErrNotFound: // Do nothing
	default: // Actual error
		l.Err(err).Send()
		return diag.FromErr(err)
	}

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

	if es.err != nil {
		l.Err(es.err).Msg("Error setting state value")
		return diag.FromErr(err)
	}
	d.SetId(email)
	l.Info().Msg("Successfully created or updated rollbar_user resource")
	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Id()
	userID := d.Get("user_id").(int)
	l := log.With().
		Str("email", email).
		Int("userID", userID).
		Logger()
	l.Info().Msg("Reading user resource")
	c := meta.(*client.RollbarApiClient)
	es := errSetter{d: d}
	teamIDs := make(map[int]bool)
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
	if userID != 0 {
		es.Set("status", "registered")
		u, err := c.ReadUser(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		es.Set("username", u.Username)
		ids, err := c.ListUserTeams(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		for _, id := range ids {
			teamIDs[id] = true
		}
	}

	// Add pending invitations to team IDs
	invitations, err := c.FindInvitations(email)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, inv := range invitations {
		if inv.Status == "pending" {
			teamIDs[inv.TeamID] = true
		}
	}

	// Flatten the teamIDs map and set state
	var tids []int
	for id, _ := range teamIDs {
		tids = append(tids, id)
	}
	es.Set("team_ids", tids)

	if es.err != nil {
		l.Err(es.err).Msg("Error setting state")
		return diag.FromErr(err)
	}
	l.Info().Msg("Successfully read user resource")
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Get("email").(string)
	teamIDs := d.Get("teams").([]int)
	l := log.With().
		Str("email", email).
		Ints("teamIDs", teamIDs).
		Logger()
	l.Info().Msg("Updating resource rollbar_user")
	return resourceUserCreateOrUpdate(ctx, d, meta)
}

func resourceUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	email := d.Id()
	l := log.With().
		Str("email", email).
		Logger()
	l.Info().Msg("Deleting user resource")
	c := meta.(*client.RollbarApiClient)

	// If user ID is known, remove user from teams
	userID, _ := c.FindUserID(email)
	if userID != 0 {
		teamIDs, err := c.ListUserTeams(userID)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		for _, teamID := range teamIDs {
			err := c.RemoveUserFromTeam(userID, teamID)
			if err != nil {
				l.Err(err).Send()
				return diag.FromErr(err)
			}
		}
	}

	// Cancel user's invitations
	invitations, err := c.FindInvitations(email)
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

	l.Info().Msg("Successfully deleted user resource")
	return nil
}

/*
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

*/
