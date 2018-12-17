package cloudflare

import "github.com/hashicorp/terraform/helper/schema"

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
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCloudflareWorkersKVNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
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
