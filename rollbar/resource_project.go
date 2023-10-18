/*
 * Copyright (c) 2023 Rollbar, Inc.
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
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		DeleteContext: resourceProjectDelete,
		UpdateContext: resourceProjectUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Description: "The human readable name for the project",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional
			"team_ids": {
				Description: "IDs of the teams assigned to the project",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			// Computed
			"account_id": {
				Description: "ID of the account that owns the project",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"date_created": {
				Description: "Date the project was created",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"date_modified": {
				Description: "Date the project was last modified",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"status": {
				Description: "Status of the project",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"timezone": {
				Description: "Timezone for the project",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"time_format": {
				Description: "Time format for the project",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	timezone := d.Get("timezone").(string)
	timeFormat := d.Get("time_format").(string)

	if timezone == "" {
		timezone = timeZoneDefault
	}
	if timeFormat == "" {
		timeFormat = timeformatDefault
	}

	l := log.With().Str("name", name).Logger()
	l.Info().Msg("Creating new Rollbar project resource")

	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]

	client.Mutex.Lock()
	setResourceHeader(rollbarProject, c)
	p, err := c.CreateProject(name, timezone, timeFormat)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	l.Debug().Interface("project", p).Msg("CreateProject() result")
	projectID := p.ID
	l = l.With().Int("project_id", projectID).Logger()
	d.SetId(strconv.Itoa(projectID))

	// A set of four default access tokens are automagically created by Rollbar
	// when creating a new project.  However we only want access tokens that are
	// explicitly created and managed by Terraform.  Therefore we delete the
	// default tokens for our new project.
	expectedTokenNames := map[string]bool{
		"read":             true,
		"write":            true,
		"post_client_item": true,
		"post_server_item": true,
	}
	client.Mutex.Lock()
	tokens, err := c.ListProjectAccessTokens(projectID)
	client.Mutex.Unlock()
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	for _, t := range tokens {
		// Sanity check
		expected := expectedTokenNames[t.Name]
		if !expected {
			err = fmt.Errorf("unexpected token name in default tokens")
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		// Deletion
		client.Mutex.Lock()
		err = c.DeleteProjectAccessToken(projectID, t.AccessToken)
		client.Mutex.Unlock()
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		l.Debug().
			Str("name", t.Name).
			Msg("Successfully deleted a default access token")
	}

	// Team assignments
	teamIDsSet := d.Get("team_ids").(*schema.Set)
	for _, teamIDiface := range teamIDsSet.List() {
		teamID := teamIDiface.(int)
		l = l.With().Int("team_id", teamID).Logger()
		client.Mutex.Lock()
		err = c.AssignTeamToProject(teamID, projectID)
		client.Mutex.Unlock()
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	}

	l.Debug().Msg("Successfully created Rollbar project resource")
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectID := mustGetID(d)
	l := log.With().
		Int("projectID", projectID).
		Logger()
	l.Info().Msg("Reading Rollbar project resource")

	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]

	client.Mutex.Lock()
	setResourceHeader(rollbarProject, c)
	proj, err := c.ReadProject(projectID)
	client.Mutex.Unlock()

	if err == client.ErrNotFound {
		l.Debug().Msg("Project not found on Rollbar - removing from state")
		d.SetId("")
		return nil
	}
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}

	var mProj map[string]interface{}
	mustDecodeMapStructure(proj, &mProj)
	for k, v := range mProj {
		if k == "id" {
			continue
		}
		if k == settingsData {
			continue
		}
		mustSet(d, k, v)
	}

	for k, v := range mProj[settingsData].(map[string]interface{}) {
		if k == "timezone" && v == timeZoneDefault {
			continue
		}
		if k == "time_format" && v == timeformatDefault {
			continue
		}
		mustSet(d, k, v)
	}

	client.Mutex.Lock()
	teamIDs, err := c.FindProjectTeamIDs(projectID)
	client.Mutex.Unlock()
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	mustSet(d, "team_ids", teamIDs)

	d.SetId(strconv.Itoa(proj.ID))
	l.Debug().Msg("Successfully read Rollbar project resource from the API")
	return nil
}

// resourceProjectUpdate handles update for a `rollbar_project` resource.
func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	teamIDs := getTeamIDs(d)
	projectID := mustGetID(d)
	name := d.Get("name").(string)

	timezone := d.Get("timezone").(string)
	timeFormat := d.Get("time_format").(string)
	if timezone == "" {
		timezone = timeZoneDefault
	}
	if timeFormat == "" {
		timeFormat = timeformatDefault
	}
	l := log.With().
		Int("project_id", projectID).
		Ints("team_ids", teamIDs).
		Logger()
	l.Debug().Msg("Updating rollbar_project resource")
	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]

	client.Mutex.Lock()
	setResourceHeader(rollbarProject, c)
	err := c.UpdateProjectTeams(projectID, teamIDs)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Msg("Error updating rollbar_project resource")
		return diag.FromErr(err)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarProject, c)
	p, err := c.UpdateProject(projectID, name, timezone, timeFormat)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Msg("Error updating rollbar_project resource")
		return diag.FromErr(err)
	}

	if p.ID != projectID {
		err = errors.New("IDs are not equal")
		l.Err(err).Send()
		d.SetId("") // removing from the state
		return diag.FromErr(err)
	}

	l.Debug().Msg("Successfully updated rollbar_project resource")
	return resourceProjectRead(ctx, d, m)
}

// resourceProjectDelete handles delete for a `rollbar_project` resource.
func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectID := mustGetID(d)
	l := log.With().
		Int("projectID", projectID).
		Logger()
	l.Info().Msg("Deleting rollbar_project resource")
	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]

	client.Mutex.Lock()
	setResourceHeader(rollbarProject, c)
	err := c.DeleteProject(projectID)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Msg("Error deleting rollbar_project resource")
		return diag.FromErr(err)
	}
	l.Debug().Msg("Successfully deleted rollbar_project resource")
	return nil
}
