package cloudflare

import (
	"context"

	"crypto/sha256"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudflareWorkersKV() *schema.Resource {
	return &schema.Resource{
		Create: nil,
		Update: nil,
		Delete: nil,
		Read:   resourceCloudflareWorkersKVRead,
		Importer: &schema.ResourceImporter{
			State: nil,
		},
		Schema: map[string]*schema.Schema{
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source", "content_base64"},
			},
			"content_base64": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source", "content"},
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sha256": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content", "content_base64"},
			},
		},
	}
}

func resourceCloudflareWorkersKVRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	namespaceID := d.Get("namespaceID").(string)
	key := d.Get("key").(string)

	value, err := client.ReadWorkersKV(context.Background(), namespaceID, key)
	if err != nil {
		return err
	}

	hasher := sha256.New()
	return d.Set("sha256", string(hasher.Sum(value)))
}
