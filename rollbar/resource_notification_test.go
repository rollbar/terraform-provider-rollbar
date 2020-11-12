package rollbar

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestResourceNotificationValidateFilterType tests validation of `type`
// argument on notification filter configuration.
func TestResourceNotificationValidateFilterType(t *testing.T) {
	validTypes := []string{
		"environment",
		"level",
		"title",
		"filename",
		"context",
		"method",
		"framework",
		"path",
	}
	p := cty.Path{} // placeholder
	for _, filterType := range validTypes {
		d := resourceNotificationValidateFilterType(filterType, p)
		assert.Nil(t, d)
	}
	d := resourceNotificationValidateFilterType("invalid-filter-type", p)
	assert.IsType(t, diag.Diagnostics{}, d)
	assert.Len(t, d, 1)
	assert.Equal(t, diag.Error, d[0].Severity)
}

// TestResourceNotificationValidateFilterType tests validation of `type`
// argument on notification filter configuration.
func TestResourceNotificationValidateTrigger(t *testing.T) {
	validTriggers := []string{
		"exp_repeat_item",
		"occurrence_rate",
		"resolved_item",
		"reactivated_item",
		"new_item",
	}
	p := cty.Path{} // placeholder
	for _, trigger := range validTriggers {
		d := resourceNotificationValidateTrigger(trigger, p)
		assert.Nil(t, d)
	}
	d := resourceNotificationValidateTrigger("invalid-trigger", p)
	assert.IsType(t, diag.Diagnostics{}, d)
	assert.Len(t, d, 1)
	assert.Equal(t, diag.Error, d[0].Severity)
}
