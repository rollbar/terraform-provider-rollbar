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

var DELETE Action = "DELETE"
var UPDATE Action = "UPDATE"
var CREATE Action = "CREATE"

// resourceIntegraion constructs a resource representing a Rollbar integration.
func resourceIntegraion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationCreate,
		UpdateContext: resourceIntegrationUpdate,
		ReadContext:   resourceIntegrationRead,
		DeleteContext: resourceIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"slack": {
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
		},
	}
}

func resourceIntegrationCreateUpdateDelete(integration, channel, serviceAccountID string, enabled, showMessageButtons bool, d *schema.ResourceData, m interface{}, action Action) (diag.Diagnostics, zerolog.Logger) {

	l := log.With().Str("integration", integration).Logger()
	switch action {
	case "CREATE":
		l.Info().Msg("Creating rollbar_integration resource")
	case "UPDATE":
		l.Info().Msg("Updating rollbar_integration resource")
	case "DELETE":
		l.Info().Msg("Deleting rollbar_integration resource")

	}
	var id string
	if action == "UPDATE" || action == "DELETE" {
		id = d.Id()
		l = l.With().Str("id", id).Logger()
	}
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	intf, err := c.UpdateIntegration(integration, channel, serviceAccountID, enabled, showMessageButtons)
	if err != nil {
		l.Err(err).Send()
		if action == "CREATE" || action == "UPDATE" {
			d.SetId("") // removing from the state
		}
		return diag.FromErr(err), l
	}

	slackIntegration := intf.(*client.SlackIntegration)
	integrationID := strconv.Itoa(slackIntegration.ProjectID) + ComplexImportSeparator + integration

	if action == "UPDATE" {
		if integrationID != id {
			err = errors.New("IDs are not equal")
			l.Err(err).Send()
			d.SetId("") // removing from the state
			return diag.FromErr(err), l
		}
	}
	if action == "CREATE" {
		d.SetId(integrationID)
	}

	return nil, l
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	integration := "slack"
	l := log.With().Str("integration", integration).Logger()
	slack := parseSet(integration, d)
	if len(slack) == 0 {
		l.Debug().Msg(" rollbar_integration resource cannot be created")
		return nil
	}
	enabled := slack["enabled"].(bool)
	showMessageButtons := slack["show_message_buttons"].(bool)
	channel := slack["channel"].(string)
	serviceAccountID := slack["service_account_id"].(string)
	err, l := resourceIntegrationCreateUpdateDelete(integration, channel, serviceAccountID, enabled, showMessageButtons, d, m, "CREATE")
	if err != nil {
		return err
	}
	l.Debug().Msg("Successfully created rollbar_integration resource")
	return nil
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	integration := "slack"
	l := log.With().Str("integration", integration).Logger()
	slack := parseSet(integration, d)
	if len(slack) == 0 {
		l.Debug().Msg("rollbar_integration resource cannot be updated")
		return nil
	}
	enabled := slack["enabled"].(bool)
	channel := slack["channel"].(string)
	showMessageButtons := slack["show_message_buttons"].(bool)
	serviceAccountID := slack["service_account_id"].(string)
	err, l := resourceIntegrationCreateUpdateDelete(integration, channel, serviceAccountID, enabled, showMessageButtons, d, m, "UPDATE")
	if err != nil {
		return err
	}
	l.Debug().Msg("Successfully updated rollbar_integraion resource")
	return nil
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	integration := "slack"
	l := log.With().Str("integration", integration).Logger()
	slack := parseSet(integration, d)
	if len(slack) == 0 {
		l.Debug().Msg(" rollbar_integration resource cannot be deleted")
		return nil
	}
	channel := slack["channel"].(string)
	serviceAccountID := slack["service_account_id"].(string)
	err, l := resourceIntegrationCreateUpdateDelete(integration, channel, serviceAccountID, false,
		false, d, m, "DELETE")
	if err != nil {
		return err
	}
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
	c := m.(map[string]*client.RollbarAPIClient)[projectKeyToken]
	intf, err := c.ReadIntegration(integration)
	if err == client.ErrNotFound {
		d.SetId("")
		l.Info().Msg("Integration not found - removed from state")
		return nil
	}
	slackIntegration := intf.(*client.SlackIntegration)
	if err != nil {
		l.Err(err).Msg("error reading rollbar_integration resource")
		return diag.FromErr(err)
	}
	slack := map[string]interface{}{"enabled": slackIntegration.Settings.Enabled, "channel": slackIntegration.Settings.Channel,
		"service_account_id": slackIntegration.Settings.ServiceAccountID, "show_message_buttons": slackIntegration.Settings.ShowMessageButtons}

	mustSet(d, "slack", flattenConfig(slack))
	l.Debug().Msg("Successfully read integration resource")
	return nil
}
