package rollbar_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccTeam tests CRUD operations for a Rollbar team.
func (s *AccSuite) TestAccUser() {
	resource.ParallelTest(s.T(), resource.TestCase{
		PreCheck: func() { s.preCheck() },
		//ProviderFactories: testAccProviderFactories(),
		Providers:    s.providers,
		CheckDestroy: nil,
		Steps:        []resource.TestStep{
			//// Invalid name - failure expected
			//{
			//	Config:      s.configResourceUserInvalidname(),
			//	ExpectError: regexp.MustCompile("name cannot be blank"),
			//},
			//
			//// Initial create
			//{
			//	Config: s.configResourceUser(teamName0),
			//	Check: resource.ComposeTestCheckFunc(
			//		s.checkResourceStateSanity(rn),
			//		resource.TestCheckResourceAttr(rn, "name", teamName0),
			//		s.checkUser(rn, teamName0, "standard"),
			//	),
			//},

		},
	})
}
