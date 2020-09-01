/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package client

import "fmt"

// Project is a Rollbar project
type Project struct {
	AccountId    int    `json:"account_id" model:"account_id"`
	DateCreated  int    `json:"date_created" model:"date_created"`
	DateModified int    `json:"date_modified" model:"date_modified"`
	Id           int    `json:"id" model:"id"`
	Name         string `json:"name" model:"name"`
	//SettingsData struct {
	//	Grouping struct {
	//		AutoUpgrade    bool     `json:"auto_upgrade"`
	//		RecentVersions []string `json:"recent_versions"`
	//	} `json:"grouping"`
	//	Integrations struct {
	//		Asana       interface{} `json:"asana"`
	//		AzureDevops interface{} `json:"azuredevops"`
	//		Bitbucket   interface{} `json:"bitbucket"`
	//		/*
	//			"campfire": {},
	//			"ciscospark": {},
	//			"clubhouse": {},
	//			"datadog": {},
	//			"email": {
	//				"enabled": true
	//			},
	//			"flowdock": {},
	//			"github": {},
	//			"gitlab": {},
	//			"hipchat": {},
	//			"jira": {},
	//			"lightstep": {},
	//			"pagerduty": {},
	//			"pivotal": {},
	//			"slack": {},
	//			"sprintly": {},
	//			"trello": {},
	//			"victorops": {},
	//			"webhook": {}
	//		*/
	//	} `json:"integrations"`
	//	TimeFormat string `json:"time_format"`
	//	Timezone   string `json:"timezone"`
	//} `json:"settings_data"`
	Status string `json:"status" model:"status"`
}

type ProjectListResult struct {
	Err     int
	Message string
	Result  []Project
}

type ProjectResult struct {
	Err    int     `json:"err"`
	Result Project `json:"result"`
}

type ErrorResult struct {
	Err     int    `json:"err"`
	Message string `jason:"message"`
}

func (er ErrorResult) Error() string {
	return fmt.Sprintf("%v %v", er.Err, er.Message)
}
