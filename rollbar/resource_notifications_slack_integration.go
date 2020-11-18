package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

// resourceNotificationsSlackIntegration constructs a resource representing a
// Rollbar notifications Slack integration.
func resourceNotificationsSlackIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationsSlackIntegrationUpdate,
		ReadContext:   resourceNotificationsSlackIntegrationRead,
		UpdateContext: resourceNotificationsSlackIntegrationUpdate,
		DeleteContext: resourceNotificationsSlackIntegrationDelete,

		Schema: map[string]*schema.Schema{
			// Required
			"enabled": {
				Description: "Enable the Slack notifications globally",
				Type:        schema.TypeBool,
				Required:    true,
			},
			// FIXME: This should be the ID of the token, when that becomes available.
			//  https://github.com/rollbar/terraform-provider-rollbar/issues/73
			"project_access_token": {
				Description: "Project access token with 'write' scope",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional
			"service_account_id": {
				Description: "You can find your Service Account ID in https://rollbar.com/settings/integrations/#slack",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"channel": {
				Description: "The default Slack channel to send the messages",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"show_message_buttons": {
				Description: "Show the Slack actionable buttons",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

// resourceNotificationsSlackIntegrationUpdate updates a Rollbar notifications
// email integration resource.
func resourceNotificationsSlackIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	args := client.SlackIntegrationArgs{
		Enabled:            d.Get("enabled").(bool),
		ServiceAccountID:   d.Get("service_account_id").(int),
		Channel:            d.Get("channel").(string),
		ShowMessageButtons: d.Get("show_message_buttons").(bool),
		Token:              d.Get("project_access_token").(string),
	}
	l := log.With().
		Interface("args", args).
		Logger()
	l.Info().Msg("Updating " + resNotificationsSlack)
	c := m.(*client.RollbarApiClient)
	err := c.UpdateNotificationsSlackIntegration(args)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	d.SetId("singleton")
	l.Info().Msg("Successfully updated " + resNotificationsSlack)
	return resourceNotificationsSlackIntegrationRead(ctx, d, m)
}

// resourceNotificationsSlackIntegrationRead reads a Rollbar notifications Slack
// integration resource.
func resourceNotificationsSlackIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{{
		Severity: diag.Error,
		Summary:  "Read not yet implemented for " + resNotificationsSlack,
	}}
}

// resourceNotificationsSlackIntegrationDelete deletes a Rollbar notifications
// Slack integration resource.
func resourceNotificationsSlackIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debug().Msgf(
		"Deleting %s. This removes the resource from TF state but does not touch the API.",
		resNotificationsSlack,
	)
	d.SetId("")
	return nil
}
