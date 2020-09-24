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
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		DeleteContext: resourceProjectDelete,

		// Projects cannot be updated via API
		//UpdateContext: resourceProjectUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	log.Debug().Str("name", name).
		Msg("Creating new project")

	c := m.(*client.RollbarApiClient)
	p, err := c.CreateProject(name)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Debug().Interface("project", p).Msg("CreateProject() result")

	d.SetId(strconv.Itoa(p.Id))

	readDiags := resourceProjectRead(ctx, d, m)
	diags = append(diags, readDiags...)
	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	mProj, err := model.Map(proj)
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

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
