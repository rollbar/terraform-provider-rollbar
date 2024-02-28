/*
 * Copyright (c) 2022 Rollbar, Inc.
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

package rollbar

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

// dataSourceProjectAccessToken is a data source returning a named access token
// belonging to a Rollbar project.
func dataSourceProjectAccessToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectAccessTokenRead,

		Schema: map[string]*schema.Schema{
			// Required fields
			"project_id": {
				Description: "ID of a Rollbar project",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": {
				Description: "Name of the token",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// Computed fields
			"access_token": {
				Description: "API token",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cur_rate_limit_window_count": {
				Description: "Number of API hits that occurred in the current rate limit window",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"cur_rate_limit_window_start": {
				Description: "Time when the current rate limit window began",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"date_created": {
				Description: "Date the token was created",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"date_modified": {
				Description: "Date the token was last modified",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"rate_limit_window_count": {
				Description: "Maximum allowed API hits during a rate limit window",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"rate_limit_window_size": {
				Description: "Duration of a rate limit window",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"scopes": {
				Description: `Project access scopes for the token.  Possible values are "read", "write", "post_server_item", or "post_client_item".`,
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Description: "Status of the token",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// dataSourceProjectAccessTokenRead reads a Rollbar project access token from
// the API
func dataSourceProjectAccessTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(int)
	var name string
	name, _ = d.Get("name").(string)
	l := log.With().
		Int("project_id", projectID).
		Str("name", name).
		Logger()
	l.Debug().Msg("Reading project access token from Rollbar")

	c := m.(map[string]*client.RollbarAPIClient)[schemaKeyToken]
	c.SetHeaderDataSource(rollbarProjectAccessToken)
	tokens, err := c.ListProjectAccessTokens(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Look for a token with matching name
	var found *client.ProjectAccessToken
	for i, t := range tokens {
		if t.Name == name {
			found = &tokens[i]
		}
	}

	// Error if no token matches.
	if found == nil {
		msg := fmt.Sprintf(`could not find access token with name matching %q`, name)
		l.Error().Msg(msg)
		return diag.FromErr(errors.New(msg))
	}

	// Write the values from API to Terraform state
	tokenMap := make(map[string]interface{})
	mustDecodeMapStructure(found, &tokenMap)
	for key, value := range tokenMap {
		mustSet(d, key, value)
	}

	// Set ID based on current time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	// Success
	return nil
}
