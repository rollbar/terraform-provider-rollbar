package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

// dataSourceProjectAccessTokens is a data source for listing all project access
// tokens belonging to a Rollbar project.
func dataSourceProjectAccessTokens() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectAccessTokensRead,

		Schema: map[string]*schema.Schema{
			// Required fields
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},

			// Computed fields
			"access_tokens": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Required: true,
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
						"name": {
							Type:     schema.TypeString,
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
				},
			},
		},
	}
}

// dataSourceProjectAccessTokensRead reads project access token data from Rollbar
func dataSourceProjectAccessTokensRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projId := d.Get("project_id").(int)
	l := log.With().
		Int("projId", projId).
		Logger()
	l.Debug().Msg("Reading project access token data from Rollbar")

	c := m.(*client.RollbarApiClient)
	tokens, err := c.ListProjectAccessTokens(projId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("access_tokens", tokens)
	if err != nil {
		log.Err(err).
			Interface("access_tokens", tokens).
			Msg("Error setting resource data")
		return diag.FromErr(err)
	}

	// Set resource ID to current timestamp (every resource must have an ID or
	// it will be destroyed).
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
