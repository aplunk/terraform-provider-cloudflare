package cloudflare

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"crypto/sha256"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
	home "github.com/mitchellh/go-homedir"
)

func resourceCloudflareWorkersKV() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudflareWorkersKVCreate,
		Delete: resourceCloudflareWorkersKVDelete,
		Read:   resourceCloudflareWorkersKVRead,
		Update: resourceCloudflareWorkersKVCreate,
		Importer: &schema.ResourceImporter{
			State: resourceCloudflareWorkersKVImport,
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
				ForceNew: true,
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

	return setSha256(d, value)
}

func setSha256(d *schema.ResourceData, data []byte) error {
	hasher := sha256.New()
	return d.Set("sha256", string(hasher.Sum(data)))
}

func getValue(d *schema.ResourceData) ([]byte, error) {
	if source, ok := d.GetOk("source"); ok {
		fname, err := home.Expand(source.(string))
		if err != nil {
			return nil, err
		}

		return ioutil.ReadFile(fname)
	}

	if content, ok := d.GetOk("content"); ok {
		return []byte(content.(string)), nil
	}

	if contentB64, ok := d.GetOk("content_base64"); ok {
		return base64.StdEncoding.DecodeString(contentB64.(string))
	}

	return nil, fmt.Errorf("source, content, or content_base64 must be set")
}

func resourceCloudflareWorkersKVCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	namespaceID := d.Get("namespace_id").(string)
	key := d.Get("key").(string)
	value, err := getValue(d)
	if err != nil {
		return err
	}

	if err = setSha256(d, value); err != nil {
		return err
	}

	_, err = client.CreateWorkersKV(context.Background(), namespaceID, key, value)
	return err
}

func resourceCloudflareWorkersKVDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	namespaceID := d.Get("namespace_id").(string)
	key := d.Get("key").(string)

	_, err := client.DeleteWorkersKV(context.Background(), namespaceID, key)
	return err
}

func resourceCloudflareWorkersKVImport(d *schema.ResourceData, meta interface{}) (result []*schema.ResourceData, err error) {
	client := meta.(*cloudflare.API)

	namespaceID := d.Get("namespace_id").(string)
	key := d.Get("key").(string)

	data, err := client.ReadWorkersKV(context.Background(), namespaceID, key)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, setSha256(d, data)
}
