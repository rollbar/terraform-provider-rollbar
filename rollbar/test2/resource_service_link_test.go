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
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

func init() {
	resource.AddTestSweepers("rollbar_service_link", &resource.Sweeper{
		Name: "rollbar_service_link",
		F:    sweepResourceNotification,
	})
}

// TestServiceLinkCreate tests creating a service link
func (s *AccSuite) TestServiceLinkCreate() {
	serviceLinkResourceName := "rollbar_service_link.service_link"
	// language=hcl
	tmpl := `
		resource "rollbar_service_link" "service_link" {
          name = "%s"
     	  template = "sometemplate_new.{{ss}}"
       }
	`

	config := fmt.Sprintf(tmpl, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(serviceLinkResourceName),
					resource.TestCheckResourceAttr(serviceLinkResourceName, "name", s.randName),
					resource.TestCheckResourceAttr(serviceLinkResourceName, "template", "sometemplate_new.{{ss}}"),
				),
			},
		},
	})
}

// TestServiceLinkUpdate tests updating a service link
func (s *AccSuite) TestServiceLinkUpdate() {
	serviceLinkResourceName := "rollbar_service_link.service_link"
	// language=hcl
	tmpl1 := `
		resource "rollbar_service_link" "service_link" {
          name = "%s"
     	  template = "sometemplate.{{ss}}"
       }
	`
	config1 := fmt.Sprintf(tmpl1, s.randName)
	tmpl2 := `
		resource "rollbar_service_link" "service_link" {
          name = "%s"
     	  template = "sometemplate_new.{{ss}}"
       }
	`
	config2 := fmt.Sprintf(tmpl2, s.randName)
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck:     func() { s.preCheck() },
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(serviceLinkResourceName),
					resource.TestCheckResourceAttr(serviceLinkResourceName, "name", s.randName),
					resource.TestCheckResourceAttr(serviceLinkResourceName, "template", "sometemplate.{{ss}}"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					s.checkResourceStateSanity(serviceLinkResourceName),
					resource.TestCheckResourceAttr(serviceLinkResourceName, "name", s.randName),
					resource.TestCheckResourceAttr(serviceLinkResourceName, "template", "sometemplate_new.{{ss}}"),
				),
			},
		},
	})
}

// sweepResourceServiceLink cleans up service links.
func sweepResourceServiceLink(_ string) error {
	log.Info().Msg("Cleaning up Rollbar service links from acceptance test runs.")

	c := client.NewClient(client.DefaultBaseURL, os.Getenv("ROLLBAR_PROJECT_API_KEY"))
	serviceLinks, err := c.ListSerivceLinks()
	if err != nil {
		log.Err(err).Send()
		return err
	}
	for _, s := range serviceLinks {
		err = c.DeleteServiceLink(s.ID)
		if err != nil {
			log.Err(err).Send()
			return err
		}
	}
	return nil
}
