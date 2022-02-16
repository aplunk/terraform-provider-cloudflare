package cloudflare

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudflareDevicePolicyZoneCertificates() *schema.Resource {
	return &schema.Resource{
		Schema: resourceCloudflareDevicePolicyZoneCertificatesSchema(),
		Create: resourceCloudflareDevicePolicyZoneCertificateUpdate,
		Read:   resourceCloudflareDevicePolicyZoneCertificateRead,
		Update: resourceCloudflareDevicePolicyZoneCertificateUpdate,
		Delete: resourceCloudflareDevicePolicyZoneCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudflareDevicePolicyZoneCertificateImport,
		},
	}
}

func resourceCloudflareDevicePolicyZoneCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)
	enabled := d.Get("enable").(bool)

	log.Printf("[DEBUG] Updating Cloudflare device policy zone certificate: zone:%s enabled:%t", zoneID, enabled)

	_, err := client.UpdateDeviceClientCertificatesZone(context.Background(), zoneID, enabled)
	if err != nil {
		return fmt.Errorf("error updating Cloudflare device policy zone certificate %q: %s", zoneID, err)
	}

	d.SetId(zoneID)
	return nil
}

func resourceCloudflareDevicePolicyZoneCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	zoneID := d.Get("zone_id").(string)

	enabled, err := client.GetDeviceClientCertificatesZone(context.Background(), zoneID)
	if err != nil {
		return fmt.Errorf("error reading device policy zone certificate setting %q: %s", zoneID, err)
	}

	d.SetId(zoneID)
	d.Set("enable", enabled.Result.Enable)
	return nil
}

func resourceCloudflareDevicePolicyZoneCertificateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	zoneID := d.Id()

	log.Printf("[DEBUG] Importing Cloudflare device policy zone certificate setting: zone %s", zoneID)

	d.SetId(zoneID)
	resourceCloudflareDevicePostureRuleRead(d, meta)
	return []*schema.ResourceData{d}, nil
}

func resourceCloudflareDevicePolicyZoneCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
