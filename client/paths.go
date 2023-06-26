/*
 * Copyright (c) 2023 Rollbar, Inc.
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

package client

// Path components of Rollbar API URLs
const (
	pathProjectCreate                    = "/api/1/projects"
	pathProjectDelete                    = "/api/1/project/{projectID}"
	pathProjectList                      = "/api/1/projects"
	pathProjectRead                      = "/api/1/project/{projectID}"
	pathProjectToken                     = "/api/1/project/{projectID}/access_token/{accessToken}"
	pathProjectTokens                    = "/api/1/project/{projectID}/access_tokens"
	pathProjectTeams                     = "/api/1/project/{projectID}/teams"
	pathTeamCreate                       = "/api/1/teams"
	pathTeamRead                         = "/api/1/team/{teamID}"
	pathTeamList                         = "/api/1/teams"
	pathTeamDelete                       = "/api/1/team/{teamID}"
	pathTeamUser                         = "/api/1/team/{teamID}/user/{userID}"
	pathTeamProject                      = "/api/1/team/{teamID}/project/{projectID}"
	pathUser                             = "/api/1/user/{userID}"
	pathUserTeams                        = "/api/1/user/{userID}/teams"
	pathUsers                            = "/api/1/users"
	pathInvitation                       = "/api/1/invite/{inviteID}"
	pathInvitations                      = "/api/1/invites"
	pathTeamInvitations                  = "/api/1/team/{teamID}/invites"
	pathNotificationCreate               = "/api/1/notifications/{channel}/rules"
	pathNotificationReadOrDeleteOrUpdate = "/api/1/notifications/{channel}/rule/{notificationID}"
	pathServiceLinkCreate                = "/api/1/service_links"
	pathServiceLinkReadOrDeleteOrUpdate  = "/api/1/service_links/{id}"
	pathIntegration                      = "/api/1/notifications/{integration}"
)
