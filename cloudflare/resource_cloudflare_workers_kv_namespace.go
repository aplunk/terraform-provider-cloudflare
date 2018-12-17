package cloudflare

import (
	"context"
	"fmt"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudflareWorkersKVNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudflareWorkersKVNamespaceCreate,
		Read:   resourceCloudflareWorkersKVNamespaceRead,
		Update: resourceCloudflareWorkersKVNamespaceUpdate,
		Delete: resourceCloudflareWorkersKVNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudflareWorkersKVNamespaceImport,
		},
		Schema: map[string]*schema.Schema{
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCloudflareWorkersKVNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	title := d.Get("title").(string)
	if title == "" {
		return fmt.Errorf("missing required title field")
	}

	resp, err := client.CreateWorkersKVNamespace(context.Background(),
		&cloudflare.WorkersKVNamespaceRequest{Title: title},
	)

	if err != nil {
		return err
	}

	d.SetId(resp.Result.Title)
	return nil
}

func resourceCloudflareWorkersKVNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudflareWorkersKVNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudflareWorkersKVNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudflareWorkersKVNamespaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
