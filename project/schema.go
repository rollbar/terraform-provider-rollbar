package project

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"projects": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"name": {
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

var resourceSchema = map[string]*schema.Schema{
	//"id": {
	//	Type:     schema.TypeInt,
	//	Required: true,
	//},
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
}
