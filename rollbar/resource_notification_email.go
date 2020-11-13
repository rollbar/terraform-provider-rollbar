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
				Description: `ID of project access token with "write" scope`,
				Type:        schema.TypeInt,
				Required:    true,
			},
			"rule": {
				Description: "Notification rule",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger": {
							Description:      "Trigger for the notification",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: resourceNotificationValidateTrigger,
						},
						//"environment_filter": {
						//	Description: "Environment filter",
						//	Type:        schema.TypeSet,
						//	Optional:    true,
						//	Elem:
						//},
					},
				},
			},
		},
	}
}

func resourceNotificationEmailCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNotificationEmailRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNotificationEmailUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNotificationEmailDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
