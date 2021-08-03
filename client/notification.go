/*
 * Copyright (c) 2021 Rollbar, Inc.
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

import (
	"github.com/rs/zerolog/log"
	"strings"
)

type Notification struct {
	ID int `model:"id" mapstructure:"id"`
}

// CreateOrUpdateNotification creates a new Rollbar notification.
func (c *RollbarAPIClient) CreateOrUpdateNotification(channel string, filters, trigger, config interface{}) (*Notification, error) {
	u := c.BaseURL + pathNotificationCreate
	u = strings.ReplaceAll(u, "{channel}", channel)
	l := log.With().
		Str("channel", channel).
		Logger()
	l.Debug().Msg("Creating new notification")

	resp, err := c.Resty.R().
		SetBody([]map[string]interface{}{{"filters": filters, "trigger": trigger, "config": config}}).
		SetResult(notificationResponse{}).
		SetError(ErrorResult{}).
		Put(u)
	if err != nil {
		l.Err(err).Msg("Error creating notification")
		return nil, err
	}
	err = errorFromResponse(resp)
	if err != nil {
		l.Err(err).Send()
		return nil, err
	}
	l.Debug().Msg("Notification successfully created")
	pr := resp.Result().(*notificationResponse)
	return &pr.Result, nil

}

//// ReadProject a Rollbar project from the API. If no matching project is found,
//// returns error ErrNotFound.
//func (c *RollbarAPIClient) ReadProject(projectID int) (*Project, error) {
//	u := c.BaseURL + pathProjectRead
//
//	l := log.With().
//		Int("projectID", projectID).
//		Logger()
//	l.Debug().Msg("Reading project from API")
//
//	resp, err := c.Resty.R().
//		SetResult(projectResponse{}).
//		SetError(ErrorResult{}).
//		SetPathParams(map[string]string{
//			"projectID": strconv.Itoa(projectID),
//		}).
//		Get(u)
//	if err != nil {
//		l.Err(err).Msg("Error reading project")
//		return nil, err
//	}
//	err = errorFromResponse(resp)
//	if err != nil {
//		l.Err(err).Send()
//		return nil, err
//	}
//	pr := resp.Result().(*projectResponse)
//	// FIXME: This is a workaround for a known bug in the API
//	//  https://github.com/rollbar/terraform-provider-rollbar/issues/23
//	if pr.Result.Name == "" {
//		l.Warn().Msg("Project not found")
//		return nil, ErrNotFound
//	}
//	l.Debug().Msg("Project successfully read")
//	return &pr.Result, nil
//
//}
//
//// DeleteProject deletes a Rollbar project. If no matching project is found,
//// returns error ErrNotFound.
//func (c *RollbarAPIClient) DeleteProject(projectID int) error {
//	u := c.BaseURL + pathProjectDelete
//	l := log.With().
//		Int("projectID", projectID).
//		Logger()
//	l.Debug().Msg("Deleting project")
//
//	resp, err := c.Resty.R().
//		SetError(ErrorResult{}).
//		SetPathParams(map[string]string{
//			"projectID": strconv.Itoa(projectID),
//		}).
//		Delete(u)
//	if err != nil {
//		l.Err(err).Msg("Error deleting project")
//		return err
//	}
//	err = errorFromResponse(resp)
//	if err != nil {
//		l.Err(err).Send()
//		return err
//	}
//	l.Debug().Msg("Project successfully deleted")
//	return nil
//}
//
//// FindProjectTeamIDs finds IDs of all teams assigned to the project. Caution:
//// this is a potentially slow operation that makes multiple calls to the API.
//// https://github.com/rollbar/terraform-provider-rollbar/issues/104
//func (c *RollbarAPIClient) FindProjectTeamIDs(projectID int) ([]int, error) {
//	l := log.With().Int("project_id", projectID).Logger()
//	l.Debug().Msg("Finding teams assigned to project")
//	var projectTeamIDs []int
//
//	allTeams, err := c.ListCustomTeams()
//	if err != nil {
//		l.Err(err).Send()
//		return nil, err
//	}
//	for _, t := range allTeams {
//		teamID := t.ID
//		projectIDs, err := c.ListTeamProjectIDs(teamID)
//		if err != nil && err != ErrNotFound {
//			l.Err(err).Send()
//			return nil, err
//		}
//		for _, id := range projectIDs {
//			if id == projectID {
//				projectTeamIDs = append(projectTeamIDs, teamID)
//			}
//		}
//	}
//	count := len(projectTeamIDs)
//	l.Debug().
//		Int("team_count", count).
//		Msg("Successfully found teams assigned to project")
//	return projectTeamIDs, nil
//}
//
//// UpdateProjectTeams updates the Rollbar teams assigned to a project, assigning
//// and removing teams as necessary. Caution: this is a potentially slow
//// operation that makes multiple calls to the API.
//// https://github.com/rollbar/terraform-provider-rollbar/issues/104
//func (c *RollbarAPIClient) UpdateProjectTeams(projectID int, teamIDs []int) error {
//	l := log.With().
//		Int("project_id", projectID).
//		Ints("team_ids", teamIDs).
//		Logger()
//	l.Debug().Msg("Updating teams for project")
//
//	// Compute which teams to assign and to remove
//	var assignTeamIDs, removeTeamIDs []int
//	currentTeamIDs, err := c.FindProjectTeamIDs(projectID) // Potential slowness is here
//	if err != nil {
//		l.Err(err).Send()
//		return err
//	}
//	current := make(map[int]bool)
//	for _, id := range currentTeamIDs {
//		current[id] = true
//	}
//	desired := make(map[int]bool)
//	for _, id := range teamIDs {
//		desired[id] = true
//	}
//	for id := range current {
//		if !desired[id] {
//			removeTeamIDs = append(removeTeamIDs, id)
//		}
//	}
//	for id := range desired {
//		if !current[id] {
//			assignTeamIDs = append(assignTeamIDs, id)
//		}
//	}
//	l.Debug().
//		Ints("assign_team_ids", assignTeamIDs).
//		Ints("remove_team_ids", removeTeamIDs).
//		Msg("Teams to assign and remove")
//
//	for _, teamID := range assignTeamIDs {
//		err = c.AssignTeamToProject(teamID, projectID)
//		if err != nil {
//			l.Err(err).Send()
//			return err
//		}
//	}
//	for _, teamID := range removeTeamIDs {
//		err = c.RemoveTeamFromProject(teamID, projectID)
//		if err != nil {
//			l.Err(err).Send()
//			return err
//		}
//	}
//	return nil
//}
//
///*
// * Containers for unmarshalling API responses
// */
//
//type projectListResponse struct {
//	Err    int       `json:"err"`
//	Result []Project `json:"result"`
//}
//
//type projectResponse struct {
//	Err    int     `json:"err"`
//	Result Project `json:"result"`
//}

type notificationResponse struct {
	Err    int          `json:"err"`
	Result Notification `json:"result"`
}
