package rollbar

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
)

func (s *AccSuite) TestNotificationEmailNotImplemented() {

	// TestAccUserCreateInvite tests creating a new rollbar_user resource with an
	// invitation to email is not registered as a Rollbar user.
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

		resource "rollbar_notification_email" "foo" {
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
