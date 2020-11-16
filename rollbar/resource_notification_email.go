package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceNotificationEmail constructs a schema representing a Rollbar email
// notification.
func resourceNotificationEmail() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationEmailCreate,
		ReadContext:   resourceNotificationEmailRead,
		UpdateContext: resourceNotificationEmailUpdate,
		DeleteContext: resourceNotificationEmailDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
			"token_id": {
				Description: `ID (access token string) of project access token with "write" scope`,
				Type:        schema.TypeString,
				Required:    true,
			},
			"rule": resourceNotificationRuleSchema,
		},
	}
}

var resourceNotificationEmailNotImplementedMessage = "resource `rollbar_notification_email` not yet implemented"
var resourceNotificationEmailNotImplementedDiagnostics = diag.Diagnostics{diag.Diagnostic{
	Severity: diag.Error,
	Summary:  resourceNotificationEmailNotImplementedMessage,
}}

func resourceNotificationEmailCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNotificationEmailNotImplementedDiagnostics
}

func resourceNotificationEmailRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNotificationEmailNotImplementedDiagnostics
}

func resourceNotificationEmailUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNotificationEmailNotImplementedDiagnostics
}

func resourceNotificationEmailDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNotificationEmailNotImplementedDiagnostics
}
