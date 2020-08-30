package rollbar

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var baseSchemaProject = map[string]*schema.Schema{
	"projects": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": &schema.Schema{
					Type:     schema.TypeInt,
					Computed: true,
				},
				"name": &schema.Schema{
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
}

func dataSourceSchemaProject() map[string]*schema.Schema {
	s := baseSchemaProject
	s["id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	return s
}
func resourceSchemaProject() map[string]*schema.Schema {
	s := baseSchemaProject
	s["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return s
}
