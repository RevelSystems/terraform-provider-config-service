package revel_config_service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConfigServiceRequest struct {
	ClientName string                 `json:"ClientName"`
	Attributes map[string]interface{} `json:"Attributes"`
}

func resourceConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigCreate,
		ReadContext:   resourceConfigRead,
		UpdateContext: resourceConfigUpdate,
		DeleteContext: resourceConfigDelete,
		Schema: map[string]*schema.Schema{
			"client": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attributes_json": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config, ok := m.(Ctx)
	if !ok {
		return diag.FromErr(fmt.Errorf("can't read config"))
	}

	client := &http.Client{Timeout: time.Duration(config.Timeout) * time.Second}
	clientName := d.Get("client").(string)

	attributes_json := d.Get("attributes_json").(string)
	var new_json map[string]interface{}
	err := json.Unmarshal([]byte(attributes_json), &new_json)
	if err != nil {
		return diag.FromErr(err)
	}

	requestBody, err := json.Marshal(ConfigServiceRequest{ClientName: clientName, Attributes: new_json})
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequest("POST", config.BaseUrl, strings.NewReader(string(requestBody)))
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Authorization", config.Token)
	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	if r.StatusCode == http.StatusCreated {
		d.SetId(clientName)
		return diags
	}

	return diag.FromErr(fmt.Errorf("resourceConfigCreate: unexpected status code returned %v", r.StatusCode))
}

func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config, ok := m.(Ctx)
	if !ok {
		return diag.FromErr(fmt.Errorf("can't read config"))
	}

	client := &http.Client{Timeout: time.Duration(config.Timeout) * time.Second}
	clientName := d.Id()

	attributes_json := d.Get("attributes_json").(string)
	var new_json map[string]interface{}
	err := json.Unmarshal([]byte(attributes_json), &new_json)
	if err != nil {
		return diag.FromErr(err)
	}

	requestBody, err := json.Marshal(ConfigServiceRequest{ClientName: clientName, Attributes: new_json})
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/%s", config.BaseUrl, clientName), strings.NewReader(string(requestBody)),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Authorization", config.Token)
	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	if r.StatusCode == http.StatusOK {
		return diags
	}

	return diag.FromErr(fmt.Errorf("resourceConfigUpdate: unexpected status code returned %v", r.StatusCode))
}

func resourceConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config, ok := m.(Ctx)
	if !ok {
		return diag.FromErr(fmt.Errorf("can't read config"))
	}

	client := &http.Client{Timeout: time.Duration(config.Timeout) * time.Second}
	clientName := d.Id()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", config.BaseUrl, clientName), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Authorization", config.Token)
	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	if r.StatusCode == http.StatusNoContent {
		return diags
	}

	return diag.FromErr(fmt.Errorf("resourceConfigDelete: unexpected status code returned %v", r.StatusCode))
}
