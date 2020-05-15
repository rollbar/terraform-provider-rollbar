package rollbar

import "encoding/json"

// Project represents a project
type Project struct {
	AccountID    int    `json:"account_id"`
	ID           int    `json:"id"`
	DateCreated  int    `json:"date_created"`
	DateModified int    `json:"date_modified"`
	Name         string `json:"name"`
}

// listProjectsResponse represents the list projects response
type listProjectsResponse struct {
	Error  int `json:"err"`
	Result []*Project
}

// ListProjects lists the projects for this API Key
// https://docs.rollbar.com/reference#list-all-projects
func (c *Client) ListProjects() ([]*Project, error) {
	var data listProjectsResponse

	bytes, err := c.get("projects")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return nil, err
	}

	projects := []*Project{}
	projects = append(projects, data.Result...)

	return projects, nil
}

// GetProjectByName returns the first project from the list-projects
// call that matches a given name. If there is no matching project
// returns nil.
func (c *Client) GetProjectByName(name string) (*Project, error) {
	projects, err := c.ListProjects()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == name {
			return project, nil
		}
	}

	return nil, nil
}
