`rollbar_notification` Resource
=========================

Rollbar projects can be configured with different notification integrations (aka "channels") and rules for when to send notifications to those integration platforms.  This resource manages the rules for the project configured for the Rollbar provider.  The notification channels are enabled and disabled through the Rollbar UI.

This resource can manage notification rules for different integration channels.  See the following api documentation for more details about the arguments with respect to each channel:

* [Rollbar API Slack Notification Rules](https://docs.rollbar.com/reference/slack-notification-rules)
* [Rollbar API Pagerduty Notification Rules](https://docs.rollbar.com/reference/pagerduty-notification-rules)
* [Rollbar API Email Notification Rules](https://docs.rollbar.com/reference/email-notification-rules)
* [Rollbar API Webhook Notification Rules](https://docs.rollbar.com/reference/webhook-notification-rules)


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

# Create an email notification for the project
#
resource "rollbar_notification" "email" {
  channel = "email"
  rule {
    enabled = true
    trigger = "occurrence_rate"
    filters {
      type   = "rate"
      period = 300
      count  = 10
    }
  }
  config {
    users = ["travis.mattera@rollbar.com"]
    teams = ["test-team-example"]
  }
}
```

Argument Reference
------------------

The following arguments are supported:

* `channel` - (Required) The notification channel (eg. `slack`, `pagerduty`, `email`, `webhook`) to configure a notification rule(s) for
* `rule` - (Required) An array of expression configurations for notification rules.  Structure is [documented below](#nested_rule)
* `config` - (Required) An array of configurations for notification rules.  Structure is [documented below](#nested_config)

<a name="nested_rule"></a>The `rule` block supports:
* `enabled` - (Optional) Boolean that enables the rule notification. The default value is `true`.
* `trigger` - (Required) The category of trigger evaluations using the expressions defined in filters block(s).
* `filters` - (Required) One or more nested configuration blocks that define filter expressions.  Structure is [documented below](#nested_filters)

<a name="nested_filters"></a>The `filters` block supports:
* `path` - json path (body.field1.field2)
* `type` - (Required) The type of filter expression.
* `operation` - The comparator used in the expression evalution for the filter.
* `value` - The value to compare the triggering metric against.
* `period` - The period of time in seconds.  Allowed values `300`, `1800`, `3600`, `86400`, `60`.
* `count` - The number of distinct items or occurrences used as a threshold for the filter evaluation.


<a name="nested_config"></a>The `config` block supports:

* `users` - (Required only for Email)  A list of users to notify.
* `teams` - (Required only for Email)  A list of teams to notify.
* `message_template` - (Required only for Slack)  A template for posting messages to a Slack channel.
* `channel` - (Required only for Slack)  The Slack channel to post messages to.
* `show_message_buttons` - (Required only for Slack)  Boolean value to toggle message buttons on/off in Slack.
* `service_key` - (Required only for PagerDuty)  The Pagerduty service API key.
* `url` - (Required only for Webhook)  The Webhook URL.
* `format` - (Required only for Webhook)  The Webhook format (json or xml).

Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the notification rule


Import
------

Projects can be imported using the notification channel and ID separated by a comma, e.g.

```
$ terraform import rollbar_notification.foo email,857623
```
