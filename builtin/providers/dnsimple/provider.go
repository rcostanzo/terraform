package dnsimple

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("DNSIMPLE_EMAIL"),
				Description: "A registered DNSimple email address.",
			},

			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFunc("DNSIMPLE_TOKEN"),
				Description: "The token key for API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dnsimple_record": resourceDNSimpleRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return nil, nil
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Email: d.Get("email").(string),
		Token: d.Get("token").(string),
	}

	return config.Client()
}
