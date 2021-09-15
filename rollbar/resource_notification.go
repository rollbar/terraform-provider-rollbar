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
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

var configMap = map[string][]string{"email": {"users", "teams"},
	"slack":     {"message_template", "channel", "show_message_buttons"},
	"pagerduty": {"service_key"}}

func CustomNotificationImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	splitID := strings.Split(d.Id(), ",")
	if len(splitID) > 1 {
		mustSet(d, "channel", splitID[0])
		d.SetId(splitID[1])
	}
	return []*schema.ResourceData{d}, nil
}

// resourceNotification constructs a resource representing a Rollbar notification.
func resourceNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationCreate,
		UpdateContext: resourceNotificationUpdate,
		ReadContext:   resourceNotificationRead,
		DeleteContext: resourceNotificationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: CustomNotificationImport,
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
									"period": {
										Description: "Period",
										Type:        schema.TypeFloat,
										Optional:    true,
										Default:     0,
									},
									"count": {
										Description: "Count",
										Type:        schema.TypeFloat,
										Optional:    true,
										Default:     0,
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
							Description: "Users (email)",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"teams": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Teams (email)",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"message_template": {
							Description: "Message template (slack)",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"channel": {
							Description: "Channel (slack)",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"show_message_buttons": {
							Description: "Show message buttons (slack)",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"service_key": {
							Description: "Service key (pagerduty)",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func parseSet(setName string, d *schema.ResourceData) map[string]interface{} {
	setMap, ok := d.GetOk(setName)
	var properSetMap map[string]interface{}

	if ok {
		set := setMap.(*schema.Set).List()
		for _, s := range set {
			properSetMap, ok = s.(map[string]interface{})
			if ok {
				return properSetMap
			}
		}
	}
	return map[string]interface{}{}
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

func cleanConfig(channel string, config map[string]interface{}) map[string]interface{} {
	returnSetMap := map[string]interface{}{}
	for key, v := range config {
		if find(configMap[channel], key) {
			returnSetMap[key] = v
		}
	}
	return returnSetMap
}

func resourceNotificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	trigger, filters := parseRule(d)
	channel := d.Get("channel").(string)
	config := parseSet("config", d)
	config = cleanConfig(channel, config)
	l := log.With().Str("channel", channel).Logger()

	l.Info().Msg("Creating rollbar_notification resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	n, err := c.CreateNotification(channel, filters, trigger, config)
	if err != nil {
		l.Err(err).Send()
		d.SetId("") // removing from the state
		return diag.FromErr(err)
	}
	l = l.With().Int("id", n.ID).Logger()

	d.SetId(strconv.Itoa(n.ID))
	l.Debug().Msg("Successfully created rollbar_notification resource")

	return nil
}

func resourceNotificationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	id := mustGetID(d)
	trigger, filters := parseRule(d)
	channel := d.Get("channel").(string)
	config := parseSet("config", d)
	config = cleanConfig(channel, config)
	l := log.With().Str("channel", channel).Logger()

	l.Info().Msg("Creating rollbar_notification resource")
	l.Print(config)

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	n, err := c.UpdateNotification(id, channel, filters, trigger, config)

	if err != nil {
		l.Err(err).Send()
		d.SetId("") // removing from the state
		return diag.FromErr(err)
	}
	if n.ID != id {
		err = errors.New("IDs are not equal")
		l.Err(err).Send()
		d.SetId("") // removing from the state
		return diag.FromErr(err)
	}
	l = l.With().Int("id", n.ID).Logger()

	l.Debug().Msg("Successfully updated Rollbar notification resource")
	return nil
}

func flattenConfig(config map[string]interface{}) *schema.Set {
	var out = make([]interface{}, 0)
	m := make(map[string]interface{})
	for key, value := range config {
		m[key] = value
		out = append(out, m)
	}
	specResource := resourceNotification().Schema["config"].Elem.(*schema.Resource)
	f := schema.HashResource(specResource)
	return schema.NewSet(f, out)
}

func flattenRule(filters []interface{}, trigger string) *schema.Set {
	var out = make([]interface{}, 0)
	m := make(map[string]interface{})
	for _, filter := range filters {
		filterConv := filter.(map[string]interface{})
		filterValue := filterConv["value"]
		switch v := filterValue.(type) {
		case int:
			filterConv["value"] = strconv.Itoa(v)
		case int8:
			filterConv["value"] = strconv.FormatInt(int64(v), 10)
		case int16:
			filterConv["value"] = strconv.FormatInt(int64(v), 10)
		case int32:
			filterConv["value"] = strconv.FormatInt(int64(v), 10)
		case int64:
			filterConv["value"] = strconv.FormatInt(v, 10)
		case float32:
			filterConv["value"] = strconv.FormatFloat(float64(v), 'f', -1, 32)
		case float64:
			filterConv["value"] = strconv.FormatFloat(v, 'f', -1, 64)
		}
		filter = filterConv
	}
	m["filters"] = filters
	out = append(out, m)
	m["trigger"] = trigger
	out = append(out, m)
	specResource := resourceNotification().Schema["rule"].Elem.(*schema.Resource)
	f := schema.HashResource(specResource)
	return schema.NewSet(f, out)
}

func resourceNotificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := mustGetID(d)
	channel := d.Get("channel").(string)
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading rollbar_notification resource")
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	n, err := c.ReadNotification(id, channel)
	if err == client.ErrNotFound {
		d.SetId("")
		l.Info().Msg("Notification not found - removed from state")
		return nil
	}
	if err != nil {
		l.Err(err).Msg("error reading rollbar_notification resource")
		return diag.FromErr(err)
	}

	mustSet(d, "config", flattenConfig(n.Config))
	mustSet(d, "rule", flattenRule(n.Filters, n.Trigger))
	l.Debug().Msg("Successfully read rollbar_notification resource")
	return nil
}

func resourceNotificationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := mustGetID(d)
	channel := d.Get("channel").(string)
	l := log.With().Int("id", id).Logger()
	l.Info().Msg("Deleting rollbar_notification resource")
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	err := c.DeleteNotification(id, channel)
	if err != nil {
		l.Err(err).Msg("Error deleting rollbar_notification resource")
		return diag.FromErr(err)
	}
	l.Debug().Msg("Successfully deleted rollbar_notification resource")
	return nil
}
