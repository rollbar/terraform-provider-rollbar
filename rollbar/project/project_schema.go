package project

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var baseFields = map[string]*schema.Schema{
	//"id": &schema.Schema{
	//	Type:     schema.TypeInt,
	//	Computed: true,
	//},
	//"name": &schema.Schema{
	//	Type:     schema.TypeString,
	//	Computed: true,
	//},
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

var baseSchema = map[string]*schema.Schema{
	"projects": {
		Type:     schema.TypeList,
		Computed: true,
	},
}

func dataSourceSchemaProject() map[string]*schema.Schema {
	s := baseSchema
	f := baseFields
	f["id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	f["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	s["projects"].Elem = &schema.Resource{
		Schema: f,
	}
	return s
}
func resourceSchemaProject() map[string]*schema.Schema {
	s := baseSchema
	s["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return s
}
