package project

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// commonFields is the set of fields that will be included in both the data
// source schema and the resource schema
var commonFields = map[string]*schema.Schema{
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

// dataSourceFields is the set of fields that will be included only in the
// data source schema
var dataSourceFields = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

// resourceFields is the set of fields that will be included only in the
// resource schema.
var resourceFields = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
}

// constructSchema constructs a schema, combining common fields with fields
// specified in the argument.
func constructSchema(fields map[string]*schema.Schema) map[string]*schema.Schema {
	f := commonFields
	for k, v := range fields {
		f[k] = v
	}
	s := map[string]*schema.Schema{
		"projects": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: f,
			},
		},
	}
	return s
}

// dataSourceSchema returns the schema for a Terraform data source
func dataSourceSchema() map[string]*schema.Schema {
	return constructSchema(dataSourceFields)
}

// resourceSchema returns the schema for a Terraform resource
func resourceSchema() map[string]*schema.Schema {
	return constructSchema(resourceFields)
}
