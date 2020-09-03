/*
 * Copyright (c) 2020 Jason McVetta <jmcvetta@protonmail.com>, all rights
 * reserved.
 *
 * NO LICENSE WHATSOEVER IS GRANTED for this software without written contract
 * between author and licensee.
 */

/*
 * Package rollbar implements a Terraform provider for the Rollbar API.
 */
package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmcvetta/terraform-provider-rollbar/client"
)

const (
	schemaKeyToken = "token"
	schemaKeyUrl   = "api_url"
)

// Provider is a Terraform provider for Rollbar
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			schemaKeyToken: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ROLLBAR_TOKEN", nil),
			},
			schemaKeyUrl: {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc(
					"ROLLBAR_API_URL",
					"https://api.rollbar.com",
				),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rollbar_project": Resource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"rollbar_projects": DataSource(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure sets up authentication in a Resty HTTP client.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	token := d.Get(schemaKeyToken).(string)
	u := d.Get(schemaKeyUrl).(string)

	c, err := client.NewClient(u, token)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
