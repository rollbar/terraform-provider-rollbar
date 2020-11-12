package rollbar

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
