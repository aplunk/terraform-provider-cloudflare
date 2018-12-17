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

func resourceCloudflareWorkersKVNamespaceImport(d *schema.ResourceData, meta interface{}) (result []*schema.ResourceData, err error) {
	client := meta.(*cloudflare.API)

	resp, err := client.ListWorkersKVNamespaces(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "error importing workers kv namespaces")
	}

	id := d.Id()
	var filter func(id string) bool
	switch strings.ToLower(id) {
	case "all":
		filter = func(_ string) bool { return true }
	default:
		filter = func(nxtID string) bool { return id == nxtID }
	}

	for _, namespace := range resp.Result {
		if !filter(namespace.ID) {
			continue
		}
		data := &schema.ResourceData{}
		data.SetId(namespace.ID)
		data.Set("title", namespace.Title)
		result = append(result, data)
	}

	return result, err
}
