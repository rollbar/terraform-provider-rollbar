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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
)

// resourceTeam constructs a resource representing a Rollbar team.
func resourceNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationCreateOrUpdate,
		UpdateContext: resourceNotificationCreateOrUpdate,
		ReadContext:   resourceNotificationRead,
		DeleteContext: resourceNotificationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
			"channel": {
				Description: "Channel",
				Type:        schema.TypeString,
				Required:    true,
			},
			"rule": {
				Description: "Human readable name for the rule",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger": {
							Description: "Trigger",
							Type:        schema.TypeString,
							Required:    true,
						},
						"filters": {
							Description: "Filters",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description: "Operation",
										Type:        schema.TypeString,
										Required:    true,
									},
									"operation": {
										Description: "Operation",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"value": {
										Description: "Value",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			"config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"users": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Users",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"teams": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Teams",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func parseSet(setName string, d *schema.ResourceData) map[string]interface{} {
	setMap, ok := d.GetOk(setName)
	properSetMap := map[string]interface{}{}

	if ok {
		set := setMap.(*schema.Set).List()
		for _, s := range set {
			properSetMap, ok = s.(map[string]interface{})
			if ok {
				return properSetMap
			}
		}
	}
	return nil
}

func parseRule(d *schema.ResourceData) (string, interface{}) {
	rule := parseSet("rule", d)

	var trigger string
	var filters interface{}
	for key, value := range rule {
		if key == "trigger" {
			trigger = value.(string)
		} else {
			filters = value
		}
	}
	return trigger, filters
}

func resourceNotificationCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	trigger, filters := parseRule(d)
	channel := d.Get("channel").(string)
	config := parseSet("config", d)
	l := log.With().Str("channel", channel).Logger()

	l.Info().Msg("Creating or updating rollbar_notification resource")

	c := m.([]*client.RollbarAPIClient)[1]
	p, err := c.CreateOrUpdateNotification(channel, filters, trigger, config)
	if err != nil {
		l.Err(err).Send()
		d.SetId("") // removing from the state
		return diag.FromErr(err)
	}
	l = l.With().Int("id", p.ID).Logger()

	d.SetId(strconv.Itoa(p.ID))
	l.Debug().Msg("Successfully created or updated Rollbar notification resource")

	return nil
}

func resourceNotificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := mustGetID(d)
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading rollbar_notification resource")
	//c := m.([]*client.RollbarAPIClient)[0])
	//t, err := c.ReadTeam(id)
	//if err == client.ErrNotFound {
	//d.SetId("")
	//l.Info().Msg("Team not found - removed from state")
	//return nil
	//}
	//if err != nil {
	//	l.Err(err).Msg("error reading rollbar_team resource")
	//	return diag.FromErr(err)
	//}
	//mustSet(d, "name", t.Name)
	//mustSet(d, "account_id", t.AccountID)
	//mustSet(d, "access_level", t.AccessLevel)
	//l.Debug().Msg("Successfully read rollbar_team resource")
	return nil
}

func resourceNotificationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := mustGetID(d)

	l := log.With().Int("id", id).Logger()
	l.Info().Msg("Deleting rollbar_notification resource")
	//c := m.([]*client.RollbarAPIClient)[0])
	//err := c.DeleteTeam(id)
	//if err != nil {
	//	l.Err(err).Msg("Error deleting rollbar_team resource")
	//	return diag.FromErr(err)
	//}
	//l.Debug().Msg("Successfully deleted rollbar_team resource")
	return nil
}
