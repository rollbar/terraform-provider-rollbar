/*
 * Copyright (c) 2020 Rollbar, Inc.
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
	c := client.NewClient(token)

	projects, err := c.ListProjects()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	for _, p := range projects {
		l := log.With().
			Str("name", p.Name).
			Int("id", p.Id).
			Logger()
		if strings.HasPrefix(p.Name, "tf-acc-test-") {
			err = c.DeleteProject(p.Id)
			if err != nil {
				l.Fatal().Err(err).Send()
			}
			l.Info().Msg("Deleted project")
		}
	}

	log.Info().Msg("Cleanup complete")
}
