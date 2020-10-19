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

package client

// Path components of Rollbar API URLs
const (
	pathProjectCreate = "/api/1/projects"
	pathProjectDelete = "/api/1/project/{projectId}"
	pathProjectList   = "/api/1/projects"
	pathProjectRead   = "/api/1/project/{projectId}"
	pathPatList       = "/api/1/project/{projectId}/access_tokens"
	pathPatCreate     = "/api/1/project/{projectId}/access_tokens"
	pathPatUpdate     = "/api/1/project/{projectId}/access_token/{accessToken}"
	pathTeamCreate    = "/api/1/teams"
	pathTeamList      = "/api/1/teams"
	pathTeamRead      = "/api/1/team/{teamId}"
	pathTeamDelete    = "/api/1/team/{teamId}"
	pathInviteCreate  = "/api/1/team/{teamId}/invites"
)
