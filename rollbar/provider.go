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

// Package rollbar implements a Terraform provider for the Rollbar API.
package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"strconv"
)

const schemaKeyToken = "api_key"
const projectKeyToken = "project_api_key"
const schemaKeyBaseURL = "api_url"

// Provider is a Terraform provider for Rollbar.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			schemaKeyToken: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ROLLBAR_API_KEY", nil),
				Description: "Rollbar API authentication token. Value will be sourced from environment variable `ROLLBAR_API_KEY` if set.",
			},
			projectKeyToken: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ROLLBAR_PROJECT_API_KEY", nil),
				Description: "Rollbar API authentication token (project level). Value will be sourced from environment variable `ROLLBAR_PROJECT_API_KEY` if set.",
			},
			schemaKeyBaseURL: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ROLLBAR_API_URL", client.DefaultBaseURL),
				Description: "Base URL for the Rollbar API.  Defaults to https://api.rollbar.com.  Value will be sourced from environment variable `ROLLBAR_API_URL` if set.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rollbar_project":              resourceProject(),
			"rollbar_project_access_token": resourceProjectAccessToken(),
			"rollbar_team":                 resourceTeam(),
			"rollbar_user":                 resourceUser(),
			"rollbar_notification":         resourceNotification(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"rollbar_project":               dataSourceProject(),
			"rollbar_projects":              dataSourceProjects(),
			"rollbar_project_access_token":  dataSourceProjectAccessToken(),
			"rollbar_project_access_tokens": dataSourceProjectAccessTokens(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure sets up authentication in a Resty HTTP client.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	token := d.Get(schemaKeyToken).(string)
	projectToken := d.Get(projectKeyToken).(string)
	baseURL := d.Get(schemaKeyBaseURL).(string)
	c := client.NewClient(baseURL, token)
	pc := client.NewClient(baseURL, projectToken)
	return []*client.RollbarAPIClient{c, pc}, diags
}

/*

// errSetter sets Terraform state values until an error occurs, whereupon it
// becomes a no-op but preserves the error value.
// Based on Rob Pike's errWriter - https://blog.golang.org/errors-are-values
type errSetter struct {
	d   *schema.ResourceData
	err error
}

func (es *errSetter) Set(key string, value interface{}) {
	if es.err != nil {
		return
	}
	es.err = es.d.Set(key, value)
}

*/

// mustSet sets a value for a key in a schema, or panics on error.
func mustSet(d *schema.ResourceData, key string, value interface{}) {
	err := d.Set(key, value)
	if err != nil {
		panic(err)
	}
}

// mustGetID gets the ID of the resource as an integer, or panics if string ID
// value cannot be cast to int.
func mustGetID(d *schema.ResourceData) int {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		panic(err)
	}
	return id
}

// Decode takes an input structure and uses reflection to translate it to the
// output structure, panicking on error. Output must be a pointer to a map or
// struct.
func mustDecodeMapStructure(input, output interface{}) {
	err := mapstructure.Decode(input, &output)
	if err != nil {
		panic(err)
	}
}
