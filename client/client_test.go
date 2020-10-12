package client

import (
	"bytes"
	"github.com/rs/zerolog/log"
)

// TestClientNoToken checks that a warning message is logged when a
// RollbarApiClient is initialized without an API token.
func (s *Suite) TestClientNoToken() {
	var buf bytes.Buffer
	log.Logger = log.Logger.Output(&buf)
	NewClient("") // Valid, but probably not what you want, thus warn
	bs := buf.String()
	s.NotZero(bs)
	s.Contains(bs, "warn")
	s.Contains(bs, "Rollbar API token not set")
}
