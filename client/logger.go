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
	"github.com/rs/zerolog"
)

// restyZeroLogger implements resty.Logger on top of zerolog.Logger
type restyZeroLogger struct {
	zl zerolog.Logger
}

func (r restyZeroLogger) Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	r.zl.Error().Msg(msg)
}
func (r restyZeroLogger) Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	r.zl.Warn().Msg(msg)
}

func (r restyZeroLogger) Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	r.zl.Debug().Msg(msg)
}
