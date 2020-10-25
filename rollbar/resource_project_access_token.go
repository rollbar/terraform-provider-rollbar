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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/rollbar/terraform-provider-rollbar/client"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func resourceProjectAccessToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectAccessTokenCreate,
		ReadContext:   resourceProjectAccessTokenRead,
		DeleteContext: resourceProjectAccessTokenDelete,
		UpdateContext: resourceProjectAccessTokenUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceProjectAccessTokenImporter,
		},

		Schema: map[string]*schema.Schema{
			// Required fields
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // FIXME: https://github.com/rollbar/terraform-provider-rollbar/issues/41
			},
			"scopes": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true, // FIXME: https://github.com/rollbar/terraform-provider-rollbar/issues/41
			},

			// Optional fields
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enabled",
				ForceNew: true, // FIXME: https://github.com/rollbar/terraform-provider-rollbar/issues/41
			},
			"rate_limit_window_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"rate_limit_window_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			// Computed fields
			"access_token": {
				Type:     schema.TypeString,
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
			"cur_rate_limit_window_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cur_rate_limit_window_start": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceProjectAccessTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectId := d.Get("project_id").(int)
	name := d.Get("name").(string)
	scopesInterface := d.Get("scopes").([]interface{})
	var scopes []client.Scope
	for _, v := range scopesInterface {
		s := v.(string)
		scopes = append(scopes, client.Scope(s))
	}
	status := client.Status(d.Get("status").(string))
	size := d.Get("rate_limit_window_size").(int)
	count := d.Get("rate_limit_window_count").(int)
	l := log.With().
		Int("project_id", projectId).
		Str("name", name).
		Int("rate_limit_window_size", size).
		Int("rate_limit_window_count", count).
		Interface("scopes", scopes).
		Interface("status", status).
		Logger()
	l.Debug().Msg("Creating new project access token")

	c := m.(*client.RollbarApiClient)
	pat, err := c.CreateProjectAccessToken(client.ProjectAccessTokenCreateArgs{
		Name:                 name,
		ProjectID:            projectId,
		Scopes:               scopes,
		Status:               status,
		RateLimitWindowSize:  size,
		RateLimitWindowCount: count,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pat.AccessToken)

	return resourceProjectAccessTokenRead(ctx, d, m)
}

func resourceProjectAccessTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accessToken := d.Id()
	projectId := d.Get("project_id").(int)
	l := log.With().
		Str("accessToken", accessToken).
		Logger()
	l.Debug().Msg("Reading resource project access token")

	c := m.(*client.RollbarApiClient)
	pat, err := c.ReadProjectAccessToken(projectId, accessToken)
	if err == client.ErrNotFound {
		return handleErrNotFound(d, "project access token")
	}
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
	accessToken := d.Id()
	projectID := d.Get("project_id").(int)
	size := d.Get("rate_limit_window_size").(int)
	count := d.Get("rate_limit_window_count").(int)
	args := client.ProjectAccessTokenUpdateArgs{
		ProjectID:            projectID,
		AccessToken:          accessToken,
		RateLimitWindowSize:  size,
		RateLimitWindowCount: count,
	}
	l := log.With().Interface("args", args).Logger()
	l.Debug().Msg("Updating resource project access token")
	c := m.(*client.RollbarApiClient)
	err := c.UpdateProjectAccessToken(args)
	if err != nil {
		log.Err(err).Send()
		return diag.FromErr(err)
	}
	diags := resourceProjectAccessTokenRead(ctx, d, m)
	return diags
}

func resourceProjectAccessTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	accessToken := d.Id()
	projectId := d.Get("project_id").(int)

	l := log.With().
		Int("projectId", projectId).
		Str("accessToken", accessToken).
		Logger()
	l.Debug().Msg("Deleting resource project access token")

	c := m.(*client.RollbarApiClient)
	err := c.DeleteProjectAccessToken(projectId, accessToken)
	if err != nil {
		return diag.FromErr(err)
	}

	// FIXME: Implement this functionality when the API is ready!
	//  https://github.com/rollbar/terraform-provider-rollbar/issues/12
	diags = diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "Delete not implemented for resource `rollbar_project_access_token`",
		Detail:   "https://github.com/rollbar/terraform-provider-rollbar/issues/12",
	}}
	return diags
}

func resourceProjectAccessTokenImporter(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	l := log.With().Str("id", d.Id()).Logger()
	l.Debug().Msg("Importing resource rollbar project access token")
	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%q), expected PROJECT-ID/ACCESS-TOKEN", d.Id())
	}
	projectIdString := idParts[0]
	accessToken := idParts[1]
	projectId, err := strconv.Atoi(projectIdString)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}
	l.Debug().
		Int("project_id", projectId).
		Str("access_token", accessToken).
		Send()
	err = d.Set("project_id", projectId)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}
	d.SetId(accessToken)
	return []*schema.ResourceData{d}, nil
}
