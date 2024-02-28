/*
 * Copyright (c) 2022 Rollbar, Inc.
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestIntegrationCreate tests creating an integration
func (s *AccSuite) TestIntegrationCreate() {
	integrationResourceName := "rollbar_integration.webhook_integration"
	// language=hcl
	config := `
		resource "rollbar_integration" "webhook_integration" {
  webhook {
    enabled = true
    url     = "https://www.rollbar.com"
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
					s.checkResourceStateSanity(integrationResourceName),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook.0.url", "https://www.rollbar.com"),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook.0.enabled", "true"),
				),
			},
		},
	})
}

// TestIntegrationUpdate tests updating an integration
func (s *AccSuite) TestIntegrationUpdate() {
	integrationResourceName := "rollbar_integration.webhook_integration"
	// language=hcl
	config1 := `
	resource "rollbar_integration" "webhook_integration" {
  webhook {
    enabled = true
    url     = "https://www.rollbar.com"
  }
}
	`
	config2 := `
		resource "rollbar_integration" "webhook_integration" {
  webhook {
    enabled = false
    url     = "https://www.rollbar.com"
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
					s.checkResourceStateSanity(integrationResourceName),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook.0.url", "https://www.rollbar.com"),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook.0.enabled", "true"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(integrationResourceName),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook.0.url", "https://www.rollbar.com"),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook.0.enabled", "false"),
				),
			},
		},
	})
}
