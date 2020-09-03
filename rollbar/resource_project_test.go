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
	"github.com/jmcvetta/terraform-provider-rollbar/client"
	"testing"
)

func TestAccRollbarProject(t *testing.T) {

	rn := "rollbar_project.foo"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-acc-test-%s", randString)
	//description := fmt.Sprintf("Terraform acceptance tests %s", randString)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRollbarProjectConfig(randString),
				Check:  resource.ComposeTestCheckFunc(),
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

func testAccRollbarProjectExists(rn string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Not Found: %s", rn)
		}

		repoName := rs.Primary.ID
		if repoName == "" {
			return fmt.Errorf("No repository name is set")
		}

		c := testAccProvider.(*client.RollbarApiClient)

		org := testAccProvider.Meta().(*Owner)
		conn := org.v3client
		gotRepo, _, err := conn.Repositories.Get(context.TODO(), org.name, repoName)
		if err != nil {
			return err
		}
		*repo = *gotRepo
		return nil
	}
}
