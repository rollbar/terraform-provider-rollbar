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
	"github.com/mitchellh/mapstructure"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		DeleteContext: resourceProjectDelete,

		// Projects cannot be updated via API
		//UpdateContext: resourceProjectUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date_modified": {
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

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	l := log.With().Str("name", name).Logger()
	l.Info().Msg("Creating new Rollbar project resource")

	c := m.(*client.RollbarApiClient)
	p, err := c.CreateProject(name)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	l.Debug().Interface("project", p).Msg("CreateProject() result")
	l = l.With().Int("project_id", p.Id).Logger()
	d.SetId(strconv.Itoa(p.Id))

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
	tokens, err := c.ListProjectAccessTokens(p.Id)
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
		err = c.DeleteProjectAccessToken(p.Id, t.AccessToken)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
		l.Debug().
			Str("name", t.Name).
			Msg("Successfully deleted a default access token")
	}

	l.Info().Msg("Successfully created Rollbar project resource")
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Err(err).Msg("Error converting Id to integer")
		return diag.FromErr(err)
	}
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading Rollbar project resource")

	c := m.(*client.RollbarApiClient)
	proj, err := c.ReadProject(id)
	if err == client.ErrNotFound {
		d.SetId("")
		msg := fmt.Sprintf("Removing project %d from state because it was not found on Rollbar", id)
		l.Err(err).Msg(msg)
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "Project not found, removed from state",
			Detail:   msg,
		}}
	}
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}

	var mProj map[string]interface{}
	err = mapstructure.Decode(proj, &mProj)
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range mProj {
		if k == "id" {
			continue
		}
		err = d.Set(k, v)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.Itoa(proj.Id))
	l.Info().Msg("Successfully read Rollbar project resource from the API")
	return nil
}

/*
No need for this function until we have update support in the Rollbar API.

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceProjectRead(ctx, d, m)
}
*/

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	l := log.With().
		Int("projectId", projectId).
		Logger()
	l.Info().Msg("Deleting Rollbar project resource")
	c := m.(*client.RollbarApiClient)
	err = c.DeleteProject(projectId)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	l.Info().Msg("Successfully deleted Rollbar project resource")
	return nil
}
