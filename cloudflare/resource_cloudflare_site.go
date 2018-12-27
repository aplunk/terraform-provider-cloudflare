package cloudflare

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudflareWorkersKVSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudflareSiteCreate,
		Delete: resourceCloudflareSiteDelete,
		Read:   resourceCloudflareSiteRead,
		Update: resourceCloudflareSiteUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceCloudflareSiteImport,
		},
		CustomizeDiff: updateComputed,
		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func uploadLargeFile(pathStr, prefix string, info os.FileInfo, split int, uploader func(string, []byte) error) (chunks []string, err error) {
	fh, err := os.Open(path.Join(pathStr, info.Name()))
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	vSize := int(info.Size())
	for i := 0; i < vSize; i += split {
		data := make([]byte, vSize)
		read, err := fh.ReadAt(data, int64(i))
		if err != nil {
			return nil, err
		}

		key := fmt.Sprintf("%s_%d", prefix, i%vSize)

		if err := uploader(key, data[:read]); err != nil {
			return nil, err
		}
		chunks = append(chunks, key)
	}

	return chunks, nil
}

func resourceCloudflareSiteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	source := d.Get("source").(string)
	namespaceID := d.Get("namespace_id").(string)

	largeFiles := make(map[string][]string)

	err := filepath.Walk(source, func(pathStr string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// normalize the file key
		key := fmt.Sprintf("%s_%s", strings.Replace(pathStr, string(filepath.Separator), "_", -1), info.Name())

		split := 1024
		uploader := func(key string, value []byte) error {
			_, err := client.CreateWorkersKV(context.Background(), namespaceID, key, value)
			return err
		}

		if info.Size() > int64(split) {
			keys, err := uploadLargeFile(pathStr, key, info, split, uploader)
			if err != nil {
				return err
			}
			largeFiles[key] = keys
			return nil
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceCloudflareSiteDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudflareSiteRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudflareSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudflareSiteImport(d *schema.ResourceData, meta interface{}) (result []*schema.ResourceData, err error) {
	return nil, nil
}
