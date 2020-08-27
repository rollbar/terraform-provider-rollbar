package rollbar

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmcvetta/terraform-provider-rollbar/rollbar/client"
	"gopkg.in/jeevatkm/go-model.v1"
	"strconv"
	"time"
)

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectsRead,
		Schema: map[string]*schema.Schema{
			"projects": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_created": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"date_modified": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.RollbarApiClient)

	lp, err := c.ListProjects()
	if err != nil {
		return diag.FromErr(err)
	}

	projects, err := model.Map(lp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("projects", projects); err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID to current timestamp (every resource must have an ID or
	// it will be destroyed).
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
