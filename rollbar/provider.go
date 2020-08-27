package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmcvetta/terraform-provider-rollbar/client"
)

const tokenSchemaKey = "token"

// Provider is a Terraform provider for Rollbar
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			tokenSchemaKey: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ROLLBAR_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			//"rollbar_project": resourceProject(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"rollbar_projects": dataSourceProjects(),,
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure sets up authentication in a Resty HTTP client.
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	token := d.Get(tokenSchemaKey).(string)
	c := client.NewClient(token)
	return c, diags
}
