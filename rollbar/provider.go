package rollbar

import (
	"github.com/babbel/rollbar-go/rollbar"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ROLLBAR_APIKEY", nil),
				Description: "API Key for accessing the rollbar api.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rollbar_user": resourceUser(),
		},

		ConfigureFunc: configureProvider,
	}

}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	apiKey := d.Get("api_key").(string)
	client, err := rollbar.NewClient(apiKey)

	if err != nil {
		return nil, err
	}

	return client, nil
}
