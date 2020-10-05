/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// EnableRestyZerolog configures a Resty client to use zerolog for output
func EnableRestyZerolog(c *resty.Client) {
	l := log.Logger
	rzl := RestyZeroLogger{l}
	c.SetLogger(rzl)
}

// RestyZeroLogger implements resty.Logger on top of zerolog.Logger
type RestyZeroLogger struct {
	zl zerolog.Logger
}

func (r RestyZeroLogger) Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	r.zl.Error().Msg(msg)
}
func (r RestyZeroLogger) Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	r.zl.Warn().Msg(msg)
}

func (r RestyZeroLogger) Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	r.zl.Debug().Msg(msg)
}
