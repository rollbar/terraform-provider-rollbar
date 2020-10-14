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
	"github.com/mitchellh/mapstructure"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
)

func resourceProjectAccessToken() *schema.Resource {
	log.Fatal().Msg("Not yet implemented")
	return &schema.Resource{
		CreateContext: resourceProjectAccessTokenCreate,
		ReadContext:   resourceProjectAccessTokenRead,
		DeleteContext: resourceProjectAccessTokenDelete,
		UpdateContext: resourceProjectAccessTokenUpdate,

		Schema: map[string]*schema.Schema{
			// Required fields
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Optional fields
			"rate_limit_window_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rate_limit_window_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// FIXME: Is status field optional or computed?
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

func resourceProjectAccessTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	projectId := d.Get("project_id").(int)
	name := d.Get("name").(string)
	scopeStrings := d.Get("name").([]string)
	var scopes []client.Scope
	for _, st := range scopeStrings {
		sc := client.Scope(st)
		scopes = append(scopes, sc)
	}
	l := log.With().
		Int("project_id", projectId).
		Str("name", name).
		Interface("scopes", scopes).
		Logger()
	l.Debug().Msg("Creating new project access token")

	c := m.(*client.RollbarApiClient)
	t, err := c.CreateProjectAccessToken(client.ProjectAccessTokenArgs{
		Name:      name,
		ProjectID: projectId,
		Scopes:    scopes,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(t.AccessToken)

	readDiags := resourceProjectAccessTokenRead(ctx, d, m)
	diags = append(diags, readDiags...)
	return diags
}

func resourceProjectAccessTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Fatal().Msg("Not yet implemented")
	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Err(err).Msg("Error converting Id to integer")
		return diag.FromErr(err)
	}

	l := log.With().
		Int("id", id).
		Logger()
	l.Debug().Msg("Reading project resource")

	c := m.(*client.RollbarApiClient)
	proj, err := c.ReadProject(id)
	if err != nil {
		return diag.FromErr(err)
	}
	var mProj map[string]interface{}
	err = mapstructure.Decode(proj, &mProj)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	for k, v := range mProj {
		if k == "id" {
			continue
		}
		err = d.Set(k, v)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	}
	d.SetId(strconv.Itoa(proj.Id))

	return diags
}

func resourceProjectAccessTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Fatal().Msg("Not yet implemented")
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectAccessTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Fatal().Msg("Not yet implemented")
	var diags diag.Diagnostics

	projectId, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Err(err).Msg("Error converting Id to integer")
		return diag.FromErr(err)
	}

	l := log.With().
		Int("projectId", projectId).
		Logger()
	l.Debug().Msg("Deleting project")

	c := m.(*client.RollbarApiClient)
	err = c.DeleteProject(projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
