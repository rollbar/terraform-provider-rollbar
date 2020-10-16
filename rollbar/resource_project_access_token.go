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

package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
)

func resourceProjectAccessToken() *schema.Resource {
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
			"status": {
				Type:     schema.TypeString,
				Required: true,
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
	scopesInterface := d.Get("scopes").([]interface{})
	var scopes []client.Scope
	for _, v := range scopesInterface {
		s := v.(string)
		scopes = append(scopes, client.Scope(s))
	}
	status := client.Status(d.Get("status").(string))
	l := log.With().
		Int("project_id", projectId).
		Str("name", name).
		Interface("scopes", scopes).
		Interface("status", status).
		Logger()
	l.Debug().Msg("Creating new project access token")

	c := m.(*client.RollbarApiClient)
	pat, err := c.CreateProjectAccessToken(client.ProjectAccessTokenCreateArgs{
		Name:      name,
		ProjectID: projectId,
		Scopes:    scopes,
		Status:    status,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pat.AccessToken)

	readDiags := resourceProjectAccessTokenRead(ctx, d, m)
	diags = append(diags, readDiags...)
	return diags
}

func resourceProjectAccessTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accessToken := d.Id()
	projectId := d.Get("project_id").(int)
	l := log.With().
		Str("accessToken", accessToken).
		Logger()
	l.Debug().Msg("Reading project resource")

	c := m.(*client.RollbarApiClient)
	pat, err := c.ReadProjectAccessToken(projectId, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}
	var mPat map[string]interface{}
	err = mapstructure.Decode(pat, &mPat)
	if err != nil {
		l.Err(err).Send()
		return diag.FromErr(err)
	}
	for k, v := range mPat {
		err = d.Set(k, v)
		if err != nil {
			l.Err(err).Send()
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceProjectAccessTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Fatal().Msg("Not yet implemented")
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectAccessTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accessToken := d.Id()
	projectId := d.Get("project_id").(int)

	l := log.With().
		Int("projectId", projectId).
		Str("accessToken", accessToken).
		Logger()
	l.Debug().Msg("Deleting project")

	c := m.(*client.RollbarApiClient)
	err := c.DeleteProjectAccessToken(projectId, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
