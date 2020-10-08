/*
 * `cleanup` is a utility for deleting orphaned Rollbar projects from failed
 * acceptance test runs.
 */
package main

import (
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

func main() {
	log.Info().Msg("Cleaning up orphaned Rollbar projects from failed acceptance test runs.")

	token := os.Getenv("ROLLBAR_TOKEN")
	c, err := client.NewClient(token)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	projects, err := c.ListProjects()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	for _, p := range projects {
		l := log.With().
			Str("name", p.Name).
			Int("id", p.Id).
			Logger()
		err = c.DeleteProject(p.Id)
		if err != nil {
			l.Fatal().Err(err).Send()
		}
		if strings.HasPrefix(p.Name, "tf-acc-test-") {
			l.Info().Msg("Deleted project")
		}
	}

	log.Info().Msg("Cleanup complete")
}
