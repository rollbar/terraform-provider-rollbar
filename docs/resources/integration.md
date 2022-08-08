`rollbar_integration` Resource
=========================

Rollbar projects can be configured with different notification integrations (aka "channels") and rules for when to send notifications to those integration platforms.  This resource manages the configuration for the integrations(limited to the Slack channel at the moment) for the project configured for the Rollbar provider.

This resource can manage configuration for the Slack channel. See the following api documentation for more details about the arguments with respect to the Slack channel:

* [Rollbar API Slack Integration](https://docs.rollbar.com/reference/Notification-Channels)

Example Usage
-------------

```hcl
# Set the rollbar provider to manage a project
#
# NOTE: the account access token `api_key` is not required
# for managing project however it should be set if you want
# to use the provider for managing account-level resources
# in addition to project-level resources

provider "rollbar" {
    api_key         = "my-account-access-token"  # optional for this resource
    project_api_key = "my-project-access-token"
}

# Configure the Slack integration for the project
#

resource "rollbar_integration" "slack_integration" {
  slack {
    enabled = false
    channel =  "#demo"
    service_account_id = "1234r45"
    show_message_buttons = true
  }
}
```

Argument Reference
------------------

The following arguments are supported:

Slack:

* `enabled` - (Required) Boolean that enables the Slack notifications globally
* `service_account_id` - (Required) The Rollbar Slack service account ID configured for the account.  You can find your service account ID [here](https://rollbar.com/settings/integrations/#slack)
* `channel` - (Required) The default Slack channel name to send the messages. Requires a `#` as a prefix
* `show_message_buttons` - Boolean that enables the Slack actionable buttons

Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - Composite ID of the project ID and the integration name
