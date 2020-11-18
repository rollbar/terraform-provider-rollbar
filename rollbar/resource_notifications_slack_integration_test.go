package rollbar

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
)

// TestAccNotificationsSlackIntegrationCreateNotImplemented tests creating a
// Rollbar notifications Slack integration - it expects an error because this
// resource's read method is not yet implemented.
func (s *AccSuite) TestAccNotificationsSlackIntegrationCreateNotImplemented() {
	// language=hcl
	tmpl := `
		resource "rollbar_project" "test" {
			name = "%s"
		}

		resource "rollbar_project_access_token" "test" {
			project_id = rollbar_project.test.id
			name = "test_%s"
			scopes = ["write"]
		}

		resource "rollbar_notifications_slack_integration" "singleton" {
			enabled = true
			project_access_token = rollbar_project_access_token.test.access_token
			show_message_buttons = true
		}
	`
	config := fmt.Sprintf(tmpl, s.randName, s.randName)
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
