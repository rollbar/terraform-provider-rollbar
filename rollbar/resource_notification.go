package rollbar

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceNotification constructs a schema representing a set of Rollbar
// notification rules.
func resourceNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationCreate,
		ReadContext:   resourceNotificationRead,
		UpdateContext: resourceNotificationUpdate,
		DeleteContext: resourceNotificationDelete,

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
			"rule":  resourceNotificationRuleSchema,
			"email": resourceNotificationEmailSchema,
			"slack": resourceNotificationSlackSchema,
		},
	}
}

var resourceNotificationEmailNotImplementedMessage = "resource `rollbar_notification_email` not yet implemented"
var resourceNotificationNotImplementedDiagnostics = diag.Diagnostics{diag.Diagnostic{
	Severity: diag.Error,
	Summary:  resourceNotificationEmailNotImplementedMessage,
}}

func resourceNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNotificationNotImplementedDiagnostics
}

func resourceNotificationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNotificationNotImplementedDiagnostics
}

func resourceNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNotificationNotImplementedDiagnostics
}

func resourceNotificationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceNotificationNotImplementedDiagnostics
}

// resourceNotificationValidateFilterType validates the `type` argument for a
// notification filter resource.
func resourceNotificationValidateFilterType(v interface{}, p cty.Path) diag.Diagnostics {
	s := v.(string)
	switch s {
	case "environment", "level", "title", "filename", "context", "method", "framework", "path":
		return nil
	default:
		d := diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf(`Invalid filter type: "%s"`, s),
			Detail:        `Must be "environment", "level", "title", "filename", "context", "method", "framework", or "path"`,
			AttributePath: p,
		}
		return diag.Diagnostics{d}
	}
}

// resourceNotificationValidateTrigger validates the `trigger` argument for a
// notification.
func resourceNotificationValidateTrigger(v interface{}, p cty.Path) diag.Diagnostics {
	s := v.(string)
	switch s {
	case "exp_repeat_item", "occurrence_rate", "resolved_item", "reactivated_item", "new_item":
		return nil
	default:
		d := diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf(`Invalid trigger: "%s"`, s),
			Detail:        `Must be "exp_repeat_item", "occurrence_rate", "resolved_item", "reactivated_item", or "new_item"`,
			AttributePath: p,
		}
		return diag.Diagnostics{d}
	}
}

var resourceNotificationRuleSchema = &schema.Schema{
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
			"environment_filter": {
				Description: "Environment filter",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        operationValueFilterSchema([]string{"eq", "neq"}),
			},
			"level_filter": {
				Description: "Level filter",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        operationValueFilterSchema([]string{"eq", "gte"}),
			},
			"title_filter": {
				Description: "Title filter",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        operationValueFilterSchema([]string{"within", "nwithin", "regex", "nregex"}),
			},
			"filename_filter": {
				Description: "Filename filter",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        operationValueFilterSchema([]string{"within", "nwithin", "regex", "nregex"}),
			},
			"context_filter": {
				Description: "Context filter",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        operationValueFilterSchema([]string{"startswith", "eq", "neq"}),
			},
			"method_filter": {
				Description: "Method filter",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        operationValueFilterSchema([]string{"within", "nwithin", "regex", "nregex"}),
			},
			"framework_filter": {
				Description: "Framework filter",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        operationValueFilterSchema([]string{"eq"}),
			},
		},
	},
}

// operationValueFilterSchema constructs a resource schema for Rollbar
// notification filters based on a combination of `operation` and `value`
// arguments.
func operationValueFilterSchema(allowedOperators []string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operation": {
				Description:      "Filter operation",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateStringInSlice(allowedOperators),
			},
			"value": {
				Description: "Filter value",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}

}

// resourceNotificationEmailSchema constructs a schema for a Rollbar
// notification rule email config.
var resourceNotificationEmailSchema = &schema.Schema{
	Description: "Email config",
	Type:        schema.TypeSet,
	Optional:    true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"teams": {
				Description: "List of team names to send emails to",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"users": {
				Description: "List of usernames or email addresses to send emails to",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

// resourceNotificationSlackSchema constructs a schema for a Rollbar
// notification rule Slack config.
var resourceNotificationSlackSchema = &schema.Schema{
	Description: "Slack config",
	Type:        schema.TypeSet,
	Optional:    true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"message_template": {
				Description: "Custom template for the Slack message",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"show_message_buttons": {
				Description: "Show the Slack actionable buttons",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"channel": {
				Description: "Slack channel to send the messages",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	},
}
