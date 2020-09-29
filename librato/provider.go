package librato

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/heroku/go-librato/librato"
)

// Provider returns a schema.Provider for Librato.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIBRATO_URL", nil),
				Description: "The librato API URL to use for all requests.",
			},

			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIBRATO_EMAIL", nil),
				Description: "The email address for the Librato account.",
			},

			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIBRATO_TOKEN", nil),
				Description: "The auth token for the Librato account.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"librato_space":       resourceLibratoSpace(),
			"librato_space_chart": resourceLibratoSpaceChart(),
			"librato_metric":      resourceLibratoMetric(),
			"librato_alert":       resourceLibratoAlert(),
			"librato_service":     resourceLibratoService(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	email := d.Get("email").(string)
	token := d.Get("token").(string)

	if u := d.Get("url").(string); u != "" {
		u, err := url.Parse(u)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse librato URL: %s", err)
		}
		return librato.NewClientWithBaseURL(u, email, token), nil
	}
	return librato.NewClient(email, token), nil
}
