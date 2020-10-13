package rollbar

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

// dataSourceProjectAccessToken is a data source returning a named access token
// belonging to a Rollbar project.
func dataSourceProjectAccessToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectAccessTokenRead,

		Schema: map[string]*schema.Schema{
			// Required fields
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed fields
			"access_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cur_rate_limit_window_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cur_rate_limit_window_start": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rate_limit_window_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rate_limit_window_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceProjectAccessTokenRead reads a Rollbar project access token from
// the API
func dataSourceProjectAccessTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projId := d.Get("project_id").(int)
	var name string
	name, _ = d.Get("name").(string)
	l := log.With().
		Int("project_id", projId).
		Str("name", name).
		Logger()
	l.Debug().Msg("Reading project access token from Rollbar")

	c := m.(*client.RollbarApiClient)
	tokens, err := c.ListProjectAccessTokens(projId)
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
	// FIXME: Is an error appropriate for this situation?  Is there some better
	//  way to say we didn't find anything that matches, but there was no
	//  internal error?
	if found == nil {
		msg := "could not find access token with name matching name"
		l.Error().Msg(msg)
		return diag.FromErr(fmt.Errorf(msg))
	}

	// Write the values from API to Terraform state
	tokenMap := make(map[string]interface{})
	err = mapstructure.Decode(found, &tokenMap)
	if err != nil {
		return diag.FromErr(err)
	}
	for key, value := range tokenMap {
		err = d.Set(key, value)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Set ID based on current time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	// Success
	return nil
}