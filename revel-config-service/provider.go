package revel_config_service

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Ctx struct {
	Token   string
	BaseUrl string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CONFIG_SERVICE_TOKEN", nil),
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONFIG_SERVICE_BASE_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"configuration": resourceConfig(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	url := d.Get("base_url").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return Ctx{
		Token:   token,
		BaseUrl: url,
	}, diags
}
