/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package client

// Path components of Rollbar API URLs
const (
	pathProjectCreate = "/api/1/projects"
	pathProjectDelete = "/api/1/project/{projectId}"
	pathProjectList   = "/api/1/projects"
	pathProjectRead   = "/api/1/project/{projectId}"
	pathPatList       = "/api/1/project/{projectId}/access_tokens"
	pathPatCreate     = "/api/1/project/{projectId}/access_tokens"
	pathTeamCreate    = "/api/1/teams"
	pathTeamList      = "/api/1/teams"
	pathTeamRead      = "/api/1/team/{teamId}"
	pathTeamDelete    = "/api/1/team/{teamId}"
)
