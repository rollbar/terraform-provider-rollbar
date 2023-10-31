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
	"errors"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

var configMap = map[string][]string{
	"email":     {"users", "teams"},
	"slack":     {"message_template", "channel", "show_message_buttons"},
	"pagerduty": {"service_key"},
	"webhook":   {"url", "format"},
}

var emailDailySummaryConfigList = []string{"summary_time", "environments", "send_only_if_data", "min_item_level"}

func CustomNotificationImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	splitID := strings.Split(d.Id(), ComplexImportSeparator)
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
			"project_api_key": {
				Description: "Overrides the project_api_key defined in the provider",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
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
									"path": {
										Description: "Path",
										Type:        schema.TypeString,
										Optional:    true,
									},
									"type": {
										Description: "Type",
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
						"summary_time": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Summary Time (email daily summary only)",
						},
						"send_only_if_data": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Send only if data (email daily summary only)",
						},
						"environments": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Environments (email daily summary only)",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"min_item_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Min item level (email daily summary only)",
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
						"url": {
							Description: "URL (webhook)",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"format": {
							Description: "Format (webhook)",
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

func parseRule(d *schema.ResourceData) (trigger string, filters interface{}) {
	rule := parseSet("rule", d)
	for key, value := range rule {
		if key == "trigger" {
			trigger = value.(string)
		} else {
			filters = value
		}
	}
	return trigger, filters
}

func cleanConfig(channel, trigger string, config map[string]interface{}) map[string]interface{} {
	returnSetMap := map[string]interface{}{}
	for key, v := range config {
		if find(configMap[channel], key) {
			returnSetMap[key] = v
		}
	}
	switch channel {
	case "email":
		if trigger == "daily_summary" {
			for key, v := range config {
				if find(emailDailySummaryConfigList, key) {
					returnSetMap[key] = v
				}
			}
		}
	case "slack":
		if trigger == "deploy" || trigger == "new_version" || trigger == "exp_repeat_item" {
			delete(returnSetMap, "show_message_buttons")
		}
	}
	return returnSetMap
}

func resourceNotificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	trigger, filters := parseRule(d)
	channel := d.Get("channel").(string)
	project_api_key := d.Get("project_api_key").(string)
	config := parseSet("config", d)
	config = cleanConfig(channel, trigger, config)
	l := log.With().Str("channel", channel).Logger()

	l.Info().Msg("Creating rollbar_notification resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarNotification, c)
	n, err := c.CreateNotification(channel, filters, trigger, config)
	client.Mutex.Unlock()

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
	project_api_key := d.Get("project_api_key").(string)
	config := parseSet("config", d)
	config = cleanConfig(channel, trigger, config)
	l := log.With().Str("channel", channel).Logger()

	l.Info().Msg("Creating rollbar_notification resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarNotification, c)
	n, err := c.UpdateNotification(id, channel, filters, trigger, config)
	client.Mutex.Unlock()

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

	l.Debug().Msg("Successfully updated rollbar_notification resource")
	return nil
}

func flattenConfig(config map[string]interface{}) *schema.Set {
	var out = make([]interface{}, 0)
	out = append(out, config)
	specResource := resourceNotification().Schema["config"].Elem.(*schema.Resource)
	f := schema.HashResource(specResource)
	set := schema.NewSet(f, out)
	return set
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
	}
	m["filters"] = filters
	out = append(out, m)
	m["trigger"] = trigger
	specResource := resourceNotification().Schema["rule"].Elem.(*schema.Resource)
	f := schema.HashResource(specResource)
	set := schema.NewSet(f, out)
	return set
}

func resourceNotificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := mustGetID(d)
	channel := d.Get("channel").(string)
	project_api_key := d.Get("project_api_key").(string)
	l := log.With().
		Int("id", id).
		Logger()
	l.Info().Msg("Reading rollbar_notification resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarNotification, c)
	n, err := c.ReadNotification(id, channel)
	client.Mutex.Unlock()

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
	project_api_key := d.Get("project_api_key").(string)
	l := log.With().Int("id", id).Logger()
	l.Info().Msg("Deleting rollbar_notification resource")

	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarNotification, c)
	err := c.DeleteNotification(id, channel)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Msg("Error deleting rollbar_notification resource")
		return diag.FromErr(err)
	}
	l.Debug().Msg("Successfully deleted rollbar_notification resource")
	return nil
}
