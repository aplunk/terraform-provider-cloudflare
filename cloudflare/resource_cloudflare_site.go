package cloudflare

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
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
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"chunk_size": {
				Type:    schema.TypeInt,
				Default: 1024,
			},
		},
	}
}

type uploader func(key string, value []byte) error

func uploadFile(pathStr, prefix string, info os.FileInfo, split int, uploadKV uploader) (keys []string, err error) {
	fh, err := os.Open(pathStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer fh.Close()

	vSize := int(info.Size())
	if vSize == 0 {
		return nil, errors.Errorf("refusing to upload empty file:%s", info.Name())
	}

	for i := 0; i < vSize; i += split {
		data := make([]byte, vSize)
		read, err := fh.ReadAt(data, int64(i))
		if err != nil && err != io.EOF {
			return nil, errors.WithStack(err)
		}

		key := prefix
		if i != 0 {
			key = fmt.Sprintf("%s_%d", prefix, int(vSize/i))
		}

		if err := uploadKV(key, data[:read]); err != nil {
			return nil, errors.WithStack(err)
		}
		keys = append(keys, key)

	}

	return keys, nil
}

func uploadSite(namespaceID, source string, limit int, uploadKV uploader) (map[string][]string, error) {
	largeFiles := make(map[string][]string)
	return largeFiles, filepath.Walk(source, func(pathStr string, info os.FileInfo, err error) error {
		// fail early if an error is passed in
		if err != nil {
			return errors.WithStack(err)
		}

		// unable to upload directories
		if info.IsDir() {
			return nil
		}

		// normalize the file key
		key := strings.Replace(pathStr, string(filepath.Separator), "_", -1)

		// upload large files in chunks returning a mapping of the chunks which will
		// become a manifest enabling reconstructing the original file
		if info.Size() > int64(limit) {
			chunks, err := uploadFile(pathStr, key, info, limit, uploadKV)
			if err != nil {
				return errors.WithStack(err)
			}
			largeFiles[key] = chunks
			return nil
		}

		// files smaller than the limit can be uploaded without returning a mapping
		_, err = uploadFile(pathStr, key, info, limit, uploadKV)
		return errors.WithStack(err)
	})
}

func resourceCloudflareSiteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	source := d.Get("source").(string)
	namespaceID := d.Get("namespace_id").(string)
	chunkSize := d.Get("chunk_size").(int)

	uploader := func(key string, value []byte) error {
		_, err := client.CreateWorkersKV(context.Background(), namespaceID, key, value)
		return err
	}

	manifest, err := uploadSite(namespaceID, source, chunkSize, uploader)
	_ = manifest
	return err
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
