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

/*
 * Package rollbar implements a Terraform provider for the Rollbar API.
 */
package rollbar

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
)

const schemaKeyToken = "api_key"

// Provider is a Terraform provider for Rollbar.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			schemaKeyToken: {
				Type:     schema.TypeString,
				Optional: true,
				// FIXME: Should the environment variable be ROLLBAR_API_KEY to
				//  match the name of this field?
				DefaultFunc: schema.EnvDefaultFunc("ROLLBAR_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rollbar_project":              resourceProject(),
			"rollbar_project_access_token": resourceProjectAccessToken(),
			"rollbar_team":                 resourceTeam(),
			"rollbar_user":                 resourceUser(),
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
	c := client.NewClient(token)
	c.Resty.GetClient().Transport = http.DefaultTransport
	return c, diags
}

// handleErrNotFound handles an ErrNotFound when reading a resource, by removing
// the resource from state and returning a Diagnostics object.
func handleErrNotFound(d *schema.ResourceData, resourceName string) diag.Diagnostics {
	id := d.Id()
	d.SetId("")
	tmpl := `Removing %s %s from state because it was not found on Rollbar`
	detail := fmt.Sprintf(tmpl, resourceName, id)
	log.Warn().Msg(detail)
	tmpl = "%s not found, removed from state"
	summary := fmt.Sprintf(tmpl, strings.ToTitle(resourceName))
	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  summary,
		Detail:   detail,
	}}
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
