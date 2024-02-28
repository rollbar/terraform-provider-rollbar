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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamRead,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Description:   "Team ID",
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Description: "Name of the team",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"account_id": {
				Description: "Account ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"access_level": {
				Description: "The team's access level",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var team client.Team
	var l zerolog.Logger
	teamID, ok := d.GetOk("team_id")
	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	c.SetHeaderDataSource(rollbarTeam)

	if ok {
		l = log.With().
			Int("id", teamID.(int)).
			Logger()
		l.Debug().Msg("Reading Team from Rollbar by ID")
		respTeam, err := c.ReadTeam(teamID.(int))
		if err != nil {
			return diag.Errorf("Team not found by ID: %v", err)
		}
		team = respTeam
	} else {
		name, nameOk := d.GetOk("name")
		if !nameOk {
			return diag.Errorf("Data Source requires either \"name\" or \"team_id\"")
		}
		l = log.With().
			Str("name", name.(string)).
			Logger()
		l.Debug().Msg("Reading team from Rollbar by name")

		teams, err := c.ListTeams()
		if err != nil {
			return diag.FromErr(err)
		}

		t, err := findTeamByName(teams, name.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		team = t
	}
	d.SetId(strconv.FormatInt(int64(team.ID), 10))
	_ = d.Set("team_id", team.ID)
	_ = d.Set("name", team.Name)
	_ = d.Set("access_level", team.AccessLevel)
	_ = d.Set("account_id", team.AccountID)
	return nil
}

func findTeamByName(teams []client.Team, name string) (client.Team, error) {
	for _, team := range teams {
		if team.Name == name {
			return team, nil
		}
	}
	return client.Team{}, fmt.Errorf("Team not found by name: %s", name)
}
