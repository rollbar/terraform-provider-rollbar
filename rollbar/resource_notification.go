package rollbar

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog/log"
	"strconv"
)

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
				Description: "Filter operation",
				Type:        schema.TypeString,
				Required:    true,
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					s := v.(string)

					// Check if operator is allowed
					quotedOperators := make([]string, len(allowedOperators)) // Used for error message below
					for i, operator := range allowedOperators {
						if operator == s {
							return nil
						}
						quotedOperators[i] = strconv.Quote(allowedOperators[i])
					}

					// Operator was not allowed, so construct error message
					detail := "Must be "
					opCount := len(allowedOperators)
					switch opCount {
					case 1:
						detail = detail + quotedOperators[0]
					case 2:
						detail = detail + quotedOperators[0] + " or " + quotedOperators[1]
					default:
						for i := 0; i < opCount-1; i++ {
							detail = detail + quotedOperators[i] + ", "
						}
						detail = detail + "or " + quotedOperators[opCount-1]
					}
					d := diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "Invalid filter operation",
						Detail:        detail,
						AttributePath: p,
					}
					log.Error().Interface("diagnostic", d).Send()
					return diag.Diagnostics{d}
				},
			},
			"value": {
				Description: "Filter value",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}

}