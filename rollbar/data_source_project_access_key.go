package rollbar

import (
	"fmt"

	"github.com/babbel/rollbar-go/rollbar"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceProjectAccessToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectAccessTokenRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed values
			"access_token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectAccessTokenRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(int)

	client := meta.(*rollbar.Client)
	accessToken, err := client.GetProjectAccessTokenByProjectIDAndName(projectID, name)
	if err != nil {
		return err
	}
	if accessToken == nil {
		d.SetId("")
		return fmt.Errorf("could not find project access token by project ID %d name %s", projectID, name)
	}

	id := fmt.Sprintf("%d-%s", accessToken.ProjectID, accessToken.Name)
	d.SetId(id)
	d.Set("access_token", accessToken.AccessToken)
	d.Set("status", accessToken.Status)
	d.Set("date_created", accessToken.DateCreated)

	return nil
}
