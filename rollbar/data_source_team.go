/*
 * Copyright (c) 2021 Rollbar, Inc.
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
)

// dataSourceTeam is a data source returning a Rollbar team.
func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Human readable name for the team",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Computed values
			"id": {
				Description: "ID of the team",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"access_level": {
				Description: `The team's access level.  Must be "standard", "light", or "view".  Defaults to "standard".`,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"account_id": {
				Description: "ID of account that owns the team",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceTeamRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := data.Get("name").(string)

	log.Debug().Msg("Reading team list from API")
	c := m.(*client.RollbarAPIClient)
	teams, err := c.ListTeams()

	if err != nil {
		return diag.FromErr(err)
	}

	var team client.Team
	var found bool
	for _, p := range teams {
		if p.Name == name {
			found = true
			team = p
		}
	}

	if !found {
		data.SetId("")
		return diag.FromErr(fmt.Errorf("no team with the name %s found", name))
	}

	id := fmt.Sprintf("%d", team.ID)
	data.SetId(id)
	mustSet(data, "name", team.Name)
	mustSet(data, "account_id", team.AccountID)
	mustSet(data, "access_level", team.AccessLevel)
	log.Debug().Msg("Successfully read rollbar_team resource")

	return nil
}
