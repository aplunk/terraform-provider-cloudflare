package cloudflare

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudflareDevicePolicyZoneCertificatesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"enable": {
			Type:     schema.TypeBool,
			Required: true,
		},
	}
}
