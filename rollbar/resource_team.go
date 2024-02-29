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

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

// resourceTeam constructs a resource representing a Rollbar team.
func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		DeleteContext: resourceTeamDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Description: "Human readable name for the team",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			// Optional
			"access_level": {
				Description:      `The team's access level.  Must be "standard", "light", or "view".  Defaults to "standard".`,
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "standard",
				ForceNew:         true,
				ValidateDiagFunc: resourceTeamValidateAccessLevel,
			},

			// Computed
			"account_id": {
				Description: "ID of account that owns the team",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceTeamValidateAccessLevel(v interface{}, p cty.Path) diag.Diagnostics {
	s := v.(string)
	switch s {
	case "standard", "light", "view":
		return nil
	default:
		summary := fmt.Sprintf(`Invalid access_level: %q`, s)
		d := diag.Diagnostic{
			Severity:      diag.Error,
			AttributePath: p,
			Summary:       summary,
			Detail:        `Must be "standard", "light", or "view"`,
		}
		return diag.Diagnostics{d}
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	level := d.Get("access_level").(string)
	l := log.With().Str("name", name).Str("access_level", level).Logger()
	l.Info().Msg("Creating rollbar_team resource")
	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	c.SetHeaderResource(rollbarTeam)
	t, err := c.CreateTeam(name, level)

	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	teamID := t.ID
	l = l.With().Int("teamID", teamID).Logger()
	d.SetId(strconv.Itoa(teamID))
	l.Debug().Int("id", teamID).Msg("Successfully created rollbar_team resource")
	return resourceTeamRead(ctx, d, m)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := mustGetID(d)
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading rollbar_team resource")
	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	c.SetHeaderResource(rollbarTeam)
	t, err := c.ReadTeam(id)

	if err == client.ErrNotFound {
		d.SetId("")
		l.Err(err).Msg("Team not found - removed from state")
		return nil
	}
	if err != nil {
		l.Err(err).Msg("error reading rollbar_team resource")
		return diag.FromErr(err)
	}
	mustSet(d, "name", t.Name)
	mustSet(d, "account_id", t.AccountID)
	mustSet(d, "access_level", t.AccessLevel)
	l.Debug().Msg("Successfully read rollbar_team resource")
	return nil
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := mustGetID(d)

	l := log.With().Int("id", id).Logger()
	l.Info().Msg("Deleting rollbar_team resource")
	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	c.SetHeaderResource(rollbarTeam)
	err := c.DeleteTeam(id)

	if err != nil {
		l.Err(err).Msg("Error deleting rollbar_team resource")
		return diag.FromErr(err)
	}
	l.Debug().Msg("Successfully deleted rollbar_team resource")
	return nil
}
