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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Action string

var (
	DELETE Action = "DELETE"
	UPDATE Action = "UPDATE"
	CREATE Action = "CREATE"
)

// resourceIntegraion constructs a resource representing a Rollbar integration.
func resourceIntegraion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationCreate,
		UpdateContext: resourceIntegrationUpdate,
		ReadContext:   resourceIntegrationRead,
		DeleteContext: resourceIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"project_api_key": {
				Description: "Overrides the project_api_key defined in the provider",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			client.EMAIL: {
				Description: "Email integration",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Enabled",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"scrub_params": {
							Description: "Scrub params",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
			client.PAGERDUTY: {
				Description: "PagerDuty integration",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Enabled",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"service_key": {
							Description: "Service key",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			client.SLACK: {
				Description: "Slack integration",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Enabled",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"channel": {
							Description: "Channel",
							Type:        schema.TypeString,
							Required:    true,
						},
						"service_account_id": {
							Description: "Service account ID",
							Type:        schema.TypeString,
							Required:    true,
						},
						"show_message_buttons": {
							Description: "Show message buttons",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
			client.WEBHOOK: {
				Description: "Webhook integration",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Description: "Enabled",
							Type:        schema.TypeBool,
							Required:    true,
						},
						"url": {
							Description: "URL",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func flattenIntegration(integration string, body map[string]interface{}) *schema.Set {
	var out = make([]interface{}, 0)
	out = append(out, body)

	specResource := resourceIntegraion().Schema[integration].Elem.(*schema.Resource)
	f := schema.HashResource(specResource)
	return schema.NewSet(f, out)
}

func resourcePreCheck(d *schema.ResourceData) (string, error) {
	var integrationCount int
	var validIntegration string
	for integration := range client.Integrations {
		if _, ok := d.GetOk(integration); ok {
			validIntegration = integration
			integrationCount++
		}
	}
	if integrationCount > 1 {
		return "", errors.New("only one integration allowed per resource")
	}
	return validIntegration, nil
}

func setBodyMapFromMap(integration string, properIntgr map[string]interface{}, toDelete bool) (bodyMap map[string]interface{}) {
	switch integration {
	case client.EMAIL:
		enabled := properIntgr["enabled"].(bool)
		scrubParams := properIntgr["scrub_params"].(bool)
		bodyMap = map[string]interface{}{"enabled": enabled, "scrub_params": scrubParams}

	case client.PAGERDUTY:
		enabled := properIntgr["enabled"].(bool)
		serviceKey := properIntgr["service_key"].(string)
		bodyMap = map[string]interface{}{"enabled": enabled, "service_key": serviceKey}

	case client.SLACK:
		enabled := properIntgr["enabled"].(bool)
		showMessageButtons := properIntgr["show_message_buttons"].(bool)
		channel := properIntgr["channel"].(string)
		serviceAccountID := properIntgr["service_account_id"].(string)
		bodyMap = map[string]interface{}{"channel": channel, "service_account_id": serviceAccountID,
			"enabled": enabled, "show_message_buttons": showMessageButtons}

	case client.WEBHOOK:
		enabled := properIntgr["enabled"].(bool)
		url := properIntgr["url"].(string)
		bodyMap = map[string]interface{}{"enabled": enabled, "url": url}
	}
	if toDelete {
		bodyMap["enabled"] = false
	}
	return bodyMap
}

func setBodyMapFromInterface(integration string, intf interface{}, toDelete bool) (bodyMap map[string]interface{}) {
	switch integration {
	case client.EMAIL:
		emailIntegration := intf.(*client.EmailIntegration)
		bodyMap = map[string]interface{}{"enabled": emailIntegration.Settings.Enabled,
			"scrub_params": emailIntegration.Settings.ScrubParams}

	case client.PAGERDUTY:
		pagerDutyIntegration := intf.(*client.PagerDutyIntegration)
		bodyMap = map[string]interface{}{"enabled": pagerDutyIntegration.Settings.Enabled,
			"service_key": pagerDutyIntegration.Settings.ServiceKey}

	case client.SLACK:
		slackIntegration := intf.(*client.SlackIntegration)
		bodyMap = map[string]interface{}{"enabled": slackIntegration.Settings.Enabled, "channel": slackIntegration.Settings.Channel,
			"service_account_id": slackIntegration.Settings.ServiceAccountID, "show_message_buttons": slackIntegration.Settings.ShowMessageButtons}

	case client.WEBHOOK:
		webhookIntegration := intf.(*client.WebhookIntegration)
		bodyMap = map[string]interface{}{"enabled": webhookIntegration.Settings.Enabled,
			"url": webhookIntegration.Settings.URL}
	}
	if toDelete {
		bodyMap["enabled"] = false
	}
	return bodyMap
}
func resourceIntegrationCreateUpdateDelete(integration string, bodyMap map[string]interface{}, d *schema.ResourceData, m interface{}, action Action) (zerolog.Logger, diag.Diagnostics) {
	l := log.With().Str("integration", integration).Logger()
	project_api_key := d.Get("project_api_key").(string)
	switch action {
	case CREATE:
		l.Info().Msg("Creating rollbar_integration resource")
	case UPDATE:
		l.Info().Msg("Updating rollbar_integration resource")
	case DELETE:
		l.Info().Msg("Deleting rollbar_integration resource")

	}
	var id string
	if action == UPDATE || action == DELETE {
		id = d.Id()
		l = l.With().Str("id", id).Logger()
	}
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarIntegration, c)
	intf, err := c.UpdateIntegration(integration, bodyMap)
	client.Mutex.Unlock()

	if err != nil {
		l.Err(err).Send()
		if action == CREATE || action == UPDATE {
			d.SetId("") // removing from the state
		}
		return l, diag.FromErr(err)
	}
	var projectID int64
	switch integration {
	case client.EMAIL:
		i := intf.(*client.EmailIntegration)
		projectID = i.ProjectID
	case client.PAGERDUTY:
		i := intf.(*client.PagerDutyIntegration)
		projectID = i.ProjectID
	case client.SLACK:
		i := intf.(*client.SlackIntegration)
		projectID = i.ProjectID
	case client.WEBHOOK:
		i := intf.(*client.WebhookIntegration)
		projectID = i.ProjectID
	}
	integrationID := strconv.FormatInt(projectID, 10) + ComplexImportSeparator + integration

	if action == UPDATE {
		if integrationID != id {
			err = errors.New("IDs are not equal")
			l.Err(err).Send()
			d.SetId("") // removing from the state
			return l, diag.FromErr(err)
		}
	}
	if action == CREATE {
		d.SetId(integrationID)
	}
	return l, nil
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	var integration string
	if integration, err = resourcePreCheck(d); err != nil {
		return diag.FromErr(err)
	}
	properIntgr := parseSet(integration, d)
	bodyMap := setBodyMapFromMap(integration, properIntgr, false)
	l, e := resourceIntegrationCreateUpdateDelete(integration, bodyMap, d, m, CREATE)
	if e != nil {
		return e
	}
	l.Debug().Msg("Successfully created rollbar_integration resource")
	return nil
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	var integration string
	if integration, err = resourcePreCheck(d); err != nil {
		return diag.FromErr(err)
	}
	properIntgr := parseSet(integration, d)
	bodyMap := setBodyMapFromMap(integration, properIntgr, false)
	l, e := resourceIntegrationCreateUpdateDelete(integration, bodyMap, d, m, UPDATE)
	if e != nil {
		return e
	}
	l.Debug().Msg("Successfully updated rollbar_integration resource")
	return nil
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	var integration string
	if integration, err = resourcePreCheck(d); err != nil {
		return diag.FromErr(err)
	}
	properIntgr := parseSet(integration, d)
	bodyMap := setBodyMapFromMap(integration, properIntgr, true)
	l, e := resourceIntegrationCreateUpdateDelete(integration, bodyMap, d, m, DELETE)
	if e != nil {
		return e
	}
	d.SetId("") // removing from the state
	l.Debug().Msg("Successfully deleted rollbar_integraion resource")
	return nil
}

func resourceIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	l := log.With().
		Str("id", id).
		Logger()
	spl := strings.Split(id, ComplexImportSeparator)
	integration := spl[1]
	l.Info().Msg("Reading rollbar_integration resource")

	project_api_key := d.Get("project_api_key").(string)
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	if len(project_api_key) > 0 {
		c = client.NewClient(c.BaseURL, project_api_key)
	}

	client.Mutex.Lock()
	setResourceHeader(rollbarIntegration, c)
	intf, err := c.ReadIntegration(integration)
	client.Mutex.Unlock()

	if err == client.ErrNotFound {
		d.SetId("")
		l.Info().Msg("Integration not found - removed from state")
		return nil
	}
	if err != nil {
		l.Err(err).Msg("error reading rollbar_integration resource")
		return diag.FromErr(err)
	}
	bodyMap := setBodyMapFromInterface(integration, intf, false)
	mustSet(d, integration, flattenIntegration(integration, bodyMap))
	l.Debug().Msg("Successfully read integration resource")
	return nil
}
