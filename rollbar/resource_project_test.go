/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package rollbar

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
	"testing"
)

// TestAccRollbarProject tests creation and deletion of a Rollbar project.
func TestAccRollbarProject(t *testing.T) {

	rn := "rollbar_project.foo"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-acc-test-%s", randString)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRollbarProjectConfig(randString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "name", name),
					testAccRollbarProjectExists(rn, name),
					testAccRollbarProjectInProjectList(rn),
				),
			},
		},
	})
}
func testAccRollbarProjectConfig(randString string) string {
	return fmt.Sprintf(`
		resource "rollbar_project" "foo" {
		  name         = "tf-acc-test-%s"
		}
	`, randString)
}

// testAccRollbarProjectExists tests that the newly created project exists
func testAccRollbarProjectExists(rn string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check terraform config is sane
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not Found: %s", rn)
		}
		idString := rs.Primary.ID
		if idString == "" {
			return fmt.Errorf("No project ID is set")
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return err
		}

		// Check that project exists
		c := testAccProvider.Meta().(*client.RollbarApiClient)
		proj, err := c.ReadProject(id)
		if err != nil {
			return err
		}
		if proj.Name != name {
			return fmt.Errorf("project name from API does not match project name in Terraform config")
		}

		// Success
		return nil
	}
}

// testAccRollbarProjectInProjectList tests that the newly created project is
// present in the list of all projects.
func testAccRollbarProjectInProjectList(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check terraform config is sane
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not Found: %s", rn)
		}
		idString := rs.Primary.ID
		if idString == "" {
			return fmt.Errorf("No project ID is set")
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			return err
		}

		// Check that project exists
		c := testAccProvider.Meta().(*client.RollbarApiClient)
		projList, err := c.ListProjects()
		if err != nil {
			return err
		}
		found := false
		for _, proj := range projList {
			if proj.Id == id {
				found = true
			}
		}
		if !found {
			msg := "Project not found in project list"
			log.Debug().Int("id", id).Msg(msg)
			return fmt.Errorf(msg)
		}

		// Success
		return nil
	}
}
