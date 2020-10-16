package rollbar

import (
	"fmt"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rollbar/terraform-provider-rollbar/client"
	//"github.com/rs/zerolog/log"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed values
			"account_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date_created": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	c := meta.(*client.RollbarApiClient)
	pl, err := c.ListProjects()
	if err != nil {
		return err
	}

	var project client.Project
	var found bool
	for _, p := range pl {
		if p.Name == name {
			found = true
			project = p
		}
	}
	if !found {
		d.SetId("")
		return fmt.Errorf("no project with the name %s found", name)
	}

	id := fmt.Sprintf("%d", project.Id)
	d.SetId(id)
	err = d.Set("account_id", project.AccountId)
	if err != nil {
		return err
	}
	err = d.Set("date_created", project.DateCreated)
	if err != nil {
		return err
	}

	return nil
}
