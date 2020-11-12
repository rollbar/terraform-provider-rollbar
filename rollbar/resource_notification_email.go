package rollbar

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceNotificationEmail constructs a schema representing a Rollbar email
// notification.
func resourceNotificationEmail() *schema.Resource {
	return &schema.Resource{
		//CreateContext: resourceEmailCreate,
		//ReadContext:   resourceEmailRead,
		//DeleteContext: resourceEmailDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required
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
						"filters": {
							Description: "Notification filters",
							Type:        schema.TypeSet,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Description:      "Filter type",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: resourceNotificationValidateFilterType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
