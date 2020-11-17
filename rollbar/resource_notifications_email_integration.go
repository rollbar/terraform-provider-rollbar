package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

// resourceNotificationsEmailIntegration constructs a resource representing a
// Rollbar notifications email integration.
func resourceNotificationsEmailIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationsEmailIntegrationUpdate,
		ReadContext:   resourceNotificationsEmailIntegrationRead,
		UpdateContext: resourceNotificationsEmailIntegrationUpdate,
		DeleteContext: resourceNotificationsEmailIntegrationDelete,

		Schema: map[string]*schema.Schema{
			// Required
			"enabled": {
				Description: "Enable the Email notifications globally",
				Type:        schema.TypeBool,
				Required:    true,
			},

			// Optional
			"include_request_params": {
				Description: "Whether to include request parameters",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

// resourceNotificationsEmailIntegrationUpdate updates a Rollbar notifications
// email integration resource.
func resourceNotificationsEmailIntegrationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	enabled := d.Get("enabled").(bool)
	includeRequestParams := d.Get("include_request_params").(bool)
	l := log.With().
		Bool("enabled", enabled).
		Bool("include_request_params", includeRequestParams).
		Logger()
	l.Info().Msg("Updating " + resNotificationsEmail)
	c := m.(*client.RollbarApiClient)
	err := c.UpdateNotificationsEmailIntegration(enabled, includeRequestParams)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	d.SetId("singleton")
	l.Info().Msg("Successfully updated " + resNotificationsEmail)
	return resourceNotificationsEmailIntegrationRead(ctx, d, m)
}

// resourceNotificationsEmailIntegrationRead reads a Rollbar notifications email
// integration resource.
func resourceNotificationsEmailIntegrationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{{
		Severity: diag.Error,
		Summary:  "Read not yet implemented for " + resNotificationsEmail,
	}}
}

// resourceNotificationsEmailIntegrationDelete deletes a Rollbar notifications
// email integration resource.
func resourceNotificationsEmailIntegrationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debug().Msgf(
		"Deleting %s. This removes the resource from TF state but does not touch the API.",
		resNotificationsEmail,
	)
	d.SetId("")
	return nil
}
