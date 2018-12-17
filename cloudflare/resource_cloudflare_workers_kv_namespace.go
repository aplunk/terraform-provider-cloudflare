package cloudflare

import (
	"context"
	"fmt"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceCloudflareWorkersKVNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudflareWorkersKVNamespaceCreate,
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

func resourceCloudflareWorkersKVNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	id := d.Get("id").(string)
	title := d.Get("title").(string)

	_, err := client.UpdateWorkersKVNamespace(context.Background(), id,
		&cloudflare.WorkersKVNamespaceRequest{Title: title},
	)

	if err != nil {
		return errors.Wrap(err, "error updating worker kv namespace")
	}

	return nil
}

func resourceCloudflareWorkersKVNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	id := d.Get("id").(string)

	_, err := client.DeleteWorkersKVNamespace(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP status 404") {
			return nil
		}

		return errors.Wrap(err, "error deleting workers kv namespace")
	}

	return nil
}

func resourceCloudflareWorkersKVNamespaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
