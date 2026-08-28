`rollbar_integration` Resource
=========================

Rollbar projects can be configured with different notification integrations (aka "channels") and rules for when to send notifications to those integration platforms.  This resource manages the configuration for the integrations(limited to the Slack, Webhook, Email and PagerDuty channels at the moment) for the project configured for the Rollbar provider.

This resource can manage configuration for the Slack, Webhook, Email and PagerDuty channels. See the following api documentation for more details about the arguments with respect to these channels:

* [Rollbar API Slack Integration](https://docs.rollbar.com/reference/put_api-1-notifications-slack), [Rollbar Slack Integration Documentation](https://docs.rollbar.com/docs/slack)
* [Rollbar API Webhook Integration](https://docs.rollbar.com/reference/put_api-1-notifications-webhook), [Rollbar Webhook Integration Documentation](https://docs.rollbar.com/docs/webhook)
* [Rollbar API PagerDuty Integration](https://docs.rollbar.com/reference/put_api-1-notifications-pagerduty), [Rollbar PagerDuty Integration Documentation](https://docs.rollbar.com/docs/pagerduty)
* [Rollbar API Email Integration](https://docs.rollbar.com/reference/put_api-1-notifications-email)

Example Usage
-------------

**Option 1: Using a project access token**

```hcl
provider "rollbar" {
    api_key         = "my-account-access-token"  # optional for this resource
    project_api_key = "my-project-access-token"
}

resource "rollbar_integration" "slack_integration" {
  slack {
    enabled = false
    channel =  "#demo"
    service_account_id = "1234r45"
    show_message_buttons = true
  }
}

resource "rollbar_integration" "webhook_integration" {
  webhook {
    enabled = true
    url = "https://www.example.com"
  }
}

resource "rollbar_integration" "email_integration" {
  email {
    enabled = true
    scrub_params = true
  }
}

resource "rollbar_integration" "pagerduty_integration" {
  pagerduty {
    enabled = false
    service_key = "123456789"
  }
}
```

**Option 2: Using an account access token with `project_id`**

If you prefer not to configure a separate project access token, you can use the account access token and specify `project_id` on each resource:

```hcl
provider "rollbar" {
    api_key = "my-account-access-token"
}

resource "rollbar_integration" "slack_integration" {
  project_id = 123456
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

* `project_id` - (Optional) The Rollbar project ID. When set, the account access token (`api_key`) is used with this project ID instead of requiring a separate project access token (`project_api_key`).

Slack:

* `enabled` - (Required) Boolean that enables the Slack notifications globally
* `service_account_id` - (Required) The Rollbar Slack service account ID configured for the account.  You can find your service account ID [here](https://rollbar.com/settings/integrations/#slack)
* `channel` - (Required) The default Slack channel name to send the messages. Requires a `#` as a prefix
* `show_message_buttons` - Boolean that enables the Slack actionable buttons

Webhook:

* `enabled` - (Required) Boolean that enables the Webhook notifications globally
* `url` - (Required) URL for the webhook.

Email:

* `enabled` - (Required) Boolean that enables the Email notifications globally
* `scrub_params` - Optional Boolean that enables scrubbing param values (when set to true)

PagerDuty:

* `enabled` - (Required) Boolean that enables the PagerDuty notifications globally
* `service_key` - (Required) PagerDuty service key linked to PagerDuty account

Attribute Reference
-------------------

In addition to all arguments above, the following attributes are exported:

* `id` - Composite ID of the project ID and the integration name
