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
func resourcePreCheck(d *schema.ResourceData) error {
	var integrationCount int
	for integration := range client.Integrations {
		if _, ok := d.GetOk(integration); ok {
			integrationCount++
		}
	}
	if integrationCount > 1 {
		return errors.New("only one integration allowed per resource")
	}
	return nil
}

func setBodyMapFromMap(integration string, properIntgr map[string]interface{}) (bodyMap map[string]interface{}) {
	switch integration {
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
	return bodyMap
}

func setBodyMapFromInterface(integration string, intf interface{}) (bodyMap map[string]interface{}) {
	switch integration {
	case client.SLACK:
		slackIntegration := intf.(*client.SlackIntegration)
		bodyMap = map[string]interface{}{"enabled": slackIntegration.Settings.Enabled, "channel": slackIntegration.Settings.Channel,
			"service_account_id": slackIntegration.Settings.ServiceAccountID, "show_message_buttons": slackIntegration.Settings.ShowMessageButtons}

	case client.WEBHOOK:
		webhookIntegration := intf.(*client.WebhookIntegration)
		bodyMap = map[string]interface{}{"enabled": webhookIntegration.Settings.Enabled,
			"url": webhookIntegration.Settings.URL}
	}
	return bodyMap
}
func resourceIntegrationCreateUpdateDelete(integration string, bodyMap map[string]interface{}, d *schema.ResourceData, m interface{}, action Action) (diag.Diagnostics, zerolog.Logger) {

	l := log.With().Str("integration", integration).Logger()
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
	intf, err := c.UpdateIntegration(integration, bodyMap)
	if err != nil {
		l.Err(err).Send()
		if action == CREATE || action == UPDATE {
			d.SetId("") // removing from the state
		}
		return diag.FromErr(err), l
	}
	var projectID int64
	switch integration {
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
			return diag.FromErr(err), l
		}
	}
	if action == CREATE {
		d.SetId(integrationID)
	}

	return nil, l
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := resourcePreCheck(d); err != nil {
		return diag.FromErr(err)
	}
	for integration := range client.Integrations {
		l := log.With().Str("integration", integration).Logger()

		properIntgr := parseSet(integration, d)
		if len(properIntgr) == 0 {
			l.Debug().Msg("no rollbar_integration for " + integration)
			continue
		}
		bodyMap := setBodyMapFromMap(integration, properIntgr)
		err, l := resourceIntegrationCreateUpdateDelete(integration, bodyMap, d, m, CREATE)
		if err != nil {
			return err
		}
		l.Debug().Msg("Successfully created rollbar_integration resource")
		return nil
	}
	return nil
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := resourcePreCheck(d); err != nil {
		return diag.FromErr(err)
	}
	for integration := range client.Integrations {
		l := log.With().Str("integration", integration).Logger()
		properIntgr := parseSet(integration, d)
		if len(properIntgr) == 0 {
			l.Debug().Msg("no rollbar_integration resource updates for " + integration)
			continue
		}
		bodyMap := setBodyMapFromMap(integration, properIntgr)
		err, l := resourceIntegrationCreateUpdateDelete(integration, bodyMap, d, m, UPDATE)
		if err != nil {
			return err
		}
		l.Debug().Msg("Successfully updated rollbar_integration resource")
		return nil
	}
	return nil
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	for integration := range client.Integrations {
		l := log.With().Str("integration", integration).Logger()
		properIntgr := parseSet(integration, d)
		if len(properIntgr) == 0 {
			l.Debug().Msg("no rollbar_integration resource deletes for " + integration)
			continue
		}
		bodyMap := setBodyMapFromMap(integration, properIntgr)
		err, l := resourceIntegrationCreateUpdateDelete(integration, bodyMap, d, m, DELETE)
		if err != nil {
			return err
		}
		l.Debug().Msg("Successfully deleted rollbar_integraion resource")
		return nil
	}
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
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	intf, err := c.ReadIntegration(integration)
	if err == client.ErrNotFound {
		d.SetId("")
		l.Info().Msg("Integration not found - removed from state")
		return nil
	}
	if err != nil {
		l.Err(err).Msg("error reading rollbar_integration resource")
		return diag.FromErr(err)
	}
	bodyMap := setBodyMapFromInterface(integration, intf)
	mustSet(d, integration, flattenConfig(bodyMap))
	l.Debug().Msg("Successfully read integration resource")
	return nil
}
