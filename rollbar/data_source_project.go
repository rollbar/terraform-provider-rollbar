/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmcvetta/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"gopkg.in/jeevatkm/go-model.v1"
	"strconv"
	"time"
)

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectsRead,
		Schema: map[string]*schema.Schema{
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
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

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Debug().Msg("Reading projects data from API")
	var diags diag.Diagnostics
	c := m.(*client.RollbarApiClient)

	lp, err := c.ListProjects()
	if err != nil {
		return diag.FromErr(err)
	}

	projects := make([]map[string]interface{}, 0)
	for _, v := range lp {
		m, err := model.Map(v)
		if err != nil {
			log.Err(err).
				Interface("lp", lp).
				Msg("Error converting to map")
			return diag.FromErr(err)
		}
		projects = append(projects, m)
	}

	if err := d.Set("projects", projects); err != nil {
		log.Err(err).
			Interface("projects", projects).
			Msg("Error setting resource data")
		return diag.FromErr(err)
	}

	// Set resource ID to current timestamp (every resource must have an ID or
	// it will be destroyed).
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	log.Warn().Msg("Successfully read project list from API.")

	return diags
}
