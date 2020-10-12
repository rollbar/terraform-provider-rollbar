package rollbar

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// TestAccRollbarProjectsDataSource tests listing of all projects with
// `rollbar_projects` data source.
func TestAccRollbarProjectsDataSource(t *testing.T) {

	rn := "data.rollbar_projects.all"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-acc-test-%s", randString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRollbarProjectsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "projects.#"),
					resource.TestCheckResourceAttr(rn, "projects.#", "1"),
					resource.TestCheckResourceAttr(rn, "projects.0.name", name),
				),
			},
		},
	})
}

func testAccRollbarProjectsDataSourceConfig(projName string) string {
	return fmt.Sprintf(`
		resource "rollbar_project" "test" {
		  name         = "%s"
		}
		
		data "rollbar_projects" "all" {
			depends_on = [rollbar_project.test]
		}
	`, projName)
}
