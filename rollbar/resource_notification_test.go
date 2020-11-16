package rollbar

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"regexp"
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

// TestOperationValueFilterSchema tests validation in the schema resource
// constructed by operationValueFilterSchema().
func TestOperationValueFilterSchema(t *testing.T) {
	checkOps := func(validOperations []string, expectedErrorDetail string) {
		ovfs := operationValueFilterSchema(validOperations)
		p := cty.GetAttrPath("foo").GetAttr("bar") // Placeholder
		validationFunc := ovfs.Schema["operation"].ValidateDiagFunc
		for _, op := range validOperations {
			diags := validationFunc(op, p)
			assert.Nil(t, diags)
		}
		diags := validationFunc("invalid-operation", p)
		assert.NotNil(t, diags)
		assert.Len(t, diags, 1)
		d := diags[0]
		assert.Equal(t, d.Severity, diag.Error)
		assert.Equal(t, expectedErrorDetail, d.Detail)
	}

	// Check single and multiple valid operations
	checkOps(
		[]string{"foo"},
		`Must be "foo".`+"\nPath: foo.bar",
	)
	checkOps(
		[]string{"foo", "bar"},
		`Must be "foo" or "bar".`+"\nPath: foo.bar",
	)
	checkOps(
		[]string{"foo", "bar", "baz"},
		`Must be "foo", "bar", or "baz".`+"\nPath: foo.bar",
	)
}

func (s *AccSuite) TestNotificationNotImplemented() {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
		  name = "%s"
		}

		resource "rollbar_project_access_token" "write_token" {
		  name = "test-write-token"
		  project_id = rollbar_project.test.id
		  scopes = ["write"]
		}

		resource "rollbar_notification" "foo" {
		  token_id = rollbar_project_access_token.write_token.id

		  rule {
			trigger = "new_item"
			environment_filter {
			  operation = "eq"
			  value = "foo"
			}
			level_filter {
			  operation = "gte"
			  value = "bar"
			}
			title_filter {
			  operation = "within"
			  value = "baz"
			}
			filename_filter {
			  operation = "regex"
			  value = "spam"
			}
			context_filter {
			  operation = "startswith"
			  value = "eggs"
			}
			method_filter {
			  operation = "nregex"
			  value = "foo"
			}
			framework_filter {
			  operation = "eq"
			  value = "bar"
			}
		  }

		}
	`
	config := fmt.Sprintf(tmpl, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("not yet implemented"),
			},
		},
	})
}
