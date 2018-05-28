package rollbar

import (
	"fmt"
	"github.com/babbel/rollbar-go/rollbar"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserImporter,
		},

		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	email := d.Get("email").(string)
	teamID := d.Get("team_id").(int)

	log.Printf("[INFO] Inviting user with email: %s", email)
	client := **m.(**rollbar.Client)
	resp, err := client.InviteUser(teamID, email)
	if err != nil {
		return err
	}
	// Use the email as an id.
	id := resp.Result.ToEmail

	d.SetId(id)
	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	var invites []string
	invited := false
	userPresent := false
	email := d.Id()
	teamID := d.Get("team_id").(int)
	client := **m.(**rollbar.Client)

	listInvites, err := client.ListInvites(teamID)

	if err != nil {
		return err
	}

	listUsers, err := client.ListUsers()

	if err != nil {
		return err
	}
	// This logic is needed so that we can connect the user was invited.
	// Check if there's an invite for the user
	for _, invite := range listInvites.Result {
		// Find the corresponding invite with the provided email.
		if invite.ToEmail == email {
			// Append all the invites into a slice.
			invites := append(invites, invite.Status)
			// Get the last invite (they are usually sequential).
			lastInv := invites[len(invites)-1]
			// If the invitation is pending that means that the user is invited.
			if lastInv == "pending" {
				invited = true
			}
		}
	}
	// Check if the user is present in the team.
	for _, user := range listUsers.Result.Users {
		if user.Email == email {
			userPresent = true
		}
	}

	if userPresent == false {
		if invited == false {
			d.SetId("")
			return nil
		}
	}

	d.Set("email", email)
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	email := d.Id()
	teamID := d.Get("team_id").(int)
	client := **m.(**rollbar.Client)

	client.RemoveUserTeam(email, teamID)

	d.SetId("")
	return nil
}

func resourceUserImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Import needs to be done with 2 values email and team id which will be split.
	sParts := strings.Split(d.Id(), ":")

	if len(sParts) != 2 {
		return nil, fmt.Errorf("Invalid ID specified. Supplied ID must be written as <email>:<team_id>")
	}

	teamIDInt, err := strconv.Atoi(sParts[1])

	if err != nil {
		return nil, err
	}

	d.Set("team_id", teamIDInt)
	d.SetId(sParts[0])

	return []*schema.ResourceData{d}, nil
}
