package client

// Project is a Rollbar project
type Project struct {
	AccountId    int    `json:"account_id"`
	DateCreated  int    `json:"date_created"`
	DateModified int    `json:"date_modified"`
	Id           int    `json:"id"`
	Name         string `json:"name"`
	SettingsData struct {
		Grouping struct {
			AutoUpgrade    bool     `json:"auto_upgrade"`
			RecentVersions []string `json:"recent_versions"`
		} `json:"grouping"`
		Integrations struct {
			Asana       interface{} `json:"asana"`
			AzureDevops interface{} `json:"azuredevops"`
			Bitbucket   interface{} `json:"bitbucket"`
			/*
				"campfire": {},
				"ciscospark": {},
				"clubhouse": {},
				"datadog": {},
				"email": {
					"enabled": true
				},
				"flowdock": {},
				"github": {},
				"gitlab": {},
				"hipchat": {},
				"jira": {},
				"lightstep": {},
				"pagerduty": {},
				"pivotal": {},
				"slack": {},
				"sprintly": {},
				"trello": {},
				"victorops": {},
				"webhook": {}
			*/
		} `json:"integrations"`
		TimeFormat string `json:"time_format"`
		Timezone   string `json:"timezone"`
	} `json:"settings_data"`
	Status string `json:"status"`
}
