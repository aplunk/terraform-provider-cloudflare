package cloudflare

import (
	"context"
	"fmt"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

const (
	allKey   = "all"
	titleKey = "title"
)

func resourceCloudflareWorkersKVNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudflareWorkersKVNamespaceCreate,
		Update: resourceCloudflareWorkersKVNamespaceUpdate,
		Delete: resourceCloudflareWorkersKVNamespaceDelete,
		Read:   resourceCloudflareWorkersKVNamespaceRead,
		Importer: &schema.ResourceImporter{
			State: resourceCloudflareWorkersKVNamespaceImport,
		},
		Schema: map[string]*schema.Schema{
			titleKey: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCloudflareWorkersKVNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	title := d.Get(titleKey).(string)
	if title == "" {
		return fmt.Errorf("missing required title field")
	}

	resp, err := client.CreateWorkersKVNamespace(context.Background(),
		&cloudflare.WorkersKVNamespaceRequest{Title: title},
	)

	if err != nil {
		return err
	}

	d.SetId(resp.Result.ID)
	return nil
}

func resourceCloudflareWorkersKVNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	title := d.Get(titleKey).(string)

	_, err := client.UpdateWorkersKVNamespace(context.Background(), d.Id(),
		&cloudflare.WorkersKVNamespaceRequest{Title: title},
	)

	d.Set(titleKey, title)

	if err != nil {
		return errors.Wrap(err, "error updating worker kv namespace")
	}

	return nil
}

func resourceCloudflareWorkersKVNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	_, err := client.DeleteWorkersKVNamespace(context.Background(), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "HTTP status 404") {
			return nil
		}

		return errors.Wrap(err, "error deleting workers kv namespace")
	}

	return nil
}

func makeNamespaceFilter(d *schema.ResourceData) func(string) bool {
	id := d.Id()
	switch strings.ToLower(id) {
	case allKey:
		return func(_ string) bool { return true }
	default:
		return func(nxtID string) bool { return id == nxtID }
	}
}

func resourceCloudflareWorkersKVNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	resp, err := client.ListWorkersKVNamespaces(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "HTTP status 404") {
			return nil
		}

		return errors.Wrap(err, "error deleting workers kv namespace")
	}

	filter := makeNamespaceFilter(d)

	for _, namespace := range resp.Result {
		if !filter(namespace.ID) {
			continue
		}
		data := &schema.ResourceData{}
		data.SetId(namespace.ID)
		data.Set(titleKey, namespace.Title)
		break
	}
	return nil
}

func resourceCloudflareWorkersKVNamespaceImport(d *schema.ResourceData, meta interface{}) (result []*schema.ResourceData, err error) {
	client := meta.(*cloudflare.API)

	resp, err := client.ListWorkersKVNamespaces(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "error importing workers kv namespaces")
	}

	filter := makeNamespaceFilter(d)

	for _, namespace := range resp.Result {
		if !filter(namespace.ID) {
			continue
		}
		data := &schema.ResourceData{}
		data.SetId(namespace.ID)
		data.Set(titleKey, namespace.Title)
		data.SetType("cloudflare_workers_kv_namespace")
		result = append(result, data)
	}

	return result, err
}
