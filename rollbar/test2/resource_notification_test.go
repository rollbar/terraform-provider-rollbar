/*
 * Copyright (c) 2024 Rollbar, Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package test2

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

func init() {
	resource.AddTestSweepers("rollbar_notification", &resource.Sweeper{
		Name: "rollbar_notification",
		F:    sweepResourceNotification,
	})
}

// TestNotificationCreate tests creating a notification
func (s *AccSuite) TestNotificationCreate() {
	notificationResourceName := "rollbar_notification.webhook_notification"
	// language=hcl
	config := `
		resource "rollbar_notification" "webhook_notification" {
  rule  {
    filters {
        type =  "environment"
        operation =  "eq"
        value = "production"
    }
    filters {
       type = "framework"
       operation = "eq"
       value = 13
    }
   trigger = "new_item"
  }
  channel = "webhook"
  config  {
     url = "https://www.rollbar.com"
     format = "json"
  }
}
	`
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(notificationResourceName),
					resource.TestCheckResourceAttr(notificationResourceName, "channel", "webhook"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.filters.0.type", "environment"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.filters.1.type", "framework"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.trigger", "new_item"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.url", "https://www.rollbar.com"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.format", "json"),
				),
			},
		},
	})
}

// TestNotificationCreateDisabledRule tests creating a disbaled notification
func (s *AccSuite) TestNotificationCreateDisabledRule() {
	notificationResourceName := "rollbar_notification.webhook_notification"
	// language=hcl
	config := `
		resource "rollbar_notification" "webhook_notification" {
  rule  {
	enabled = false
    filters {
        type =  "environment"
        operation =  "eq"
        value = "production"
    }
    filters {
       type = "framework"
       operation = "eq"
       value = 13
    }
   trigger = "new_item"
  }
  channel = "webhook"
  config  {
     url = "https://www.rollbar.com"
     format = "json"
  }
}
	`
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(notificationResourceName),
					resource.TestCheckResourceAttr(notificationResourceName, "channel", "webhook"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.filters.0.type", "environment"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.filters.1.type", "framework"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.trigger", "new_item"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.enabled", "false"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.url", "https://www.rollbar.com"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.format", "json"),
				),
			},
		},
	})
}

// TestNotificationUpdate tests updating a notification
func (s *AccSuite) TestNotificationUpdate() {
	notificationResourceName := "rollbar_notification.webhook_notification"
	// language=hcl
	config1 := `
		resource "rollbar_notification" "webhook_notification" {
  rule  {
	enabled = true
    filters {
        type =  "environment"
        operation =  "eq"
        value = "production"
    }
    filters {
       type = "framework"
       operation = "eq"
       value = 13
    }
   trigger = "new_item"
  }
  channel = "webhook"
  config  {
     url = "https://www.rollbar.com"
     format = "json"
  }
}
	`
	config2 := `
		resource "rollbar_notification" "webhook_notification" {
  rule  {
 	enabled = false
    filters {
        type =  "environment"
        operation =  "eq"
        value = "production"
    }
    filters {
       type = "framework"
       operation = "eq"
       value = 13
    }
   trigger = "new_item"
  }
  channel = "webhook"
  config  {
     url = "https://www.rollbar.com"
     format = "xml"
  }
}
	`
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(notificationResourceName),
					resource.TestCheckResourceAttr(notificationResourceName, "channel", "webhook"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.filters.0.type", "environment"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.filters.1.type", "framework"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.trigger", "new_item"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.url", "https://www.rollbar.com"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.format", "json"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(notificationResourceName),
					resource.TestCheckResourceAttr(notificationResourceName, "channel", "webhook"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.filters.0.type", "environment"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.filters.1.type", "framework"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.trigger", "new_item"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.enabled", "false"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.url", "https://www.rollbar.com"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.format", "xml"),
				),
			},
		},
	})
}

func (s *AccSuite) TestNotificationCreateSpecialEmail() {
	notificationResourceName := "rollbar_notification.email_notification"
	// language=hcl
	config := `
       resource "rollbar_notification" "email_notification" {
          rule  {
             trigger = "daily_summary"
          }
          channel = "email"
          config  {
	         summary_time = 2
	         min_item_level = "critical"
             send_only_if_data = true
			 environments = ["production", "staging"]
          }
       }`

	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(notificationResourceName),
					resource.TestCheckResourceAttr(notificationResourceName, "channel", "email"),
					resource.TestCheckResourceAttr(notificationResourceName, "rule.0.trigger", "daily_summary"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.summary_time", "2"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.min_item_level", "critical"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.send_only_if_data", "true"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.environments.0", "production"),
					resource.TestCheckResourceAttr(notificationResourceName, "config.0.environments.1", "staging"),
				),
			},
		},
	})
}

// sweepResourceNotification cleans up notifications
func sweepResourceNotification(_ string) error {
	log.Info().Msg("Cleaning up Rollbar notifications from acceptance test runs.")

	c := client.NewClient(client.DefaultBaseURL, os.Getenv("ROLLBAR_PROJECT_API_KEY"))
	notifications, err := c.ListNotifications("webhook")
	if err != nil {
		log.Err(err).Send()
		return err
	}
	for _, n := range notifications {
		err = c.DeleteNotification(n.ID, "webhook")
		if err != nil {
			log.Err(err).Send()
			return err
		}
	}

	notifications, err = c.ListNotifications("email")
	if err != nil {
		log.Err(err).Send()
		return err
	}
	for _, n := range notifications {
		err = c.DeleteNotification(n.ID, "email")
		if err != nil {
			log.Err(err).Send()
			return err
		}
	}
	return nil
}
