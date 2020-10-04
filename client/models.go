/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package client

import "fmt"

type ErrorResult struct {
	Err     int    `json:"err"`
	Message string `jason:"message"`
}

func (er ErrorResult) Error() string {
	return fmt.Sprintf("%v %v", er.Err, er.Message)
}
